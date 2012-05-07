/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	// "fmt"
	// "github.com/banthar/gl"
	"github.com/iand/perlin"
	"math"
	"math/rand"
	"time"
)

type World struct {
	mobs   []Mob
	chunks map[uint64]*Chunk
}

type Chunk struct {
	x, z         uint16
	Blocks       [CHUNK_WIDTH * CHUNK_WIDTH * CHUNK_HEIGHT]byte
	vertexBuffer *VertexBuffer
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

	xc, yc, zc := chunkCoordsFromWorld(PLAYER_START_X, self.GroundLevel(PLAYER_START_X, PLAYER_START_Z), PLAYER_START_Z)

	self.GenerateChunk(xc, yc, zc)

	// wolf := new(Wolf)
	// wolf.Init(200, 25, 19, float32(self.FindSurface(25, 19)))
	// self.mobs = append(self.mobs, wolf)

}

func (self *World) GroundLevel(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(8*NOISE_SCALE), float64(z-MAP_DIAM)/(8*NOISE_SCALE), worldSeed, 0.4, 4)
	ground := uint16(noise*(CHUNK_HEIGHT/3.0) + CHUNK_HEIGHT/2.0)
	if ground > CHUNK_HEIGHT {
		ground = CHUNK_HEIGHT
	}
	return ground
}

func (self *World) SoilThickness(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE), float64(z-MAP_DIAM)/(NOISE_SCALE), worldSeed, 0.6, 8)
	if noise < -1 {
		noise = -1
	}
	return uint16(noise*4 + 4)
}

func (self *World) Precipitation(x uint16, z uint16) float64 {
	return perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE), float64(z-MAP_DIAM)/(NOISE_SCALE), worldSeed, 0.6, 12)
}

func (self *World) Drainage(x uint16, z uint16) float64 {
	return perlin.Noise2D(float64(x-MAP_DIAM)/(6*NOISE_SCALE), float64(z-MAP_DIAM)/(6*NOISE_SCALE), worldSeed, 0.4, 12)
}

func (self *World) Rocks(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(16*NOISE_SCALE), float64(z-MAP_DIAM)/(16*NOISE_SCALE), worldSeed, 0.999, 8)
	noise = noise*12 - 20

	if noise < 0 {
		noise = 0
	}
	return uint16(noise)
}

func (self *World) Feature1(x uint16, z uint16) float64 {
	return perlin.Noise2D(float64(x-MAP_DIAM)/(30*NOISE_SCALE), float64(z-MAP_DIAM)/(90*NOISE_SCALE), worldSeed, 0.99, 10)
}

func (self *World) Feature2(x uint16, z uint16) float64 {
	return perlin.Noise2D(float64(x-MAP_DIAM)/(90*NOISE_SCALE), float64(z-MAP_DIAM)/(30*NOISE_SCALE), worldSeed, 0.99, 10)
}

func (self *World) GenerateChunk(cx uint16, cy uint16, cz uint16) *Chunk {
	startTicks := time.Now().UnixNano()
	var chunk Chunk
	chunk.Init(cx, cy, cz)
	self.chunks[chunkIndex(cx, cy, cz)] = &chunk

	println("Generating chunk at x:", cx, " y:", cy, " z:", cz)

	xw := cx * CHUNK_WIDTH
	zw := cz * CHUNK_WIDTH

	for x := uint16(0); x < CHUNK_WIDTH; x++ {
		for z := uint16(0); z < CHUNK_WIDTH; z++ {
			ground := self.GroundLevel(x+xw, z+zw)
			soil := uint16(float64(self.SoilThickness(x+xw, z+zw))*1 - (float64(ground) / CHUNK_HEIGHT))

			if soil > ground {
				soil = ground
			}

			for y := uint16(0); y < ground-soil; y++ {
				chunk.Set(x, y, z, BLOCK_STONE)
			}
			for y := ground - soil; y < ground; y++ {
				chunk.Set(x, y, z, BLOCK_DIRT)
			}

			for y := ground; y < ground+self.Rocks(x+xw, z+zw); y++ {
				chunk.Set(x, y, z, BLOCK_STONE)
			}

		}
	}

	for x := uint16(xw); x < xw+CHUNK_WIDTH; x++ {
		for z := uint16(zw); z < zw+CHUNK_WIDTH; z++ {
			y := self.FindSurface(x, z)
			if self.Precipitation(x, z) > TREE_PRECIPITATION_MIN && rand.Intn(100) < TREE_DENSITY_PCT {

				if y > 1 && y < treeLine && chunk.At(x-xw, y-1, z-zw) == BLOCK_DIRT {
					self.GrowTree(x, y, z)

				}
			} else {
				feature1 := self.Feature1(x, z)
				feature2 := self.Feature2(x, z)

				if feature1 > 0.8 && feature2 > 0.8 {
					// 	self.Set(x, y, z, BLOCK_STONE)
					// 	self.Set(x, y+1, z, BLOCK_STONE)
					// 	self.Set(x, y+2, z, BLOCK_STONE)
					// 	self.Set(x, y+3, z, BLOCK_STONE)
					// 	self.Set(x, y+4, z, BLOCK_STONE)
					// 	self.Set(x, y+5, z, BLOCK_STONE)
					// 	self.Set(x, y+6, z, BLOCK_STONE)
					// } else if feature1 > 0.7 && feature2 > 0.7 {
					// 	self.Set(x, y, z, BLOCK_STONE)
					// 	self.Set(x, y+1, z, BLOCK_STONE)
					// 	self.Set(x, y+2, z, BLOCK_STONE)
					// 	self.Set(x, y+3, z, BLOCK_STONE)
					// 	self.Set(x, y+4, z, BLOCK_STONE)
					// } else if feature1 > 0.6 && feature2 > 0.6 {
					// 	self.Set(x, y, z, BLOCK_STONE)
					// 	self.Set(x, y+1, z, BLOCK_STONE)
					// 	self.Set(x, y+2, z, BLOCK_STONE)
				}
			}
		}
	}

	console.chunkGenerationTime = time.Now().UnixNano() - startTicks
	return &chunk

}

