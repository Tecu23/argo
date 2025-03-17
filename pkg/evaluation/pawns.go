// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// evaluatePawn returns bonuses and penalties for a specific pawn
func (e *Evaluator) evaluatePawn(b *board.Board, sq int) (mg, eg int) {
	if isIsolated(b, sq) {
		if doubleIsolated(b, sq) {
			mg -= 11
			eg -= 56
		} else {
			mg -= 5
			eg -= 15
		}
	} else if isBackward(b, sq) {
		mg -= 9
		eg -= 24
	}

	if isDoubled(b, sq) {
		mg -= 11
		eg -= 56
	}

	suppCount := supported(b, sq)
	isPh := isPhalanx(b, sq)

	if suppCount > 0 || isPh {

		seed := []int{0, 7, 8, 12, 29, 48, 86}
		r := 8 - sq/8
		bonus := 0

		if r >= 2 || r <= 7 {

			var op int
			if isOpposed(b, sq) {
				op = 1
			}

			ph := 0
			if isPh {
				ph = 1
			}

			bonus = seed[r-1]*(2+ph-op) + 21*suppCount
		}

		mg += bonus
		eg += bonus * (8 - (sq / 8) - 3) / 4
	}

	weakUnopposedBonus := 0
	if weakUnopposedPawn(b, sq) {
		weakUnopposedBonus = 1
	}
	mg -= 13 * weakUnopposedBonus
	eg -= 27 * weakUnopposedBonus

	weakLeverBonus := 0
	if isWeakLever(b, sq) {
		weakLeverBonus = 1
	}
	eg -= 56 * weakLeverBonus

	blockedIdx := blocked(b, sq)
	mg += []int{0, -11, -3}[blockedIdx]
	eg += []int{0, -4, 4}[blockedIdx]

	return mg, eg
}

// doubleIsolated is a penalty if a double pawn is stopped only
// by a single opponent pawn on the same file.
func doubleIsolated(b *board.Board, sq int) bool {
	// Should return if square doesn't contain a white pawn
	// But because we only call double isolated on white pawns this is not needed

	file := sq % 8
	rank := sq / 8

	// NOTE: No need for this as we check if pawn isolated first

	// // Return early is the pawn is not isolated
	// if !isolated(b, sq) {
	// 	return false
	// }

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
func isIsolated(b *board.Board, sq int) bool {
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
		aheadRankMask := RankMasks[rank-2]

		leftAttack := leftFile && b.Bitboards[BP]&aheadRankMask&leftFileMask != 0
		rightAttack := rightFile && b.Bitboards[BP]&aheadRankMask&rightFileMask != 0

		if leftAttack || rightAttack {
			return true
		}
	}

	return false
}

