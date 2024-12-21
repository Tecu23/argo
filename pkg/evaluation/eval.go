package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func MiddleGameEvaluation(b *board.Board, noWinnable bool) int {
	score := 0
	mirror := b.Mirror()

	score += PieceValueMg(b) - PieceValueMg(mirror)
	score += PsqtMg(b) - PsqtMg(mirror)
	score += ImbalanceTotal(b, mirror)
	score += PawnsMg(b) - PawnsMg(mirror)
	score += PiecesMg(b) - PiecesMg(mirror)
	score += MobilityMg(b) - MobilityMg(mirror)
	score += ThreatsMg(b) - ThreatsMg(mirror)
	score += PassedMg(b) - PassedMg(mirror)
	score += Space(b) - Space(mirror)
	score += KingMg(b) - KingMg(mirror)

	if !noWinnable {
		score += WinnableTotalMg(b, score)
	}

	return score
}

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
	return 0
}

func PawnsMg(b *board.Board) int {
	return 0
}

func PiecesMg(b *board.Board) int {
	return 0
}

func MobilityMg(b *board.Board) int {
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
