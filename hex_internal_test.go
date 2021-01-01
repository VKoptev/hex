package hex

import (
	"fmt"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = New(i, -i)
	}
}

func BenchmarkHex_NeighborsAtDistance(b *testing.B) {
	r := make(map[int][]Hex)
	for _, d := range []int{1, 6, 14, 20, 100} {
		r[d] = make([]Hex, AreaAtDistance(d)-1)
	}
	for _, d := range []int{1, 6, 14, 20, 100} {
		d := d
		b.Run(fmt.Sprintf("distance-%d", d), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ZE.NeighborsAtDistance(d, r[d])
			}
		})
	}
}

func TestHex_Add(t *testing.T) {
	t.Parallel()
	tt := [][3][3]int{
		{{0, 0, 0}, {1, 1, -2}, {1, 1, -2}},
		{{10, 10, -20}, {1, 1, -2}, {11, 11, -22}},
	}
	for _, test := range tt {
		h1 := New(test[0][0], test[0][1])
		h2 := New(test[1][0], test[1][1])
		r := h1.Add(h2)

		if r.Q() != test[2][0] || r.R() != test[2][1] || r.S() != test[2][2] {
			t.Errorf("error in test: %+v", test)
		}
	}
}

func TestHex_RingAtDistance(t *testing.T) {
	t.Parallel()
	tt := []struct {
		d int
		e []Hex
	}{
		{d: 1, e: []Hex{New(-1, 1), New(0, 1), New(1, 0), New(1, -1), New(0, -1), New(-1, 0)}},
		{
			d: 10,
			e: []Hex{
				New(-10, 10), New(-9, 10), New(-8, 10), New(-7, 10), New(-6, 10),
				New(-5, 10), New(-4, 10), New(-3, 10), New(-2, 10), New(-1, 10),
				New(0, 10), New(1, 9), New(2, 8), New(3, 7), New(4, 6),
				New(5, 5), New(6, 4), New(7, 3), New(8, 2), New(9, 1),
				New(10, 0), New(10, -1), New(10, -2), New(10, -3), New(10, -4),
				New(10, -5), New(10, -6), New(10, -7), New(10, -8), New(10, -9),
				New(10, -10), New(9, -10), New(8, -10), New(7, -10), New(6, -10),
				New(5, -10), New(4, -10), New(3, -10), New(2, -10), New(1, -10),
				New(0, -10), New(-1, -9), New(-2, -8), New(-3, -7), New(-4, -6),
				New(-5, -5), New(-6, -4), New(-7, -3), New(-8, -2), New(-9, -1),
				New(-10, 0), New(-10, 1), New(-10, 2), New(-10, 3), New(-10, 4),
				New(-10, 5), New(-10, 6), New(-10, 7), New(-10, 8), New(-10, 9),
			},
		},
	}

	for _, test := range tt {
		r := make([]Hex, RingLenAtDistance(test.d))
		ZE.RingAtDistance(test.d, r)
		if !CompareSlices(r, test.e) {
			t.Errorf("d=%d r=%v e=%v", test.d, r, test.e)
		}
	}
}

func TestHex_NeighborsAtDistance(t *testing.T) {
	t.Parallel()
	tt := []struct {
		d int
		e []Hex
	}{
		{d: 1, e: []Hex{New(-1, 1), New(0, 1), New(1, 0), New(1, -1), New(0, -1), New(-1, 0)}},
		{d: 2, e: []Hex{
			New(-1, 1), New(0, 1), New(1, 0), New(1, -1), New(0, -1), New(-1, 0),
			New(-2, 2), New(-1, 2), New(0, 2), New(1, 1), New(2, 0), New(2, -1),
			New(2, -2), New(1, -2), New(0, -2), New(-1, -1), New(-2, 0), New(-2, 1),
		}},
	}
	for _, test := range tt {
		r := make([]Hex, AreaAtDistance(test.d)-1)
		ZE.NeighborsAtDistance(test.d, r)
		if !CompareSlices(r, test.e) {
			t.Errorf("d=%d r=%v e=%v", test.d, r, test.e)
		}
	}
}

func TestHex_Line(t *testing.T) {
	t.Parallel()
	tt := []struct {
		h, t Hex
		s    LingSign
		e    []Hex
	}{
		{
			h: ZE,
			t: EE.Mul(3),
			s: LSPlus,
			e: []Hex{ZE, EE, EE.Mul(2), EE.Mul(3)},
		},
		{
			h: SW.Add(SE).Add(SW),
			t: ZE,
			s: LSPlus,
			e: []Hex{SW.Add(SE).Add(SW), SW.Add(SE), SW, ZE},
		},
		{
			h: ZE,
			t: SW.Add(SE),
			s: LSPlus,
			e: []Hex{ZE, SW, SW.Add(SE)},
		},
		{
			h: ZE,
			t: SE.Add(SW), // eq SW.Add(SE)
			s: LSMinus,
			e: []Hex{ZE, SE, SE.Add(SW)},
		},
		{
			h: ZE,
			t: NE.Add(NW),
			s: LSPlus,
			e: []Hex{ZE, NE, NE.Add(NW)},
		},
		{
			h: ZE,
			t: NW.Add(NE), // eq NE.Add(NW)
			s: LSMinus,
			e: []Hex{ZE, NW, NW.Add(NE)},
		},
		{
			h: SE.Add(SW),
			t: ZE,
			s: LSPlus,
			e: []Hex{SE.Add(SW), SE, ZE},
		},
		{
			h: SW.Add(SE), // eq SE.Add(SW)
			t: ZE,
			s: LSMinus,
			e: []Hex{SW.Add(SE), SW, ZE},
		},
	}

	for _, test := range tt {
		l := test.h.Line(test.t, test.s)
		if !CompareSlices(l, test.e) {
			t.Errorf("h=%s t=%s l=%v expected=%v", test.h.String(), test.t.String(), l, test.e)
		}
	}
}

func TestHex_Direction(t *testing.T) {
	t.Parallel()
	t.Run("mul-10", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 6; i++ {
			h := directions[i].Mul(10)
			d := ZE.Direction(h, LSZero)
			if d != i {
				t.Errorf("d=%v expected=%v", d, i)
			}

			inv := InverseDirection(i)
			d = h.Direction(ZE, LSZero)
			if d != inv {
				t.Errorf("d=%v expected=%v", d, inv)
			}
		}
	})

	t.Run("plus-minus", func(t *testing.T) {
		t.Parallel()
		ss := []LingSign{LSPlus, LSMinus}
		hh := make([]Hex, RingLenAtDistance(2))
		ZE.RingAtDistance(2, hh)

		for i := range hh {
			if i%2 == 0 {
				continue
			}
			for j := range ss {
				e := NormalizeDirection(4 + (i / 2) + j)
				d := ZE.Direction(hh[i], ss[j])
				if d != e {
					t.Errorf("h=%v ls=%v d=%v e=%v", hh[i], ss[j], d, e)
				}
			}
		}
	})
}
