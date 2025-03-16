// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// ThreatsEvaluation evaluates the threats in a certain position
func (e *Evaluator) ThreatsEvaluation(b *board.Board) (mg, eg int) {
	mirror := b.Mirror()
	mg, eg = 0, 0

	hangingBonus := hanging(b)
	mg += 69 * hangingBonus
	eg += 36 * hangingBonus

	if kingThreat(b) {
		mg += 24
		eg += 80
	}

	pawnPushBonus := pawnPushThreat(b)
	mg += 48 * pawnPushBonus
	eg += 39 * pawnPushBonus

	threatSafeBonus := threatSafePawn(b)
	mg += 173 * threatSafeBonus
	eg += 94 * threatSafeBonus

	sliderOnQueenBonus := sliderOnQueen(b, mirror)
	mg += 60 * sliderOnQueenBonus
	eg += 18 * sliderOnQueenBonus

	knightOnQueenBonus := knightOnQueen(b, mirror)
	mg += 16 * knightOnQueenBonus
	eg += 11 * knightOnQueenBonus

	restrictedBonus := restricted(b)
	mg += 7 * restrictedBonus
	eg += 7 * restrictedBonus

	mg += 14 * weakQueenProtection(b)

	for sq := A8; sq <= H1; sq++ {
		minorThreatIdx := minorThreat(b, sq)
		rookThreatIdx := rookThreat(b, sq)

		mg += []int{0, 5, 57, 77, 88, 79, 0}[minorThreatIdx]
		eg += []int{0, 32, 41, 56, 119, 161, 0}[minorThreatIdx]

		mg += []int{0, 3, 37, 42, 0, 58, 0}[rookThreatIdx]
		eg += []int{0, 46, 68, 60, 38, 41, 0}[rookThreatIdx]
	}

	return mg, eg
}

// hanging returns weak enemies not defended by opponent or non-pawn weak
// enemies attacked twice
func hanging(b *board.Board) int {
	weakEnemiesCount := 0
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		if weakEnemies(b, sq) == 0 {
			continue
		}

		if !b.Bitboards[BP].Test(sq) && attack(b, sq) > 1 {
			weakEnemiesCount++
			continue
		}

		mirror := b.Mirror()
		file := sq % 8
		rank := sq / 8
		if attack(mirror, (7-rank)*8+file) == 0 {
			weakEnemiesCount++
			continue
		}
	}
	return weakEnemiesCount
}

// weakEnemies returns enemies not defended by a pawn and under our attack
func weakEnemies(b *board.Board, sq int) int {
	rank := sq / 8
	file := sq % 8

	if b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 && rank > 0 {
		return 0
	}

	if b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 && rank > 0 {
		return 0
	}

	if attack(b, sq) == 0 {
		return 0
	}

	mirror := b.Mirror()
	if attack(b, sq) <= 1 && attack(mirror, (7-rank)*8+file) > 1 {
		return 0
	}

	return 1
}

// attack counts the number of attacks on square by all pieces. For bishop and rook
// x-ray attacks are included. For pawns pins or en-passant are ignored.
func attack(b *board.Board, sq int) int {
	score := 0
	score += PawnAttack(b, sq)
	score += KingAttack(b, sq)
	score += KnightAttack(b, sq, -1)
	score += BishopXrayAttack(b, sq, -1)
	score += RookXrayAttack(b, sq, -1)
	score += QueenAttack(b, sq, -1)

	return score
}

// kingThreat returns if the king is in threat
func kingThreat(b *board.Board) bool {
	threats := 0
	blackBB := b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		if weakEnemies(b, sq) == 0 {
			continue
		}

		if KingAttack(b, sq) == 0 {
			continue
		}

		threats++

	}

	if threats > 0 {
		return true
	}

	return false
}

