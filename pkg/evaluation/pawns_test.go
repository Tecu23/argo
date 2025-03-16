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
			res := optimizedDoubleIsolated(&b, tt.sq)
			if res != tt.result {
				t.Errorf("Pawn Evaluation failed, %s: got %v, want %v", tt.name, res, tt.result)
			}
		})
	}
}
