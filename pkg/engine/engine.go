// Package engine contains the engine class
package engine

import (
	"context"
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
}

func NewEngine(options Options) *Engine {
	return &Engine{
		Options: options,
	}
}

func (e *Engine) Prepare() {}

func (e *Engine) Search(ctx context.Context, searchParams SearchParams) SearchInfo {
	return SearchInfo{}
}

func (e *Engine) Clear() {}
