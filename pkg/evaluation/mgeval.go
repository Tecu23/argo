package evaluation

import (
	"fmt"
	"math"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	evalhelpers "github.com/Tecu23/argov2/pkg/evaluation/helpers"
)

func PieceValueMg(b *board.Board) int {
	PawnBonus := 124
	KnightBonus := 781
	BishopBonus := 825
	RookBonus := 1276
	QueenBonus := 2538

	pawnScore := b.Bitboards[WP].Count() * PawnBonus
	knightScore := b.Bitboards[WN].Count() * KnightBonus
	bishopScore := b.Bitboards[WB].Count() * BishopBonus
	rookScore := b.Bitboards[WR].Count() * RookBonus
	queenScore := b.Bitboards[WQ].Count() * QueenBonus

	return pawnScore + knightScore + bishopScore + rookScore + queenScore
}

func PsqtMg(b *board.Board) int {
	bonus := [][][]int{
		{
			{-175, -92, -74, -73},
			{-77, -41, -27, -15},
			{-61, -17, 6, 12},
			{-35, 8, 40, 49},
			{-34, 13, 44, 51},
			{-9, 22, 58, 53},
			{-67, -27, 4, 37},
			{-201, -83, -56, -26},
		},
		{
			{-53, -5, -8, -23},
			{-15, 8, 19, 4},
			{-7, 21, -5, 17},
			{-5, 11, 25, 39},
			{-12, 29, 22, 31},
			{-16, 6, 1, 11},
			{-17, -14, 5, 0},
			{-48, 1, -14, -23},
		},
		{
			{-31, -20, -14, -5},
			{-21, -13, -8, 6},
			{-25, -11, -1, 3},
			{-13, -5, -4, -6},
			{-27, -15, -4, 3},
			{-22, -2, 6, 12},
			{-2, 12, 16, 18},
			{-17, -19, -1, 9},
		},
		{
			{3, -5, -5, 4},
			{-3, 5, 8, 12},
			{-3, 6, 13, 7},
			{4, 5, 9, 8},
			{0, 14, 12, 5},
			{-4, 10, 6, 8},
			{-5, 6, 10, 8},
			{-2, -2, 1, -2},
		},
		{
			{271, 327, 271, 198},
			{278, 303, 234, 179},
			{195, 258, 169, 120},
			{164, 190, 138, 98},
			{154, 179, 105, 70},
			{123, 145, 81, 31},
			{88, 120, 65, 33},
			{59, 89, 45, -1},
		},
	}

	pBonus := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{3, 3, 10, 19, 16, 19, 7, -5},
		{-9, -15, 11, 15, 32, 22, 5, -22},
		{-4, -23, 6, 20, 40, 17, 4, -8},
		{13, 0, -13, 1, 11, -2, -13, 5},
		{5, -12, -7, 22, -8, -5, -15, -8},
		{-7, 7, -3, -13, 5, -16, 10, -8},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}

	score := 0

	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += pBonus[7-rank][file]
	}

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq := knightBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += bonus[0][7-rank][min(file, 7-file)]
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += bonus[1][7-rank][min(file, 7-file)]
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq := rookBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += bonus[2][7-rank][min(file, 7-file)]
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq := queenBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += bonus[3][7-rank][min(file, 7-file)]
	}

	kingBB := b.Bitboards[WK]
	for kingBB != 0 {
		sq := kingBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		score += bonus[4][7-rank][min(file, 7-file)]
	}

	return score
}

func ImbalanceTotal(b *board.Board, mirror *board.Board) int {
	score := 0

	score += Imbalance(b) - Imbalance(mirror)
	score += BishopPair(b) - BishopPair(mirror)

	return (score / 16) << 0
}

