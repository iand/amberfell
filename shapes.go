/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	//	"fmt"
	"github.com/banthar/gl"
	"image"
	_ "image/png"
	"os"
)

type Vertex struct {
	p [3]float32 // Position
	t [2]float32 // Texture coordinate
	n [3]float32 // Normal
	c [4]float32 // Colour
	e [4]float32 // Emissive light
}

type TriangleIndex struct {
	i1 uint32
	i2 uint32
	i3 uint32
}

type VertexBuffer struct {
	vertices    []Vertex
	indices     []TriangleIndex
	vertexCount int
	indexCount  int
	texture     *gl.Texture
}

func NewVertexBuffer(capacity uint32, texture *gl.Texture) *VertexBuffer {
	var v VertexBuffer
	v.vertices = make([]Vertex, capacity, capacity)
	v.indices = make([]TriangleIndex, capacity, capacity)
	v.texture = texture
	return &v
}

func (self *VertexBuffer) Reset() {
	self.vertexCount = 0
	self.indexCount = 0
}

func (self *VertexBuffer) AddFace(face uint8, texture uint16, selected bool, shade int, x1, y1, z1, tx1, ty1, x2, y2, z2, tx2, ty2 float32) {
	if self.vertexCount >= VERTEX_BUFFER_CAPACITY-4 {
		// TODO: log a warning about overflowing buffer
		return
	}

	c := COLOURS[shade]
	if selected {
		c = COLOUR_HIGH
	}

	pos := Vectorf{float64(x1), float64(y1), float64(z1)}

	e := LightLevel(pos, NORMALS[face])

	vc := self.vertexCount

	if x1 == x2 {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+2] = Vertex{p: [3]float32{x1, y2, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+3] = Vertex{p: [3]float32{x1, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
	} else if y1 == y2 {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+2] = Vertex{p: [3]float32{x2, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+3] = Vertex{p: [3]float32{x2, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
	} else {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+2] = Vertex{p: [3]float32{x2, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
		self.vertices[vc+3] = Vertex{p: [3]float32{x2, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c, e: e}
	}

	self.vertexCount += 4
	ic := self.indexCount
	self.indices[ic] = TriangleIndex{uint32(vc), uint32(vc) + 1, uint32(vc) + 2}
	self.indices[ic+1] = TriangleIndex{uint32(vc) + 2, uint32(vc) + 3, uint32(vc)}
	self.indexCount += 2
}

func (self *VertexBuffer) RenderDirect(clip bool) {
	self.texture.Bind(gl.TEXTURE_2D)

	if clip {
		cutoff := float32(-0.00006) // magic!

		planes32 := viewport.ClipPlanes()

		gl.Begin(gl.QUADS)
		for i := 0; i < self.vertexCount; i += 4 {
			draw := true
			// Cull back faces
			// Use dot product of front viewport plane with normal of face
			dot := planes32[5][0]*self.vertices[i].n[0] + planes32[5][1]*self.vertices[i].n[1] + planes32[5][2]*self.vertices[i].n[2]
			if dot <= 0 {
				draw = false
			} else {
				for p := 0; p < 6; p++ {
					dist1 := planes32[p][0]*self.vertices[i].p[0] + planes32[p][1]*self.vertices[i].p[1] + planes32[p][2]*self.vertices[i+1].p[2] + planes32[p][3]
					dist2 := planes32[p][0]*self.vertices[i+1].p[0] + planes32[p][1]*self.vertices[i+2].p[1] + planes32[p][2]*self.vertices[i+2].p[2] + planes32[p][3]
					dist3 := planes32[p][0]*self.vertices[i+2].p[0] + planes32[p][1]*self.vertices[i+3].p[1] + planes32[p][2]*self.vertices[i+3].p[2] + planes32[p][3]
					dist4 := planes32[p][0]*self.vertices[i+3].p[0] + planes32[p][1]*self.vertices[i+4].p[1] + planes32[p][2]*self.vertices[i+4].p[2] + planes32[p][3]
					if dist1 <= cutoff && dist2 <= cutoff && dist3 <= cutoff && dist4 <= cutoff {
						draw = false
						break
					}
				}
			}

			if draw {
				for j := i; j < i+4; j++ {
					gl.Normal3f(self.vertices[j].n[0], self.vertices[j].n[1], self.vertices[j].n[2])
					gl.TexCoord2f(self.vertices[j].t[0], self.vertices[j].t[1])
					gl.Color4f(self.vertices[j].c[0], self.vertices[j].c[1], self.vertices[j].c[2], self.vertices[j].c[3])

					gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{self.vertices[j].e[0], self.vertices[j].e[1], self.vertices[j].e[2], self.vertices[j].e[3]})
					gl.Vertex3f(self.vertices[j].p[0], self.vertices[j].p[1], self.vertices[j].p[2])
				}
				console.vertices += 4
			} else {
				console.culledVertices += 4
			}
		}
		gl.End()
	} else {
		gl.Begin(gl.QUADS)
		for i := 0; i < self.vertexCount; i++ {
			gl.Normal3f(self.vertices[i].n[0], self.vertices[i].n[1], self.vertices[i].n[2])
			gl.TexCoord2f(self.vertices[i].t[0], self.vertices[i].t[1])
			gl.Color4f(self.vertices[i].c[0], self.vertices[i].c[1], self.vertices[i].c[2], self.vertices[i].c[3])
			gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{self.vertices[i].e[0], self.vertices[i].e[1], self.vertices[i].e[2], self.vertices[i].e[3]})
			gl.Vertex3f(self.vertices[i].p[0], self.vertices[i].p[1], self.vertices[i].p[2])
		}
		gl.End()

	}

	self.texture.Unbind(gl.TEXTURE_2D)
}

// func LoadMapTextures() {

// 	var file, err = os.Open("tiles.png")
// 	var img image.Image
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	if img, _, err = image.Decode(file); err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < 10; i++ {
// 		for j := 0; j < 16; j++ {
// 			textureIndex := uint16(i*16 + j)
// 			textures[textureIndex] = imageSectionToTexture(img, image.Rect(TILE_WIDTH*j, TILE_WIDTH*i, TILE_WIDTH*j+TILE_WIDTH, TILE_WIDTH*i+TILE_WIDTH))
// 		}
// 	}
// }

func LoadPlayerTextures() {

	var file, err = os.Open("res/player.png")
	var img image.Image
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if img, _, err = image.Decode(file); err != nil {
		panic(err)
	}

	unit := 12
	hatFront := image.NewRGBA(image.Rect(0, 0, unit, unit))
	for x := 0; x < unit; x++ {
		for y := 0; y < unit; y++ {
			hatFront.Set(x, y, img.At(x, y))
		}
	}
	textures[TEXTURE_HAT_FRONT] = imageSectionToTexture(img, image.Rect(0, 0, unit, unit))
	textures[TEXTURE_HAT_LEFT] = imageSectionToTexture(img, image.Rect(unit+1, 0, 2*unit, unit))
	textures[TEXTURE_HAT_BACK] = imageSectionToTexture(img, image.Rect(2*unit+1, 0, 3*unit, unit))
	textures[TEXTURE_HAT_RIGHT] = imageSectionToTexture(img, image.Rect(3*unit+1, 0, 4*unit, unit))
	textures[TEXTURE_HAT_TOP] = imageSectionToTexture(img, image.Rect(4*unit+1, 0, 5*unit, unit))

	textures[TEXTURE_HEAD_FRONT] = imageSectionToTexture(img, image.Rect(0, unit+1, unit, 2*unit))
	textures[TEXTURE_HEAD_LEFT] = imageSectionToTexture(img, image.Rect(unit+1, unit+1, 2*unit, 2*unit))
	textures[TEXTURE_HEAD_BACK] = imageSectionToTexture(img, image.Rect(2*unit+1, unit+1, 3*unit, 2*unit))
	textures[TEXTURE_HEAD_RIGHT] = imageSectionToTexture(img, image.Rect(3*unit+1, unit+1, 4*unit, 2*unit))
	textures[TEXTURE_HEAD_BOTTOM] = imageSectionToTexture(img, image.Rect(4*unit+1, unit+1, 5*unit, 2*unit))

	textures[TEXTURE_TORSO_FRONT] = imageSectionToTexture(img, image.Rect(0, 2*unit+1, 2*unit, 5*unit+unit/4))
	textures[TEXTURE_TORSO_LEFT] = imageSectionToTexture(img, image.Rect(2*unit+1, 2*unit+1, 3*unit, 5*unit+unit/4))
	textures[TEXTURE_TORSO_BACK] = imageSectionToTexture(img, image.Rect(3*unit+1, 2*unit+1, 5*unit, 5*unit+unit/4))
	textures[TEXTURE_TORSO_RIGHT] = imageSectionToTexture(img, image.Rect(5*unit+1, 2*unit+1, 6*unit, 5*unit+unit/4))
	textures[TEXTURE_TORSO_TOP] = imageSectionToTexture(img, image.Rect(32, 64, 55, 75))

	textures[TEXTURE_LEG] = imageSectionToTexture(img, image.Rect(0, 64, 11, 105))
	textures[TEXTURE_LEG_SIDE] = imageSectionToTexture(img, image.Rect(12, 64, 22, 105))
	textures[TEXTURE_ARM] = imageSectionToTexture(img, image.Rect(23, 57, 31, 96))
	textures[TEXTURE_ARM_TOP] = imageSectionToTexture(img, image.Rect(56, 64, 64, 72))
	textures[TEXTURE_BRIM] = imageSectionToTexture(img, image.Rect(31, 76, 49, 78))
	textures[TEXTURE_HAND] = imageSectionToTexture(img, image.Rect(23, 97, 31, 105))

}

func imageSectionToTexture(img image.Image, r image.Rectangle) *gl.Texture {
	rgba := image.NewRGBA(image.Rect(0, 0, r.Max.X-r.Min.X, r.Max.Y-r.Min.Y))
	for x := r.Min.X; x < r.Max.X+1; x++ {
		for y := r.Min.Y; y < r.Max.Y+1; y++ {
			rgba.Set(x-r.Min.X, y-r.Min.Y, img.At(x, y))
		}
	}

	return imageToTexture(rgba)
}

func imageToTexture(rgba *image.RGBA) *gl.Texture {
	rect := rgba.Bounds()
	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rect.Max.X, rect.Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
	texture.Unbind(gl.TEXTURE_2D)

	return &texture
}

func loadTexture(filename string) *gl.Texture {

	var img image.Image
	var file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if img, _, err = image.Decode(file); err != nil {
		panic(err)
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	for x := 0; x < rect.Max.X; x++ {
		for y := 0; y < rect.Max.Y; y++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, rect.Max.X, rect.Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
	texture.Unbind(gl.TEXTURE_2D)

	return &texture

}

func TerrainCube(vertexBuffer *VertexBuffer, x float32, y float32, z float32, neighbours [18]uint16, blockid uint16, selectedFace uint8) {

	block := items[uint16(blockid)]
	var visible [6]bool
	var shadeLevels [6]int

	for i := 0; i < 6; i++ {
		if items[neighbours[i]].transparent {
			visible[i] = true
		}
	}

	if items[neighbours[EAST_FACE]].transparent {
		visible[EAST_FACE] = true
		if !items[neighbours[DIR_NE]].transparent {
			shadeLevels[EAST_FACE]++
		}
		if !items[neighbours[DIR_SE]].transparent {
			shadeLevels[EAST_FACE]++
		}
		if !items[neighbours[DIR_UE]].transparent {
			shadeLevels[EAST_FACE]++
		}
		if !items[neighbours[DIR_DE]].transparent {
			shadeLevels[EAST_FACE]++
		}
	}

	if items[neighbours[WEST_FACE]].transparent {
		visible[WEST_FACE] = true
		if !items[neighbours[DIR_NW]].transparent {
			shadeLevels[WEST_FACE]++
		}
		if !items[neighbours[DIR_SW]].transparent {
			shadeLevels[WEST_FACE]++
		}
		if !items[neighbours[DIR_UW]].transparent {
			shadeLevels[WEST_FACE]++
		}
		if !items[neighbours[DIR_DW]].transparent {
			shadeLevels[WEST_FACE]++
		}
	}

	if items[neighbours[NORTH_FACE]].transparent {
		visible[NORTH_FACE] = true
		if !items[neighbours[DIR_NW]].transparent {
			shadeLevels[NORTH_FACE]++
		}
		if !items[neighbours[DIR_NE]].transparent {
			shadeLevels[NORTH_FACE]++
		}
		if !items[neighbours[DIR_UN]].transparent {
			shadeLevels[NORTH_FACE]++
		}
		if !items[neighbours[DIR_DN]].transparent {
			shadeLevels[NORTH_FACE]++
		}
	}

	if items[neighbours[SOUTH_FACE]].transparent {
		visible[SOUTH_FACE] = true
		if !items[neighbours[DIR_SW]].transparent {
			shadeLevels[SOUTH_FACE]++
		}
		if !items[neighbours[DIR_SE]].transparent {
			shadeLevels[SOUTH_FACE]++
		}
		if !items[neighbours[DIR_US]].transparent {
			shadeLevels[SOUTH_FACE]++
		}
		if !items[neighbours[DIR_DS]].transparent {
			shadeLevels[SOUTH_FACE]++
		}
	}

	if items[neighbours[UP_FACE]].transparent {
		visible[UP_FACE] = true
		if !items[neighbours[DIR_UN]].transparent {
			shadeLevels[UP_FACE]++
		}
		if !items[neighbours[DIR_UE]].transparent {
			shadeLevels[UP_FACE]++
		}
		if !items[neighbours[DIR_US]].transparent {
			shadeLevels[UP_FACE]++
		}
		if !items[neighbours[DIR_UW]].transparent {
			shadeLevels[UP_FACE]++
		}
	}

	if items[neighbours[DOWN_FACE]].transparent {
		visible[DOWN_FACE] = true
		if !items[neighbours[DIR_UN]].transparent {
			shadeLevels[DOWN_FACE]++
		}
		if !items[neighbours[DIR_UE]].transparent {
			shadeLevels[DOWN_FACE]++
		}
		if !items[neighbours[DIR_US]].transparent {
			shadeLevels[DOWN_FACE]++
		}
		if !items[neighbours[DIR_UW]].transparent {
			shadeLevels[DOWN_FACE]++
		}
	}

	switch blockid {
	case BLOCK_CAMPFIRE:
		Pile(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
	case BLOCK_LOG_SLAB:
		if neighbours[NORTH_FACE] != BLOCK_AIR {
			if neighbours[EAST_FACE] != BLOCK_AIR {
				if neighbours[SOUTH_FACE] != BLOCK_AIR {
					if neighbours[WEST_FACE] != BLOCK_AIR {
						// Blocks to all four sides
						SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
					} else {
						// Blocks to north, east, south
						SlabTee(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, shadeLevels, selectedFace)
					}
				} else if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, east, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to north, east
					SlabCorner(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
				}

			} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, south, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to north, south
					SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to the north and west
				SlabCorner(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)

			} else {
				// Just a block to the north
				SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
			}

		} else if neighbours[EAST_FACE] != BLOCK_AIR {
			if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to east, south, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to east, south
					SlabCorner(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, shadeLevels, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to east, west
				SlabLine(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
			} else {
				// Just a block to the east
				SlabLine(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
			}
		} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
			if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to south, west
				SlabCorner(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, shadeLevels, selectedFace)
			} else {
				// Just a block to the south
				SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
			}
		} else if neighbours[WEST_FACE] != BLOCK_AIR {
			// Just a block to the west
			SlabLine(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, shadeLevels, selectedFace)
		} else {
			// Lone block
			SlabSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
		}

	case BLOCK_LOG_WALL:
		if neighbours[NORTH_FACE] != BLOCK_AIR {
			if neighbours[EAST_FACE] != BLOCK_AIR {
				if neighbours[SOUTH_FACE] != BLOCK_AIR {
					if neighbours[WEST_FACE] != BLOCK_AIR {
						// Blocks to all four sides
						WallCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
					} else {
						// Blocks to north, east, south
						WallTee(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, shadeLevels, selectedFace)
					}
				} else if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, east, west
					WallTee(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to north, east
					WallCorner(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
				}

			} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, south, west
					WallTee(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to north, south
					WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to the north and west
				WallCorner(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)

			} else {
				// Just a block to the north
				WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
			}

		} else if neighbours[EAST_FACE] != BLOCK_AIR {
			if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to east, south, west
					WallTee(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, shadeLevels, selectedFace)
				} else {
					// Blocks to east, south
					WallCorner(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, shadeLevels, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to east, west
				WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
			} else {
				// Just a block to the east
				WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
			}
		} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
			if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to south, west
				WallCorner(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, shadeLevels, selectedFace)
			} else {
				// Just a block to the south
				WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, shadeLevels, selectedFace)
			}
		} else if neighbours[WEST_FACE] != BLOCK_AIR {
			// Just a block to the west
			WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
		} else {
			// Lone block
			WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)
		}

	default:
		Cuboid2(vertexBuffer, x, y, z, 1, 1, 1, block, visible, shadeLevels, selectedFace)

	}

}

func RenderQuads(v []Vertex) {
	gl.Begin(gl.QUADS)
	for i := 0; i < len(v); i++ {
		gl.Normal3f(v[i].n[0], v[i].n[1], v[i].n[2])
		gl.TexCoord2f(v[i].t[0], v[i].t[1])
		gl.Color4f(v[i].c[0], v[i].c[1], v[i].c[2], v[i].c[3])
		gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{v[i].e[0], v[i].e[1], v[i].e[2], v[i].e[3]})
		gl.Vertex3f(v[i].p[0], v[i].p[1], v[i].p[2])
	}
	gl.End()
}

func Cuboid2(vertexBuffer *VertexBuffer, x float32, y float32, z float32, bw float64, bh float64, bd float64, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2

	// East face
	if visible[EAST_FACE] {

		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+d, y+h, z-w, 1.0, 1.0,
			x+d, y-h, z+w, 0.0, 0.0)

	}

	// West Face
	if visible[WEST_FACE] {

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x-d, y+h, z-w, 1.0, 1.0,
			x-d, y-h, z+w, 0.0, 0.0)
	}

	// North Face
	if visible[NORTH_FACE] {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+d, y+h, z-w, 1.0, 1.0,
			x-d, y-h, z-w, 0.0, 0.0)
	}

	// South Face
	if visible[SOUTH_FACE] {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+d, y+h, z+w, 1.0, 1.0,
			x-d, y-h, z+w, 0.0, 0.0)

	}

	// Up Face
	if visible[UP_FACE] {

		//  -d/-w   -d/0   -d/+w
		//
		//  0/-w     0/0     0/+w
		//
		//  +d/-w   +d/0   +d/+w

		// +d/-w     0/-w   -d/-w
		// +d/0      0/0    -d/0
		// +d/+w     0/+w   -d/+w

		// Texture
		// 0.0/1.0    0.0/0.5   0.0/0.0
		// 0.5/1.0    0.5/0.5   0.5/0.0
		// 1.0/1.0    1.0/0.5   1.0/0.0

		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+d, y+h, z-w, 1.0, 1.0,
			x-d, y+h, z+w, 0.0, 0.0)

		// vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		// 	x+d, y+h, z+0, 1.0, 0.5,
		// 	x+0, y+h, z+w, 0.5, 0.0)

		// vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		// 	x+0, y+h, z+0, 0.5, 0.5,
		// 	x-d, y+h, z+w, 0.0, 0.0)

		// vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		// 	x+0, y+h, z-w, 0.5, 1.0,
		// 	x-d, y+h, z+0, 0.0, 0.5)
	}

	// Down Face
	if visible[DOWN_FACE] {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+d, y-h, z-w, 1.0, 1.0,
			x-d, y-h, z+w, 0.0, 0.0)
	}

}