func (self *World) FindSurface(x uint16, z uint16) uint16 {
	y := self.GroundLevel(x, z)
	if self.At(x, y, z) == BLOCK_AIR {
		for y > 0 && self.At(x, y, z) == BLOCK_AIR {
			y--
		}
		y++
	} else {
		for y < CHUNK_HEIGHT && self.At(x, y, z) != BLOCK_AIR {
			y++
		}

	}

	return y
}

// Gets the chunk for a given x/z block coordinate
func (self *World) GetChunkForBlock(x uint16, y uint16, z uint16) (*Chunk, uint16, uint16, uint16) {
	cx, cy, cz := chunkCoordsFromWorld(x, y, z)

	chunk := self.GetChunk(cx, cy, cz)

	ox := x - cx*CHUNK_WIDTH
	oy := y - cy*CHUNK_HEIGHT
	oz := z - cz*CHUNK_WIDTH

	return chunk, ox, oy, oz

}

// Gets the chunk for a given x/z block coordinate
// x = 0, z = 0 is in the top left of the home chunk
func (self *World) GetChunk(cx uint16, cy uint16, cz uint16) *Chunk {
	chunk, ok := self.chunks[chunkIndex(cx, cy, cz)]
	if !ok {
		chunk = self.GenerateChunk(cx, cy, cz)
	}
	return chunk

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

func (self *World) RandomSquare(x1 uint16, z1 uint16, radius uint16) (x uint16, z uint16) {

	x = uint16(rand.Intn(int(radius))) + x1 - radius/2
	z = uint16(rand.Intn(int(radius))) + z1 - radius/2
	return
}

// north/south = -/+ z
// east/west = +/- x
// up/down = +/- y

func (self *World) Grow(x uint16, y uint16, z uint16, n int, s int, w int, e int, u int, d int, texture byte) {
	if y > 0 && x < MAP_DIAM && self.At(x+1, y-1, z) != 0 && rand.Intn(100) < e {
		self.Set(x+1, y, z, texture)
		self.Grow(x+1, y, z, n, s, 0, e-2, u, d, texture)
	}
	if y > 0 && x > 0 && self.At(x-1, y-1, z) != 0 && rand.Intn(100) < w {
		self.Set(x-1, y, z, texture)
		self.Grow(x-1, y, z, n, s, w-2, 0, u, d, texture)
	}
	if y > 0 && z < MAP_DIAM && self.At(x, y-1, z+1) != 0 && rand.Intn(100) < s {
		self.Set(x, y, z+1, texture)
		self.Grow(x, y, z+1, 0, s-2, w, e, u, d, texture)
	}
	if y > 0 && z > 0 && self.At(x, y-1, z-1) != 0 && rand.Intn(100) < n {
		self.Set(x, y, z-1, texture)
		self.Grow(x, y, z-1, n-2, 0, w, e, u, d, texture)
	}
	if y < MAP_DIAM && rand.Intn(100) < u {
		self.Set(x, y+1, z, texture)
		self.Grow(x, y+1, z, n, s, w, e, u-2, 0, texture)
	}
	if y > 0 && rand.Intn(100) < d {
		self.Set(x, y-1, z, texture)
		self.Grow(x, y-1, z, n, s, w, e, 0, d-2, texture)
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
	} else if self.GroundLevel(x-1, z) > y {
		neighbours[WEST_FACE] = BLOCK_DIRT
	} else {
		neighbours[WEST_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x+1, y, z) {
		neighbours[EAST_FACE] = uint16(self.At(x+1, y, z))
	} else if self.GroundLevel(x+1, z) > y {
		neighbours[EAST_FACE] = BLOCK_DIRT
	} else {
		neighbours[EAST_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y, z-1) {
		neighbours[NORTH_FACE] = uint16(self.At(x, y, z-1))
	} else if self.GroundLevel(x, z-1) > y {
		neighbours[NORTH_FACE] = BLOCK_DIRT
	} else {
		neighbours[NORTH_FACE] = BLOCK_AIR
	}

	if self.ChunkLoadedFor(x, y, z+1) {
		neighbours[SOUTH_FACE] = uint16(self.At(x, y, z+1))
	} else if self.GroundLevel(x, z+1) > y {
		neighbours[SOUTH_FACE] = BLOCK_DIRT
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

	playerRect := Rect{x: float64(mp[XAXIS]) + mob.Velocity()[XAXIS]*dt, y: float64(mp[ZAXIS]) + mob.Velocity()[ZAXIS]*dt, sizex: mob.W(), sizey: mob.D()}

	// collisionCandidates := make([]Side, 0)

	if self.Atv(ip.North()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] < 0 && ip.North().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.North()[ZAXIS])+0.5+playerRect.sizey/2, 0)
		}
	}

	if self.Atv(ip.South()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] > 0 && ip.South().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.South()[ZAXIS])-0.5-playerRect.sizey/2, 0)
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

	pxmin, _, pzmin := chunkCoordsFromWorld(uint16(center[XAXIS]-float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]-float64(viewRadius)))
	pxmax, _, pzmax := chunkCoordsFromWorld(uint16(center[XAXIS]+float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]+float64(viewRadius)))

	for px := pxmin; px <= pxmax; px++ {
		for pz := pzmin; pz <= pzmax; pz++ {
			self.GetChunk(px, 0, pz).Render(selectedBlockFace)
		}
	}

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
	yw = y
	return
}

