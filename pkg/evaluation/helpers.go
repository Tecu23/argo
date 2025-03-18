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

// RookXrayAttack counts the number of rooks that attack sq
// If a sq2 is sent, the function evaluates if the rook at sq2 attacks sq
func RookXrayAttack(b *board.Board, sq, sq2 int) int {
	// Get occupied squares excluding queens
	occupancy := b.Occupancies[color.BOTH] & ^(b.Bitboards[WQ] | b.Bitboards[BQ] | b.Bitboards[WR])

	if sq2 != -1 {
		if !b.Bitboards[WR].Test(sq2) {
			return 0
		}

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

		// Check if sq2 rook attacks sq
		rookAttacks := attacks.GetRookAttacks(sq2, occupancy)
		if rookAttacks.Test(sq) {
			return 1
		}
		return 0
	}

	// Count white bishops attacking sq
	score := 0
	whiteRooks := b.Bitboards[WR]

	for whiteRooks != 0 {
		rookSq := whiteRooks.FirstOne()
		rookAttacks := attacks.GetRookAttacks(rookSq, occupancy)

		if rookAttacks.Test(sq) {
			// Check pin direction
			dir := PinnedDirection(b, rookSq)
			if dir == 0 {
				score++
			} else {
				// Get diagonal direction from bishop to sq
				rank := sq / 8
				file := sq % 8
				rank2 := rookSq / 8
				file2 := rookSq % 8

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

// QueenAttack counts the number of queens that attack sq
// If a sq2 is sent, the function evaluates if the queen at sq2 attacks sq
func QueenAttack(b *board.Board, sq, sq2 int) int {
	// Get occupied squares excluding queens
	occupancy := b.Occupancies[color.BOTH]

	if sq2 != -1 {
		if !b.Bitboards[WQ].Test(sq2) {
			return 0
		}

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

		// Check if sq2 rook attacks sq
		queenAttacks := attacks.GetQueenAttacks(sq2, occupancy)
		if queenAttacks.Test(sq) {
			return 1
		}
		return 0
	}

	// Count white bishops attacking sq
	score := 0
	whiteQueens := b.Bitboards[WQ]

	for whiteQueens != 0 {
		queenSq := whiteQueens.FirstOne()
		queenAttacks := attacks.GetQueenAttacks(queenSq, occupancy)

		if queenAttacks.Test(sq) {
			// Check pin direction
			dir := PinnedDirection(b, queenSq)
			if dir == 0 {
				score++
			} else {
				// Get diagonal direction from bishop to sq
				rank := sq / 8
				file := sq % 8
				rank2 := queenSq / 8
				file2 := queenSq % 8

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

func oldQueenAttack(b *board.Board, sq int, sq2 int) int {
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

// PinnedDirection returns the direction that the piece on sq is pinned to the king.
// If the piece is not pinned returns 0
func PinnedDirection(b *board.Board, sq int) int {
	if !b.Occupancies[color.BOTH].Test(sq) {
		return 0
	}

	c := 1
	if b.Occupancies[color.BLACK].Test(sq) {
		c = -1
	}

	kingBB := b.Bitboards[WK]
	kingSq := kingBB.FirstOne()

	dir := getDirection(sq, kingSq)
	if dir == DirNone {
		return 0
	}

	rank1, file1 := sq/8, sq%8

	// Calculate direction vectors based on direction
	var dx, dy int
	switch dir {
	case DirWest:
		dx, dy = -1, 0
	case DirNorth:
		dx, dy = 0, -1
	case DirEast:
		dx, dy = 1, 0
	case DirSouth:
		dx, dy = 0, 1
	case DirNorthWest:
		dx, dy = -1, -1
	case DirNorthEast:
		dx, dy = 1, -1
	case DirSouthWest:
		dx, dy = -1, 1
	case DirSouthEast:
		dx, dy = 1, 1
	}

	// Check for pieces between our piece and the king
	// Start one square away from our piece in the direction of the king
	for d := 1; ; d++ {
		checkRank := rank1 + d*dy
		checkFile := file1 + d*dx
		checkSq := checkRank*8 + checkFile

		// If we've reached the king, stop checking
		if checkSq == kingSq {
			break
		}

		// If there's a piece between us and the king, no pin possible
		if b.Occupancies[color.BOTH].Test(checkSq) {
			return 0
		}
	}

	// Check in the opposite direction for potential pinning pieces
	dx = -dx
	dy = -dy

	// Check if there are attacking pieces that could create a pin
	isDiagonal := (dir == DirNorthWest || dir == DirNorthEast ||
		dir == DirSouthWest || dir == DirSouthEast)

	for d := 1; d < 8; d++ {
		checkRank := rank1 + d*dy
		checkFile := file1 + d*dx

		// Check if we're still on the board
		if checkRank < 0 || checkRank > 7 || checkFile < 0 || checkFile > 7 {
			break
		}

		checkSq := checkRank*8 + checkFile

		// If we hit a piece, check if it's a potential pinner
		if b.Occupancies[color.BOTH].Test(checkSq) {
			if isDiagonal {
				// On diagonal, check for bishop or queen
				if b.Bitboards[BQ].Test(checkSq) || b.Bitboards[BB].Test(checkSq) {
					return abs(dx+dy*3) * c
				}
			} else {
				// On rank/file, check for rook or queen
				if b.Bitboards[BQ].Test(checkSq) || b.Bitboards[BR].Test(checkSq) {
					return abs(dx+dy*3) * c
				}
			}
			// If we hit a non-pinning piece, no pin possible
			break
		}
	}

	return 0
}

// getDirection is a helper function that returns the direction between 2 pieces
func getDirection(sq1, sq2 int) int {
	r1, f1 := sq1/8, sq1%8
	r2, f2 := sq2/8, sq2%8

	rDiff := r2 - r1
	fDiff := f2 - f1

	// If squares are the same, no direction
	if rDiff == 0 && fDiff == 0 {
		return DirNone
	}

	// Check horizontal alignment (same rank)
	if rDiff == 0 {
		if fDiff > 0 {
			return DirEast
		}
		return DirWest
	}

	// Check vertical alignment (same file)
	if fDiff == 0 {
		if rDiff > 0 {
			return DirSouth
		}
		return DirNorth
	}

	// Check diagonal alignment (abs of rank diff equals abs of file diff)
	if abs(rDiff) == abs(fDiff) {
		if rDiff > 0 && fDiff > 0 {
			return DirSouthEast
		} else if rDiff > 0 && fDiff < 0 {
			return DirSouthWest
		} else if rDiff < 0 && fDiff > 0 {
			return DirNorthEast
		}
		return DirNorthWest
	}

	// Not aligned
	return DirNone
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
