package nnue

import (
	"fmt"
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq()

	hash.Init()

	err := LoadWeights("../../default.net")
	if err != nil {
		// This will cause the tests to fail immediately if weights don't load
		panic(fmt.Sprintf("Failed to load NNUE weights: %v", err))
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		eval int
	}{
		{
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			name: "Starting Position",
			eval: 105,
		},
		{fen: "4k3/8/8/8/8/8/8/4K3 w - - 0 1", name: "Empty Board with Just Kings", eval: 8},
		{fen: "4k3/8/8/8/8/8/4P3/4K3 w - - 0 1", name: "White Pawn Advantage", eval: 99},
		{fen: "4k3/8/8/8/8/8/4N3/4K3 w - - 0 1", name: "White Knight Advantage", eval: 166},
		{fen: "4k3/8/8/8/8/8/4Q3/4K3 w - - 0 1", name: "White Queen Advantage", eval: 398},
		{fen: "8/8/8/8/8/8/8/4K2k w - - 0 1", name: "Kings in Proximity", eval: 10},
		{fen: "k7/8/8/8/8/8/8/K7 w - - 0 1", name: "Kings in Corners", eval: -1},
		{fen: "3k4/8/8/8/8/8/8/3K4 w - - 0 1", name: "Kings in Center Files", eval: 8},
		{fen: "4k3/8/8/8/8/3P1P2/4P3/4K3 w - - 0 1", name: "Pawn Chain Structure", eval: 373},
		{fen: "4k3/8/8/8/8/3P4/4P3/4K3 w - - 0 1", name: "Isolated Pawns", eval: 236},
		{
			fen:  "4k3/8/8/8/8/8/4N3/3QK3 w - - 0 1",
			name: "Piece Coordination and King Safety",
			eval: 1017,
		},
		{fen: "4k3/8/8/8/8/8/4R3/4K3 w - - 0 1", name: "Rook Placement", eval: 290},
		{fen: "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", name: "Advanced Pawn", eval: 86},
		{fen: "3k4/8/8/8/8/8/8/R3K2R w KQ - 0 1", name: "Castling Rights Consideration", eval: 606},
		{fen: "4k3/8/8/8/8/2N5/3B4/4K3 w - - 0 1", name: "Minor Piece Coordination", eval: 335},
		{fen: "4k3/8/8/8/8/2B5/3B4/4K3 w - - 0 1", name: "Bishop Pair", eval: 209},
		{
			fen:  "r1bqkbnr/ppp2ppp/2np4/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			name: "Early Opening Position",
			eval: 147,
		},
		{
			fen:  "r1bk3r/p2pBpNp/n4n2/1p1NP2P/6P1/3P4/P1P1K3/q5b1 w - - 0 1",
			name: "Complex Tactical Position",
			eval: -1325,
		},
	}

	var result int
	evaluator := NewEvaluator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)

			evaluator.Reset(&b)

			result = evaluator.Evaluate(&b)
			if result != tt.eval {
				t.Errorf("Middle Game Evaluation, %s: got %v, want %v", tt.name, result, tt.eval)
			}
		})
	}
}