func Imbalance(b *board.Board) int {
	score := 0

	qo := [][]int{
		{0},
		{40, 38},
		{32, 255, -62},
		{0, 104, 4, 0},
		{-26, -2, 47, 105, -208},
		{-189, 24, 117, 133, -134, -6},
	}
	qt := [][]int{
		{0},
		{36, 0},
		{9, 63, 0},
		{59, 65, 42, 0},
		{46, 39, 24, -24, 0},
		{97, 100, -42, 137, 268, 0},
	}

	bishop := []int{0, 0}
	bishop[0] = b.Bitboards[BB].Count()
	bishop[1] = b.Bitboards[WB].Count()

	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		_ = pawnBB.FirstOne()

		score += qt[1][1] * b.Bitboards[BP].Count()
		score += qo[1][1] * b.Bitboards[WP].Count()

		if bishop[0] > 1 {
			score += qt[1][0]
		}

		if bishop[1] > 1 {
			score += qo[1][0]
		}
	}

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		_ = knightBB.FirstOne()

		score += qt[2][1] * b.Bitboards[BP].Count()
		score += qo[2][1] * b.Bitboards[WP].Count()

		score += qt[2][2] * b.Bitboards[BN].Count()
		score += qo[2][2] * b.Bitboards[WN].Count()

		if bishop[0] > 1 {
			score += qt[2][0]
		}

		if bishop[1] > 1 {
			score += qo[2][0]
		}
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		_ = bishopBB.FirstOne()

		score += qt[3][1] * b.Bitboards[BP].Count()
		score += qo[3][1] * b.Bitboards[WP].Count()

		score += qt[3][2] * b.Bitboards[BN].Count()
		score += qo[3][2] * b.Bitboards[WN].Count()

		score += qt[3][3] * b.Bitboards[BB].Count()
		score += qo[3][3] * b.Bitboards[WB].Count()

		if bishop[0] > 1 {
			score += qt[3][0]
		}

		if bishop[1] > 1 {
			score += qo[3][0]
		}
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		_ = rookBB.FirstOne()

		score += qt[4][1] * b.Bitboards[BP].Count()
		score += qo[4][1] * b.Bitboards[WP].Count()

		score += qt[4][2] * b.Bitboards[BN].Count()
		score += qo[4][2] * b.Bitboards[WN].Count()

		score += qt[4][3] * b.Bitboards[BB].Count()
		score += qo[4][3] * b.Bitboards[WB].Count()

		score += qt[4][4] * b.Bitboards[BR].Count()
		score += qo[4][4] * b.Bitboards[WR].Count()

		if bishop[0] > 1 {
			score += qt[4][0]
		}

		if bishop[1] > 1 {
			score += qo[4][0]
		}
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		_ = queenBB.FirstOne()

		score += qt[5][1] * b.Bitboards[BP].Count()
		score += qo[5][1] * b.Bitboards[WP].Count()

		score += qt[5][2] * b.Bitboards[BN].Count()
		score += qo[5][2] * b.Bitboards[WN].Count()

		score += qt[5][3] * b.Bitboards[BB].Count()
		score += qo[5][3] * b.Bitboards[WB].Count()

		score += qt[5][4] * b.Bitboards[BR].Count()
		score += qo[5][4] * b.Bitboards[WR].Count()

		score += qt[5][5] * b.Bitboards[BQ].Count()
		score += qo[5][5] * b.Bitboards[WQ].Count()

		if bishop[0] > 1 {
			score += qt[5][0]
		}

		if bishop[1] > 1 {
			score += qo[5][0]
		}
	}

	return score
}

func BishopPair(b *board.Board) int {
	if b.Bitboards[WB].Count() < 2 {
		return 0
	}

	return 1438
}

func PawnsMg(b *board.Board) int {
	score := 0

	pawnsBB := b.Bitboards[WP]

	for pawnsBB != 0 {
		sq := pawnsBB.FirstOne()

		if DoubleIsolated(b, sq) {
			score -= 11
		} else if Isolated(b, sq) {
			score -= 5
		} else if Backward(b, sq) {
			score -= 9
		}

		if Doubled(b, sq) {
			score -= 11
		}

		if Connected(b, sq) {
			score += ConnectedBonus(b, sq)
		}

		score -= 13 * WeakUnopposedPawn(b, sq)
		score += []int{0, -11, -3}[Blocked(b, sq)]
	}

	return score
}

