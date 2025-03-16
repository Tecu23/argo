// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PawnsEvaluation returns the evaluation for pawns and pawn structure
func (e *Evaluator) PawnsEvaluation(b *board.Board) (mg, eg int) {
	pawnsBB := b.Bitboards[WP]
	for pawnsBB != 0 {
		sq := pawnsBB.FirstOne()

		if doubleIsolated(b, sq) {
			mg -= 11
			eg -= 56
		} else if isolated(b, sq) {
			mg -= 5
			eg -= 15
		} else if backward(b, sq) {
			mg -= 9
			eg -= 24
		}

		if doubled(b, sq) {
			mg -= 11
			eg -= 56
		}

		if connected(b, sq) {
			connBonus := connectedBonus(b, sq)
			mg += connBonus
			eg += connBonus * (8 - (sq / 8) - 3) / 4
		}

		weakUnopposedBonus := weakUnopposedPawn(b, sq)
		mg -= 13 * weakUnopposedBonus
		eg -= 27 * weakUnopposedBonus

		eg -= 56 * weakLever(b, sq)

		blockedIdx := blocked(b, sq)
		mg += []int{0, -11, -3}[blockedIdx]
		eg += []int{0, -4, 4}[blockedIdx]

	}

	return mg, eg
}

// doubleIsolated is a penalty if a double pawn is stopped only
// by a single opponent pawn on the same file.
func doubleIsolated(b *board.Board, sq int) bool {
	// Should return if square doesn't contain a white pawn
	// But because we only call double isolated on white pawns this is not needed

	file := sq % 8
	rank := sq / 8

	// Return early is the pawn is not isolated
	if !isolated(b, sq) {
		return false
	}

	leftFile := file > 0
	rightFile := file < 7

	// Create maks for adjancent files
	var leftFileMask, rightFileMask bitboard.Bitboard
	if leftFile {
		leftFileMask = FileMasks[file-1]
	}

	if rightFile {
		rightFileMask = FileMasks[file+1]
	}

	// Get file mask for current file
	fileMask := FileMasks[file]

	// Check for another white pawn behind this one (greater rank values on same file)
	var behindRanksMask bitboard.Bitboard
	for r := rank + 1; r < 8; r++ {
		behindRanksMask |= RankMasks[r]
	}

	// If no pawn behind, then the pawn is not doubled
	whitePawnsBehind := b.Bitboards[WP] & fileMask & behindRanksMask
	if whitePawnsBehind == 0 {
		return false
	}

	// Create a mask for ranks in front of our pawn
	var frontRankMask bitboard.Bitboard
	for r := 0; r < rank; r++ {
		frontRankMask |= RankMasks[r]
	}

	// Check for black pawns in front on the same file
	blackPawnsOnFileInFront := b.Bitboards[BP] & fileMask & frontRankMask
	if blackPawnsOnFileInFront == 0 {
		return false
	}

	// Check for enemy pawns on adjacent files (any rank)
	blackPawnsOnAdjacentFiles := (leftFile && (b.Bitboards[BP]&leftFileMask != 0)) ||
		(rightFile && (b.Bitboards[BP]&rightFileMask != 0))

	// If there's a doubled white pawn, black pawns in front on the same file,
	// and no black pawns on adjacent files, return true
	return !blackPawnsOnAdjacentFiles
}

// Isolated checks if pawn is isolated. In chess, an isolated pawn is pawn
// which has no friendly pawn on an adjacent files
func isolated(b *board.Board, sq int) bool {
	file := sq % 8

	leftFile := file > 0
	rightFile := file < 7

	// Create maks for adjancent files
	var leftFileMask, rightFileMask bitboard.Bitboard
	if leftFile {
		leftFileMask = FileMasks[file-1]
	}

	if rightFile {
		rightFileMask = FileMasks[file+1]
	}

	// Check if pawn is isolated
	adjancentFilePawns := (leftFile && (b.Bitboards[WP]&leftFileMask) != 0) ||
		(rightFile && (b.Bitboards[WP]&rightFileMask) != 0)

	if adjancentFilePawns {
		return false
	}

	return true
}

func optimizedBackward(b *board.Board, sq int) bool {
	// Should return if square doesn't contain a white pawn
	// But because we only call double isolated on white pawns this is not needed

	rank := sq / 8
	file := sq % 8

	// Check if there are friendly pawns on adjanced files at or ahead of this pawn
	leftFile := file > 0
	rightFile := file < 7

	// Create masks for adjacent files
	var leftFileMask, rightFileMask bitboard.Bitboard
	if leftFile {
		leftFileMask = FileMasks[file-1]
	}

	if rightFile {
		rightFileMask = FileMasks[file+1]
	}

	// Create mask for ranks at or ahead of our pawn
	var forwardRanksMask bitboard.Bitboard
	for r := rank; r < 8; r++ {
		forwardRanksMask |= RankMasks[r]
	}

	// Check for white pawns on adjacent files at or ahead of this pawn
	whitePawnsAheadLeft := leftFile && (b.Bitboards[WP]&leftFileMask&forwardRanksMask != 0)
	whitePawnsAheadRight := rightFile && (b.Bitboards[WP]&rightFileMask&forwardRanksMask != 0)

	if whitePawnsAheadLeft || whitePawnsAheadRight {
		return false
	}

	// Check if the pawn is blocked or can be captured when advanced

	// Position one square ahead
	squareAhead := (rank-1)*8 + file
	if rank > 0 && b.Bitboards[BP].Test(squareAhead) {
		return true
	}

	// Positions that could attack the square ahead diagonally
	if rank > 1 {
		leftDiagAttacker := (rank-2)*8 + file - 1
		rightDiagAttacker := (rank-2)*8 + file + 1

		leftAttack := leftFile && b.Bitboards[BP].Test(leftDiagAttacker)
		rightAttack := rightFile && b.Bitboards[BP].Test(rightDiagAttacker)

		if leftAttack || rightAttack {
			return true
		}
	}

	return false
}

