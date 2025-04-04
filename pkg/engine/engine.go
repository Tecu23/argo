// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package engine contains the engine class
package engine

import (
	"context"
	"time"

	"github.com/Tecu23/argov2/internal/history"
	"github.com/Tecu23/argov2/internal/reduction"
	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/nnue"
)

type mainLine struct {
	moves []move.Move
	score int
	depth int
	nodes int64
}

type Engine struct {
	nodes          int64
	Options        Options
	mainLine       mainLine
	start          time.Time
	progress       func(SearchInfo)
	timeManager    *timeManager
	cancel         context.CancelFunc
	evaluator      nnue.Evaluator
	tt             *TranspositionTable
	reductionTable *reduction.Table
	historyTable   *history.HistoryTable
	killerMoves    [MaxDepth][MaxKillers]move.Move
}

func NewEngine(options Options) *Engine {
	return &Engine{
		Options:        options,
		evaluator:      *nnue.NewEvaluator(),
		tt:             NewTranspositionTable(32),
		reductionTable: reduction.New(),
		historyTable:   history.New(),
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

	e.timeManager = newTimeManager(ctx, e.start, params.Limits, &currentBoard)
	defer e.timeManager.Close()

	// Start actual search
	return e.search(ctx, &currentBoard, e.timeManager)
}

func (e *Engine) Clear() {
	e.tt.Clear()
}

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
