// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func TestDoubleIsolated(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		result bool
		sq     int
	}{
		// Complex Position
		{
			name:   "Single Double Isolated White Pawn on D4 #1",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: true,
			sq:     D4,
		},
		{
			name:   "Single Double Isolated White Pawn on D4 #2",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D3,
		},
		{
			name:   "Single Double Isolated White Pawn on D4 #3",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     A2,
		},
		{
			name:   "Double Double Isolated White Pawn on D4 #4",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP1P2PP/RNBQKBNR w KQkq - 0 2",
			result: true,
			sq:     D4,
		},
		{
			name:   "Double Double Isolated White Pawn on D4 #5",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP1P2PP/RNBQKBNR w KQkq - 0 2",
			result: true,
			sq:     D3,
		},
		{
			name:   "Double Double Isolated White Pawn on D4 #6",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP1P2PP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     G2,
		},
		// Additional test cases
		{
			name:   "No black pawns in front, not double isolated",
			fen:    "rnbqkbnr/pp3ppp/8/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D4,
		},
		{
			name:   "Black pawn on adjacent file (C6), not double isolated",
			fen:    "rnbqkbnr/pp3ppp/2pp4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D4,
		},
		{
			name:   "Black pawn on adjacent file (E6), not double isolated",
			fen:    "rnbqkbnr/pp3ppp/3pp3/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D4,
		},
		{
			name:   "Isolated but not doubled, not double isolated",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/8/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D4,
		},
		{
			name:   "Three pawns in a file, top pawn is double isolated",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP1P2PP/RNBQKBNR b KQkq - 0 2",
			result: true,
			sq:     D4,
		},
		{
			name:   "Three pawns in a file, middle pawn is not double isolated",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/8/3P4/PP1P2PP/RNBQKBNR b KQkq - 0 2",
			result: true,
			sq:     D3,
		},
		{
			name:   "Edge file pawn (A-file) is double isolated",
			fen:    "rnbqkbnr/5ppp/p1pp4/8/P7/P7/1P3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     A4,
		},
		{
			name:   "Edge file pawn (H-file) is double isolated",
			fen:    "rnbqkbnr/ppp3p1/3p3p/8/7P/7P/PPP3P1/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     H4,
		},
		{
			name:   "No other pawns on the board, not double isolated",
			fen:    "rnbqkbnr/8/8/8/3P4/8/8/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "No other white pawns, empty square is not double isolated",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/8/8/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			var res bool
			if !isolated(&b, tt.sq) {
				res = false
			} else {
				res = doubleIsolated(&b, tt.sq)
			}
			if res != tt.result {
				t.Errorf("Pawn Evaluation failed, %s: got %v, want %v", tt.name, res, tt.result)
			}
		})
	}
}

func TestIsolated(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		result bool
		sq     int
	}{
		// Basic test cases
		{
			name:   "Isolated white pawn on D4",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: true,
			sq:     D4,
		},
		{
			name:   "Isolated white pawn on D3",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: true,
			sq:     D3,
		},
		{
			name:   "Non-isolated white pawn on A2 (adjacent to B2)",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PP3PPP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     A2,
		},
		{
			name:   "Non-isolated white pawn on D4 with pawn on C2",
			fen:    "rnbqkbnr/pp3ppp/3p4/8/3P4/3P4/PPP3PP/RNBQKBNR w KQkq - 0 2",
			result: false,
			sq:     D4,
		},
		// Edge cases
		{
			name:   "Edge file pawn on A4 is isolated",
			fen:    "rnbqkbnr/1ppppppp/8/8/P7/8/2PPPPPP/RNBQKBNR w KQkq - 0 1",
			result: true,
			sq:     A4,
		},
		{
			name:   "Edge file pawn on A4 is not isolated with pawn on B-file",
			fen:    "rnbqkbnr/pppppppp/8/8/PB6/8/1PPPPPPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     A4,
		},
		{
			name:   "Edge file pawn on H4 is isolated",
			fen:    "rnbqkbnr/ppppppp1/8/8/7P/8/PPPPPP2/RNBQKBNR w KQkq - 0 1",
			result: true,
			sq:     H4,
		},
		{
			name:   "Edge file pawn on H4 is not isolated with pawn on G-file",
			fen:    "rnbqkbnr/pppppppp/8/8/6PP/8/PPPPPP11/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     H4,
		},
		// Special cases
		{
			name:   "Square with no pawn shouldn't be considered isolated",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "Isolated pawn with friendly pawn on same file",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		// Complex positions
		{
			name:   "Multiple pawns that seem isolated",
			fen:    "rnbqkbnr/pppppppp/8/8/1P1P1P2/8/P1P1P1PP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "Pawn chain (not isolated)",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/2P5/PP2PPPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "Pawn with diagonal neighbor only (non-isolated)",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/4P3/PPP2PPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "Empty square should not be isolated",
			fen:    "rnbqkbnr/pppppppp/8/8/4P3/8/PPP3PP/RNB1KBNR b KQkq e3 0 1",
			result: true,
			sq:     E4,
		},
		{
			name:   "Queen on square is not an isolated pawn",
			fen:    "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNB1KBNR w KQkq - 0 1",
			result: false,
			sq:     D4,
		},
		{
			name:   "Starting position for white pawns",
			fen:    "rnbqkbnr/ppp1pppp/8/3p4/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			res := isolated(&b, tt.sq)
			if res != tt.result {
				t.Errorf("Pawn Evaluation failed, %s: got %v, want %v", tt.name, res, tt.result)
			}
		})
	}
}

