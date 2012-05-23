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

func BenchmarkNewChunk(b *testing.B) {
	b.StopTimer()
	TheWorld = NewWorld()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		chunk := NewChunk(1, 1, false, [2]uint16{0, 0})
		chunk.At(0, 0, 0)
	}
}

func BenchmarkRender(b *testing.B) {
	b.StopTimer()
	TheWorld = NewWorld()
	chunk := NewChunk(20, 20, false, [2]uint16{0, 0})

	var adjacents [4]*Chunk
	adjacents[NORTH_FACE] = NewChunk(20, 19, false, [2]uint16{0, 0})
	adjacents[SOUTH_FACE] = NewChunk(20, 21, false, [2]uint16{0, 0})
	adjacents[EAST_FACE] = NewChunk(21, 20, false, [2]uint16{0, 0})
	adjacents[WEST_FACE] = NewChunk(19, 20, false, [2]uint16{0, 0})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		chunk.clean = false
		chunk.Render(adjacents, nil, nil)
	}
}