func Cuboid(pos Vectorf, bw float64, bh float64, bd float64, etexture *gl.Texture, wtexture *gl.Texture, ntexture *gl.Texture, stexture *gl.Texture, utexture *gl.Texture, dtexture *gl.Texture, selectedFace uint8) {

	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2

	// East face
	if etexture != nil {

		c := COLOUR_WHITE
		if selectedFace == EAST_FACE {
			c = COLOUR_HIGH
		}

		e := LightLevel(pos, NORMALS[EAST_FACE])

		v := []Vertex{
			{p: [3]float32{d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[EAST_FACE], c: c, e: e},
			{p: [3]float32{d, h, -w}, t: [2]float32{1.0, 0.0}, n: NORMALS[EAST_FACE], c: c, e: e},
			{p: [3]float32{d, h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[EAST_FACE], c: c, e: e},
			{p: [3]float32{d, -h, w}, t: [2]float32{0.0, 1.0}, n: NORMALS[EAST_FACE], c: c, e: e},
		}

		etexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		etexture.Unbind(gl.TEXTURE_2D)
	}

	// West Face
	if wtexture != nil {
		c := COLOUR_WHITE
		if selectedFace == WEST_FACE {
			c = COLOUR_HIGH
		}

		e := LightLevel(pos, NORMALS[WEST_FACE])
		v := []Vertex{
			{p: [3]float32{-d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[WEST_FACE], c: c, e: e},
			{p: [3]float32{-d, -h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[WEST_FACE], c: c, e: e},
			{p: [3]float32{-d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[WEST_FACE], c: c, e: e},
			{p: [3]float32{-d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[WEST_FACE], c: c, e: e},
		}

		wtexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		wtexture.Unbind(gl.TEXTURE_2D)

	}

	// North Face
	if ntexture != nil {
		c := COLOUR_WHITE
		if selectedFace == NORTH_FACE {
			c = COLOUR_HIGH
		}

		e := LightLevel(pos, NORMALS[NORTH_FACE])
		v := []Vertex{
			{p: [3]float32{-d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[NORTH_FACE], c: c, e: e},
			{p: [3]float32{-d, h, -w}, t: [2]float32{1.0, 0.0}, n: NORMALS[NORTH_FACE], c: c, e: e},
			{p: [3]float32{d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[NORTH_FACE], c: c, e: e},
			{p: [3]float32{d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[NORTH_FACE], c: c, e: e},
		}

		ntexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		ntexture.Unbind(gl.TEXTURE_2D)

	}

	// South Face
	if stexture != nil {
		c := COLOUR_WHITE
		if selectedFace == SOUTH_FACE {
			c = COLOUR_HIGH
		}

		e := LightLevel(pos, NORMALS[SOUTH_FACE])
		v := []Vertex{
			{p: [3]float32{-d, -h, w}, t: [2]float32{0.0, 1.0}, n: NORMALS[SOUTH_FACE], c: c, e: e},
			{p: [3]float32{d, -h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[SOUTH_FACE], c: c, e: e},
			{p: [3]float32{d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[SOUTH_FACE], c: c, e: e},
			{p: [3]float32{-d, h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[SOUTH_FACE], c: c, e: e},
		}

		stexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		stexture.Unbind(gl.TEXTURE_2D)
	}

	// Up Face
	if utexture != nil {

		c := COLOUR_WHITE
		if selectedFace == UP_FACE {
			c = COLOUR_HIGH
		}
		e := LightLevel(pos, NORMALS[UP_FACE])

		//  -d/-w   -d/0   -d/+w
		//
		//  0/-w     0/0     0/+w
		//
		//  +d/-w   +d/0   +d/+w

		// +d/-w     0/-w   -d/-w
		// +d/0      0/0    -d/0
		// +d/+w     0/+w   -d/+w

		// Texture
		// 0.0/1.0    0.0/0.5   0.0/0.0
		// 0.5/1.0    0.5/0.5   0.5/0.0
		// 1.0/1.0    1.0/0.5   1.0/0.0

		v := []Vertex{
			{p: [3]float32{-d, h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},
			{p: [3]float32{d, h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},
			{p: [3]float32{d, h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},
			{p: [3]float32{-d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},

			// Vertex{p: [3]float32{d, h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, -w}, t: [2]float32{0.0, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{d, h, 0}, t: [2]float32{0.5, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},

			// Vertex{p: [3]float32{d, h, 0}, t: [2]float32{0.5, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, w}, t: [2]float32{1.0, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{d, h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[UP_FACE], c: c, e: e},

			// Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{-d, h, 0}, t: [2]float32{0.5, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{-d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, w}, t: [2]float32{1.0, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},

			// Vertex{p: [3]float32{0, h, -w}, t: [2]float32{0.0, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{-d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{-d, h, 0}, t: [2]float32{0.5, 0.0}, n: NORMALS[UP_FACE], c: c, e: e},
			// Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c, e: e},
		}

		utexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		utexture.Unbind(gl.TEXTURE_2D)

	}

	// Down Face
	if dtexture != nil {
		c := COLOUR_WHITE
		if selectedFace == DOWN_FACE {
			c = COLOUR_HIGH
		}

		e := LightLevel(pos, NORMALS[DOWN_FACE])
		v := []Vertex{
			{p: [3]float32{-d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[DOWN_FACE], c: c, e: e},
			{p: [3]float32{d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[DOWN_FACE], c: c, e: e},
			{p: [3]float32{d, -h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[DOWN_FACE], c: c, e: e},
			{p: [3]float32{-d, -h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[DOWN_FACE], c: c, e: e},
		}

		dtexture.Bind(gl.TEXTURE_2D)
		RenderQuads(v)
		dtexture.Unbind(gl.TEXTURE_2D)

	}

}

func CheckGLError() {
	glerr := gl.GetError()
	if glerr&gl.INVALID_ENUM == gl.INVALID_ENUM {
		println("gl error: INVALID_ENUM")
	}
	if glerr&gl.INVALID_VALUE != 0 {
		println("gl error: INVALID_VALUE")
	}
	if glerr&gl.INVALID_OPERATION != 0 {
		println("gl error: INVALID_OPERATION")
	}
	if glerr&gl.STACK_OVERFLOW != 0 {
		println("gl error: STACK_OVERFLOW")
	}
	if glerr&gl.STACK_UNDERFLOW != 0 {
		println("gl error: STACK_UNDERFLOW")
	}
	if glerr&gl.OUT_OF_MEMORY != 0 {
		println("gl error: OUT_OF_MEMORY")
	}
	if glerr&gl.TABLE_TOO_LARGE != 0 {
		println("gl error: TABLE_TOO_LARGE ")
	}

	if glerr != gl.NO_ERROR {
		panic("Got an OpenGL Error")
	}

}

func Line(v Vectorf) {

	gl.Begin(gl.LINE)
	gl.Vertex3f(0, 0, 0)
	gl.Vertex3f(float32(v[XAXIS]), float32(v[YAXIS]), float32(v[ZAXIS]))
	gl.End()
}

func WallSingle(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	if visible[EAST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p4, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 1.0, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)
	}

	if visible[WEST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p1, y+p4, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 1.0, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)
	}

	if visible[NORTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 1.0, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p1, 1.0/3, 0.0)
	}

	if visible[SOUTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 1.0, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p2, y+p1, z+p4, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)
	}

	if visible[UP_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p1, y+p4, z+p2, 1.0, 2.0/3,
			x+p4, y+p4, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p4, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)
	}

	if visible[DOWN_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p1, z+p4, 1.0/3, 0.0)
	}

}

func WallTee(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	WallSingle(vertexBuffer, x, y, z, orient, block, visible, shadeLevels, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//  XXX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p2, y+p4, z+p2, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p2, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 2.0/3, 1.0,
			x+p2, y+p4, z+p3, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p2, y+p4, z+p3, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p2, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_WEST {
		//
		//  XXX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p3, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p3, 2.0/3, 1.0,
			x+p2, y+p4, z+p4, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   X
		//   XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p3, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p3, y+p1, z+p3, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p3, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	}

}

func WallCorner(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	if orient == ORIENT_EAST {
		//   X
		//   XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 1,
			x+p3, y+p4, z+p2, 1.0/3, 0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 1,
			x+p2, y+p4, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p3, y+p1, z+p2, 1.0/3, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p2, y+p1, z+p3, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 1.0/3, 1.0,
			x+p2, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 1.0/3, 1.0,
			x+p2, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 2.0/3, 1.0,
			x+p3, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}

	} else if orient == ORIENT_WEST {
		//   
		//  XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p2, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p3, 1.0/3, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p3, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 1.0/3, 1.0,
			x+p2, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   
		//   XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p3, 1.0/3, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p2, 2.0/3, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p2, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p3, y+p1, z+p3, 1.0/3, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	}

}

func WallCross(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p1, z+p4, 2.0/3, 1.0,
		x+p3, y+p4, z+p3, 1.0/3, 0.0)
	vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p1, z+p1, 2.0/3, 1.0,
		x+p3, y+p4, z+p2, 1.0/3, 0.0)

	if visible[EAST_FACE] {
		// Can never actually be visible
		// vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE], 
		// 	x+p4, y+p1, z+p2, 2.0/3, 1.0,
		// 	x+p4, y+p4, z+p3, 1.0/3, 0.0)
	}

	vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p1, z+p4, 2.0/3, 1.0,
		x+p2, y+p4, z+p3, 1.0/3, 0.0)
	vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p1, z+p1, 2.0/3, 1.0,
		x+p2, y+p4, z+p2, 1.0/3, 0.0)

	if visible[WEST_FACE] {
		// Can never actually be visible
		// vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE], 
		// 	x+p1, y+p1, z+p2, 2.0/3, 1.0,
		// 	x+p1, y+p4, z+p3, 1.0/3, 0.0)

	}

	vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p1, y+p1, z+p2, 2.0/3, 1.0,
		x+p2, y+p4, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p3, y+p1, z+p2, 2.0/3, 1.0,
		x+p4, y+p4, z+p2, 1.0/3, 0.0)

	if visible[NORTH_FACE] {
		// Can never actually be visible

		// vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE], 
		// 	x+p2, y+p1, z+p1, 2.0/3, 1.0,
		// 	x+p3, y+p4, z+p1, 1.0/3, 0.0)
	}

	vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p1, y+p1, z+p3, 2.0/3, 1.0,
		x+p2, y+p4, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p3, y+p1, z+p3, 2.0/3, 1.0,
		x+p4, y+p4, z+p3, 1.0/3, 0.0)

	if visible[SOUTH_FACE] {
		// Can never actually be visible

		// vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE], 
		// 	x+p2, y+p1, z+p4, 2.0/3, 1.0,
		// 	x+p3, y+p4, z+p4, 1.0/3, 0.0)

	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p1, y+p4, z+p2, 1.0, 2.0/3,
			x+p4, y+p4, z+p3, 0.0, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p4, z+p2, 1.0/3, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p4, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p4, z+p3, 1.0/3, 1.0/3)
	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p1, z+p2, 1.0/3, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p1, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p1, z+p3, 1.0/3, 1.0/3)
	}

}

func SlabSingle(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	_ = p4
	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)

	vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p2, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	if visible[DOWN_FACE] {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 1.0/3, 1.0/3)
	}

}

func SlabLine(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)

	if visible[EAST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p4, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 1.0, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)
	}

	if visible[WEST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p1, y+p3, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 1.0, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)
	}

	if visible[NORTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p1, 1.0/3, 0.0)
	}

	if visible[SOUTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 1.0, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p2, y+p1, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 1.0/3, 0.0)
	}

	if visible[DOWN_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p1, z+p4, 1.0/3, 0.0)
	}

}

