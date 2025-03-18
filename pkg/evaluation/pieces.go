// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"math"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func (e *Evaluator) evaluateKnight(b *board.Board, sq int) (mg, eg int) {
	minorBonus := minorBehindPawn(b, sq)
	mg += 18 * minorBonus
	eg += 3 * minorBonus

	kingProtectorBonus := kingProtector(b, sq)
	mg -= 8 * kingProtectorBonus
	eg -= 9 * kingProtectorBonus

	return mg, eg
}

func (e *Evaluator) evaluateBishop(b *board.Board, sq int) (mg, eg int) {
	minorBonus := minorBehindPawn(b, sq)
	mg += 18 * minorBonus
	eg += 3 * minorBonus

	bishopPawnsBonus := bishopPawns(b, sq)
	mg -= 3 * bishopPawnsBonus
	eg -= 7 * bishopPawnsBonus

	bishopXrayBonus := bishopXrayPawns(b, sq)
	mg -= 4 * bishopXrayBonus
	eg -= 5 * bishopXrayBonus

	kingProtectorBonus := kingProtector(b, sq)
	mg -= 6 * kingProtectorBonus
	eg -= 9 * kingProtectorBonus

	// mg += 24 * bishopOnKingRing(b, sq)
	mg += 45 * longDiagonalBishop(b, sq)

	return mg, eg
}

func (e *Evaluator) evaluateRook(b *board.Board, sq int) (mg, eg int) {
	rookOnQueenBonus := rookOnQueenFile(b, sq)
	mg += 6 * rookOnQueenBonus
	eg += 11 * rookOnQueenBonus

	rookOnFileIdx := rookOnFile(b, sq)
	mg += []int{0, 19, 48}[rookOnFileIdx]
	eg += []int{0, 7, 29}[rookOnFileIdx]

	// mg += 16 * rookOnKingRing(b, sq)

	// factor := 2
	// if uint(b.Castlings)&ShortW != 0 || uint(b.Castlings)&LongW != 0 {
	// 	factor = 1
	// }

	// trappedRookBonus := trappedRook(b, mirror, sq)
	// mg -= trappedRookBonus * 55 * factor
	// eg -= trappedRookBonus * 13 * factor

	return mg, eg
}

func (e *Evaluator) evaluateQueen(b *board.Board, sq int) (mg, eg int) {
	// weakQueenBonus := weakQueen(b, sq)
	// mg -= 56 * weakQueenBonus
	// eg -= 15 * weakQueenBonus

	queenInfiltrationBonus := queenInfiltration(b, sq)
	mg -= 2 * queenInfiltrationBonus
	eg += 14 * queenInfiltrationBonus

	return mg, eg
}

// minorBehindPawn return whether the bishop/knight is begind a pawn
func minorBehindPawn(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if rank > 0 &&
		!b.Bitboards[WP].Test((rank-1)*8+file) &&
		!b.Bitboards[BP].Test((rank-1)*8+file) {
		return 0
	}

	return 1
}

// kingProtector adds penalties and bonuses for pieces, depending on the distance
// from the own king
func kingProtector(b *board.Board, sq int) int {
	return kingDistance(b, sq)
}

// kingDistance counts the distance to our king
func kingDistance(b *board.Board, sq int) int {
	kingBB := b.Bitboards[WK]
	kingSq := kingBB.FirstOne()

	return max(abs((kingSq/8)-(sq/8)), abs((kingSq%8)-(sq%8)))
}

// bishopPawns returns the number of pawns on the same color square
// as the bishop multiplied by one of our blocked pawns in the center files C,D,E or F
// NOTE: Could this be improved by keeping track of pawns on white sq or pawns on black sq?
func bishopPawns(b *board.Board, sq int) int {
	score := 0
	c := (sq/8 + sq%8) % 2
	blocked := 0

	pawnsBB := b.Bitboards[WP]
	for pawnsBB != 0 {
		pawnSq := pawnsBB.FirstOne()

		pawnRank := pawnSq / 8
		pawnFile := pawnSq % 8

		if (pawnFile+pawnRank)%2 == c {
			score++
		}

		if pawnFile > 1 && pawnFile < 6 {
			squareInFront := (pawnRank-1)*8 + pawnFile

			if squareInFront >= 0 && b.Occupancies[color.BOTH].Test(squareInFront) {
				blocked++
			}
		}
	}

	pawnAttack := 1
	if PawnAttack(b, sq) > 0 {
		pawnAttack = 0
	}

	score = score * (blocked + pawnAttack)

	return score
}

