// Package engine contains the engine class
package engine

import (
	"context"
	"time"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/move"
)

type mainLine struct {
	moves []move.Move
	score int
	depth int
	nodes int64
}

type Engine struct {
	Options     Options
	mainLine    mainLine
	start       time.Time
	progress    func(SearchInfo)
	timeManager timeManager
	cancel      context.CancelFunc
}

const (
	maxDepth   = 100
	infinity   = 50000
	mateScore  = 49000
	mateHeight = 48000
)

func NewEngine(options Options) *Engine {
	return &Engine{
		Options: options,
	}
}

func (e *Engine) Prepare() {}

// Search is the main entry point for starting a search
func (e *Engine) Search(ctx context.Context, params SearchParams) SearchInfo {
	e.start = time.Now()
	e.mainLine = mainLine{}
	e.progress = params.Progress

	// Get current position
	currentBoard := params.Boards[len(params.Boards)-1]

	// Setup time manager
	done := ctx.Done()
	tm := newTimeManager(params.Limits, currentBoard.Side == color.BLACK, done, e.cancel)

	// Start actual search
	return e.search(ctx, currentBoard, tm)
}

func (e *Engine) Clear() {}

// createSearchInfo creates a SearchInfo struct from current engine state
func (e *Engine) createSearchInfo() SearchInfo {
	return SearchInfo{
		Score: UciScore{
			Centipawns: e.mainLine.score,
			Mate:       0,
		},
		Depth:    e.mainLine.depth,
		Nodes:    e.mainLine.nodes,
		Time:     time.Since(e.start),
		MainLine: e.mainLine.moves,
	}
}