// backward returns is a pawn is backward. It happens when the pawn is behind
// all the pawns of the same color on the adjancent files and cannot be safely advanced
func backward(b *board.Board, sq int) bool {
	if !b.Bitboards[WP].Test(sq) {
		return false
	}

	rank := sq / 8
	file := sq % 8

	for y := rank; y < 8; y++ {
		if (b.Bitboards[WP].Test(y*8+file-1) && file > 0) ||
			(file < 7 && b.Bitboards[WP].Test(y*8+file+1)) {
			return false
		}
	}

	if (b.Bitboards[BP].Test((rank-2)*8+file-1) && file > 0) ||
		(b.Bitboards[BP].Test((rank-2)*8+file+1) && file < 7) ||
		b.Bitboards[BP].Test((rank-1)*8+file) {
		return true
	}

	return false
}

// doubled checks if pawn is doubled. A doubled pawn is a pawn which has another friendly
// pawn on the same file but here we attach doubled pawn penalty only if pawn which has
// another friendly pawn on square directly behind that pawn and is not supported
func doubled(b *board.Board, sq int) bool {
	if !b.Bitboards[WP].Test(sq) {
		return false
	}

	rank := sq / 8
	file := sq % 8

	if !b.Bitboards[WP].Test((rank+1)*8 + file) {
		return false
	}

	if b.Bitboards[WP].Test((rank+1)*8+file-1) && file > 0 {
		return false
	}

	if b.Bitboards[WP].Test((rank+1)*8+file+1) && file < 7 {
		return false
	}

	return true
}

// connected checks is pawn is supported or phalanx
func connected(b *board.Board, sq int) bool {
	if supported(b, sq) > 0 || phalanx(b, sq) > 0 {
		return true
	}

	return false
}

// supported counts the number of pawns supporting this pawn. The pawn is supported
// if a friendly pawn is exacly in the adjancent file of the pawn and directly behind it
func supported(b *board.Board, sq int) int {
	if !b.Bitboards[WP].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	score := 0

	if b.Bitboards[WP].Test((rank+1)*8+file-1) && file > 0 {
		score++
	}

	if b.Bitboards[WP].Test((rank+1)*8+file+1) && file < 7 {
		score++
	}

	return score
}

// phalanx flag is set if there is friendly pawn on adjancent file and same rank
func phalanx(b *board.Board, sq int) int {
	if !b.Bitboards[WP].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if (b.Bitboards[WP].Test(rank*8+file-1) && file > 0) ||
		(b.Bitboards[WP].Test(rank*8+file+1) && file < 7) {
		return 1
	}

	return 0
}

// connectedBonus is the bonus for connected pawns
func connectedBonus(b *board.Board, sq int) int {
	if !connected(b, sq) {
		return 0
	}

	rank := 8 - sq/8

	seed := []int{0, 7, 8, 12, 29, 48, 86}
	op := opposed(b, sq)
	ph := phalanx(b, sq)
	su := supported(b, sq)
	if rank < 2 || rank > 7 {
		return 0
	}

	return seed[rank-1]*(2+ph-op) + 21*su
}

// opposed flag is set if there is opponent opposing pawn on the same file
// to prevent it from advancing
func opposed(b *board.Board, sq int) int {
	if !b.Bitboards[WP].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	for y := 0; y < rank; y++ {
		if b.Bitboards[BP].Test(y*8 + file) {
			return 1
		}
	}

	return 0
}

// weakUnopposedPawn checks if out pawn is weak and unopposed
func weakUnopposedPawn(b *board.Board, sq int) int {
	if opposed(b, sq) > 0 {
		return 0
	}
	score := 0

	if isolated(b, sq) {
		score++
	} else if backward(b, sq) {
		score++
	}

	return score
}

// blocked bonus for blocked pawns on the 5th or 6th rank
func blocked(b *board.Board, sq int) int {
	if !b.Bitboards[WP].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if rank != 2 && rank != 3 {
		return 0
	}

	if !b.Bitboards[BP].Test((rank-1)*8 + file) {
		return 0
	}

	return 4 - rank
}

// weakLever adds a penalty for unsupported pawns attacked twice by enemy pawns
func weakLever(b *board.Board, sq int) int {
	if !b.Bitboards[WP].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if !b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
		return 0
	}

	if !b.Bitboards[BP].Test((rank-1)*8+file+1) && rank > 0 && file < 7 {
		return 0
	}

	if b.Bitboards[WP].Test((rank+1)*8+file-1) && rank < 7 && file > 0 {
		return 0
	}

	if b.Bitboards[WP].Test((rank+1)*8+file+1) && rank < 7 && file < 7 {
		return 0
	}

	return 1
}
