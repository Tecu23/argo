// Package transposition keeps all the logic for working with a transposition table
package transposition

import (
	"sync"

	"github.com/Tecu23/argov2/pkg/move"
)

// Entry flags to indicate the type of score stored
const (
	EXACT = iota // Exact score
	ALPHA        // Upper bound (failed low)
	BETA         // Lower bound (falied high)
)

// TTEntry represents a single entry in the transposition table
type TTEntry struct {
	Key   uint64    // Zobrist hash of position
	Depth int       // Depth of the search that produced this entry
	Flag  int       // Type of node (exact, alpha, beta)
	Score int       // Score of the position
	Move  move.Move // Best move for this position
	Age   uint8     // Used for replacement strategy
}

// Table stores previously evaluated positions
type Table struct {
	entries []TTEntry
	size    int
	age     uint8
	mutex   sync.RWMutex
}

// NewTable creates a new transposition table with the specified size (in MB)
func NewTable(sizeInMB int) *Table {
	// Calculate number of entries based on size
	entrySize := 16 // approximate size of TTEntry in bytes
	numEntries := (sizeInMB * 1024 * 1024) / entrySize

	// Ensure size is a power of 2 for efficient modulo operation
	size := 1
	for size < numEntries {
		size *= 2
	}

	return &Table{
		entries: make([]TTEntry, size),
		size:    size,
		age:     0,
	}
}

// Clear removes all entries from the table
func (t *Table) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.entries = make([]TTEntry, t.size)
}

// NewSearch increments the age counter for a new search iteration
func (t *Table) NewSearch() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.age = (t.age + 1) % 255 // Avoid overflow and keep age in uint8 range
}

// Store adds or updates an entry in the transposition table
func (t *Table) Store(key uint64, depth int, flag int, score int, move move.Move) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Calculate index using modulo. Since size is power of 2, we can bitwise AND
	index := int(key & uint64(t.size-1))

	// Replacement strategy: Always replace if current entry is from an older search
	// or if new entry has greater or equal depth
	entry := &t.entries[index]

	if entry.Key == 0 || // Empty slot
		entry.Age != t.age || // From an older search
		depth >= entry.Depth { // Equal or deeper search

		entry.Key = key
		entry.Depth = depth
		entry.Flag = flag
		entry.Score = score
		entry.Move = move
		entry.Age = t.age
	}
}

// Probe looks up a position in the transposition table
// Returns (found, score, move, flag)
func (t *Table) Probe(key uint64, depth int, alpha int, beta int) (bool, int, move.Move, int) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	index := int(key & uint64(t.size-1))
	entry := t.entries[index]

	// Check if we have a valid entry for this position
	if entry.Key == key {
		// Found the position, but check if depth is sufficient
		if entry.Depth >= depth {
			// Adjust mate scores based on current distance from root
			score := entry.Score

			// If it's an exact score, we can return it directly
			if entry.Flag == EXACT {
				return true, score, entry.Move, EXACT
			}

			// For bounds, only return if they cause a cutoff
			if entry.Flag == ALPHA && score <= alpha {
				return true, alpha, entry.Move, ALPHA
			}

			if entry.Flag == BETA && score >= beta {
				return true, beta, entry.Move, BETA
			}
		}

		// Return the move even if we can't use the score
		return false, 0, entry.Move, -1
	}

	return false, 0, move.NoMove, -1
}

// GetHashFull returns the percentage of the table that is currently used
func (t *Table) GetHashFull() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Sample 1000 entries to estimate fullness
	step := t.size / 1000
	if step == 0 {
		step = 1
	}

	count := 0
	for i := 0; i < t.size; i += step {
		if t.entries[i].Key != 0 && t.entries[i].Age == t.age {
			count++
		}
	}

	// Return percentage (0-1000)
	return (count * 1000) / (t.size / step)
}

// Size returns the number of entries in the table
func (t *Table) Size() int {
	return t.size
}
