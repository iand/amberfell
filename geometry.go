/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type Ray struct {
	origin *Vectorf
	dir    *Vectorf
}

type Box struct {
	min    *Vectorf
	max    *Vectorf
}

// Ported from http://tog.acm.org/resources/GraphicsGems/gems/RayBox.c
func (self *Ray) HitsBox(box *Box) bool {
	const RIGHT = 	0
	const LEFT = 	1
	const MIDDLE = 	2

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
	if inside {
		return true
	}	

	// Calculate T distances to candidate planes
	for i := 0; i < 3; i++ {
		if quadrant[i] != MIDDLE && self.dir[i] == 0 {
			maxT[i] = (candidatePlane[i]-self.origin[i]) / self.dir[i]
		} else {
			maxT[i] = -1.0
		}
	}


	// Get largest of the maxT's for final choice of intersection
	whichPlane = 0;
	for i := 1; i < 3; i++ {
		if maxT[whichPlane] < maxT[i] {
			whichPlane = i
		}
	}

	// Check final candidate actually inside box
	if maxT[whichPlane] < 0 {
	 return false
	}

	for i := 0; i < 3; i++ {
		if whichPlane != i {
			coord[i] = self.origin[i] + maxT[whichPlane] * self.dir[i]
			if coord[i] < box.min[i] || coord[i] > box.max[i] {
				return false
			}
		} else {
			coord[i] = candidatePlane[i]
		}
	}
	return true	// ray hits box


}
