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

func (a *Vectorf) Normalize() *Vectorf {
	mag := a.Magnitude()
	return &Vectorf{a[0] / mag, a[1] / mag, a[2] / mag}
}

func (a Vectorf) String() string {
	return fmt.Sprintf("Vector(%f,%f,%f)", a[0], a[1], a[2])
}

type Vectori [3]int16

type Rect struct {
	x, z, sizex, sizez float64
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
	self[XAXIS] += dx
	self[YAXIS] += dy
	self[ZAXIS] += dz
}

func (self *Vectori) String() string {
	return fmt.Sprintf("[x:%d, y:%d, z:%d]", self[XAXIS], self[YAXIS], self[ZAXIS])
}

func (r1 Rect) Intersects(r2 Rect) bool {
	if r2.x >= r1.x && r2.x <= r1.x+r1.sizex && r2.z >= r1.z && r2.z <= r1.z+r1.sizez {
		return true
	}
	if r2.x+r2.sizex >= r1.x && r2.x+r2.sizex <= r1.x+r1.sizex && r2.z >= r1.z && r2.z <= r1.z+r1.sizez {
		return true
	}
	if r2.x+r2.sizex >= r1.x && r2.x+r2.sizex <= r1.x+r1.sizex && r2.z+r2.sizez >= r1.z && r2.z+r2.sizez <= r1.z+r1.sizez {
		return true
	}
	if r2.x >= r1.x && r2.x <= r1.x+r1.sizex && r2.z+r2.sizez >= r1.z && r2.z+r2.sizez <= r1.z+r1.sizez {
		return true
	}
	return false

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
type Matrix4 [16]float64

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

type MatrixError string

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

func (a *Matrix4) Rotatex(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply( &Matrix4{ 1, 0, 0, 0,   0, c, s, 0,   0, -s, c, 0,   0, 0, 0, 1} )
}

func (a *Matrix4) Rotatey(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply( &Matrix4{ c, 0, -s, 0,   0, 1, 0, 0,   s, 0, c, 0,   0, 0, 0, 1} )
}

func (a *Matrix4) Rotatez(angle float64) *Matrix4 {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	return a.Multiply( &Matrix4{ c, s, 0, 0,   -s, c, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1} )
}

func (a *Matrix4) Scale(scale float64) *Matrix4 {
	return a.Multiply( &Matrix4{ scale, 0, 0, 0,   0, scale, 0, 0,   0, 0, scale, 0,   0, 0, 0, 1} )
}

func (a *Matrix4) Transform(v *Vectorf, w float64) *Vectorf {
	return &Vectorf{ 

		v[0]*a[0] + v[1]*a[4] + v[2]*a[8]  + w*a[12],
		v[0]*a[1] + v[1]*a[5] + v[2]*a[9]  + w*a[13],
		v[0]*a[2] + v[1]*a[6] + v[2]*a[10] + w*a[14]}
}