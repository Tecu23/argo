package tables

var PawnOpeningTable = [64]int{
	// 0, 0, 0, 0, 0, 0, 0, 0, // 8th rank
	// 50, 50, 50, 50, 50, 50, 50, 50, // 7th rank (promotion potential)
	// 10, 10, 20, 30, 30, 20, 10, 10, // 6th rank
	// 5, 5, 10, 25, 25, 10, 5, 5, // 5th rank
	// 0, 0, 0, 20, 20, 0, 0, 0, // 4th rank
	// 5, -5, -10, 0, 0, -10, -5, 5, // 3rd rank
	// 5, 10, 10, -20, -20, 10, 10, 5, // 2nd rank
	// 0, 0, 0, 0, 0, 0, 0, 0, // 1st rank
	0, 0, 0, 0, 0, 0, 0, 0,
	45, 52, 42, 43, 28, 34, 19, 9,
	-14, -3, 7, 14, 35, 50, 15, -6,
	-27, -6, -8, 13, 16, 4, -3, -25,
	-32, -28, -7, 5, 7, -1, -15, -30,
	-29, -25, -12, -12, -1, -5, 6, -17,
	-34, -23, -27, -18, -14, 10, 13, -22,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PawnEndgameTable = [64]int{
	// 0, 0, 0, 0, 0, 0, 0, 0,
	// 80, 80, 80, 80, 80, 80, 80, 80,
	// 50, 50, 50, 50, 50, 50, 50, 50,
	// 30, 30, 30, 30, 30, 30, 30, 30,
	// 20, 20, 20, 20, 20, 20, 20, 20,
	// 10, 10, 10, 10, 10, 10, 10, 10,
	// 5, 5, 5, 5, 5, 5, 5, 5,
	// 0, 0, 0, 0, 0, 0, 0, 0,
	// // EG Pawn PST
	0, 0, 0, 0, 0, 0, 0, 0,
	77, 74, 63, 53, 59, 60, 72, 77,
	17, 11, 11, 11, 11, -6, 14, 8,
	-3, -14, -18, -31, -29, -25, -20, -18,
	-12, -14, -24, -31, -29, -28, -27, -28,
	-22, -20, -25, -20, -21, -24, -34, -34,
	-16, -22, -11, -19, -13, -23, -32, -34,
	0, 0, 0, 0, 0, 0, 0, 0,
}