func Blocked(b *board.Board, sq int) int {
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

func WeakUnopposedPawn(b *board.Board, sq int) int {
	if Opposed(b, sq) > 0 {
		return 0
	}
	score := 0

	if Isolated(b, sq) {
		score++
	} else if Backward(b, sq) {
		score++
	}

	return score
}

func ConnectedBonus(b *board.Board, sq int) int {
	if !Connected(b, sq) {
		return 0
	}

	rank := 8 - sq/8

	seed := []int{0, 7, 8, 12, 29, 48, 86}
	op := Opposed(b, sq)
	ph := Phalanx(b, sq)
	su := Supported(b, sq)
	if rank < 2 || rank > 7 {
		return 0
	}

	return seed[rank-1]*(2+ph-op) + 21*su
}

func Opposed(b *board.Board, sq int) int {
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

func Connected(b *board.Board, sq int) bool {
	if Supported(b, sq) > 0 || Phalanx(b, sq) > 0 {
		return true
	}

	return false
}

func Phalanx(b *board.Board, sq int) int {
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

func Supported(b *board.Board, sq int) int {
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

func Doubled(b *board.Board, sq int) bool {
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

func Backward(b *board.Board, sq int) bool {
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

// Doubled Isolated is a penalty if a double pawn is stopped only
// by a single opponent pawn on the same file.
func DoubleIsolated(b *board.Board, sq int) bool {
	if !b.Bitboards[WP].Test(sq) {
		return false
	}

	if Isolated(b, sq) {
		obe, eop, ene := 0, 0, 0

		rank := sq / 8
		file := sq % 8

		for y := 0; y < 8; y++ {
			if y > rank && b.Bitboards[WP].Test(y*8+file) {
				obe++
			}

			if y < rank && b.Bitboards[BP].Test(y*8+file) {
				eop++
			}

			if (file > 0 && b.Bitboards[BP].Test(y*8+file-1)) ||
				(b.Bitboards[BP].Test(y*8+file+1) && file < 7) {
				ene++
			}
		}

		if obe > 0 && ene == 0 && eop > 0 {
			return true
		}

	}

	return false
}

func Isolated(b *board.Board, sq int) bool {
	file := sq % 8

	if !b.Bitboards[WP].Test(sq) {
		return false
	}

	for y := 0; y < 8; y++ {
		if (b.Bitboards[WP].Test(y*8+file-1) && file > 0) ||
			(b.Bitboards[WP].Test(y*8+file+1) && file < 7) {
			return false
		}
	}

	return true
}

func PiecesMg(b *board.Board) int {
	score := 0

	// scores for Knight, Bishop, Rook, Queen
	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq := knightBB.FirstOne()

		score += []int{0, 31, -7, 30, 56}[OutpostTotal(b, sq)]
		score += 18 * MinorBehindPawn(b, sq)
		score -= 8 * KingProtector(b, sq)

	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()

		score += []int{0, 31, -7, 30, 56}[OutpostTotal(b, sq)]
		score += 18 * MinorBehindPawn(b, sq)
		score -= 3 * BishopPawns(b, sq)
		score -= 4 * BishopXrayPawns(b, sq)
		score += 24 * BishopOnKingRing(b, sq)
		score -= 6 * KingProtector(b, sq)
		score += 45 * LongDiagonalBishop(b, sq)

	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq := rookBB.FirstOne()

		score += 6 * RookOnQueenFile(b, sq)
		score += 16 * RookOnKingRing(b, sq)

		score += []int{0, 19, 48}[RookOnFile(b, sq)]

		factor := 2
		if uint(b.Castlings)&ShortW != 0 || uint(b.Castlings)&LongW != 0 {
			factor = 1
		}

		score -= TrappedRook(b, sq) * 55 * factor
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq := queenBB.FirstOne()

		score -= 56 * WeakQueen(b, sq)
		score -= 2 * QueenInfiltration(b, sq)
	}

	return score
}

func LongDiagonalBishop(b *board.Board, sq int) int {
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

func KingProtector(b *board.Board, sq int) int {
	return KingDistance(b, sq)
}

func KingDistance(b *board.Board, sq int) int {
	kingBB := b.Bitboards[WK]
	kingSq := kingBB.FirstOne()

	return max(abs((kingSq/8)-(sq/8)), abs((kingSq%8)-(sq%8)))
}

func QueenInfiltration(b *board.Board, sq int) int {
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

	if PawnAttacksSpan(b, sq) > 0 {
		return 0
	}

	return 1
}

func PawnAttacksSpan(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8
	mirror := b.Mirror()

	for y := 0; y < rank; y++ {
		if b.Bitboards[BP].Test(y*8+file-1) && file > 0 &&
			(y == file-1 || (b.Bitboards[WP].Test((y+1)*8+file-1) && file > 0 && !Backward(mirror, (7-y)*8+file-1))) {
			return 1
		}

		if b.Bitboards[BP].Test(y*8+file+1) && file < 7 &&
			(y == file-1 || (b.Bitboards[WP].Test((y+1)*8+file+1) && file < 7 && !Backward(mirror, (7-y)*8+file+1))) {
			return 1
		}
	}

	return 0
}

// WeakQueen returns a penalty if any relative pin or discovered attack
// against the queen
func WeakQueen(b *board.Board, sq int) int {
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
				fmt.Println("here 1")
				return 1
			}

			if b == BB && (ix != 0 && iy != 0) && count == 1 {
				fmt.Println("here 2", b, ix, iy, count)
				return 1
			}

			if b != Empty {
				count++
			}

		}
	}

	return 0
}

// TrappedRook penalizes the took when is trapped by the king, even more
// if the king cannot castle
func TrappedRook(b *board.Board, sq int) int {
	if RookOnFile(b, sq) > 0 {
		return 0
	}

	if Mobility(b, sq) > 3 {
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
func RookOnFile(b *board.Board, sq int) int {
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

// BishopOnKingRing gives bonus for bishops that are alligned with the
// enemy kingring.
func BishopOnKingRing(b *board.Board, sq int) int {
	if KingAttackersCount(b, sq) > 0 {
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

			if KingRing(b, y*8+x, false) > 0 {
				return 1
			}

			if b.Bitboards[BP].Test(y*8+x) || b.Bitboards[WP].Test(y*8+x) {
				break
			}
		}
	}
	return 0
}

// RookOnKingRing gives bonus for rooks that are alligned with the enemy
// king ring
func RookOnKingRing(b *board.Board, sq int) int {
	if KingAttackersCount(b, sq) > 0 {
		return 0
	}

	file := sq % 8

	for y := 0; y < 8; y++ {
		if KingRing(b, y*8+file, false) > 0 {
			return 1
		}
	}

	return 0
}

// KingAttackersCount returns the number of pieces of the given color which
// attack a square in the kingring of the enemy king. For pawns we count the
// number of attacked squares in kingring
func KingAttackersCount(b *board.Board, sq int) int {
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

			if file+dir >= 0 && file+dir <= 7 && KingRing(b, (rank-1)*8+file+dir, true) > 0 {
				score = score + fr
			}
		}
		return int(math.Round(float64(score)))
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if KingRing(b, y*8+x, false) > 0 {
				if evalhelpers.KnightAttack(b, y*8+x, sq) > 0 ||
					evalhelpers.BishopXrayAttack(b, y*8+x, sq) > 0 ||
					evalhelpers.RookXrayAttack(b, y*8+x, sq) > 0 ||
					evalhelpers.QueenAttack(b, y*8+x, sq) > 0 {
					return 1
				}
			}
		}
	}
	return 0
}

// KingRing is square occupied by king and 8 squares around the king. Squares
// defended by two pawns are removed from king ring
func KingRing(b *board.Board, sq int, full bool) int {
	rank := sq / 8
	file := sq % 8

	if !full && b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 && file > 0 &&
		b.Bitboards[BP].Test((rank-1)*8+file-1) {
		return 0
	}

	for ix := -2; ix <= 2; ix++ {
		for iy := -2; iy <= 2; iy++ {
			if ix+file < 0 || ix+file > 7 || iy+rank < 0 || iy+rank > 7 {
				continue
			}

			if b.Bitboards[BK].Test(
				(rank+iy)*8+file+ix,
			) && (ix >= -1 && ix <= 1 || file+ix == 0 || file+ix == 7) &&
				(iy >= -1 && iy <= 1 || rank+iy == 0 || rank+iy == 7) {
				return 1
			}
		}
	}

	return 0
}

// RookOnQueenFile is a simple bonus for a rook that is on the same file as any queen
func RookOnQueenFile(b *board.Board, sq int) int {
	file := sq % 8

	for y := 0; y < 8; y++ {
		if b.Bitboards[WQ].Test(y*8+file) || b.Bitboards[BQ].Test(y*8+file) {
			return 1
		}
	}

	return 0
}

// BishopXrayPawns is a penalty for all enemy pawns xrayed by our bishop
func BishopXrayPawns(b *board.Board, sq int) int {
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

// BishopPawns returns the number of pawns on the same color square
// as the bishop multiplied by one of our blocked pawns in the center files C,D,E or F
func BishopPawns(b *board.Board, sq int) int {
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
	if evalhelpers.PawnAttack(b, sq) > 0 {
		pawnAttack = 0
	}

	score = score * (blocked + pawnAttack)

	return score
}

// MinorBehindPawn return whether the bishop/knight is begind a pawn
func MinorBehindPawn(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if !b.Bitboards[WP].Test((rank-1)*8+file) && !b.Bitboards[BP].Test((rank-1)*8+file) {
		return 0
	}

	return 1
}

// TODO: To be implemented
func OutpostTotal(b *board.Board, sq int) int {
	return 0
}

// MobilityMg returnes the mobility bonus for middlegame
func MobilityMg(b *board.Board) int {
	score := 0
	sq := 0

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq = knightBB.FirstOne()

		score += MobilityBonus(b, sq, true)
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq = bishopBB.FirstOne()

		score += MobilityBonus(b, sq, true)
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq = rookBB.FirstOne()

		score += MobilityBonus(b, sq, true)
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq = queenBB.FirstOne()

		score += MobilityBonus(b, sq, true)
	}
	return score
}

// MobilityBonus attaches bonuses for middlegame and endgame by piece type and Mobility
func MobilityBonus(b *board.Board, sq int, isMiddleGame bool) int {
	bonus := [][]int{
		{-81, -56, -31, -16, 5, 11, 17, 20, 25},
		{-59, -23, -3, 13, 24, 42, 54, 57, 65, 73, 78, 86, 88, 97},
		{-78, -17, 23, 39, 70, 99, 103, 121, 134, 139, 158, 164, 168, 169, 172},
		{
			-48, -30, -7, 19, 40, 55, 59, 75, 78, 96, 96, 100, 121,
			127, 131, 133, 136, 141, 147, 150, 151, 168, 168, 171, 182, 182, 192, 219,
		},
	}

	if isMiddleGame {
		bonus = [][]int{
			{-62, -53, -12, -4, 3, 13, 22, 28, 33},
			{-48, -20, 16, 26, 38, 51, 55, 63, 63, 68, 81, 81, 91, 98},
			{-60, -20, 2, 3, 3, 11, 22, 31, 40, 40, 41, 48, 57, 57, 62},
			{
				-30, -12, -8, -9, 20, 23, 23, 35, 38, 53, 64, 65, 65, 66, 67,
				67, 72, 72, 77, 79, 93, 108, 108, 108, 110, 114, 114, 116,
			},
		}
	}

	if b.Bitboards[WN].Test(sq) {
		return bonus[0][Mobility(b, sq)]
	}

	if b.Bitboards[WB].Test(sq) {
		return bonus[1][Mobility(b, sq)]
	}

	if b.Bitboards[WR].Test(sq) {
		return bonus[2][Mobility(b, sq)]
	}

	if b.Bitboards[WQ].Test(sq) {
		return bonus[3][Mobility(b, sq)]
	}

	return 0
}

// Mobility is the number of attacked squares in the Mobility area. For queens squares
// defended by opponent knight, bishop or rook are ignored. For minor pieces squares
// occupied by our queen are ignored
func Mobility(b *board.Board, sq int) int {
	if !b.Bitboards[WN].Test(sq) && !b.Bitboards[WB].Test(sq) && !b.Bitboards[WR].Test(sq) &&
		!b.Bitboards[WQ].Test(sq) {
		return 0
	}

	score := 0

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if !MobilityArea(b, y*8+x) {
				continue
			}

			if b.Bitboards[WN].Test(sq) && evalhelpers.KnightAttack(b, y*8+x, sq) > 0 &&
				b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WB].Test(sq) && evalhelpers.BishopXrayAttack(b, y*8+x, sq) > 0 &&
				b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WR].Test(sq) && evalhelpers.RookXrayAttack(b, y*8+x, sq) > 0 {
				score++
			}
			if b.Bitboards[WQ].Test(sq) && evalhelpers.QueenAttack(b, y*8+x, sq) > 0 {
				score++
			}
		}
	}
	return score
}

