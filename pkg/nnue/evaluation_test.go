// Copyright (C) 2025 Tecu23
// Port of Koivisto evaluation, licensed under GNU GPL v3

package nnue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
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

	err := InitializeNNUE()
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

func TestProcessMoveAndEvaluate(t *testing.T) {
	testCases := []struct {
		Name          string
		StartFEN      string
		Moves         []move.Move
		ExpectedEvals []int
	}{
		{
			Name:     "Opening Sequence",
			StartFEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			Moves: []move.Move{
				move.EncodeMove(E2, E4, WP, move.DoublePawnPush, 0),
				move.EncodeMove(E7, E5, BP, move.DoublePawnPush, 0),
				move.EncodeMove(G1, F3, WN, move.Quiet, 0),
				move.EncodeMove(B8, C6, BN, move.Quiet, 0),
			},
			ExpectedEvals: []int{
				105,
				-13,
				97,
				-54,
				100,
			},
		},
		{
			Name:     "Captures and Special Moves",
			StartFEN: "r3k2r/ppp2ppp/2n1bn2/3pP3/3P4/2N2N2/PPP2PPP/R3K2R w KQkq d6 0 1",
			Moves: []move.Move{
				move.EncodeMove(E5, D6, WP, move.EnPassant, BP),
				move.EncodeMove(E8, G8, BK, move.KingCastle, 0),
				move.EncodeMove(E1, G1, WK, move.KingCastle, 0),
				move.EncodeMove(D6, C7, WP, move.Capture, BP),
			},
			ExpectedEvals: []int{
				-605,
				607,
				-630,
				618,
				-306,
			},
		},
		{
			Name:     "King Movement and En Passant",
			StartFEN: "rnbqk2r/ppp1bppp/3p1n2/4p3/4P3/3B1N2/PPPP1PPP/RNBQK2R w KQkq - 0 1",
			Moves: []move.Move{
				move.EncodeMove(E1, G1, WK, move.KingCastle, 0),
				move.EncodeMove(C8, G3, BB, move.Quiet, 0),
				move.EncodeMove(F3, E5, WN, move.Capture, BP),
				move.EncodeMove(F6, E4, BN, move.Capture, WP),
				move.EncodeMove(D3, E4, WB, move.Capture, BN),
			},
			ExpectedEvals: []int{
				-106,
				67,
				288,
				-354,
				375,
				-804,
			},
		},
	}

	e := NewEvaluator()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b, _ := board.ParseFEN(tc.StartFEN)
			e.Reset(&b)

			// Test initial position evaluation
			initialEval := e.Evaluate(&b)
			t.Logf(
				"Initial position evaluation: %d (expected: %d)",
				initialEval,
				tc.ExpectedEvals[0],
			)
			assert.InDelta(
				t,
				tc.ExpectedEvals[0],
				initialEval,
				5,
				"Initial evaluation should match expected value (±5)",
			)

			// Apply each move and check evaluation
			for i, mv := range tc.Moves {
				// Process the move
				if !b.MakeMove(mv, board.AllMoves) {
					fmt.Println("SKIPPED MOVE")
				}
				e.ProcessMove(&b, mv)

				// Evaluate the position
				eval := e.Evaluate(&b)
				t.Logf(
					"After move %s: evaluation: %d (expected: %d)",
					mv,
					eval,
					tc.ExpectedEvals[i+1],
				)

				// Check if evaluation is close to expected
				assert.InDelta(t, tc.ExpectedEvals[i+1], eval, 5,
					"Evaluation after %s should match expected value (±5)", mv)

			}
			e.PopAccumulation()
		})
	}
}

func BenchmarkEval(b *testing.B) {
	board, _ := board.ParseFEN("r1bqkbnr/ppp2ppp/2np4/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4")
	evaluator := NewEvaluator()

	evaluator.Reset(&board)
	for i := 0; i < b.N; i++ {
		evaluator.Evaluate(&board)
	}
}
