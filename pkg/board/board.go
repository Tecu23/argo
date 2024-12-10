// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/util"
)

// Variables to keep the move flags
const (
	AllMoves     = 0
	OnlyCaptures = 1
)

// Board represents the state of a chess board, including piece placement,
// occupancy bitboards, side to move, en passant square, rule 50 counter,
// and castling rights.
type Board struct {
	Bitboards   [12]bitboard.Bitboard
	Occupancies [3]bitboard.Bitboard
	Side        color.Color
	EnPassant   int
	Rule50      uint8
	Castlings
}

// Reset restores the board to an initial empty state and sets defaults.
func (b *Board) Reset() {
	b.Side = color.WHITE
	b.EnPassant = -1
	b.Castlings = 0

	for i := 0; i < 12; i++ {
		b.Bitboards[i] = 0
	}

	for i := 0; i < 3; i++ {
		b.Occupancies[i] = 0
	}
}

// CopyBoard creates a copy of the board's current state.
func (b Board) CopyBoard() Board {
	boardCopy := b
	return boardCopy
}

// TakeBack restores the board state from a previously saved copy.
// Useful for undoing moves (like a "pop" from a stack).
func (b *Board) TakeBack(cpy Board) {
	*b = cpy
}

// SetSq places a piece on a given square (or clears it if piece == constants.Empty).
// It updates both piece bitboards and occupancy bitboards accordingly.
func (b *Board) SetSq(piece, sq int) {
	pieceColor := util.PcColor(piece)

	// If there is a piece on the square, remove it first
	if b.Occupancies[color.BOTH].Test(sq) {
		for p := constants.WP; p <= constants.BK; p++ {
			if b.Bitboards[p].Test(sq) {
				b.Bitboards[p].Clear(sq)
			}
		}

		b.Occupancies[color.BOTH].Clear(sq)
		b.Occupancies[color.WHITE].Clear(sq)
		b.Occupancies[color.BLACK].Clear(sq)
	}

	// If setting an empty piece, we just return after clearing.
	if piece == constants.Empty {
		return
	}

	b.Bitboards[piece].Set(sq)

	if pieceColor == color.WHITE {
		b.Occupancies[color.WHITE].Set(sq)
	} else {
		b.Occupancies[color.BLACK].Set(sq)
	}

	// Update BOTH occupancy as union of WHITE and BLACK
	b.Occupancies[color.BOTH] |= b.Occupancies[color.WHITE]
	b.Occupancies[color.BOTH] |= b.Occupancies[color.BLACK]
}

// IsSquareAttacked checks if a given square is attacked by the specified side (WHITE or BLACK).
// It uses various precomputed attack arrays to check if any piece of that color can attack 'sq'.
func (b *Board) IsSquareAttacked(sq int, side color.Color) bool {
	if side == color.WHITE {
		// Check White's pawns, knights, king, bishops, rooks, queen attacks
		if attacks.PawnAttacks[color.BLACK][sq]&b.Bitboards[constants.WP] != 0 {
			return true
		}

		if attacks.KnightAttacks[sq]&b.Bitboards[constants.WN] != 0 {
			return true
		}

		if attacks.KingAttacks[sq]&b.Bitboards[constants.WK] != 0 {
			return true
		}

		bishopAttacks := attacks.GetBishopAttacks(sq, b.Occupancies[color.BOTH])

		if bishopAttacks&b.Bitboards[constants.WB] != 0 {
			return true
		}
		rookAttacks := attacks.GetRookAttacks(sq, b.Occupancies[color.BOTH])

		if rookAttacks&b.Bitboards[constants.WR] != 0 {
			return true
		}
		queenAttacks := attacks.GetQueenAttacks(sq, b.Occupancies[color.BOTH])

		if queenAttacks&b.Bitboards[constants.WQ] != 0 {
			return true
		}
	} else {
		// Check Black's pawns, knights, king, bishops, rooks, queen attacks
		if attacks.PawnAttacks[color.WHITE][sq]&b.Bitboards[constants.BP] != 0 {
			return true
		}

		if attacks.KnightAttacks[sq]&b.Bitboards[constants.BN] != 0 {
			return true
		}

		if attacks.KingAttacks[sq]&b.Bitboards[constants.BK] != 0 {
			return true
		}

		bishopAttacks := attacks.GetBishopAttacks(sq, b.Occupancies[color.BOTH])

		if bishopAttacks&b.Bitboards[constants.BB] != 0 {
			return true
		}
		rookAttacks := attacks.GetRookAttacks(sq, b.Occupancies[color.BOTH])

		if rookAttacks&b.Bitboards[constants.BR] != 0 {
			return true
		}
		queenAttacks := attacks.GetQueenAttacks(sq, b.Occupancies[color.BOTH])

		if queenAttacks&b.Bitboards[constants.BQ] != 0 {
			return true
		}
	}
	return false
}

