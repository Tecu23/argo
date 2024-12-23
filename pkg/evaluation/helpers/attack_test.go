package evalhelpers

import (
	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(Bishop)
	attacks.InitSliderPiecesAttacks(Rook)

	util.InitFen2Sq()

	hash.Init()
}

// // Additional test specifically for PawnAttack function
// func TestPawnAttack(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		fen      string
// 		square   int
// 		expected int
// 	}{
// 		// {
// 		// 	name:     "No Pawn Attacks",
// 		// 	fen:      "8/8/8/8/4B3/8/8/8 w - - 0 1",
// 		// 	square:   E4,
// 		// 	expected: 0,
// 		// },
// 		{
// 			name:     "Single Pawn Attack",
// 			fen:      "8/8/3P4/4B3/8/8/8/8 w - - 0 1",
// 			square:   D6,
// 			expected: 0,
// 		},
// 		// {
// 		// 	name:     "Double Pawn Attack",
// 		// 	fen:      "8/8/3P1P2/4B3/8/8/8/8 w - - 0 1",
// 		// 	square:   F6,
// 		// 	expected: 2,
// 		// },
// 		// {
// 		// 	name:     "Edge Square Test Left",
// 		// 	fen:      "8/8/1P6/B7/8/8/8/8 w - - 0 1",
// 		// 	square:   B6,
// 		// 	expected: 1,
// 		// },
// 		// {
// 		// 	name:     "Edge Square Test Right",
// 		// 	fen:      "8/8/6P1/7B/8/8/8/8 w - - 0 1",
// 		// 	square:   G6,
// 		// 	expected: 1,
// 		// },
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			b.PrintBoard()
// 			result := PawnAttack(&b, tt.square)
// 			if result != tt.expected {
// 				t.Errorf("%s: got %v, want %v", tt.name, result, tt.expected)
// 			}
// 		})
// 	}
// }
