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
		0:       NewMatrix(1, 0, 0, 0,   0, 1, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1),
		90:      NewMatrix(1, 0, 0, 0,   0, 0, 1, 0,   0, -1, 0, 0,   0, 0, 0, 1),
		180:     NewMatrix(1, 0, 0, 0,   0, -1, 0, 0,   0, 0, -1, 0,   0, 0, 0, 1),
		270:     NewMatrix(1, 0, 0, 0,   0, 0, -1, 0,   0, 1, 0, 0,   0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatex(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}

func TestMatrixRotatey(t *testing.T) {
	cases := map[float64]*Matrix4{
		0:       NewMatrix(1, 0, 0, 0,   0, 1, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1),
		90:      NewMatrix(0, 0, -1, 0,   0, 1, 0, 0,   1, 0, 0, 0,   0, 0, 0, 1),
		180:     NewMatrix(-1, 0, 0, 0,   0, 1, 0, 0,   0, 0, -1, 0,   0, 0, 0, 1),
		270:     NewMatrix(0, 0, 1, 0,   0, 1, 0, 0,   -1, 0, 0, 0,   0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatey(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}

func TestMatrixRotatez(t *testing.T) {
	cases := map[float64]*Matrix4{
		0:       NewMatrix(1, 0, 0, 0,   0, 1, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1),
		90:      NewMatrix(0, 1, 0, 0,   -1, 0, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1),
		180:     NewMatrix(-1, 0, 0, 0,   0, -1, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1),
		270:     NewMatrix(0, -1, 0, 0,   1, 0, 0, 0,   0, 0, 1, 0,   0, 0, 0, 1)}

	for angle, expected := range cases {
		actual := NewIdentity().Rotatez(angle)
		if !expected.Equals(actual, 3) {
			t.Errorf("Expected %s but got %s for angle %d", expected, actual, angle)
		}

	}
}