// backward returns is a pawn is backward. It happens when the pawn is behind
// all the pawns of the same color on the adjancent files and cannot be safely advanced
func isBackward(b *board.Board, sq int) bool {
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
func isDoubled(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	if rank >= 7 {
		return false
	}

	if !b.Bitboards[WP].Test((rank+1)*8 + file) {
		return false
	}

	if file > 0 && b.Bitboards[WP].Test((rank+1)*8+file-1) {
		return false
	}

	if file < 7 && b.Bitboards[WP].Test((rank+1)*8+file+1) {
		return false
	}

	return true
}

// supported counts the number of pawns supporting this pawn. The pawn is supported
// if a friendly pawn is exacly in the adjancent file of the pawn and directly behind it
func supported(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if rank >= 7 {
		return 0
	}

	score := 0

	if file > 0 && b.Bitboards[WP].Test((rank+1)*8+file-1) {
		score++
	}

	if file < 7 && b.Bitboards[WP].Test((rank+1)*8+file+1) {
		score++
	}

	return score
}

// isPhalanx flag is set if there is friendly pawn on adjancent file and same rank
func isPhalanx(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	if (file > 0 && b.Bitboards[WP].Test(rank*8+file-1)) ||
		(file < 7 && b.Bitboards[WP].Test(rank*8+file+1)) {
		return true
	}

	return false
}

// opposed flag is set if there is opponent opposing pawn on the same file
// to prevent it from advancing
func isOpposed(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	fileMask := FileMasks[file]

	var frontRanksMask bitboard.Bitboard
	for r := 0; r < rank; r++ {
		frontRanksMask |= RankMasks[r]
	}

	blackPawnsInFront := b.Bitboards[BP] & fileMask & frontRanksMask

	if blackPawnsInFront > 0 {
		return true
	}

	return false
}

// weakUnopposedPawn checks if out pawn is weak and unopposed
func weakUnopposedPawn(b *board.Board, sq int) bool {
	if isOpposed(b, sq) {
		return false
	}

	if isIsolated(b, sq) {
		return true
	}

	if isBackward(b, sq) {
		return true
	}

	return false
}

// blocked bonus for blocked pawns on the 5th or 6th rank
func blocked(b *board.Board, sq int) int {
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

// isWeakLever adds a penalty for unsupported pawns attacked twice by enemy pawns
func isWeakLever(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	if rank > 0 && file > 0 && !b.Bitboards[BP].Test((rank-1)*8+file-1) {
		return false
	}

	if rank > 0 && file < 7 && !b.Bitboards[BP].Test((rank-1)*8+file+1) {
		return false
	}

	if rank < 7 && file > 0 && b.Bitboards[WP].Test((rank+1)*8+file-1) {
		return false
	}

	if rank < 7 && file < 7 && b.Bitboards[WP].Test((rank+1)*8+file+1) {
		return false
	}

	return true
}

// func (e *Evaluator) optimizedEvaluatePawn(b *board.Board, sq int) (mg, eg int) {
// 	file := sq % 8
// 	rank := sq / 8
//
// 	leftFile := file > 0
// 	rightFile := file < 7
//
// 	// Create masks for adjancent files once per pawn
// 	var leftFileMask, rightFileMask bitboard.Bitboard
// 	if leftFile {
// 		leftFileMask = FileMasks[file-1]
// 	}
//
// 	if rightFile {
// 		rightFileMask = FileMasks[file+1]
// 	}
//
// 	// Check for isolated pawns
// 	adjancentFilePawns := (leftFile && (b.Bitboards[WP]&leftFileMask) != 0) ||
// 		(rightFile && (b.Bitboards[WP]&rightFileMask) != 0)
//
// 	if !adjancentFilePawns {
// 		isDoubleIsolated := false
//
// 		// Check for double isolated condition only if pawn is already isolated
// 		fileMask := FileMasks[file]
//
// 		// Check for another white pawn behind this one (greater rank values on same file)
// 		var behindRanksMask bitboard.Bitboard
// 		for r := rank + 1; r < 8; r++ {
// 			behindRanksMask |= RankMasks[r]
// 		}
//
// 		// If no pawn behind, then the pawn is not doubled
// 		whitePawnsBehind := b.Bitboards[WP] & fileMask & behindRanksMask
// 		if whitePawnsBehind != 0 {
// 			// Create a mask for ranks in front of our pawn
// 			var frontRankMask bitboard.Bitboard
// 			for r := 0; r < rank; r++ {
// 				frontRankMask |= RankMasks[r]
// 			}
//
// 			// Check for black pawns in front on the same file
// 			blackPawnsOnFileInFront := b.Bitboards[BP] & fileMask & frontRankMask
// 			if blackPawnsOnFileInFront != 0 {
// 				// Check for enemy pawns on adjacent files (any rank)
// 				blackPawnsOnAdjacentFiles := (leftFile && (b.Bitboards[BP]&leftFileMask != 0)) ||
// 					(rightFile && (b.Bitboards[BP]&rightFileMask != 0))
//
// 				// If there's a doubled white pawn, black pawns in front on the same file,
// 				// and no black pawns on adjacent files, return true
// 				isDoubleIsolated = !blackPawnsOnAdjacentFiles
// 			}
// 		}
//
// 		if isDoubleIsolated {
// 			mg -= 11
// 			eg -= 56
// 		} else {
// 			mg -= 5
// 			eg -= 15
// 		}
// 	} else if backward(b, sq) {
// 		mg -= 9
// 		eg -= 24
// 	}
//
// 	// Check for doubled pawns - more efficient implementation
// 	if isDoubled(b, sq) {
// 		mg -= 11
// 		eg -= 56
// 	}
//
// 	// Connected pawns bonus
// 	supportedCount := countSupported(b, sq)
// 	phalanxCount := countPhalanx(b, sq)
// 	isOpp := isOpposed(b, sq)
//
// 	if supportedCount > 0 || phalanxCount > 0 {
// 		// Calculate connected bonus
// 		connBonus := calculateConnectedBonus(rank, isOpposed, phalanxCount, supportedCount)
//
// 		mg += connBonus
// 		eg += connBonus * (8 - rank - 3) / 4
// 	}
//
// 	// Weak unopposed pawn penalty
// 	if !isOpp && (!has)
//
// 	return mg, eg
// }
