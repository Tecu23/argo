package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PieceValueEg returns the material evaluation for an endgame position
func PieceValueEg(b *board.Board) int {
	PawnBonus := 206
	KnightBonus := 854
	BishopBonus := 915
	RookBonus := 1380
	QueenBonus := 2682

	pawnScore := b.Bitboards[WP].Count() * PawnBonus
	knightScore := b.Bitboards[WN].Count() * KnightBonus
	bishopScore := b.Bitboards[WB].Count() * BishopBonus
	rookScore := b.Bitboards[WR].Count() * RookBonus
	queenScore := b.Bitboards[WQ].Count() * QueenBonus

	return pawnScore + knightScore + bishopScore + rookScore + queenScore
}

// PsqtEg returns the table bonuses endgame evaluation
func PsqtEg(b *board.Board) int {
	bonus := [][][]int{
		{
			{-96, -65, -49, -21},
			{-67, -54, -18, 8},
			{-40, -27, -8, 29},
			{-35, -2, 13, 28},
			{-45, -16, 9, 39},
			{-51, -44, -16, 17},
			{-69, -50, -51, 12},
			{-100, -88, -56, -17},
		},
		{
			{-57, -30, -37, -12},
			{-37, -13, -17, 1},
			{-16, -1, -2, 10},
			{-20, -6, 0, 17},
			{-17, -1, -14, 15},
			{-30, 6, 4, 6},
			{-31, -20, -1, 1},
			{-46, -42, -37, -24},
		},
		{
			{-9, -13, -10, -9},
			{-12, -9, -1, -2},
			{6, -8, -2, -6},
			{-6, 1, -9, 7},
			{-5, 8, 7, -6},
			{6, 1, -7, 10},
			{4, 5, 20, -5},
			{18, 0, 19, 13},
		},
		{
			{-69, -57, -47, -26},
			{-55, -31, -22, -4},
			{-39, -18, -9, 3},
			{-23, -3, 13, 24},
			{-29, -6, 9, 21},
			{-38, -18, -12, 1},
			{-50, -27, -24, -8},
			{-75, -52, -43, -36},
		},
		{
			{1, 45, 85, 76},
			{53, 100, 133, 135},
			{88, 130, 169, 175},
			{103, 156, 172, 172},
			{96, 166, 199, 199},
			{92, 172, 184, 191},
			{47, 121, 116, 131},
			{11, 59, 73, 78},
		},
	}

	pBonus := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{-10, -6, 10, 0, 14, 7, -5, -19},
		{-10, -10, -10, 4, 4, 3, -6, -4},
		{6, -2, -8, -4, -13, -12, -10, -9},
		{10, 5, 4, -5, -5, -5, 14, 9},
		{28, 20, 21, 28, 30, 7, 6, 13},
		{0, -11, 12, 21, 25, 19, 4, 7},
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

// PawnsEg returns endgame evaluation for pawns
func PawnsEg(b *board.Board) int {
	score := 0

	pawnsBB := b.Bitboards[WP]
	for pawnsBB != 0 {
		sq := pawnsBB.FirstOne()

		if DoubleIsolated(b, sq) {
			score -= 56
		} else if Isolated(b, sq) {
			score -= 15
		} else if Backward(b, sq) {
			score -= 24
		}

		if Doubled(b, sq) {
			score -= 56
		}

		if Connected(b, sq) {
			score += ConnectedBonus(b, sq) * ((8 - (sq / 8) - 3) / 4)
		}

		score -= 27 * WeakUnopposedPawn(b, sq)
		score += []int{0, -4, 4}[Blocked(b, sq)]
	}
	return score
}

// PiecesEg returns endgame bonuses and penalties to the pieces
// of a given color and type. Mobility not included here
func PiecesEg(b *board.Board) int {
	score := 0

	// scores for Knight, Bishop, Rook, Queen
	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq := knightBB.FirstOne()

		score += []int{0, 22, 36, 23, 36}[OutpostTotal(b, sq)]
		score += 3 * MinorBehindPawn(b, sq)
		score -= 9 * KingProtector(b, sq)

	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()

		score += []int{0, 22, 36, 23, 36}[OutpostTotal(b, sq)]
		score += 3 * MinorBehindPawn(b, sq)
		score -= 7 * BishopPawns(b, sq)
		score -= 5 * BishopXrayPawns(b, sq)
		score -= 9 * KingProtector(b, sq)

	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq := rookBB.FirstOne()

		score += 11 * RookOnQueenFile(b, sq)
		score += []int{0, 7, 29}[RookOnFile(b, sq)]

		factor := 2
		if uint(b.Castlings)&ShortW != 0 || uint(b.Castlings)&LongW != 0 {
			factor = 1
		}

		score -= TrappedRook(b, sq) * 13 * factor
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq := queenBB.FirstOne()

		score -= 15 * WeakQueen(b, sq)
		score += 14 * QueenInfiltration(b, sq)
	}

	return score
}

