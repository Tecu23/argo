package move

import (
	"testing"

	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq()
}

// TestMoveEncoding tests the EncodeMove function and all getter methods
func TestMoveEncoding(t *testing.T) {
	testCases := []struct {
		name       string
		source     int
		target     int
		piece      int
		promoted   int
		capture    int
		doublePush int
		enpassant  int
		castling   int
	}{
		{
			name:       "Simple pawn move",
			source:     12, // e2
			target:     20, // e4
			piece:      1,  // pawn
			promoted:   0,  // no promotion
			capture:    0,  // no capture
			doublePush: 1,  // double push
			enpassant:  0,  // no en passant
			castling:   0,  // no castling
		},
		{
			name:       "Pawn capture with promotion",
			source:     51, // d7
			target:     58, // e8
			piece:      1,  // pawn
			promoted:   5,  // queen promotion
			capture:    1,  // capture
			doublePush: 0,  // no double push
			enpassant:  0,  // no en passant
			castling:   0,  // no castling
		},
		{
			name:       "En passant capture",
			source:     35, // e5
			target:     42, // f6
			piece:      1,  // pawn
			promoted:   0,  // no promotion
			capture:    1,  // capture
			doublePush: 0,  // no double push
			enpassant:  1,  // en passant
			castling:   0,  // no castling
		},
		{
			name:       "Kingside castling",
			source:     4, // e1
			target:     6, // g1
			piece:      6, // king
			promoted:   0, // no promotion
			capture:    0, // no capture
			doublePush: 0, // no double push
			enpassant:  0, // no en passant
			castling:   1, // castling
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode move
			move := EncodeMove(
				tc.source,
				tc.target,
				tc.piece,
				tc.promoted,
				tc.capture,
				tc.doublePush,
				tc.enpassant,
				tc.castling,
			)

			// Test all getters
			if got := move.GetSource(); got != tc.source {
				t.Errorf("GetSource() = %v, want %v", got, tc.source)
			}
			if got := move.GetTarget(); got != tc.target {
				t.Errorf("GetTarget() = %v, want %v", got, tc.target)
			}
			if got := move.GetPiece(); got != tc.piece {
				t.Errorf("GetPiece() = %v, want %v", got, tc.piece)
			}
			if got := move.GetPromoted(); got != tc.promoted {
				t.Errorf("GetPromoted() = %v, want %v", got, tc.promoted)
			}
			if got := move.GetCapture(); got != tc.capture {
				t.Errorf("GetCapture() = %v, want %v", got, tc.capture)
			}
			if got := move.GetDoublePush(); got != tc.doublePush {
				t.Errorf("GetDoublePush() = %v, want %v", got, tc.doublePush)
			}
			if got := move.GetEnpassant(); got != tc.enpassant {
				t.Errorf("GetEnpassant() = %v, want %v", got, tc.enpassant)
			}
			if got := move.GetCastling(); got != tc.castling {
				t.Errorf("GetCastling() = %v, want %v", got, tc.castling)
			}
		})
	}
}

// TestMoveBitMasks tests that the bit masks correctly isolate move components
func TestMoveBitMasks(t *testing.T) {
	testCases := []struct {
		name string
		move Move
		want struct {
			source     int
			target     int
			piece      int
			promoted   int
			capture    int
			doublePush int
			enpassant  int
			castling   int
		}
	}{
		{
			name: "Test all bits set",
			move: Move(0xFFFFFF), // All relevant bits set
			want: struct {
				source     int
				target     int
				piece      int
				promoted   int
				capture    int
				doublePush int
				enpassant  int
				castling   int
			}{
				source:     0x3f,
				target:     0x3f,
				piece:      0xf,
				promoted:   0xf,
				capture:    1,
				doublePush: 1,
				enpassant:  1,
				castling:   1,
			},
		},
		{
			name: "Test alternating bits",
			move: Move(0x555555), // Alternating bits
			want: struct {
				source     int
				target     int
				piece      int
				promoted   int
				capture    int
				doublePush int
				enpassant  int
				castling   int
			}{
				source:     0x15,
				target:     0x15,
				piece:      0x5,
				promoted:   0x5,
				capture:    0,
				doublePush: 1,
				enpassant:  0,
				castling:   1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			move := tc.move
			if got := move.GetSource(); got != tc.want.source {
				t.Errorf("GetSource() = %v, want %v", got, tc.want.source)
			}
			// ... similar checks for other components
		})
	}
}

// TestMoveString tests the String() method for move representation
func TestMoveString(t *testing.T) {
	testCases := []struct {
		name    string
		move    Move
		wantStr string
	}{
		{
			name:    "E2 to E4",
			move:    EncodeMove(E2, E4, 1, 0, 0, 1, 0, 0),
			wantStr: "e2e4",
		},
		{
			name:    "G1 to F3",
			move:    EncodeMove(G1, F3, 2, 0, 0, 0, 0, 0),
			wantStr: "g1f3",
		},
		{
			name:    "E7 to E8 with promotion",
			move:    EncodeMove(E7, E8, 1, 5, 0, 0, 0, 0),
			wantStr: "e7e8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.move.String(); got != tc.wantStr {
				t.Errorf("Move.String() = %v, want %v", got, tc.wantStr)
			}
		})
	}
}

// TestNoMove tests the NoMove constant
func TestNoMove(t *testing.T) {
	if NoMove != Move(0) {
		t.Errorf("NoMove should be 0, got %v", NoMove)
	}

	// Test all getters return 0 for NoMove
	if got := NoMove.GetSource(); got != 0 {
		t.Errorf("NoMove.GetSource() = %v, want 0", got)
	}
	if got := NoMove.GetTarget(); got != 0 {
		t.Errorf("NoMove.GetTarget() = %v, want 0", got)
	}
	if got := NoMove.GetPiece(); got != 0 {
		t.Errorf("NoMove.GetPiece() = %v, want 0", got)
	}
	if got := NoMove.GetPromoted(); got != 0 {
		t.Errorf("NoMove.GetPromoted() = %v, want 0", got)
	}
	if got := NoMove.GetCapture(); got != 0 {
		t.Errorf("NoMove.GetCapture() = %v, want 0", got)
	}
	if got := NoMove.GetDoublePush(); got != 0 {
		t.Errorf("NoMove.GetDoublePush() = %v, want 0", got)
	}
	if got := NoMove.GetEnpassant(); got != 0 {
		t.Errorf("NoMove.GetEnpassant() = %v, want 0", got)
	}
	if got := NoMove.GetCastling(); got != 0 {
		t.Errorf("NoMove.GetCastling() = %v, want 0", got)
	}
}

// BenchmarkMoveEncoding benchmarks the move encoding process
func BenchmarkMoveEncoding(b *testing.B) {
	source := 12    // e2
	target := 28    // e4
	piece := 1      // pawn
	promoted := 0   // no promotion
	capture := 0    // no capture
	doublePush := 1 // double push
	enpassant := 0  // no en passant
	castling := 0   // no castling

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeMove(source, target, piece, promoted, capture, doublePush, enpassant, castling)
	}
}

// BenchmarkMoveGetters benchmarks the move getter methods
func BenchmarkMoveGetters(b *testing.B) {
	move := EncodeMove(12, 28, 1, 0, 0, 1, 0, 0) // e2e4 double push

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		move.GetSource()
		move.GetTarget()
		move.GetPiece()
		move.GetPromoted()
		move.GetCapture()
		move.GetDoublePush()
		move.GetEnpassant()
		move.GetCastling()
	}
}
