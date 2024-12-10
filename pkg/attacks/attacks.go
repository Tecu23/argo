package attacks

func initAllAttacks() {}

// func InitSliderPiecesAttacks(piece int) {
// 	// loop over 64 board squares
// 	for sq := constants.A8; sq <= constants.H1; sq++ {
// 		// init bishop & rook masks
// 		BishopMasks[sq] = GenerateBishopAttacks(sq)
// 		RookMasks[sq] = GenerateRookAttacks(sq)
//
// 		// init current mask
// 		var attackMask bitboard.Bitboard
//
// 		if piece == constants.Bishop {
// 			attackMask = GenerateBishopAttacks(sq)
// 		} else {
// 			attackMask = GenerateRookAttacks(sq)
// 		}
//
// 		// count attack mask bits
// 		bitCount := attackMask.Count()
//
// 		// occupancy variations count
// 		occupancyVariations := 1 << bitCount
//
// 		// loop over occupancy variations
// 		for count := 0; count < occupancyVariations; count++ {
// 			if piece == constants.Bishop {
// 				occupancy := SetOccupancy(count, bitCount, attackMask)
//
// 				magicIndex := occupancy * BishopMagicNumbers[sq] >> (64 - BishopRelevantBits[sq])
// 				BishopAttacks[sq][magicIndex] = GenerateBishopAttacksOnTheFly(sq, occupancy)
// 			} else {
//
// 				occupancy := SetOccupancy(count, bitCount, attackMask)
//
// 				magicIndex := occupancy * RookMagicNumbers[sq] >> (64 - RookRelevantBits[sq])
// 				RookAttacks[sq][magicIndex] = GenerateRookAttacksOnTheFly(sq, occupancy)
// 			}
// 		}
// 	}
// }
