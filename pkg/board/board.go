// Package board contains the board representation and all board helper functions.
// This package will handle move geration
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
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
	Castlings   Castlings
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

// SetSq places a piece on a given square (or clears it if piece == Empty).
// It updates both piece bitboards and occupancy bitboards accordingly.
func (b *Board) SetSq(piece, sq int) {
	pieceColor := util.PcColor(piece)

	// If there is a piece on the square, remove it first
	if b.Occupancies[color.BOTH].Test(sq) {
		for p := WP; p <= BK; p++ {
			if b.Bitboards[p].Test(sq) {
				b.Bitboards[p].Clear(sq)
			}
		}

		b.Occupancies[color.BOTH].Clear(sq)
		b.Occupancies[color.WHITE].Clear(sq)
		b.Occupancies[color.BLACK].Clear(sq)
	}

	// If setting an empty piece, we just return after clearing.
	if piece == Empty {
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
		if attacks.PawnAttacks[color.BLACK][sq]&b.Bitboards[WP] != 0 {
			return true
		}

		if attacks.KnightAttacks[sq]&b.Bitboards[WN] != 0 {
			return true
		}

		if attacks.KingAttacks[sq]&b.Bitboards[WK] != 0 {
			return true
		}

		bishopAttacks := attacks.GetBishopAttacks(sq, b.Occupancies[color.BOTH])

		if bishopAttacks&b.Bitboards[WB] != 0 {
			return true
		}
		rookAttacks := attacks.GetRookAttacks(sq, b.Occupancies[color.BOTH])

		if rookAttacks&b.Bitboards[WR] != 0 {
			return true
		}
		queenAttacks := attacks.GetQueenAttacks(sq, b.Occupancies[color.BOTH])

		if queenAttacks&b.Bitboards[WQ] != 0 {
			return true
		}
	} else {
		// Check Black's pawns, knights, king, bishops, rooks, queen attacks
		if attacks.PawnAttacks[color.WHITE][sq]&b.Bitboards[BP] != 0 {
			return true
		}

		if attacks.KnightAttacks[sq]&b.Bitboards[BN] != 0 {
			return true
		}

		if attacks.KingAttacks[sq]&b.Bitboards[BK] != 0 {
			return true
		}

		bishopAttacks := attacks.GetBishopAttacks(sq, b.Occupancies[color.BOTH])

		if bishopAttacks&b.Bitboards[BB] != 0 {
			return true
		}
		rookAttacks := attacks.GetRookAttacks(sq, b.Occupancies[color.BOTH])

		if rookAttacks&b.Bitboards[BR] != 0 {
			return true
		}
		queenAttacks := attacks.GetQueenAttacks(sq, b.Occupancies[color.BOTH])

		if queenAttacks&b.Bitboards[BQ] != 0 {
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
			b.SetSq(Empty, src)
			if clr == color.WHITE {
				b.SetSq(Empty, tgt+8)
			} else {
				b.SetSq(Empty, tgt-8)
			}

			b.SetSq(pc, tgt)

			b.Side = b.Side.Opp()

			// Check for checks to ensure move legality
			var kingPos int
			if b.Side == color.WHITE {
				if b.Bitboards[BK] == 0 {
					b.TakeBack(copyB)
					return false
				}
				kingPos = b.Bitboards[BK].FirstOne()
			} else {
				if b.Bitboards[WK] == 0 {
					b.TakeBack(copyB)
					return false
				}
				kingPos = b.Bitboards[WK].FirstOne()
			}
			if b.IsSquareAttacked(kingPos, b.Side) {
				// take back
				b.TakeBack(copyB)
				return false
			}
			if b.Side == color.WHITE {
				b.Bitboards[BK].Set(kingPos)
			} else {
				b.Bitboards[WK].Set(kingPos)
			}

			return true
		}

		// Handle castling
		if cast != 0 {
			switch tgt {
			// WHITE Short Castle
			case G1:
				b.SetSq(Empty, H1)
				b.SetSq(WR, F1)
			// WHITE Long Castle
			case C1:
				b.SetSq(Empty, A1)
				b.SetSq(WR, D1)
			// BLACK Short Castle
			case G8:
				b.SetSq(Empty, H8)
				b.SetSq(BR, F8)
			// BLACK Long Castle
			case C8:
				b.SetSq(Empty, A8)
				b.SetSq(BR, D8)
			}
		}

		// Double push pawn update
		if dblPwn != 0 {
			if clr == color.WHITE {
				b.EnPassant = src - 8
				// b.Key ^= EnpassantKeys[src+N]
			} else {
				b.EnPassant = src + 8
				// b.Key ^= EnpassantKeys[src+S]
			}
		}

		b.SetSq(Empty, src)

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
			if b.Bitboards[BK] == 0 {
				b.TakeBack(copyB)
				return false
			}
			kingPos = b.Bitboards[BK].FirstOne()
		} else {
			if b.Bitboards[WK] == 0 {
				b.TakeBack(copyB)
				return false
			}
			kingPos = b.Bitboards[WK].FirstOne()
		}
		if b.IsSquareAttacked(kingPos, b.Side) {
			// take back
			b.TakeBack(copyB)
			return false
		}
		if b.Side == color.WHITE {
			b.Bitboards[BK].Set(kingPos)
		} else {
			b.Bitboards[WK].Set(kingPos)
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
func (b *Board) ParseMove(moveString string) (Board, bool) {
	newB := b.CopyBoard()
	moves := b.GenerateMoves()

	src := util.Fen2Sq[moveString[:2]]
	tgt := util.Fen2Sq[moveString[2:4]]

	tmpMove := move.NoMove

	for cnt := 0; cnt < len(moves); cnt++ {
		mv := moves[cnt]

		if mv.GetSource() == src && mv.GetTarget() == tgt {
			prom := mv.GetPromoted()

			if prom != 0 {
				// Check if promotion matches requested piece
				if (prom == WQ || prom == BQ) && moveString[4] == 'q' {
					tmpMove = mv
					break
				}
				if (prom == WR || prom == BR) && moveString[4] == 'r' {
					tmpMove = mv
					break
				}
				if (prom == WB || prom == BB) && moveString[4] == 'b' {
					tmpMove = mv
					break
				}
				if (prom == WN || prom == BN) && moveString[4] == 'n' {
					tmpMove = mv
					break
				}
				continue // continue the loop on wrong promotions
			}
			// If no promotion needed or matches, return this move
			tmpMove = mv
			break
		}
	}

	if tmpMove == move.NoMove {
		return Board{}, false
	}

	newB.MakeMove(tmpMove, AllMoves)
	return newB, true
}

// PrintBoard prints the board state in a human-readable format with ranks and files.
func (b Board) PrintBoard() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if file == 0 {
				fmt.Printf("%d  ", 8-rank)
			}
			piece := -1

			// loop over all piece bitboards
			for bb := WP; bb <= BK; bb++ {
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
