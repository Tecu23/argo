package attacks

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

// BishopMasks and RookMasks store bitboards representing the "attack mask" from each square
// for bishops and rooks respectively. These masks identify which squares affect the piece’s
// moves from a given square. They are used to extract the subset of occupancy bits that matter.
var (
	bishopMasks [64]bitboard.Bitboard
	rookMasks   [64]bitboard.Bitboard
)

// BishopAttacks and RookAttacks are precomputed tables indexed by square and a transformed occupancy index.
// After combining the occupancy with a magic number and shifting, we get an index into these tables, which
// returns a precomputed bitboard of attacked squares. This allows very fast attack retrieval.
var (
	BishopAttacks [64][512]bitboard.Bitboard
	RookAttacks   [64][4096]bitboard.Bitboard
)

// InitMagic prints out the magic numbers for rooks and bishops after they are found.
// It uses findMagicNumbers() to attempt to find a magic number that uniquely maps occupancy
// patterns to distinct indices. Once found, these can be hardcoded into the engine.
func InitMagic() {
	fmt.Printf("const Bitboard rookMagics[64] = {\n")

	// loop over 64 board squares
	for sq := constants.A1; sq <= constants.H8; sq++ {
		fmt.Printf("    %x,\n", findMagicNumbers(sq, rookRelevantBits[sq], constants.Rook))
	}

	fmt.Printf("};\n\nconst U64 bishop_magics[64] = {\n")

	// loop over 64 board squares
	for sq := constants.A1; sq <= constants.H8; sq++ {
		fmt.Printf("    %x,\n", findMagicNumbers(sq, bishopRelevantBits[sq], constants.Bishop))
	}

	fmt.Printf("};\n\n")
}

// generateRandomUint32Number and generateRandomUint64Number produce pseudo-random numbers
// using a simple XOR shift algorithm. They are used during the magic number finding process.
func generateRandomUint32Number(randomState *uint32) uint32 {
	number := *randomState

	number ^= (number << 13)
	number ^= (number >> 17)
	number ^= (number << 5)

	*randomState = number

	return number
}

func generateRandomUint64Number(randomState *uint32) uint64 {
	var n1, n2, n3, n4 uint64

	n1 = uint64(generateRandomUint32Number(randomState)) & 0xFFFF
	n2 = uint64(generateRandomUint32Number(randomState)) & 0xFFFF
	n3 = uint64(generateRandomUint32Number(randomState)) & 0xFFFF
	n4 = uint64(generateRandomUint32Number(randomState)) & 0xFFFF

	return n1 | (n2 << 16) | (n3 << 32) | (n4 << 48)
}

// generateMagicNumber returns a candidate magic number by AND-ing three random 64-bit numbers.
// The idea is to produce a "dense" number with unpredictable bits, hoping to yield a good magic.
func generateMagicNumber(randomState *uint32) uint64 {
	firstNum := generateRandomUint64Number(randomState)
	secondNum := generateRandomUint64Number(randomState)
	thirdNum := generateRandomUint64Number(randomState)
	return firstNum & secondNum & thirdNum
}