func SlabCross(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2
	_ = p2

	if visible[EAST_FACE] {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p4, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p3, z+p4, 0.0, 2.0/3)
	}

	if visible[WEST_FACE] {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p1, y+p3, z+p4, 0.0, 2.0/3)

	}

	if visible[NORTH_FACE] {
		vertexBuffer.AddFace(NORTH_FACE, block.texture2, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p3, z+p1, 0.0, 2.0/3)
	}

	if visible[SOUTH_FACE] {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture2, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p4, z+p4, 1.0, 1.0,
			x+p4, y+p3, z+p4, 0.0, 2.0/3)
	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.texture1, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p4, z+p4, 0.0, 0.0)
	}

	// underside usually visible
	vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
		x+p1, y+p3, z+p1, 1.0, 1.0,
		x+p4, y+p3, z+p4, 0.0, 0.0)

}

func SlabCorner(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, shadeLevels, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//   XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 1.0/3, 0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p3, y+p1, z+p2, 1.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 0.0, 2.0/3)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}

	} else if orient == ORIENT_WEST {
		//   
		//  XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   
		//   XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p3, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	}
}

func SlabTee(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabLine(vertexBuffer, x, y, z, orient, block, visible, shadeLevels, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//  XXX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p2, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_WEST {
		//
		//  XXX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   X
		//   XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
				x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	}

}

