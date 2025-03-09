package tables

// BishopTable represents piece-square values for bishops
var BishopTable = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20, // 8th rank
	-10, 0, 0, 0, 0, 0, 0, -10, // 7th rank
	-10, 0, 5, 10, 10, 5, 0, -10, // 6th rank
	-10, 5, 5, 10, 10, 5, 5, -10, // 5th rank
	-10, 0, 10, 10, 10, 10, 0, -10, // 4th rank
	-10, 10, 10, 10, 10, 10, 10, -10, // 3rd rank
	-10, 5, 0, 0, 0, 0, 5, -10, // 2nd rank
	-20, -10, -10, -10, -10, -10, -10, -20, // 1st rank
}

// Bishops:
// - Bonus for diagonal activity
// - Penalty for corner squares
// - Encourages control of central diagonals
// - Small bonus for development
