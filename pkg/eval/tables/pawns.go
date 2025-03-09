package tables

// PawnTable represents piece-square values for pawns from White's perspective.
// Positive values are good for White, negative values are good for Black.
var PawnTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0, // 8th rank
	50, 50, 50, 50, 50, 50, 50, 50, // 7th rank (promotion potential)
	10, 10, 20, 30, 30, 20, 10, 10, // 6th rank
	5, 5, 10, 25, 25, 10, 5, 5, // 5th rank
	0, 0, 0, 20, 20, 0, 0, 0, // 4th rank
	5, -5, -10, 0, 0, -10, -5, 5, // 3rd rank
	5, 10, 10, -20, -20, 10, 10, 5, // 2nd rank
	0, 0, 0, 0, 0, 0, 0, 0, // 1st rank
}

// Pawns:
// - Encouraged to advance (higher values in ranks 4-7)
// - Penalized for staying on starting squares
// - Center control bonus (e4, d4, e5, d5)
// - Slightly favors pawn chain formations
