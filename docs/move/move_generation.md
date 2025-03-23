# Move Generation System

## Overview

Move generation in the ArGO chess engine is a critical component that directly
impacts search performance. The system generates all legal moves for a given
position using bitboards and pre-computed attack tables. Move generation must
be not only correct but extremely fast, as it's called millions of times during
a typical search.

## Move Encoding

Moves are encoded as 32-bit unsigned integers to efficiently store all
necessary information:

```txt
0000 0000 0000 0000 0000 0000 0011 1111   from square            (6 bits)     0x0000003F
0000 0000 0000 0000 0000 1111 1100 0000   to square              (6 bits)     0x00000FC0
0000 0000 0000 0000 1111 0000 0000 0000   moving piece           (4 bits)     0x0000F000
0000 0000 0000 1111 0000 0000 0000 0000   move type              (4 bits)     0x000F0000
0000 0000 1111 0000 0000 0000 0000 0000   captured piece         (4 bits)     0x00F00000
1111 1111 0000 0000 0000 0000 0000 0000   score bits (for ordering) (8 bits)  0xFF000000
```

This compact encoding allows for:

- 64 possible source and destination squares
- 16 piece types (white/black pieces)
- 16 move types (quiet, capture, castling, en passant, promotions)
- 16 captured piece types
- Move scoring for move ordering

## Move Types

The engine defines various move types as constants:

```go
const (
    Quiet                  Type = iota // 0: quiet (non-capturing, non-special) move
    DoublePawnPush                     // 1: double pawn push (pawn advances two squares)
    KingCastle                         // 2: kingside castling move
    QueenCastle                        // 3: queenside castling move
    Capture                            // 4: capture move (non-en passant)
    EnPassant                          // 5: en passant capture move
    // ... more types for promotions ...
    KnightPromotion                    // 8: knight promotion (non-capture)
    BishopPromotion                    // 9: bishop promotion (non-capture)
    RookPromotion                      // 10: rook promotion (non-capture)
    QueenPromotion                     // 11: queen promotion (non-capture)
    KnightPromotionCapture             // 12: knight promotion with capture
    BishopPromotionCapture             // 13: bishop promotion with capture
    RookPromotionCapture               // 14: rook promotion with capture
    QueenPromotionCapture              // 15: queen promotion with capture
)
```

## Move Generation Process

The move generation process includes several steps:

### 1. Main Generation Function

```go
func (b *Board) GenerateMoves() []move.Move {
    result := make([]move.Move, 0, 10)

    // Generate moves for each piece type
    for piece := WP; piece <= BK; piece++ {
        // Check if piece belongs to the side to move
        if (piece <= WK && b.SideToMove == color.WHITE) ||
           (piece >= BP && b.SideToMove == color.BLACK) {
            // Generate piece-specific moves
            // ... (pawn moves, knight moves, etc.)
        }
    }

    return result
}
```

### 2. Piece-Specific Generation

Different piece types have specialized generation logic:

#### Pawns

```go
// For pawns, handle:
// - Single and double pushes
// - Captures (diagonal attacks)
// - En passant captures
// - Promotions
```

#### Knights

```go
// For knights:
// 1. Get knight bitboard
// 2. For each knight, get pre-computed attacks
// 3. Remove squares occupied by friendly pieces
// 4. Add moves for remaining attack squares
```

#### Sliding Pieces (Bishops, Rooks, Queens)

```go
// For sliding pieces:
// 1. Get piece bitboard
// 2. For each piece, get attacks using magic bitboard lookup
// 3. Remove squares occupied by friendly pieces
// 4. Add moves for remaining attack squares
```

#### Kings

```go
// For kings:
// 1. Get king bitboard
// 2. Get pre-computed king attacks
// 3. Remove squares occupied by friendly pieces
// 4. Add moves for remaining attack squares
// 5. Check castling rights and add castling moves if available
```

### 3. Castling Move Generation

Castling requires special handling:

