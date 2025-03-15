package evalhelpers

import (
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	attacks.InitPawnAttacks()
	util.InitFen2Sq() // Make sure square mappings are initialized
	hash.Init()
}

func TestPawnAttackFunctions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test pawn attacks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := board.NewBoard()

			for sq := A8; sq <= H1; sq++ {
				method1 := PawnAttack(b, sq)
				method2 := PawnAttackTable(b, sq)

				if method1 != method2 {
					t.Errorf(
						"Pawn Attack is not the same for square %s, method: %d, table: %d",
						util.Sq2Fen[sq],
						method1,
						method2,
					)
				}
			}
		})
	}
}
