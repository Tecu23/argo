// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
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

	for piece := WP; piece <= BK; piece++ {
		bitboard = b.Bitboards[piece]

		// generate pawns and castling moves depending of size
		if b.Side == color.WHITE {
			if piece == WP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq + N

					// quiet pawn moves
					if !(targetSq < A1) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= A7 && sourceSq <= H7 {
							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									WQ,
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
									WR,
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
									WB,
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
									WN,
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
							if (sourceSq >= A2 && sourceSq <= H2) && !b.Occupancies[color.BOTH].Test(targetSq+N) {
								movelist.AddMove(move.EncodeMove(sourceSq, targetSq+N, piece, 0, 0, 1, 0, 0))
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.WHITE][sourceSq] & b.Occupancies[color.BLACK]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if sourceSq >= A7 && sourceSq <= H7 {

							movelist.AddMove(
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									WQ,
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
									WR,
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
									WB,
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
									WN,
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
			if piece == WK {
				// King side castling is available
				if uint(b.Castlings)&ShortW != 0 {
					// make sure square between king and king's rook are empty
					if !b.Occupancies[color.BOTH].Test(F1) &&
						!b.Occupancies[color.BOTH].Test(G1) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(E1, color.BLACK) &&
							!b.IsSquareAttacked(F1, color.BLACK) {
							movelist.AddMove(
								move.EncodeMove(E1, G1, piece, 0, 0, 0, 0, 1),
							)
						}
					}
				}

				// Queen side castling is available
				if uint(b.Castlings)&LongW != 0 {
					// make sure square between king and queens's rook are empty
					if !b.Occupancies[color.BOTH].Test(D1) &&
						!b.Occupancies[color.BOTH].Test(C1) &&
						!b.Occupancies[color.BOTH].Test(B1) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(E1, color.BLACK) &&
							!b.IsSquareAttacked(D1, color.BLACK) {
							movelist.AddMove(
								move.EncodeMove(E1, C1, piece, 0, 0, 0, 0, 1),
							)
						}
					}
				}
			}

		} else {
			if piece == BP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq + S

					// quiet pawn moves
					if !(targetSq < 0 || targetSq > H8) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= A2 && sourceSq <= H2 {
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BQ, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BR, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BB, 0, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BN, 0, 0, 0, 0))
						} else {
							// one square ahead move
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, 0, 0, 0, 0, 0))

							// two square ahead move
							if (sourceSq >= A7 && sourceSq <= H7) && !b.Occupancies[color.BOTH].Test(targetSq+S) {
								movelist.AddMove(move.EncodeMove(sourceSq, targetSq+S, piece, 0, 0, 1, 0, 0))
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.BLACK][sourceSq] & b.Occupancies[color.WHITE]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if sourceSq >= A2 && sourceSq <= H2 {

							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BQ, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BR, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BB, 1, 0, 0, 0))
							movelist.AddMove(move.EncodeMove(sourceSq, targetSq, piece, BN, 1, 0, 0, 0))
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
			if piece == BK {
				// King side castling is available
				if uint(b.Castlings)&ShortB != 0 {
					// make sure square between king and king's rook are empty
					if !b.Occupancies[color.BOTH].Test(F8) && !b.Occupancies[color.BOTH].Test(G8) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(E8, color.WHITE) && !b.IsSquareAttacked(F8, color.WHITE) {
							movelist.AddMove(move.EncodeMove(E8, G8, piece, 0, 0, 0, 0, 1))
						}
					}
				}

				// Queen side castling is available
				if uint(b.Castlings)&LongB != 0 {
					// make sure square between king and queens's rook are empty
					if !b.Occupancies[color.BOTH].Test(D8) && !b.Occupancies[color.BOTH].Test(C8) &&
						!b.Occupancies[color.BOTH].Test(B8) {
						// make sure king and the f1 square are not under attack
						if !b.IsSquareAttacked(E8, color.WHITE) && !b.IsSquareAttacked(D8, color.WHITE) {
							movelist.AddMove(move.EncodeMove(E8, C8, piece, 0, 0, 0, 0, 1))
						}
					}
				}
			}
		}

		// generate knight moves
		if b.Side == color.WHITE {
			if piece == WN {
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
			if piece == BN {
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
			if piece == WB {
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
			if piece == BB {
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
			if piece == WR {
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
			if piece == BR {
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
			if piece == WQ {
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
			if piece == BQ {
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
			if piece == WK {
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
			if piece == BK {
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