// bishopXrayPawns is a penalty for all enemy pawns xrayed by our bishop
func bishopXrayPawns(b *board.Board, sq int) int {
	count := 0

	rank := sq / 8
	file := sq % 8

	pawnsBB := b.Bitboards[BP]
	for pawnsBB != 0 {
		pawnSq := pawnsBB.FirstOne()

		pawnRank := pawnSq / 8
		pawnFile := pawnSq % 8

		if abs(file-pawnFile) == abs(rank-pawnRank) {
			count++
		}
	}

	return count
}

// bishopOnKingRing gives bonus for bishops that are alligned with the
// enemy kingring.
func bishopOnKingRing(b *board.Board, sq int) int {
	if kingAttackersCount(b, sq) > 0 {
		return 0
	}
	factor1, factor2 := 0, 0

	rank := sq / 8
	file := sq % 8

	for i := 0; i < 4; i++ {
		factor1, factor2 = 0, 0
		if i > 1 {
			factor1 = 1
		}

		if i%2 == 0 {
			factor2 = 1
		}

		ix := factor1*2 - 1
		iy := factor2*2 - 1

		for d := 1; d < 8; d++ {
			x := file + d*ix
			y := rank + d*iy

			if x < 0 || x > 7 || y < 0 || y > 7 {
				break
			}

			if kingRing(b, y*8+x, false) > 0 {
				return 1
			}

			if b.Bitboards[BP].Test(y*8+x) || b.Bitboards[WP].Test(y*8+x) {
				break
			}
		}
	}
	return 0
}

// longDiagonalBishop is a bonus for bishop on a long diagonal which
// can "see" both center square
func longDiagonalBishop(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8
	if file-rank != 0 && file-(7-rank) != 0 {
		return 0
	}

	if min(file, 7-file) > 2 {
		return 0
	}

	x1, y1 := file, rank

	for i := min(x1, 7-x1); i < 4; i++ {
		if b.Bitboards[BP].Test(y1*8 + x1) {
			return 0
		}

		if b.Bitboards[WP].Test(y1*8 + x1) {
			return 0
		}

		if x1 < 4 {
			x1++
		} else {
			x1--
		}

		if y1 < 4 {
			y1++
		} else {
			y1--
		}
	}

	return 1
}

// rookOnQueenFile is a simple bonus for a rook that is on the same file as any queen
func rookOnQueenFile(b *board.Board, sq int) int {
	file := sq % 8

	for y := 0; y < 8; y++ {
		if b.Bitboards[WQ].Test(y*8+file) || b.Bitboards[BQ].Test(y*8+file) {
			return 1
		}
	}

	return 0
}

// rookOnKingRing gives bonus for rooks that are alligned with the enemy
// king ring
func rookOnKingRing(b *board.Board, sq int) int {
	if kingAttackersCount(b, sq) > 0 {
		return 0
	}

	file := sq % 8

	for y := 0; y < 8; y++ {
		if kingRing(b, y*8+file, false) > 0 {
			return 1
		}
	}

	return 0
}

// trappedRook penalizes the took when is trapped by the king, even more
// if the king cannot castle
func trappedRook(b *board.Board, mirror *board.Board, sq int) int {
	if rookOnFile(b, sq) > 0 {
		return 0
	}

	if mobility(b, mirror, sq) > 3 {
		return 0
	}

	kingBB := b.Bitboards[WK]
	kingSq := kingBB.FirstOne()

	kx := kingSq % 8

	if kx < 4 != ((sq % 8) < kx) {
		return 0
	}

	return 1
}

// RookOnFile returns whether the took is on open / semi-open file
func rookOnFile(b *board.Board, sq int) int {
	open := 1
	file := sq % 8

	for y := 0; y < 8; y++ {
		if b.Bitboards[WP].Test(y*8 + file) {
			return 0
		}
		if b.Bitboards[BP].Test(y*8 + file) {
			open = 0
		}
	}

	return open + 1
}