// pawnPushThreat returns the number of pawns that can be safely pushed
// and attack and enemy piece
func pawnPushThreat(b *board.Board) int {
	score := 0
	for sq := A8; sq <= H1; sq++ {
		rank := sq / 8
		file := sq % 8

		if !b.Occupancies[color.BLACK].Test(sq) {
			continue
		}

		for ix := -1; ix <= 1; ix += 2 {
			if b.Bitboards[WP].Test((rank+2)*8+file+ix) &&
				file+ix >= 0 && file+ix <= 7 && rank+2 <= 7 &&
				!b.Occupancies[color.BOTH].Test((rank+1)*8+file+ix) &&
				(!b.Bitboards[BP].Test(rank*8+file+ix-1) && file+ix-1 >= 0 && file+ix-1 <= 7) &&
				(!b.Bitboards[BP].Test(rank*8+file+ix+1) && file+ix+1 >= 0 && file+ix+1 <= 7) &&
				(attack(b, (rank+1)*8+file+ix) > 0 || attack(b.Mirror(), (6-rank)*8+file+ix) == 0) {
				score++
			}

			if rank == 3 &&
				(b.Bitboards[WP].Test((rank+3)*8+file+ix) && file+ix >= 0 && file+ix <= 7 && rank+3 <= 7) &&
				(!b.Occupancies[color.BOTH].Test((rank+2)*8+file+ix) && rank+2 <= 7) &&
				(!b.Occupancies[color.BOTH].Test((rank+1)*8+file+ix) && rank+1 <= 7) &&
				(!b.Bitboards[BP].Test(rank*8+file+ix-1) && file+ix-1 >= 0 && file+ix-1 <= 7) &&
				(!b.Bitboards[BP].Test(rank*8+file+ix+1) && file+ix+1 >= 0 && file+ix+1 <= 7) &&
				(attack(b, (rank+1)*8+file+ix) > 0 || attack(b.Mirror(), (6-rank)*8+file+ix) == 0) {
				score++
			}
		}

	}
	return score
}

// threatSafePawn returns the non-pawn enemies attacked by a safe pawn
func threatSafePawn(b *board.Board) int {
	score := 0
	blackBB := b.Bitboards[BP] ^ b.Occupancies[color.BLACK]
	for blackBB != 0 {
		sq := blackBB.FirstOne()

		rank := sq / 8
		file := sq % 8

		if PawnAttack(b, sq) == 0 {
			continue
		}

		if (safePawn(b, (rank+1)*8+file-1) && file > 0 && rank < 7) ||
			(safePawn(b, (rank+1)*8+file+1) && file < 7 && rank < 7) {
			score++
		}
	}
	return score
}

// safePawn returns whether or not our pawn is not attacked or is defended
func safePawn(b *board.Board, sq int) bool {
	rank := sq / 8
	file := sq % 8

	if b.Bitboards[WP].Test(sq) {
		return false
	}

	if attack(b, sq) > 0 {
		return true
	}

	if attack(b.Mirror(), (7-rank)*8+file) == 0 {
		return true
	}

	return false
}

// sliderOnQueen adds a bonus for safe slider attack threats on opponent queen
func sliderOnQueen(b *board.Board, mirror *board.Board) int {
	score := 0

	if queenCount(mirror) != 1 {
		return 0
	}

	rank, file := 0, 0

	for sq := A8; sq <= H1; sq++ {
		rank = sq / 8
		file = sq % 8

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			continue
		}

		if b.Bitboards[BP].Test((rank-1)*8+file+1) && rank > 0 && file < 7 {
			continue
		}

		if attack(b, sq) <= 1 {
			continue
		}

		if !mobilityArea(b, mirror, sq) {
			continue
		}

		diagonal := QueenAttackDiagonal(mirror, (7-rank)*8+file, -1)
		v := 1
		if queenCount(b) == 0 {
			v = 2
		}

		if diagonal != 0 && BishopXrayAttack(b, sq, -1) != 0 {
			score += v
		}

		if diagonal == 0 &&
			RookXrayAttack(b, sq, -1) != 0 &&
			QueenAttack(mirror, (7-rank)*8+file, -1) != 0 {
			score += v
		}
	}

	return score
}

// queenCount returns the number of white queens
func queenCount(b *board.Board) int {
	return b.Bitboards[WQ].Count()
}

