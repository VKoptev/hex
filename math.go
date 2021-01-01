package hex

import "math"

func absInt(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func roundInt(x float64) int {
	return int(round(x))
}

func round(x float64) float64 {
	const (
		mask  = 0x7FF
		shift = 64 - 11 - 1
		bias  = 1023

		signMask = 1 << 63
		fracMask = (1 << shift) - 1
		halfMask = 1 << (shift - 1)
		one      = bias << shift
	)

	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask

	switch {
	case e < bias:
		// Round abs(x)<1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= one // +-1
		}
	case e < bias+shift:
		// Round any abs(x)>=1 containing a fractional component [0,1).
		e -= bias
		bits += halfMask >> e
		bits &^= fracMask >> e
	}

	return math.Float64frombits(bits)
}
