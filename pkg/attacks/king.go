package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

var KingAttacks [64]bitboard.Bitboard

func InitKingAttacks() {
	for sq := constants.A8; sq <= constants.H1; sq++ {
		KingAttacks[sq] = generateKingAttacks(sq)
	}
}

func generateKingAttacks(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	b := bitboard.Bitboard(0)
	b.Set(square)

	attacks |= (b & ^constants.FileA) << constants.NW
	attacks |= b << constants.N
	attacks |= (b & ^constants.FileH) << constants.NE

	attacks |= (b & ^constants.FileH) << constants.E

	attacks |= (b & ^constants.FileA) >> constants.E

	attacks |= (b & ^constants.FileH) >> constants.NW
	attacks |= b >> constants.N
	attacks |= (b & ^constants.FileA) >> constants.NE

	return attacks
}