// findMagicNumbers attempts to find a magic number for the given square and piece type (bishop or rook).
// It:
// 1. Generates all occupancy variations for the square using SetOccupancy().
// 2. Precomputes attacks for each occupancy variation.
// 3. Tries random candidates for the magic number.
// 4. For each candidate, checks if it creates a perfect mapping of occupancy to attacks with no collisions.
// 5. Returns the first magic number that succeeds.
func findMagicNumbers(square, relevantBits int, piece int) bitboard.Bitboard {
	// define occupancies array
	occupancy := [4096]bitboard.Bitboard{}

	// define attacks array
	attacks := [4096]bitboard.Bitboard{}

	// define used indices array
	usedAttacks := [4096]bitboard.Bitboard{}

	var maskAttacks bitboard.Bitboard
	// mask piece attack
	if piece == constants.Bishop {
		maskAttacks = generateBishopPossibleBlockers(square)
	} else {
		maskAttacks = generateRookPossibleBlockers(square)
	}

	// occupancy variations
	occupancyVariations := 1 << relevantBits

	// loop over the number of occupancy variations
	for count := 0; count < occupancyVariations; count++ {
		// init occupancies
		occupancy[count] = SetOccupancy(count, relevantBits, maskAttacks)

		// init attacks
		if piece == constants.Bishop {
			attacks[count] = generateBishopAttacks(square, occupancy[count])
		} else {
			attacks[count] = generateRookAttacks(square, occupancy[count])
		}
	}
	randomState := uint32(1804289383)

	// test magic numbers
	for randomCount := 0; randomCount < 100000000; randomCount++ {

		// init magic number candidate
		magic := bitboard.Bitboard(generateMagicNumber(&randomState))

		// skip testing magic number if innappropriate
		if bitboard.Bitboard((maskAttacks*magic)&0xFF00000000000000).Count() < 6 {
			continue
		}

		// reset used attacks array
		usedAttacks = [4096]bitboard.Bitboard{}

		// init count & fail flag
		count, fail := 0, false

		// test magic index
		for count, fail = 0, false; !fail && count < occupancyVariations; count++ {
			// generate magic index
			magicIndex := int((occupancy[count] * magic) >> (64 - relevantBits))

			// if got free index
			if usedAttacks[magicIndex] == 0 {
				// assign corresponding attack map
				usedAttacks[magicIndex] = attacks[count]
			} else if usedAttacks[magicIndex] != attacks[count] {
				fail = true
			}
		}

		// return magic if it works
		if !fail {
			return magic
		}
	}

	// on fail
	fmt.Printf("***Failed***\n")
	return bitboard.Bitboard(0)
}

// SetOccupancy generates an occupancy bitboard (a subset of the given attackMask)
// based on the binary pattern of 'index'.
//
// Parameters:
// - index: Used as a bit pattern to decide which squares from the attackMask are included.
// - bitsInMask: The number of set bits (ones) in attackMask. We will iterate over each set bit.
// - attackMask: A bitboard representing a set of squares (e.g., a ray of attacks).
//
// Process:
//  1. Convert the 'attackMask' into an ordered list of squares by repeatedly extracting the
//     least significant set bit using FirstOne().
//  2. For each extracted bit position (square), check the corresponding bit in 'index' (starting from 0).
//  3. If the bit in 'index' is set, include that square in the final occupancy bitboard.
//
// In other words, 'index' acts like a binary index into all subsets of the bits in attackMask.
// For example, if 'attackMask' has N set bits, then 'index' ranges from 0 to 2^N - 1, allowing
// you to select any subset of those bits.
func SetOccupancy(index, bitsInMask int, attackMask bitboard.Bitboard) bitboard.Bitboard {
	// Initialize an empty occupancy bitboard.
	occupancy := bitboard.Bitboard(0)

	// Iterate over the number of bits in the mask.
	// Each iteration removes one LSB (least significant bit) from attackMask.
	for count := 0; count < bitsInMask; count++ {
		// Extract the position of the least significant set bit in attackMask.
		square := attackMask.FirstOne()

		// Check if the corresponding bit in 'index' is set.
		// If so, add that square to the occupancy bitboard.
		if index&(1<<count) != 0 {
			occupancy |= (1 << square)
		}
	}

	// Return the generated occupancy subset.
	return occupancy
}

