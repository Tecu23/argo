package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
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
					b.SetSq(constants.Empty, sq)
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
			b.Side = color.WHITE
		} else if remaining[0] == "b" {
			b.Side = color.BLACK
		} else {
			b.Side = color.WHITE
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
	b.Rule50 = 0
	if len(remaining) > 3 {
		cnt, err := strconv.Atoi(remaining[3])
		if err != nil {
			b.Rule50 = 0
		}

		b.Rule50 = uint8(cnt)
	}

	return b, nil
}
