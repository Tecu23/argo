package move

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq()
}

func TestEncodeMove(t *testing.T) {
	tests := []struct {
		name      string
		source    int
		target    int
		piece     int
		moveType  Type
		captured  int
		expected  Move
		expectStr string
	}{
		{
			name:      "Pawn Push e2-e4",
			source:    util.FenToSq("e2"),
			target:    util.FenToSq("e4"),
			piece:     WP,
			moveType:  DoublePawnPush,
			captured:  0,
			expectStr: "e2e4",
		},
		{
			name:      "Knight Development b1-c3",
			source:    util.FenToSq("b1"),
			target:    util.FenToSq("c3"),
			piece:     WN,
			moveType:  Quiet,
			captured:  0,
			expectStr: "b1c3",
		},
		{
			name:      "Pawn Capture e4xd5",
			source:    util.FenToSq("e4"),
			target:    util.FenToSq("d5"),
			piece:     WP,
			moveType:  Capture,
			captured:  BP,
			expectStr: "e4d5",
		},
		{
			name:      "King Castle O-O",
			source:    util.FenToSq("e1"),
			target:    util.FenToSq("g1"),
			piece:     WK,
			moveType:  KingCastle,
			captured:  0,
			expectStr: "e1g1",
		},
		{
			name:      "Queen Castle O-O-O",
			source:    util.FenToSq("e1"),
			target:    util.FenToSq("c1"),
			piece:     WK,
			moveType:  QueenCastle,
			captured:  0,
			expectStr: "e1c1",
		},
		{
			name:      "Queen Promotion",
			source:    util.FenToSq("e7"),
			target:    util.FenToSq("e8"),
			piece:     WP,
			moveType:  QueenPromotion,
			captured:  0,
			expectStr: "e7e8q",
		},
		{
			name:      "Queen Promotion with Capture",
			source:    util.FenToSq("d7"),
			target:    util.FenToSq("e8"),
			piece:     WP,
			moveType:  QueenPromotionCapture,
			captured:  BQ,
			expectStr: "d7e8q",
		},
		{
			name:      "En Passant",
			source:    util.FenToSq("e5"),
			target:    util.FenToSq("d6"),
			piece:     WP,
			moveType:  EnPassant,
			captured:  BP,
			expectStr: "e5d6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := EncodeMove(tt.source, tt.target, tt.piece, tt.moveType, tt.captured)

			// Test that all parts were encoded correctly
			assert.Equal(t, tt.source, m.GetSourceSquare())
			assert.Equal(t, tt.target, m.GetTargetSquare())
			assert.Equal(t, tt.piece, m.GetMovingPiece())
			assert.Equal(t, tt.captured, m.GetCapturedPiece())

			// Test string representation
			assert.Equal(t, tt.expectStr, m.String())
		})
	}
}

