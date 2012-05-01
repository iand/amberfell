/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	// "fmt"
	// "math"
	"testing"
)

func TestMatrixAdd(t *testing.T) {
	a, b := NewIdentity(), NewIdentity()
	expected := &Matrix4{2.0, 0, 0, 0, 0, 2.0, 0, 0, 0, 0, 2.0, 0, 0, 0, 0, 2.0}
	actual := a.Add(b)
	if !actual.Equals(expected, 8) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMatrixMultiply(t *testing.T) {
	a := NewMatrix(320.0, 0.0, 0.0, 319.5, 0.0, 240.0, 0.0, 239.5, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0)
	b := NewMatrix(1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.12, -0.62, 0.0, 0.0, 0.2, 0.0)
	expected := NewMatrix(320.0, 0.0, 0.0, 319.5, 0.0, 240.0, 0.0, 239.5, 0.0, 0.0, 1.12, -0.62, 0.0, 0.0, 0.2, 0.0)
	actual := a.Multiply(b)
	if !actual.Equals(expected, 8) {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

}

func TestDet3(t *testing.T) {
	a := &[9]float64{0, 0, 319.5, 240, 0, 239.5, 0, 1, 0}
	expected := float64(76680.0000000000)
	actual := det3(a)
	if actual != expected {
		t.Errorf("Expected %18.13f but got %18.13f", expected, actual)
	}
	identity := &[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}
	expected_identity := float64(1)
	actual_identity := det3(identity)
	if actual_identity != expected_identity {
		t.Errorf("Expected %18.13f but got %18.13f", expected_identity, actual_identity)
	}
}

func TestMatrixDet(t *testing.T) {
	a := NewMatrix(320.0, 0.0, 0.0, 319.5, 0.0, 240.0, 0.0, 239.5, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0)
	b := NewMatrix(1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.12, -0.62, 0.0, 0.0, 0.2, 0.0)

	cases := map[*Matrix4]float64{
		a:             float64(76800),
		b:             float64(0.124),
		NewIdentity(): float64(1),
		a.Transpose(): a.Det(),
		a.Multiply(b): a.Det() * b.Det()}

	for m, expected_det := range cases {
		actual_det := m.Det()
		if Round(actual_det, 9) != Round(expected_det, 9) {
			t.Errorf("Expected determinant of %18.13f but got %18.13f for matrix %s", expected_det, actual_det, m)
		}

	}
}

func TestMatrixInverse(t *testing.T) {
	a := NewMatrix(320.0, 0.0, 0.0, 319.5, 0.0, 240.0, 0.0, 239.5, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0)
	b := NewMatrix(1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.12, -0.62, 0.0, 0.0, 0.2, 0.0)

	a_inv := NewMatrix(0.003125, 0, 0, -0.9984375, 0, 0.00416667, 0, -0.99791667, 0, 0, 1, 0, 0, 0, 0, 1)
	b_inv := NewMatrix(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 9.91270558e-17, 5, 0, 0, -1.61290323e+00, 9.03225806e+00)

	cases := map[*Matrix4]*Matrix4{
		a:     a_inv,
		b:     b_inv,
		a_inv: a}

	for m, expected_inverse := range cases {
		actual_inverse, _ := m.Inverse()
		if !expected_inverse.Equals(actual_inverse, 3) {
			t.Errorf("Expected inverse of %s but got %s for matrix %s", expected_inverse, actual_inverse, m)
		}

	}
}

func TestMatrixRotatex(t *testing.T) {
	cases := map[float64]*Matrix4{
		0:   NewMatrix(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1),
		90:  NewMatrix(1, 0, 0, 0, 0, 0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 1),
		180: NewMatrix(1, 0, 0, 0, 0, -1, 0, 0, 0, 0, -1, 0, 0, 0, 0, 1),
		270: NewMatrix(1, 0, 0, 0, 0, 0, -1, 0, 0, 1, 0, 0, 0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatex(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}

func TestMatrixRotatey(t *testing.T) {
	cases := map[float64]*Matrix4{
		0:   NewMatrix(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1),
		90:  NewMatrix(0, 0, -1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1),
		180: NewMatrix(-1, 0, 0, 0, 0, 1, 0, 0, 0, 0, -1, 0, 0, 0, 0, 1),
		270: NewMatrix(0, 0, 1, 0, 0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatey(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}

func TestMatrixRotatez(t *testing.T) {
	cases := map[float64]*Matrix4{
		0:   NewMatrix(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1),
		90:  NewMatrix(0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1),
		180: NewMatrix(-1, 0, 0, 0, 0, -1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1),
		270: NewMatrix(0, -1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatez(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}

func TestHitsBox(t *testing.T) {

	ray := Ray{&Vectorf{0, 0, 0}, &Vectorf{1, 0, 0}}
	box := Box{&Vectorf{2, -1, -1}, &Vectorf{3, 1, 1}}
	expected := true

	actual, _ := ray.HitsBox(&box)
	if expected != actual {
		if expected {
			t.Errorf("Ray %s did not hit box %s as expected", ray, box)
		} else {
			t.Errorf("Ray %s did not miss box %s as expected", ray, box)
		}
	}
}

func TestHitsBox2(t *testing.T) {
	cases := map[Ray]bool{
		Ray{&Vectorf{-8.667095, 13.666898, 16.215151}, &Vectorf{0.851651, -0.422618, -0.309976}}: true}

	box := Box{&Vectorf{8.5, 4.5, 9.5}, &Vectorf{9.5, 5.5, 10.5}}

	for ray, expected := range cases {
		actual, _ := ray.HitsBox(&box)
		if expected != actual {
			if expected {
				t.Errorf("Ray %s did not hit box %s as expected", ray, box)
			} else {
				t.Errorf("Ray %s did not miss box %s as expected", ray, box)
			}
		}
	}

	// Ray origin: {-8.667095, 13.666898, 16.215151
	// Ray norm: 0.851651, -0.422618, -0.309976
	// box: {%!s(*main.Vectorf=&[8.5 3.5 9.5]) %!s(*main.Vectorf=&[9.5 6.5 10.5])}
	// Ray origin: -8.557143, 13.787739, 16.352494
	// Ray norm: 0.851651, -0.422618, -0.309976
	// box: {%!s(*main.Vectorf=&[8.5 3.5 9.5]) %!s(*main.Vectorf=&[9.5 6.5 10.5])}
	// Ray origin: -8.474961, 14.014316, 16.269373
	// Ray norm: 0.851651, -0.422618, -0.309976
	// box: {%!s(*main.Vectorf=&[8.5 3.5 9.5]) %!s(*main.Vectorf=&[9.5 6.5 10.5])}
	// Ray origin: 24.588263, 13.833055, 0.603553
	// Ray norm: -0.774813, -0.422618, 0.470168
	// box: {%!s(*main.Vectorf=&[8.5 3.5 9.5]) %!s(*main.Vectorf=&[9.5 6.5 10.5])}

}
