/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"math"
)

type Vectorf [3]float64
type Vectori [3]uint16

type Rect struct {
	x, y, sizex, sizey float64
}

type Ray struct {
	origin *Vectorf
	dir    *Vectorf
}

type Box struct {
	min *Vectorf
	max *Vectorf
}

type Matrix4 [16]float64

type MatrixError string

func NewMatrix(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p float64) *Matrix4 {
	return &Matrix4{a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p}
}

func NewIdentity() *Matrix4 {
	return &Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func NewTranslation(x, y, z float64) *Matrix4 {
	return &Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1}
}

func (a Vectorf) Add(b Vectorf) Vectorf {
	return Vectorf{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func (a Vectorf) Minus(b Vectorf) Vectorf {
	return Vectorf{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func (a Vectorf) Dot(b Vectorf) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func (a Vectorf) Cross(b Vectorf) Vectorf {
	return Vectorf{a[1]*b[2] - b[1]*a[2], a[2]*b[0] - b[2]*a[0], a[0]*b[1] - b[0]*a[1]}
}

func (a Vectorf) Scale(f float64) Vectorf {
	return Vectorf{a[0] * f, a[1] * f, a[2] * f}
}

func (a Vectorf) Magnitude() float64 {
	return math.Sqrt(math.Pow(a[0], 2) + math.Pow(a[1], 2) + math.Pow(a[2], 2))
}

func (a Vectorf) Angle(b Vectorf) float64 {
	amag := a.Magnitude()
	bmag := b.Magnitude()
	if amag == 0 || bmag == 0 {
		return 0
	}
	return math.Acos(a.Dot(b) / (amag * bmag))
}

// Assumes a and b are normalized
func (a Vectorf) AngleNormalized(b Vectorf) float64 {
	return math.Acos(a.Dot(b))
}

func (a Vectorf) Normalize() Vectorf {
	mag := a.Magnitude()
	return Vectorf{a[0] / mag, a[1] / mag, a[2] / mag}
}

func (a Vectorf) String() string {
	return fmt.Sprintf("Vector(%f,%f,%f)", a[0], a[1], a[2])
}

func (wp Vectori) North() Vectori { return Vectori{wp[XAXIS], wp[YAXIS], wp[ZAXIS] - 1} }
func (wp Vectori) South() Vectori { return Vectori{wp[XAXIS], wp[YAXIS], wp[ZAXIS] + 1} }
func (wp Vectori) East() Vectori  { return Vectori{wp[XAXIS] + 1, wp[YAXIS], wp[ZAXIS]} }
func (wp Vectori) West() Vectori  { return Vectori{wp[XAXIS] - 1, wp[YAXIS], wp[ZAXIS]} }
func (wp Vectori) Up() Vectori    { return Vectori{wp[XAXIS], wp[YAXIS] + 1, wp[ZAXIS]} }
func (wp Vectori) Down() Vectori  { return Vectori{wp[XAXIS], wp[YAXIS] - 1, wp[ZAXIS]} }

func (wp Vectori) HRect() Rect {
	return Rect{float64(wp[XAXIS]) - 0.5, float64(wp[ZAXIS]) - 0.5, 1, 1}
}

func (self *Vectori) Adjust(dx int16, dy int16, dz int16) {
	self[XAXIS] = uint16(int32(self[XAXIS]) + int32(dx))
	self[YAXIS] = uint16(int32(self[YAXIS]) + int32(dy))
	self[ZAXIS] = uint16(int32(self[ZAXIS]) + int32(dz))
}

func (self *Vectori) Translate(dx int16, dy int16, dz int16) Vectori {
	return Vectori{uint16(int32(self[XAXIS]) + int32(dx)),
		uint16(int32(self[YAXIS]) + int32(dy)),
		uint16(int32(self[ZAXIS]) + int32(dz))}
}

func (self *Vectori) String() string {
	return fmt.Sprintf("[x:%d, y:%d, z:%d]", self[XAXIS], self[YAXIS], self[ZAXIS])
}

func (self *Vectori) Vectorf() Vectorf {
	return Vectorf{float64(self[XAXIS]), float64(self[YAXIS]), float64(self[ZAXIS])}
}

func (self *Vectori) Equals(b *Vectori) bool {
	return self[XAXIS] == b[XAXIS] && self[YAXIS] == b[YAXIS] && self[ZAXIS] == b[ZAXIS]
}

func (r1 Rect) Intersects(r2 Rect) bool {
	if r2.x >= r1.x && r2.x <= r1.x+r1.sizex && r2.y >= r1.y && r2.y <= r1.y+r1.sizey {
		return true
	}
	if r2.x+r2.sizex >= r1.x && r2.x+r2.sizex <= r1.x+r1.sizex && r2.y >= r1.y && r2.y <= r1.y+r1.sizey {
		return true
	}
	if r2.x+r2.sizex >= r1.x && r2.x+r2.sizex <= r1.x+r1.sizex && r2.y+r2.sizey >= r1.y && r2.y+r2.sizey <= r1.y+r1.sizey {
		return true
	}
	if r2.x >= r1.x && r2.x <= r1.x+r1.sizex && r2.y+r2.sizey >= r1.y && r2.y+r2.sizey <= r1.y+r1.sizey {
		return true
	}
	return false

}

func (self *Rect) Contains(x, y float64) bool {
	return x >= self.x && x < self.x+self.sizex && y >= self.y && y < self.y+self.sizey
}

// Stored column first like opengl
//
//      | 0  4  8  12 |
//      |             |
//      | 1  5  9  13 |
//  M = |             |
//      | 2  6  10 14 |
//      |             |
//      | 3  7  11 15 |
//

func (a *Matrix4) String() string {
	return fmt.Sprintf("[%18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f %18.13f]", a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14], a[15])
}

func (a *Matrix4) Float32() *[16]float32 {
	var ret [16]float32
	for i := 0; i < 16; i++ {
		ret[i] = float32(a[i])
	}
	return &ret
}

func (a *Matrix4) Equals(b *Matrix4, precision int) bool {
	return (Round(a[0], precision) == Round(b[0], precision) &&
		Round(a[1], precision) == Round(b[1], precision) &&
		Round(a[2], precision) == Round(b[2], precision) &&
		Round(a[3], precision) == Round(b[3], precision) &&
		Round(a[4], precision) == Round(b[4], precision) &&
		Round(a[5], precision) == Round(b[5], precision) &&
		Round(a[6], precision) == Round(b[6], precision) &&
		Round(a[7], precision) == Round(b[7], precision) &&
		Round(a[8], precision) == Round(b[8], precision) &&
		Round(a[9], precision) == Round(b[9], precision) &&
		Round(a[10], precision) == Round(b[10], precision) &&
		Round(a[11], precision) == Round(b[11], precision) &&
		Round(a[12], precision) == Round(b[12], precision) &&
		Round(a[13], precision) == Round(b[13], precision) &&
		Round(a[14], precision) == Round(b[14], precision))

}

func (a *Matrix4) Transpose() *Matrix4 {
	return &Matrix4{
		a[0], a[4], a[8], a[12],
		a[1], a[5], a[9], a[13],
		a[2], a[6], a[10], a[14],
		a[3], a[7], a[11], a[15]}
}

func (a *Matrix4) Add(b *Matrix4) *Matrix4 {
	return &Matrix4{
		a[0] + b[0], a[1] + b[1], a[2] + b[2], a[3] + b[3],
		a[4] + b[4], a[5] + b[5], a[6] + b[6], a[7] + b[7],
		a[8] + b[8], a[9] + b[9], a[10] + b[10], a[11] + b[11],
		a[12] + b[12], a[13] + b[13], a[14] + b[14], a[15] + b[15]}
}

func (a *Matrix4) Subtract(b *Matrix4) *Matrix4 {
	return &Matrix4{
		a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3],
		a[4] - b[4], a[5] - b[5], a[6] - b[6], a[7] - b[7],
		a[8] - b[8], a[9] - b[9], a[10] - b[10], a[11] - b[11],
		a[12] - b[12], a[13] - b[13], a[14] - b[14], a[15] - b[15]}
}

func (a *Matrix4) Multiply(b *Matrix4) *Matrix4 {
	var m Matrix4

	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			m[i*4+j] = b[i*4]*a[j] +
				b[i*4+1]*a[4+j] +
				b[i*4+2]*a[8+j] +
				b[i*4+3]*a[12+j]
		}
	}

	return &m
}

