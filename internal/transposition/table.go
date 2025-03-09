package transposition

import (
	"sync"

	"github.com/Tecu23/argov2/pkg/move"
)

// TTFlag represents the type of score stored in a transposition table entry
type TTFlag int

const (
	TTExact TTFlag = iota // Exact score
	TTAlpha               // Upper bound (fail low)
	TTBeta                // Lower bound (fail high)
)

// TTEntry represents an entry in the transposition table
type TTEntry struct {
	Hash     uint64      // Position hash
	Depth    int         // Search depth
	Score    int         // Evaluation score
	Flag     TTFlag      // Type of score (exact, alpha, beta)
	BestMove move.Move   // Best move for this position
	Age      uint8       // Used for replacement strategy
	PV       []move.Move // Principal variation for this position
}

// Table implements a hash table for storing previously searched positions
type Table struct {
	entries    []TTEntry  // The table entries
	size       int        // Number of entries in the table
	currentAge uint8      // Current age for replacement strategy
	mutex      sync.Mutex // Mutex for thread safety
}

// New creates a new transposition table with the specified size in MB
func New(sizeInMB int) *Table {
	// Calculate number of entries that fit in the specified size
	// Each entry size depends on the struct size (can be approximated)
	entrySize := 32 // Approximate size in bytes (adjust based on actual size)
	numEntries := (sizeInMB * 1024 * 1024) / entrySize

	return &Table{
		entries:    make([]TTEntry, numEntries),
		size:       numEntries,
		currentAge: 0,
	}
}

// getIndex calculates the index in the table for a given hash
func (tt *Table) getIndex(hash uint64) int {
	return int(hash % uint64(tt.size))
}

// Probe looks up a position in the transposition table
func (tt *Table) Probe(hash uint64) (TTEntry, bool) {
	index := tt.getIndex(hash)

	tt.mutex.Lock()
	entry := tt.entries[index]
	tt.mutex.Unlock()

	// Check if this is the right position
	if entry.Hash == hash {
		return entry, true
	}

	return TTEntry{}, false
}

// Store saves a position to the transposition table
func (tt *Table) Store(
	hash uint64,
	score, depth int,
	flag TTFlag,
	bestMove move.Move,
) {
	tt.StoreEntry(TTEntry{
		Hash:     hash,
		Depth:    depth,
		Score:    score,
		Flag:     flag,
		BestMove: bestMove,
		Age:      tt.currentAge,
		PV:       nil, // No PV stored in this simplified version
	})
}

// StoreEntry saves a complete entry to the transposition table
func (tt *Table) StoreEntry(entry TTEntry) {
	index := tt.getIndex(entry.Hash)

	tt.mutex.Lock()
	defer tt.mutex.Unlock()

	existing := tt.entries[index]

	// Replacement strategy: always replace if hash matches
	// Otherwise, replace if new entry is deeper or older
	if existing.Hash == 0 || // Empty slot
		existing.Hash == entry.Hash || // Same position
		entry.Depth >= existing.Depth || // Deeper search
		existing.Age != tt.currentAge { // Old entry from previous search
		tt.entries[index] = entry
	}
}

// NewSearch increments the age counter for the replacement strategy
func (tt *Table) NewSearch() {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()

	tt.currentAge = (tt.currentAge + 1) % 255
}

// Clear completely empties the transposition table
func (tt *Table) Clear() {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()

	// Create a new array to ensure garbage collection of old entries
	tt.entries = make([]TTEntry, tt.size)
	tt.currentAge = 0
}

// GetHashFull returns the percentage of the table that is in use
func (tt *Table) GetHashFull() int {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()

	used := 0
	sampleSize := 1000 // Check a sample rather than the entire table

	if tt.size < sampleSize {
		sampleSize = tt.size
	}

	for i := 0; i < sampleSize; i++ {
		if tt.entries[i].Hash != 0 {
			used++
		}
	}

	return (used * 1000) / sampleSize
}
