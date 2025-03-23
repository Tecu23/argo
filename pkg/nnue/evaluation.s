// Copyright (C) 2025 Tecu23
// Port of Koivisto evaluation, licensed under GNU GPL v3

// File: evaluation.s
#include "textflag.h"

// func computeScoreASM(accActive, accInactive []int16, hiddenWeights []int16, hiddenBias int32) int32
TEXT ·computeScoreASM(SB), NOSPLIT, $0-40
	// Input parameters:
	// accActive     +0(FP)
	// accActive_len +8(FP)
	// accActive_cap +16(FP)
	// accInactive   +24(FP)
	// accInactive_len +32(FP)
	// accInactive_cap +40(FP)
	// hiddenWeights +48(FP)
	// hiddenWeights_len +56(FP)
	// hiddenWeights_cap +64(FP)
	// hiddenBias    +72(FP)
	// Return value:  +80(FP)

	// Load arguments
	MOVQ accActive+0(FP), SI      // Load accActive slice pointer
	MOVQ accActive_len+8(FP), CX  // Load length of accActive (HiddenSize)
	MOVQ accInactive+24(FP), DI   // Load accInactive slice pointer
	MOVQ hiddenWeights+48(FP), R8 // Load hiddenWeights slice pointer
	MOVL hiddenBias+72(FP), AX    // Load hidden bias

	// Initialize sum with the bias value
	MOVL AX, R10

	// Calculate number of 8-element chunks to process (CX / 8)
	MOVQ CX, R9
	SHRQ $3, R9 // R9 = CX / 8 (divide by 8)

	// Initialize accumulator for vectorized sums
	PXOR X7, X7 // X7 = 0 (32-bit accumulator)

	// Process 8 elements at a time using SIMD
	XORQ R11, R11 // R11 = loop counter (i) for chunks

vector_loop:
	CMPQ R11, R9         // Compare i with chunk count
	JGE  vector_loop_end // If i >= chunks, exit loop

	// Calculate offset for this chunk
	MOVQ R11, R12
	SHLQ $4, R12  // R12 = i * 16 (8 elements * 2 bytes per int16)

	// Load 8 elements from accActive into X0
	MOVOU (SI)(R12*1), X0 // X0 = 8 consecutive int16 from accActive

	// Apply ReLU to X0 (max(0, X0))
	PXOR   X2, X2 // X2 = 0
	PMAXSW X2, X0 // X0 = max(0, X0) (ReLU applied to all 8 elements)

	// Load corresponding weights for active side
	MOVOU (R8)(R12*1), X1 // X1 = 8 consecutive int16 from hiddenWeights

	// Multiply and accumulate
	PMADDWL X1, X0 // X0 = X0 * X1 (multiply and add adjacent pairs)
	PADDD   X0, X7 // X7 += X0 (accumulate results)

	// Calculate offset for inactive side weights
	MOVQ CX, R13
	SHLQ $1, R13  // R13 = HiddenSize * 2 (in bytes)
	ADDQ R12, R13 // R13 = offset for inactive weights

	// Load 8 elements from accInactive into X3
	MOVOU (DI)(R12*1), X3 // X3 = 8 consecutive int16 from accInactive

	// Apply ReLU to X3 (max(0, X3))
	PXOR   X5, X5 // X5 = 0
	PMAXSW X5, X3 // X3 = max(0, X3) (ReLU applied to all 8 elements)

	// Load corresponding weights for inactive side
	MOVOU (R8)(R13*1), X4 // X4 = 8 consecutive int16 from hiddenWeights[HiddenSize+i]

	// Multiply and accumulate
	PMADDWL X4, X3 // X3 = X3 * X4 (multiply and add adjacent pairs)
	PADDD   X3, X7 // X7 += X3 (accumulate results)

	INCQ R11         // i++
	JMP  vector_loop // Continue loop

vector_loop_end:
	// Horizontal sum of X7 (4 int32 values in X7 -> 1 int32 value)
	PSHUFD $0xE, X7, X6 // Shuffle to add upper 2 elements to lower 2 elements
	PADDD  X6, X7       // X7 now has sum in lower 2 elements
	PSHUFD $0x1, X7, X6 // Shuffle to add second element to first element
	PADDD  X6, X7       // X7 now has total sum in lowest element

	// Extract the lowest 32 bits from X7 and add to accumulated sum
	MOVD X7, AX
	ADDL AX, R10 // Add SIMD results to R10

	// Calculate start index for remaining elements
	MOVQ R9, R11
	SHLQ $3, R11 // R11 = chunks * 8

	// Handle remaining elements one by one
	CMPQ R11, CX // Compare i with HiddenSize
	JGE  done    // If i >= HiddenSize, exit loop

remainder_loop:
	// Process active side
	MOVW  (SI)(R11*2), AX // Load accActive[i] into AX
	TESTW AX, AX          // Test if AX > 0
	JLE   skip_active     // If AX <= 0, skip active contribution

	// Convert int16 to int32 and multiply by hiddenWeights[i]
	MOVWLSX AX, AX          // Sign-extend AX from 16-bit to 32-bit
	MOVW    (R8)(R11*2), DX // Load hiddenWeights[i]
	MOVWLSX DX, DX          // Sign-extend DX from 16-bit to 32-bit
	IMULL   DX, AX          // AX = accActive[i] * hiddenWeights[i]
	ADDL    AX, R10         // sum += AX

skip_active:
	// Process inactive side
	MOVW  (DI)(R11*2), AX // Load accInactive[i] into AX
	TESTW AX, AX          // Test if AX > 0
	JLE   skip_inactive   // If AX <= 0, skip inactive contribution

	// Convert int16 to int32 and multiply by hiddenWeights[i+HiddenSize]
	MOVWLSX AX, AX          // Sign-extend AX from 16-bit to 32-bit
	MOVQ    R11, R13        // Copy i to R13
	ADDQ    CX, R13         // R13 = i + HiddenSize
	MOVW    (R8)(R13*2), DX // Load hiddenWeights[i+HiddenSize]
	MOVWLSX DX, DX          // Sign-extend DX from 16-bit to 32-bit
	IMULL   DX, AX          // AX = accInactive[i] * hiddenWeights[i+HiddenSize]
	ADDL    AX, R10         // sum += AX

skip_inactive:
	INCQ R11            // i++
	CMPQ R11, CX        // Compare i with HiddenSize
	JL   remainder_loop // If i < HiddenSize, continue loop

done:
	MOVL R10, ret+80(FP) // Store the result
	RET

