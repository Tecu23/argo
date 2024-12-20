package engine

import (
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/constants"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

var engine *Engine

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq()
	hash.Init()

	opts := NewOptions()
	engine = NewEngine(opts)
}

// func TestTranspositionTable(t *testing.T) {
// 	options := NewOptions()
// 	e := NewEngine(options)
//
// 	positions := []struct {
// 		fen   string
// 		depth int
// 	}{
// 		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5},
// 		{"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 1", 4},
// 	}
//
// 	for _, pos := range positions {
// 		b, _ := board.ParseFEN(pos.fen)
//
// 		// First search - should fill TT
// 		firstScore := e.Search(context.Background(), SearchParams{
// 			Limits: LimitsType{
// 				Depth: pos.depth,
// 			},
// 			Boards: []board.Board{b},
// 		})
//
// 		// Second search - should be much faster due to TT hits
// 		start := time.Now()
// 		secondScore := e.Search(context.Background(), SearchParams{
// 			Limits: LimitsType{
// 				Depth: pos.depth,
// 			},
// 			Boards: []board.Board{b},
// 		})
// 		secondTime := time.Since(start)
//
// 		if firstScore.Score != secondScore.Score {
// 			t.Errorf("Scores differ: first=%v, second=%v", firstScore, secondScore)
// 		}
//
// 		// Add logging to see improvement
// 		t.Logf("Position: %s, Depth: %d, Time: %v", pos.fen, pos.depth, secondTime)
// 	}
// }

func TestMoveGenConsistency(t *testing.T) {
	_ = SearchInfo{}
	b, _ := board.ParseFEN(StartPosition)
	b.PrintBoard()
	for depth := 2; depth <= 5; depth++ {
		nodes1 := board.PerftTest(&b, depth)
		nodes2 := board.PerftTest(&b, depth)
		if nodes1 != nodes2 {
			t.Errorf("Inconsistent node counts at depth %d: %d vs %d",
				depth, nodes1, nodes2)
		}
	}
}

func TestEvaluationSymmetry(t *testing.T) {
	positions := []struct {
		name string
		fen  string
		desc string
	}{
		{
			name: "Starting Position",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			desc: "Initial board state",
		},
		{
			name: "Early Opening",
			fen:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			desc: "After 1.e4",
		},
		{
			name: "Complex Middle Game",
			fen:  "r1bq1rk1/ppp2ppp/2np1n2/2b1p3/2B1P3/2NP1N2/PPP2PPP/R1BQ1RK1 w - - 0 1",
			desc: "Symmetric pawn structure with developed pieces",
		},
		{
			name: "Isolated Pawns",
			fen:  "8/3p4/8/1P6/8/8/8/8 w - - 0 1",
			desc: "Testing pawn structure evaluation",
		},
		{
			name: "Doubled Pawns",
			fen:  "8/3p4/3p4/8/8/2P5/2P5/8 w - - 0 1",
			desc: "Testing doubled pawns evaluation",
		},
		{
			name: "Passed Pawns",
			fen:  "8/1p6/8/8/6P1/8/8/8 w - - 0 1",
			desc: "Testing passed pawn evaluation",
		},
		{
			name: "Connected Passed Pawns",
			fen:  "8/2pp4/8/8/8/8/3PP3/8 w - - 0 1",
			desc: "Testing connected passed pawns",
		},
		{
			name: "Rook on Open File",
			fen:  "8/8/8/3p4/8/8/3P4/R7 w - - 0 1",
			desc: "Testing rook positioning evaluation",
		},
		{
			name: "Bishop Pair",
			fen:  "8/8/8/8/8/8/8/BB6 w - - 0 1",
			desc: "Testing bishop pair bonus",
		},
		{
			name: "Knight Outpost",
			fen:  "8/2p5/8/3N4/8/8/2P5/8 w - - 0 1",
			desc: "Testing knight positioning",
		},
		{
			name: "Queen Mobility",
			fen:  "8/8/8/3p4/8/4Q3/8/8 w - - 0 1",
			desc: "Testing piece mobility evaluation",
		},
		{
			name: "King Safety",
			fen:  "8/pppp4/8/8/8/8/PPPP4/K7 w - - 0 1",
			desc: "Testing king safety evaluation",
		},
		{
			name: "Rook Behind Passed Pawn",
			fen:  "8/4p3/8/8/8/4R3/8/8 w - - 0 1",
			desc: "Testing rook behind passed pawn",
		},
		{
			name: "Blocked Center",
			fen:  "8/3pp3/4p3/3P4/4P3/8/8/8 w - - 0 1",
			desc: "Testing blocked center evaluation",
		},
		{
			name: "Complex Pawn Chain",
			fen:  "8/pp6/2p5/3p4/4P3/5P2/6PP/8 w - - 0 1",
			desc: "Testing pawn chain evaluation",
		},
	}

	for _, test := range positions {
		t.Run(test.name, func(t *testing.T) {
			b, _ := board.ParseFEN(test.fen)
			mirrored := b.Mirror()

			// Get evaluations
			eval1 := engine.evaluator.Evaluate(&b)
			eval2 := -engine.evaluator.Evaluate(mirrored)

			if test.desc == "Testing rook positioning evaluation" {
				b.PrintBoard()
				mirrored.PrintBoard()
			}

			// Test exact symmetry
			if eval1 != eval2 {
				t.Errorf(
					"%s: Asymmetric evaluation\nOriginal: %d\nMirrored: %d\nDiff: %d\nFEN: %s\nDesc: %s",
					test.name,
					eval1,
					eval2,
					eval1-eval2,
					test.fen,
					test.desc,
				)
			}
		})
	}
}

// Helper to test evaluation stability through moves
func TestEvaluationSymmetryThroughMoves(t *testing.T) {
	// Test position with a few moves
	type testMove struct {
		from, to string
		evalDiff int // expected evaluation difference
	}

	tests := []struct {
		startFen string
		moves    []testMove
	}{
		{
			startFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []testMove{
				{"e2", "e4", 10}, // center pawn push
				{"d2", "d4", 10}, // second center pawn
				{"g1", "f3", 20}, // develop knight
			},
		},
	}

	for _, test := range tests {
		b, _ := board.ParseFEN(test.startFen)

		for _, move := range test.moves {
			// Make move on both boards
			newB, _ := b.ParseMove(move.from + move.to)
			newMirr := newB.Mirror()

			// Check evaluation symmetry
			eval1 := engine.evaluator.Evaluate(&newB)
			eval2 := -engine.evaluator.Evaluate(newMirr)

			if eval1 != eval2 {
				t.Errorf("Asymmetric evaluation after move %s-%s:\nOriginal: %d\nMirrored: %d",
					move.from, move.to, eval1, eval2)
			}
		}
	}
}