// bishopRelevantBits and rookRelevantBits specify how many squares are considered "relevant"
// for occupancy calculation for each square on the board. These bits define how large the
// occupancy subsets are and, hence, how large the attack tables must be.
//
// For example, a rook in the corner has fewer relevant squares (because fewer squares affect its moves)
// while a rook in the center has more relevant squares. Similar logic applies to bishops.
//
// The values in these arrays determine how many bits we shift and how large the lookup tables are for each square.
var bishopRelevantBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var rookRelevantBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// bishopMagicNumbers and rookMagicNumbers arrays store the found magic numbers for each square.
// Once these magic numbers are found (via the findMagicNumbers function or precomputed offline),
// they can be hardcoded to speed up the initialization of the engine.
var bishopMagicNumbers = [64]bitboard.Bitboard{
	0xc085080200420200,
	0x60014902028010,
	0x401240100c201,
	0x580ca104020080,
	0x8434052000230010,
	0x102080208820420,
	0x2188410410403024,
	0x40120805282800,
	0x4420410888208083,
	0x1049494040560,
	0x6090100400842200,
	0x1000090405002001,
	0x48044030808c409,
	0x20802080384,
	0x2012008401084008,
	0x9741088200826030,
	0x822000400204c100,
	0x14806004248220,
	0x30200101020090,
	0x148150082004004,
	0x6020402112104,
	0x4001000290080d22,
	0x2029100900400,
	0x804203145080880,
	0x60a10048020440,
	0xc08080b20028081,
	0x1009001420c0410,
	0x101004004040002,
	0x1004405014000,
	0x10029a0021005200,
	0x4002308000480800,
	0x301025015004800,
	0x2402304004108200,
	0x480110c802220800,
	0x2004482801300741,
	0x400400820a60200,
	0x410040040040,
	0x2828080020011000,
	0x4008020050040110,
	0x8202022026220089,
	0x204092050200808,
	0x404010802400812,
	0x422002088009040,
	0x180604202002020,
	0x400109008200,
	0x2420042000104,
	0x40902089c008208,
	0x4001021400420100,
	0x484410082009,
	0x2002051108125200,
	0x22e4044108050,
	0x800020880042,
	0xb2020010021204a4,
	0x2442100200802d,
	0x10100401c4040000,
	0x2004a48200c828,
	0x9090082014000,
	0x800008088011040,
	0x4000000a0900b808,
	0x900420000420208,
	0x4040104104,
	0x120208c190820080,
	0x4000102042040840,
	0x8002421001010100,
}

var rookMagicNumbers = [64]bitboard.Bitboard{
	0x11800040001481a0,
	0x2040400010002000,
	0xa280200308801000,
	0x100082005021000,
	0x280280080040006,
	0x200080104100200,
	0xc00040221100088,
	0xe00072200408c01,
	0x2002045008600,
	0xa410804000200089,
	0x4081002000401102,
	0x2000c20420010,
	0x800800400080080,
	0x40060010041a0009,
	0x441004442000100,
	0x462800080004900,
	0x80004020004001,
	0x1840420021021081,
	0x8020004010004800,
	0x940220008420010,
	0x2210808008000400,
	0x24808002000400,
	0x803604001019a802,
	0x520000440081,
	0x802080004000,
	0x1200810500400024,
	0x8000100080802000,
	0x2008080080100480,
	0x8000404002040,
	0xc012040801104020,
	0xc015000900240200,
	0x20040200208041,
	0x1080004000802080,
	0x400081002110,
	0x30002000808010,
	0x2000100080800800,
	0x2c0800400800800,
	0x1004800400800200,
	0x818804000210,
	0x340082000a45,
	0x8520400020818000,
	0x2008900460020,
	0x100020008080,
	0x601001000a30009,
	0xc001000408010010,
	0x2040002008080,
	0x11008218018c0030,
	0x20c0080620011,
	0x400080002080,
	0x8810040002500,
	0x400801000200080,
	0x2402801000080480,
	0x204040280080080,
	0x31044090200801,
	0x40c10830020400,
	0x442800100004080,
	0x10080002d005041,
	0x134302820010a2c2,
	0x6202001080200842,
	0x1820041000210009,
	0x1002001008210402,
	0x2000108100402,
	0x10310090a00b824,
	0x800040100944822,
}
