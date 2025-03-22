// File: evaluation.s
#include "textflag.h"

// computeScoreASM calculates the dot product with ReLU and bias
// func computeScoreASM(accActive, accInactive []int16, hiddenWeights []int16, hiddenBias int32) int32
TEXT ·computeScoreASM(SB), 4, $0 // 4 = NOSPLIT flag
	// Load parameters
	MOVQ accActive+0(FP), SI      // Active side accumulator
	MOVQ accInactive+24(FP), DI   // Inactive side accumulator
	MOVQ hiddenWeights+48(FP), R8 // Hidden weights
	MOVQ accActive+8(FP), CX      // Length (HiddenSize)
	MOVL hiddenBias+72(FP), AX    // Hidden bias

	// Initialize sum with bias
	MOVL AX, R9 // R9 will hold our sum

	// Process 2 elements per iteration to reduce loop overhead
	SHRQ $1, CX    // Divide length by 2
	MOVQ CX, R13   // Save half length for later
	JZ   remainder // If length is 0 after division, go to remainder

	// Initialize index
	XORQ R10, R10 // R10 = 0 (loop index for pairs)

loop_pairs:
	// Check if we're done with pairs
	CMPQ R10, R13  // Compare index with half-length
	JGE  remainder // If index >= half-length, process remainder

	// Compute indexes for the pair
	MOVQ R10, R14
	SHLQ $1, R14     // R14 = 2*i (index for first element in pair)
	LEAQ 1(R14), R15 // R15 = 2*i+1 (index for second element in pair)

	// Process accActive[2*i]
	MOVWQZX (SI)(R14*2), AX // Load accActive[2*i]
	TESTW   AX, AX          // Check if accActive[2*i] > 0
	JLE     skip_active1    // Skip if accActive[2*i] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	MOVWQZX (R8)(R14*2), AX // Load hiddenWeights[2*i]
	MOVQ    AX, R12
	SHLQ    $48, R12        // Shift left by 48 bits
	SARQ    $48, R12        // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accActive[2*i] * hiddenWeights[2*i]
	ADDL  R11, R9  // Add to sum

skip_active1:
	// Process accActive[2*i+1]
	MOVWQZX (SI)(R15*2), AX // Load accActive[2*i+1]
	TESTW   AX, AX          // Check if accActive[2*i+1] > 0
	JLE     skip_active2    // Skip if accActive[2*i+1] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	MOVWQZX (R8)(R15*2), AX // Load hiddenWeights[2*i+1]
	MOVQ    AX, R12
	SHLQ    $48, R12        // Shift left by 48 bits
	SARQ    $48, R12        // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accActive[2*i+1] * hiddenWeights[2*i+1]
	ADDL  R11, R9  // Add to sum

skip_active2:
	// Process accInactive[2*i]
	MOVWQZX (DI)(R14*2), AX // Load accInactive[2*i]
	TESTW   AX, AX          // Check if accInactive[2*i] > 0
	JLE     skip_inactive1  // Skip if accInactive[2*i] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	// Calculate the offset for hiddenWeights for inactive side
	MOVQ    accActive+8(FP), R12 // Get original length (HiddenSize)
	ADDQ    R14, R12             // Add index (2*i) to get offset
	MOVWQZX (R8)(R12*2), AX      // Load hiddenWeights[2*i+HiddenSize]

	MOVQ AX, R12
	SHLQ $48, R12 // Shift left by 48 bits
	SARQ $48, R12 // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accInactive[2*i] * hiddenWeights[2*i+HiddenSize]
	ADDL  R11, R9  // Add to sum

skip_inactive1:
	// Process accInactive[2*i+1]
	MOVWQZX (DI)(R15*2), AX // Load accInactive[2*i+1]
	TESTW   AX, AX          // Check if accInactive[2*i+1] > 0
	JLE     skip_inactive2  // Skip if accInactive[2*i+1] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	// Calculate the offset for hiddenWeights for inactive side
	MOVQ    accActive+8(FP), R12 // Get original length (HiddenSize)
	ADDQ    R15, R12             // Add index (2*i+1) to get offset
	MOVWQZX (R8)(R12*2), AX      // Load hiddenWeights[2*i+1+HiddenSize]

	MOVQ AX, R12
	SHLQ $48, R12 // Shift left by 48 bits
	SARQ $48, R12 // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accInactive[2*i+1] * hiddenWeights[2*i+1+HiddenSize]
	ADDL  R11, R9  // Add to sum

skip_inactive2:
	INCQ R10        // index++
	JMP  loop_pairs // Continue with next pair

remainder:
	// Check if there's a remaining element (odd length)
	MOVQ  accActive+8(FP), R13 // Reload original length
	MOVQ  R13, R14
	ANDQ  $1, R14              // R14 = length % 2 (1 if odd, 0 if even)
	TESTQ R14, R14
	JZ    done                 // If length is even, we're done

	// Process the last element
	MOVQ R13, R14
	DECQ R14      // R14 = length - 1 (index of last element)

	// Process accActive[last]
	MOVWQZX (SI)(R14*2), AX  // Load accActive[last]
	TESTW   AX, AX           // Check if accActive[last] > 0
	JLE     skip_active_last // Skip if accActive[last] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	MOVWQZX (R8)(R14*2), AX // Load hiddenWeights[last]
	MOVQ    AX, R12
	SHLQ    $48, R12        // Shift left by 48 bits
	SARQ    $48, R12        // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accActive[last] * hiddenWeights[last]
	ADDL  R11, R9  // Add to sum

skip_active_last:
	// Process accInactive[last]
	MOVWQZX (DI)(R14*2), AX // Load accInactive[last]
	TESTW   AX, AX          // Check if accInactive[last] > 0
	JLE     done            // Skip if accInactive[last] <= 0

	// Sign-extend int16 to int32 manually
	MOVQ AX, R11
	SHLQ $48, R11 // Shift left by 48 bits
	SARQ $48, R11 // Shift right arithmetically to sign-extend

	// Calculate the offset for hiddenWeights for inactive side
	MOVQ    accActive+8(FP), R12 // Get original length (HiddenSize)
	ADDQ    R14, R12             // Add index (last) to get offset
	MOVWQZX (R8)(R12*2), AX      // Load hiddenWeights[last+HiddenSize]

	MOVQ AX, R12
	SHLQ $48, R12 // Shift left by 48 bits
	SARQ $48, R12 // Shift right arithmetically to sign-extend

	IMULQ R12, R11 // accInactive[last] * hiddenWeights[last+HiddenSize]
	ADDL  R11, R9  // Add to sum

done:
	MOVL R9, ret+80(FP) // Store final sum as return value
	RET

