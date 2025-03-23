// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package engine

// import (
// 	. "github.com/Tecu23/argov2/internal/types"
// )

// func TestTranspositionTable(t *testing.T) {
// 	options := NewOptions()
// 	e := NewEngine(options)
//
// 	positions := []struct {
// 		fen   string
// 		depth int
// 	}{
// 		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5},
// 		{"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 0 1", 4},
// 	}
//
// 	for _, pos := range positions {
// 		b, _ := board.ParseFEN(pos.fen)
//
// 		// First search - should fill TT
// 		firstScore := e.Search(context.Background(), SearchParams{
// 			Limits: LimitsType{
// 				Depth: pos.depth,
// 			},
// 			Boards: []board.Board{b},
// 		})
//
// 		// Second search - should be much faster due to TT hits
// 		start := time.Now()
// 		secondScore := e.Search(context.Background(), SearchParams{
// 			Limits: LimitsType{
// 				Depth: pos.depth,
// 			},
// 			Boards: []board.Board{b},
// 		})
// 		secondTime := time.Since(start)
//
// 		if firstScore.Score != secondScore.Score {
// 			t.Errorf("Scores differ: first=%v, second=%v", firstScore, secondScore)
// 		}
//
// 		// Add logging to see improvement
// 		t.Logf("Position: %s, Depth: %d, Time: %v", pos.fen, pos.depth, secondTime)
// 	}
// }
