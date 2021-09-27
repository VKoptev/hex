// Package hex contains geometry of hexagonal world.
package hex

import (
	"fmt"
)

const (
	// DirectionCount is const to be used for arrays and calculations.
	DirectionCount = 6
)

// LineSign is alias for float64 to determine it from other real values.
type LineSign float64

// Line signs.
// LSPlus means counterclockwise path with ray connecting start hex with end one being traced on edge.
// LSMinus means clockwise path with ray connecting start hex with end one being traced on edge.
// LSZero means that there are no such ray or unknown direction.
const (
	LSPlus  LineSign = 1
	LSMinus LineSign = -1
	LSZero  LineSign = 0
)

// String implements fmt.Stringer interface.
func (ls LineSign) String() string {
	switch ls {
	case LSPlus:
		return "line sign plus"
	case LSMinus:
		return "line sign minus"
	case LSZero:
		return "line sign zero"
	}

	return fmt.Sprintf("%f", ls)
}

// Hex is coordinates at Hexagonal Grid.
type Hex struct {
	q, r, s int
}

// Directions.
var (
	ZE         = Hex{0, 0, 0}
	EE         = Hex{1, 0, -1}
	NE         = Hex{1, -1, 0}
	NW         = Hex{0, -1, 1}
	WW         = Hex{-1, 0, 1}
	SW         = Hex{-1, 1, 0}
	SE         = Hex{0, 1, -1}
	directions = [...]Hex{EE, NE, NW, WW, SW, SE}
)

// New returns Hex object. It's convenient to use object but neither pointer.
func New(q, r int) Hex {
	return Hex{q: q, r: r, s: -q - r}
}

// NewWithArray returns new Hex object represented as array.
func NewWithArray(a [2]int) Hex {
	return New(a[0], a[1])
}

// Q returns Hex.q coordinate.
func (h Hex) Q() int {
	return h.q
}

// R returns Hex.r coordinate.
func (h Hex) R() int {
	return h.r
}

// S returns Hex.s coordinate.
func (h Hex) S() int {
	return h.s
}

// Equal returns true if Hexes h and t are equal.
func (h Hex) Equal(t Hex) bool {
	return h.q == t.q && h.r == t.r && h.s == t.s
}

// Add returns Hex that equals sum of h and t.
func (h Hex) Add(t Hex) Hex {
	return New(h.q+t.q, h.r+t.r)
}

// Sub returns Hex that equals difference of h and t.
func (h Hex) Sub(t Hex) Hex {
	return New(h.q-t.q, h.r-t.r)
}

// Mul returns Hex that equals h multiplied by k.
func (h Hex) Mul(k int) Hex {
	return New(h.q*k, h.r*k)
}

// Len returns radius-vector of Hex.
func (h Hex) Len() int {
	// nolint:gomnd
	return (absInt(h.q) + absInt(h.r) + absInt(h.s)) / 2
}

// Distance returns distance between h and t.
func (h Hex) Distance(t Hex) int {
	return h.Sub(t).Len()
}

// Neighbor returns neighbor Hex of h to direction d.
func (h Hex) Neighbor(d int) Hex {
	return h.Add(Direction(d))
}

// Direction returns index of direction.
func (h Hex) Direction(t Hex, sign LineSign) int {
	l := h.Line(t, sign)

	if len(l) <= 1 {
		// Top direction
		return DirectionCount
	}

	r := l[1].Sub(l[0])

	for i := 0; i < DirectionCount; i++ {
		if directions[i] == r {
			return i
		}
	}

	// Error
	return -1
}

// RingAtDistance returns ring.
func (h Hex) RingAtDistance(d int, res []Hex) {
	if d <= 0 {
		return
	}

	t := h.Add(SW.Mul(d))

	for i := 0; i < DirectionCount; i++ {
		for j := 0; j < d; j++ {
			res[i*d+j] = t
			t = t.Neighbor(i)
		}
	}
}

// NeighborsAtDistance returns all neighbors Hex of h on distance d.
func (h Hex) NeighborsAtDistance(d int, res []Hex) {
	for i := 1; i <= d; i++ {
		j := AreaAtDistance(i-1) - 1
		k := AreaAtDistance(i) - 1
		h.RingAtDistance(i, res[j:k])
	}
}

// Line returns Hexes from h to t.
func (h Hex) Line(t Hex, sign LineSign) (res []Hex) {
	n := h.Distance(t)
	step := 1 / float64(maxInt(n, 1))

	for i := 0; i <= n; i++ {
		res = append(res, hexLerp(h, t, step*float64(i), float64(sign)).round())
	}

	return res
}

// String returns string representation of Hex. It implements Stringer interface.
func (h Hex) String() string {
	return fmt.Sprintf("%v", h.Array())
}

// Title returns string representation of Hex.
func (h Hex) Title() string {
	switch h {
	case ZE:
		return "zero"
	case EE:
		return "east"
	case NE:
		return "north-east"
	case NW:
		return "north-west"
	case WW:
		return "west"
	case SW:
		return "south-west"
	case SE:
		return "south-east"
	}

	return fmt.Sprintf("%d,%d", h.q, h.r)
}

// Array returns base type representation of Hex.
func (h Hex) Array() [2]int {
	return [2]int{h.q, h.r}
}

// CompareSlices returns true if slices a and b have equal length and elements at same positions.
func CompareSlices(a, b []Hex) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// InverseDirection returns inverter direction for specified.
func InverseDirection(d int) int {
	return (((DirectionCount + (d % DirectionCount)) % DirectionCount) + DirectionCount/2) % DirectionCount
}

// NormalizeDirection returns normalized value direction for specified.
func NormalizeDirection(d int) int {
	return (DirectionCount + (d % DirectionCount)) % DirectionCount
}

// PrettyDirection returns pretty representation of direction.
func PrettyDirection(d int) string {
	if d == DirectionCount {
		return ZE.Title()
	}

	return Direction(d).Title()
}

// RingLenAtDistance returns length of ring at specified distance.
func RingLenAtDistance(r int) int {
	// nolint:gomnd
	return 6 * r
}

// AreaAtDistance returns amount of all hexes limited by ring at specified distance
// including center hex and ring.
func AreaAtDistance(r int) int {
	return 1 + 3*r*(r+1)
}

// Direction returns hex by direction.
func Direction(d int) Hex {
	return directions[NormalizeDirection(d)]
}
