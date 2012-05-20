/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"time"
)

type Chunk struct {
	x, z              uint16
	Blocks            [CHUNK_WIDTH * CHUNK_WIDTH * CHUNK_HEIGHT]byte
	vertexBuffer      *VertexBuffer
	clean             bool
	standingStoneProb float64
	featuresLoaded    bool
}

func NewChunk(cx uint16, cy uint16, cz uint16, hasAmberfell bool, amberfellCoords [2]uint16) *Chunk {
	startTicks := time.Now().UnixNano()
	var chunk Chunk
	chunk.Init(cx, cy, cz)

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
					size := TheWorld.Ore(x+xw, z+zw, occurrence.itemid, occurrence.occurrence)
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

	TheWorld.GenerateChunkFeatures(&chunk)

	console.chunkGenerationTime = time.Now().UnixNano() - startTicks
	if console.visible {
		println("Generating chunk at x:", cx, " y:", cy, " z:", cz, " in ", console.chunkGenerationTime/1e6)
	}
	return &chunk
}

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