func TestMoveProperties(t *testing.T) {
	tests := []struct {
		name           string
		move           Move
		isCapture      bool
		isPromotion    bool
		isDoublePush   bool
		isCastle       bool
		isKingCastle   bool
		isQueenCastle  bool
		isEnPassant    bool
		promotionPiece int
	}{
		{
			name:           "Pawn Push",
			move:           EncodeMove(util.FenToSq("e2"), util.FenToSq("e3"), WP, Quiet, 0),
			isCapture:      false,
			isPromotion:    false,
			isDoublePush:   false,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: 0,
		},
		{
			name: "Double Pawn Push",
			move: EncodeMove(
				util.FenToSq("e2"),
				util.FenToSq("e4"),
				WP,
				DoublePawnPush,
				0,
			),
			isCapture:      false,
			isPromotion:    false,
			isDoublePush:   true,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: 0,
		},
		{
			name:           "Normal Capture",
			move:           EncodeMove(util.FenToSq("e4"), util.FenToSq("d5"), WP, Capture, BP),
			isCapture:      true,
			isPromotion:    false,
			isDoublePush:   false,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: 0,
		},
		{
			name:           "King Castle",
			move:           EncodeMove(util.FenToSq("e1"), util.FenToSq("g1"), WK, KingCastle, 0),
			isCapture:      false,
			isPromotion:    false,
			isDoublePush:   false,
			isCastle:       true,
			isKingCastle:   true,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: 0,
		},
		{
			name:           "Queen Castle",
			move:           EncodeMove(util.FenToSq("e1"), util.FenToSq("c1"), WK, QueenCastle, 0),
			isCapture:      false,
			isPromotion:    false,
			isDoublePush:   false,
			isCastle:       true,
			isKingCastle:   false,
			isQueenCastle:  true,
			isEnPassant:    false,
			promotionPiece: 0,
		},
		{
			name: "Queen Promotion",
			move: EncodeMove(
				util.FenToSq("e7"),
				util.FenToSq("e8"),
				WP,
				QueenPromotion,
				0,
			),
			isCapture:      false,
			isPromotion:    true,
			isDoublePush:   false,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: WQ,
		},
		{
			name: "Knight Promotion with Capture",
			move: EncodeMove(
				util.FenToSq("d7"),
				util.FenToSq("e8"),
				WP,
				KnightPromotionCapture,
				BR,
			),
			isCapture:      true,
			isPromotion:    true,
			isDoublePush:   false,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    false,
			promotionPiece: WN,
		},
		{
			name:           "En Passant",
			move:           EncodeMove(util.FenToSq("e5"), util.FenToSq("d6"), WP, EnPassant, BP),
			isCapture:      true,
			isPromotion:    false,
			isDoublePush:   false,
			isCastle:       false,
			isKingCastle:   false,
			isQueenCastle:  false,
			isEnPassant:    true,
			promotionPiece: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isCapture, tt.move.IsCapture())
			assert.Equal(t, tt.isPromotion, tt.move.IsPromotion())
			assert.Equal(t, tt.isDoublePush, tt.move.IsDoublePawnPush())
			assert.Equal(t, tt.isCastle, tt.move.IsCastle())
			assert.Equal(t, tt.isQueenCastle, tt.move.IsQueenCastle())
			assert.Equal(t, tt.isEnPassant, tt.move.IsEnPassant())

			if tt.isPromotion {
				assert.Equal(t, tt.promotionPiece, tt.move.GetPromotionPiece())
			}
		})
	}
}

func TestPieceInfo(t *testing.T) {
	tests := []struct {
		name              string
		move              Move
		movingPiece       int
		movingPieceType   int
		movingPieceColor  int
		capturedPiece     int
		capturedPieceType int
	}{
		{
			name: "White Pawn Move",
			move: EncodeMove(
				util.FenToSq("e2"),
				util.FenToSq("e4"),
				WP,
				DoublePawnPush,
				0,
			),
			movingPiece:       WP,
			movingPieceType:   Pawn,
			movingPieceColor:  0, // white
			capturedPiece:     0,
			capturedPieceType: 0,
		},
		{
			name:              "White Knight Captures Black Pawn",
			move:              EncodeMove(util.FenToSq("d4"), util.FenToSq("e6"), WN, Capture, BP),
			movingPiece:       WN,
			movingPieceType:   Knight,
			movingPieceColor:  0, // white
			capturedPiece:     BP,
			capturedPieceType: Pawn,
		},
		{
			name:              "Black Queen Captures White Rook",
			move:              EncodeMove(util.FenToSq("d8"), util.FenToSq("d1"), BQ, Capture, WR),
			movingPiece:       BQ,
			movingPieceType:   Queen,
			movingPieceColor:  1, // black
			capturedPiece:     WR,
			capturedPieceType: Rook,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.movingPiece, tt.move.GetMovingPiece())
			assert.Equal(t, tt.movingPieceType, tt.move.GetMovingPieceType())
			assert.Equal(t, tt.movingPieceColor, tt.move.GetMovingPieceColor())
			assert.Equal(t, tt.capturedPiece, tt.move.GetCapturedPiece())
			assert.Equal(t, tt.capturedPieceType, tt.move.GetCapturedPieceType())
		})
	}
}