```go
// Check kingside castling
if (b.Castlings & ShortW) != 0 {
    // Verify squares between king and rook are empty
    if !b.Occupancies[color.BOTH].Test(F1) && !b.Occupancies[color.BOTH].Test(G1) {
        // Verify king and traversed squares are not under attack
        if !b.IsSquareAttacked(E1, color.BLACK) && !b.IsSquareAttacked(F1, color.BLACK) {
            // Add castling move
            result = append(result, move.EncodeMove(E1, G1, WK, move.KingCastle, 0))
        }
    }
}
```

### 4. Legal Move Filtering

ArGO generates "pseudo-legal" moves first, then filters out illegal moves
that would leave the king in check:

```go
func (b *Board) MakeMove(m move.Move, moveFlag int) bool {
    // Preserve board state for potential unmake
    copyB := b.CopyBoard()

    // Make the move on the board
    // ... (update bitboards, piece placements, etc.)

    // Check if the move leaves the own king in check
    kingPos := findKingPosition(b.SideToMove.Opp())
    if b.IsSquareAttacked(kingPos, b.SideToMove) {
        // Illegal move - restore the board state
        b.TakeBack(copyB)
        return false
    }

    // Legal move - update additional board state
    // ... (update castling rights, en passant, etc.)
    return true
}
```

## Attack Detection

A critical component of move generation is determining which squares are
attacked. This is used for both move legality checking and tactical evaluation:

```go
func (b *Board) IsSquareAttacked(sq int, attackingSide color.Color) bool {
    // Check pawn attacks
    if attackingSide == color.WHITE {
        if attacks.PawnAttacks[color.BLACK][sq] & b.Bitboards[WP] != 0 {
            return true
        }
    } else {
        if attacks.PawnAttacks[color.WHITE][sq] & b.Bitboards[BP] != 0 {
            return true
        }
    }

    // Check knight attacks
    knightBB := b.Bitboards[WN]
    if attackingSide == color.BLACK {
        knightBB = b.Bitboards[BN]
    }
    if attacks.KnightAttacks[sq] & knightBB != 0 {
        return true
    }

    // Check king attacks
    kingBB := b.Bitboards[WK]
    if attackingSide == color.BLACK {
        kingBB = b.Bitboards[BK]
    }
    if attacks.KingAttacks[sq] & kingBB != 0 {
        return true
    }

    // Check bishop/queen attacks (diagonal)
    bishopBB := b.Bitboards[WB]
    queenBB := b.Bitboards[WQ]
    if attackingSide == color.BLACK {
        bishopBB = b.Bitboards[BB]
        queenBB = b.Bitboards[BQ]
    }
    bishopAttacks := attacks.GetBishopAttacks(sq, b.Occupancies[color.BOTH])
    if bishopAttacks & (bishopBB | queenBB) != 0 {
        return true
    }

    // Check rook/queen attacks (horizontal/vertical)
    rookBB := b.Bitboards[WR]
    if attackingSide == color.BLACK {
        rookBB = b.Bitboards[BR]
    }
    rookAttacks := attacks.GetRookAttacks(sq, b.Occupancies[color.BOTH])
    if rookAttacks & (rookBB | queenBB) != 0 {
        return true
    }

    return false
}
```

## Capture-Only Move Generation

For quiescence search, the engine needs to generate only capturing moves.
This uses a specialized function for better performance:

```go
func (b *Board) GenerateCaptures() []move.Move {
    result := make([]move.Move, 0, 10)

    // Generate only capturing moves for each piece type
    // ... (similar to regular move generation but only for captures)

    return result
}
```

## Staged Move Generation

In actual search, moves can be generated in stages to improve move ordering:

1. **Hash Move**: The best move from the transposition table
2. **Good Captures**: Captures that appear good based on MVV-LVA
3. **Killer Moves**: Quiet moves that caused beta cutoffs at the same ply
4. **Quiet Moves**: Remaining non-capturing moves ordered by history

This approach is more efficient than generating all moves at once
and then sorting them.

## Performance Optimizations

ArGO's move generation includes several optimizations:

### 1. Bitboard Operations

Using bitboard operations allows many potential moves to be generated simultaneously:

```go
// Get all knight moves in one operation
knightBB := b.Bitboards[WN]
knightMoves := bitboard.Bitboard(0)

// For each knight
for knightBB != 0 {
    fromSq := knightBB.FirstOne()

    // Get attack bitboard from pre-computed table
    attacks := attacks.KnightAttacks[fromSq]

    // Remove friendly piece locations
    attacks &= ^b.Occupancies[color.WHITE]

    // Now 'attacks' contains all valid destination squares
    // ...
}
```

