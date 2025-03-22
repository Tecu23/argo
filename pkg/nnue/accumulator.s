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
