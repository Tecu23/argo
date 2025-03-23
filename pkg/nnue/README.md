# Chess Engine Evaluation Package

## Overview

This package implements a chess position evaluation function ported from
[Koivisto UCI](https://github.com/Luecx/Koivisto) chess engine. The original C++
implementation has been carefully translated to Golang while preserving
the evaluation logic and performance characteristics.

## Features

- Complete port of the [Koivisto](https://github.com/Luecx/Koivisto)
  NNUE evaluation function
- Assembly optimizations for critical calculation paths on AMD64 architecture
- Designed for high-performance chess analysis

## Platform Requirements

**Important**: This engine is specifically designed for **Linux AMD64** systems only.

Due to the low-level assembly optimizations implemented for performance,
this package will not function on other operating systems or CPU architectures.
The assembly code targets x86-64 instructions available on AMD64
processors running Linux.

## Implementation Details

The core evaluation logic follows the original C++ implementation while
taking advantage of Golang's features. Performance-critical sections
have been optimized with hand-tuned assembly code to ensure maximum calculation
speed for position evaluation.

## Usage

```golang
    import "github.com/tecu23/argov2/pkg/nnue"

    // Initialize weights of the NNUE
    if err := nnue.InitializeNNUE(); err != nil {
        log.Fatalf("Error initializing NNUE: %v", err)
    }

    // Create a new evaluator
    evaluator := nnue.NewEvaluator()

    // First reset the evaluator with the starting board position
    // so the weights are correctly handled
    evaluator.Reset(board)

    // Evaluate the current board position
    evaluator.Evaluate(board)
```

For efficient updating during search

```golang

    // Process a single move on the current board position
    evaluator.ProcessMove(board, move)

    // Evaluate the updated board position
    evaluator.Evaluate(board)

    // Reset the evaluator by popping the last move made
    evaluator.PopAccumulation()
```

## License

This project is licensed under the GNU General Public License v3.0.

As this is a derivative work based on a GPL-licensed chess engine,
the entire codebase must remain under the GPL v3 license. This means:

1. You are free to use, modify, and distribute this software
2. If you distribute this software, you must provide the complete source code
3. Any modifications must also be released under the GPL v3
4. The full license text must be included with any distribution

See the [LICENSE](LICENSE) file for the complete text of the GNU GPL v3 license.

## Acknowledgments

All credit for the original evaluation logic goes to Kim Kahre and Finn Eggers,
creators of the Koivisto UCI chess engine. This port maintains the spirit and
structure of their implementation while adapting it for the Go language with
platform-specific optimizations.

The original source code can be found
at:[Koivisto](https://github.com/Luecx/Koivisto)