func TestBackward(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		result bool
		sq     int
	}{
		{
			name:   "Backward pawn behind friendly pawns on adjacent files",
			fen:    "rnbqkbnr/pppppppp/8/8/2P1P3/3P4/PP3PPP/RNBQKBNR w KQkq - 0 1",
			result: false, // D3 is not backward because there are friendly pawns ahead on C4 and E4
			sq:     D3,
		},
		{
			name:   "Backward pawn with no friendly pawns on adjacent files",
			fen:    "rnbqkbnr/pp1ppppp/8/1p1P4/1P6/2P5/P3PPPP/RNBQKBNR b KQkq - 0 1",
			result: true, // D3 is backward, no friendly pawns ahead on adjacent files and advance is unsafe
			sq:     C3,
		},
		{
			name:   "Non-backward pawn despite no friendly pawns, can safely advance",
			fen:    "rnbqkbnr/pp2pppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			result: false, // D3 can safely advance, so not backward
			sq:     D3,
		},
		{
			name:   "Backward pawn blocked directly in front",
			fen:    "rnbqkbnr/ppp2ppp/3p4/2P1p3/8/3P4/PP3PPP/RNBQKBNR w KQkq - 0 1",
			result: true, // D3 is blocked by black pawn directly in front
			sq:     D3,
		},
		// Edge file cases
		{
			name:   "Edge file (A-file) backward pawn",
			fen:    "rnbqkbnr/pppppppp/8/8/1P6/P7/1P1PPPPP/RNBQKBNR w KQkq - 0 1",
			result: false, // A3 is not backward because there's a friendly pawn on B4
			sq:     A3,
		},
		{
			name:   "Edge file (A-file) backward pawn, threatened by diagonal attacker",
			fen:    "rnbqkbnr/2pppppp/p7/1p6/8/P7/2PPPPPP/RNBQKBNR w KQkq - 0 1",
			result: true, // A3 is backward, advance is unsafe due to diagonal attacker
			sq:     A3,
		},
		{
			name:   "Edge file (H-file) backward pawn",
			fen:    "rnbqkbnr/pppppppp/8/8/6P1/7P/PPPPPPP1/RNBQKBNR w KQkq - 0 1",
			result: false, // H3 is not backward because there's a friendly pawn on G4
			sq:     H3,
		},
		// Complex scenarios
		{
			name:   "Potential backward pawn with diagonal attackers",
			fen:    "rnbqkbnr/ppp1p1pp/5p2/3p4/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			result: false,
			sq:     D3,
		},
		{
			name:   "Pawn that appears backward but has support",
			fen:    "rnbqkbnr/pppppppp/8/8/2P1P3/3P4/PP3PPP/RNBQKBNR w KQkq - 0 1",
			result: false, // D3 has support from pawns on C4 and E4
			sq:     D3,
		},
		{
			name:   "Pawn on rank 2 can be backward",
			fen:    "rnbqkbnr/ppp1pppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			result: false, // D4 can advance safely
			sq:     D4,
		},
		{
			name:   "Pawn on rank 2 with direct blocker",
			fen:    "rnbqkbnr/ppp1pppp/3p4/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
			result: false, // D2 is backward, blocked directly
			sq:     D4,
		},
		// // Non-pawn cases
		// {
		// 	name:   "Empty square is not backward",
		// 	fen:    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// 	result: false,
		// 	sq:     D2,
		// },
		// {
		// 	name:   "Other piece (not pawn) is not backward",
		// 	fen:    "rnbqkbnr/ppp1pppp/8/3p4/5P2/4P3/PPP3PP/RNBQKB1R w KQkq d6 0 2",
		// 	result: true,
		// 	sq:     E3,
		// },
		// // Special cases
		// {
		// 	name:   "Pawn on last rank (theoretical)",
		// 	fen:    "rnbqkPnr/pppppppp/8/8/8/8/PPPPP1PP/RNBQKBNR w KQkq - 0 1",
		// 	result: false, // F8 pawn can't advance further
		// 	sq:     F8,
		// },
		// {
		// 	name:   "Advanced pawn with diagonal attackers",
		// 	fen:    "rnbqkbnr/ppp1pppp/8/2P1p3/4P3/3P4/PP3PPP/RNBQKBNR b KQkq e3 0 1",
		// 	result: true, // D3 is backward due to diagonal attacker on E5
		// 	sq:     D3,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			res := optimizedBackward(&b, tt.sq)
			if res != tt.result {
				t.Errorf("Pawn Evaluation failed, %s: got %v, want %v", tt.name, res, tt.result)
			}
		})
	}
}