//
//       | A B C |
//   M = | D E F |
//       | G H I |
//
//   det M = A * (EI - HF) - B * (DI - GF) + C * (DH - GE)
//
func det3(a *[9]float64) float64 {
	return a[0]*(a[4]*a[8]-a[5]*a[7]) - a[3]*(a[1]*a[8]-a[2]*a[7]) + a[6]*(a[1]*a[5]-a[2]*a[4])
}

func (a *Matrix4) Det() float64 {
	var det float64

	if a[0] != 0 {
		det += a[0] * det3(&[9]float64{a[5], a[6], a[7], a[9], a[10], a[11], a[13], a[14], a[15]})
	}
	if a[1] != 0 {
		det -= a[1] * det3(&[9]float64{a[4], a[6], a[7], a[8], a[10], a[11], a[12], a[14], a[15]})
	}
	if a[2] != 0 {
		det += a[2] * det3(&[9]float64{a[4], a[5], a[7], a[8], a[9], a[11], a[12], a[13], a[15]})
	}
	if a[3] != 0 {
		det -= a[3] * det3(&[9]float64{a[4], a[5], a[6], a[8], a[9], a[10], a[12], a[13], a[14]})
	}

	return det
}

func (e MatrixError) Error() string {
	return string(e)
}

func (a *Matrix4) Inverse() (*Matrix4, error) {
	var b Matrix4
	det := a.Det()
	if det == 0 {
		return a, MatrixError("Matrix is singular and has no inverse")
	}

	b[0] = (a[5]*a[10]*a[15] + a[9]*a[14]*a[7] + a[13]*a[6]*a[11] - a[5]*a[14]*a[11] - a[9]*a[6]*a[15] - a[13]*a[10]*a[7]) / det
	b[4] = (a[4]*a[14]*a[11] + a[8]*a[6]*a[15] + a[12]*a[10]*a[7] - a[4]*a[10]*a[15] - a[8]*a[14]*a[7] - a[12]*a[6]*a[11]) / det

	b[8] = (a[4]*a[9]*a[15] + a[8]*a[13]*a[7] + a[12]*a[5]*a[11] - a[4]*a[13]*a[11] - a[8]*a[5]*a[15] - a[12]*a[9]*a[7]) / det
	b[12] = (a[4]*a[13]*a[10] + a[8]*a[5]*a[14] + a[12]*a[9]*a[6] - a[4]*a[9]*a[14] - a[8]*a[13]*a[6] - a[12]*a[5]*a[10]) / det

	b[1] = (a[1]*a[14]*a[11] + a[9]*a[2]*a[15] + a[13]*a[10]*a[3] - a[1]*a[10]*a[15] - a[9]*a[14]*a[3] - a[13]*a[2]*a[11]) / det
	b[5] = (a[0]*a[10]*a[15] + a[8]*a[14]*a[3] + a[12]*a[2]*a[11] - a[0]*a[14]*a[11] - a[8]*a[2]*a[15] - a[12]*a[10]*a[3]) / det
	b[9] = (a[0]*a[13]*a[11] + a[8]*a[1]*a[15] + a[12]*a[9]*a[3] - a[0]*a[9]*a[15] - a[8]*a[13]*a[3] - a[12]*a[1]*a[11]) / det
	b[13] = (a[0]*a[9]*a[14] + a[8]*a[13]*a[2] + a[12]*a[1]*a[10] - a[0]*a[13]*a[10] - a[8]*a[1]*a[14] - a[12]*a[9]*a[2]) / det

	b[2] = (a[1]*a[6]*a[15] + a[5]*a[14]*a[3] + a[13]*a[2]*a[7] - a[1]*a[14]*a[7] - a[5]*a[2]*a[15] - a[13]*a[6]*a[3]) / det
	b[6] = (a[0]*a[14]*a[7] + a[4]*a[2]*a[15] + a[12]*a[6]*a[3] - a[0]*a[6]*a[15] - a[4]*a[14]*a[3] - a[12]*a[2]*a[7]) / det
	b[10] = (a[0]*a[5]*a[15] + a[4]*a[13]*a[3] + a[12]*a[1]*a[7] - a[0]*a[13]*a[7] - a[4]*a[1]*a[15] - a[12]*a[5]*a[3]) / det
	b[14] = (a[0]*a[13]*a[6] + a[4]*a[1]*a[14] + a[12]*a[5]*a[2] - a[0]*a[5]*a[14] - a[4]*a[13]*a[2] - a[12]*a[1]*a[6]) / det

	b[3] = (a[1]*a[10]*a[7] + a[5]*a[2]*a[11] + a[9]*a[6]*a[3] - a[1]*a[6]*a[11] - a[5]*a[10]*a[3] - a[9]*a[2]*a[7]) / det
	b[7] = (a[0]*a[6]*a[11] + a[4]*a[10]*a[3] + a[8]*a[2]*a[7] - a[0]*a[10]*a[7] - a[4]*a[2]*a[11] - a[8]*a[6]*a[3]) / det
	b[11] = (a[0]*a[9]*a[7] + a[4]*a[1]*a[11] + a[8]*a[5]*a[3] - a[0]*a[5]*a[11] - a[4]*a[9]*a[3] - a[8]*a[1]*a[7]) / det
	b[15] = (a[0]*a[5]*a[10] + a[4]*a[9]*a[2] + a[8]*a[1]*a[6] - a[0]*a[9]*a[6] - a[4]*a[1]*a[10] - a[8]*a[5]*a[2]) / det

	return &b, nil
}

