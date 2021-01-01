package hex

import (
	"math"
)

var (
	epsilon = floatHex{q: 1e-6, r: 2e-6, s: -3e-6}
)

type floatHex struct {
	q, r, s float64
}

func newFloat(q, r float64) floatHex {
	return floatHex{q: q, r: r, s: -q - r}
}

func newFloatFromHex(h Hex) floatHex {
	return newFloat(float64(h.q), float64(h.r))
}

func (fh floatHex) round() Hex {
	q, r, s := roundInt(fh.q), roundInt(fh.r), roundInt(fh.s)
	dq, dr, ds := math.Abs(fh.q-float64(q)), math.Abs(fh.r-float64(r)), math.Abs(fh.s-float64(s))

	if dq > dr && dq > ds {
		q = -r - s
	} else if dr > ds {
		r = -q - s
	}

	return New(q, r)
}

func (fh floatHex) add(ft floatHex) floatHex {
	return newFloat(fh.q+ft.q, fh.r+ft.r)
}

func (fh floatHex) mul(i float64) floatHex {
	return newFloat(fh.q*i, fh.r*i)
}

func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

func hexLerp(a, b Hex, t, sign float64) floatHex {
	v := b.Sub(a)

	if v.Len() == 0 {
		return newFloat(0, 0)
	}

	// rotate epsilon vector without normalize to e.
	x, y := float64(v.q), float64(v.r)

	eps := newFloat(x*epsilon.q-y*epsilon.r, x*epsilon.r+y*epsilon.q)
	aa := newFloatFromHex(a).add(eps.mul(sign))
	bb := newFloatFromHex(b).add(eps.mul(sign))

	return newFloat(
		lerp(aa.q, bb.q, t),
		lerp(aa.r, bb.r, t),
	)
}
