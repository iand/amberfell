/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/gl"
	"math"
	"math/rand"
)

type World struct {
	GroundLevel uint16
	mobs        []Mob
	chunks      map[uint64]*Chunk
}

type Chunk struct {
	x, y, z      uint16
	Blocks       [CHUNK_WIDTH * CHUNK_WIDTH * CHUNK_HEIGHT]byte
	vertexBuffer VertexBuffer
	clean        bool
}

type Side struct {
	x, x1, x2, z, z1, z2, dir, y float64
}

type BlockFace struct {
	pos  Vectori
	face uint8
}

type InteractingBlockFace struct {
	blockFace *BlockFace
	hitCount  uint8
}

func (self *World) Init() {

	self.chunks = make(map[uint64]*Chunk)

	xc, yc, zc := chunkCoordsFromWorld(32760, GroundLevel, 32760)

	println("1")
	self.GenerateChunk(xc, yc, zc)
	println("2")
	// self.GenerateChunk(0, 0, 1)
	// self.GenerateChunk(0, 0, -1)
	// self.GenerateChunk(-1, 0, 0)
	// self.GenerateChunk(-1, 0, 1)
	// self.GenerateChunk(-1, 0, -1)
	// self.GenerateChunk(1, 0, 0)
	// self.GenerateChunk(1, 0, 1)
	// self.GenerateChunk(1, 0, -1)

	self.Set(32760, GroundLevel, 32760, 1)   // stone
	self.Set(32760, GroundLevel+1, 32760, 1) // stone
	self.Set(32760, GroundLevel+2, 32760, 1) // stone
	self.Set(32760, GroundLevel+3, 32760, 1) // stone

	var iw, id uint16

	numFeatures := rand.Intn(21)
	for i := 0; i < numFeatures; i++ {
		iw, id = self.RandomSquare()

		self.Set(iw, GroundLevel, id, 1) // stone
		self.Grow(iw, GroundLevel, id, 45, 45, 45, 40, 20, 4, byte(rand.Intn(2))+1)
	}

	for i := 0; i < 40; i++ {
		iw, id = self.RandomSquare()
		self.GrowTree(iw, self.FindSurface(iw, id), id)
	}

	// wolf := new(Wolf)
	// wolf.Init(200, 25, 19, float32(self.FindSurface(25, 19)))
	// self.mobs = append(self.mobs, wolf)

}

// A chunk is a 24 x 24 x 48 set of blocks
// x is east/west offset from World Origin
// z is south/north offset from World Origin
func (self *World) GenerateChunk(x uint16, y uint16, z uint16) *Chunk {
	var chunk Chunk
	chunk.Init(x, y, z)
	println("Generating chunk at x:", x, " y:", y, " z:", z)
	var iw, id, ih uint16
	for iw = 0; iw < CHUNK_WIDTH; iw++ {
		for id = 0; id < CHUNK_WIDTH; id++ {
			for ih = 0; ih < CHUNK_HEIGHT; ih++ {
				if chunk.y*CHUNK_HEIGHT+ih <= GroundLevel {
					chunk.Set(iw, ih, id, 2) // dirt
				} else {
					chunk.Set(iw, ih, id, 0) // dirt
				}
			}
		}
	}

	self.chunks[chunkIndex(x, y, z)] = &chunk
	return &chunk

}

// Gets the chunk for a given x/z block coordinate
// x = 0, z = 0 is in the top left of the home chunk
func (self *World) GetChunkForBlock(x uint16, y uint16, z uint16) (*Chunk, uint16, uint16, uint16) {
	cx, cy, cz := chunkCoordsFromWorld(x, y, z)

	chunk, ok := self.chunks[chunkIndex(cx, cy, cz)]
	if !ok {
		chunk = self.GenerateChunk(cx, cy, cz)
	}

	ox := x - cx*CHUNK_WIDTH
	oy := y - cy*CHUNK_HEIGHT
	oz := z - cz*CHUNK_WIDTH

	return chunk, ox, oy, oz

}

