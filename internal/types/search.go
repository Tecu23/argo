// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package types
package types

import (
	"time"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
)

type UciScore struct {
	Centipawns int
	Mate       int
}

type SearchInfo struct {
	Score    UciScore
	Depth    int
	Nodes    int64
	Time     time.Duration
	MainLine []move.Move
}

type LimitsType struct {
	Ponder         bool
	Infinite       bool
	WhiteTime      int
	BlackTime      int
	WhiteIncrement int
	BlackIncrement int
	MoveTime       int
	MovesToGo      int
	Depth          int
	Nodes          int
	Mate           int
}

type SearchParams struct {
	Boards   []board.Board
	Limits   LimitsType
	Progress func(si SearchInfo)
}
