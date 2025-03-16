// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// KingEvaluation returns bonuses and penalties for attacks on enemy king
func (e *Evaluator) KingEvaluation(b *board.Board) (mg, eg int) {
	mg, eg = 0, 0

	mg -= shelterStrength(b, -1)
	mg += shelterStorm(b, -1)

	whiteMobilityMG, _ := e.MobilityEvaluation(b)
	blackMobilityMG, _ := e.MobilityEvaluation(b.Mirror())
	kd := kingDanger(b, whiteMobilityMG, blackMobilityMG)
	mg += (kd * kd / 4096) << 0
	eg += (kd / 16) << 0

	mg += 8 * flankAttack(b)

	eg -= 16 * kingPawnDistance(b)
	eg += endgameShelter(b, -1)

	if pawnlessFlank(b) {
		mg += 17
		eg += 95
	}

	return mg, eg
}

// kingDanger returns the danger that the king is in. The initial value
// is based on the number and types of the enemy's attacking pieces, the
// number of attacked and undefended squares around our king and the
// quality of the pawn shelter
func kingDanger(b *board.Board, wMobMG, bMobMG int) int {
	mirror := b.Mirror()
	score := 0

	count := 0
	weight := 0
	kingAttacksCount := 0
	weak := 0
	unsafeChecksCount := 0
	blockersForKingCount := 0

	whiteBB := b.Occupancies[color.WHITE]
	for whiteBB != 0 {
		sq := whiteBB.FirstOne()
		count += kingAttackersCount(b, sq)
		weight += kingAttackersWeight(b, sq)
		kingAttacksCount += kingAttacks(b, sq)
	}

	for sq := A8; sq <= H1; sq++ {
		unsafeChecksCount += unsafeChecks(b, sq)
		weak += weakBonus(b, sq)
	}

	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()
		blockersForKingCount += blockersForKing(b, mirror, sq)
	}

	kingFlankAttack := flankAttack(b)
	kingFlankDefense := flankDefense(b)

	noQueen := 1
	if queenCount(b) > 0 {
		noQueen = 0
	}

	knightBonusFactor := 0
	if knightDefender(b.Mirror()) > 0 {
		knightBonusFactor = 1
	}

	score = count*weight +
		69*kingAttacksCount +
		185*weak -
		100*knightBonusFactor +
		148*unsafeChecksCount +
		98*blockersForKingCount -
		4*kingFlankDefense +
		((3 * kingFlankAttack * kingFlankAttack / 8) << 0) -
		873*noQueen -
		((6 * (shelterStrength(b, -1) - shelterStorm(b, -1)) / 8) << 0) +
		(wMobMG - bMobMG) +
		37 +
		int((772 * min(safeCheck(b, -1, 3), 1.45))) +
		int((1084 * min(safeCheck(b, -1, 2), 1.75))) +
		int((645 * min(safeCheck(b, -1, 1), 1.50))) +
		int((792 * min(safeCheck(b, -1, 0), 1.62)))

	if score > 100 {
		return score
	}

	return 0
}

// weakBonus returns if the king has weak squares
func weakBonus(b *board.Board, sq int) int {
	if weakSquares(b, sq) == 0 {
		return 0
	}

	if kingRing(b, sq, false) == 0 {
		return 0
	}

	return 1
}

// weakSquares returns attacked squares defended at most once
// by our queen or king
func weakSquares(b *board.Board, sq int) int {
	if attack(b, sq) > 0 {
		mirror := b.Mirror()

		rank := sq / 8
		file := sq % 8

		attack := attack(mirror, (7-rank)*8+file)
		if attack >= 2 {
			return 0
		}

		if attack == 0 {
			return 1
		}

		if KingAttack(mirror, (7-rank)*8+file) > 0 ||
			QueenAttack(mirror, (7-rank)*8+file, -1) > 0 {
			return 1
		}
	}

	return 0
}

// kingAttackersWeight is the sum of the "weights" of the pieces of
// the given color which attack a square in the king ring of the enemy king
func kingAttackersWeight(b *board.Board, sq int) int {
	if kingAttackersCount(b, sq) > 0 {
		if b.Bitboards[WP].Test(sq) {
			return 0
		} else if b.Bitboards[WN].Test(sq) {
			return 81
		} else if b.Bitboards[WB].Test(sq) {
			return 52
		} else if b.Bitboards[WR].Test(sq) {
			return 44
		} else if b.Bitboards[WQ].Test(sq) {
			return 10
		}
	}
	return 0
}

