package board

import (
	"fmt"
	"strconv"
	"strings"

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

	b.calculateHash()

	return b, nil
}

// Mirror returns a new board that's flipped vertically (white pieces become black and vice versa)
func (b *Board) Mirror() *Board {
	// Create a new board
	mirrored := NewBoard()

	mirrored.Rule50 = b.Rule50
	mirrored.MoveNumber = b.MoveNumber

	for pc := WP; pc <= BK; pc++ {
		oppPc := util.OppositeColorPiece(pc)
		bb := b.Bitboards[pc]

		var sq int
		var oppSq int

		for bb.Count() > 0 {
			sq = bb.FirstOne()
			oppSq = sq ^ 56

			mirrored.SetSq(oppPc, oppSq)
		}
	}

	// Mirror castling rights
	var newCastling Castlings
	if uint(b.Castlings)&ShortW != 0 {
		newCastling |= Castlings(ShortB)
	}
	if uint(b.Castlings)&LongW != 0 {
		newCastling |= Castlings(LongB)
	}
	if uint(b.Castlings)&ShortB != 0 {
		newCastling |= Castlings(ShortW)
	}
	if uint(b.Castlings)&LongB != 0 {
		newCastling |= Castlings(LongW)
	}

	mirrored.Castlings = newCastling

	// Mirror en passant square if exists
	if b.EnPassant != -1 {
		mirrored.EnPassant = b.EnPassant ^ 56 // Flip rank (0-7 becomes 7-0)
	}

	// Switch side to move
	mirrored.Side = b.Side.Opp()

	// Update hash
	mirrored.hash = mirrored.calculateHash()

	return mirrored
}

func (b *Board) GetPieceCountForSide(piece int, clr color.Color) int {
	switch piece {
	case Pawn:
		if clr == color.BLACK {
			return b.Bitboards[BP].Count()
		} else {
			return b.Bitboards[WP].Count()
		}
	case Knight:
		if clr == color.BLACK {
			return b.Bitboards[BN].Count()
		} else {
			return b.Bitboards[WN].Count()
		}
	case Bishop:
		if clr == color.BLACK {
			return b.Bitboards[BB].Count()
		} else {
			return b.Bitboards[WB].Count()
		}
	case Rook:
		if clr == color.BLACK {
			return b.Bitboards[BR].Count()
		} else {
			return b.Bitboards[WR].Count()
		}
	case Queen:
		if clr == color.BLACK {
			return b.Bitboards[BQ].Count()
		} else {
			return b.Bitboards[WQ].Count()
		}
	case King:
		if clr == color.BLACK {
			return b.Bitboards[BK].Count()
		} else {
			return b.Bitboards[WK].Count()
		}
	}

	return 0
}

func (b *Board) OppositeBishops() bool {
	bishopCount := b.GetPieceCountForSide(Bishop, color.WHITE)
	oppBishopCount := b.GetPieceCountForSide(Bishop, color.BLACK)

	if bishopCount != 1 || oppBishopCount != 1 {
		return false
	}

	c := []int{0, 0}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()
		c[0] = sq % 2
	}

	bishopBB = b.Bitboards[BB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()
		c[1] = sq % 2
	}

	return c[0] != c[1]
}

// Candidate passed checks if pawn is passed or candidate passer. A pawn is passed
// if one of the three following conditions is true:
//
//	(a) there is no stoppers except some levers
//	(b) the only stoppers are the leverPush, but we outnumber them
//	(c) there is only one front stopper which can be levered.
//
// If there is a pawn of our color in the same file in front of a current pawn
// it's no longer counts as passed.
func (b *Board) CandidatePassed(side color.Color) int {
	count := 0
	pawns := b.Bitboards[WP]
	if side == color.BLACK {
		pawns = b.Bitboards[BP]
	}

	tmpSq := 0

	// Iterate through all squares
	for sq := 0; sq < 64; sq++ {
		tmpSq = sq
		if !pawns.Test(sq) {
			continue
		}

		if side == color.BLACK {
			tmpSq = tmpSq ^ 56
		}

		if b.IsPassedPawn(sq) {
			count++
		}
	}

	return count
}

// IsPassedPawn checks if a white pawn on a given square is passed
func (b *Board) IsPassedPawn(square int) bool {
	file := square % 8
	rank := square / 8

	// Must be a white pawn
	if !b.Bitboards[WP].Test(square) {
		return false
	}
	ty1 := 8
	ty2 := 8

	// Loop over the remaining ranks to check for other pawns
	for y := rank - 1; y >= 0; y-- {

		sq := y*8 + file
		if b.Bitboards[WP].Test(sq) {
			return false
		}

		if b.Bitboards[BP].Test(sq) {
			ty1 = y
		}

		if file > 0 {
			sq := y*8 + (file - 1)
			if b.Bitboards[BP].Test(sq) {
				ty2 = y
			}
		}
		// Check right file if it exists
		if file < 7 {
			sq := y*8 + (file + 1)
			if b.Bitboards[BP].Test(sq) {
				ty2 = y
			}
		}
	}

	if ty1 == 8 && ty2 >= rank-1 {
		return true
	}

	if ty2 < rank-2 || ty1 < rank-1 {
		return false
	}

	if ty2 >= rank && ty1 == rank-1 && rank < 4 {
		if b.Bitboards[WP].Test((rank+1)*8+file-1) &&
			!b.Bitboards[BP].Test(rank*8+file-1) &&
			!b.Bitboards[BP].Test((rank-1)*8+file-2) {
			return true
		}

		if b.Bitboards[WP].Test((rank+1)*8+file+1) &&
			!b.Bitboards[BP].Test(rank*8+file+1) &&
			!b.Bitboards[BP].Test((rank-1)*8+file+2) {
			return true
		}
	}

	if b.Bitboards[BP].Test((rank-1)*8 + file) {
		return false
	}

	lever := 0
	if b.Bitboards[BP].Test((rank-1)*8 + file - 1) {
		lever++
	}
	if b.Bitboards[BP].Test((rank-1)*8 + file + 1) {
		lever++
	}

	leverpush := 0
	if b.Bitboards[BP].Test((rank-2)*8 + file - 1) {
		leverpush++
	}
	if b.Bitboards[BP].Test((rank-2)*8 + file + 1) {
		leverpush++
	}

	phalanx := 0
	if b.Bitboards[WP].Test(rank*8 + file - 1) {
		phalanx++
	}
	if b.Bitboards[WP].Test(rank*8 + file + 1) {
		phalanx++
	}

	if lever-countSupportingPawns(b, rank*8+file) > 1 {
		return false
	}

	if leverpush-phalanx > 0 {
		return false
	}

	if lever > 0 && leverpush > 0 {
		return false
	}

	return true
}

// Helper function to count pawns supporting a square
func countSupportingPawns(b *Board, square int) int {
	file := square % 8
	rank := square / 8
	count := 0

	// Check bottom-left supporter
	if file > 0 && rank < 7 {
		if b.Bitboards[WP].Test(square + 7) {
			count++
		}
	}
	// Check bottom-right supporter
	if file < 7 && rank < 7 {
		if b.Bitboards[WP].Test(square + 9) {
			count++
		}
	}
	return count
}

func (b *Board) PieceCount(side color.Color) int {
	if side == color.WHITE {
		return b.Occupancies[color.WHITE].Count()
	}

	return b.Occupancies[color.BLACK].Count()
}