func (self *World) At(x uint16, y uint16, z uint16) byte {
	chunk, ox, oy, oz := self.GetChunkForBlock(x, y, z)
	return chunk.At(ox, oy, oz)
}

func (self *World) Atv(v [3]uint16) byte {
	return self.At(v[XAXIS], v[YAXIS], v[ZAXIS])
}

func (self *World) Set(x uint16, y uint16, z uint16, b byte) {
	chunk, ox, oy, oz := self.GetChunkForBlock(x, y, z)
	chunk.Set(ox, oy, oz, b)
}

func (self *World) Setv(v Vectori, b byte) {
	chunk, ox, oy, oz := self.GetChunkForBlock(v[XAXIS], v[YAXIS], v[ZAXIS])
	chunk.Set(ox, oy, oz, b)
}

func (self *World) RandomSquare() (x uint16, z uint16) {
	radius := int(120)

	x = uint16(rand.Intn(radius) + 32760 - radius/2)
	z = uint16(rand.Intn(radius) + 32760 - radius/2)
	return
}

// north/south = -/+ z
// east/west = +/- x
// up/down = +/- y

func (self *World) Grow(x uint16, y uint16, z uint16, n int, s int, w int, e int, u int, d int, texture byte) {
	if self.At(x+1, y-1, z) != 0 && rand.Intn(100) < e {
		self.Set(x+1, y, z, texture)
		self.Grow(x+1, y, z, n, s, 0, e, u, d, texture)
	}
	if self.At(x-1, y-1, z) != 0 && rand.Intn(100) < w {
		self.Set(x-1, y, z, texture)
		self.Grow(x-1, y, z, n, s, w, 0, u, d, texture)
	}
	if self.At(x, y-1, z+1) != 0 && rand.Intn(100) < s {
		self.Set(x, y, z+1, texture)
		self.Grow(x, y, z+1, 0, s, w, e, u, d, texture)
	}
	if self.At(x, y-1, z-1) != 0 && rand.Intn(100) < n {
		self.Set(x, y, z-1, texture)
		self.Grow(x, y, z-1, n, 0, w, e, u, d, texture)
	}
	if rand.Intn(100) < u {
		self.Set(x, y+1, z, texture)
		self.Grow(x, y+1, z, n, s, w, e, u, 0, texture)
	}
	if rand.Intn(100) < d {
		self.Set(x, y-1, z, texture)
		self.Grow(x, y-1, z, n, s, w, e, 0, d, texture)
	}
}

func (self *World) HasVisibleFaces(neighbours [6]uint16) bool {

	switch neighbours[WEST_FACE] {
	case BLOCK_AIR, BLOCK_LOG_WALL, BLOCK_LOG_SLAB:
		return true
	}

	switch neighbours[EAST_FACE] {
	case BLOCK_AIR, BLOCK_LOG_WALL, BLOCK_LOG_SLAB:
		return true
	}

	switch neighbours[NORTH_FACE] {
	case BLOCK_AIR, BLOCK_LOG_WALL, BLOCK_LOG_SLAB:
		return true
	}

	switch neighbours[SOUTH_FACE] {
	case BLOCK_AIR, BLOCK_LOG_WALL, BLOCK_LOG_SLAB:
		return true
	}

	switch neighbours[UP_FACE] {
	case BLOCK_AIR, BLOCK_LOG_WALL, BLOCK_LOG_SLAB:
		return true
	}

	return false
}

