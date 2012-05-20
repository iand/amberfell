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
	mobs             []Mob
	chunks           map[uint64]*Chunk
	amberfell        map[uint64][2]uint16
	lightSources     map[Vectori]LightSource
	timedObjects     map[Vectori]TimedObject
	containerObjects map[Vectori]ContainerObject
	generatorObjects map[Vectori]GeneratorObject
	genseed          int64
	lastSimulated    int64
	campfires        map[Vectori]*CampFire
}

type Chunk struct {
	x, z              uint16
	Blocks            [CHUNK_WIDTH * CHUNK_WIDTH * CHUNK_HEIGHT]byte
	vertexBuffer      *VertexBuffer
	clean             bool
	standingStoneProb float64
	featuresLoaded    bool
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

func NewWorld() *World {
	world := &World{}

	world.genseed = worldSeed
	world.chunks = make(map[uint64]*Chunk)
	world.amberfell = make(map[uint64][2]uint16)
	world.timedObjects = make(map[Vectori]TimedObject)
	world.containerObjects = make(map[Vectori]ContainerObject)
	world.lightSources = make(map[Vectori]LightSource)
	world.generatorObjects = make(map[Vectori]GeneratorObject)
	world.campfires = make(map[Vectori]*CampFire)
	world.lastSimulated = time.Now().UnixNano()
	xc, yc, zc := chunkCoordsFromWorld(PLAYER_START_X, world.GroundLevel(PLAYER_START_X, PLAYER_START_Z), PLAYER_START_Z)

	world.GenerateAmberfell()
	world.GenerateChunk(xc, yc, zc)
	return world
}

func (self *World) GroundLevel(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(4*NOISE_SCALE), float64(z-MAP_DIAM)/(4*NOISE_SCALE), worldSeed, 1.4, 1.2, 4)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < -1.0 {
		noise = -1.0
	}
	if noise < 0.0 {
		noise /= 10
	}

	ground := uint16(SEA_LEVEL + ((CHUNK_HEIGHT-SEA_LEVEL)*0.9)*(noise+0.1)/1.1)
	return ground
}

func (self *World) SoilThickness(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(2*NOISE_SCALE), float64(z-MAP_DIAM)/(2*NOISE_SCALE), worldSeed, 1.8, 1.6, 8)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < -1.0 {
		noise = -1.0
	}

	return uint16(noise*2.5 + 2.5)
}

func (self *World) Precipitation(x uint16, z uint16) float64 {
	return perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE), float64(z-MAP_DIAM)/(NOISE_SCALE), worldSeed, 2.0, 0.6, 1)
}

// func (self *World) Drainage(x uint16, z uint16) float64 {
// 	return perlin.Noise2D(float64(x-MAP_DIAM)/(6*NOISE_SCALE), float64(z-MAP_DIAM)/(6*NOISE_SCALE), worldSeed, 0.4, 12)
// }

func (self *World) Rocks(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/2), float64(z-MAP_DIAM)/(NOISE_SCALE/2), worldSeed, 1.5, 3.0, 12)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < 0.8 {
		noise = 0
	}

	noise = noise * 5

	return uint16(noise)
}

func (self *World) Ore(x uint16, z uint16, blockid byte, occcurrence float64) uint16 {
	xloc := (float64(x) + MAP_DIAM*float64(blockid)) / (NOISE_SCALE / 2)
	zloc := (float64(z) + MAP_DIAM*float64(blockid)) / (NOISE_SCALE / 2)
	noise := perlin.Noise2D(xloc, zloc, worldSeed, 2.4, 1.8, 4)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < occcurrence {
		noise = 0
	} else {
		noise = 5 * (noise - occcurrence) / (1 - occcurrence)
	}
	return uint16(noise)
}

func (self *World) Coal(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/2), float64(z-MAP_DIAM)/(NOISE_SCALE/2), worldSeed, 1.9, 1.2, 12)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < 0.55 {
		noise = 0
	}

	noise = noise * 4

	return uint16(noise)
}

func (self *World) Iron(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/2.2), float64(z-MAP_DIAM)/(NOISE_SCALE/2.2), worldSeed, 2.5, 1.9, 6)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < 0.56 {
		noise = 0
	}

	noise = noise * 3

	return uint16(noise)
}