// MobilityArea  do not include in mobility area squares protected by enemy pawns,
// or occupied by our blocked pawns or king. Pawns blocked or on ranks 2 and 3
// will be excluded from the mobility area. Also excludes blockers for king from
// mobility area - blockers for king can't really move until king moves (in most cases)
// so logic behind it is the same as behind excluding king square from mobility area.
func MobilityArea(b *board.Board, sq int) bool {
	if b.Bitboards[WK].Test(sq) {
		return false
	}

	if b.Bitboards[WQ].Test(sq) {
		return false
	}

	rank := sq / 8
	file := sq % 8

	if b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 {
		return false
	}

	if b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 {
		return false
	}

	if b.Bitboards[WP].Test(sq) && rank < 4 || b.Occupancies[color.BOTH].Test((rank-1)*8+file) {
		return false
	}

	mirror := b.Mirror()

	if BlockersForKing(mirror, (7-rank)*8+file) > 0 {
		return false
	}

	return true
}

// BlockersForKing returns if a particular piece on a particular square is a blocker
// for the king for a pin
func BlockersForKing(b *board.Board, sq int) int {
	mirror := b.Mirror()
	rank := sq / 8
	if evalhelpers.PinnedDirection(mirror, (7-rank)*8+(sq%8)) > 0 {
		return 1
	}

	return 0
}

