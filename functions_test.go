/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	// "fmt"
	// "fmt"
	// "math"
	"testing"
)

func TestRound(t *testing.T) {

	cases := map[float64]float64{
		float64(0):    float64(0),
		float64(0.1):  float64(0),
		float64(0.5):  float64(1),
		float64(-0.1): float64(0),
		float64(-0.5): float64(-1),
	}

	for val, expected := range cases {
		actual := Round(val, 0)
		if actual != expected {
			t.Errorf("Expected %f but got %f", expected, actual)
		}
	}
}
