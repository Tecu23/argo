# Go Chess Engine Project

This project is a Go-based chess engine that implements core
functionalities such as board representation, move generation,
and attack precomputation using bitboards and magic bitboards for
fast sliding piece move lookups. It supports various chess concepts,
including castling, en passant, and pawn promotions, and is designed with
a modular architecture that makes it easy to extend and improve.

## Features

- **Bitboard Representation**:  
   Every piece placement and occupancy is represented using 64-bit integers
  (bitboards), allowing efficient bitwise operations to check, set, and clear squares.
- **Magic Bitboards for Sliding Pieces**:  
   Rooks, bishops, and queens use magic numbers and precomputed tables for O(1)
  attack generation. This significantly speeds up move generation, particularly in
  complex positions.
- **Precomputed Attacks**:  
   Static arrays store precomputed attacks for knights, kings, and pawns.
  Sliding pieces (rooks, bishops, queens) leverage magic bitboard indexing.
  This allows constant-time retrieval of attack patterns.
- **Comprehensive Move Generation**:  
   The engine can generate all pseudo-legal moves for both sides,
  including special moves:
  - **Pawn Moves**: Single and double pushes, captures, en passant, and promotions.
  - **Castling**: Detects available castling rights and generates castling moves
    if king and rook haven’t moved and squares are not attacked.
  - **Promotions**: Automatically generates all promotion piece options
    (Q, R, B, N) when a pawn reaches the last rank.
- **Legality Checks**:  
   While the engine primarily generates pseudo-legal moves, it does integrate
  checks to ensure that the king is not placed in check by a move. Moves that
  would leave one’s own king in check are rejected.
- **FEN Parsing**:  
   The engine can initialize board states from Forsyth-Edwards Notation (FEN)
  strings, enabling easy testing and integration with other chess tools.
- **Performance Testing (Perft)**:  
   A perft (performance test) method is included, allowing validation of
  move generation correctness and performance tuning. Perft counts the total
  number of leaf nodes (positions) reached at a given search depth.

## Project Structure

```bash
.
├── cmd/
│   └── engine/        # Main application entry point (e.g., a UCI loop)
├── pkg/
│   ├── attacks/       # Precomputed attacks and magic bitboards logic
│   ├── bitboard/      # Bitboard type and basic bit manipulation helpers
│   ├── board/         # Board state representation, move generation, FEN parsing
│   ├── color/         # Color enumeration (white, black, both)
│   ├── constants/     # Chess constants (piece indices, directions, squares)
│   ├── move/          # Move encoding, move lists, and utility functions
│   ├── util/          # Miscellaneous utilities (timing, indexing maps)
│   └── ... (other optional packages)
└── go.mod             # Go module file

```

## Key Concepts

### Bitboards

A 64-bit integer represents the 8x8 chessboard, with one bit per square.
This structure makes checking occupancy, attacks, and legal moves very efficient
through bitwise operations.

### Magic Bitboards

Magic bitboards use carefully chosen constants (magic numbers) to map
board occupancy configurations to unique indices, enabling O(1) lookup of
sliding piece attacks. This allows fast generation of rook, bishop, and
queen moves by simply indexing into precomputed arrays.

### Move Generation

The engine iterates over the pieces of the side to move, calculates their
attack squares using the precomputed tables, and generates moves based on
current occupancy. Special rules like en passant and castling are integrated
into this process.

### Testing & Validation

- **Perft**: By running `PerftTest` at various depths and comparing results
  to known correct perft values, we validate the correctness of move generation.
- **Unit Tests**: Although not explicitly shown, you can add unit tests to verify
  bitboard operations, move legality checks, and specific edge cases in code.

## Building

You can build the project using `go build` commands or a provided `Makefile`.
For example, to build a binary for your platform:

```sh
make
```

To specify a version or target OS/ARCH:

```sh
make VERSION=1.0.0
make GOOS=windows GOARCH=amd64 build
```

## Usage

- **FEN Parsing**:  
   Use `ParseFEN` to set the board to a given position:

```go
  b := board.Board{}
  b.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
```

- **Move Generation**:  
   Create a `Movelist` and call `GenerateMoves`:

  ```go
  var moves move.Movelist
  b.GenerateMoves(&moves)
  for _, mv := range moves {
     fmt.Println(mv.String())
  }
  ```

- **Checking Attacks**:

  ```go
  attacked := b.IsSquareAttacked(constants.E4, color.BLACK)
  fmt.Println("Is E4 attacked by Black?", attacked)
  ```

## Roadmap

- **UCI Interface**:  
   Integrate a UCI (Universal Chess Interface) loop to allow external GUIs
  to communicate with the engine.
- **Search & Evaluation**:  
   Implement alpha-beta search with iterative deepening, transposition tables,
  and evaluation heuristics to make the engine play competitive chess.
- **Optimizations**:  
   Add parallel search, advanced evaluation terms, and opening books or endgame tablebases.

## Contributing

Contributions are welcome! If you find a bug or have a suggestion,
please open an issue or submit a pull request. Before contributing,
ensure your code is tested and formatted with standard Go tools (`go fmt`).

## License

This project is licensed under the MIT License. See the LICENSE file for details.