// queenInfiltration is a bonus for queen on weak square in enemy camp,
// Idea is that queen feels much better when it can't be kicked away now or later
// by pawn moves, especially in endgame
func queenInfiltration(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if rank > 3 {
		return 0
	}

	if b.Bitboards[BP].Test((rank-1)*file+1) && file < 7 {
		return 0
	}

	if b.Bitboards[BP].Test((rank-1)*file-1) && file > 0 {
		return 0
	}

	if pawnAttacksSpan(b, sq) > 0 {
		return 0
	}

	return 1
}

// pawnAttacksSpan compute additional span if pawn is not backward nor blocked
func pawnAttacksSpan(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8
	mirror := b.Mirror()

	for y := 0; y < rank; y++ {
		if b.Bitboards[BP].Test(y*8+file-1) && file > 0 &&
			(y == file-1 || (b.Bitboards[WP].Test((y+1)*8+file-1) && file > 0 && !isBackward(mirror, (7-y)*8+file-1))) {
			return 1
		}

		if b.Bitboards[BP].Test(y*8+file+1) && file < 7 &&
			(y == file-1 || (b.Bitboards[WP].Test((y+1)*8+file+1) && file < 7 && !isBackward(mirror, (7-y)*8+file+1))) {
			return 1
		}
	}

	return 0
}

// weakQueen returns a penalty if any relative pin or discovered attack
// against the queen
func weakQueen(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8
	for i := 0; i < 8; i++ {
		factor := 0
		if i > 3 {
			factor = 1
		}

		ix := (i+factor)%3 - 1
		iy := (((i + factor) / 3) << 0) - 1
		count := 0

		for d := 1; d < 8; d++ {
			if rank+d*iy < 0 || rank+d*iy > 7 || file+d*ix < 0 || file+d*ix > 7 {
				continue
			}
			b := b.GetPieceAt((rank+d*iy)*8 + (file + d*ix))

			if b == BR && (ix == 0 || iy == 0) && count == 1 {
				return 1
			}

			if b == BB && (ix != 0 && iy != 0) && count == 1 {
				return 1
			}

			if b != Empty {
				count++
			}

		}
	}

	return 0
}

// kingAttackersCount returns the number of pieces of the given color which
// attack a square in the kingring of the enemy king. For pawns we count the
// number of attacked squares in kingring
func kingAttackersCount(b *board.Board, sq int) int {
	if !b.Occupancies[color.WHITE].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if b.Bitboards[WP].Test(sq) {
		score := 0.0

		for dir := -1; dir <= 1; dir += 2 {
			fr := 1.0

			if b.Bitboards[WP].Test(rank*8 + file + dir*2) {
				fr = 0.5
			}

			if file+dir >= 0 && file+dir <= 7 && kingRing(b, (rank-1)*8+file+dir, true) > 0 {
				score = score + fr
			}
		}
		return int(math.Round(float64(score)))
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if kingRing(b, y*8+x, false) > 0 {
				if KnightAttack(b, y*8+x, sq) > 0 ||
					BishopXrayAttack(b, y*8+x, sq) > 0 ||
					RookXrayAttack(b, y*8+x, sq) > 0 ||
					QueenAttack(b, y*8+x, sq) > 0 {
					return 1
				}
			}
		}
	}
	return 0
}

// PERF:kingRing is square occupied by king and 8 squares around the king. Squares
// defended by two pawns are removed from king ring
func kingRing(b *board.Board, sq int, full bool) int {
	rank := sq / 8
	file := sq % 8

	if !full && rank > 0 && file < 7 && file > 0 &&
		b.Bitboards[BP].Test((rank-1)*8+file+1) &&
		b.Bitboards[BP].Test((rank-1)*8+file-1) {
		return 0
	}

	bb := b.Bitboards[BK]
	kingSq := bb.FirstOne()

	fileMask := FileMasks[file]
	rankMask := RankMasks[rank]

	if fileMask&rankMask&KingRingPatterns[kingSq] != 0 {
		return 1
	}

	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
