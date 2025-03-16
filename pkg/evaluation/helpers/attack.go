package evaluationhelpers

import (
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

		if b.Bitboards[WN].Test((rank+iy)*8+file+ix) &&
			(rank+iy >= 0 && rank+iy <= 7 && file+ix >= 0 && file+ix <= 7) &&
			(sq2 == -1 || file2 == file+ix && rank2 == rank+iy) &&
			Pinned(b, (rank+iy)*8+file+ix) == 0 {
			score++
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

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
