// Package board contains the board representation and all board functions
package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

type Board struct {
	Bitboards   [12]bitboard.Bitboard
	Occupancies [3]bitboard.Bitboard
	Side        color.Color
	EnPassant   int
	Rule50      uint8
	Castlings
}

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

func (b Board) CopyBoard() Board {
	boardCopy := b

	return boardCopy
}

// TakeBack should restore the board state to a copy
func (b *Board) TakeBack(copy Board) {
	*b = copy
}

// SetSq should set a square sq to a particular piece pc
func (b *Board) SetSq(piece, sq int) {
	pieceColor := util.PcColor(piece)

	if b.Occupancies[color.BOTH].Test(sq) {
		// need to remove the piece
		for p := constants.WP; p <= constants.BK; p++ {
			if b.Bitboards[p].Test(sq) {
				b.Bitboards[p].Clear(sq)
			}
		}

		b.Occupancies[color.BOTH].Clear(sq)
		b.Occupancies[color.WHITE].Clear(sq)
		b.Occupancies[color.BLACK].Clear(sq)
	}

	if piece == constants.Empty {
		return
	}
	// b.Key ^= PieceKeys[piece][sq]
	b.Bitboards[piece].Set(sq)

	if pieceColor == color.WHITE {
		b.Occupancies[color.WHITE].Set(sq)
	} else {
		b.Occupancies[color.BLACK].Set(sq)
	}

	b.Occupancies[color.BOTH] |= b.Occupancies[color.WHITE]
	b.Occupancies[color.BOTH] |= b.Occupancies[color.BLACK]
}

// PrintBoard should print the current position of the board
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

// ParseFEN should parse a FEN string and retrieve the board
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -
func (b *Board) ParseFEN(FEN string) {
	b.Reset()

	fenIdx := 0
	sq := 0

	// parsing the FEN from the start and setting the from top to bottom
	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {

			char := string(FEN[fenIdx])
			fenIdx++

			if char == "/" {
				continue
			}

			// if we find a number we should skip that many squares from our current board
			if i, err := strconv.Atoi(char); err == nil {
				for j := 0; j < i; j++ {
					b.SetSq(constants.Empty, sq)
					sq++
				}
				continue
			}

			// if we find an invalid piece we skip
			if !strings.Contains(util.PcFen, char) {
				fmt.Printf("Invalid piece %s try next one", char)
				// log.Errorf("error string invalid piece %s try next one", char)
				continue
			}

			b.SetSq(util.Fen2pc(char), sq)

			sq++
		}
	}

	remaining := strings.Split(strings.TrimSpace(FEN[fenIdx:]), " ")

	// Setting the Side to Move
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

	// Checking for castling
	b.Castlings = 0
	if len(remaining) > 1 {
		b.Castlings = ParseCastlings(remaining[1])
	}

	// En Passant
	b.EnPassant = -1
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			b.EnPassant = util.Fen2Sq[remaining[2]]
		}
	}

	// Cheking for 50 move rule
	b.Rule50 = 0
	if len(remaining) > 3 {
		cnt, err := strconv.Atoi(remaining[3])
		if err != nil {
			b.Rule50 = 0
		}

		b.Rule50 = uint8(cnt)
	}
}
