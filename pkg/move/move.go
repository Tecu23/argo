// Package move contains the move and move list representation and all move helper functions.
// Move is represented as 64 bit unsigned integer where some bits represent some part of the move
// The first 6 bits keep the source square, the next 6 bits keep the target square and so on...
package move

import (
	"fmt"
	"unicode"

	"github.com/Tecu23/argov2/pkg/util"
)

/*
         binary move bits representaion                    hexadecimal constants

   0000 0000 0000 0000 0011 1111 source square              0x3f
   0000 0000 0000 1111 1100 0000 target square              0xfc0
   0000 0000 1111 0000 0000 0000 piece                      0xf000
   0000 1111 0000 0000 0000 0000 promoted piece             0xf0000
   0001 0000 0000 0000 0000 0000 capture flag               0x100000
   0010 0000 0000 0000 0000 0000 double push flag           0x200000
   0100 0000 0000 0000 0000 0000 enpassant capture flag     0x400000
   1000 0000 0000 0000 0000 0000 castling flag              0x800000
*/

// Move encoding uses a 64-bit integer to store source, target, piece, promoted piece,
// and flags for capture, double push, en passant, and castling.
//
// Bit layout commented above shows how each part is stored in the move integer.
const (
	SourceMask     = 0x3f
	TargetMask     = 0xfc0
	PieceMask      = 0xf000
	PromotedMask   = 0xf0000
	CaptureMask    = 0x100000
	DoublePushMask = 0x200000
	EnpassantMask  = 0x400000
	CastlingMask   = 0x800000

	SourceShift     = 0
	TargetShift     = 6
	PieceShift      = 12
	PromotedShift   = 16
	CaptureShift    = 20
	DoublePushShift = 21
	EnpassantShift  = 22
	CastlingShift   = 23

	NoMove = Move(0)
)

// Move is a 64-bit unsigned integer that encodes a chess move.
type Move uint64

// EncodeMove creates a move from components: source, target, piece, promoted piece, capture flag,
// double-push flag, en passant flag, and castling flag.
func EncodeMove(
	source, target, piece, promoted, capture, doublePush, enpassant, castling int,
) Move {
	move := Move(
		(source) | (target << TargetShift) | (piece << PieceShift) |
			(promoted << PromotedShift) | (capture << CaptureShift) |
			(doublePush << DoublePushShift) | (enpassant << EnpassantShift) |
			(castling << CastlingShift),
	)

	return move
}

// GetSource should retrieve the source square of a move
func (m Move) GetSource() int {
	return int(m & SourceMask)
}

// GetTarget should retrieve the target square of a move
func (m Move) GetTarget() int {
	return int(m&TargetMask) >> TargetShift
}

// GetPiece should retrieve the piece that is moved
func (m Move) GetPiece() int {
	return int(m&PieceMask) >> PieceShift
}

// GetPromoted should retrieve the promoted piece if it exists
func (m Move) GetPromoted() int {
	return int(m&PromotedMask) >> PromotedShift
}

// GetCapture should retrieve the capture flag
func (m Move) GetCapture() int {
	return int(m&CaptureMask) >> CaptureShift
}

// GetDoublePush should retrieve the double push flag
func (m Move) GetDoublePush() int {
	return int(m&DoublePushMask) >> DoublePushShift
}

// GetEnpassant should retrieve the en passant flag
func (m Move) GetEnpassant() int {
	return int(m&EnpassantMask) >> EnpassantShift
}

// GetCastling should retrieve the castling flah
func (m Move) GetCastling() int {
	return int(m&CastlingMask) >> CastlingShift
}

// String prints the move in algebraic notation (e.g. "e2e4").
func (m Move) String() string {
	return fmt.Sprintf(
		"%s%s",
		util.Sq2Fen[m.GetSource()],
		util.Sq2Fen[m.GetTarget()],
	)
}

// PrintMove prints a detailed move description including promoted piece, capture, etc.
func (m Move) PrintMove() {
	fmt.Printf(
		"%s%s",
		util.Sq2Fen[m.GetSource()],
		util.Sq2Fen[m.GetTarget()],
	)

	if m.GetPromoted() != 0 {
		fmt.Printf("%c ", unicode.ToLower(rune(util.ASCIIPieces[m.GetPromoted()])))
	} else {
		fmt.Printf("  ")
	}

	fmt.Printf("   %c ", util.ASCIIPieces[m.GetPiece()])
	fmt.Printf("       %d ", m.GetCapture())
	fmt.Printf("        %d ", m.GetDoublePush())
	fmt.Printf("        %d ", m.GetEnpassant())
	fmt.Printf("         %d\n", m.GetCastling())
}