// kingAttacks is the number of attacks by the given color to squares directly
// adjancent to the enemy king. Pieces which attack more than one square are
// counted multuple times. For instance, If there is a white knight on g5 and
// black's king is on g8, this white knight adds 2.
func kingAttacks(b *board.Board, sq int) int {
	if b.Bitboards[WP].Test(sq) || b.Bitboards[WK].Test(sq) {
		return 0
	}

	if kingAttackersCount(b, sq) == 0 {
		return 0
	}

	score := 0

	kingBB := b.Bitboards[BK]
	kingSq := kingBB.FirstOne()

	kingRank := kingSq / 8
	kingFile := kingSq % 8

	for x := kingFile - 1; x <= kingFile+1; x++ {
		for y := kingRank - 1; y <= kingRank+1; y++ {
			if x >= 0 && y >= 0 && x <= 7 && y <= 7 && (x != kingFile || y != kingRank) {
				score += KnightAttack(b, y*8+x, sq)
				score += BishopXrayAttack(b, y*8+x, sq)
				score += RookXrayAttack(b, y*8+x, sq)
				score += QueenAttack(b, y*8+x, sq)
			}
		}
	}
	return score
}

// unsafeChecks returns unsafe checks
func unsafeChecks(b *board.Board, sq int) int {
	if check(b, sq, 0) && safeCheck(b, -1, 0) == 0 {
		return 1
	}

	if check(b, sq, 1) && safeCheck(b, -1, 1) == 0 {
		return 1
	}

	if check(b, sq, 2) && safeCheck(b, -1, 2) == 0 {
		return 1
	}
	return 0
}

// check returns possible checks by knight, bishop, rook or queen. Defending
// queen is not considered as check blocker
func check(b *board.Board, sq int, t int) bool {
	rank := sq / 8
	file := sq % 8

	if (RookXrayAttack(b, sq, -1) > 0 && (t == -1 || t == 2 || t == 4)) ||
		(QueenAttack(b, sq, -1) > 0 && (t == -1 || t == 3)) {

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
				if b.Bitboards[BK].Test((rank+d*iy)*8+file+d*ix) &&
					file+d*ix >= 0 && file+d*ix <= 7 && rank+d*iy >= 0 && rank+d*iy <= 7 {
					return true
				}

				if (!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) &&
					b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix)) &&
					file+d*ix >= 0 && file+d*ix <= 7 && rank+d*iy >= 0 && rank+d*iy <= 7 {
					break
				}
			}
		}
	}

	if (BishopXrayAttack(b, sq, -1) > 0 && (t == -1 || t == 1 || t == 4)) ||
		(QueenAttack(b, sq, -1) > 0 && (t == -1 || t == 3)) {

		factor1, factor2 := 0, 0

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
				if b.Bitboards[BK].Test((rank+d*iy)*8+file+d*ix) &&
					file+d*ix >= 0 && file+d*ix <= 7 && rank+d*iy >= 0 && rank+d*iy <= 7 {
					return true
				}

				if (!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) &&
					b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix)) &&
					file+d*ix >= 0 && file+d*ix <= 7 && rank+d*iy >= 0 && rank+d*iy <= 7 {
					break
				}
			}
		}
	}

	if KnightAttack(b, sq, -1) > 0 && (t == -1 || t == 0 || t == 4) {
		if (b.Bitboards[BK].Test((rank+1)*8+file+2) && rank < 7 && file < 6) ||
			(b.Bitboards[BK].Test((rank-1)*8+file+2) && rank > 0 && file < 6) ||
			(b.Bitboards[BK].Test((rank+2)*8+file+1) && rank < 6 && file < 7) ||
			(b.Bitboards[BK].Test((rank-2)*8+file+1) && rank > 1 && file < 7) ||
			(b.Bitboards[BK].Test((rank+1)*8+file-2) && rank < 7 && file > 1) ||
			(b.Bitboards[BK].Test((rank-1)*8+file-2) && rank > 0 && file > 1) ||
			(b.Bitboards[BK].Test((rank+2)*8+file-1) && rank < 6 && file > 0) ||
			(b.Bitboards[BK].Test((rank-2)*8+file-1) && rank > 1 && file > 0) {
			return true
		}
	}

	return false
}

// safeCheck analyses the sage enemy”t give a rook check: we count them only if they are from squares from
// which we can't give a queen check, because queen checks are more valuable
func safeCheck(b *board.Board, sq int, t int) float32 {
	score := float32(0.0)
	if sq == -1 {
		for sq := A8; sq <= H1; sq++ {
			if b.Occupancies[color.WHITE].Test(sq) {
				continue
			}

			rank := sq / 8
			file := sq % 8

			if !check(b, sq, t) {
				continue
			}

			mirror := b.Mirror()

			if t == 3 && safeCheck(b, sq, 2) > 0 {
				continue
			}

			if t == 1 && safeCheck(b, sq, 3) > 0 {
				continue
			}

			if (attack(mirror, (7-rank)*8+file) == 0 ||
				(weakSquares(b, sq) > 0 && attack(b, sq) > 1)) &&
				(t != 3 || QueenAttack(mirror, (7-rank)*8+file, -1) == 0) {
				score += 1.0
			}
		}
	} else {
		rank := sq / 8
		file := sq % 8

		if !check(b, sq, t) {
			return 0.0
		}

		mirror := b.Mirror()

		if t == 3 && safeCheck(b, sq, 2) > 0 {
			return 0.0
		}

		if t == 1 && safeCheck(b, sq, 3) > 0 {
			return 0.0
		}

		if (attack(mirror, (7-rank)*8+file) == 0 ||
			(weakSquares(b, sq) > 0 && attack(b, sq) > 1)) &&
			(t != 3 || QueenAttack(mirror, (7-rank)*8+file, -1) == 0) {
			score += 1.0
		}
	}

	return score
}

