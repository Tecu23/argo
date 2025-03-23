// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package engine

import (
	"unsafe"

	"github.com/Tecu23/argov2/pkg/move"
)

type TTFlag uint8

const (
	TTExact TTFlag = iota // Exact score (within alpha-beta window)
	TTAlpha               // Upper bound (failed low, score <= alpha)
	TTBeta                // Lower bound (failed high, score >= beta)
)

// Each position gets an entry in the table
type TTEntry struct {
	Key      uint64    // Zobrist hash of the position
	Depth    int       // How Deep we searched
	Score    int       // Position evaluation
	Flag     TTFlag    // Type of score (exact/upper/lower bound)
	BestMove move.Move // Best move found
	Age      uint8     // when this entry was created
}

type TranspositionTable struct {
	entries []TTEntry
	size    int
	age     uint8
}

func NewTranspositionTable(sizeInMB int) *TranspositionTable {
	entrySize := unsafe.Sizeof(TTEntry{})
	entriesCount := (sizeInMB * 1024 * 1024) / int(entrySize)
	return &TranspositionTable{
		entries: make([]TTEntry, entriesCount),
		size:    entriesCount,
		age:     0,
	}
}

func (tt *TranspositionTable) NewSearch() {
	tt.age++ // Increment age for new search
}

func (tt *TranspositionTable) Clear() {
	for i := range tt.entries {
		tt.entries[i] = TTEntry{}
	}
}

func (tt *TranspositionTable) Store(key uint64, score, depth int, flag TTFlag, bestMove move.Move) {
	index := key % uint64(tt.size)
	entry := &tt.entries[index]

	// Replacement strategy
	if entry.Key == 0 || // Empty slot
		entry.Age < tt.age || // Older entry
		depth >= entry.Depth { // Deeper search
		entry.Key = key
		entry.Score = score
		entry.Depth = depth
		entry.Flag = flag
		entry.BestMove = bestMove
		entry.Age = tt.age
	}
}

func (tt *TranspositionTable) Probe(key uint64) (TTEntry, bool) {
	index := key % uint64(tt.size)
	entry := tt.entries[index]

	if entry.Key == key {
		return entry, true
	}
	return TTEntry{}, false
}
