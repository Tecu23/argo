package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

func GetQueenAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
	queenAttacks := bitboard.Bitboard(0)

	bishopOccupancies := occupancy
	rookOccupancies := occupancy

	bishopOccupancies &= bishopMasks[sq]
	bishopOccupancies *= bishopMagicNumbers[sq]
	bishopOccupancies >>= 64 - bishopRelevantBits[sq]

	rookOccupancies &= rookMasks[sq]
	rookOccupancies *= rookMagicNumbers[sq]
	rookOccupancies >>= 64 - rookRelevantBits[sq]

	queenAttacks = BishopAttacks[sq][bishopOccupancies] | RookAttacks[sq][rookOccupancies]

	return queenAttacks
}
