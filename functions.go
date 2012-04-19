/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"math"
)

// relative coordinates range from -3 to +3
func RelativeCoordinateToBlockId(dx int16, dy int16, dz int16) (id uint16) {
	id = 0
	id |= uint16(dx + 3)
	id |= uint16(dy+3) << 3
	id |= uint16(dz+3) << 6
	return
}

func BlockIdToRelativeCoordinate(id uint16) (dx int16, dy int16, dz int16) {
	dx = int16(id&0x0007 - 3)
	dy = int16((id&0x0038)>>3 - 3)
	dz = int16((id&0x01C0)>>6 - 3)
	return
}

func IntPosition(pos Vectorf) Vectori {
	return Vectori{int16(math.Floor(pos[XAXIS] + 0.5)),
		int16(math.Floor(pos[YAXIS] + 0.5)),
		int16(math.Floor(pos[ZAXIS] + 0.5))}
}

// Finds the surface level for a given x, z coordinate
func FindSurface(x int16, z int16) (y int16) {
	y = GroundLevel
	if TheWorld.At(x, y, z) == BLOCK_AIR {
		for y > 0 && TheWorld.At(x, y, z) == BLOCK_AIR {
			y--
		}
	} else {
		for TheWorld.At(x, y, z) != BLOCK_AIR {
			y++
		}
	}

	return
}
