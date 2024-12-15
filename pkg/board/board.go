// Package board contains the board representation and all board helper functions.
// This package will handle move geration
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/internal/hash"
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
	MoveNumber  int
	hash        uint64
}

// Reset restores the board to an initial empty state and sets defaults.
func (b *Board) Reset() {
	b.Side = color.WHITE
	b.EnPassant = -1
	b.Castlings = 0
	b.Rule50 = 0
	b.MoveNumber = 0

	for i := 0; i < 12; i++ {
		b.Bitboards[i] = 0
	}

	for i := 0; i < 3; i++ {
		b.Occupancies[i] = 0
	}

	b.calculateHash()
}

func NewBoard() *Board {
	b := &Board{}
	b.Reset()
	b.hash = b.calculateHash()
	return b
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

// SetSq should set a square sq to a particular piece pc
func (b *Board) SetSq(piece, sq int) {
	pieceColor := util.PcColor(piece)

	if b.Occupancies[color.BOTH].Test(sq) {
		// need to remove the piece
		for p := WP; p <= BK; p++ {
			if b.Bitboards[p].Test(sq) {
				b.Bitboards[p].Clear(sq)
				b.hash ^= hash.HashTable.PieceSquare[p*64+sq]
			}
		}

		b.Occupancies[color.BOTH].Clear(sq)
		b.Occupancies[color.WHITE].Clear(sq)
		b.Occupancies[color.BLACK].Clear(sq)
	}

	if piece == Empty {
		return
	}
	b.hash ^= hash.HashTable.PieceSquare[piece*64+sq]
	b.Bitboards[piece].Set(sq)

	if pieceColor == color.WHITE {
		b.Occupancies[color.WHITE].Set(sq)
	} else {
		b.Occupancies[color.BLACK].Set(sq)
	}

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

// TODO: REDO THIS for better performance

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

		// If there was an en passant square, remove it from hash
		if b.EnPassant != -1 {
			file := b.EnPassant % 8
			b.hash ^= hash.HashTable.EnPassant[file]
		}
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

			b.hash ^= hash.HashTable.Side
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

			if b.Side == color.BLACK {
				b.MoveNumber++
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
				b.hash ^= hash.HashTable.EnPassant[(src-8)%8]
			} else {
				b.EnPassant = src + 8
				b.hash ^= hash.HashTable.EnPassant[(src+8)%8]
			}
		}

		b.SetSq(Empty, src)

		if prom != 0 {
			b.SetSq(prom, tgt)
		} else {
			b.SetSq(pc, tgt)
		}

		// Update castling rights if necessary
		oldCastling := b.Castlings
		b.Castlings &= Castlings(CastlingRights[src])
		b.Castlings &= Castlings(CastlingRights[tgt])

		// Update hash for changed castling rights
		if oldCastling != b.Castlings {
			if uint(oldCastling)&ShortW != uint(b.Castlings)&ShortW {
				b.hash ^= hash.HashTable.Castling[0]
			}
			if uint(oldCastling)&LongW != uint(b.Castlings)&LongW {
				b.hash ^= hash.HashTable.Castling[1]
			}
			if uint(oldCastling)&ShortB != uint(b.Castlings)&ShortB {
				b.hash ^= hash.HashTable.Castling[2]
			}
			if uint(oldCastling)&LongB != uint(b.Castlings)&LongB {
				b.hash ^= hash.HashTable.Castling[3]
			}
		}

		// change side
		b.hash ^= hash.HashTable.Side
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

	if b.Side == color.BLACK {
		b.MoveNumber++
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

// InCheck determines if the current side to move is in check
func (b *Board) InCheck() bool {
	var kingPos int
	var kingBB bitboard.Bitboard

	// Find king position for the side to move
	if b.Side == color.WHITE {
		if b.Bitboards[WK] == 0 {
			return false
		}

		kingBB = b.Bitboards[WK]
		kingPos = kingBB.FirstOne()
		return b.IsSquareAttacked(kingPos, color.BLACK)
	}

	if b.Bitboards[BK] == 0 {
		return false
	}
	kingBB = b.Bitboards[BK]
	kingPos = kingBB.FirstOne()
	return b.IsSquareAttacked(kingPos, color.WHITE)
}

// IsCheckmate determines if the current position is checkmate
func (b *Board) IsCheckmate() bool {
	// If not in check, can't be checkmate
	if !b.InCheck() {
		return false
	}

	// Generate all possible moves
	moves := b.GenerateMoves()

	// Try each move to see if it gets us out of check
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if b.MakeMove(mv, AllMoves) {
			b.TakeBack(copyB)
			return false
		}
		b.TakeBack(copyB)
	}

	// No legal moves found while in check => checkmate
	return true
}

// IsStalemate determines if the current position is stalemate
func (b *Board) IsStalemate() bool {
	// If in check, can't be stalemate
	if b.InCheck() {
		return false
	}

	// Generate all possible moves
	moves := b.GenerateMoves()

	// Try each move to see if it gets us out of check
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if b.MakeMove(mv, AllMoves) {
			b.TakeBack(copyB)
			return false // Found a legal move, not stalemate
		}
		b.TakeBack(copyB)
	}

	// No legal moves found while not in check => stalemate
	return true
}

