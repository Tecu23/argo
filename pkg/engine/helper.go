package engine

func winIn(height int) int {
	return MateScore - height
}

func lossIn(height int) int {
	return -MateScore + height
}