func Pile(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block Item, visible [6]bool, shadeLevels [6]int, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	if visible[EAST_FACE] {
		vertexBuffer.AddFace(EAST_FACE, block.texture1, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p4, y+p2, z+p1, 1.0, 1.0,
			x+p4, y+p1, z+p4, 0.0, 2.0/3)

		vertexBuffer.AddFace(EAST_FACE, block.texture2, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
			x+p4, y+p3, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p2, z+p3, 1.0/3, 0.0)

	}

	if visible[WEST_FACE] {
		vertexBuffer.AddFace(WEST_FACE, block.texture1, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p1, y+p2, z+p1, 1.0, 1.0,
			x+p1, y+p1, z+p4, 0.0, 2.0/3)

		vertexBuffer.AddFace(WEST_FACE, block.texture2, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
			x+p1, y+p3, z+p2, 2.0/3, 2.0/3,
			x+p1, y+p2, z+p3, 1.0/3, 0.0)

	}

	if visible[NORTH_FACE] {
		vertexBuffer.AddFace(NORTH_FACE, block.texture1, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p1, y+p2, z+p1, 1.0, 1.0,
			x+p4, y+p1, z+p1, 0.0, 2.0/3)

		vertexBuffer.AddFace(NORTH_FACE, block.texture2, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
			x+p2, y+p3, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p2, z+p1, 1.0/3, 0.0)

	}

	if visible[SOUTH_FACE] {
		vertexBuffer.AddFace(SOUTH_FACE, block.texture1, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p1, y+p2, z+p4, 1.0, 1.0,
			x+p4, y+p1, z+p4, 0.0, 2.0/3)

		vertexBuffer.AddFace(SOUTH_FACE, block.texture2, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
			x+p2, y+p3, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p2, z+p4, 1.0/3, 0.0)
	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.texture1, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p1, y+p2, z+p1, 1.0, 1.0,
			x+p4, y+p2, z+p4, 0.0, 0.0)
	}

	if visible[DOWN_FACE] {
		vertexBuffer.AddFace(DOWN_FACE, block.texture1, selectedFace == DOWN_FACE, shadeLevels[DOWN_FACE],
			x+p1, y+p1, z+p1, 1.0, 1.0,
			x+p4, y+p1, z+p4, 0.0, 0.0)
	}

	vertexBuffer.AddFace(EAST_FACE, block.texture2, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p4, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(EAST_FACE, block.texture2, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p3, z+p1, 2.0/3, 2.0/3,
		x+p3, y+p2, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(EAST_FACE, block.texture2, selectedFace == EAST_FACE, shadeLevels[EAST_FACE],
		x+p3, y+p3, z+p3, 2.0/3, 2.0/3,
		x+p3, y+p2, z+p4, 1.0/3, 0.0)

	vertexBuffer.AddFace(WEST_FACE, block.texture2, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
		x+p2, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(WEST_FACE, block.texture2, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p3, z+p1, 2.0/3, 2.0/3,
		x+p2, y+p2, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(WEST_FACE, block.texture2, selectedFace == WEST_FACE, shadeLevels[WEST_FACE],
		x+p2, y+p3, z+p3, 2.0/3, 2.0/3,
		x+p2, y+p2, z+p4, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.texture2, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.texture2, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p1, y+p3, z+p2, 2.0/3, 2.0/3,
		x+p2, y+p2, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.texture2, selectedFace == NORTH_FACE, shadeLevels[NORTH_FACE],
		x+p3, y+p3, z+p2, 2.0/3, 2.0/3,
		x+p4, y+p2, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.texture2, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.texture2, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p1, y+p3, z+p3, 2.0/3, 2.0/3,
		x+p2, y+p2, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.texture2, selectedFace == SOUTH_FACE, shadeLevels[SOUTH_FACE],
		x+p3, y+p3, z+p3, 2.0/3, 2.0/3,
		x+p4, y+p2, z+p3, 1.0/3, 0.0)

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
			x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p4, z+p3, 1.0/3, 1.0/3)
	}

	vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		x+p2, y+p3, z+p1, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p2, 1.0/3, 1.0/3)

	vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		x+p2, y+p3, z+p3, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p4, 1.0/3, 1.0/3)

	vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		x+p1, y+p3, z+p2, 2.0/3, 2.0/3,
		x+p2, y+p3, z+p3, 1.0/3, 1.0/3)

	vertexBuffer.AddFace(UP_FACE, block.texture2, selectedFace == UP_FACE, shadeLevels[UP_FACE],
		x+p3, y+p3, z+p2, 2.0/3, 2.0/3,
		x+p4, y+p3, z+p3, 1.0/3, 1.0/3)
}

func LightLevel(pos Vectorf, normal [3]float32) [4]float32 {
	n64 := Vectorf{float64(normal[0]), float64(normal[1]), float64(normal[2])}
	lightLevel := 0

	for _, lightSource := range lightSources {
		distance := uint16(pos.Minus(lightSource.pos).Magnitude())
		dir := lightSource.pos.Minus(pos)
		if distance < 2 || dir.Dot(n64) > 0 {
			if distance <= lightSource.intensity {
				lightLevel += int(lightSource.intensity - distance)
			}
		}
	}

	if lightLevel > MAX_LIGHT_LEVEL {
		lightLevel = MAX_LIGHT_LEVEL
	} else if lightLevel < 0 {
		lightLevel = 0
	}
	return LIGHT_LEVELS[lightLevel]
}
