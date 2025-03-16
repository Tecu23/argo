// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"math"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// Space computes the space evaluation for a given side. The Space are
// bonus is multiplied by a weight: number of our pieces minus two times
// number of open files. The aim is to improve play on game opening
func Space(b *board.Board) int {
	if nonPawnMaterial(b, color.WHITE)+nonPawnMaterial(b.Mirror(), color.WHITE) < 12222 {
		return 0
	}
	score := 0.0
	pieceCount := b.Occupancies[color.WHITE].Count()

	for sq := A8; sq <= H1; sq++ {
		blockedCount := 0
		spacearea := 0

		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {

				if b.Bitboards[WP].Test(y*8+x) &&
					(b.Bitboards[BP].Test((y-1)*8+x) ||
						(b.Bitboards[BP].Test((y-2)*8+x-1) &&
							b.Bitboards[BP].Test((y-2)*8+x+1))) {
					blockedCount++
				}

				if b.Bitboards[BP].Test(y*8+x) &&
					(b.Bitboards[WP].Test((y+1)*8+x) ||
						(b.Bitboards[WP].Test((y+2)*8+x-1) &&
							b.Bitboards[WP].Test((y+2)*8+x+1))) {
					blockedCount++
				}
			}
		}
		spacearea += SpaceArea(b, sq)
		weight := pieceCount - 3 + min(blockedCount, 9)
		score += float64(spacearea * weight * weight)

	}
	return int(math.Floor(score / 16))
}

// SpaceArea returns the number of safe squares available for minor pieces
// on the central four files on ranks 2 to 4. Safe squares one, two or three
// squares behind a friendly pawn are counted twice
func SpaceArea(b *board.Board, sq int) int {
	score := 0

	rank := sq / 8
	file := sq % 8

	if ((8-rank) >= 2 && (8-rank) <= 4 && file+1 >= 3 && file+1 <= 6) &&
		!b.Bitboards[WP].Test(sq) &&
		(!b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0) &&
		(!b.Bitboards[BP].Test((rank-1)*8+file+1) && rank > 0 && file < 7) {
		score++

		if ((b.Bitboards[WP].Test((rank-1)*8+file) && rank > 0) ||
			(b.Bitboards[WP].Test((rank-2)*8+file) && rank > 1) ||
			(b.Bitboards[WP].Test((rank-3)*8+file) && rank > 2)) &&
			attack(b.Mirror(), (7-rank)*8+file) == 0 {
			score++
		}
	}

	return score
}
