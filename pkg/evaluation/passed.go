// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PassedPawnEvaluation returns bonuses for passed pawns. Scale down bonuses for candidate passers
// which need more than one pawn push to become passed, or have a pawn in front of them.
func (e *Evaluator) PassedPawnEvaluation(b *board.Board) (mg, eg int) {
	mg, eg = 0, 0

	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()

		if passedLeverable(b, sq) == 0 {
			continue
		}

		eg += kingProximity(b, sq)

		passRank := passedRank(b, sq)
		mg += []int{0, 10, 17, 15, 62, 168, 276}[passRank]
		eg += []int{0, 28, 33, 41, 72, 177, 260}[passRank]

		passedBlockBonus := passedBlock(b, sq)
		mg += passedBlockBonus
		eg += passedBlockBonus

		passedFileBonus := passedFile(b, sq)
		mg -= 11 * passedFileBonus
		eg -= 8 * passedFileBonus
	}

	return mg, eg
}

// passedLeverable returns candidate passers without candidate passers w/o
// feasible lever
func passedLeverable(b *board.Board, sq int) int {
	if !b.IsPassedPawn(sq) {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if !b.Bitboards[BP].Test((rank-1)*8+file) && rank > 0 {
		return 1
	}

	mirror := b.Mirror()
	for i := -1; i <= 1; i += 2 {
		sq1 := rank*8 + file + i
		sq2 := (7-rank)*8 + file + i

		if (b.Bitboards[WP].Test((rank+1)*8+file+i) && rank < 7 && file+i >= 0 && file+i <= 7) &&
			(!b.Occupancies[color.BLACK].Test(rank+8+file+i) && file+i >= 0 && file+i <= 7) &&
			(attack(b, sq1) > 0 || attack(mirror, sq2) <= 1) {
			return 1
		}
	}

	return 0
}

// passedRank is a bonus according to the rank of a passed pawn
func passedRank(b *board.Board, sq int) int {
	if passedLeverable(b, sq) == 0 {
		return 0
	}

	rank := sq / 8

	return 7 - rank
}

// passedFile is a bonus according to the file of a passed pawn
func passedFile(b *board.Board, sq int) int {
	if passedLeverable(b, sq) == 0 {
		return 0
	}

	file := sq % 8

	return min(file, 7-file)
}

// passedBlock adds bonus if passed pawn is free to advance. Bonus is
// adjusted based on attacked and defended status of block square and
// entire path in front of path
func passedBlock(b *board.Board, sq int) int {
	if passedLeverable(b, sq) == 0 {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if 8-rank < 4 {
		return 0
	}

	if b.Occupancies[color.BOTH].Test((rank-1)*8+file) && rank > 0 {
		return 0
	}

	r := 7 - rank
	w := 0
	if r > 2 {
		w = 5*r - 13
	}

	mirror := b.Mirror()

	defended, unsafe, wunsafe, defended1, unsafe1 := 0, 0, 0, 0, 0

	for y := rank - 1; y >= 0; y-- {
		if attack(b, y*8+file) > 0 {
			defended++
		}

		if attack(mirror, (7-y)*8+file) > 0 {
			unsafe++
		}

		if attack(mirror, (7-y)*8+file-1) > 0 && file > 0 {
			wunsafe++
		}

		if attack(mirror, (7-y)*8+file+1) > 0 && file < 7 {
			wunsafe++
		}

		if y == rank-1 {
			defended1 = defended
			unsafe1 = unsafe
		}
	}

	for y := rank + 1; y < 8; y++ {
		if b.Bitboards[WR].Test(y*8+file) || b.Bitboards[WQ].Test(y*8+file) {
			defended1 = rank
			defended = rank
		}

		if b.Bitboards[BR].Test(y*8+file) || b.Bitboards[BQ].Test(y*8+file) {
			unsafe1 = rank
			unsafe = rank
		}
	}

	k := 0

	if unsafe == 0 && wunsafe == 0 {
		k = 35
	} else if unsafe == 0 {
		k = 20
	} else if unsafe1 == 0 {
		k = 9
	}

	if defended1 != 0 {
		k += 5
	}

	return k * w
}

// kingProximity is an endgame bonus based on the king's proximity.
// If block square is not the queening square then consider also a second push
func kingProximity(b *board.Board, sq int) int {
	if passedLeverable(b, sq) == 0 {
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
