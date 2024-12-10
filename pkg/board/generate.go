// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
)

// TODO: Split this method into captures/quiet moves for all pieces to debug more easily
// TODO: Fix Issue where a black pawn could capture a black king ???????

// GenerateMoves generates all pseudo-legal moves for the current board position and
// stores them in movelist. It handles pawns, knights, bishops, rooks, queens, kings,
// and also pawn promotions, en passant, castling, etc.
//
// This is a large function with complex logic, so comments are only at a high level.
//
// Steps (high-level):
// 1. Identify side to move.
// 2. For each piece type of that color, find all pieces on the board.
// 3. Generate moves according to the piece movement rules and occupancy of the board.
// 4. For pawns: handle one and two-step advances, captures, en passant, and promotions.
// 5. For kings: handle castling if available and king moves.
// 6. For knights, bishops, rooks, queens: use precomputed attack tables for sliding pieces and direct attacks for knights/king.
// 7. Add each valid generated move to movelist. This function doesn't check for check legality - that is done in MakeMove.
func (b *Board) GenerateMoves(movelist *move.Movelist) {
	sourceSq, targetSq := 0, 0
	bitboard, attcks := bitboard.Bitboard(0), bitboard.Bitboard(0)

	for piece := constants.WP; piece <= constants.BK; piece++ {
		bitboard = b.Bitboards[piece]

		// generate pawns and castling moves depending of size
		if b.Side == color.WHITE {
			if piece == constants.WP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq + constants.N

					// quiet pawn moves
					if !(targetSq < constants.A1) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= constants.A7 && sourceSq <= constants.H7 {
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WQ,
									0,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WR,
									0,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WB,
									0,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WN,
									0,
									0,
									0,
									0,
								),
							)
						} else {

							// one square ahead move
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))

							// two square ahead move
							if (sourceSq >= constants.A2 && sourceSq <= constants.H2) && !b.Occupancies[color.BOTH].Test(targetSq+constants.N) {
								movelist.AddMove(move.EncodeMove(sourceSq, targetSq+constants.N, piece, 0, 0, 1, 0, 0))
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.WHITE][sourceSq] & b.Occupancies[color.BLACK]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if sourceSq >= constants.A7 && sourceSq <= constants.H7 {

							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WQ,
									1,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WR,
									1,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WB,
									1,
									0,
									0,
									0,
								),
							)
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									constants.WN,
									1,
									0,
									0,
									0,
								),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}

					// generate EnPassant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.WHITE][sourceSq] & (1 << b.EnPassant)

						// check enpassant capture
						if enpassantAttacks != 0 {
							// init enpassant capture target square
							targetEnpassant := enpassantAttacks.FirstOne()
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetEnpassant, piece, 0, 1, 0, 1, 0),
							)
						}
					}
				}
			}
			// Castlings moves
			if piece == constants.WK {
				// King side castling is available
				if uint(b.Castlings)&constants.ShortW != 0 {
					// make sure square between king and king's rook are empty
					if !b.Occupancies[color.BOTH].Test(constants.F1) &&
						!b.Occupancies[color.BOTH].Test(constants.G1) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(constants.E1, color.BLACK) &&
							!b.IsSquareAttacked(constants.F1, color.BLACK) {
							movelist.AddMove(
								move.EncodeMove(constants.E1, constants.G1, piece, 0, 0, 0, 0, 1),
							)
						}
					}
				}

				// Queen side castling is available
				if uint(b.Castlings)&constants.LongW != 0 {
					// make sure square between king and queens's rook are empty
					if !b.Occupancies[color.BOTH].Test(constants.D1) &&
						!b.Occupancies[color.BOTH].Test(constants.C1) &&
						!b.Occupancies[color.BOTH].Test(constants.B1) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(constants.E1, color.BLACK) &&
							!b.IsSquareAttacked(constants.D1, color.BLACK) {
							movelist.AddMove(
								move.EncodeMove(constants.E1, constants.C1, piece, 0, 0, 0, 0, 1),
							)
						}
					}
				}
			}

		} else {
			if piece == constants.BP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq + constants.S

					// quiet pawn moves
					if !(targetSq < 0 || targetSq > constants.H8) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= constants.A2 && sourceSq <= constants.H2 {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BQ, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BR, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BB, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BN, 0, 0, 0, 0))
						} else {
							// one square ahead move
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))

							// two square ahead move
							if (sourceSq >= constants.A7 && sourceSq <= constants.H7) && !b.Occupancies[color.BOTH].Test(targetSq+constants.S) {
								movelist.AddMove(move.EncodeMove(sourceSq, targetSq+constants.S, piece, 0, 0, 1, 0, 0))
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.BLACK][sourceSq] & b.Occupancies[color.WHITE]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if sourceSq >= constants.A2 && sourceSq <= constants.H2 {

							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BQ, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BR, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BB, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, constants.BN, 1, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}

					// generate EnPassant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.BLACK][sourceSq] & (1 << b.EnPassant)

						// check enpassant capture
						if enpassantAttacks != 0 {
							// init enpassant capture target square
							targetEnpassant := enpassantAttacks.FirstOne()
							movelist.AddMove(move.EncodeMove(sourceSq, targetEnpassant, piece, 0, 1, 0, 1, 0))
						}
					}
				}
			}

			// Castlings moves
			if piece == constants.BK {
				// King side castling is available
				if uint(b.Castlings)&constants.ShortB != 0 {
					// make sure square between king and king's rook are empty
					if !b.Occupancies[color.BOTH].Test(constants.F8) && !b.Occupancies[color.BOTH].Test(constants.G8) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(constants.E8, color.WHITE) && !b.IsSquareAttacked(constants.F8, color.WHITE) {
							movelist.AddMove(move.EncodeMove(constants.E8, constants.G8, piece, 0, 0, 0, 0, 1))
						}
					}
				}

				// Queen side castling is available
				if uint(b.Castlings)&constants.LongB != 0 {
					// make sure square between king and queens's rook are empty
					if !b.Occupancies[color.BOTH].Test(constants.D8) && !b.Occupancies[color.BOTH].Test(constants.C8) &&
						!b.Occupancies[color.BOTH].Test(constants.B8) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(constants.E8, color.WHITE) && !b.IsSquareAttacked(constants.D8, color.WHITE) {
							movelist.AddMove(move.EncodeMove(constants.E8, constants.C8, piece, 0, 0, 0, 0, 1))
						}
					}
				}
			}
		}

		// generate knight moves
		if b.Side == color.WHITE {
			if piece == constants.WN {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KnightAttacks[sourceSq] & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		} else {
			if piece == constants.BN {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KnightAttacks[sourceSq] & (^b.Occupancies[color.BLACK])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.WHITE].Test(targetSq) {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		}

		// generate bihop moves
		if b.Side == color.WHITE {
			if piece == constants.WB {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetBishopAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		} else {
			if piece == constants.BB {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetBishopAttacks(sourceSq, b.Occupancies[color.BOTH]) & (^b.Occupancies[color.BLACK])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.WHITE].Test(targetSq) {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		}

		// generate rook moves
		if b.Side == color.WHITE {
			if piece == constants.WR {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetRookAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		} else {
			if piece == constants.BR {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetRookAttacks(sourceSq, b.Occupancies[color.BOTH]) & (^b.Occupancies[color.BLACK])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.WHITE].Test(targetSq) {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		}

		// generate queen moves
		if b.Side == color.WHITE {
			if piece == constants.WQ {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetQueenAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		} else {
			if piece == constants.BQ {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.GetQueenAttacks(sourceSq, b.Occupancies[color.BOTH]) & (^b.Occupancies[color.BLACK])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.WHITE].Test(targetSq) {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		}

		// generate king moves
		if b.Side == color.WHITE {
			if piece == constants.WK {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KingAttacks[sourceSq] & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							movelist.AddMove(
								move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0),
							)
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		} else {
			if piece == constants.BK {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KingAttacks[sourceSq] & (^b.Occupancies[color.BLACK])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.WHITE].Test(targetSq) {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))
						} else {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 1, 0, 0, 0))
						}
					}
				}
			}
		}
	}
}