// knightOnQueen returns a bonus for safe knight attack threaths on
// opponent queen
func knightOnQueen(b *board.Board, mirror *board.Board) int {
	blackQueen := b.Bitboards[BQ]
	blackQueenSq := blackQueen.FirstOne()

	blackQueenRank := blackQueenSq / 8
	blackQueenFile := blackQueenSq % 8

	if queenCount(mirror) != 1 {
		return 0
	}

	rank, file := 0, 0

	score := 0

	for sq := A8; sq <= H1; sq++ {

		rank = sq / 8
		file = sq % 8

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			continue
		}

		if b.Bitboards[BP].Test((rank-1)*8+file-1) && rank > 0 && file > 0 {
			continue
		}

		if attack(b, sq) <= 1 && attack(mirror, (7-rank)*8+file) > 1 {
			continue
		}

		if !mobilityArea(b, mirror, sq) {
			continue
		}

		if KnightAttack(b, sq, -1) == 0 {
			continue
		}

		v := 1
		if queenCount(b) == 0 {
			v = 2
		}

		if abs(blackQueenFile-file) == 2 && abs(blackQueenRank-rank) == 1 {
			score += v
		}

		if abs(blackQueenFile-file) == 1 && abs(blackQueenRank-rank) == 2 {
			score += v
		}
	}

	return score
}

// restricted returns a bonus for restricing their pieces moves
func restricted(b *board.Board) int {
	restricted := 0

	for sq := A8; sq <= H1; sq++ {
		rank := sq / 8
		file := sq % 8

		mirror := b.Mirror()

		if attack(b, sq) == 0 {
			continue
		}

		if attack(mirror, (7-rank)*8+file) == 0 {
			continue
		}

		if PawnAttack(mirror, (7-rank)*8+file) > 0 {
			continue
		}

		if attack(mirror, (7-rank)*8+file) > 1 && attack(b, sq) == 1 {
			continue
		}

		restricted++
	}

	return restricted
}

// weakQueenProtection adds an additional bonus if weak piece is only
// protected by a queen
func weakQueenProtection(b *board.Board) int {
	weakPieces := 0

	for sq := A8; sq <= H1; sq++ {
		if !b.Occupancies[color.BLACK].Test(sq) {
			continue
		}

		if weakEnemies(b, sq) == 0 {
			continue
		}

		rank := sq / 8
		file := sq % 8

		if QueenAttack(b.Mirror(), (7-rank)*8+file, -1) == 0 {
			continue
		}

		weakPieces++
	}
	return weakPieces
}

// minorThreat returns the threat type for knight and bishop attacking pieces
func minorThreat(b *board.Board, sq int) int {
	if !b.Occupancies[color.BLACK].Test(sq) {
		return 0
	}

	if KnightAttack(b, sq, -1) == 0 &&
		BishopXrayAttack(b, sq, -1) == 0 {
		return 0
	}

	rank := sq / 8
	file := sq % 8

	if (b.Bitboards[BP].Test(sq) ||
		!((b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 && rank > 0) ||
			(b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 && rank > 0) ||
			(attack(b, sq) <= 1 && attack(b.Mirror(), (7-rank)*8+file) > 1))) &&
		weakEnemies(b, sq) == 0 {

		return 0
	}

	if b.Bitboards[BP].Test(sq) {
		return 1
	}

	if b.Bitboards[BN].Test(sq) {
		return 2
	}

	if b.Bitboards[BB].Test(sq) {
		return 3
	}

	if b.Bitboards[BR].Test(sq) {
		return 4
	}

	if b.Bitboards[BQ].Test(sq) {
		return 5
	}

	return 6
}

// rookThreat return the threat type for attacked by rook pieces
func rookThreat(b *board.Board, sq int) int {
	if !b.Occupancies[color.BLACK].Test(sq) {
		return 0
	}

	if weakEnemies(b, sq) == 0 {
		return 0
	}

	if RookXrayAttack(b, sq, -1) == 0 {
		return 0
	}

	if b.Bitboards[BP].Test(sq) {
		return 1
	}

	if b.Bitboards[BN].Test(sq) {
		return 2
	}

	if b.Bitboards[BB].Test(sq) {
		return 3
	}

	if b.Bitboards[BR].Test(sq) {
		return 4
	}

	if b.Bitboards[BQ].Test(sq) {
		return 5
	}

	return 6
}
