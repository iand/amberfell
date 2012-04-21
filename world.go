/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
	"math/rand"
	// "math"
	// "fmt"   

)

type World struct {
	GroundLevel int16
	mobs        []Mob
	chunks      map[int16]*Chunk
}

type Chunk struct {
	x, y, z int16
	Blocks  []byte
}

type Side struct {
	x, x1, x2, z, z1, z2, dir, y float64
}




func (world *World) Init() {

	world.chunks = make(map[int16]*Chunk)

	world.GenerateChunk(0, 0, 0)
	world.GenerateChunk(0, 0, 1)
	world.GenerateChunk(0, 0, -1)
	world.GenerateChunk(-1, 0, 0)
	world.GenerateChunk(-1, 0, 1)
	world.GenerateChunk(-1, 0, -1)
	world.GenerateChunk(1, 0, 0)
	world.GenerateChunk(1, 0, 1)
	world.GenerateChunk(1, 0, -1)

	var iw, id int16

	numFeatures := rand.Intn(20)
	for i := 0; i < numFeatures; i++ {
		iw, id = world.RandomSquare()

		world.Set(iw, GroundLevel, id, 1) // stone
		world.Grow(iw, GroundLevel, id, 45, 45, 45, 52, 10, 10, 1)
	}
	iw, id = world.RandomSquare()

	world.Set(iw, GroundLevel, id, 0) // air
	world.Grow(iw, GroundLevel, id, 20, 20, 20, 20, 0, 30, 0)

	wolf := new(Wolf)
	wolf.Init(120, 17, 19, GroundLevel+1)
	world.mobs = append(world.mobs, wolf)

}

// A chunk is a 24 x 24 x 48 set of blocks
// x is east/west offset from World Origin
// z is south/north offset from World Origin
func (world *World) GenerateChunk(x int16, y int16, z int16) *Chunk {
	var chunk Chunk
	chunk.Init(x, y, z)
	println("Generating chunk at x:", x, " y:", y, " z:", z)
	var iw, id, ih int16
	for iw = 0; iw < CHUNK_WIDTH; iw++ {
		for id = 0; id < CHUNK_WIDTH; id++ {
			for ih = 0; ih <= GroundLevel; ih++ {
				chunk.Set(iw, ih, id, 2) // dirt
			}
			for ih = GroundLevel + 1; ih < CHUNK_HEIGHT; ih++ {
				chunk.Set(iw, ih, id, 0) // air
			}
		}
	}

	world.chunks[chunkIndex(x, y, z)] = &chunk
	return &chunk

}

// Gets the chunk for a given x/z block coordinate
// x = 0, z = 0 is in the top left of the home chunk
func (world *World) GetChunkForBlock(x int16, y int16, z int16) (*Chunk, int16, int16, int16) {
	cx := x / CHUNK_WIDTH
	cy := y / CHUNK_HEIGHT
	cz := z / CHUNK_WIDTH
	//println("cx:", cx, "cz:", cz)

	chunk, ok := world.chunks[chunkIndex(cx, cy, cz)]
	if !ok {
		chunk = world.GenerateChunk(cx, cy, cz)
	}

	ox := x - cx*CHUNK_WIDTH
	if ox < 0 {
		ox += CHUNK_WIDTH
	}

	oy := y - cy*CHUNK_HEIGHT
	if oy < 0 {
		oy += CHUNK_HEIGHT
	}

	oz := z - cz*CHUNK_WIDTH
	if oz < 0 {
		oz += CHUNK_WIDTH
	}

	return chunk, ox, oy, oz

}

func (world *World) At(x int16, y int16, z int16) byte {
	//println("x:", x, " y:", y, "z:", z)
	chunk, ox, oy, oz := world.GetChunkForBlock(x, y, z)
	//println("ox:", ox, " y:", y, "oz:", oz)
	return chunk.At(ox, oy, oz)
}

func (world *World) Atv(v Vectori) byte {
	return world.At(v[XAXIS], v[YAXIS], v[ZAXIS])
}

func (world *World) Set(x int16, y int16, z int16, b byte) {
	chunk, ox, oy, oz := world.GetChunkForBlock(x, y, z)
	chunk.Set(ox, oy, oz, b)
}

func (world *World) Setv(v Vectori, b byte) {
	chunk, ox, oy, oz := world.GetChunkForBlock(v[XAXIS], v[YAXIS], v[ZAXIS])
	chunk.Set(ox, oy, oz, b)
}

func (world *World) RandomSquare() (x int16, z int16) {
	x = int16(rand.Intn(40) - 20)
	z = int16(rand.Intn(40) - 20)
	return
}

// north/south = -/+ z
// east/west = +/- x
// up/down = +/- y

