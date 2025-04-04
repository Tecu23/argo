// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

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

// source, target, piece, promoted, captured, captureFlag, doublePush, enpassant, castlingFlag, castlingType int,
func (b *Board) GenerateMoves() []move.Move {
	result := make([]move.Move, 0, 10)
	sourceSq, targetSq := 0, 0
	bitboard, attcks := bitboard.Bitboard(0), bitboard.Bitboard(0)

	for piece := WP; piece <= BK; piece++ {
		bitboard = b.Bitboards[piece]

		// generate pawns and castling moves depending of size
		if b.SideToMove == color.WHITE {
			if piece == WP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq - 8

					// quiet pawn moves
					if !(targetSq < A8) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= A7 && sourceSq <= H7 {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.QueenPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.RookPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.BishopPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.KnightPromotionCapture,
									0,
								),
							)
						} else {

							// one square ahead move
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)

							// two square ahead move
							if (sourceSq >= A2 && sourceSq <= H2) && !b.Occupancies[color.BOTH].Test(targetSq-8) {
								result = append(
									result,
									move.EncodeMove(sourceSq, targetSq-8, piece, move.DoublePawnPush, 0),
								)
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.WHITE][sourceSq] & b.Occupancies[color.BLACK]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()
						captured := b.GetPieceAt(targetSq)

						if sourceSq >= A7 && sourceSq <= H7 {

							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.QueenPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.RookPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.BishopPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.KnightPromotionCapture,
									captured,
								),
							)
						} else {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}

					// generate EnPassant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.WHITE][sourceSq] & (1 << b.EnPassant)

						// check enpassant capture
						if enpassantAttacks != 0 {
							// init enpassant capture target square
							targetSq := enpassantAttacks.FirstOne()
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.EnPassant, BP),
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
							result = append(
								result,
								move.EncodeMove(E1, G1, piece, move.KingCastle, 0),
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
							result = append(
								result,
								move.EncodeMove(E1, C1, piece, move.QueenCastle, 0),
							)
						}
					}
				}
			}

		} else {
			if piece == BP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()
					targetSq = sourceSq + 8

					// quiet pawn moves
					if !(targetSq < 0 || targetSq > H1) && !b.Occupancies[color.BOTH].Test(targetSq) {
						// pawn promotion
						if sourceSq >= A2 && sourceSq <= H2 {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.QueenPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.RookPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.BishopPromotion, 0),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.KnightPromotion, 0),
							)
						} else {
							// one square ahead move
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)

							// two square ahead move
							if (sourceSq >= A7 && sourceSq <= H7) && !b.Occupancies[color.BOTH].Test(targetSq+8) {
								result = append(
									result,
									move.EncodeMove(sourceSq, targetSq+8, piece, move.DoublePawnPush, 0),
								)
							}
						}
					}

					// init pawn attacks bb
					attcks = attacks.PawnAttacks[color.BLACK][sourceSq] & b.Occupancies[color.WHITE]

					// generate pawn captures
					for attcks != 0 {
						targetSq = attcks.FirstOne()

						captured := b.GetPieceAt(targetSq)
						if sourceSq >= A2 && sourceSq <= H2 {

							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.QueenPromotionCapture, captured),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.RookPromotionCapture, captured),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.BishopPromotionCapture, captured),
							)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.KnightPromotionCapture, captured),
							)
						} else {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}

					// generate EnPassant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.BLACK][sourceSq] & (1 << b.EnPassant)

						// check enpassant capture
						if enpassantAttacks != 0 {
							// init enpassant capture target square
							targetSq := enpassantAttacks.FirstOne()
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.EnPassant, WP),
							)
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
							result = append(
								result,
								move.EncodeMove(E8, G8, piece, move.KingCastle, 0),
							)
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
							result = append(
								result,
								move.EncodeMove(E8, C8, piece, move.QueenCastle, 0),
							)
						}
					}
				}
			}
		}

		// generate knight moves
		if b.SideToMove == color.WHITE {
			if piece == WN {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KnightAttacks[sourceSq] & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}
				}
			}
		}

		// generate bihop moves
		if b.SideToMove == color.WHITE {
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}
				}
			}
		}

		// generate rook moves
		if b.SideToMove == color.WHITE {
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}
				}
			}
		}

		// generate queen moves
		if b.SideToMove == color.WHITE {
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}
				}
			}
		}

		// generate king moves
		if b.SideToMove == color.WHITE {
			if piece == WK {
				for bitboard != 0 {

					sourceSq = bitboard.FirstOne()

					// init piece attacks
					attcks = attacks.KingAttacks[sourceSq] & (^b.Occupancies[color.WHITE])

					for attcks != 0 {
						targetSq = attcks.FirstOne()

						if !b.Occupancies[color.BLACK].Test(targetSq) {
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
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
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Quiet, 0),
							)
						} else {
							captured := b.GetPieceAt(targetSq)
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
					}
				}
			}
		}
	}

	return result
}