func (chunk *Chunk) Init(x uint16, y uint16, z uint16) {
	chunk.x = x
	// chunk.y = y
	chunk.z = z
	chunk.vertexBuffer = NewVertexBuffer(VERTEX_BUFFER_CAPACITY, terrainTexture)
}

func (chunk *Chunk) At(x uint16, y uint16, z uint16) byte {
	return chunk.Blocks[blockIndex(x, y, z)]
}

func (chunk *Chunk) Set(x uint16, y uint16, z uint16, b byte) {
	chunk.Blocks[blockIndex(x, y, z)] = b
	chunk.clean = false
}

func (self *Chunk) PreRender(selectedBlockFace *BlockFace) {
	if !self.clean || (selectedBlockFace != nil && selectedBlockFace.pos[XAXIS] >= self.x*CHUNK_WIDTH && selectedBlockFace.pos[XAXIS] < (self.x+1)*CHUNK_WIDTH &&
		/*selectedBlockFace.pos[YAXIS] >= self.y*CHUNK_HEIGHT && selectedBlockFace.pos[YAXIS] < (self.y+1)*CHUNK_HEIGHT && */
		selectedBlockFace.pos[ZAXIS] >= self.z*CHUNK_WIDTH && selectedBlockFace.pos[ZAXIS] < (self.z+1)*CHUNK_WIDTH) {
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
						yw := y
						zw := self.z*CHUNK_WIDTH + z

						neighbours := TheWorld.Neighbours(xw, yw, zw)

						if TheWorld.HasVisibleFaces(neighbours) {

							selectedFace := uint8(FACE_NONE)
							if selectedBlockFace != nil && xw == selectedBlockFace.pos[XAXIS] && yw == selectedBlockFace.pos[YAXIS] && zw == selectedBlockFace.pos[ZAXIS] {
								selectedFace = selectedBlockFace.face
							}

							TerrainCube(self.vertexBuffer, float32(xw), float32(yw), float32(zw), neighbours, blockid, selectedFace)
						}
					}
				}
			}
		}

		self.clean = true
	}

}

func (self *Chunk) Render(selectedBlockFace *BlockFace) {
	self.PreRender(selectedBlockFace)

	self.vertexBuffer.RenderDirect()

	// fmt.Printf("Chunk ticks: %4.0f\n", float64(t.GetTicks())/1e6)

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
