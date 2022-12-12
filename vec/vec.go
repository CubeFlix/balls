// vec/vec.go
// Package vec provides vectors.

package vec

// Vector 2D.
type Vec2 struct {
	X float64
	Y float64
}

// Add vector 2Ds.
func AddVec2(vecs ...Vec2) Vec2 {
	x := float64(0)
	y := float64(0)
	for i := range vecs {
		x += vecs[i].X
		y += vecs[i].Y
	}
	return NewVec2(x, y)
}

// New vector 2D.
func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

// Add two vector 2Ds.
func (v Vec2) Add(x Vec2) Vec2 {
	return NewVec2(v.X+x.X, v.Y+x.Y)
}

// Subtract two vector 2Ds.
func (v Vec2) Sub(x Vec2) Vec2 {
	return NewVec2(v.X-x.X, v.Y-x.Y)
}

// Multiply by a scalar.
func (v Vec2) MulScalar(x float64) Vec2 {
	return NewVec2(v.X*x, v.Y*x)
}
