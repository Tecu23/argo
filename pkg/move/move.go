// Package move contains the move and move list representation with
// all move helper functions. Move is represented as 32 bit unsigned integer
package move

import (
	"fmt"
	"strings"
	"unicode"

	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

/*

   Binary move bits representation and corresponding hexadecimal constants:

   Bits layout (from LSB to MSB):
   -----------------------------------------------------------------------------
   | Bits 0-5   | Bits 6-11  | Bits 12-15  | Bits 16-19  | Bits 20-23  | Bits 24-31  |
   -----------------------------------------------------------------------------
   | from square| to square  | moving piece| move type   | captured piece| score info  |
   -----------------------------------------------------------------------------


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
	SourceMask        = 0x0000003F // Mask for source square (6 bits)
	TargetMask        = 0x00000FC0 // Mask for target square (6 bits)
	MovingPieceMask   = 0x0000F000 // Mask for moving piece (4 bits)
	MoveTypeMask      = 0x000F0000 // Mask for the move type information (4 bits)
	CapturedPieceMask = 0x00F00000 // Mask for the captured piece (4 bits)

	// Bit Masks for individual move type flags
	PromotionMask = 0x8 // Bit flag for promotion moves within the move type field
	CaptureMask   = 0x4 // Bit flag for capture move within the move type field
	SpecialMask   = 0x3 // Additional bits for special move types

	// Bit shift constants for each move component
	SourceShift        = 0  // No shift needed for source square (bits 0-5)
	TargetShift        = 6  // Target square starts at bit 6
	MovingPieceShift   = 12 // Moving piece start at bit 12
	MoveTypeShift      = 16 // Move type information starts at bit 16
	CapturedPieceShift = 20 // Captured piece starts at bit 20
	ScoreInfoShift     = 24 // Score information starts at bit 24

	NoMove = Move(0) // A constant representing a null or no move.
)

// Type keeps the type of move (castle, promotion, enpassant capture...)
type Type uint8

// Score keeps the score of each move associated with a move.
type Score uint32

// Different move types are defined as constants. The numeric values are used when encoding
// the move type bits in the 32-bit move representation.
const (
	Quiet                  Type = iota // 0: quiet (non-capturing, non-special) move
	DoublePawnPush                     // 1: double pawn push (pawn advances two squares)
	KingCastle                         // 2: kingside castling move
	QueenCastle                        // 3: queenside castling move
	Capture                            // 4: capture move (non-en passant)
	EnPassant                          // 5: en passant capture move
	_                                  // skip (6) unused
	_                                  // skip (7) unused
	KnightPromotion                    // 8: knight promotion (non-capture)
	BishopPromotion                    // 9: bishop promotion (non-capture)
	RookPromotion                      // 10: rook promotion (non-capture)
	QueenPromotion                     // 11: queen promotion (non-capture)
	KnightPromotionCapture             // 12: knight promotion with capture
	BishopPromotionCapture             // 13: bishop promotion with capture
	RookPromotionCapture               // 14: rook promotion with capture
	QueenPromotionCapture              // 15: queen promotion with capture
)

// Move is a 32-bit unsigned integer that encodes all details of a chess move.
type Move uint32

// EncodeMove creates an encoded move from its individual components:
// - source: index of the source square (0-63)
// - target: index of the target square (0-63)
// - piece: the piece being moved (with additional color/type encoding)
// - t: the move type (using the predefined Type constants)
// - captured: the piece captured, if any (0 if no capture)
// The function packs each component into the proper bit field and returns the encoded move.
func EncodeMove(source, target, piece int, t Type, captured int) Move {
	move := Move(
		(source) | (target << TargetShift) | (piece << MovingPieceShift) |
			(int(t) << MoveTypeShift) | (captured << CapturedPieceShift),
	)
	return move
}

// GetSourceSquare extracts the source square index from the move by
// applying the SourceMask.
func (m Move) GetSourceSquare() int {
	return int(m & SourceMask)
}

// GetTargetSquare extracts the target square index by applying the
// TargetMask and then shifting right.
func (m Move) GetTargetSquare() int {
	return int(m&TargetMask) >> TargetShift
}

// GetMovingPiece extracts the encoded moving piece information by applying
// the MovingPieceMask and shifting right by MovingPieceShift.
func (m Move) GetMovingPiece() int {
	return int(m&MovingPieceMask) >> MovingPieceShift
}

// GetMovingPieceType extracts the type (e.g., pawn, knight, etc.) of the moving piece.
// It assumes that the piece type is stored in the lower 3 bits of the moving piece field.
func (m Move) GetMovingPieceType() int {
	// return int(m>>MovingPieceShift) & int(createMask(3))
	return util.GetPieceType(m.GetMovingPiece())
}

// GetMovingPieceColor extracts the color information from the moving piece field.
// It uses the bit that represents the color (e.g., white or black) – here assumed to be bit 3.
func (m Move) GetMovingPieceColor() int {
	// return int((m>>MovingPieceShift)&0x8) >> 3
	return util.GetPieceColor(m.GetMovingPiece())
}

// GetCapturedPiece extracts the captured piece (if any) by applying the CapturedPieceMask
// and shifting the result to the right by CapturedPieceShift.
func (m Move) GetCapturedPiece() int {
	return int(m&CapturedPieceMask) >> CapturedPieceShift
}

// GetCapturedPieceType extracts the type of the captured piece (e.g., pawn, knight, etc.).
// It assumes that the piece type is stored in the lower 3 bits of the captured piece field.
func (m Move) GetCapturedPieceType() int {
	// return int(m>>CapturedPieceShift) & int(createMask(3))
	return util.GetPieceType(m.GetCapturedPiece())
}

// IsPromotion checks if the move is a promotion by testing the promotion flag bit in the move type field.
func (m Move) IsPromotion() bool {
	return (m & 0x80000) != 0
}

// GetPromotionPiece calculates the promoted piece by combining the promotion information from the move type
// bits with the moving piece data. It adds an offset (+1) to correctly map to the piece value.
func (m Move) GetPromotionPiece() int {
	return int((m&0x30000)>>MoveTypeShift) + m.GetMovingPiece() + 1
}

// GetPromotionPieceType returns the type of piece the pawn is promoted to by isolating the promotion bits
// and applying the proper offset.
func (m Move) GetPromotionPieceType() int {
	// return int((m&0x30000)>>MoveTypeShift) + 1
	return util.GetPieceType(m.GetPromotionPiece())
}

// GetMoveType extracts the move type from the encoded move.
// Note: It multiplies the move by MoveTypeMask and then shifts right by MoveTypeShift,
// which is intended to isolate the move type bits.
func (m Move) GetMoveType() Type {
	return Type(m & MoveTypeMask >> MoveTypeShift)
}

// IsCapture checks whether the move is a capture by testing the capture flag bit in the move type field.
func (m Move) IsCapture() bool {
	return int(m&0x40000) != 0
}

// IsDoublePawnPush checks if the move corresponds to a pawn moving two squares forward.
func (m Move) IsDoublePawnPush() bool {
	return m.GetMoveType() == DoublePawnPush
}

// IsCastle returns true if the move is a castling move (either kingside or queenside).
func (m Move) IsCastle() bool {
	t := m.GetMoveType()
	return t == KingCastle || t == QueenCastle
}

// IsQueenCastle specifically checks for a queenside castling move.
func (m Move) IsQueenCastle() bool {
	return m.GetMoveType() == QueenCastle
}

// IsEnPassant checks if the move is an en passant capture.
func (m Move) IsEnPassant() bool {
	return m.GetMoveType() == EnPassant
}

// SetScore encodes a move's score into the upper score bits. It first clears the score portion
// and then inserts the new score shifted left by ScoreInfoShift.
func (m *Move) SetScore(score int) {
	scoreMask := createMask(8) << ScoreInfoShift
	*m = Move(uint32(*m) & ^scoreMask)  // Clear score bits
	*m |= Move(score << ScoreInfoShift) // Set new score bits
}

// GetScore extracts the score information from the move by shifting right by ScoreInfoShift.
func (m Move) GetScore() int {
	return int(m >> ScoreInfoShift)
}

// String returns the move in algebraic notation (e.g. "e2e4").
// For promotions, it appends the promoted piece's letter (lowercase).
func (m Move) String() string {
	if m == NoMove {
		return "0000"
	}

	sPromotion := ""
	// Determine if the move includes a promotion and select the appropriate piece letter.
	prom := m.GetPromotionPiece()
	if m.IsPromotion() {
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
		util.Sq2Fen[m.GetSourceSquare()], // Convert source square index to algebraic notation
		util.Sq2Fen[m.GetTargetSquare()], // Convert target square index to algebraic notation
		sPromotion,                       // Append promotion piece if applicable
	)
}

// PrintMove prints a detailed description of the move, including the source and target squares,
// promoted piece (if any), captured piece (if any), and flags indicating special move types.
func (m Move) PrintMove() {
	// Print the move's start and end squares in algebraic notation.
	fmt.Printf(
		"%s%s",
		util.Sq2Fen[m.GetSourceSquare()],
		util.Sq2Fen[m.GetTargetSquare()],
	)

	// If the move is a promotion, print the promoted piece in lowercase.
	if m.IsPromotion() {
		fmt.Printf("%c ", unicode.ToLower(rune(util.PieceIdentifier[m.GetPromotionPiece()])))
	} else {
		fmt.Printf("  ")
	}

	// If a piece was captured, print the captured piece in lowercase.
	if m.GetCapturedPiece() != 0 {
		fmt.Printf("%c ", unicode.ToLower(rune(util.PieceIdentifier[m.GetCapturedPiece()])))
	} else {
		fmt.Printf("  ")
	}

	// Print the moving piece and flags for capture, double pawn push, en passant, and castling.
	fmt.Printf("   %c ", util.PieceIdentifier[m.GetMovingPiece()])
	fmt.Printf("       %v ", m.IsCapture())
	fmt.Printf("        %v ", m.IsDoublePawnPush())
	fmt.Printf("        %v ", m.IsEnPassant())
	fmt.Printf("         %v\n", m.IsCastle())
}

// SAN returns the move in Standard Algebraic Notation (SAN).
// It handles castling, captures, pawn moves (with promotion and en passant) as well as moves
// by other pieces, assembling the notation string accordingly.
func (m Move) SAN() string {
	t := m.GetMoveType()
	pcType := m.GetMovingPieceType()
	from := m.GetSourceSquare()
	to := m.GetTargetSquare()

	// Handle castling moves first.
	if m.IsCastle() {
		if t == QueenCastle {
			return "O-O-O"
		}
		return "O-O"
	}

	prefix, midfix, postfix := "", "", ""

	// If the move is a capture, include an 'x' in the notation.
	if m.IsCapture() {
		midfix = "x"
	}

	// Convert the target square index to algebraic notation.
	postfix = util.Sq2Fen[to]

	// Handle pawn moves separately, as they have special notation rules.
	if pcType == Pawn {
		// For pawn captures, include the originating file letter.
		if m.IsCapture() {
			file := from % 8
			prefix = string(util.FileIdentifier[file])
		}

		// Append "e.p." if the move is an en passant capture.
		if m.IsEnPassant() {
			postfix += " e.p."
		}

		// Append promotion notation if the pawn is being promoted.
		if m.IsPromotion() {
			postfix += "="
			postfix += strings.ToUpper(string(util.PieceIdentifier[m.GetPromotionPiece()]))
		}
	} else {
		// For non-pawn moves, prefix the piece letter in uppercase.
		prefix = strings.ToUpper(string(util.PieceIdentifier[pcType]))
	}

	return prefix + midfix + postfix
}