// ThreatsMg returns the bonuses for middlegame threats
func ThreatsMg(b *board.Board) int {
	score := 0

	score += 69 * Hanging(b)
	if KingThreat(b) {
		score += 24
	}

	score += 48 * PawnPushThreat(b)
	score += 173 * ThreatSafePawn(b)
	score += 60 * SliderOnQueen(b)
	score += 16 * KnightOnQueen(b)
	score += 7 * Restricted(b)
	score += 14 * WeakQueenProtection(b)

	for sq := A8; sq <= H1; sq++ {
		score += []int{0, 5, 57, 77, 88, 79, 0}[MinorThreat(b, sq)]
		score += []int{0, 3, 37, 42, 0, 58, 0}[RookThreat(b, sq)]
	}

	return score
}

// Hanging returns weak enemies not defended by opponent or non-pawn weak
// enemies attacked twice
func Hanging(b *board.Board) int {
	weakEnemies := 0
	// Loop Through all enemies and check if there are any
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		if WeakEnemies(b, sq) == 0 {
			continue
		}

		if !b.Bitboards[BP].Test(sq) && Attack(b, sq) > 1 {
			weakEnemies++
		}

		mirror := b.Mirror()
		if Attack(mirror, (7-(sq/8))*8+(sq%8)) == 0 {
			weakEnemies++
		}
	}
	return weakEnemies
}

// WeakEnemies returns enemies not defended by a pawn and under our attack
func WeakEnemies(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 {
		return 0
	}

	if b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 {
		return 0
	}

	if Attack(b, sq) == 0 {
		return 0
	}

	mirror := b.Mirror()
	if Attack(b, sq) <= 1 && Attack(mirror, (7-rank)*8+file) > 1 {
		return 0
	}

	return 1
}

// Attack counts the number of attacks on square by all pieces. For bishop and rook
// x-ray attacks are included. For pawns pins or en-passant are ignored.
func Attack(b *board.Board, sq int) int {
	score := 0
	score += evalhelpers.PawnAttack(b, sq)
	score += evalhelpers.KingAttack(b, sq)
	score += evalhelpers.KnightAttack(b, sq, -1)
	score += evalhelpers.BishopXrayAttack(b, sq, -1)
	score += evalhelpers.RookXrayAttack(b, sq, -1)
	score += evalhelpers.QueenAttack(b, sq, -1)

	return score
}

