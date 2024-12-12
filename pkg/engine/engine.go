// Package engine contains the engine class
package engine

import (
	"context"
	"fmt"
	"time"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/move"
)

type mainLine struct {
	moves []move.Move
	score int
	depth int
	nodes int64
}

type Engine struct {
	Options  Options
	mainLine mainLine
	start    time.Time
	progress func(SearchInfo)
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

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) SearchInfo {
	e.start = time.Now()
	e.mainLine = mainLine{}
	e.progress = searchParams.Progress

	fmt.Println("Called search")

	for {
		fmt.Println("Searching...")
		time.Sleep(2 * time.Second)
	}

	return SearchInfo{}
}

func (e *Engine) Clear() {}