func (self *World) Copper(x uint16, z uint16) uint16 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/2.5), float64(z-MAP_DIAM)/(NOISE_SCALE/2.5), worldSeed, 3.1, 2.0, 6)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < 0.54 {
		noise = 0
	}

	noise = noise * 3

	return uint16(noise)
}

func (self *World) Feature1(x uint16, z uint16) float64 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/4), float64(z-MAP_DIAM)/(NOISE_SCALE/4), worldSeed, 6, 8, 14)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < -1.0 {
		noise = -1.0
	}
	return noise
}

func (self *World) Feature2(x uint16, z uint16) float64 {
	noise := perlin.Noise2D(float64(x-MAP_DIAM)/(NOISE_SCALE/4), float64(z-MAP_DIAM)/(NOISE_SCALE/4), worldSeed, 3, 5, 7)
	if noise > 1.0 {
		noise = 1.0
	}
	if noise < -1.0 {
		noise = -1.0
	}
	return noise
}

func (self *World) GenerateChunk(cx uint16, cy uint16, cz uint16) *Chunk {
	startTicks := time.Now().UnixNano()
	var chunk Chunk
	chunk.Init(cx, cy, cz)
	index := chunkIndex(cx, cy, cz)

	self.chunks[index] = &chunk

	amberfellCoords, hasAmberfell := self.amberfell[index]

	xw := cx * CHUNK_WIDTH
	zw := cz * CHUNK_WIDTH

	chunk.standingStoneProb = 0.0
	for x := uint16(0); x < CHUNK_WIDTH; x++ {
		for z := uint16(0); z < CHUNK_WIDTH; z++ {
			ground := self.GroundLevel(x+xw, z+zw)

			soil := ground + uint16(float64(self.SoilThickness(x+xw, z+zw))*(1-((float64(ground)-CHUNK_HEIGHT/2)/(CHUNK_HEIGHT/2))))
			rocks := ground // + self.Rocks(x+xw, z+zw)

			upper := ground

			if rocks > upper {
				upper = rocks
			}

			if soil > upper {
				upper = soil
			}

			if hasAmberfell && amberfellCoords[0] == x && amberfellCoords[1] == z {
				for y := uint16(0); y < upper; y++ {
					chunk.Set(x, y, z, BLOCK_AMBERFELL_SOURCE)
				}
			} else {
				for y := uint16(0); y < upper; y++ {
					if y >= rocks && y <= soil {
						chunk.Set(x, y, z, BLOCK_DIRT)
					} else {
						chunk.Set(x, y, z, BLOCK_STONE)
					}

				}

				for _, occurrence := range ORE_DISTRIBUTIONS {
					surface := upper - 1
					size := self.Ore(x+xw, z+zw, occurrence.itemid, occurrence.occurrence)
					if size > 0 {
						if size > 2 {
							surface++
						}
						for y := surface; y > surface-size && y > 0; y-- {
							chunk.Set(x, y, z, occurrence.itemid)
						}
						chunk.standingStoneProb += 0.000001 * occurrence.occurrence
						break
					}
				}

			}
		}
	}

	self.GenerateChunkFeatures(&chunk)

	console.chunkGenerationTime = time.Now().UnixNano() - startTicks
	if console.visible {
		println("Generating chunk at x:", cx, " y:", cy, " z:", cz, " in ", console.chunkGenerationTime/1e6)
	}
	return &chunk

}

