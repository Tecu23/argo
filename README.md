# ArGO Chess Engine

ArGO is a powerful chess engine written in Go, implementing modern chess
AI techniques including a Neural Network Efficiently Updated (NNUE) evaluation function.

## Features

- UCI protocol compliant for compatibility with chess GUIs
- Advanced move ordering and search techniques
  - Alpha-beta pruning with principal variation search
  - Transposition table
  - Move ordering heuristics (MVV-LVA, killer moves, history heuristics)
  - Late move reduction
- Neural network evaluation (NNUE)
  - Efficient incremental updates
  - Assembly-optimized for maximum performance on AMD64
- Magic bitboards for fast move generation
- Time management with dynamic adjustment based on position complexity

## System Requirements

- **Platform**: Linux AMD64 only
- The NNUE evaluation system uses custom assembly code optimized
  specifically for x86-64 instruction set on Linux.

## Building from Source

```bash
# Clone the repository
git clone https://github.com/Tecu23/argov2.git

# Navigate to the directory
cd argov2

# Build the project
make build-linux
```

## Usage

### Basic Usage

```bash
# Start ArGO in UCI mode
./argo
```

### Debug Mode

```bash
# Start with debug output, does not do anything at the moment
./argo -debug
```

### UCI Commands

ArGO implements the standard Universal Chess Interface (UCI) protocol.
Here are some common commands:

- `uci` - Identify the engine
- `isready` - Check if the engine is ready to receive commands
- `position [fen <fenstring> | startpos] [moves <move1> <move2> ...]` - Set
  up a position
- `go [depth <x> | movetime <x> | wtime <x> btime <x> winc <x> binc <x>]` -
  Start searching
- `stop` - Stop the current search
- `quit` - Exit the program

Example:

```sh
position startpos moves e2e4 e7e5
go depth 10
```

## Architecture

ArGO's source code is organized into several key packages:

- `attacks` - Contains pre-computed attack tables for all pieces
- `bitboard` - Implements bitboard operations for efficient board representation
- `board` - Chess board representation and move generation
- `engine` - Search algorithms and engine control
- `nnue` - Neural network position evaluation
- `move` - Move encoding and manipulation
- `uci` - UCI protocol implementation

## Chess Engine Evaluation Package

The NNUE (Efficiently Updated Neural Network) evaluation system is a
port of the [Koivisto](https://github.com/Luecx/Koivisto) chess engine's
evaluation function, translated from C++ to Go with assembly optimizations for AMD64.

### Implementation Details

The evaluation function uses a two-layer neural network with:

- Input features based on piece positions relative to king locations
- Hidden layer with 512 neurons
- Efficient incremental updates that avoid recomputing the entire network
- Assembly-optimized critical calculations for maximum performance

## License

This project is licensed under the GNU General Public License v3.0. As this is
a derivative work based on a GPL-licensed chess engine, the entire codebase must
remain under the GPL v3 license. This means:

1. You are free to use, modify, and distribute this software
2. If you distribute this software, you must provide the complete source code
3. Any modifications must also be released under the GPL v3
4. The full license text must be included with any distribution

## Acknowledgments

The NNUE evaluation function is based on the work of Kim Kahre and Finn Eggers,
creators of the Koivisto UCI chess engine. This port maintains the spirit and
structure of their implementation while adapting it for the Go language with
platform-specific optimizations.

The original source code can be found at: [Koivisto](https://github.com/Luecx/Koivisto)
