package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

func initAllAttacks() {}

func InitSliderPiecesAttacks(piece int) {
	// loop over 64 board squares
	for sq := constants.A8; sq <= constants.H1; sq++ {
		// init bishop & rook masks
		bishopMasks[sq] = generateBishopPossibleBlockers(sq)
		rookMasks[sq] = generateRookPossibleBlockers(sq)

		// init current mask
		var attackMask bitboard.Bitboard

		if piece == constants.Bishop {
			attackMask = generateBishopPossibleBlockers(sq)
		} else {
			attackMask = generateRookPossibleBlockers(sq)
		}

		// count attack mask bits
		bitCount := attackMask.Count()

		// occupancy variations count
		occupancyVariations := 1 << bitCount

		// loop over occupancy variations
		for count := 0; count < occupancyVariations; count++ {
			if piece == constants.Bishop {
				occupancy := SetOccupancy(count, bitCount, attackMask)

				magicIndex := occupancy * bishopMagicNumbers[sq] >> (64 - bishopRelevantBits[sq])
				BishopAttacks[sq][magicIndex] = generateBishopAttacks(sq, occupancy)
			} else {

				occupancy := SetOccupancy(count, bitCount, attackMask)

				magicIndex := occupancy * rookMagicNumbers[sq] >> (64 - rookRelevantBits[sq])
				RookAttacks[sq][magicIndex] = generateRookAttacks(sq, occupancy)
			}
		}
	}
}