func (self *World) GenerateChunkFeatures(chunk *Chunk) {
	if !chunk.featuresLoaded {
		cx := chunk.x
		cz := chunk.z
		if _, ok := self.chunks[chunkIndex(cx+1, 0, cz)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx-1, 0, cz)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx, 0, cz+1)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx, 0, cz-1)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx+1, 0, cz+1)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx+1, 0, cz-1)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx-1, 0, cz+1)]; !ok {
			return
		} else if _, ok := self.chunks[chunkIndex(cx-1, 0, cz-1)]; !ok {
			return
		}

		xw := cx * CHUNK_WIDTH
		zw := cz * CHUNK_WIDTH

		for x := uint16(xw); x < xw+CHUNK_WIDTH; x++ {
			for z := uint16(zw); z < zw+CHUNK_WIDTH; z++ {
				y := self.FindSurface(x, z)

				if rand.Float64() < chunk.standingStoneProb {
					// Place a standing stone
					chunk.standingStoneProb = 0
					chunk.Set(x-xw, y, z-zw, BLOCK_STONE)
					chunk.Set(x-xw, y+1, z-zw, BLOCK_STONE)
					chunk.Set(x-xw, y+2, z-zw, BLOCK_CARVED_STONE)
				}

				if self.Precipitation(x, z) > TREE_PRECIPITATION_MIN && rand.Intn(100) < TREE_DENSITY_PCT {

					if y > 1 && y < treeLine && chunk.At(x-xw, y-1, z-zw) == BLOCK_DIRT {
						self.GrowTree(x, y, z)

					}
				} else if self.Precipitation(x, z) > BUSH_PRECIPITATION_MIN && rand.Intn(100) < BUSH_DENSITY_PCT {
					self.Set(x, y, z, BLOCK_BUSH)
				} else {
					// feature1 := self.Feature1(x, z)
					// feature2 := self.Feature2(x, z)

					// if feature1 > 0.8 && feature2 > 0.8 {
					// 	self.Set(x, y, z, BLOCK_LEAVES)
					// 	self.Set(x, y+1, z, BLOCK_LEAVES)
					// 	self.Set(x, y+2, z, BLOCK_LEAVES)
					// 	self.Set(x, y+3, z, BLOCK_LEAVES)
					// 	self.Set(x, y+4, z, BLOCK_LEAVES)
					// 	self.Set(x, y+5, z, BLOCK_LEAVES)
					// 	self.Set(x, y+6, z, BLOCK_LEAVES)
					// } else if feature1 > 0.7 && feature2 > 0.7 {
					// 	self.Set(x, y, z, BLOCK_LEAVES)
					// 	self.Set(x, y+1, z, BLOCK_LEAVES)
					// 	self.Set(x, y+2, z, BLOCK_LEAVES)
					// 	self.Set(x, y+3, z, BLOCK_LEAVES)
					// 	self.Set(x, y+4, z, BLOCK_LEAVES)
					// } else if feature1 > 0.6 && feature2 > 0.6 {
					// 	self.Set(x, y, z, BLOCK_STONE)
					// 	self.Set(x, y+1, z, BLOCK_LEAVES)
					// 	self.Set(x, y+2, z, BLOCK_LEAVES)
					// }
				}
			}
		}

		chunk.featuresLoaded = true
	}
}

func (self *World) GenerateAmberfell() {
	for x := uint16(0); x < MAP_DIAM; x += 256 {
		for z := uint16(0); z < MAP_DIAM; z += 256 {
			pos := uint16(self.GenNext() >> 16)
			xw := x + pos%256
			zw := z + pos/256

			if xw > PLAYER_START_X-250 && xw < PLAYER_START_X+250 &&
				zw > PLAYER_START_Z-250 && zw < PLAYER_START_Z+250 {

				println("Amberfell: ", xw, zw)
			}

			index := chunkIndexFromWorld(xw, 0, zw)
			cx, _, cz := chunkCoordsFromindex(index)

			xo := xw - cx*CHUNK_WIDTH
			zo := zw - cz*CHUNK_WIDTH
			self.amberfell[index] = [2]uint16{xo, zo}

		}
	}

}