func (world *World) Grow(x int16, y int16, z int16, n int, s int, w int, e int, u int, d int, texture byte) {
	if (y == 0 || world.At(x+1, y-1, z) != 0) && rand.Intn(100) < e {
		world.Set(x+1, y, z, texture)
		world.Grow(x+1, y, z, n, s, 0, e, u, d, texture)
	}
	if (y == 0 || world.At(x-1, y-1, z) != 0) && rand.Intn(100) < w {
		world.Set(x-1, y, z, texture)
		world.Grow(x-1, y, z, n, s, w, 0, u, d, texture)
	}
	if (y == 0 || world.At(x, y-1, z+1) != 0) && rand.Intn(100) < s {
		world.Set(x, y, z+1, texture)
		world.Grow(x, y, z+1, 0, s, w, e, u, d, texture)
	}
	if (y == 0 || world.At(x, y-1, z-1) != 0) && rand.Intn(100) < n {
		world.Set(x, y, z-1, texture)
		world.Grow(x, y, z-1, n, 0, w, e, u, d, texture)
	}
	if y < CHUNK_HEIGHT-1 && rand.Intn(100) < u {
		world.Set(x, y+1, z, texture)
		world.Grow(x, y+1, z, n, s, w, e, u, 0, texture)
	}
	if y > 0 && rand.Intn(100) < d {
		world.Set(x, y-1, z, texture)
		world.Grow(x, y-1, z, n, s, w, e, 0, d, texture)
	}
}

func (world *World) AirNeighbours(x int16, z int16, y int16) (n, s, w, e, u, d bool) {

	if world.ChunkLoadedFor(x-1, y, z) && world.At(x-1, y, z) == BLOCK_AIR {
		w = true
	}
	if world.ChunkLoadedFor(x+1, y, z) && world.At(x+1, y, z) == BLOCK_AIR {
		e = true
	}
	if world.ChunkLoadedFor(x, y, z-1) && world.At(x, y, z-1) == BLOCK_AIR {
		n = true
	}
	if world.ChunkLoadedFor(x, y, z+1) && world.At(x, y, z+1) == BLOCK_AIR {
		s = true
	}
	if world.ChunkLoadedFor(x, y+1, z) && world.At(x, y+1, z) == BLOCK_AIR {
		u = true
	}
	return
}

func (world *World) AirNeighbour(x int16, z int16, y int16, face int) bool {
	if face == UP_FACE && world.ChunkLoadedFor(x, y+1, z) && world.At(x, y+1, z) == BLOCK_AIR { 
		return true 
	}
	if face == NORTH_FACE && world.ChunkLoadedFor(x, y, z-1) && world.At(x, y, z-1) == BLOCK_AIR { 
		return true 
	}
	if face == SOUTH_FACE && world.ChunkLoadedFor(x, y, z+1) && world.At(x, y, z+1) == BLOCK_AIR { 
		return true 
	}
	if face == EAST_FACE && world.ChunkLoadedFor(x+1, y, z) && world.At(x+1, y, z) == BLOCK_AIR { 
		return true 
	}
	if face == WEST_FACE && world.ChunkLoadedFor(x-1, y, z) && world.At(x-1, y, z) == BLOCK_AIR { 
		return true 
	}
	if face == DOWN_FACE && world.ChunkLoadedFor(x, y-1, z) && world.At(x, y-1, z) == BLOCK_AIR { 
		return true 
	}
	return false
}


// lineRectCollide( line, rect )
//
// Checks if an axis-aligned line and a bounding box overlap.
// line = { z, x1, x2 } or line = { x, z1, z2 }
// rect = { x, z, size }

func lineRectCollide(line Side, rect Rect) (ret bool) {
	if line.z != 0 {
		ret = rect.z > line.z-rect.sizez/2 && rect.z < line.z+rect.sizez/2 && rect.x > line.x1-rect.sizex/2 && rect.x < line.x2+rect.sizex/2
	} else {
		ret = rect.x > line.x-rect.sizex/2 && rect.x < line.x+rect.sizex/2 && rect.z > line.z1-rect.sizez/2 && rect.z < line.z2+rect.sizez/2
	}
	return
}

// rectRectCollide( r1, r2 )
//
// Checks if two rectangles (x1, y1, x2, y2) overlap.

func rectRectCollide(r1 Side, r2 Side) bool {
	if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 {
		return true
	}
	if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 {
		return true
	}
	if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
		return true
	}
	if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
		return true
	}
	return false
}

