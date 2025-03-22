package move

import (
	"testing"

	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq()
}

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name         string
		source       int
		target       int
		piece        int
		promoted     int
		captured     int
		captureFlag  int
		doublePush   int
		enpassant    int
		castlingFlag int
		castlingType int
	}{
		{
			name:         "pawn push",
			source:       E2,
			target:       E3,
			piece:        WP,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "pawn double push",
			source:       E2,
			target:       E4,
			piece:        WP,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   1,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "pawn capture",
			source:       E2,
			target:       F3,
			piece:        WP,
			promoted:     0,
			captured:     BP,
			captureFlag:  1,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "pawn promotion",
			source:       E7,
			target:       E8,
			piece:        WP,
			promoted:     WQ,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "pawn promotion with capture",
			source:       E7,
			target:       F8,
			piece:        WP,
			promoted:     WQ,
			captured:     BR,
			captureFlag:  1,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "en passant capture",
			source:       E5,
			target:       F6,
			piece:        WP,
			promoted:     0,
			captured:     BP,
			captureFlag:  1,
			doublePush:   0,
			enpassant:    1,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "white kingside castling",
			source:       E1,
			target:       G1,
			piece:        WK,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 1,
			castlingType: int(WhiteKingCastle),
		},
		{
			name:         "white queenside castling",
			source:       E1,
			target:       C1,
			piece:        WK,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 1,
			castlingType: int(WhiteQueenCastle),
		},
		{
			name:         "black kingside castling",
			source:       E8,
			target:       G8,
			piece:        BK,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 1,
			castlingType: int(BlackKingCastle),
		},
		{
			name:         "black queenside castling",
			source:       E8,
			target:       C8,
			piece:        BK,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 1,
			castlingType: int(BlackQueenCastle),
		},
		{
			name:         "regular knight move",
			source:       B1,
			target:       C3,
			piece:        WN,
			promoted:     0,
			captured:     0,
			captureFlag:  0,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
		{
			name:         "knight capture",
			source:       B1,
			target:       C3,
			piece:        WN,
			promoted:     0,
			captured:     BP,
			captureFlag:  1,
			doublePush:   0,
			enpassant:    0,
			castlingFlag: 0,
			castlingType: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode the move
			move := EncodeMove(
				tt.source, tt.target, tt.piece, tt.promoted, tt.captured,
				tt.captureFlag, tt.doublePush, tt.enpassant, tt.castlingFlag, tt.castlingType,
			)

			// Test all the getters to make sure they extract the correct values
			if got := move.GetSourceSquare(); got != tt.source {
				t.Errorf("GetSourceSquare() = %v, want %v", got, tt.source)
			}
			if got := move.GetTargetSquare(); got != tt.target {
				t.Errorf("GetTargetSquare() = %v, want %v", got, tt.target)
			}
			if got := move.GetMovingPiece(); got != tt.piece {
				t.Errorf("GetMovingPiece() = %v, want %v", got, tt.piece)
			}
			if got := move.GetPromotedPiece(); got != tt.promoted {
				t.Errorf("GetPromotedPiece() = %v, want %v", got, tt.promoted)
			}
			if got := move.GetCapturedPiece(); got != tt.captured {
				t.Errorf("GetCapturedPiece() = %v, want %v", got, tt.captured)
			}
			if got := move.IsCapture(); got != (tt.captureFlag != 0) {
				t.Errorf("IsCapture() = %v, want %v", got, tt.captureFlag != 0)
			}
			if got := move.IsDoublePush(); got != (tt.doublePush != 0) {
				t.Errorf("IsDoublePush() = %v, want %v", got, tt.doublePush != 0)
			}
			if got := move.IsEnPassant(); got != (tt.enpassant != 0) {
				t.Errorf("IsEnPassant() = %v, want %v", got, tt.enpassant != 0)
			}
			if got := move.IsCastle(); got != (tt.castlingFlag != 0) {
				t.Errorf("IsCastle() = %v, want %v", got, tt.castlingFlag != 0)
			}
			if tt.castlingFlag != 0 {
				if got := move.GetCastleType(); got != CastleType(tt.castlingType) {
					t.Errorf("GetCastleType() = %v, want %v", got, CastleType(tt.castlingType))
				}
			}
		})
	}
}