// IsInsufficientMaterial checks if there are enough pieces left for checkmate
func (b *Board) IsInsufficientMaterial() bool {
	// Get piece counts
	whitePieceCount := (b.Bitboards[WP] | b.Bitboards[WR] | b.Bitboards[WQ]).Count()
	blackPieceCount := (b.Bitboards[BP] | b.Bitboards[BR] | b.Bitboards[BQ]).Count()

	// If any pawns, rooks, or queens exist, there's sufficient material
	if whitePieceCount > 0 || blackPieceCount > 0 {
		return false
	}

	whiteKnights := b.Bitboards[WN].Count()
	blackKnights := b.Bitboards[BN].Count()
	whiteBishops := b.Bitboards[WB].Count()
	blackBishops := b.Bitboards[BB].Count()

	// King vs King
	if whiteKnights == 0 && blackKnights == 0 && whiteBishops == 0 && blackBishops == 0 {
		return true
	}

	// King + minor piece vs King
	if whiteKnights+whiteBishops <= 1 && blackKnights+blackBishops == 0 {
		return true
	}
	if blackKnights+blackBishops <= 1 && whiteKnights+whiteBishops == 0 {
		return true
	}

	// King + 2 knights vs King
	if whiteKnights == 2 && whiteBishops == 0 && blackKnights == 0 && blackBishops == 0 {
		return true
	}
	if blackKnights == 2 && blackBishops == 0 && whiteKnights == 0 && whiteBishops == 0 {
		return true
	}

	return false
}

func (b *Board) calculateHash() uint64 {
	tmpHash := uint64(0)

	// Hash Pieces
	for piece := WP; piece <= BK; piece++ {
		bb := b.Bitboards[piece]
		for bb != 0 {
			square := bb.FirstOne()
			tmpHash ^= hash.HashTable.PieceSquare[piece*64+square]
		}
	}

	// Hash Castling rights
	if uint(b.Castlings)&ShortW != 0 {
		tmpHash ^= hash.HashTable.Castling[0]
	}
	if uint(b.Castlings)&LongW != 0 {
		tmpHash ^= hash.HashTable.Castling[1]
	}
	if uint(b.Castlings)&ShortB != 0 {
		tmpHash ^= hash.HashTable.Castling[2]
	}
	if uint(b.Castlings)&LongB != 0 {
		tmpHash ^= hash.HashTable.Castling[3]
	}

	// Hash en passant
	if b.EnPassant != -1 {
		file := b.EnPassant % 8
		tmpHash ^= hash.HashTable.EnPassant[file]
	}

	// Hash side to move
	if b.Side == color.WHITE {
		tmpHash ^= hash.HashTable.Side
	}

	return tmpHash
}

// Update hash incrementally when making moves
func (b *Board) updateHashForMove(m move.Move) {
	from := m.GetSource()
	to := m.GetTarget()
	piece := m.GetPiece()
	capture := m.GetCapture()
	promotion := m.GetPromoted()

	// Remove piece from source square
	b.hash ^= hash.HashTable.PieceSquare[piece*64+from]

	// Add piece to destination square
	if promotion != 0 {
		b.hash ^= hash.HashTable.PieceSquare[promotion*64+to]
	} else {
		b.hash ^= hash.HashTable.PieceSquare[piece*64+to]
	}

	// Handle Capture
	if capture != 0 {
		b.hash ^= hash.HashTable.PieceSquare[piece*64+to]
	}

	// Update en passant
	if b.EnPassant != -1 {
		file := b.EnPassant % 8
		b.hash ^= hash.HashTable.EnPassant[file]
	}

	// Handle new en passant
	if m.GetDoublePush() != 0 {
		file := to % 8
		b.hash ^= hash.HashTable.EnPassant[file]
	}

	// Update castling rights
	oldRights := b.Castlings
	newRights := b.Castlings & Castlings(CastlingRights[from]) & Castlings(CastlingRights[to])
	if oldRights != newRights {
		if uint(oldRights)&ShortW != uint(newRights)&ShortW {
			b.hash ^= hash.HashTable.Castling[0]
		}
		if uint(oldRights)&LongW != uint(newRights)&LongW {
			b.hash ^= hash.HashTable.Castling[1]
		}
		if uint(oldRights)&ShortB != uint(newRights)&LongB {
			b.hash ^= hash.HashTable.Castling[2]
		}
		if uint(oldRights)&LongB != uint(newRights)&LongB {
			b.hash ^= hash.HashTable.Castling[3]
		}
	}

	// Switch side to move
	b.hash ^= hash.HashTable.Side
}

// Hash method to get current hash
func (b *Board) Hash() uint64 {
	return b.hash
}