func (self *World) GenNext() int32 {
	self.genseed = (self.genseed*25214903917 + 11) % (1 << 48)

	return int32(self.genseed >> 16)
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

func (self *World) InvalidateRadius(x uint16, z uint16, r uint16) {
	chunks := make(map[uint64]bool, 10)

	for xc := x - r; xc < x+r; xc++ {
		for zc := z - r; zc < z+r; zc++ {
			chunks[chunkIndexFromWorld(xc, 0, zc)] = true
		}
	}

	for i := range chunks {
		chunk, ok := self.chunks[i]
		if ok {
			chunk.clean = false
		}
	}

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

func (self *World) HasVisibleFaces(neighbours [18]uint16) bool {

	for i := 0; i < 6; i++ {
		if items[neighbours[i]].transparent {
			return true
		}
	}

	return false
}

func (self *World) ApproxBlockAt(x uint16, y uint16, z uint16) uint16 {

	if self.ChunkLoadedFor(x, y, z) {
		return uint16(self.At(x, y, z))
	} else if self.GroundLevel(x, z) > y {
		return BLOCK_DIRT
	}
	return BLOCK_AIR
}

func (self *World) Neighbours(pos Vectori) (neighbours [18]uint16) {
	x := pos[XAXIS]
	y := pos[YAXIS]
	z := pos[ZAXIS]

	neighbours[WEST_FACE] = self.ApproxBlockAt(x-1, y, z)
	neighbours[EAST_FACE] = self.ApproxBlockAt(x+1, y, z)
	neighbours[NORTH_FACE] = self.ApproxBlockAt(x, y, z-1)
	neighbours[SOUTH_FACE] = self.ApproxBlockAt(x, y, z+1)
	neighbours[DOWN_FACE] = self.ApproxBlockAt(x, y-1, z)
	neighbours[UP_FACE] = self.ApproxBlockAt(x, y+1, z)

	// neighbours[DIR_NE] = self.ApproxBlockAt(x+1, y, z-1)
	// neighbours[DIR_SE] = self.ApproxBlockAt(x+1, y, z+1)
	// neighbours[DIR_SW] = self.ApproxBlockAt(x-1, y, z+1)
	// neighbours[DIR_NW] = self.ApproxBlockAt(x-1, y, z-1)

	// neighbours[DIR_UN]  = self.ApproxBlockAt(x  , y+1, z-1)
	// neighbours[DIR_UE]  = self.ApproxBlockAt(x+1, y+1, z)
	// neighbours[DIR_US]  = self.ApproxBlockAt(x  , y+1, z+1)
	// neighbours[DIR_UW]  = self.ApproxBlockAt(x-1, y+1, z)

	// neighbours[DIR_DN]  = self.ApproxBlockAt(x  , y-1, z-1)
	// neighbours[DIR_DE]  = self.ApproxBlockAt(x+1, y-1, z)
	// neighbours[DIR_DS]  = self.ApproxBlockAt(x  , y-1, z+1)
	// neighbours[DIR_DW]  = self.ApproxBlockAt(x-1, y-1, z)

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
		mob.Setvx(mob.Velocity()[XAXIS] / (1.0 + 2*dt))
		mob.Setvy(mob.Velocity()[YAXIS] - 18*dt)
		mob.Setvz(mob.Velocity()[ZAXIS] / (1.0 + 2*dt))
	} else {
		mob.Setvx(mob.Velocity()[XAXIS] / (1.0 + 12*dt))
		mob.Setvz(mob.Velocity()[ZAXIS] / (1.0 + 12*dt))
	}

	playerRect := Rect{x: float64(mp[XAXIS]) + mob.Velocity()[XAXIS]*dt, y: float64(mp[ZAXIS]) + mob.Velocity()[ZAXIS]*dt, sizex: mob.W(), sizey: mob.D()}

	// collisionCandidates := make([]Side, 0)

	if self.Atv(ip.North()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] < 0 && ip.North().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.North()[ZAXIS])+0.5+playerRect.sizey/2, 0)

			if self.Atv(ip.North().Up()) == BLOCK_AIR {
				mob.Setvy(4)
			}
		}
	}

	if self.Atv(ip.South()) != BLOCK_AIR {
		if mob.Velocity()[ZAXIS] > 0 && ip.South().HRect().Intersects(playerRect) {
			mob.Snapz(float64(ip.South()[ZAXIS])-0.5-playerRect.sizey/2, 0)
			if self.Atv(ip.South().Up()) == BLOCK_AIR {
				mob.Setvy(4)
			}
		}
	}

	if self.Atv(ip.East()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] > 0 && ip.East().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.East()[XAXIS])-0.5-playerRect.sizex/2, 0)
			if self.Atv(ip.East().Up()) == BLOCK_AIR {
				mob.Setvy(4)
			}
		}
	}

	if self.Atv(ip.West()) != BLOCK_AIR {
		if mob.Velocity()[XAXIS] < 0 && ip.West().HRect().Intersects(playerRect) {
			mob.Snapx(float64(ip.West()[XAXIS])+0.5+playerRect.sizex/2, 0)
			if self.Atv(ip.West().Up()) == BLOCK_AIR {
				mob.Setvy(4)
			}
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

func (self World) ChunkLoadedFor(x uint16, y uint16, z uint16) bool {
	cx := x / CHUNK_WIDTH
	cy := y / CHUNK_HEIGHT
	cz := z / CHUNK_WIDTH

	_, ok := self.chunks[chunkIndex(cx, cy, cz)]
	return ok
}

func (self *World) Draw(center Vectorf, selectedBlockFace *BlockFace) {
	for _, mob := range self.mobs {
		mob.Draw(center, selectedBlockFace)
	}

	pxmin, _, pzmin := chunkCoordsFromWorld(uint16(center[XAXIS]-float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]-float64(viewRadius)))
	pxmax, _, pzmax := chunkCoordsFromWorld(uint16(center[XAXIS]+float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]+float64(viewRadius)))

	for px := pxmin; px <= pxmax; px++ {
		for pz := pzmin; pz <= pzmax; pz++ {
			if chunk, ok := self.chunks[chunkIndex(px, 0, pz)]; ok {
				chunk.Render(selectedBlockFace)
			}
		}
	}

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
		if rand.Float64()*100 < BEESNEST_DENSITY_PCT {
			pos := Vectori{newx, newy, newz}
			nest := NewBeesNest(pos)
			TheWorld.timedObjects[pos] = nest
			TheWorld.containerObjects[pos] = nest

			self.Set(newx, newy, newz, BLOCK_BEESNEST)
		} else {
			self.Set(newx, newy, newz, BLOCK_LEAVES)

		}

	}
}

