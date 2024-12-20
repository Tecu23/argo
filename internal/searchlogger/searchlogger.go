package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
)

// Debug flags for different aspects of the search
type SearchDebug struct {
	enabled         bool
	logFile         *os.File
	logSearchInfo   bool // Log general search progress
	logMoveOrdering bool // Log move ordering decisions
	logPruning      bool // Log pruning decisions
	logTT           bool // Log transposition table hits/misses
	logQuiescence   bool // Log quiescence search
	depth           int  // Current search depth
	ply             int  // Current ply
}

// SearchLogger handles debug output for the search
type SearchLogger struct {
	debug SearchDebug
	buf   strings.Builder
}

func NewSearchLogger(filename string) (*SearchLogger, error) {
	logger := &SearchLogger{}
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		logger.debug.logFile = f
		logger.debug.enabled = true
		logger.debug.logMoveOrdering = true
		logger.debug.logTT = true
		logger.debug.logPruning = true
		logger.debug.enabled = true
	}
	return logger, nil
}

func (l *SearchLogger) LogSearchNode(
	b *board.Board,
	depth, ply, alpha, beta, score int,
	mv move.Move,
) {
	if !l.debug.enabled {
		return
	}

	l.buf.Reset()
	fmt.Fprintf(&l.buf, "d:%2d p:%2d [%6d,%6d] ", depth, ply, alpha, beta)

	// Add indentation based on ply
	for i := 0; i < ply; i++ {
		l.buf.WriteString("  ")
	}

	// Log the move if present
	if mv != move.NoMove {
		fmt.Fprintf(&l.buf, "%s ", mv.String())
	}

	// Log score and position details
	fmt.Fprintf(&l.buf, "score:%6d hash:%016x\n", score, b.Hash())

	if l.debug.logFile != nil {
		l.debug.logFile.WriteString(l.buf.String())
	}
}

func (l *SearchLogger) LogMoveOrdering(ply int, moves []move.Move, scores []int) {
	if !l.debug.logMoveOrdering || !l.debug.enabled {
		return
	}

	l.buf.Reset()
	fmt.Fprintf(&l.buf, "Move ordering at ply %d:\n", ply)
	for i, m := range moves {
		fmt.Fprintf(&l.buf, "  %s: %d\n", m.String(), scores[i])
	}
	l.buf.WriteString("\n")

	if l.debug.logFile != nil {
		l.debug.logFile.WriteString(l.buf.String())
	}
}

func (l *SearchLogger) LogPruning(reason string, depth, beta, score int) {
	if !l.debug.logPruning || !l.debug.enabled {
		return
	}

	l.buf.Reset()
	fmt.Fprintf(&l.buf, "Pruning at depth %d: %s (β=%d, score=%d)\n",
		depth, reason, beta, score)

	if l.debug.logFile != nil {
		l.debug.logFile.WriteString(l.buf.String())
	}
}
