package tables

// KnightTable represents piece-square values for knights
var KnightTable = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50, // 8th rank
	-40, -20, 0, 0, 0, 0, -20, -40, // 7th rank
	-30, 0, 10, 15, 15, 10, 0, -30, // 6th rank
	-30, 5, 15, 20, 20, 15, 5, -30, // 5th rank
	-30, 0, 15, 20, 20, 15, 0, -30, // 4th rank
	-30, 5, 10, 15, 15, 10, 5, -30, // 3rd rank
	-40, -20, 0, 5, 5, 0, -20, -40, // 2nd rank
	-50, -40, -30, -30, -30, -30, -40, -50, // 1st rank
}

// Knights:
// - Strong bonus for central squares
// - Heavy penalties for edge squares
// - Encourages development to active squares
// - Slight bonus for forward positioning