func TestMoveScore(t *testing.T) {
	m := EncodeMove(util.FenToSq("e2"), util.FenToSq("e4"), WP, DoublePawnPush, 0)

	// Test default score is 0
	assert.Equal(t, 0, m.GetScore())

	// Test setting score
	scores := []int{100, 200, 75, 24}
	for _, score := range scores {
		m.SetScore(score)
		assert.Equal(t, score, m.GetScore())
	}
}

func TestSANNotation(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		expected string
	}{
		{
			name:     "Pawn Push",
			move:     EncodeMove(util.FenToSq("e2"), util.FenToSq("e4"), WP, DoublePawnPush, 0),
			expected: "e4",
		},
		{
			name:     "Knight Development",
			move:     EncodeMove(util.FenToSq("g1"), util.FenToSq("f3"), WN, Quiet, 0),
			expected: "Nf3",
		},
		{
			name:     "Pawn Capture",
			move:     EncodeMove(util.FenToSq("e4"), util.FenToSq("d5"), WP, Capture, BP),
			expected: "exd5",
		},
		{
			name:     "Queen Capture",
			move:     EncodeMove(util.FenToSq("d1"), util.FenToSq("d8"), WQ, Capture, BQ),
			expected: "Qxd8",
		},
		{
			name:     "King Castle",
			move:     EncodeMove(util.FenToSq("e1"), util.FenToSq("g1"), WK, KingCastle, 0),
			expected: "O-O",
		},
		{
			name:     "Queen Castle",
			move:     EncodeMove(util.FenToSq("e1"), util.FenToSq("c1"), WK, QueenCastle, 0),
			expected: "O-O-O",
		},
		{
			name:     "Queen Promotion",
			move:     EncodeMove(util.FenToSq("e7"), util.FenToSq("e8"), WP, QueenPromotion, 0),
			expected: "e8=Q",
		},
		{
			name: "Knight Promotion with Capture",
			move: EncodeMove(
				util.FenToSq("d7"),
				util.FenToSq("e8"),
				WP,
				KnightPromotionCapture,
				BR,
			),
			expected: "dxe8=N",
		},
		{
			name:     "En Passant",
			move:     EncodeMove(util.FenToSq("e5"), util.FenToSq("d6"), WP, EnPassant, BP),
			expected: "exd6 e.p.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.move.SAN())
		})
	}
}

func TestNoMove(t *testing.T) {
	m := NoMove

	assert.Equal(t, "0000", m.String())
	assert.Equal(t, 0, m.GetSourceSquare())
	assert.Equal(t, 0, m.GetTargetSquare())
	assert.Equal(t, 0, m.GetMovingPiece())
	assert.Equal(t, 0, m.GetCapturedPiece())
	assert.Equal(t, 0, m.GetScore())
	assert.False(t, m.IsCapture())
	assert.False(t, m.IsPromotion())
	assert.False(t, m.IsDoublePawnPush())
	assert.False(t, m.IsCastle())
	assert.False(t, m.IsQueenCastle())
	assert.False(t, m.IsEnPassant())
}

func TestEdgeCases(t *testing.T) {
	// Test with extreme values that might overflow bit masks
	source := 63 // Max square index (h8)
	target := 0  // Min square index (a1)
	piece := WK  // King
	moveType := Quiet
	captured := BQ // Queen

	m := EncodeMove(source, target, piece, moveType, captured)

	assert.Equal(t, source, m.GetSourceSquare())
	assert.Equal(t, target, m.GetTargetSquare())
	assert.Equal(t, piece, m.GetMovingPiece())
	assert.Equal(t, captured, m.GetCapturedPiece())
}