// KingThreat returns if the king is in threat
func KingThreat(b *board.Board) bool {
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		if WeakEnemies(b, sq) == 0 {
			return false
		}

		if evalhelpers.KingAttack(b, sq) == 0 {
			return false
		}
	}
	return true
}

// PawnPushThreat returns the number of pawns that can be safely pushed
// and attack and enemy piece
func PawnPushThreat(b *board.Board) int {
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		rank := sq / 8
		file := sq % 8

		for ix := -1; ix <= 1; ix += 2 {
			if b.Bitboards[WP].Test((rank+2)*8+file+ix) &&
				file+ix >= 0 && file+ix <= 7 && rank+2 <= 7 &&
				b.Occupancies[color.BOTH].Test((rank+1)*8+file+ix) &&
				!(b.Bitboards[BP].Test(rank*8+file+ix-1) && file+ix-1 >= 0 && file+ix-1 <= 7) &&
				!(b.Bitboards[BP].Test(rank*8+file+ix+1) && file+ix+1 >= 0 && file+ix+1 <= 7) &&
				(Attack(b, (rank+1)*8+file+ix) > 0 || Attack(b.Mirror(), (6-rank)*8+file+ix) == 0) {
				return 1
			}

			if file == 3 && b.Bitboards[WP].Test((rank+3)*8+file+ix) &&
				file+ix >= 0 && file+ix <= 7 && rank+3 <= 7 &&
				b.Occupancies[color.BOTH].Test((rank+1)*8+file+ix) &&
				b.Occupancies[color.BOTH].Test((rank+2)*8+file+ix) &&
				!(b.Bitboards[BP].Test(rank*8+file+ix-1) && file+ix-1 >= 0 && file+ix-1 <= 7) &&
				!(b.Bitboards[BP].Test(rank*8+file+ix+1) && file+ix+1 >= 0 && file+ix+1 <= 7) &&
				(Attack(b, (rank+1)*8+file+ix) > 0 || Attack(b.Mirror(), (6-rank)*8+file+ix) == 0) {
				return 1
			}
		}

	}
	return 0
}

// ThreatSafePawn returns the non-pawn enemies attacked by a safe pawn
func ThreatSafePawn(b *board.Board) int {
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		rank := sq / 8
		file := sq % 8

		if evalhelpers.PawnAttack(b, sq) == 0 {
			return 0
		}

		if (SafePawn(b, (rank+1)*8+file-1) && file > 0 && rank < 7) ||
			(SafePawn(b, (rank+1)*8+file+1) && file < 7 && rank < 7) {
			return 1
		}
	}
	return 0
}

// SafePawn returns whether or not our pawn is not attacked or is defended
func SafePawn(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	if b.Bitboards[WP].Test(sq) {
		return false
	}

	if Attack(b, sq) > 0 {
		return true
	}

	if Attack(b.Mirror(), (7-rank)*8+file) == 0 {
		return true
	}

	return false
}

// SliderOnQueen adds a bonus for safe slider attack threats on opponent queen
func SliderOnQueen(b *board.Board) int {
	mirror := b.Mirror()

	if QueenCount(mirror) != 1 {
		return 0
	}

	rank, file := 0, 0
	sq := 0

	bb := b.Bitboards[WB] | b.Bitboards[WR] | b.Bitboards[WQ]
	for bb != 0 {
		sq = bb.FirstOne()

		rank = sq / 8
		file = sq % 8

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			return 0
		}

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			return 0
		}

		if Attack(b, sq) <= 1 {
			return 0
		}

		if !MobilityArea(b, sq) {
			return 0
		}

		diagonal := evalhelpers.QueenAttackDiagonal(mirror, (7-rank)*8+file, -1)
		v := 1
		if QueenCount(b) == 0 {
			v = 2
		}

		if diagonal != 0 && evalhelpers.BishopXrayAttack(b, sq, -1) != 0 {
			return v
		}

		if diagonal == 0 && evalhelpers.RookXrayAttack(b, sq, -1) != 0 &&
			evalhelpers.QueenAttack(mirror, (7-rank)*8+file, -1) != 0 {
			return v
		}
	}

	return 0
}

// QueenCount returns the number of white queens
func QueenCount(b *board.Board) int {
	return b.Bitboards[WQ].Count()
}

// KnightOnQueen returns a bonus for safe knight attack threaths on
// opponent queen
func KnightOnQueen(b *board.Board) int {
	mirror := b.Mirror()

	blackQueen := b.Bitboards[BQ]
	blackQueenSq := blackQueen.FirstOne()

	blackQueenRank := blackQueenSq / 8
	blackQueenFile := blackQueenSq % 8

	if QueenCount(mirror) != 1 {
		return 0
	}

	sq, rank, file := 0, 0, 0

	bb := b.Bitboards[WN]
	for bb != 0 {

		sq = bb.FirstOne()

		rank = sq / 8
		file = sq % 8

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			return 0
		}

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			return 0
		}

		if Attack(b, sq) <= 1 && Attack(mirror, (7-rank)*8+file) > 1 {
			return 0
		}

		if !MobilityArea(b, sq) {
			return 0
		}

		if evalhelpers.KnightAttack(b, sq, -1) == 0 {
			return 0
		}

		v := 1
		if QueenCount(b) == 0 {
			v = 2
		}

		if abs(blackQueenFile-file) == 2 && abs(blackQueenRank-rank) == 1 {
			return v
		}

		if abs(blackQueenFile-file) == 1 && abs(blackQueenRank-rank) == 2 {
			return v
		}
	}

	return 0
}