// MobilityEg returns a mobility bonus for endgame
func MobilityEg(b *board.Board) int {
	score := 0
	sq := 0

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq = knightBB.FirstOne()
		score += MobilityBonus(b, sq, false)
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq = bishopBB.FirstOne()
		score += MobilityBonus(b, sq, false)
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq = rookBB.FirstOne()
		score += MobilityBonus(b, sq, false)
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq = queenBB.FirstOne()
		score += MobilityBonus(b, sq, false)
	}
	return score
}

// ThreatsEg is a endgame threats bonus
func ThreatsEg(b *board.Board) int {
	score := 0

	score += 36 * Hanging(b)

	if KingThreat(b) {
		score += 89
	}

	score += 39 * PawnPushThreat(b)
	score += 94 * ThreatSafePawn(b)
	score += 18 * SliderOnQueen(b)
	score += 11 * KnightOnQueen(b)
	score += 7 * Restricted(b)

	for sq := A8; sq <= H1; sq++ {
		score += []int{0, 32, 41, 56, 119, 161, 0}[MinorThreat(b, sq)]
		score += []int{0, 46, 68, 60, 38, 41, 0}[RookThreat(b, sq)]
	}

	return score
}

// PassedEg endgame bonuses for passed pawns. Scale down bonus for
// candidate passers which need more than one pawn push to become
// passed, or have a pawn in front of them.
func PassedEg(b *board.Board) int {
	finalScore := 0

	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()

		if PassedLeverable(b, sq) == 0 {
			continue
		}

		score := 0
		score += KingProximity(b, sq)
		score += []int{0, 28, 33, 41, 72, 177, 260}[PassedRank(b, sq)]
		score += PassedBlock(b, sq)
		score -= 8 * PassedFile(b, sq)

		finalScore += score
	}

	return finalScore
}

// KingProximity is an endgame bonus based on the king's proximity.
// If block square is not the queening square then consider also a second push
func KingProximity(b *board.Board, sq int) int {
	if PassedLeverable(b, sq) == 0 {
		return 0
	}
	rank := sq / 8
	file := sq % 8

	r := (8 - rank) - 1
	w := 0
	if r > 2 {
		w = 5*r - 13
	}
	score := 0

	if w <= 0 {
		return 0
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if b.Bitboards[BK].Test(y*8 + x) {
				score += ((min(max(abs(y-rank+1), abs(x-file)), 5) * 19 / 4) << 0) * w
			}

			if b.Bitboards[WK].Test(y*8 + x) {
				score -= min(max(abs(y-rank+1), abs(x-file)), 5) * 2 * w

				if rank > 1 {
					score -= min(max(abs(y-rank+2), abs(x-file)), 5) * w
				}
			}
		}
	}
	return score
}

// KingEg assigns endgame bonuses and penalties for attacks on enemy king
func KingEg(b *board.Board) int {
	score := 0

	score -= 16 * KingPawnDistance(b)
	score += EndgameShelter(b, -1)

	if PawnlessFlank(b) {
		score += 95
	}
	kd := KingDanger(b)
	score += kd / 16
	return score
}

// KingPawnDistance is the minimal distance of our king to our pawns
func KingPawnDistance(b *board.Board) int {
	v := 6
	kx, ky := 0, 0
	// px, py := 0, 0

	kingBB := b.Bitboards[WK]
	kingSq := kingBB.FirstOne()

	kx = kingSq % 8
	ky = kingSq / 8

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			dist := max(abs(x-kx), abs(y-ky))

			if b.Bitboards[WP].Test(y*8+x) && dist < v {
				// px = x
				// py = y
				v = dist
			}
		}
	}

	return v
}

// EndgameShelter adds an endgame component to the blockedstorm penalty
// so that the penalty applies more uniformly through the game.
func EndgameShelter(b *board.Board, square int) int {
	w := 0
	e := 0
	s := 1024

	rank, file := 0, 0

	for sq := A8; sq <= H1; sq++ {
		rank = sq / 8
		file = sq % 8
		if b.Bitboards[BK].Test(sq) ||
			uint(b.Castlings)&ShortB != 0 && file == 6 && rank == 0 ||
			uint(b.Castlings)&LongB != 0 && file == 2 && rank == 0 {

			w1 := StrengthSquare(b, sq)
			s1 := StormSquare(b, sq, false)
			e1 := StormSquare(b, sq, true)

			if s1-w1 < s-w {
				w = w1
				s = s1
				e = e1
			}
		}
	}

	if square == -1 {
		return e
	}

	return 0
}

// WinnableTotalEg returns end game winnable
func WinnableTotalEg(b *board.Board, score int) int {
	if score == -1 {
		score = EndGameEvaluation(b, true)
	}

	factor := 0
	if score > 0 {
		factor = 1
	} else if score < 0 {
		factor = -1
	}

	return factor * max(Winnable(b), -abs(score))
}
