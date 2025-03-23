package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

// ParseFEN sets the board state according to a given FEN string.
// It places pieces, sets side to move, castling rights, and en passant square.
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -
func ParseFEN(FEN string) (Board, error) {
	b := Board{}
	b.Reset()

	fenIdx := 0
	sq := 0

	// Parse the ranks from top (rank=7) to bottom (rank=0)
	for row := 0; row < 8; row++ {
		for sq = row * 8; sq < row*8+8; {
			char := string(FEN[fenIdx])
			fenIdx++

			if char == "/" {
				continue
			}

			// If char is a digit, skip that many squares
			if i, err := strconv.Atoi(char); err == nil {
				for j := 0; j < i; j++ {
					b.SetSq(Empty, sq)
					sq++
				}
				continue
			}

			// Otherwise, it should be a piece character
			if !strings.Contains(util.PcFen, char) {
				return Board{}, fmt.Errorf(
					"parse fen failed: %s",
					fmt.Sprintf("Invalid piece %s try next one", char),
				)
			}

			b.SetSq(util.Fen2pc(char), sq)
			sq++
		}
	}

	remaining := strings.Split(strings.TrimSpace(FEN[fenIdx:]), " ")

	// Set side to move
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			b.SideToMove = color.WHITE
		} else if remaining[0] == "b" {
			b.SideToMove = color.BLACK
		} else {
			b.SideToMove = color.WHITE
			return Board{}, fmt.Errorf(
				"parse fen failed: %s",
				fmt.Sprintf("%s invalid side to move color", remaining[0]),
			)
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
	b.HalfMoveClock = 0
	if len(remaining) > 3 {
		cnt, err := strconv.Atoi(remaining[3])
		if err != nil {
			b.HalfMoveClock = 0
		}

		b.HalfMoveClock = uint8(cnt)
	}

	b.calculateHash()

	return b, nil
}

// Mirror returns a new board that's flipped vertically (white pieces become black and vice versa)
func (b *Board) Mirror() *Board {
	// Create a new board
	mirrored := &Board{}

	mirrored.HalfMoveClock = b.HalfMoveClock
	mirrored.FullMoveCounter = b.FullMoveCounter

	mirroredHash := uint64(0)
	for pc := WP; pc <= BK; pc++ {
		oppPc := util.OppositeColorPiece(pc)
		bb := b.Bitboards[pc]

		var sq int
		var oppSq int

		for bb != 0 {
			sq = bb.FirstOne()
			oppSq = sq ^ 56

			mirrored.SetSq(oppPc, oppSq)
			mirroredHash ^= hash.HashTable.PieceSquare[oppPc*64+oppSq]
		}
	}

	// Mirror castling rights
	var newCastling Castlings

	if uint(b.Castlings)&ShortB != 0 {
		newCastling |= Castlings(ShortW)
		mirroredHash ^= hash.HashTable.Castling[0]
	}
	if uint(b.Castlings)&LongB != 0 {
		newCastling |= Castlings(LongW)
		mirroredHash ^= hash.HashTable.Castling[1]
	}
	if uint(b.Castlings)&ShortW != 0 {
		newCastling |= Castlings(ShortB)
		mirroredHash ^= hash.HashTable.Castling[2]
	}
	if uint(b.Castlings)&LongW != 0 {
		newCastling |= Castlings(LongB)
		mirroredHash ^= hash.HashTable.Castling[3]
	}

	mirrored.Castlings = newCastling

	// Mirror en passant square if exists
	if b.EnPassant != -1 {
		mirrored.EnPassant = b.EnPassant ^ 56 // Flip rank (0-7 becomes 7-0)
		file := mirrored.EnPassant % 8
		mirroredHash ^= hash.HashTable.EnPassant[file]
	}

	// Switch side to move
	mirrored.SideToMove = b.SideToMove.Opp()
	if mirrored.SideToMove == color.WHITE {
		mirroredHash ^= hash.HashTable.Side
	}

	// NOTE: This will not be the same as calculate hash but we do not
	// need it to be the same. We only need the mirror position for eval
	mirrored.hash = mirroredHash
	return mirrored
}

func (b *Board) GetPieceCountForSide(piece int, clr color.Color) int {
	switch piece {
	case Pawn:
		if clr == color.BLACK {
			return b.Bitboards[BP].Count()
		}
		return b.Bitboards[WP].Count()
	case Knight:
		if clr == color.BLACK {
			return b.Bitboards[BN].Count()
		}
		return b.Bitboards[WN].Count()
	case Bishop:
		if clr == color.BLACK {
			return b.Bitboards[BB].Count()
		}
		return b.Bitboards[WB].Count()
	case Rook:
		if clr == color.BLACK {
			return b.Bitboards[BR].Count()
		}
		return b.Bitboards[WR].Count()
	case Queen:
		if clr == color.BLACK {
			return b.Bitboards[BQ].Count()
		}
		return b.Bitboards[WQ].Count()
	case King:
		if clr == color.BLACK {
			return b.Bitboards[BK].Count()
		}
		return b.Bitboards[WK].Count()
	}

	return 0
}

func (b *Board) PieceCount(side color.Color) int {
	if side == color.WHITE {
		return b.Occupancies[color.WHITE].Count()
	}

	return b.Occupancies[color.BLACK].Count()
}

func (b *Board) GetPieceBB(color, piece int) bitboard.Bitboard {
	if color == 0 {
		switch piece {
		case Pawn:
			return b.Bitboards[WP]
		case Bishop:
			return b.Bitboards[WB]
		case Knight:
			return b.Bitboards[WN]
		case Rook:
			return b.Bitboards[WR]
		case Queen:
			return b.Bitboards[WQ]
		case King:
			return b.Bitboards[WK]
		default:
			return Empty
		}
	}

	switch piece {
	case Pawn:
		return b.Bitboards[BP]
	case Bishop:
		return b.Bitboards[BB]
	case Knight:
		return b.Bitboards[BN]
	case Rook:
		return b.Bitboards[BR]
	case Queen:
		return b.Bitboards[BQ]
	case King:
		return b.Bitboards[BK]
	default:
		return Empty
	}
}
