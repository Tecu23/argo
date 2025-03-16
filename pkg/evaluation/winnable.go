// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func (e *Evaluator) WinnableEvaluation(b *board.Board, mirror *board.Board, mg, eg int) (int, int) {
	return e.winnableMg(b, mirror, mg), e.winnableEg(b, mirror, eg)
}

func (e *Evaluator) winnableMg(b *board.Board, mirror *board.Board, mg int) int {
	if mg == -1 {
		tmg, _ := e.EvaluateOneSide(b, mirror, true)
		return tmg
	}

	factor := 0
	if mg > 0 {
		factor = 1
	} else if mg < 0 {
		factor = -1
	}

	return factor * max(min(winnable(b)+50, 0), -abs(mg))
}

func (e *Evaluator) winnableEg(b *board.Board, mirror *board.Board, eg int) int {
	if eg == -1 {
		_, teg := e.EvaluateOneSide(b, mirror, true)
		return teg
	}

	factor := 0
	if eg > 0 {
		factor = 1
	} else if eg < 0 {
		factor = -1
	}

	return factor * max(winnable(b), -abs(eg))
}

// winnable computes the winnable correction value for this position, i.e. second
// order bonus/malus based on the known attacking/defending status of the players
func winnable(b *board.Board) int {
	pawns := 0
	kx := []int{0, 0}
	ky := []int{0, 0}
	flanks := []int{0, 0}

	for x := 0; x < 8; x++ {
		open := []int{0, 0}
		for y := 0; y < 8; y++ {
			if b.Bitboards[WP].Test(y*8 + x) {
				open[0] = 1
				pawns++
			} else if b.Bitboards[BP].Test(y*8 + x) {
				open[1] = 1
				pawns++
			}

			if b.Bitboards[WK].Test(y*8 + x) {
				kx[0] = x
				ky[0] = y

			} else if b.Bitboards[BK].Test(y*8 + x) {
				kx[1] = x
				ky[1] = y
			}
		}

		if open[0]+open[1] > 0 {
			if x < 4 {
				flanks[0] = 1
			} else {
				flanks[1] = 1
			}
		}
	}

	mirror := b.Mirror()

	passedCount := b.CandidatePassed() + mirror.CandidatePassed()
	bothFlanks := flanks[0] != 0 && flanks[1] != 0

	outflanking := abs(kx[0]-kx[1]) - abs(ky[0]-ky[1])
	purePawn := nonPawnMaterial(b, color.WHITE)+nonPawnMaterial(mirror, color.WHITE) == 0
	almostWinnable := outflanking < 0 && !bothFlanks
	infiltration := ky[0] < 4 || ky[1] > 3

	score := -110

	score += 9 * passedCount
	score += 12 * pawns
	score += 9 * outflanking
	if infiltration {
		score += 24
	}

	if bothFlanks {
		score += 21
	}

	if purePawn {
		score += 51
	}

	if almostWinnable {
		score -= 43
	}

	return score
}
