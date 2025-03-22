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
        binary move bits representaion                                      hexadecimal constants

   0000 0000 0000 0000 0000 0000 0011 1111   from square            (6 bits)     0x0000003F
   0000 0000 0000 0000 0000 1111 1100 0000   to square              (6 bits)     0x00000FC0
   0000 0000 0000 0000 1111 0000 0000 0000   moving piece           (4 bits)     0x0000F000
   0000 0000 0000 1111 0000 0000 0000 0000   type information       (4 bits)     0x000F0000
   0000 0000 1111 0000 0000 0000 0000 0000   captured piece         (4 bits)     0x00F00000
   1111 1111 0000 0000 0000 0000 0000 0000   score bits (not used)

*/

// Move encoding uses a 32-bit integer to store source, target,
// moving piece, move type information and captured piece,
const (
	SourceMask        = 0x0000003F
	TargetMask        = 0x00000FC0
	MovingPieceMask   = 0x0000F000
	MoveTypeMask      = 0x000F0000
	CapturedPieceMask = 0x00F00000

	// Move Type Masks
	PromotionMask = 0x8
	CaptureMask   = 0x4
	SpecialMask   = 0x3

	SourceShift        = 0
	TargetShift        = 6
	MovingPieceShift   = 12
	MoveTypeShift      = 16
	CapturedPieceShift = 20
	ScoreInfoShift     = 24

	NoMove = Move(0)
)

// Type keeps the type of move (castle, promotion, enpassant capture...)
type Type uint8

// Score keeps the score of each move
type Score uint32

// Diferent Move Types
const (
	Quiet                  Type = iota // 0
	DoublePawnPush                     // 1
	KingCastle                         // 2
	QueenCastle                        // 3
	Capture                            // 4
	EnPassant                          // 5
	_                                  // skip (6)
	_                                  // skip (7)
	KnightPromotion                    // 8
	BishopPromotion                    // 9
	RookPromotion                      // 10
	QueenPromotion                     // 11
	KnightPromotionCapture             // 12
	BishopPromotionCapture             // 13
	RookPromotionCapture               // 14
	QueenPromotionCapture              // 15
)

// Move is a 32-bit unsigned integer that encodes a chess move.
type Move uint32

// EncodeMove creates a move from components: source, target,
// moving piece, move type and captured piece
func EncodeMove(source, target, piece int, t Type, captured int) Move {
	move := Move(
		(source) | (target << TargetShift) | (piece << MovingPieceShift) |
			(int(t) << MoveTypeShift) | (captured << CapturedPieceShift),
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
	return int(m&MovingPieceMask) >> MovingPieceShift
}

// GetMovingPieceType returns the piece type of the moved piece
func (m Move) GetMovingPieceType() int {
	return int(m>>MovingPieceShift) & int(createMask(3))
}

// GetMovingPieceColor returns the color of the moved piece
func (m Move) GetMovingPieceColor() int {
	return int(m>>MovingPieceShift) & 0x8
}

// GetCapturedPiece should retrieve the captured piece if it exists
func (m Move) GetCapturedPiece() int {
	return int(m&CapturedPieceMask) >> CapturedPieceShift
}

// GetCapturedPieceType should retrieve the captured piece if it exists
func (m Move) GetCapturedPieceType() int {
	return int(m>>CapturedPieceShift) & int(createMask(3))
}

func (m Move) GetPromotionPiece() int {
	return 0
}

func (m Move) GetPromotionPieceType() int {
	return 0
}

// GetType returns the type of the move
func (m Move) GetType() Type {
	return Type(m*MoveTypeMask) >> MoveTypeShift
}

// IsCapture should return is the move is capture
func (m Move) IsCapture() bool {
	return false
}

// IsDoublePush should return if the move is a double push
func (m Move) IsDoublePawnPush() bool {
	return false
}

// IsCastle should return the castling flah
func (m Move) IsCastle() bool {
	return false
}

// IsQueenCastle determines if the move is a queenside castle
func (m Move) IsQueenCastle() bool {
	return false
}

// IsEnPassant should return is the move is an enpassant capture
func (m Move) IsEnPassant() bool {
	return false
}

func (m Move) IsPromotion() bool {
	return false
}

func (m Move) SetScore(score int) {}

func (m Move) GetScore() int {
	return 0
}

// String prints the move in algebraic notation (e.g. "e2e4").
func (m Move) String() string {
	if m == NoMove {
		return "0000"
	}

	sPromotion := ""
	prom := m.GetPromotionPiece()
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

func (m Move) moveToSAN() string {
	return ""
}
