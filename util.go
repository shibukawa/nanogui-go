package nanogui

import "math"

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

func clampF(a, min, max float32) float32 {
	if a > max {
		return max
	} else if a < min {
		return min
	}
	return a
}

func toB(condition bool, a, b uint8) uint8 {
	if condition {
		return a
	}
	return b
}
func toI(condition bool, a, b int) int {
	if condition {
		return a
	}
	return b
}

func toF(condition bool, a, b float32) float32 {
	if condition {
		return a
	}
	return b
}

func maxFs(v float32, values ...float32) float32 {
	max := v
	for _, value := range values {
		if max < value {
			max = value
		}
	}
	return max
}

func minFs(v float32, values ...float32) float32 {
	min := v
	for _, value := range values {
		if min > value {
			min = value
		}
	}
	return min
}

func sqrtF(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

func sinCosF(a float32) (float32, float32) {
	sin, cos := math.Sincos(float64(a))
	return float32(sin), float32(cos)
}

func absF(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func floorF(a float32) float32 {
	return float32(math.Floor(float64(a)))
}
