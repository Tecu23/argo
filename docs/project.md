# Project Structure and Architecture

## Overview

ArGO is a modern chess engine written in Go, designed with a clean, modular
architecture that separates concerns and promotes maintainability. This document
outlines the high-level architecture, package organization, and design principles
used throughout the codebase.

## Directory Structure

```txt
argov2/
├── cmd/                    # Command-line applications
│   └── argo/               # Main chess engine executable
│       └── main.go         # Entry point
├── internal/               # Private implementation details
│   ├── hash/               # Zobrist hashing implementation
│   ├── history/            # Move history tables
│   ├── reduction/          # Late move reduction tables
│   └── types/              # Common type definitions
└── pkg/                    # Public packages
    ├── attacks/            # Attack generation and lookup tables
    ├── bitboard/           # Bitboard representation and operations
    ├── board/              # Chess board representation
    ├── color/              # Color definitions and utilities
    ├── constants/          # Game constants
    ├── engine/             # Search and evaluation
    ├── move/               # Move generation and representation
    ├── nnue/               # Neural network evaluation
    ├── uci/                # UCI protocol implementation
    └── util/               # Utility functions
```

## Key Components

### Core Chess Primitives

1. **Bitboard** (`pkg/bitboard`): Fundamental data structure for board representation
   using 64-bit integers

   - `Bitboard` type with bit manipulation methods
   - Bitwise operations for efficient piece movement and attacks

2. **Board** (`pkg/board`): Complete chess board representation

   - Maintains game state (pieces, castling rights, en passant, etc.)
   - Generates legal moves
   - Validates and applies moves

3. **Move** (`pkg/move`): Compact representation of chess moves
   - 32-bit encoding of source, target, piece, and flags
   - Methods for extracting move properties
   - Move validation and algebraic notation conversion

### Chess Logic

1. **Attacks** (`pkg/attacks`): Pre-computed attack patterns

   - Magic bitboards for sliding pieces (bishop, rook, queen)
   - Knight, king, and pawn attack tables
   - Square attack detection

2. **Constants** (`pkg/constants`): Game constants

   - Square indices (A1, B1, etc.)
   - File and rank definitions
   - Piece type constants
   - Castling rights flags

3. **Color** (`pkg/color`): Color-related functionality
   - WHITE/BLACK constants
   - Methods for switching sides

### Engine Core

1. **Engine** (`pkg/engine`): Search and evaluation

   - Alpha-beta search with Principal Variation Search
   - Quiescence search for tactical positions
   - Move ordering for search efficiency
   - Time management
   - Transposition table

2. **NNUE** (`pkg/nnue`): Neural network position evaluation

   - Efficiently updatable neural network
   - Feature transformation
   - Accumulator handling
   - ASM-optimized calculations

3. **UCI** (`pkg/uci`): Protocol implementation
   - UCI command parsing and handling
   - Position setup and search control
   - Engine configuration
   - Information formatting

## Architectural Patterns

### 1. Dependency Flow

ArGO follows a clean dependency flow:

```txt
              ┌─────────┐
              │  main   │
              └────┬────┘
                   │
              ┌────▼────┐
              │   uci   │
              └────┬────┘
                   │
              ┌────▼────┐
              │ engine  │
              └────┬────┘
                   │
         ┌─────────┴─────────┐
  ┌──────▼──────┐     ┌──────▼──────┐
  │    board    │     │    nnue     │
  └──────┬──────┘     └─────────────┘
         │
  ┌──────▼──────┐
  │   attacks   │
  └──────┬──────┘
         │
  ┌──────▼──────┐
  │  bitboard   │
  └─────────────┘
```

Higher-level components depend on lower-level ones, not vice versa,
creating a clean hierarchy.

### 2. Composition Over Inheritance

ArGO relies heavily on composition rather than inheritance:

```go
// Board composes multiple bitboards
type Board struct {
    Bitboards   [12]bitboard.Bitboard
    Occupancies [3]bitboard.Bitboard
    // ...
}

// Engine composes search components
type Engine struct {
    evaluator      nnue.Evaluator
    tt             *TranspositionTable
    reductionTable *reduction.Table
    historyTable   *history.HistoryTable
    // ...
}
```

### 3. Clean Interfaces

Interfaces are used to define clear contracts between components:

```go
// Engine interface defines contract for UCI interaction
type Engine interface {
    Prepare()
    Clear()
    Search(ctx context.Context, searchParams SearchParams) SearchInfo
}

// Option interface for configurable engine parameters
type Option interface {
    UciName() string
    UciString() string
    Set(s string) error
}
```

## Performance Considerations

### 1. Memory Management

ArGO balances memory usage and performance:

- **Statically Sized Arrays**: Pre-allocated arrays for attack tables and
  transposition tables
- **Memory Alignment**: Structure fields aligned for optimal cache usage
- **Minimizing Allocations**: Reusing move lists and other frequently allocated structures

### 2. Critical Path Optimization

Performance-critical code paths are carefully optimized:

- **Assembly Implementation**: Key NNUE operations implemented in assembly
  for maximum performance
- **Bitboard Operations**: Utilizing processor-specific bit manipulation instructions
- **Move Ordering**: Sophisticated move ordering to maximize alpha-beta pruning efficiency
- **Cache-Friendly Layouts**: Data structures designed to minimize cache misses

### 3. Concurrency Model

ArGO's concurrency approach:

- **Context-Based Cancellation**: Using Go's context package for clean search termination
- **Channel Communication**: Between search and UCI interface
- **Thread-Safety Considerations**: In shared data structures like
  the transposition table

## Module Dependencies

ArGO minimizes external dependencies:

- **Standard Library**: Heavy reliance on Go's standard library
- **Assembly**: For performance-critical NNUE operations
- **No External Libraries**: Self-contained implementation without third-party dependencies

## Build Process

The build process is kept simple:

```bash
# Standard build
go build -o argo

# Optimized build
go build -o argo -ldflags="-s -w" -gcflags="-N -l"

# Build with specific architecture optimizations
GOARCH=amd64 go build -o argo
```

## Testing Approach

The testing strategy includes:

1. **Perft Testing**: Verifying move generation correctness by counting
   positions at depth
2. **Search Validation**: Testing search on known positions with expected outcomes
3. **NNUE Validation**: Comparing evaluation output against reference values
4. **UCI Command Testing**: Ensuring UCI protocol compliance

## Future Architecture Extensions

The architecture is designed to accommodate future enhancements:

1. **Parallel Search**: Adding concurrent search capabilities
2. **Syzygy Tablebase**: Integration with endgame tablebases
3. **Transposition Table Persistence**: Saving and loading hash tables between sessions
4. **NNUE Retraining**: Infrastructure for training updated neural networks

## Design Principles

Throughout ArGO's development, several key principles were followed:

1. **Performance**: Critical code paths optimized for speed
2. **Clarity**: Clean code organization and naming conventions
3. **Modularity**: Clear separation of concerns between packages
4. **Testability**: Components designed to be independently testable
5. **Extensibility**: Architecture allows for future enhancements

These principles have resulted in a chess engine that balances performance
with maintainability, making ArGO both powerful and accessible to developers
interested in chess programming.
