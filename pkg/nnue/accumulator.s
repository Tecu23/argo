// File: accumulator.s
#include "textflag.h"

// TEXT ·setUnsetPieceASM(SB),NOSPLIT,$0
TEXT ·setUnsetPieceASM(SB), NOSPLIT, $0
	MOVQ input+0(FP), SI         // input slice data pointer
	MOVQ output+24(FP), DI       // output slice data pointer
	MOVQ weightsSet+48(FP), AX   // weights to add
	MOVQ weightsUnset+72(FP), BX // weights to subtract
	MOVQ input+8(FP), CX         // input slice length

	XORQ R8, R8 // index = 0

loop:
	CMPQ R8, CX // compare index with length
	JGE  done   // if index >= length, we're done

	// Process 8 int16 values at once with SSE2
	CMPQ CX, R8
	JL   remainder // if fewer than 8 elements left, handle remainder
	SUBQ $8, CX
	CMPQ CX, R8
	JL   remainder
	ADDQ $8, CX

	// Load 8 int16 values (16 bytes) into XMM registers
	MOVOU (SI)(R8*2), X0 // input
	MOVOU (AX)(R8*2), X1 // weightsSet
	MOVOU (BX)(R8*2), X2 // weightsUnset

	// Perform operations: output = input + weightsSet - weightsUnset
	PADDW X1, X0 // X0 = input + weightsSet
	PSUBW X2, X0 // X0 = (input + weightsSet) - weightsUnset

	// Store result back to memory
	MOVOU X0, (DI)(R8*2)

	ADDQ $8, R8 // index += 8
	JMP  loop

remainder:
	CMPQ R8, CX // check if we're done
	JGE  done

	// Handle remaining elements one by one
	MOVW (SI)(R8*2), DX  // load input value
	MOVW (AX)(R8*2), R9  // load weightsSet value
	MOVW (BX)(R8*2), R10 // load weightsUnset value

	ADDW R9, DX  // DX += weightsSet
	SUBW R10, DX // DX -= weightsUnset

	MOVW DX, (DI)(R8*2) // store result

	INCQ R8        // index++
	JMP  remainder

done:
	RET

// TEXT ·setUnsetUnsetPieceASM(SB),NOSPLIT,$0
TEXT ·setUnsetUnsetPieceASM(SB), NOSPLIT, $0
	MOVQ input+0(FP), SI    // input slice data pointer
	MOVQ output+24(FP), DI  // output slice data pointer
	MOVQ set+48(FP), AX     // weights to add
	MOVQ unset1+72(FP), BX  // weights to subtract 1
	MOVQ unset2+72(FP), R11 // weights to subtract 2
	MOVQ input+8(FP), CX    // input slice length

	XORQ R8, R8 // index = 0

loop:
	CMPQ R8, CX // compare index with length
	JGE  done   // if index >= length, we're done

	// Process 8 int16 values at once with SSE2
	CMPQ CX, R8
	JL   remainder // if fewer than 8 elements left, handle remainder
	SUBQ $8, CX
	CMPQ CX, R8
	JL   remainder
	ADDQ $8, CX

	// Load 8 int16 values (16 bytes) into XMM registers
	MOVOU (SI)(R8*2), X0  // input
	MOVOU (AX)(R8*2), X1  // weightsSet
	MOVOU (BX)(R8*2), X2  // weightsUnset
	MOVOU (R11)(R8*2), X3 // weightsUnset

	// Perform operations: output = input + weightsSet - weightsUnset
	PADDW X1, X0 // X0 = input + weightsSet
	PSUBW X2, X0 // X0 = (input + weightsSet) - weightsUnset1
	PSUBW X3, X0 // X0 = (input + weightsSet - weightUnset1) - weightsUnset2

	// Store result back to memory
	MOVOU X0, (DI)(R8*2)

	ADDQ $8, R8 // index += 8
	JMP  loop

remainder:
	CMPQ R8, CX // check if we're done
	JGE  done

	// Handle remaining elements one by one
	MOVW (SI)(R8*2), DX   // load input value
	MOVW (AX)(R8*2), R9   // load weightsSet value
	MOVW (BX)(R8*2), R10  // load weightsUnset1 value
	MOVW (R11)(R8*2), R12 // load weightsUnset2 value

	ADDW R9, DX  // DX += weightsSet
	SUBW R10, DX // DX -= weightsUnset1
	SUBW R12, DX // DX -= weightsUnset2

	MOVW DX, (DI)(R8*2) // store result

	INCQ R8        // index++
	JMP  remainder

done:
	RET

// func addWeightsToAccumulatorASM(add bool, src, target, weights []int16)
TEXT ·addWeightsToAccumulatorASM(SB), NOSPLIT, $0
	MOVBQZX add+0(FP), AX            // Load boolean flag 'add' into AX
	MOVQ    src_data+8(FP), SI       // src slice data pointer
	MOVQ    src_len+16(FP), CX       // src slice length
	MOVQ    target_data+32(FP), DI   // target slice data pointer
	MOVQ    weights_data+56(FP), R10 // weights slice data pointer

	XORQ R8, R8 // index = 0

	// Check if we should add or subtract
	TESTQ AX, AX
	JZ    subtract_loop // If add == false, jump to subtract loop

add_loop:
	CMPQ R8, CX // compare index with length
	JGE  done   // if index >= length, we're done

	// Process 8 int16 values at once with SSE2 if possible
	CMPQ CX, R8
	JL   add_remainder // if fewer than 8 elements left, handle remainder
	SUBQ $8, CX
	CMPQ CX, R8
	JL   add_remainder
	ADDQ $8, CX

	// Load 8 int16 values (16 bytes) into XMM registers
	MOVOU (SI)(R8*2), X0  // src
	MOVOU (R10)(R8*2), X1 // weights

	// Perform addition: target = src + weights
	PADDW X1, X0 // X0 = src + weights

	// Store result back to memory
	MOVOU X0, (DI)(R8*2)

	ADDQ $8, R8   // index += 8
	JMP  add_loop

add_remainder:
	CMPQ R8, CX // check if we're done
	JGE  done

	// Handle remaining elements one by one
	MOVW (SI)(R8*2), DX  // load src value
	MOVW (R10)(R8*2), R9 // load weights value
	ADDW R9, DX          // DX += weights
	MOVW DX, (DI)(R8*2)  // store result

	INCQ R8            // index++
	JMP  add_remainder

subtract_loop:
	CMPQ R8, CX // compare index with length
	JGE  done   // if index >= length, we're done

	// Process 8 int16 values at once with SSE2 if possible
	CMPQ CX, R8
	JL   sub_remainder // if fewer than 8 elements left, handle remainder
	SUBQ $8, CX
	CMPQ CX, R8
	JL   sub_remainder
	ADDQ $8, CX

	// Load 8 int16 values (16 bytes) into XMM registers
	MOVOU (SI)(R8*2), X0  // src
	MOVOU (R10)(R8*2), X1 // weights

	// Perform subtraction: target = src - weights
	PSUBW X1, X0 // X0 = src - weights

	// Store result back to memory
	MOVOU X0, (DI)(R8*2)

	ADDQ $8, R8        // index += 8
	JMP  subtract_loop

sub_remainder:
	CMPQ R8, CX // check if we're done
	JGE  done

	// Handle remaining elements one by one
	MOVW (SI)(R8*2), DX  // load src value
	MOVW (R10)(R8*2), R9 // load weights value
	SUBW R9, DX          // DX -= weights
	MOVW DX, (DI)(R8*2)  // store result

	INCQ R8            // index++
	JMP  sub_remainder

done:
	RET
