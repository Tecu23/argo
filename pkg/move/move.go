// Package move contains the move and move list representation and all move helper functions.
// Move is represented as 64 bit unsigned integer where some bits represent some part of the move
// The first 6 bits keep the source square, the next 6 bits keep the target square and so on...
package move

import (
	"fmt"
	"unicode"

	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

/*
        binary move bits representaion                             hexadecimal constants

   0000 0000 0000 0000 0000 0000 0011 1111   source square              0x0000003F
   0000 0000 0000 0000 0000 1111 1100 0000   target square              0x00000FC0
   0000 0000 0000 0000 1111 0000 0000 0000   moving piece               0x0000F000
   0000 0000 0000 1111 0000 0000 0000 0000   promoted piece             0x000F0000
   0000 0000 1111 0000 0000 0000 0000 0000   captured piece             0x00F00000
   0000 0001 0000 0000 0000 0000 0000 0000   capture flag               0x01000000
   0000 0010 0000 0000 0000 0000 0000 0000   double push flag           0x02000000
   0000 0100 0000 0000 0000 0000 0000 0000   enpassant flag             0x04000000
   0000 1000 0000 0000 0000 0000 0000 0000   castling flag              0x08000000
   0011 0000 0000 0000 0000 0000 0000 0000   castling type              0x30000000
   1100 0000 0000 0000 0000 0000 0000 0000   unused bits                0xC0000000

*/

// Move encoding uses a 32-bit integer to store source, target, piece, promoted piece,
// capured piece  and flags for capture, double push, en passant, and castling.
const (
	SourceMask        = 0x0000003F
	TargetMask        = 0x00000FC0
	PieceMask         = 0x0000F000
	PromotedPieceMask = 0x000F0000
	CapturedPieceMask = 0x00F00000
	CaptureFlagMask   = 0x01000000
	DoublePushMask    = 0x02000000
	EnPasssantMask    = 0x04000000

	CastlingFlagMask = 0x08000000
	CastlingTypeMask = 0x30000000

	SourceShift         = 0
	TargetShift         = 6
	PieceShift          = 12
	PromotedShift       = 16
	CapturedPieceShift  = 20
	CaptureFlagShift    = 24
	DoublePushFlagShift = 25
	EnPassantShift      = 26
	CastlingFlagShift   = 27
	CastlingTypeShift   = 28

	NoMove = Move(0)
)

// CastleType are the types of castleling available on a chess board
type CastleType uint8

// Castling Types
const (
	WhiteKingCastle  CastleType = 0
	WhiteQueenCastle CastleType = 1
	BlackKingCastle  CastleType = 2
	BlackQueenCastle CastleType = 3
)

// Move is a 32-bit unsigned integer that encodes a chess move.
type Move uint32

// EncodeMove creates a move from components: source, target, piece, promoted piece, capture flag,
// double-push flag, en passant flag, and castling flag.
func EncodeMove(
	source, target, piece, promoted, captured, captureFlag, doublePush, enpassant, castlingFlag, castlingType int,
) Move {
	move := Move(
		(source) | (target << TargetShift) | (piece << PieceShift) |
			(promoted << PromotedShift) | (captured << CapturedPieceShift) |
			(captureFlag << CaptureFlagShift) | (doublePush << DoublePushFlagShift) |
			(enpassant << EnPassantShift) | (castlingFlag << CastlingFlagShift) |
			(castlingType << CastlingTypeShift),
	)

	return move
}

// GetSourceSquare should retrieve the source square of a move
func (m Move) GetSourceSquare() int {
	return int(m & SourceMask)
}

// GetTargetSquare should retrieve the target square of a move
func (m Move) GetTargetSquare() int {
	return int(m&TargetMask) >> TargetShift
}

// GetMovingPiece should retrieve the piece that is moved
func (m Move) GetMovingPiece() int {
	return int(m&PieceMask) >> PieceShift
}

// GetPromotedPiece should retrieve the promoted piece if it exists
func (m Move) GetPromotedPiece() int {
	return int(m&PromotedPieceMask) >> PromotedShift
}

// GetCapturedPiece should retrieve the captured piece if it exists
func (m Move) GetCapturedPiece() int {
	return int(m&CapturedPieceMask) >> CapturedPieceShift
}

// IsCapture should return is the move is capture
func (m Move) IsCapture() bool {
	return int(m&CaptureFlagMask)>>CaptureFlagShift != 0
}

// IsDoublePush should return if the move is a double push
func (m Move) IsDoublePush() bool {
	return int(m&DoublePushMask)>>DoublePushFlagShift != 0
}

// IsEnPassant should return is the move is an enpassant capture
func (m Move) IsEnPassant() bool {
	return int(m&EnPasssantMask)>>EnPassantShift != 0
}

// IsCastle should return the castling flah
func (m Move) IsCastle() bool {
	return int(m&CastlingFlagMask)>>CastlingFlagShift != 0
}

// GetCastleType should return the type of castling
func (m Move) GetCastleType() CastleType {
	return CastleType(int(m&CastlingTypeMask) >> CastlingTypeShift)
}

// String prints the move in algebraic notation (e.g. "e2e4").
func (m Move) String() string {
	if m == NoMove {
		return "0000"
	}

	sPromotion := ""
	prom := m.GetPromotedPiece()
	if prom != 0 {
		if prom == WQ || prom == BQ {
			sPromotion = "q"
		} else if prom == WR || prom == BR {
			sPromotion = "r"
		} else if prom == WB || prom == BB {
			sPromotion = "b"
		} else if prom == WN || prom == BN {
			sPromotion = "n"
		}
	}
	return fmt.Sprintf(
		"%s%s%s",
		util.Sq2Fen[m.GetSourceSquare()],
		util.Sq2Fen[m.GetTargetSquare()],
		sPromotion,
	)
}

// PrintMove prints a detailed move description including promoted piece, capture, etc.
func (m Move) PrintMove() {
	fmt.Printf(
		"%s%s",
		util.Sq2Fen[m.GetSourceSquare()],
		util.Sq2Fen[m.GetTargetSquare()],
	)

	if m.GetPromotedPiece() != 0 {
		fmt.Printf("%c ", unicode.ToLower(rune(util.ASCIIPieces[m.GetPromotedPiece()])))
	} else {
		fmt.Printf("  ")
	}

	if m.GetCapturedPiece() != 0 {
		fmt.Printf("%c ", unicode.ToLower(rune(util.ASCIIPieces[m.GetCapturedPiece()])))
	} else {
		fmt.Printf("  ")
	}

	fmt.Printf("   %c ", util.ASCIIPieces[m.GetMovingPiece()])
	fmt.Printf("       %v ", m.IsCapture())
	fmt.Printf("        %v ", m.IsDoublePush())
	fmt.Printf("        %v ", m.IsEnPassant())
	fmt.Printf("         %v\n", m.IsCastle())
	fmt.Printf("         %d\n", m.GetCastleType())
}

// IsQueenCastle determines if the move is a queenside castle
func (m Move) IsQueenCastle() bool {
	// First check if it's a castling move at all
	if !m.IsCastle() {
		return false
	}

	casType := m.GetCastleType()

	return casType == WhiteQueenCastle || casType == BlackQueenCastle
}
