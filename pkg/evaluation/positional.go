package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/evaluation/tables"
	"github.com/Tecu23/argov2/pkg/util"
)

const (
	Opening = 0
	Endgame = 1
)

type positionalEvaluator struct {
	pawnTables   [2][64]int
	knightTables [2][64]int
	bishopTables [2][64]int
	rookTables   [2][64]int
	queenTables  [2][64]int
	kingTables   [2][64]int
}

func newPositionalEvaluator() *positionalEvaluator {
	return &positionalEvaluator{
		pawnTables:   [2][64]int{tables.PawnOpeningTable, tables.PawnEndgameTable},
		knightTables: [2][64]int{tables.KingOpeningTable, tables.KnightEndgameTable},
		bishopTables: [2][64]int{tables.BishopOpeningTable, tables.BishopEndgameTable},
		rookTables:   [2][64]int{tables.RookOpeningTable, tables.RookEndgameTable},
		queenTables:  [2][64]int{tables.QueenOpeningTable, tables.QueenEndgameTable},
		kingTables:   [2][64]int{tables.KingOpeningTable, tables.KingEndgameTable},
	}
}

func (p *positionalEvaluator) Evaluate(b *board.Board) int {
	var mgScore, egScore int
	phase := calculateGamePhase(b)

	// Evaluate each piece on the board
	for sq := 0; sq < 64; sq++ {
		piece := b.GetPieceAt(sq)
		if piece == Empty {
			continue
		}

		pieceType := piece % 6     // Get piece type without color
		clr := util.PcColor(piece) // Get color (0 for White, 1 for Black)

		// Get the appropriate piece-square table scores
		var mgPosScore, egPosScore int

		// For Black pieces, mirror the square vertically
		evalSq := sq
		if clr == color.BLACK {
			evalSq = sq ^ 56 // Flip square vertically
		}

		// Get position scores based on piece type
		switch pieceType {
		case Pawn:
			mgPosScore = p.pawnTables[Opening][evalSq]
			egPosScore = p.pawnTables[Endgame][evalSq]
		case Knight:
			mgPosScore = p.knightTables[Opening][evalSq]
			egPosScore = p.knightTables[Endgame][evalSq]
		case Bishop:
			mgPosScore = p.bishopTables[Opening][evalSq]
			egPosScore = p.bishopTables[Endgame][evalSq]
		case Rook:
			mgPosScore = p.rookTables[Opening][evalSq]
			egPosScore = p.rookTables[Endgame][evalSq]
		case Queen:
			mgPosScore = p.queenTables[Opening][evalSq]
			egPosScore = p.queenTables[Endgame][evalSq]
		case King:
			mgPosScore = p.kingTables[Opening][evalSq]
			egPosScore = p.kingTables[Endgame][evalSq]
		}

		// Adjust score based on color
		if clr == color.BLACK {
			mgPosScore = -mgPosScore
			egPosScore = -egPosScore
		}

		// Add to phase-specific scores
		mgScore += mgPosScore
		egScore += egPosScore
	}

	// Interpolate between middlegame and endgame scores
	finalScore := interpolateScore(mgScore, egScore, phase)

	// Return score relative to side to move
	if b.Side == color.BLACK {
		return -finalScore
	}
	return finalScore
}

// calculateGamePhase determines the game phase based on remaining material
func calculateGamePhase(b *board.Board) int {
	// Phase weights for different pieces
	const (
		PawnPhase   = 0
		KnightPhase = 1
		BishopPhase = 1
		RookPhase   = 2
		QueenPhase  = 4
		TotalPhase  = 24 // 2*(1*4 + 1*4 + 2*2 + 4*1) for all pieces except pawns and kings
	)

	phase := TotalPhase

	// Count remaining material to determine phase
	for sq := 0; sq < 64; sq++ {
		piece := b.GetPieceAt(sq)
		if piece == Empty {
			continue
		}

		switch piece % 6 {
		case Knight:
			phase -= KnightPhase
		case Bishop:
			phase -= BishopPhase
		case Rook:
			phase -= RookPhase
		case Queen:
			phase -= QueenPhase
		}
	}

	// Convert to 256 scale for smooth interpolation
	return (phase*256 + (TotalPhase / 2)) / TotalPhase
}

// interpolateScore combines middlegame and endgame scores based on phase
func interpolateScore(mgScore, egScore, phase int) int {
	// phase 256 = middlegame, phase 0 = endgame
	return (mgScore*phase + egScore*(256-phase)) / 256
}