func TestNoMove(t *testing.T) {
	if NoMove.GetSourceSquare() != 0 || NoMove.GetTargetSquare() != 0 {
		t.Errorf("NoMove should have source and target as 0, got %v, %v",
			NoMove.GetSourceSquare(), NoMove.GetTargetSquare())
	}

	if NoMove.String() != "0000" {
		t.Errorf("NoMove.String() = %v, want 0000", NoMove.String())
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		expected string
	}{
		{
			name:     "e2 to e4",
			move:     EncodeMove(E2, E4, WP, 0, 0, 0, 1, 0, 0, 0),
			expected: "e2e4",
		},
		{
			name:     "e7 to e8 queen promotion",
			move:     EncodeMove(E7, E8, WP, WQ, 0, 0, 0, 0, 0, 0),
			expected: "e7e8q",
		},
		{
			name:     "e7 to f8 queen promotion capture",
			move:     EncodeMove(E7, F8, WP, WQ, BR, 1, 0, 0, 0, 0),
			expected: "e7f8q",
		},
		{
			name:     "e7 to e8 rook promotion",
			move:     EncodeMove(E7, E8, WP, WR, 0, 0, 0, 0, 0, 0),
			expected: "e7e8r",
		},
		{
			name:     "e7 to e8 bishop promotion",
			move:     EncodeMove(E7, E8, WP, WB, 0, 0, 0, 0, 0, 0),
			expected: "e7e8b",
		},
		{
			name:     "e7 to e8 knight promotion",
			move:     EncodeMove(E7, E8, WP, WN, 0, 0, 0, 0, 0, 0),
			expected: "e7e8n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.move.String(); got != tt.expected {
				t.Errorf("Move.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsQueenCastle(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		expected bool
	}{
		{
			name:     "white kingside castling - not queenside",
			move:     EncodeMove(E1, G1, WK, 0, 0, 0, 0, 0, 1, int(WhiteKingCastle)),
			expected: false,
		},
		{
			name:     "white queenside castling",
			move:     EncodeMove(E1, C1, WK, 0, 0, 0, 0, 0, 1, int(WhiteQueenCastle)),
			expected: true,
		},
		{
			name:     "black kingside castling - not queenside",
			move:     EncodeMove(E8, G8, BK, 0, 0, 0, 0, 0, 1, int(BlackKingCastle)),
			expected: false,
		},
		{
			name:     "black queenside castling",
			move:     EncodeMove(E8, C8, BK, 0, 0, 0, 0, 0, 1, int(BlackQueenCastle)),
			expected: true,
		},
		{
			name:     "regular move - not castling",
			move:     EncodeMove(E1, E2, WK, 0, 0, 0, 0, 0, 0, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.move.IsQueenCastle(); got != tt.expected {
				t.Errorf("IsQueenCastle() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBitMasks(t *testing.T) {
	// Test all masks to ensure they correctly extract values
	source := E2      // e2
	target := E4      // e4
	piece := WP       // white pawn
	promoted := WQ    // white queen
	captured := BP    // black pawn
	captureFlag := 1  // is capture
	doublePush := 1   // is double push
	enpassant := 1    // is en passant
	castlingFlag := 1 // is castling
	castlingType := int(WhiteKingCastle)

	move := EncodeMove(
		source, target, piece, promoted, captured,
		captureFlag, doublePush, enpassant, castlingFlag, castlingType,
	)

	// Test bit manipulation directly
	if int(move&SourceMask) != source {
		t.Errorf("Source mask extraction failed: got %v, want %v", int(move&SourceMask), source)
	}
	if int(move&TargetMask)>>TargetShift != target {
		t.Errorf(
			"Target mask extraction failed: got %v, want %v",
			int(move&TargetMask)>>TargetShift,
			target,
		)
	}
	if int(move&PieceMask)>>PieceShift != piece {
		t.Errorf(
			"Piece mask extraction failed: got %v, want %v",
			int(move&PieceMask)>>PieceShift,
			piece,
		)
	}
	if int(move&PromotedPieceMask)>>PromotedShift != promoted {
		t.Errorf(
			"Promoted piece mask extraction failed: got %v, want %v",
			int(move&PromotedPieceMask)>>PromotedShift,
			promoted,
		)
	}
	if int(move&CapturedPieceMask)>>CapturedPieceShift != captured {
		t.Errorf(
			"Captured piece mask extraction failed: got %v, want %v",
			int(move&CapturedPieceMask)>>CapturedPieceShift,
			captured,
		)
	}
	if int(move&CaptureFlagMask)>>CaptureFlagShift != captureFlag {
		t.Errorf(
			"Capture flag mask extraction failed: got %v, want %v",
			int(move&CaptureFlagMask)>>CaptureFlagShift,
			captureFlag,
		)
	}
	if int(move&DoublePushMask)>>DoublePushFlagShift != doublePush {
		t.Errorf(
			"Double push mask extraction failed: got %v, want %v",
			int(move&DoublePushMask)>>DoublePushFlagShift,
			doublePush,
		)
	}
	if int(move&EnPasssantMask)>>EnPassantShift != enpassant {
		t.Errorf(
			"En passant mask extraction failed: got %v, want %v",
			int(move&EnPasssantMask)>>EnPassantShift,
			enpassant,
		)
	}
	if int(move&CastlingFlagMask)>>CastlingFlagShift != castlingFlag {
		t.Errorf(
			"Castling flag mask extraction failed: got %v, want %v",
			int(move&CastlingFlagMask)>>CastlingFlagShift,
			castlingFlag,
		)
	}
	if int(move&CastlingTypeMask)>>CastlingTypeShift != castlingType {
		t.Errorf(
			"Castling type mask extraction failed: got %v, want %v",
			int(move&CastlingTypeMask)>>CastlingTypeShift,
			castlingType,
		)
	}
}

// This test ensures that maximum values for each field don't interfere with each other
func TestBoundaryValues(t *testing.T) {
	// Maximum possible values for each field
	source := 63      // Maximum 6-bit value
	target := 63      // Maximum 6-bit value
	piece := 15       // Maximum 4-bit value
	promoted := 15    // Maximum 4-bit value
	captured := 15    // Maximum 4-bit value
	captureFlag := 1  // 1-bit value
	doublePush := 1   // 1-bit value
	enpassant := 1    // 1-bit value
	castlingFlag := 1 // 1-bit value
	castlingType := 3 // 2-bit value

	move := EncodeMove(
		source, target, piece, promoted, captured,
		captureFlag, doublePush, enpassant, castlingFlag, castlingType,
	)

	// Verify all fields are correctly preserved
	if got := move.GetSourceSquare(); got != source {
		t.Errorf("GetSourceSquare() with max values = %v, want %v", got, source)
	}
	if got := move.GetTargetSquare(); got != target {
		t.Errorf("GetTargetSquare() with max values = %v, want %v", got, target)
	}
	if got := move.GetMovingPiece(); got != piece {
		t.Errorf("GetMovingPiece() with max values = %v, want %v", got, piece)
	}
	if got := move.GetPromotedPiece(); got != promoted {
		t.Errorf("GetPromotedPiece() with max values = %v, want %v", got, promoted)
	}
	if got := move.GetCapturedPiece(); got != captured {
		t.Errorf("GetCapturedPiece() with max values = %v, want %v", got, captured)
	}
	if !move.IsCapture() {
		t.Errorf("IsCapture() with max values = %v, want true", move.IsCapture())
	}
	if !move.IsDoublePush() {
		t.Errorf("IsDoublePush() with max values = %v, want true", move.IsDoublePush())
	}
	if !move.IsEnPassant() {
		t.Errorf("IsEnPassant() with max values = %v, want true", move.IsEnPassant())
	}
	if !move.IsCastle() {
		t.Errorf("IsCastle() with max values = %v, want true", move.IsCastle())
	}
	if got := move.GetCastleType(); got != CastleType(castlingType) {
		t.Errorf("GetCastleType() with max values = %v, want %v", got, CastleType(castlingType))
	}
}
