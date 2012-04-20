/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
  "testing"
)

// func TestHitsBox(t *testing.T) {

// 	ray := Ray{ &Vectorf{0, 0, 0}, &Vectorf{1, 0, 0} }
// 	box := Box{ &Vectorf{2, -1, -1}, &Vectorf{3, 1, 1} }
// 	expected := true

// 	actual := ray.HitsBox(&box)
// 	if expected != actual {
// 		if expected {
// 			t.Errorf("Ray %s did not hit box %s as expected", ray, box)
// 		} else {
// 			t.Errorf("Ray %s did not miss box %s as expected", ray, box)
// 		}
// 	}
// }

func TestHitsBox2(t *testing.T) {
	cases := map[Ray]bool {
		Ray{ &Vectorf{-8.667095, 13.666898, 16.215151}, &Vectorf{0.851651, -0.422618, -0.309976} } : true }

	box := Box{ &Vectorf{8.5,4.5,9.5}, &Vectorf{9.5,5.5,10.5} }

	for ray, expected := range cases {
		actual := ray.HitsBox(&box)
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