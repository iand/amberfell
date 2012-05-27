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

func BenchmarkAllNeighbours(b *testing.B) {
	b.StopTimer()
	TheWorld = NewWorld()

	x := uint16(2*CHUNK_WIDTH + CHUNK_WIDTH/2)
	z := uint16(2*CHUNK_WIDTH + CHUNK_WIDTH/2)
	y := TheWorld.GroundLevel(x, z)

	center := Vectori{x, y, z}
	px, pz := chunkCoordsFromWorld(x, z)
	TheWorld.GetChunk(px, pz)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		neighbours := TheWorld.AllNeighbours(center)
		_ = neighbours[WEST_FACE]
	}
}