func (self *World) Neighbours(x uint16, y uint16, z uint16) (neighbours [6]uint16) {

	if self.ChunkLoadedFor(x-1, y, z) {
		neighbours[WEST_FACE] = uint16(self.At(x-1, y, z))
	} else {
		neighbours[WEST_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x+1, y, z) {
		neighbours[EAST_FACE] = uint16(self.At(x+1, y, z))
	} else {
		neighbours[EAST_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y, z-1) {
		neighbours[NORTH_FACE] = uint16(self.At(x, y, z-1))
	} else {
		neighbours[NORTH_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y, z+1) {
		neighbours[SOUTH_FACE] = uint16(self.At(x, y, z+1))
	} else {
		neighbours[SOUTH_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y+1, z) {
		neighbours[UP_FACE] = uint16(self.At(x, y+1, z))
	} else {
		neighbours[UP_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y-1, z) {
		neighbours[DOWN_FACE] = uint16(self.At(x, y-1, z))
	} else {
		neighbours[DOWN_FACE] = BLOCK_AIR
	}

	return
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

func (self *World) ApplyForces(mob Mob, dt float64) {
	// mobBounds := mob.DesiredBoundingBox(dt)
	mp := mob.Position()
	ip := IntPosition(mp)

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

	playerRect := Rect{x: float64(mp[XAXIS]) + mob.Velocity()[XAXIS]*dt, z: float64(mp[ZAXIS]) + mob.Velocity()[ZAXIS]*dt, sizex: mob.W(), sizez: mob.D()}

	// collisionCandidates := make([]Side, 0)

	if self.Atv(ip.North()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] < 0 && ip.North().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.North()[ZAXIS])+0.5+playerRect.sizez/2, 0)
		}
	}

	if self.Atv(ip.South()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] > 0 && ip.South().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.South()[ZAXIS])-0.5-playerRect.sizez/2, 0)
		}
	}

	if self.Atv(ip.East()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] > 0 && ip.East().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.East()[XAXIS])-0.5-playerRect.sizex/2, 0)
		}
	}

	if self.Atv(ip.West()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] < 0 && ip.West().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.West()[XAXIS])+0.5+playerRect.sizex/2, 0)
		}
	}

	mob.SetFalling(true)
	if self.Atv(ip.Down()) != BLOCK_AIR {
		mob.SetFalling(false)
		if mob.Velocity()[YAXIS] < 0 {
			mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
		}
	} else {
		if self.Atv(ip.Down().North()) != BLOCK_AIR {
			if ip.Down().North().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if self.Atv(ip.Down().South()) != BLOCK_AIR {
			if ip.Down().South().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if self.Atv(ip.Down().East()) != BLOCK_AIR {
			if ip.Down().East().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}

		if self.Atv(ip.Down().West()) != BLOCK_AIR {
			if ip.Down().West().HRect().Intersects(playerRect) {
				mob.Snapy(float64(ip.Down()[YAXIS])+1, 0)
				mob.SetFalling(false)
			}
		}
	}

}

func (self *World) Simulate(dt float64) {
	for _, v := range self.mobs {
		v.Act(dt)
		self.ApplyForces(v, dt)
		v.Update(dt)
	}

}

func (self World) ChunkLoadedFor(x uint16, y uint16, z uint16) bool {
	cx := x / CHUNK_WIDTH
	cy := y / CHUNK_HEIGHT
	cz := z / CHUNK_WIDTH

	_, ok := self.chunks[chunkIndex(cx, cy, cz)]
	return ok
}

func (self *World) Draw(center Vectorf, selectedBlockFace *BlockFace) {
	for _, v := range self.mobs {
		v.Draw(center, selectedBlockFace)
	}

	metrics.cubecount = 0
	metrics.vertices = 0

	terrainTexture.Bind(gl.TEXTURE_2D)

	px := uint16(center[XAXIS]) / CHUNK_WIDTH
	py := uint16(center[YAXIS]) / CHUNK_HEIGHT
	pz := uint16(center[ZAXIS]) / CHUNK_WIDTH

	c1, ok := self.chunks[chunkIndex(px, py, pz)]
	if !ok {
		c1 = self.GenerateChunk(px, py, pz)
	}
	c1.Render(selectedBlockFace)
	RenderQuadIndex(&c1.vertexBuffer)
	metrics.vertices += c1.vertexBuffer.vertexCount

	c2, ok := self.chunks[chunkIndex(px+1, py, pz)]
	if !ok {
		c2 = self.GenerateChunk(px+1, py, pz)
	}
	c2.Render(selectedBlockFace)
	RenderQuadIndex(&c2.vertexBuffer)
	metrics.vertices += c2.vertexBuffer.vertexCount

	c3, ok := self.chunks[chunkIndex(px-1, py, pz)]
	if !ok {
		c3 = self.GenerateChunk(px-1, py, pz)
	}
	c3.Render(selectedBlockFace)
	RenderQuadIndex(&c3.vertexBuffer)
	metrics.vertices += c3.vertexBuffer.vertexCount

	c4, ok := self.chunks[chunkIndex(px, py, pz+1)]
	if !ok {
		c4 = self.GenerateChunk(px, py, pz+1)
	}
	c4.Render(selectedBlockFace)
	RenderQuadIndex(&c4.vertexBuffer)
	metrics.vertices += c4.vertexBuffer.vertexCount

	c5, ok := self.chunks[chunkIndex(px, py, pz-1)]
	if !ok {
		c5 = self.GenerateChunk(px, py, pz-1)
	}
	c5.Render(selectedBlockFace)
	RenderQuadIndex(&c5.vertexBuffer)
	metrics.vertices += c5.vertexBuffer.vertexCount

	c6, ok := self.chunks[chunkIndex(px+1, py, pz+1)]
	if !ok {
		c6 = self.GenerateChunk(px+1, py, pz+1)
	}
	c6.Render(selectedBlockFace)
	RenderQuadIndex(&c6.vertexBuffer)
	metrics.vertices += c6.vertexBuffer.vertexCount

	c7, ok := self.chunks[chunkIndex(px+1, py, pz-1)]
	if !ok {
		c7 = self.GenerateChunk(px+1, py, pz-1)
	}
	c7.Render(selectedBlockFace)
	RenderQuadIndex(&c7.vertexBuffer)
	metrics.vertices += c7.vertexBuffer.vertexCount

	c8, ok := self.chunks[chunkIndex(px-1, py, pz+1)]
	if !ok {
		c8 = self.GenerateChunk(px-1, py, pz+1)
	}
	c8.Render(selectedBlockFace)
	RenderQuadIndex(&c8.vertexBuffer)
	metrics.vertices += c8.vertexBuffer.vertexCount

	c9, ok := self.chunks[chunkIndex(px-1, py, pz-1)]
	if !ok {
		c9 = self.GenerateChunk(px-1, py, pz-1)
	}
	c9.Render(selectedBlockFace)
	RenderQuadIndex(&c9.vertexBuffer)
	metrics.vertices += c9.vertexBuffer.vertexCount
}

// Finds the surface level for a given x, z coordinate
func (self *World) FindSurface(x uint16, z uint16) (y uint16) {
	y = GroundLevel
	if self.At(x, y, z) == BLOCK_AIR {
		for y > 0 && self.At(x, y, z) == BLOCK_AIR {
			y--
		}
	} else {
		for self.At(x, y, z) != BLOCK_AIR {
			y++
		}
	}

	return
}

func chunkCoordsFromWorld(x uint16, y uint16, z uint16) (cx uint16, cy uint16, cz uint16) {
	cx = uint16(math.Floor(float64(x) / CHUNK_WIDTH))
	cy = uint16(math.Floor(float64(y) / CHUNK_HEIGHT))
	cz = uint16(math.Floor(float64(z) / CHUNK_WIDTH))

	return
}

func chunkIndexFromWorld(x uint16, y uint16, z uint16) uint64 {
	cx, cy, cz := chunkCoordsFromWorld(x, y, z)
	return chunkIndex(cx, cy, cz)
}

func chunkIndex(cx uint16, cy uint16, cz uint16) uint64 {
	return uint64(cz)<<32 | uint64(cy)<<16 | uint64(cx)
}

func chunkCoordsFromindex(index uint64) (cx uint16, cy uint16, cz uint16) {
	cx = uint16(index)
	cy = uint16(index >> 16)
	cz = uint16(index >> 32)

	return
}

func blockIndex(x uint16, y uint16, z uint16) uint16 {
	return CHUNK_WIDTH*CHUNK_WIDTH*y + CHUNK_WIDTH*x + z
}

// **************************************************************
// CHUNKS
// **************************************************************

func (c Chunk) WorldCoords(x uint16, y uint16, z uint16) (xw uint16, yw uint16, zw uint16) {
	xw = c.x*CHUNK_WIDTH + x
	zw = c.z*CHUNK_WIDTH + z
	yw = c.y*CHUNK_HEIGHT + y
	return
}

func (chunk *Chunk) Init(x uint16, y uint16, z uint16) {
	chunk.x = x
	chunk.y = y
	chunk.z = z
}

func (chunk *Chunk) At(x uint16, y uint16, z uint16) byte {
	return chunk.Blocks[blockIndex(x, y, z)]
}

func (chunk *Chunk) Set(x uint16, y uint16, z uint16, b byte) {
	chunk.Blocks[blockIndex(x, y, z)] = b
	chunk.clean = false
}

func (self *Chunk) Render(selectedBlockFace *BlockFace) {

	if self.clean && !(selectedBlockFace != nil && selectedBlockFace.pos[XAXIS] >= self.x*CHUNK_WIDTH && selectedBlockFace.pos[XAXIS] < (self.x+1)*CHUNK_WIDTH &&
		selectedBlockFace.pos[YAXIS] >= self.y*CHUNK_HEIGHT && selectedBlockFace.pos[YAXIS] < (self.y+1)*CHUNK_HEIGHT &&
		selectedBlockFace.pos[ZAXIS] >= self.z*CHUNK_WIDTH && selectedBlockFace.pos[ZAXIS] < (self.z+1)*CHUNK_WIDTH) {
		return
	}
	t := Timer{}
	t.Start()
	self.vertexBuffer.Reset()
	var x, y, z uint16
	for x = 0; x < CHUNK_WIDTH; x++ {
		for z = 0; z < CHUNK_WIDTH; z++ {
			for y = 0; y < CHUNK_HEIGHT; y++ {

				var blockid byte = self.Blocks[blockIndex(x, y, z)]
				if blockid != 0 {
					xw := self.x*CHUNK_WIDTH + x
					yw := self.y*CHUNK_HEIGHT + y
					zw := self.z*CHUNK_WIDTH + z

					// xw, yw, zw := self.WorldCoords(x, y, z)
					//neighbours := TheWorld.Neighbours(xw, yw, zw)
					var neighbours [6]uint16
					// if x == 0 || x == CHUNK_WIDTH-1 || z == 0 || z == CHUNK_WIDTH-1 || y == 0 || y == CHUNK_HEIGHT-1 {
					neighbours = TheWorld.Neighbours(xw, yw, zw)
					// } else {

					// 	neighbours = [6]uint16{
					// 		uint16(self.Blocks[blockIndex(x+1, y, z)]),
					// 		uint16(self.Blocks[blockIndex(x-1, y, z)]),
					// 		uint16(self.Blocks[blockIndex(x, y, z-1)]),
					// 		uint16(self.Blocks[blockIndex(x, y, z+1)]),
					// 		uint16(self.Blocks[blockIndex(x, y+1, z)]),
					// 		uint16(self.Blocks[blockIndex(x, y-1, z)]),
					// 	}
					// }

					if TheWorld.HasVisibleFaces(neighbours) {

						selectedFace := uint8(FACE_NONE)
						if selectedBlockFace != nil && xw == selectedBlockFace.pos[XAXIS] && yw == selectedBlockFace.pos[YAXIS] && zw == selectedBlockFace.pos[ZAXIS] {
							selectedFace = selectedBlockFace.face
						}

						TerrainCube(&self.vertexBuffer, float32(xw), float32(yw), float32(zw), neighbours, blockid, selectedFace)
						metrics.cubecount++
					}
				}
			}
		}
	}

	metrics.vertices += self.vertexBuffer.vertexCount
	self.clean = true

	fmt.Printf("Chunk ticks: %4.0f\n", float64(t.GetTicks())/1e6)

}

func (self *World) GrowTree(x uint16, y uint16, z uint16) {
	self.Set(x, y, z, BLOCK_TRUNK)
	self.Set(x, y+1, z, BLOCK_TRUNK)
	self.Set(x, y+2, z, BLOCK_TRUNK)
	self.Set(x, y+3, z, BLOCK_TRUNK)
	self.Set(x+1, y+3, z, BLOCK_LEAVES)
	self.Set(x-1, y+3, z, BLOCK_LEAVES)
	self.Set(x, y+3, z+1, BLOCK_LEAVES)
	self.Set(x, y+3, z-1, BLOCK_LEAVES)

	self.GrowBranch(x, y+3, z, NORTH_FACE, 50)
	self.GrowBranch(x, y+3, z, EAST_FACE, 50)
	self.GrowBranch(x, y+3, z, WEST_FACE, 50)
	self.GrowBranch(x, y+3, z, SOUTH_FACE, 50)

	if rand.Intn(100) < 50 {
		self.Set(x, y+4, z, BLOCK_TRUNK)
		self.Set(x, y+5, z, BLOCK_TRUNK)
		self.GrowBranch(x, y+5, z, NORTH_FACE, 50)
		self.GrowBranch(x, y+5, z, EAST_FACE, 50)
		self.GrowBranch(x, y+5, z, WEST_FACE, 50)
		self.GrowBranch(x, y+5, z, SOUTH_FACE, 50)
	}

	if rand.Intn(100) < 30 {
		self.Set(x, y+6, z, BLOCK_TRUNK)
		self.Set(x, y+7, z, BLOCK_TRUNK)
		self.GrowBranch(x, y+7, z, NORTH_FACE, 50)
		self.GrowBranch(x, y+7, z, EAST_FACE, 50)
		self.GrowBranch(x, y+7, z, WEST_FACE, 50)
		self.GrowBranch(x, y+7, z, SOUTH_FACE, 50)
	}
}

func (self *World) GrowBranch(x uint16, y uint16, z uint16, face uint8, chance int) {
	newx, newy, newz := x, y, z
	if face == NORTH_FACE {
		newz = z - 1
	} else if face == SOUTH_FACE {
		newz = z + 1
	} else if face == WEST_FACE {
		newx = x - 1
	} else if face == EAST_FACE {
		newx = x + 1
	} else if face == UP_FACE {
		newy = y + 1
	} else if face == DOWN_FACE {
		newy = y - 1
	}
	if rand.Intn(100) < chance {
		self.Set(newx, newy, newz, BLOCK_TRUNK)
		if face != SOUTH_FACE {
			if rand.Intn(100) < 50 {
				self.GrowBranch(newx, newy, newz, NORTH_FACE, chance/3)
			} else {
				self.Set(newx, newy, newz-1, BLOCK_LEAVES)
			}
		}

		if face != NORTH_FACE {
			if rand.Intn(100) < 50 {
				self.GrowBranch(newx, newy, newz, SOUTH_FACE, chance/3)
			} else {
				self.Set(newx, newy, newz+1, BLOCK_LEAVES)
			}
		}

		if face != EAST_FACE {
			if rand.Intn(100) < 50 {
				self.GrowBranch(newx, newy, newz, WEST_FACE, chance/3)
			} else {
				self.Set(newx-1, newy, newz, BLOCK_LEAVES)
			}
		}

		if face != WEST_FACE {
			if rand.Intn(100) < 50 {
				self.GrowBranch(newx, newy, newz, EAST_FACE, chance/3)
			} else {
				self.Set(newx+1, newy, newz, BLOCK_LEAVES)
			}
		}

		if face != DOWN_FACE {
			if rand.Intn(100) < 30 {
				self.GrowBranch(newx, newy, newz, UP_FACE, chance/3)
			} else {
				self.Set(newx, newy+1, newz, BLOCK_LEAVES)
			}
		}
		if face != UP_FACE {
			if rand.Intn(100) < 50 {
				self.GrowBranch(newx, newy, newz, DOWN_FACE, chance/3)
			} else {
				self.Set(newx, newy-1, newz, BLOCK_LEAVES)
			}
		}
	} else {
		self.Set(newx, newy, newz, BLOCK_LEAVES)

	}
}