func (self *World) Simulate() {
	dt := float64(time.Now().UnixNano()-self.lastSimulated) / 1.e9
	self.lastSimulated = time.Now().UnixNano()

	self.ApplyForces(ThePlayer, dt)
	ThePlayer.Update(dt)

	for _, mob := range self.mobs {
		mob.Act(dt)
		self.ApplyForces(mob, dt)
		mob.Update(dt)
	}

	// Despawn
	for i := len(self.mobs) - 1; i >= 0; i-- {
		if ThePlayer.position.Minus(self.mobs[i].Position()).Magnitude() > float64(viewRadius)*3 {
			self.mobs = append(self.mobs[:i], self.mobs[i+1:]...)
		}
	}

	if len(self.mobs) < 10 {
		if rand.Float64() < 0.1*dt {
			angle := rand.Float64() * 2 * math.Pi
			distance := (1 + rand.Float64()) * float64(viewRadius)

			x := ThePlayer.position[XAXIS] + math.Cos(angle)*distance
			z := ThePlayer.position[ZAXIS] + -math.Sin(angle)*distance
			self.SpawnWolfPack(x, z)
		}
	}

	self.UpdateObjects(dt)

}

func (self *World) UpdateObjects(dt float64) {
	for key, obj := range self.timedObjects {
		if obj.Update(dt) {
			delete(self.timedObjects, key)
		}

	}
}

func (self *World) SpawnWolfPack(x float64, z float64) {

	size := rand.Intn(4) + rand.Intn(4)
	for i := 0; i < size; i++ {
		wx := uint16(x + rand.Float64()*8 - 4)
		wz := uint16(z + rand.Float64()*8 - 4)
		wolf := NewWolf(180, wx, self.FindSurface(wx, wz), wz)
		self.mobs = append(self.mobs, wolf)
	}

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
	if !self.featuresLoaded {
		TheWorld.GenerateChunkFeatures(self)
	}

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

					var blockid uint16 = uint16(self.Blocks[blockIndex(x, y, z)])
					if blockid != 0 {

						pos := Vectori{self.x*CHUNK_WIDTH + x, y, self.z*CHUNK_WIDTH + z}
						neighbours := TheWorld.Neighbours(pos)

						if TheWorld.HasVisibleFaces(neighbours) {

							selectedFace := uint8(FACE_NONE)
							if selectedBlockFace != nil && pos.Equals(&selectedBlockFace.pos) {
								selectedFace = selectedBlockFace.face
							}

							TerrainCube(self.vertexBuffer, pos, neighbours, blockid, selectedFace)
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

	self.vertexBuffer.RenderDirect(true)

	// fmt.Printf("Chunk ticks: %4.0f\n", float64(t.GetTicks())/1e6)

}
