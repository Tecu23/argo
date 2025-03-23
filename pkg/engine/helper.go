// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package engine

import (
	"sort"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/nnue"
)

func winIn(height int) int {
	return MateScore - height
}

func lossIn(height int) int {
	return -MateScore + height
}

type MoveScore struct {
	move  move.Move
	score int
}

func (e *Engine) updateKillers(mv move.Move, ply int) {
	if ply >= MaxDepth {
		return
	}
	// Don't store captures as killer moves
	if mv.IsCapture() {
		return
	}

	// Don't store a move that's already a killer at this ply
	for i := 0; i < MaxKillers; i++ {
		if e.killerMoves[ply][i] == mv {
			return
		}
	}

	// Shift existing killers and insert new one at first position
	for i := MaxKillers - 1; i > 0; i-- {
		e.killerMoves[ply][i] = e.killerMoves[ply][i-1]
	}
	e.killerMoves[ply][0] = mv
}

func (e *Engine) orderMoves(
	moves []move.Move,
	b *board.Board,
	ttMove move.Move,
	ply int,
) []move.Move {
	scores := make([]MoveScore, len(moves))
	stm := b.SideToMove

	for i, mv := range moves {
		score := 0

		// TT move gets highest priority
		if mv == ttMove {
			score = 2_000_000
		} else if mv.IsCapture() {
			// MVV-LVA scoring
			victim := b.GetPieceAt(mv.GetTargetSquare())
			aggressor := b.GetPieceAt(mv.GetSourceSquare())
			score = 1_000_000 + (nnue.GetPieceValue(victim) - nnue.GetPieceValue(aggressor)/10)
		} else {
			for j := 0; j < MaxKillers; j++ {
				if mv == e.killerMoves[ply][j] {
					score = 900_000 - j*1000
					break
				}
			}

			if score == 0 {
				score = e.historyTable.Get(stm, mv.GetSourceSquare(), mv.GetTargetSquare())
			}
		}

		scores[i] = MoveScore{mv, score}
	}

	// Sort moves by score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Extract sorted moves
	sortedMoves := make([]move.Move, len(moves))
	for i, ms := range scores {
		sortedMoves[i] = ms.move
	}

	return sortedMoves
}

func adjustScore(score, ply int) int {
	if score >= MateScore-MaxDepth {
		return score - ply
	}
	if score <= -MateScore+MaxDepth {
		return score + ply
	}
	return score
}