// Restricted returns a bonus for restricing their pieces moves
func Restricted(b *board.Board) int {
	restricted := 0

	for sq := A8; sq <= H8; sq++ {
		if Attack(b, sq) == 0 {
			continue
		}

		rank := sq / 8
		file := sq % 8

		mirror := b.Mirror()

		if Attack(mirror, (7-rank)*8+file) == 0 {
			continue
		}

		if evalhelpers.PawnAttack(mirror, (7-rank)*8+file) > 0 {
			continue
		}

		if Attack(mirror, (7-rank)*8+file) > 1 && Attack(b, sq) == 1 {
			continue
		}

		restricted++
	}

	return restricted
}

// WeakQueenProtection adds an additional bonus if weak piece is only
// protected by a queen
func WeakQueenProtection(b *board.Board) int {
	weakPieces := 0

	for sq := A8; sq <= H8; sq++ {
		if WeakEnemies(b, sq) == 0 {
			continue
		}

		if evalhelpers.QueenAttack(b.Mirror(), (7-(sq/8)*8+(sq%8)), -1) == 0 {
			continue
		}
		weakPieces++
	}
	return weakPieces
}

// MinorThreat returns the threat type for knight and bishop attacking pieces
func MinorThreat(b *board.Board, sq int) int {
	if !b.Occupancies[color.BLACK].Test(sq) {
		return 0
	}

	if evalhelpers.KnightAttack(b, sq, -1) == 0 && evalhelpers.BishopXrayAttack(b, sq, -1) == 0 {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if (b.Bitboards[BP].Test(sq) ||
		!((b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 && rank > 0) ||
			(b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 && rank > 0) ||
			(Attack(b, sq) <= 1 && Attack(b.Mirror(), (7-rank)*8+file) > 1))) &&
		WeakEnemies(b, sq) == 0 {

		return 0
	}

	if b.Bitboards[BP].Test(sq) {
		return 1
	}

	if b.Bitboards[BN].Test(sq) {
		return 2
	}

	if b.Bitboards[BB].Test(sq) {
		return 3
	}

	if b.Bitboards[BR].Test(sq) {
		return 4
	}

	if b.Bitboards[BQ].Test(sq) {
		return 5
	}

	return 6
}

// RookThreat return the threat type for attacked by rook pieces
func RookThreat(b *board.Board, sq int) int {
	if !b.Occupancies[color.BLACK].Test(sq) {
		return 0
	}

	if WeakEnemies(b, sq) == 0 {
		return 0
	}

	if evalhelpers.RookXrayAttack(b, sq, -1) == 0 {
		return 0
	}

	if b.Bitboards[BP].Test(sq) {
		return 1
	}

	if b.Bitboards[BN].Test(sq) {
		return 2
	}

	if b.Bitboards[BB].Test(sq) {
		return 3
	}

	if b.Bitboards[BR].Test(sq) {
		return 4
	}

	if b.Bitboards[BQ].Test(sq) {
		return 5
	}

	return 6
}

// PassedMg return middlegame bonuses for passed pawn. Scale
// down bonus for candidate passers which need more than one pawn
// push to become passed, or have a pawn in from of them
func PassedMg(b *board.Board) int {
	finalScore := 0

	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()

		if PassedLeverable(b, sq) == 0 {
			continue
		}

		score := 0

		score += []int{0, 10, 17, 15, 62, 168, 276}[PassedRank(b, sq)]
		score += PassedBlock(b, sq)
		score -= PassedFile(b, sq)

		finalScore += score
	}

	return finalScore
}

// PassedLeverable returns candidate passers without candidate passers w/o
// feasible lever
func PassedLeverable(b *board.Board, sq int) int {
	return 0
}

func PassedRank(b *board.Board, sq int) int {
	return 0
}

func PassedFile(b *board.Board, sq int) int {
	return 0
}

func PassedBlock(b *board.Board, sq int) int {
	return 0
}