// GenerateCaptures generates all capturing moves for the current position,
// including pawn promotions which are considered tactical moves.
func (b *Board) GenerateCaptures() []move.Move {
	var result []move.Move
	sourceSq, targetSq := 0, 0
	bitboard, attcks := bitboard.Bitboard(0), bitboard.Bitboard(0)

	for piece := WP; piece <= BK; piece++ {
		bitboard = b.Bitboards[piece]

		// Generate pawn captures
		if b.SideToMove == color.WHITE {
			if piece == WP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()

					// Pawn captures including promotions
					attcks = attacks.PawnAttacks[color.WHITE][sourceSq] & b.Occupancies[color.BLACK]

					for attcks != 0 {
						targetSq = attcks.FirstOne()
						captured := b.GetPieceAt(targetSq)
						// Handle promotion captures
						if sourceSq >= A7 && sourceSq <= H7 {
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.QueenPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.RookPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.BishopPromotionCapture,
									captured,
								),
							)
							result = append(
								result,
								move.EncodeMove(
									sourceSq,
									targetSq,
									piece,
									move.KnightPromotionCapture,
									captured,
								),
							)
						} else {
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
						attcks &= attcks - 1 // clear LSB
					}

					// En passant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.WHITE][sourceSq] & (1 << b.EnPassant)
						if enpassantAttacks != 0 {
							targetSq := enpassantAttacks.FirstOne()
							result = append(
								result,
								move.EncodeMove(sourceSq, targetSq, piece, move.EnPassant, BP),
							)
						}
					}
					bitboard &= bitboard - 1 // clear LSB
				}
			}
		} else { // Black pawns
			if piece == BP {
				for bitboard != 0 {
					sourceSq = bitboard.FirstOne()

					// Pawn captures including promotions
					attcks = attacks.PawnAttacks[color.BLACK][sourceSq] & b.Occupancies[color.WHITE]

					for attcks != 0 {
						targetSq = attcks.FirstOne()
						captured := b.GetPieceAt(targetSq)
						// Handle promotion captures
						if sourceSq >= A2 && sourceSq <= H2 {
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.QueenPromotionCapture, captured),
							)
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.RookPromotionCapture, captured),
							)
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.BishopPromotionCapture, captured),
							)
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.KnightPromotionCapture, captured),
							)
						} else {
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
							)
						}
						attcks &= attcks - 1 // clear LSB
					}

					// En passant captures
					if b.EnPassant != -1 {
						enpassantAttacks := attacks.PawnAttacks[color.BLACK][sourceSq] & (1 << b.EnPassant)
						if enpassantAttacks != 0 {
							targetSq := enpassantAttacks.FirstOne()
							result = append(result,
								move.EncodeMove(sourceSq, targetSq, piece, move.EnPassant, WP),
							)
						}
					}
					bitboard &= bitboard - 1 // clear LSB
				}
			}
		}

		// Generate knight captures
		if (b.SideToMove == color.WHITE && piece == WN) ||
			(b.SideToMove == color.BLACK && piece == BN) {
			for bitboard != 0 {
				sourceSq = bitboard.FirstOne()
				// Get attacks that hit enemy pieces
				if b.SideToMove == color.WHITE {
					attcks = attacks.KnightAttacks[sourceSq] & b.Occupancies[color.BLACK]
				} else {
					attcks = attacks.KnightAttacks[sourceSq] & b.Occupancies[color.WHITE]
				}

				// Generate captures
				for attcks != 0 {
					targetSq = attcks.FirstOne()
					captured := b.GetPieceAt(targetSq)
					result = append(
						result,
						move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
					)
					attcks &= attcks - 1
				}
				bitboard &= bitboard - 1
			}
		}

		// Generate bishop captures
		if (b.SideToMove == color.WHITE && piece == WB) ||
			(b.SideToMove == color.BLACK && piece == BB) {
			for bitboard != 0 {
				sourceSq = bitboard.FirstOne()
				if b.SideToMove == color.WHITE {
					attcks = attacks.GetBishopAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & b.Occupancies[color.BLACK]
				} else {
					attcks = attacks.GetBishopAttacks(sourceSq, b.Occupancies[color.BOTH]) & b.Occupancies[color.WHITE]
				}

				for attcks != 0 {
					targetSq = attcks.FirstOne()
					captured := b.GetPieceAt(targetSq)
					result = append(
						result,
						move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
					)
					attcks &= attcks - 1
				}
				bitboard &= bitboard - 1
			}
		}

		// Generate rook captures
		if (b.SideToMove == color.WHITE && piece == WR) ||
			(b.SideToMove == color.BLACK && piece == BR) {
			for bitboard != 0 {
				sourceSq = bitboard.FirstOne()
				if b.SideToMove == color.WHITE {
					attcks = attacks.GetRookAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & b.Occupancies[color.BLACK]
				} else {
					attcks = attacks.GetRookAttacks(sourceSq, b.Occupancies[color.BOTH]) & b.Occupancies[color.WHITE]
				}

				for attcks != 0 {
					targetSq = attcks.FirstOne()
					captured := b.GetPieceAt(targetSq)
					result = append(
						result,
						move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
					)
					attcks &= attcks - 1
				}
				bitboard &= bitboard - 1
			}
		}

		// Generate queen captures
		if (b.SideToMove == color.WHITE && piece == WQ) ||
			(b.SideToMove == color.BLACK && piece == BQ) {
			for bitboard != 0 {
				sourceSq = bitboard.FirstOne()
				if b.SideToMove == color.WHITE {
					attcks = attacks.GetQueenAttacks(
						sourceSq,
						b.Occupancies[color.BOTH],
					) & b.Occupancies[color.BLACK]
				} else {
					attcks = attacks.GetQueenAttacks(sourceSq, b.Occupancies[color.BOTH]) & b.Occupancies[color.WHITE]
				}

				for attcks != 0 {
					targetSq = attcks.FirstOne()
					captured := b.GetPieceAt(targetSq)
					result = append(
						result,
						move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
					)
					attcks &= attcks - 1
				}
				bitboard &= bitboard - 1
			}
		}

		// Generate king captures
		if (b.SideToMove == color.WHITE && piece == WK) ||
			(b.SideToMove == color.BLACK && piece == BK) {
			for bitboard != 0 {
				sourceSq = bitboard.FirstOne()
				if b.SideToMove == color.WHITE {
					attcks = attacks.KingAttacks[sourceSq] & b.Occupancies[color.BLACK]
				} else {
					attcks = attacks.KingAttacks[sourceSq] & b.Occupancies[color.WHITE]
				}

				for attcks != 0 {
					targetSq = attcks.FirstOne()
					captured := b.GetPieceAt(targetSq)
					result = append(
						result,
						move.EncodeMove(sourceSq, targetSq, piece, move.Capture, captured),
					)
					attcks &= attcks - 1
				}
				bitboard &= bitboard - 1
			}
		}
	}

	return result
}
