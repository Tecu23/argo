package evaluation

import (
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PawnAttack counts the number of attacks on square by pawn
// Pins or en-passant attacks are not considered here
func PawnAttack(b *board.Board, sq int) int {
	v := 0
	file := sq % 8
	rank := sq / 8
	if b.Bitboards[WP].Test((rank+1)*8+file-1) && file > 0 && rank < 7 {
		v++
	}

	if b.Bitboards[WP].Test((rank+1)*8+file+1) && file < 7 && rank < 7 {
		v++
	}
	return v
}

// KnightAttack counts the number of knights that attack sq
// If a sq2 is sent, the function evaluates if the knight at sq2 attacks sq
func KnightAttack(b *board.Board, sq int, sq2 int) int {
	if sq2 != -1 {
		// Check if sq2 has a white knight
		if !b.Bitboards[WN].Test(sq2) {
			return 0
		}

		// Check if knight is pinned
		if Pinned(b, sq2) != 0 {
			return 0 // Pinned knights can't attack
		}

		// Check if knight at sq2 attacks sq
		if attacks.KnightAttacks[sq2].Test(sq) {
			return 1
		}
		return 0
	}

	score := 0

	knightsBB := (b.Bitboards[WN] & attacks.KnightAttacks[sq])
	for knightsBB != 0 {
		knightsSq := knightsBB.FirstOne()
		if Pinned(b, knightsSq) == 0 {
			score++
		}
	}

	return score
}

// BishopXrayAttack counts the number of bishops that attack sq
// If a sq2 is sent, the function evaluates if the bishop at sq2 attacks sq
func BishopXrayAttack(b *board.Board, sq, sq2 int) int {
	// Get occupied squares excluding queens (which we can x-ray through)
	occupancy := b.Occupancies[color.BOTH] & ^(b.Bitboards[WQ] | b.Bitboards[BQ])

	// If sq2 is provided, check if a bishop on sq2 attacks sq
	if sq2 != -1 {
		// Check if sq2 has a white bishop
		if !b.Bitboards[WB].Test(sq2) {
			return 0
		}

		// Check pin direction
		dir := PinnedDirection(b, sq2)
		if dir != 0 {
			// Get diagonal direction from sq2 to sq
			rank := sq / 8
			file := sq % 8
			rank2 := sq2 / 8
			file2 := sq2 % 8

			// Calculate direction factors
			ix := sign(file - file2)
			iy := sign(rank - rank2)

			// Check if pin allows movement in this direction
			if abs(ix+iy*3) != dir {
				return 0
			}
		}

		// Check if sq2 bishop attacks sq
		bishopAttacks := attacks.GetBishopAttacks(sq2, occupancy)
		if bishopAttacks.Test(sq) {
			return 1
		}
		return 0
	}

	// Count white bishops attacking sq
	score := 0
	whiteBishops := b.Bitboards[WB]

	for whiteBishops != 0 {
		bishopSq := whiteBishops.FirstOne()
		bishopAttacks := attacks.GetBishopAttacks(bishopSq, occupancy)

		if bishopAttacks.Test(sq) {
			// Check pin direction
			dir := PinnedDirection(b, bishopSq)
			if dir == 0 {
				score++
			} else {
				// Get diagonal direction from bishop to sq
				rank := sq / 8
				file := sq % 8
				rank2 := bishopSq / 8
				file2 := bishopSq % 8

				ix := sign(file - file2)
				iy := sign(rank - rank2)

				if abs(ix+iy*3) == dir {
					score++
				}
			}
		}
	}

	return score
}

func oldBishopXrayAttack(b *board.Board, sq int, sq2 int) int {
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
				(sq2 == -1 || file2 == file+d*ix && rank2 == rank+d*iy) {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) {
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
				(file+d*ix >= 0) && (file+d*ix <= 7) &&
				(sq2 == -1 || file2 == file+d*ix && rank2 == rank+d*iy) {

				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WR].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				!b.Bitboards[BQ].Test((rank+d*iy)*8+file+d*ix) {
				break
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
				(file+d*ix >= 0) && (file+d*ix <= 7) &&
				(sq2 == -1 || file2 == file+d*ix && rank2 == rank+d*iy) {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8 + file + d*ix) {
				break
			}
		}
	}

	return score
}

// QueenAttackDiagonal counts number of attacks on square by queen only with
// diagonal direction
func QueenAttackDiagonal(b *board.Board, sq int, sq2 int) int {
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

		if ix == 0 || iy == 0 {
			continue
		}

		for d := 1; d < 8; d++ {
			if b.Bitboards[WQ].Test((rank+d*iy)*8+file+d*ix) &&
				(file+d*ix >= 0 && file+d*ix <= 7 && rank+d*iy >= 0 && rank+d*iy <= 7) &&
				(sq2 == -1 || file2 == file+d*ix && rank2 == rank+d*iy) {
				dir := PinnedDirection(b, (rank+d*iy)*8+file+d*ix)

				if dir == 0 || abs(ix+iy*3) == dir {
					score++
				}
			}

			if b.Occupancies[color.BOTH].Test((rank+d*iy)*8 + file + d*ix) {
				break
			}
		}
	}

	return score
}

// KingAttack counts the number of attacks on a square by the king
func KingAttack(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	factor := 0
	for i := 0; i < 8; i++ {
		factor = 0
		if i > 3 {
			factor = 1
		}

		ix := (i+factor)%3 - 1
		iy := (((i + factor) / 3) << 0) - 1
		if b.Bitboards[WK].Test((rank+iy)*8+file+ix) &&
			file+ix >= 0 && file+ix <= 7 &&
			rank+iy >= 0 && rank+iy <= 7 {
			return 1
		}
	}
	return 0
}

func Pinned(b *board.Board, sq int) int {
	return PinnedDirection(b, sq)
}

func PinnedDirection(b *board.Board, sq int) int {
	if !b.Occupancies[color.BOTH].Test(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8
	c := 1

	if b.Occupancies[color.BLACK].Test(sq) {
		c = -1
	}

	for i := 0; i < 8; i++ {
		factor := 0
		if i > 3 {
			factor = 1
		}

		ix := (i+factor)%3 - 1
		iy := (((i + factor) / 3) << 0) - 1

		king := false

		for d := 1; d < 8; d++ {
			if b.Bitboards[WK].Test((rank+d*iy)*8 + file + d*ix) {
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

// Helper function to get sign of a number
func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}