func (a *Matrix4) Translation(x, y, z float64) *Matrix4 {
	return a.Multiply(&Matrix4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, x, y, z, 1})
}

func (a *Matrix4) Rotatex(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply(&Matrix4{1, 0, 0, 0, 0, c, s, 0, 0, -s, c, 0, 0, 0, 0, 1})
}

func (a *Matrix4) Rotatey(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply(&Matrix4{c, 0, -s, 0, 0, 1, 0, 0, s, 0, c, 0, 0, 0, 0, 1})
}

func (a *Matrix4) Rotatez(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply(&Matrix4{c, s, 0, 0, -s, c, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
}

func (a *Matrix4) Scale(scale float64) *Matrix4 {
	return a.Multiply(&Matrix4{scale, 0, 0, 0, 0, scale, 0, 0, 0, 0, scale, 0, 0, 0, 0, 1})
}

func (a *Matrix4) Transform(v *Vectorf, w float64) *Vectorf {
	return &Vectorf{

		v[0]*a[0] + v[1]*a[4] + v[2]*a[8] + w*a[12],
		v[0]*a[1] + v[1]*a[5] + v[2]*a[9] + w*a[13],
		v[0]*a[2] + v[1]*a[6] + v[2]*a[10] + w*a[14]}
}

// Ported from http://tog.acm.org/resources/GraphicsGems/gems/RayBox.c
func (self *Ray) HitsBox(box *Box) (hit bool, face uint8) {
	const RIGHT = 0
	const LEFT = 1
	const MIDDLE = 2

	inside := true
	quadrant := [3]int{}
	var whichPlane int
	var candidatePlane, maxT, coord Vectorf

	// Find candidate planes
	for i := 0; i < 3; i++ {
		if self.origin[i] < box.min[i] {
			quadrant[i] = LEFT
			candidatePlane[i] = box.min[i]
			inside = false
		} else if self.origin[i] > box.max[i] {
			quadrant[i] = RIGHT
			candidatePlane[i] = box.max[i]
			inside = false
		} else {
			quadrant[i] = MIDDLE
		}
	}

	// Ray origin inside bounding box 
	if !inside {

		// Calculate T distances to candidate planes
		for i := 0; i < 3; i++ {
			if quadrant[i] != MIDDLE && self.dir[i] != 0 {
				maxT[i] = (candidatePlane[i] - self.origin[i]) / self.dir[i]
			} else {
				maxT[i] = -1.0
			}
		}

		// Get largest of the maxT's for final choice of intersection
		whichPlane = 0
		for i := 1; i < 3; i++ {
			if maxT[whichPlane] < maxT[i] {
				whichPlane = i
			}
		}

		// println("whichPlane", whichPlane)
		// println("maxT[whichPlane]", maxT[whichPlane])

		// Check final candidate actually inside box
		if maxT[whichPlane] < 0 {
			hit = false
			return
		}

		for i := 0; i < 3; i++ {
			if whichPlane != i {
				coord[i] = self.origin[i] + maxT[whichPlane]*self.dir[i]
				// println("coord[", i, "]=", coord[i])
				if coord[i] < box.min[i] || coord[i] > box.max[i] {
					hit = false
					return
				}
			} else {
				coord[i] = candidatePlane[i]
			}
		}
	}

	if whichPlane == 0 {
		if quadrant[whichPlane] == LEFT {
			face = WEST_FACE
		} else {
			face = EAST_FACE
		}

	} else if whichPlane == 1 {
		if quadrant[whichPlane] == LEFT {
			face = DOWN_FACE
		} else {
			face = UP_FACE
		}
	} else {
		if quadrant[whichPlane] == LEFT {
			face = NORTH_FACE
		} else {
			face = SOUTH_FACE
		}
	}

	hit = true
	return

}

// lineRectCollide( line, rect )
//
// Checks if an axis-aligned line and a bounding box overlap.
// line = { z, x1, x2 } or line = { x, z1, z2 }
// rect = { x, z, size }

func lineRectCollide(line Side, rect Rect) (ret bool) {
	if line.z != 0 {
		ret = rect.y > line.y-rect.sizey/2 && rect.y < line.y+rect.sizey/2 && rect.x > line.x1-rect.sizex/2 && rect.x < line.x2+rect.sizex/2
	} else {
		ret = rect.x > line.x-rect.sizex/2 && rect.x < line.x+rect.sizex/2 && rect.y > line.z1-rect.sizey/2 && rect.y < line.z2+rect.sizey/2
	}
	return
}

// rectRectCollide( r1, r2 )
//
// Checks if two rectangles (x1, y1, x2, y2) overlap.

func rectRectCollide(r1 Side, r2 Side) bool {
	if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 {
		return true
	}
	if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 {
		return true
	}
	if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
		return true
	}
	if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
		return true
	}
	return false
}