### 2. Pre-computed Attack Tables

All attack patterns are pre-computed at initialization:

```go
func InitKnightAttacks() {
    for sq := A8; sq <= H1; sq++ {
        KnightAttacks[sq] = generateKnightAttacks(sq)
    }
}
```

### 3. Move Reuse

When possible, the engine reuses previously generated moves rather
than regenerating them.

### 4. Memory Management

The move generation system avoids unnecessary memory allocations:

```go
// Pre-allocate with a reasonable capacity
result := make([]move.Move, 0, 10)
```

### 5. Pseudolegal Move Generation

ArGO generates pseudolegal moves first, then filters out illegal
ones during the `MakeMove` step. This is more efficient than
checking legality during generation.

## Special Position Handling

### 1. Check Detection

The engine has an optimized function to detect if the king is in check:

```go
func (b *Board) InCheck() bool {
    var kingPos int
    var kingBB bitboard.Bitboard

    // Find king position for the side to move
    if b.SideToMove == color.WHITE {
        kingBB = b.Bitboards[WK]
        kingPos = kingBB.FirstOne()
        return b.IsSquareAttacked(kingPos, color.BLACK)
    }

    kingBB = b.Bitboards[BK]
    kingPos = kingBB.FirstOne()
    return b.IsSquareAttacked(kingPos, color.WHITE)
}
```

### 2. Checkmate and Stalemate Detection

The engine can detect terminal positions:

```go
func (b *Board) IsCheckmate() bool {
    if !b.InCheck() {
        return false
    }

    // Generate all possible moves
    moves := b.GenerateMoves()

    // Try each move to see if it gets us out of check
    for _, mv := range moves {
        copyB := b.CopyBoard()
        if b.MakeMove(mv, AllMoves) {
            b.TakeBack(copyB)
            return false
        }
        b.TakeBack(copyB)
    }

    // No legal moves while in check => checkmate
    return true
}

func (b *Board) IsStalemate() bool {
    if b.InCheck() {
        return false
    }

    // Generate all possible moves
    moves := b.GenerateMoves()

    // Try each move to see if it's legal
    for _, mv := range moves {
        copyB := b.CopyBoard()
        if b.MakeMove(mv, AllMoves) {
            b.TakeBack(copyB)
            return false // Found a legal move, not stalemate
        }
        b.TakeBack(copyB)
    }

    // No legal moves while not in check => stalemate
    return true
}
```

### 3. Insufficient Material

The engine can detect draws by insufficient material:

```go
func (b *Board) IsInsufficientMaterial() bool {
    // Get piece counts
    whitePieceCount := (b.Bitboards[WP] | b.Bitboards[WR] | b.Bitboards[WQ]).Count()
    blackPieceCount := (b.Bitboards[BP] | b.Bitboards[BR] | b.Bitboards[BQ]).Count()

    // If any pawns, rooks, or queens exist, there's sufficient material
    if whitePieceCount > 0 || blackPieceCount > 0 {
        return false
    }

    // Count minor pieces
    whiteKnights := b.Bitboards[WN].Count()
    blackKnights := b.Bitboards[BN].Count()
    whiteBishops := b.Bitboards[WB].Count()
    blackBishops := b.Bitboards[BB].Count()

    // Check draw conditions (K vs K, K+minor vs K, etc.)
    // ...

    return true
}
```

## Validation

Move generation is validated using perft (performance test) function that
counts the number of positions reachable at a given depth:

```go
func PerftDriver(b *Board, depth int) int64 {
    if depth == 0 {
        return 1
    }

    nodes := int64(0)
    moves := b.GenerateMoves()

    for _, mv := range moves {
        copyB := b.CopyBoard()
        if !b.MakeMove(mv, AllMoves) {
            continue
        }

        nodes += PerftDriver(b, depth-1)
        b.TakeBack(copyB)
    }

    return nodes
}
```

This function is compared against known-correct node counts for standard test
positions to verify move generation correctness.