// MakeMove attempts to make a move on the board. It updates board state (bitboards,
// occupancy, castling, en passant) and returns true if successful. If the move leaves
// own king in check, it reverts and returns false. moveFlag determines whether to consider
// only captures or all moves.
func (b *Board) MakeMove(m move.Move, moveFlag int) bool {
	// If moveFlag == OnlyCaptures, we only proceed if this move is a capture move.
	// Otherwise, we handle all moves.
	// The code here includes logic for en passant, castling, double pushes, promotions,
	// and ensures legality by checking king safety.
	// Due to complexity, only high-level comments are given here.
	if moveFlag == AllMoves {
		// preserve board state
		copyB := b.CopyBoard()

		src := m.GetSource()
		tgt := m.GetTarget()
		pc := m.GetPiece()
		clr := util.PcColor(pc)
		prom := m.GetPromoted()
		// capt := m.GetCapture()
		dblPwn := m.GetDoublePush()
		ep := m.GetEnpassant()
		cast := m.GetCastling()

		b.EnPassant = -1

		// Handle en passant
		if ep != 0 {
			b.SetSq(constants.Empty, src)
			if clr == color.WHITE {
				b.SetSq(constants.Empty, tgt+constants.S)
			} else {
				b.SetSq(constants.Empty, tgt+constants.N)
			}

			b.SetSq(pc, tgt)

			b.Side = b.Side.Opp()

			// Check for checks to ensure move legality
			var kingPos int
			if b.Side == color.WHITE {
				if b.Bitboards[constants.BK] == 0 {
					b.TakeBack(copyB)
					return false
				}
				kingPos = b.Bitboards[constants.BK].FirstOne()
			} else {
				if b.Bitboards[constants.WK] == 0 {
					b.TakeBack(copyB)
					return false
				}
				kingPos = b.Bitboards[constants.WK].FirstOne()
			}
			if b.IsSquareAttacked(kingPos, b.Side) {
				// take back
				b.TakeBack(copyB)
				return false
			}
			if b.Side == color.WHITE {
				b.Bitboards[constants.BK].Set(kingPos)
			} else {
				b.Bitboards[constants.WK].Set(kingPos)
			}

			return true
		}

		// Handle castling
		if cast != 0 {
			switch tgt {
			// WHITE Short Castle
			case constants.G1:
				b.SetSq(constants.Empty, constants.H1)
				b.SetSq(constants.WR, constants.F1)
			// WHITE Long Castle
			case constants.C1:
				b.SetSq(constants.Empty, constants.A1)
				b.SetSq(constants.WR, constants.D1)
			// BLACK Short Castle
			case constants.G8:
				b.SetSq(constants.Empty, constants.H8)
				b.SetSq(constants.BR, constants.F8)
			// BLACK Long Castle
			case constants.C8:
				b.SetSq(constants.Empty, constants.A8)
				b.SetSq(constants.BR, constants.D8)
			}
		}

		// Double push pawn update
		if dblPwn != 0 {
			if clr == color.WHITE {
				b.EnPassant = src + constants.N
				// b.Key ^= EnpassantKeys[src+constants.N]
			} else {
				b.EnPassant = src + constants.S
				// b.Key ^= EnpassantKeys[src+constants.S]
			}
		}

		b.SetSq(constants.Empty, src)

		if prom != 0 {
			b.SetSq(prom, tgt)
		} else {
			b.SetSq(pc, tgt)
		}

		// Update castling rights if necessary
		b.Castlings &= Castlings(CastlingRights[src])
		b.Castlings &= Castlings(CastlingRights[tgt])

		// change side
		b.Side = b.Side.Opp()

		// Check if own king is in check after the move
		var kingPos int
		if b.Side == color.WHITE {
			if b.Bitboards[constants.BK] == 0 {
				b.TakeBack(copyB)
				return false
			}
			kingPos = b.Bitboards[constants.BK].FirstOne()
		} else {
			if b.Bitboards[constants.WK] == 0 {
				b.TakeBack(copyB)
				return false
			}
			kingPos = b.Bitboards[constants.WK].FirstOne()
		}
		if b.IsSquareAttacked(kingPos, b.Side) {
			// take back
			b.TakeBack(copyB)
			return false
		}
		if b.Side == color.WHITE {
			b.Bitboards[constants.BK].Set(kingPos)
		} else {
			b.Bitboards[constants.WK].Set(kingPos)
		}
	} else { // capture moves
		if m.GetCapture() != 0 {
			return b.MakeMove(m, AllMoves)
		}
		return false // 0 means don't make it
	}
	return true
}

