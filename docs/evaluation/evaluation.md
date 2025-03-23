# NNUE Evaluation System

## Introduction

The ArGO chess engine uses an Efficiently Updated Neural Network (NNUE)
evaluation function, adapted from the Koivisto chess engine. NNUE represents
a breakthrough in chess position evaluation, combining the pattern recognition
capabilities of neural networks with the efficiency needed for deep search.

## NNUE Architecture

ArGO's NNUE implementation uses a two-layer neural network:

1. **Input Layer**: Feature vectors based on piece positions relative to kings
2. **Hidden Layer**: 512 neurons with ReLU activation function
3. **Output Layer**: Single scalar value representing the position evaluation

The network is structured to allow for efficient incremental updates when
making/unmaking moves - hence the name "Efficiently Updated".

## Feature Representation

### Input Features

The input layer uses a feature representation based on:

```txt
Input Features = Piece Type × Color × Square × King Square Bucket
```

Where:

- **Piece Type**: 6 piece types (pawn, knight, bishop, rook, queen, king)
- **Color**: 2 colors (white, black)
- **Square**: 64 squares on the board
- **King Square Bucket**: 16 buckets/zones for king placement

This results in 6 × 2 × 64 × 16 = 12,288 potential features, though only a
small subset is active in any given position.

### King Bucketing

King positions are "bucketed" to reduce the network size:

```go
func KingSquareIndex(kingSquare, perspective int) int {
    // Maps king position to one of 16 buckets for more efficient representation
}
```

This bucketing system captures king safety patterns without requiring a separate
network for every possible king position.

## Efficient Updates

The core innovation of NNUE is incremental updates, avoiding the need to recompute
the entire network for each move. This is implemented through several key classes:

### Accumulator

The `Accumulator` struct tracks the summation of input weights for all active features:

```go
type Accumulator struct {
    Summation [2][HiddenSize]int16 // [color][neuron] stores the sum for each hidden neuron
}
```

When pieces move, rather than recomputing the entire sum:

1. **Subtract weights** for features that are no longer active
2. **Add weights** for new features that become active

For example, when a piece moves from square A to B:

```go
func SetUnsetPiece(input, output *Accumulator, side int, set, unset FeatureIndex) {
    // Get indices for the features
    idx1 := set.Get(side)
    idx2 := unset.Get(side)

    // Use assembly-optimized function to update the accumulator
    setUnsetPieceASM(
        input.Summation[side][:],
        output.Summation[side][:],
        InputWeights[idx1][:],
        InputWeights[idx2][:],
    )
}
```

### Special Handling for Different Move Types

The `ProcessMove` function handles different move types with specialized update functions:

- Standard moves → `SetUnsetPiece`
- Captures → `SetUnsetUnsetPiece`
- Castling → `SetSetUnsetUnsetPiece`
- Promotions → Special handling

### Assembly Optimization

Performance-critical sections use hand-tuned assembly code for maximum speed:

```assembly
// func addWeightsToAccumulatorASM(add bool, src, target, weights []int16)
TEXT ·addWeightsToAccumulatorASM(SB), NOSPLIT, $0
    // ... assembly code for efficient weight updates
```

These optimizations use SIMD instructions to process multiple neurons in parallel,
providing a significant speed boost on AMD64 architectures.

## Evaluation Process

During search, position evaluation follows this process:

1. **Initialize Accumulators**: At the beginning, compute full accumulators for the
   starting position
2. **Incremental Updates**: For each move in the search tree, update accumulators
   incrementally
3. **ReLU Activation**: Apply ReLU (max(0,x)) to all hidden neurons
4. **Output Calculation**: Compute the final evaluation score through
   the output layer
5. **Scaling**: Apply phase-dependent scaling (middlegame vs endgame)

```go
func (e *Evaluator) Evaluate(b *board.Board) int {
    // Calculate game phase based on piece material
    phase := calculatePhase(b)

    // Get neural network evaluation
    rawEval := e.eval(int(b.SideToMove))

    // Scale based on game phase
    return int(
        (evaluationMgScalar - phase*(evaluationMgScalar-evaluationEgScalar)) * float64(rawEval),
    )
}
```

## Game Phase Calculation

The evaluation function scales between middlegame and endgame evaluations
based on the remaining material:

```go
var phaseValues = [5]float64{
    0.552938, 1.55294, 1.50862, 2.64379, 4.0, // Pawn, Knight, Bishop, Rook, Queen
}

// Calculate phase based on remaining material
phase := phaseSum
phase -= float64((b.Bitboards[WP] | b.Bitboards[BP]).Count()) * phaseValues[Pawn]
phase -= float64((b.Bitboards[WN] | b.Bitboards[BN]).Count()) * phaseValues[Knight]
// ... and so on for other pieces
phase /= phaseSum // Normalize to 0-1 range
```

## Caching and Performance Optimizations

Several caching mechanisms improve performance:

1. **Accumulator Table**: Caches accumulators indexed by king positions
2. **History Stack**: Maintains a stack of accumulators for position unmaking
3. **Lazy Initialization**: Avoids recomputing accumulators when possible

## Network Parameters

The neural network parameters:

- **Input Layer**: 12,288 potential input features
- **Hidden Layer**: 512 neurons
- **Output Layer**: 1 neuron
- **Weights Storage**: Quantized as int16 with scaling
  factors for accuracy/size balance
- **Memory Footprint**: ~12MB for all weights and biases

## Advantages Over Classical Evaluation

NNUE offers several advantages over traditional handcrafted evaluation functions:

1. **Pattern Recognition**: Better at recognizing subtle positional patterns
2. **Strength**: Typically provides 50-100 Elo improvement
3. **Consistency**: Fewer blind spots in evaluation
4. **Incrementality**: Updates efficiently during search
5. **Trainability**: Can be improved through self-play and supervised learning

## Training Considerations

While ArGO uses weights ported from Koivisto, NNUE networks are typically
trained with:

1. **Large Position Dataset**: Millions of high-quality positions
2. **Supervised Learning**: Using strong engine evaluations as targets
3. **Reinforcement Learning**: Fine-tuning through self-play
4. **Quantization**: Converting float weights to integer for faster inference

## Performance Impact on Search

The efficiency of NNUE evaluation allows the engine to:

1. **Search Deeper**: Evaluating positions faster enables deeper search
2. **Improve Selectivity**: Better evaluations lead to more accurate pruning decisions
3. **Handle Complex Positions**: Recognize subtle patterns that traditional
   evaluation might miss
