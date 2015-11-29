package nanogui

func minF(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampI(a, min, max int) int {
	if a > max {
		return max
	} else if a < min {
		return min
	}
	return a
}

func toI(condition bool, a, b int) int {
	if condition {
		return a
	}
	return b
}