// flankAttack finds the squares that opponent attacks in our king flank
// and the squares which they attack twice in that flank
func flankAttack(b *board.Board) int {
	score := 0
	for sq := A8; sq <= H1; sq++ {
		rank := sq / 8
		file := sq % 8
		if rank > 4 {
			continue
		}

		kingBB := b.Bitboards[BK]
		kingSq := kingBB.FirstOne()

		kingFile := kingSq % 8

		if kingFile == 0 && file > 2 {
			continue
		}

		if kingFile < 3 && file > 3 {
			continue
		}

		if kingFile >= 3 && kingFile < 5 && (file < 2 || file > 5) {
			continue
		}

		if kingFile >= 5 && file < 4 {
			continue
		}

		if kingFile == 7 && file < 5 {
			continue
		}

		a := attack(b, sq)
		if a == 0 {
			continue
		}

		if a > 1 {
			score += 2
		} else {
			score++
		}
	}

	return score
}

// flankDefense finds the squares that we defend in our king flank
func flankDefense(b *board.Board) int {
	score := 0
	for sq := A8; sq <= H1; sq++ {
		rank := sq / 8
		file := sq % 8
		if rank > 4 {
			continue
		}

		kingBB := b.Bitboards[BK]
		kingSq := kingBB.FirstOne()

		kingFile := kingSq % 8

		if kingFile == 0 && file > 2 {
			continue
		}

		if kingFile < 3 && file > 3 {
			continue
		}

		if kingFile >= 3 && kingFile < 5 && (file < 2 || file > 5) {
			continue
		}

		if kingFile >= 5 && file < 4 {
			continue
		}

		if kingFile == 7 && file < 5 {
			continue
		}

		a := attack(b.Mirror(), (7-rank)*8+file)
		if a > 0 {
			score++
		}
	}

	return score
}

// knightDefender returns the squares defended by knight near our king
func knightDefender(b *board.Board) int {
	score := 0
	for sq := A8; sq <= H1; sq++ {
		if KnightAttack(b, sq, -1) > 0 &&
			KingAttack(b, sq) > 0 {
			score++
		}
	}
	return score
}

// shelterStrength it's the shelter bonus for king position. If we can castle use
// the penalty after castling if (ShelterStrength + SheleterStorm) is smaller
func shelterStrength(b *board.Board, square int) int {
	w, s, tx := 0, 1024, -1

	rank, file := 0, 0

	for sq := A8; sq <= H1; sq++ {
		rank = sq / 8
		file = sq % 8
		if b.Bitboards[BK].Test(sq) ||
			uint(b.Castlings)&ShortB != 0 && file == 6 && rank == 0 ||
			uint(b.Castlings)&LongB != 0 && file == 2 && rank == 0 {

			w1 := strengthSquare(b, sq)
			s1 := stormSquare(b, sq, false)

			if s1-w1 < s-w {
				w = w1
				s = s1
				tx = max(1, min(6, file))
			}
		}
	}

	if square == -1 {
		return w
	}

	rank = square / 8
	file = square % 8

	if tx != -1 && b.Bitboards[BP].Test(square) &&
		file >= tx-1 && file <= tx+1 {

		for y := rank - 1; y >= 0; y-- {
			if b.Bitboards[BP].Test(y*8 + file) {
				return 0
			}
		}

		return 1
	}

	return 0
}

// strengthSquare returns king shelter square for each square on the board
func strengthSquare(b *board.Board, sq int) int {
	score := 5

	rank := sq / 8
	file := sq % 8

	kx := min(6, max(1, file))
	weakness := [][]int{
		{-6, 81, 93, 58, 39, 18, 25},
		{-43, 61, 35, -49, -29, -11, -63},
		{-10, 75, 23, -2, 32, 3, -45},
		{-39, -13, -29, -52, -48, -67, -166},
	}

	for x := kx - 1; x <= kx+1; x++ {
		us := 0
		for y := 7; y >= rank; y-- {
			if b.Bitboards[BP].Test(y*8+x) &&
				(!b.Bitboards[WP].Test((y+1)*8+x-1) || x <= 0 || y >= 7) &&
				(!b.Bitboards[WP].Test((y+1)*8+x+1) || x >= 7 || y >= 7) {
				us = y
			}
		}

		f := min(x, 7-x)
		if weakness[f][us] != 0 && f >= 0 && f <= 3 && us >= 0 && us <= 6 {
			score += weakness[f][us]
		}
	}

	return score
}

