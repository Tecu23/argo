package evalhelpers

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

func PawnAttackTable(b *board.Board, sq int) int {
	return (attacks.PawnAttacks[color.WHITE][sq] & b.Occupancies[color.BLACK]).Count()
}