func (world *World) ApplyForces(mob Mob, dt float64) {
	// mobBounds := mob.DesiredBoundingBox(dt)
	ip := IntPosition(mob.Position())

	// mobx := ip[XAXIS]
	// moby := ip[YAXIS]
	// mobz := ip[ZAXIS]

	// Gravity
	if mob.IsFalling() {
		// println("is falling")
		mob.Setvx(mob.Velocity()[XAXIS] / 1.001)
		mob.Setvy(mob.Velocity()[YAXIS] - 15*dt)
		mob.Setvz(mob.Velocity()[ZAXIS] / 1.001)
	} else {
		mob.Setvx(mob.Velocity()[XAXIS] / 1.2)
		//mob.Setvy(0)
		mob.Setvz(mob.Velocity()[ZAXIS] / 1.2)
	}

	// var dx, dz, dy int16
	// var x,  z int16

	playerRect := Rect{x: float64(mob.X()) + mob.Velocity()[XAXIS]*dt, z: float64(mob.Z()) + mob.Velocity()[ZAXIS]*dt, sizex: mob.W(), sizez: mob.D()}

	// collisionCandidates := make([]Side, 0)

	if world.Atv(ip.North()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] < 0 && ip.North().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.North()[ZAXIS])+0.5+playerRect.sizez/2, 0)
		}
	}

	if world.Atv(ip.South()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] > 0 && ip.South().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.South()[ZAXIS])-0.5-playerRect.sizez/2, 0)
		}
	}

	if world.Atv(ip.East()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] > 0 && ip.East().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.East()[XAXIS])-0.5-playerRect.sizex/2, 0)
		}
	}

	if world.Atv(ip.West()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] < 0 && ip.West().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.West()[XAXIS])+0.5+playerRect.sizex/2, 0)
		}
	}

	mob.SetFalling(true)
	if world.Atv(ip.Down()) != BLOCK_AIR {
		mob.SetFalling(false)
		if mob.Velocity()[YAXIS] < 0 {
			mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
		}
	} else {
		if world.Atv(ip.Down().North()) != BLOCK_AIR {
			if ip.Down().North().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if world.Atv(ip.Down().South()) != BLOCK_AIR {
			if ip.Down().South().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if world.Atv(ip.Down().East()) != BLOCK_AIR {
			if ip.Down().East().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if world.Atv(ip.Down().West()) != BLOCK_AIR {
			if ip.Down().West().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}
	}

}

func (world *World) Simulate(dt float64) {
	for _, v := range world.mobs {
		v.Act(dt)
		world.ApplyForces(v, dt)
		v.Update(dt)
	}

}

func (world World) ChunkLoadedFor(x int16, y int16, z int16) bool {
	cx := x / CHUNK_WIDTH
	cy := y / CHUNK_HEIGHT
	cz := z / CHUNK_WIDTH

	_, ok := world.chunks[chunkIndex(cx, cy, cz)]
	return ok
}

func (world *World) Draw(center Vectorf) {
	for _, v := range world.mobs {
		v.Draw(center)
	}

	//gl.Translatef(-float32(center[XAXIS]), -float32(center[YAXIS]), -float32(center[ZAXIS]))

	var px, py, pz = int16(center[XAXIS]), int16(center[YAXIS]), int16(center[ZAXIS])

	var x, y, z int16

	count := 0
	for x = px - 30; x < px+30; x++ {
		for z = pz - 30; z < pz+30; z++ {
			if x+z-px-pz <= ViewRadius && x+z-px-pz >= -ViewRadius {
				for y = py - 5; y < py+16; y++ {

					var blockid byte = world.At(x, y, z)
					if blockid != 0 {
						var n, s, w, e, u, d bool = world.AirNeighbours(x, z, y)
						if n || s || w || e || u || d {

							gl.PushMatrix()
							gl.Translatef(float32(x), float32(y), float32(z))
							TerrainCube(n, s, w, e, u, d, blockid)
							count++
							gl.PopMatrix()
						}
					}
				}
			}
		}
	}
	//println("Drew ", count, " cubes")

}

func chunkIndex(x int16, y int16, z int16) int16 {
	return z*CHUNK_WIDTH*CHUNK_WIDTH + x*CHUNK_WIDTH + y
}

func blockIndex(x int16, y int16, z int16) int16 {
	return CHUNK_WIDTH*CHUNK_WIDTH*y + CHUNK_WIDTH*x + z
}

// **************************************************************
// CHUNKS
// **************************************************************

func (c Chunk) WorldCoords(x int16, y int16, z int16) (xw int16, yw int16, zw int16) {
	xw = c.x*CHUNK_WIDTH + x
	zw = c.z*CHUNK_WIDTH + z
	yw = c.y*CHUNK_HEIGHT + y
	return
}

func (chunk *Chunk) Init(x int16, y int16, z int16) {
	chunk.x = x
	chunk.y = y
	chunk.z = z
	chunk.Blocks = make([]byte, CHUNK_WIDTH*CHUNK_WIDTH*CHUNK_HEIGHT)
}

func (chunk *Chunk) At(x int16, y int16, z int16) byte {
	return chunk.Blocks[blockIndex(x, y, z)]
}

func (chunk *Chunk) Set(x int16, y int16, z int16, b byte) {
	chunk.Blocks[blockIndex(x, y, z)] = b
}