// shelterStorm is a penalty for king position. If we can castle use the
// penalty after the castling if (ShelterWeakness + ShelterStorm) is smaller
func shelterStorm(b *board.Board, square int) int {
	w, s, tx := 0, 1024, -1

	rank, file := 0, 0

	for sq := A8; sq <= H1; sq++ {
		rank = sq / 8
		file = sq % 8
		if b.Bitboards[BK].Test(sq) ||
			uint(b.Castlings)&ShortB != 0 && file == 6 && rank == 0 ||
			uint(b.Castlings)&LongB != 0 && file == 2 && rank == 0 {
			w1 := strengthSquare(b, rank*8+file)
			s1 := stormSquare(b, rank*8+file, false)

			if s1-w1 < s-w {
				w = w1
				s = s1
				tx = max(1, min(6, file))
			}
		}
	}

	if square == -1 {
		return s
	}

	rank = square / 8
	file = square % 8

	if tx != -1 && (b.Bitboards[BP].Test(square) || b.Bitboards[WP].Test(square)) &&
		file >= tx-1 && file <= tx+1 {

		for y := rank - 1; y >= 0; y-- {
			if b.Occupancies[color.BOTH].Test(
				y*8+file,
			) == b.Occupancies[color.BOTH].Test(
				rank*8+file,
			) {
				return 0
			}
		}

		return 1
	}

	return 0
}

// stormSquare returns the enemy pawns for each square on board
func stormSquare(b *board.Board, sq int, isEndgame bool) int {
	score, eval := 0, 5

	rank := sq / 8
	file := sq % 8

	kx := min(6, max(1, file))
	unblockedstorm := [][]int{
		{85, -289, -166, 97, 50, 45, 50},
		{46, -25, 122, 45, 37, -10, 20},
		{-6, 51, 168, 34, -2, -22, -14},
		{-15, -11, 101, 4, 11, -15, -29},
	}
	blockedstorm := [][]int{
		{0, 0, 76, -10, -7, -4, -1},
		{0, 0, 78, 15, 10, 6, 2},
	}

	for x := kx - 1; x <= kx+1; x++ {
		us, them := 0, 0

		for y := 7; y >= rank; y-- {
			if b.Bitboards[BP].Test(y*8+x) &&
				(!b.Bitboards[WP].Test((y+1)*8+x-1) && x > 0 && y < 7) &&
				(!b.Bitboards[WP].Test((y+1)*8+x+1) && x < 7 && y < 7) {
				us = y
			}

			if b.Bitboards[WP].Test(y*8 + x) {
				them = y
			}
		}

		f := min(x, 7-x)
		if us > 0 && them == us+1 {
			score += blockedstorm[0][them]
			eval += blockedstorm[1][them]
		} else {
			score += unblockedstorm[f][them]
		}
	}

	if isEndgame {
		return eval
	}

	return score
}

// pawnlessFlank is a penalty when our king is on a pawnless flank
func pawnlessFlank(b *board.Board) bool {
	pawns := []int{0, 0, 0, 0, 0, 0, 0, 0}
	kx := 0

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if b.Bitboards[WP].Test(y*8+x) || b.Bitboards[BP].Test(y*8+x) {
				pawns[x]++
			}

			if b.Bitboards[BK].Test(y*8 + x) {
				kx = x
			}
		}
	}

	sum := 0
	if kx == 0 {
		sum = pawns[0] + pawns[1] + pawns[2]
	} else if kx < 3 {
		sum = pawns[0] + pawns[1] + pawns[2] + pawns[3]
	} else if kx < 5 {
		sum = pawns[2] + pawns[3] + pawns[4] + pawns[5]
	} else if kx < 7 {
		sum = pawns[4] + pawns[5] + pawns[6] + pawns[7]
	} else {
		sum = pawns[5] + pawns[6] + pawns[7]
	}

	if sum == 0 {
		return true
	}

	return false
}

// kingPawnDistance is the minimal distance of our king to our pawns
func kingPawnDistance(b *board.Board) int {
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

// endgameShelter adds an endgame component to the blockedstorm penalty
// so that the penalty applies more uniformly through the game.
func endgameShelter(b *board.Board, square int) int {
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

			w1 := strengthSquare(b, sq)
			s1 := stormSquare(b, sq, false)
			e1 := stormSquare(b, sq, true)

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
