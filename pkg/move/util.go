package move

func createMask(n uint) uint32 {
	return (1 << n) - 1
}

func SameMove(m1, m2 Move) bool {
	return (uint32(m1^m2) & createMask(24)) == 0
}