// ParseMove takes a move string (like "e7e8q") and returns the corresponding Move object if valid.
// It generates all moves, finds the one matching this string, and returns it. If not found, returns NoMove.
func (b *Board) ParseMove(moveString string) move.Move {
	var moves move.Movelist
	b.GenerateMoves(&moves)

	src := util.Fen2Sq[moveString[:2]]
	tgt := util.Fen2Sq[moveString[2:4]]

	for cnt := 0; cnt < len(moves); cnt++ {
		mv := moves[cnt]

		if mv.GetSource() == src && mv.GetTarget() == tgt {
			prom := mv.GetPromoted()

			if prom != 0 {
				// Check if promotion matches requested piece
				if (prom == constants.WQ || prom == constants.BQ) && moveString[4] == 'q' {
					return mv
				}
				if (prom == constants.WR || prom == constants.BR) && moveString[4] == 'r' {
					return mv
				}
				if (prom == constants.WB || prom == constants.BB) && moveString[4] == 'b' {
					return mv
				}
				if (prom == constants.WN || prom == constants.BN) && moveString[4] == 'n' {
					return mv
				}
				continue // continue the loop on wrong promotions
			}
			// If no promotion needed or matches, return this move
			return mv
		}
	}
	return move.NoMove
}

// PrintBoard prints the board state in a human-readable format with ranks and files.
func (b Board) PrintBoard() {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			if file == 0 {
				fmt.Printf("%d  ", rank+1)
			}
			piece := -1

			// loop over all piece bitboards
			for bb := constants.WP; bb <= constants.BK; bb++ {
				if b.Bitboards[bb].Test(rank*8 + file) {
					piece = bb
				}
			}

			if piece == -1 {
				fmt.Printf(" %c", '.')
			} else {
				fmt.Printf(" %c", util.ASCIIPieces[piece])
			}

		}
		fmt.Println()
	}

	fmt.Printf("\n    a b c d e f g h\n\n")

	fmt.Printf("   Side:          %s\n", b.Side.String())
	fmt.Printf("   Enpassant:     %s\n", util.Sq2Fen[b.EnPassant])
	fmt.Printf("   Half Moves:    %d\n", b.Rule50)
	fmt.Printf("   Castling:   %s\n\n", b.Castlings.String())
	// fmt.Printf(" HashKey: 0x%X\n\n", b.Key)
}

// ParseFEN sets the board state according to a given FEN string.
// It places pieces, sets side to move, castling rights, and en passant square.
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -
func (b *Board) ParseFEN(FEN string) {
	b.Reset()

	fenIdx := 0
	sq := 0

	// Parse the ranks from top (rank=7) to bottom (rank=0)
	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {

			char := string(FEN[fenIdx])
			fenIdx++

			if char == "/" {
				continue
			}

			// If char is a digit, skip that many squares
			if i, err := strconv.Atoi(char); err == nil {
				for j := 0; j < i; j++ {
					b.SetSq(constants.Empty, sq)
					sq++
				}
				continue
			}

			// Otherwise, it should be a piece character
			if !strings.Contains(util.PcFen, char) {
				fmt.Printf("Invalid piece %s try next one", char)
				continue
			}

			b.SetSq(util.Fen2pc(char), sq)
			sq++
		}
	}

	remaining := strings.Split(strings.TrimSpace(FEN[fenIdx:]), " ")

	// Set side to move
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			b.Side = color.WHITE
		} else if remaining[0] == "b" {
			b.Side = color.BLACK
		} else {
			fmt.Printf("Remaining=%v; sq=%v;  fenIx=%v;", strings.Join(remaining, " "), sq, fenIdx)
			fmt.Printf("%s invalid side to move color", remaining[0])
			b.Side = color.WHITE
		}
	}

	// Set castling rights
	b.Castlings = 0
	if len(remaining) > 1 {
		b.Castlings = ParseCastlings(remaining[1])
	}

	// Set en passant square
	b.EnPassant = -1
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			b.EnPassant = util.Fen2Sq[remaining[2]]
		}
	}

	// Set halfmove clock (for 50-move rule)
	b.Rule50 = 0
	if len(remaining) > 3 {
		cnt, err := strconv.Atoi(remaining[3])
		if err != nil {
			b.Rule50 = 0
		}

		b.Rule50 = uint8(cnt)
	}
}