// Space computes the space evaluation for a given side. The Space are
// bonus is multiplied by a weight: number of our pieces minus two times
// number of open files. The aim is to improve play on game opening
func Space(b *board.Board) int {
	if nonPawnMaterial(b, color.WHITE)+nonPawnMaterial(b.Mirror(), color.WHITE) < 12222 {
		return 0
	}

	score := 0
	pieceCount, blockedCount := 0, 0

	pieceBB := b.Occupancies[color.BOTH]
	for pieceBB != 0 {
		sq := pieceBB.FirstOne()

		if b.Occupancies[color.WHITE].Test(sq) {
			pieceCount++
		}

		rank := sq / 8
		file := sq % 8

		if b.Bitboards[WP].Test(sq) &&
			((b.Bitboards[BP].Test((rank-1)*8+file) && rank > 0) ||
				(b.Bitboards[BP].Test((rank-2)*8+file-1) && rank > 1 && file > 0) &&
					(b.Bitboards[BP].Test((rank-2)*8+file+1) && rank > 1 && file < 7)) {
			blockedCount++
		}

		if b.Bitboards[BP].Test(sq) &&
			((b.Bitboards[BP].Test((rank+1)*8+file) && rank < 7) ||
				(b.Bitboards[BP].Test((rank+2)*8+file-1) && rank < 6 && file > 0) &&
					(b.Bitboards[BP].Test((rank+2)*8+file+1) && rank < 6 && file < 7)) {
			blockedCount++
		}

		weight := pieceCount - 3 + min(blockedCount, 9)
		score += ((SpaceArea(b, sq) * weight * weight / 16) << 0)
	}
	return score
}

// SpaceArea returns the number of safe squares available for minor pieces
// on the central four files on ranks 2 to 4. Safe squares one, two or three
// squares behind a friendly pawn are counted twice
func SpaceArea(b *board.Board, sq int) int {
	score := 0

	rank := sq / 8
	file := sq % 8

	if ((8-rank) >= 2 && (8-rank) <= 4 && (8-file) >= 3 && (8-file) <= 6) &&
		b.Bitboards[WP].Test(sq) &&
		(b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0) &&
		(b.Bitboards[BP].Test((rank+1)*8+file-1) && rank < 7 && file > 0) {
		score++

		if ((b.Bitboards[WP].Test((rank-1)*8+file) && rank > 0) ||
			(b.Bitboards[WP].Test((rank-2)*8+file) && rank > 1) ||
			(b.Bitboards[WP].Test((rank-3)*8+file) && rank > 2)) &&
			Attack(b.Mirror(), (7-rank)*8+file) == 0 {
			score++
		}
	}

	return score
}

// KingMg assigns middlegame bonuses and penalties for attacks on enemy king
func KingMg(b *board.Board) int {
	score := 0

	kd := KingDanger(b)
	score -= ShelterStrength(b)
	score += ShelterStorm(b)
	score += (kd * kd / 4096) << 0
	score += 8 * FlankAttack(b)
	score += 17 * PawnlessFlank(b)
	return score
}

// KingDanger returns the danger that the king is in. The initial value
// is based on the number and types of the enemy's attacking pieces, the
// number of attacked and undefended squares around our king and the
// quality of the pawn shelter
func KingDanger(b *board.Board) int {
	score := 0
	kingBB := b.Bitboards[WK]
	for kingBB != 0 {
		sq := kingBB.FirstOne()

		count := KingAttackersCount(b, sq)
		weight := KingAttackersWeight(b, sq)
		kingAttacks := evalhelpers.KingAttack(b, sq)
		weak := WeakBonus(b, sq)
		unsafeChecks := UnsafeChecks(b, sq)
		blockersForKing := BlockersForKing(b, sq)
		kingFlankAttack := FlankAttack(b)
		kingFlankDefense := FlankDefense(b)
		noQueen := 1
		if QueenCount(b) > 0 {
			noQueen = 0
		}

		knightBonusFactor := 0
		if KnightDefender(b.Mirror()) {
			knightBonusFactor = 1
		}

		score = count*weight +
			69*kingAttacks +
			185*weak -
			100*knightBonusFactor +
			148*unsafeChecks +
			98*blockersForKing -
			4*kingFlankDefense +
			((3 * kingFlankAttack * kingFlankAttack / 8) << 0) -
			873*noQueen -
			((6 * (ShelterStrength(b) - ShelterStorm(b)) / 8) << 0) +
			MobilityMg(b) - MobilityMg(b.Mirror()) +
			37 +
			((772 * int(min(SafeCheck(b, nil, 3), 1.45))) << 0) +
			((1084 * int(min(SafeCheck(b, nil, 2), 1.75))) << 0) +
			((645 * int(min(SafeCheck(b, nil, 1), 1.50))) << 0) +
			((792 * int(min(SafeCheck(b, nil, 0), 1.62))) << 0)
	}

	if score > 100 {
		return score
	}
	return 0
}

func WeakBonus(b *board.Board, sq int) int {
	return 0
}

func KingAttackersWeight(b *board.Board, sq int) int {
	return 0
}

func UnsafeChecks(b *board.Board, sq int) int {
	return 0
}

func FlankDefense(b *board.Board) int {
	return 0
}

func KnightDefender(b *board.Board) bool {
	return false
}

func SafeCheck(b *board.Board, sq *int, factor int) float32 {
	return 0.0
}

func ShelterStrength(b *board.Board) int {
	return 0
}

func ShelterStorm(b *board.Board) int {
	return 0
}

func FlankAttack(b *board.Board) int {
	return 0
}

func PawnlessFlank(b *board.Board) int {
	return 0
}

func WinnableTotalMg(b *board.Board, score int) int {
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
