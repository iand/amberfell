/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
// "fmt"
// "math"
// "testing"
)

// func BenchmarkDraw(b *testing.B) {
// 	b.StopTimer()
// 	TheWorld = NewWorld()
// 	center := Vectorf{100, 100, 50}
// 	pxmin, pzmin := chunkCoordsFromWorld(uint16(center[XAXIS]-float64(viewRadius)), uint16(center[ZAXIS]-float64(viewRadius)))
// 	pxmax, pzmax := chunkCoordsFromWorld(uint16(center[XAXIS]+float64(viewRadius)), uint16(center[ZAXIS]+float64(viewRadius)))
// 	for px := pxmin; px <= pxmax; px++ {
// 		for pz := pzmin; pz <= pzmax; pz++ {
// 			TheWorld.GetChunk(px, pz)
// 		}
// 	}

// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for _,chunk := range TheWorld.chunks {
// 			chunk.clean = false
// 		}
// 		TheWorld.Draw(center, nil)
// 	}
// }
