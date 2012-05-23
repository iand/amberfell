/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"math"
	"time"
)

const (
	BLOCK_TYPE_MASK = 0x0F
)

type chunkindex uint32

type Chunk struct {
	x, z              uint16
	Blocks            [CHUNK_WIDTH * CHUNK_WIDTH * CHUNK_HEIGHT]Block
	vertexBuffer      *VertexBuffer
	clean             bool
	standingStoneProb float64
	featuresLoaded    bool
}

type Block struct {
	id   BlockId
	data uint8
}

func NewBlock(id BlockId, damaged bool, orientation uint8) Block {
	var data uint8

	if damaged {
		data |= 0x1
	}

	data |= (orientation << 1)
	return Block{id, data}
}

func NewBlockDefault(id BlockId) Block {
	return Block{id, 0}
}

func (self *Block) Damaged() bool {
	return (self.data & 0x1) == 0x1
}

func (self *Block) SetDamaged(damaged bool) {
	self.data &= 0xE
	if damaged {
		self.data |= 0x1
	}
}

func (self *Block) Orientation() uint8 {
	return (self.data & 0x6) >> 1
}

func chunkCoordsFromWorld(x uint16, z uint16) (cx uint16, cz uint16) {
	cx = uint16(math.Floor(float64(x) / CHUNK_WIDTH))
	cz = uint16(math.Floor(float64(z) / CHUNK_WIDTH))

	return
}

func chunkIndexFromWorld(x uint16, z uint16) chunkindex {
	cx, cz := chunkCoordsFromWorld(x, z)
	return chunkIndex(cx, cz)
}

func chunkIndex(cx uint16, cz uint16) chunkindex {
	return chunkindex(cz)<<16 | chunkindex(cx)
}

func chunkCoordsFromindex(index chunkindex) (cx uint16, cz uint16) {
	cx = uint16(index)
	cz = uint16(index >> 16)

	return
}

func blockIndex(x uint16, y uint16, z uint16) uint16 {
	return CHUNK_WIDTH*CHUNK_WIDTH*y + CHUNK_WIDTH*x + z
}

func NewChunk(cx uint16, cz uint16, hasAmberfell bool, amberfellCoords [2]uint16) *Chunk {
	startTicks := time.Now().UnixNano()
	if console.visible {
		println("Generating chunk at x:", cx, " z:", cz)
	}
	var chunk Chunk
	chunk.Init(cx, cz)

	xw := cx * CHUNK_WIDTH
	zw := cz * CHUNK_WIDTH

	chunk.standingStoneProb = 0.0
	for x := uint16(0); x < CHUNK_WIDTH; x++ {
		for z := uint16(0); z < CHUNK_WIDTH; z++ {
			ground := TheWorld.GroundLevel(x+xw, z+zw)

			soil := ground + uint16(float64(TheWorld.SoilThickness(x+xw, z+zw))*(1-((float64(ground)-CHUNK_HEIGHT/2)/(CHUNK_HEIGHT/2))))
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
					size := TheWorld.Ore(x+xw, z+zw, BlockId(occurrence.itemid), occurrence.occurrence)
					if size > 0 {
						if size > 2 {
							surface++
						}
						for y := surface; y > surface-size && y > 0; y-- {
							chunk.Set(x, y, z, BlockId(occurrence.itemid))
						}
						chunk.standingStoneProb += 0.000001 * occurrence.occurrence
						break
					}
				}

			}
		}
	}

	//TheWorld.GenerateChunkFeatures(&chunk)

	console.chunkGenerationTime = time.Now().UnixNano() - startTicks
	return &chunk
}

func (c Chunk) WorldCoords(x uint16, y uint16, z uint16) (xw uint16, yw uint16, zw uint16) {
	xw = c.x*CHUNK_WIDTH + x
	zw = c.z*CHUNK_WIDTH + z
	yw = y
	return
}

func (chunk *Chunk) Init(x uint16, z uint16) {
	chunk.x = x
	// chunk.y = y
	chunk.z = z
	chunk.vertexBuffer = NewVertexBuffer(VERTEX_BUFFER_CAPACITY, terrainTexture)
}

func (chunk *Chunk) At(x uint16, y uint16, z uint16) BlockId {
	return chunk.Blocks[blockIndex(x, y, z)].id
}

func (chunk *Chunk) AtB(x uint16, y uint16, z uint16) Block {
	return chunk.Blocks[blockIndex(x, y, z)]
}

func (chunk *Chunk) Set(x uint16, y uint16, z uint16, b BlockId) {
	chunk.SetB(x, y, z, NewBlock(b, false, ORIENT_EAST))
	chunk.clean = false
}

func (chunk *Chunk) SetB(x uint16, y uint16, z uint16, block Block) {
	chunk.Blocks[blockIndex(x, y, z)] = block
	chunk.clean = false
}

// func (self *Chunk) PreRender(selectedBlockFace *BlockFace) {
// 	if !self.featuresLoaded {
// 		// TheWorld.GenerateChunkFeatures(self)
// 	}

