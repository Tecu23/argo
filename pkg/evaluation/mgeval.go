package evaluation

import (
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

func MobilityMg(b *board.Board) int {
	return 0
}

func MobilityBonus(b *board.Board, sq int) int {
	return 0
}

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

			if b.Bitboards[WN].Test(sq) && KnightAttack(b, y*8+x, sq) > 0 &&
				b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WB].Test(sq) && BishopXrayAttack(b, y*8+x, sq) > 0 &&
				b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WR].Test(sq) && RookXrayAttack(b, y*8+x, sq) > 0 {
				score++
			}
			if b.Bitboards[WQ].Test(sq) && QueenAttack(b, y*8+x, sq) > 0 {
				score++
			}
		}
	}
	return score
}

func QueenAttack(b *board.Board, sq int, sq2 int) int {
	score := 0

	rank := sq / 8
	file := sq % 8

	rank2 := sq2 / 8
	file2 := sq2 % 8

	factor := 0

	for i := 0; i < 8; i++ {
		factor = 0

		if i > 3 {
			factor = 1
		}

		ix := (i+factor)%3 - 1
		iy := (((i + factor) / 3) << 0) - 1

		for d := 1; d < 8; d++ {
			if b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				(file2 == file+ix && rank2 == rank+iy) {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if !b.Occupancies[color.BOTH].Test((rank+d*iy)*8 + file + d*ix) {
				break
			}
		}
	}

	return score
}

func RookXrayAttack(b *board.Board, sq int, sq2 int) int {
	score := 0

	rank := sq / 8
	file := sq % 8

	rank2 := sq2 / 8
	file2 := sq2 % 8

	for i := 0; i < 4; i++ {
		ix := 0
		iy := 0
		if i == 0 {
			ix = -1
		} else if i == 1 {
			ix = 1
		}

		if i == 2 {
			iy = -1
		} else if i == 3 {
			iy = 1
		}

		for d := 1; d < 8; d++ {
			if b.Bitboards[WR].Test((rank+d*iy)*8+file+d*ix) &&
				(file2 == file+ix && rank2 == rank+iy) {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if !b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WR].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) {
				break
			}
		}
	}

	return score
}

func BishopXrayAttack(b *board.Board, sq int, sq2 int) int {
	score := 0
	factor1, factor2 := 0, 0

	rank := sq / 8
	file := sq % 8

	rank2 := sq2 / 8
	file2 := sq2 % 8

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
			if b.Bitboards[WB].Test((rank+d*iy)*8+file+d*ix) &&
				(file+d*ix >= 0) && (file+d*ix <= 7) &&
				file2 == file+d*ix &&
				rank2 == rank+d*iy {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if ((rank+d*iy)*8 + file + d*ix) == G6 {
				b.Occupancies[color.BOTH].PrintBitboard()
			}

			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WB].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) {
				break
			}
		}
	}

	return score
}

func KnightAttack(b *board.Board, sq int, sq2 int) int {
	score := 0
	factor1, factor2, factor3 := 0, 0, 0

	rank := sq / 8
	file := sq % 8

	rank2 := sq2 / 8
	file2 := sq2 % 8

	for i := 0; i < 8; i++ {
		factor1, factor2, factor3 = 0, 0, 0
		if i > 3 {
			factor1 = 1
		}
		if i%4 > 1 {
			factor2 = 1
		}

		if i%2 == 0 {
			factor3 = 1
		}

		ix := (factor1 + 1) * (factor2*2 - 1)
		iy := (2 - factor1) * (factor3*2 - 1)

		if b.Bitboards[WN].Test((rank+iy)*8+file+ix) && (file2 == file+ix && rank2 == rank+iy) &&
			Pinned(b, (rank+iy)*8+file+ix) == 0 {
			score++
		}
	}

	return score
}

func Pinned(b *board.Board, sq int) int {
	return PinnedDirection(b, sq)
}

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

func BlockersForKing(b *board.Board, sq int) int {
	mirror := b.Mirror()
	rank := sq / 8
	if PinnedDirection(mirror, (7-rank)*8+(sq%8)) > 0 {
		return 1
	}

	return 0
}

func PinnedDirection(b *board.Board, sq int) int {
	c := 1

	rank := sq / 8
	file := sq % 8

	for i := 0; i < 8; i++ {
		factor := 0
		if i > 3 {
			factor = 1
		}

		ix := (i+factor)%3 - 1
		iy := (((i + factor) / 3) << 0) - 1

		king := false

		for d := 1; d < 8; d++ {
			if b.Bitboards[BK].Test((rank+d*iy)*8 + file + d*ix) {
				king = true
			}
			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8 + file + d*ix) {
				break
			}
		}

		if king {
			for d := 1; d < 8; d++ {
				if b.Bitboards[BQ].Test((rank-d*iy)*8+file-d*ix) ||
					(b.Bitboards[BB].Test((rank-d*iy)*8+file-d*ix) && ix*iy != 0) ||
					(b.Bitboards[BR].Test((rank-d*iy)*8+file-d*ix) && ix*iy == 0) {
					return abs(ix+iy*3) * c
				}

				if b.Occupancies[color.BOTH].Test((rank-d*iy)*8 + file - d*ix) {
					break
				}
			}
		}
	}

	return 0
}

func ThreatsMg(b *board.Board) int {
	return 0
}

func PassedMg(b *board.Board) int {
	return 0
}

func Space(b *board.Board) int {
	return 0
}

func KingMg(b *board.Board) int {
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