// 	if !self.clean || (selectedBlockFace != nil && selectedBlockFace.pos[XAXIS] >= self.x*CHUNK_WIDTH && selectedBlockFace.pos[XAXIS] < (self.x+1)*CHUNK_WIDTH &&
// 		/*selectedBlockFace.pos[YAXIS] >= self.y*CHUNK_HEIGHT && selectedBlockFace.pos[YAXIS] < (self.y+1)*CHUNK_HEIGHT && */
// 		selectedBlockFace.pos[ZAXIS] >= self.z*CHUNK_WIDTH && selectedBlockFace.pos[ZAXIS] < (self.z+1)*CHUNK_WIDTH) {
// 		self.vertexBuffer.Reset()
// 		var x, y, z uint16
// 		for x = 0; x < CHUNK_WIDTH; x++ {
// 			for z = 0; z < CHUNK_WIDTH; z++ {
// 				for y = 0; y < CHUNK_HEIGHT; y++ {
// 					block := self.Blocks[blockIndex(x, y, z)]
// 					if block.id != 0 {

// 						pos := Vectori{self.x*CHUNK_WIDTH + x, y, self.z*CHUNK_WIDTH + z}
// 						neighbours := TheWorld.Neighbours(pos)

// 						if HasVisibleFaces(neighbours) {

// 							selectedFace := uint8(FACE_NONE)
// 							if selectedBlockFace != nil && pos.Equals(&selectedBlockFace.pos) {
// 								selectedFace = selectedBlockFace.face
// 							}

// 							TerrainCube(self.vertexBuffer, pos, neighbours, block, selectedFace)
// 						}
// 					}
// 				}
// 			}
// 		}

// 		self.clean = true
// 	}

// }

// func (self *Chunk) Render(selectedBlockFace *BlockFace) {
// 	// self.PreRender(selectedBlockFace)

// 	self.vertexBuffer.RenderDirect(true)

// 	// fmt.Printf("Chunk ticks: %4.0f\n", float64(t.GetTicks())/1e6)

// }

func (self *Chunk) Render(adjacents [4]*Chunk, selectedBlockFace *BlockFace, vb chan *VertexBuffer) {
	if !self.featuresLoaded {

		allAdjacentsAvailable := true
		for _, ac := range adjacents {
			if ac == nil {
				allAdjacentsAvailable = false
			}
		}

		if allAdjacentsAvailable {
			TheWorld.GenerateChunkFeatures(self, adjacents)
		}

	}

	selectionInThisChunk := selectedBlockFace != nil && selectedBlockFace.pos[XAXIS] >= self.x*CHUNK_WIDTH && selectedBlockFace.pos[XAXIS] < (self.x+1)*CHUNK_WIDTH &&
		selectedBlockFace.pos[ZAXIS] >= self.z*CHUNK_WIDTH && selectedBlockFace.pos[ZAXIS] < (self.z+1)*CHUNK_WIDTH

	if !self.clean || selectionInThisChunk {

		self.vertexBuffer.Reset()
		var x, y, z uint16
		for x = 0; x < CHUNK_WIDTH; x++ {
			for z = 0; z < CHUNK_WIDTH; z++ {
				for y = 0; y < CHUNK_HEIGHT; y++ {
					block := self.Blocks[blockIndex(x, y, z)]
					if block.id != 0 {

						pos := Vectori{self.x*CHUNK_WIDTH + x, y, self.z*CHUNK_WIDTH + z}

						var neighbours [18]BlockId

						if x > 0 {
							neighbours[WEST_FACE] = self.At(x-1, y, z)
						} else if adjacents[WEST_FACE] != nil {
							neighbours[WEST_FACE] = adjacents[WEST_FACE].At(CHUNK_WIDTH-1, y, z)
						} else {
							neighbours[WEST_FACE] = BLOCK_STONE
						}

						if x < CHUNK_WIDTH-1 {
							neighbours[EAST_FACE] = self.At(x+1, y, z)
						} else if adjacents[EAST_FACE] != nil {
							neighbours[EAST_FACE] = adjacents[EAST_FACE].At(0, y, z)
						} else {
							neighbours[EAST_FACE] = BLOCK_STONE
						}

						if z > 0 {
							neighbours[NORTH_FACE] = self.At(x, y, z-1)
						} else if adjacents[NORTH_FACE] != nil {
							neighbours[NORTH_FACE] = adjacents[NORTH_FACE].At(x, y, CHUNK_WIDTH-1)
						} else {
							neighbours[NORTH_FACE] = BLOCK_STONE
						}

						if z < CHUNK_WIDTH-1 {
							neighbours[SOUTH_FACE] = self.At(x, y, z+1)
						} else if adjacents[SOUTH_FACE] != nil {
							neighbours[SOUTH_FACE] = adjacents[SOUTH_FACE].At(x, y, 0)
						} else {
							neighbours[SOUTH_FACE] = BLOCK_STONE
						}

						if y > 0 {
							neighbours[DOWN_FACE] = self.At(x, y-1, z)
						} else {
							neighbours[DOWN_FACE] = BLOCK_STONE
						}

						if y < CHUNK_HEIGHT-1 {
							neighbours[UP_FACE] = self.At(x, y+1, z)
						} else {
							neighbours[UP_FACE] = BLOCK_AIR
						}

						if HasVisibleFaces(neighbours) {

							selectedFace := uint8(FACE_NONE)
							if selectedBlockFace != nil && pos.Equals(&selectedBlockFace.pos) {
								selectedFace = selectedBlockFace.face
							}

							TerrainCube(self.vertexBuffer, pos, neighbours, block, selectedFace)
						}
					}
				}
			}
		}

		self.clean = true
	}

	if vb != nil {
		vb <- self.vertexBuffer
	}

}
