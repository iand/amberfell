/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
	"image"
	_ "image/png"
	"os"
)

type TexturePos struct {
	x float32
	y float32
}

type TerrainBlock struct {
	id          byte
	name        string
	utexture    *gl.Texture
	dtexture    *gl.Texture
	ntexture    *gl.Texture
	stexture    *gl.Texture
	etexture    *gl.Texture
	wtexture    *gl.Texture
	textures    [6]uint16
	texpos      [6]TexturePos
	hitsNeeded  byte
	transparent bool
}

type Vertex struct {
	p       [3]float32 // Position
	t       [2]float32 // Texture coordinate
	n       [3]float32 // Normal
	c       [4]float32 // Colour
	padding [16]byte
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

func (self *VertexBuffer) AddFace(face uint8, texture uint16, selected bool, x1, y1, z1, tx1, ty1, x2, y2, z2, tx2, ty2 float32) {
	if self.vertexCount >= VERTEX_BUFFER_CAPACITY+4 {
		// TODO: log a warning about overflowing buffer
		return
	}

	c := COLOUR_WHITE
	if selected {
		c = COLOUR_HIGH
	}
	vc := self.vertexCount

	if x1 == x2 {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+2] = Vertex{p: [3]float32{x1, y2, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+3] = Vertex{p: [3]float32{x1, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
	} else if y1 == y2 {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+2] = Vertex{p: [3]float32{x2, y1, z2}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+3] = Vertex{p: [3]float32{x2, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
	} else {
		self.vertices[vc] = Vertex{p: [3]float32{x1, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+1] = Vertex{p: [3]float32{x1, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx1) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+2] = Vertex{p: [3]float32{x2, y2, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty2) / TILES_VERT}, n: NORMALS[face], c: c}
		self.vertices[vc+3] = Vertex{p: [3]float32{x2, y1, z1}, t: [2]float32{(float32((texture % TILES_HORZ)) + tx2) / TILES_HORZ, (float32((texture / TILES_HORZ)) + ty1) / TILES_VERT}, n: NORMALS[face], c: c}
	}

	self.vertexCount += 4
	ic := self.indexCount
	self.indices[ic] = TriangleIndex{uint32(vc), uint32(vc) + 1, uint32(vc) + 2}
	self.indices[ic+1] = TriangleIndex{uint32(vc) + 2, uint32(vc) + 3, uint32(vc)}
	self.indexCount += 2
}

func (self *VertexBuffer) RenderDirect() {
	self.texture.Bind(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	for i := 0; i < self.vertexCount; i++ {
		gl.Normal3f(self.vertices[i].n[0], self.vertices[i].n[1], self.vertices[i].n[2])
		gl.TexCoord2f(self.vertices[i].t[0], self.vertices[i].t[1])
		gl.Color4f(self.vertices[i].c[0], self.vertices[i].c[1], self.vertices[i].c[2], self.vertices[i].c[3])
		gl.Vertex3f(self.vertices[i].p[0], self.vertices[i].p[1], self.vertices[i].p[2])
	}
	gl.End()
	self.texture.Unbind(gl.TEXTURE_2D)
}

func LoadMapTextures() {

	var file, err = os.Open("tiles.png")
	var img image.Image
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if img, _, err = image.Decode(file); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		for j := 0; j < 16; j++ {
			textureIndex := uint16(i*16 + j)
			textures[textureIndex] = imageSectionToTexture(img, image.Rect(TILE_WIDTH*j, TILE_WIDTH*i, TILE_WIDTH*j+TILE_WIDTH, TILE_WIDTH*i+TILE_WIDTH))
		}
	}
}

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

func NewTerrainBlock(id byte, name string, u uint16, d uint16, n uint16, s uint16, e uint16, w uint16, hitsNeeded byte, transparent bool) TerrainBlock {
	return TerrainBlock{id, name,
		textures[u],
		textures[d],
		textures[n],
		textures[s],
		textures[e],
		textures[w],
		[6]uint16{e, w, n, s, u, d},
		[6]TexturePos{TexturePos{x: float32((e % TILES_HORZ)) / TILES_HORZ, y: float32((e / TILES_HORZ)) / TILES_VERT},
			TexturePos{x: float32((w % TILES_HORZ)) / TILES_HORZ, y: float32((w / TILES_HORZ)) / TILES_VERT},
			TexturePos{x: float32((n % TILES_HORZ)) / TILES_HORZ, y: float32((n / TILES_HORZ)) / TILES_VERT},
			TexturePos{x: float32((s % TILES_HORZ)) / TILES_HORZ, y: float32((s / TILES_HORZ)) / TILES_VERT},
			TexturePos{x: float32((u % TILES_HORZ)) / TILES_HORZ, y: float32((u / TILES_HORZ)) / TILES_VERT},
			TexturePos{x: float32((d % TILES_HORZ)) / TILES_HORZ, y: float32((d / TILES_HORZ)) / TILES_VERT},
		},
		hitsNeeded,
		transparent,
	}
}

func InitTerrainBlocks() {
	TerrainBlocks = make(map[uint16]TerrainBlock)
	TerrainBlocks[BLOCK_AIR] = NewTerrainBlock(BLOCK_AIR, "Air", TEXTURE_NONE, TEXTURE_NONE, TEXTURE_NONE, TEXTURE_NONE, TEXTURE_NONE, TEXTURE_NONE, STRENGTH_STONE, true)
	TerrainBlocks[BLOCK_STONE] = NewTerrainBlock(BLOCK_STONE, "Stone", TEXTURE_STONE, TEXTURE_STONE, TEXTURE_STONE, TEXTURE_STONE, TEXTURE_STONE, TEXTURE_STONE, STRENGTH_STONE, false)
	TerrainBlocks[BLOCK_DIRT] = NewTerrainBlock(BLOCK_DIRT, "Dirt", TEXTURE_DIRT_TOP, TEXTURE_DIRT, TEXTURE_DIRT, TEXTURE_DIRT, TEXTURE_DIRT, TEXTURE_DIRT, STRENGTH_DIRT, false)
	TerrainBlocks[BLOCK_TRUNK] = NewTerrainBlock(BLOCK_TRUNK, "trunk", TEXTURE_TRUNK, TEXTURE_TRUNK, TEXTURE_TRUNK, TEXTURE_TRUNK, TEXTURE_TRUNK, TEXTURE_TRUNK, STRENGTH_WOOD, false)
	TerrainBlocks[BLOCK_LEAVES] = NewTerrainBlock(BLOCK_TRUNK, "trunk", TEXTURE_LEAVES, TEXTURE_LEAVES, TEXTURE_LEAVES, TEXTURE_LEAVES, TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false)
	TerrainBlocks[BLOCK_LOG_WALL] = NewTerrainBlock(BLOCK_LOG_WALL, "log wall", TEXTURE_LOG_WALL_TOP, TEXTURE_LOG_WALL_TOP, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, STRENGTH_WOOD, true)
	TerrainBlocks[BLOCK_LOG_SLAB] = NewTerrainBlock(BLOCK_LOG_SLAB, "log slab", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, STRENGTH_WOOD, true)

}

func TerrainCube(vertexBuffer *VertexBuffer, x float32, y float32, z float32, neighbours [6]uint16, blockid byte, selectedFace uint8) {

	block := TerrainBlocks[uint16(blockid)]
	var visible [6]bool

	for i := 0; i < 6; i++ {
		if TerrainBlocks[neighbours[i]].transparent {
			visible[i] = true
		}
	}

	switch blockid {
	case BLOCK_LOG_SLAB:
		if neighbours[NORTH_FACE] != BLOCK_AIR {
			if neighbours[EAST_FACE] != BLOCK_AIR {
				if neighbours[SOUTH_FACE] != BLOCK_AIR {
					if neighbours[WEST_FACE] != BLOCK_AIR {
						// Blocks to all four sides
						SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
					} else {
						// Blocks to north, east, south
						SlabTee(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, selectedFace)
					}
				} else if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, east, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
				} else {
					// Blocks to north, east
					SlabCorner(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
				}

			} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, south, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
				} else {
					// Blocks to north, south
					SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to the north and west
				SlabCorner(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)

			} else {
				// Just a block to the north
				SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
			}

		} else if neighbours[EAST_FACE] != BLOCK_AIR {
			if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to east, south, west
					SlabTee(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, selectedFace)
				} else {
					// Blocks to east, south
					SlabCorner(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to east, west
				SlabLine(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
			} else {
				// Just a block to the east
				SlabLine(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
			}
		} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
			if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to south, west
				SlabCorner(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, selectedFace)
			} else {
				// Just a block to the south
				SlabLine(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
			}
		} else if neighbours[WEST_FACE] != BLOCK_AIR {
			// Just a block to the west
			SlabLine(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, selectedFace)
		} else {
			// Lone block
			SlabSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
		}

	case BLOCK_LOG_WALL:
		if neighbours[NORTH_FACE] != BLOCK_AIR {
			if neighbours[EAST_FACE] != BLOCK_AIR {
				if neighbours[SOUTH_FACE] != BLOCK_AIR {
					if neighbours[WEST_FACE] != BLOCK_AIR {
						// Blocks to all four sides
						WallCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
					} else {
						// Blocks to north, east, south
						WallTee(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, selectedFace)
					}
				} else if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, east, west
					WallTee(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
				} else {
					// Blocks to north, east
					WallCorner(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
				}

			} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to north, south, west
					WallTee(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
				} else {
					// Blocks to north, south
					WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to the north and west
				WallCorner(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)

			} else {
				// Just a block to the north
				WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
			}

		} else if neighbours[EAST_FACE] != BLOCK_AIR {
			if neighbours[SOUTH_FACE] != BLOCK_AIR {
				if neighbours[WEST_FACE] != BLOCK_AIR {
					// Blocks to east, south, west
					WallTee(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, selectedFace)
				} else {
					// Blocks to east, south
					WallCorner(vertexBuffer, x, y, z, ORIENT_SOUTH, block, visible, selectedFace)
				}
			} else if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to east, west
				WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
			} else {
				// Just a block to the east
				WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
			}
		} else if neighbours[SOUTH_FACE] != BLOCK_AIR {
			if neighbours[WEST_FACE] != BLOCK_AIR {
				// Blocks to south, west
				WallCorner(vertexBuffer, x, y, z, ORIENT_WEST, block, visible, selectedFace)
			} else {
				// Just a block to the south
				WallSingle(vertexBuffer, x, y, z, ORIENT_NORTH, block, visible, selectedFace)
			}
		} else if neighbours[WEST_FACE] != BLOCK_AIR {
			// Just a block to the west
			WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
		} else {
			// Lone block
			WallSingle(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)
		}

	default:
		Cuboid2(vertexBuffer, x, y, z, 1, 1, 1, block, visible, selectedFace)

	}

}

func RenderQuads(v []Vertex) {
	gl.Begin(gl.QUADS)
	for i := 0; i < len(v); i++ {
		gl.Normal3f(v[i].n[0], v[i].n[1], v[i].n[2])
		gl.TexCoord2f(v[i].t[0], v[i].t[1])
		gl.Color4f(v[i].c[0], v[i].c[1], v[i].c[2], v[i].c[3])
		gl.Vertex3f(v[i].p[0], v[i].p[1], v[i].p[2])
	}
	gl.End()
}

func Cuboid2(vertexBuffer *VertexBuffer, x float32, y float32, z float32, bw float64, bh float64, bd float64, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2

	// East face
	if visible[EAST_FACE] {

		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+d, y+h, z-w, 1.0, 1.0,
			x+d, y-h, z+w, 0.0, 0.0)

	}

	// West Face
	if visible[WEST_FACE] {

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x-d, y+h, z-w, 1.0, 1.0,
			x-d, y-h, z+w, 0.0, 0.0)
	}

	// North Face
	if visible[NORTH_FACE] {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+d, y+h, z-w, 1.0, 1.0,
			x-d, y-h, z-w, 0.0, 0.0)
	}

	// South Face
	if visible[SOUTH_FACE] {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
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

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+d, y+h, z-w, 1.0, 1.0,
			x+0, y+h, z+0, 0.5, 0.5)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+d, y+h, z+0, 1.0, 0.5,
			x+0, y+h, z+w, 0.5, 0.0)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+0, y+h, z+0, 0.5, 0.5,
			x-d, y+h, z+w, 0.0, 0.0)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+0, y+h, z-w, 0.5, 1.0,
			x-d, y+h, z+0, 0.0, 0.5)
	}

	// Down Face
	if visible[DOWN_FACE] {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+d, y-h, z-w, 1.0, 1.0,
			x-d, y-h, z+w, 0.0, 0.0)
	}

}

func Cuboid(bw float64, bh float64, bd float64, etexture *gl.Texture, wtexture *gl.Texture, ntexture *gl.Texture, stexture *gl.Texture, utexture *gl.Texture, dtexture *gl.Texture, selectedFace uint8) {

	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2

	// East face
	if etexture != nil {

		c := COLOUR_WHITE
		if selectedFace == EAST_FACE {
			c = COLOUR_HIGH
		}

		v := []Vertex{
			Vertex{p: [3]float32{d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[EAST_FACE], c: c},
			Vertex{p: [3]float32{d, h, -w}, t: [2]float32{1.0, 0.0}, n: NORMALS[EAST_FACE], c: c},
			Vertex{p: [3]float32{d, h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[EAST_FACE], c: c},
			Vertex{p: [3]float32{d, -h, w}, t: [2]float32{0.0, 1.0}, n: NORMALS[EAST_FACE], c: c},
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

		v := []Vertex{
			Vertex{p: [3]float32{-d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[WEST_FACE], c: c},
			Vertex{p: [3]float32{-d, -h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[WEST_FACE], c: c},
			Vertex{p: [3]float32{-d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[WEST_FACE], c: c},
			Vertex{p: [3]float32{-d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[WEST_FACE], c: c},
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

		v := []Vertex{
			Vertex{p: [3]float32{-d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[NORTH_FACE], c: c},
			Vertex{p: [3]float32{-d, h, -w}, t: [2]float32{1.0, 0.0}, n: NORMALS[NORTH_FACE], c: c},
			Vertex{p: [3]float32{d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[NORTH_FACE], c: c},
			Vertex{p: [3]float32{d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[NORTH_FACE], c: c},
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

		v := []Vertex{
			Vertex{p: [3]float32{-d, -h, w}, t: [2]float32{0.0, 1.0}, n: NORMALS[SOUTH_FACE], c: c},
			Vertex{p: [3]float32{d, -h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[SOUTH_FACE], c: c},
			Vertex{p: [3]float32{d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[SOUTH_FACE], c: c},
			Vertex{p: [3]float32{-d, h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[SOUTH_FACE], c: c},
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
			Vertex{p: [3]float32{d, h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, -w}, t: [2]float32{0.0, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{d, h, 0}, t: [2]float32{0.5, 1.0}, n: NORMALS[UP_FACE], c: c},

			Vertex{p: [3]float32{d, h, 0}, t: [2]float32{0.5, 1.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, w}, t: [2]float32{1.0, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{d, h, w}, t: [2]float32{1.0, 1.0}, n: NORMALS[UP_FACE], c: c},

			Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{-d, h, 0}, t: [2]float32{0.5, 0.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{-d, h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, w}, t: [2]float32{1.0, 0.5}, n: NORMALS[UP_FACE], c: c},

			Vertex{p: [3]float32{0, h, -w}, t: [2]float32{0.0, 0.5}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{-d, h, -w}, t: [2]float32{0.0, 0.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{-d, h, 0}, t: [2]float32{0.5, 0.0}, n: NORMALS[UP_FACE], c: c},
			Vertex{p: [3]float32{0, h, 0}, t: [2]float32{0.5, 0.5}, n: NORMALS[UP_FACE], c: c},
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

		v := []Vertex{
			Vertex{p: [3]float32{-d, -h, -w}, t: [2]float32{1.0, 1.0}, n: NORMALS[DOWN_FACE], c: c},
			Vertex{p: [3]float32{d, -h, -w}, t: [2]float32{0.0, 1.0}, n: NORMALS[DOWN_FACE], c: c},
			Vertex{p: [3]float32{d, -h, w}, t: [2]float32{0.0, 0.0}, n: NORMALS[DOWN_FACE], c: c},
			Vertex{p: [3]float32{-d, -h, w}, t: [2]float32{1.0, 0.0}, n: NORMALS[DOWN_FACE], c: c},
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

func WallSingle(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	if visible[EAST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p4, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 1.0, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)
	}

	if visible[WEST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p1, y+p4, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 1.0, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)
	}

	if visible[NORTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 1.0, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p1, 1.0/3, 0.0)
	}

	if visible[SOUTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 1.0, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p2, y+p1, z+p4, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)
	}

	if visible[UP_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p1, y+p4, z+p2, 1.0, 2.0/3,
			x+p4, y+p4, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p2, y+p4, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)
	}

	if visible[DOWN_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p1, z+p4, 1.0/3, 0.0)
	}

}

func WallTee(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	WallSingle(vertexBuffer, x, y, z, orient, block, visible, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//  XXX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p2, y+p4, z+p2, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p2, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 2.0/3, 1.0,
			x+p2, y+p4, z+p3, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p2, y+p4, z+p3, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p2, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_WEST {
		//
		//  XXX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p3, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p3, 2.0/3, 1.0,
			x+p2, y+p4, z+p4, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   X
		//   XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p3, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p3, y+p1, z+p3, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 1.0/3, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p3, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	}

}

func WallCorner(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	if orient == ORIENT_EAST {
		//   X
		//   XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 1,
			x+p3, y+p4, z+p2, 1.0/3, 0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 1,
			x+p2, y+p4, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p3, y+p1, z+p2, 1.0/3, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p2, y+p1, z+p3, 2.0/3, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p4, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 1.0/3, 1.0,
			x+p2, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 1.0/3, 1.0,
			x+p2, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 2.0/3, 1.0,
			x+p3, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p2, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}

	} else if orient == ORIENT_WEST {
		//   
		//  XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p2, 2.0/3, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p3, 1.0/3, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 1.0,
			x+p3, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 1.0/3, 1.0,
			x+p2, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p1, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   
		//   XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p3, 1.0/3, 1.0,
			x+p3, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p2, 2.0/3, 1.0,
			x+p2, y+p4, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p2, y+p1, z+p2, 2.0/3, 1.0,
			x+p4, y+p4, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p3, y+p1, z+p3, 1.0/3, 1.0,
			x+p4, y+p4, z+p3, 0.0, 0.0)

		if visible[UP_FACE] {
			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p4, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
				x+p2, y+p4, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p4, z+p4, 1.0/3, 1.0/3)
		}

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	}

}

func WallCross(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
		x+p3, y+p1, z+p4, 2.0/3, 1.0,
		x+p3, y+p4, z+p3, 1.0/3, 0.0)
	vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
		x+p3, y+p1, z+p1, 2.0/3, 1.0,
		x+p3, y+p4, z+p2, 1.0/3, 0.0)

	if visible[EAST_FACE] {
		// Can never actually be visible
		// vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
		// 	x+p4, y+p1, z+p2, 2.0/3, 1.0,
		// 	x+p4, y+p4, z+p3, 1.0/3, 0.0)
	}

	vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
		x+p2, y+p1, z+p4, 2.0/3, 1.0,
		x+p2, y+p4, z+p3, 1.0/3, 0.0)
	vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
		x+p2, y+p1, z+p1, 2.0/3, 1.0,
		x+p2, y+p4, z+p2, 1.0/3, 0.0)

	if visible[WEST_FACE] {
		// Can never actually be visible
		// vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
		// 	x+p1, y+p1, z+p2, 2.0/3, 1.0,
		// 	x+p1, y+p4, z+p3, 1.0/3, 0.0)

	}

	vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
		x+p1, y+p1, z+p2, 2.0/3, 1.0,
		x+p2, y+p4, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
		x+p3, y+p1, z+p2, 2.0/3, 1.0,
		x+p4, y+p4, z+p2, 1.0/3, 0.0)

	if visible[NORTH_FACE] {
		// Can never actually be visible

		// vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
		// 	x+p2, y+p1, z+p1, 2.0/3, 1.0,
		// 	x+p3, y+p4, z+p1, 1.0/3, 0.0)
	}

	vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
		x+p1, y+p1, z+p3, 2.0/3, 1.0,
		x+p2, y+p4, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
		x+p3, y+p1, z+p3, 2.0/3, 1.0,
		x+p4, y+p4, z+p3, 1.0/3, 0.0)

	if visible[SOUTH_FACE] {
		// Can never actually be visible

		// vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
		// 	x+p2, y+p1, z+p4, 2.0/3, 1.0,
		// 	x+p3, y+p4, z+p4, 1.0/3, 0.0)

	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p1, y+p4, z+p2, 1.0, 2.0/3,
			x+p4, y+p4, z+p3, 0.0, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p2, y+p4, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p4, z+p2, 1.0/3, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p2, y+p4, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p4, z+p3, 1.0/3, 1.0/3)
	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p1, z+p2, 1.0/3, 1.0/3)

		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p2, y+p1, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p1, z+p3, 1.0/3, 1.0/3)
	}

}

func SlabSingle(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	_ = p4
	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)

	vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
		x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
		x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p2, y+p3, z+p3, 1.0/3, 0.0)

	vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
		x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p2, 1.0/3, 0.0)

	vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
		x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
		x+p3, y+p3, z+p3, 1.0/3, 0.0)

	if visible[DOWN_FACE] {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 1.0/3, 1.0/3)
	}

}

func SlabLine(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)

	if visible[EAST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p4, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 1.0, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)
	}

	if visible[WEST_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p1, y+p3, z+p3, 1.0/3, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 1.0, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)
	}

	if visible[NORTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p1, 1.0/3, 0.0)
	}

	if visible[SOUTH_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 1.0, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p2, y+p1, z+p4, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 1.0/3, 0.0)
	}

	if visible[DOWN_FACE] && (orient == ORIENT_EAST || orient == ORIENT_WEST) {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+p1, y+p1, z+p2, 1.0, 2.0/3,
			x+p4, y+p1, z+p3, 0.0, 1.0/3)
	} else if orient == ORIENT_NORTH || orient == ORIENT_SOUTH {
		vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 1.0,
			x+p3, y+p1, z+p4, 1.0/3, 0.0)
	}

}

func SlabCross(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2
	_ = p2

	if visible[EAST_FACE] {
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p4, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p3, z+p4, 0.0, 2.0/3)
	}

	if visible[WEST_FACE] {
		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p1, y+p3, z+p4, 0.0, 2.0/3)

	}

	if visible[NORTH_FACE] {
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p3, z+p1, 0.0, 2.0/3)
	}

	if visible[SOUTH_FACE] {
		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p4, z+p4, 1.0, 1.0,
			x+p4, y+p3, z+p4, 0.0, 2.0/3)
	}

	if visible[UP_FACE] {
		vertexBuffer.AddFace(UP_FACE, block.textures[UP_FACE], selectedFace == UP_FACE,
			x+p1, y+p4, z+p1, 1.0, 1.0,
			x+p4, y+p4, z+p4, 0.0, 0.0)
	}

	// underside usually visible
	vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
		x+p1, y+p3, z+p1, 1.0, 1.0,
		x+p4, y+p3, z+p4, 0.0, 0.0)

}

func SlabCorner(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabCross(vertexBuffer, x, y, z, ORIENT_EAST, block, visible, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//   XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 1.0/3, 0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p3, y+p1, z+p2, 1.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 0.0, 2.0/3)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}

	} else if orient == ORIENT_WEST {
		//   
		//  XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   
		//   XX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 0.0, 0.0)

		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 0.0, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p3, y+p1, z+p3, 1.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 0.0, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p2, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 0.0, 1.0/3)

			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	}
}

func SlabTee(vertexBuffer *VertexBuffer, x float32, y float32, z float32, orient byte, block TerrainBlock, visible [6]bool, selectedFace uint8) {
	var p1, p2, p3, p4 float32 = -1.0 / 2, -1.0 / 6, 1.0 / 6, 1.0 / 2

	SlabLine(vertexBuffer, x, y, z, orient, block, visible, selectedFace)

	if orient == ORIENT_EAST {
		//   X
		//  XXX 
		//   
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p1, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p2, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_NORTH {
		//   X
		//  XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p1, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p3, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p1, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p2, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_WEST {
		//
		//  XXX 
		//   X
		vertexBuffer.AddFace(EAST_FACE, block.textures[EAST_FACE], selectedFace == EAST_FACE,
			x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p3, y+p3, z+p4, 1.0/3, 0.0)

		vertexBuffer.AddFace(WEST_FACE, block.textures[WEST_FACE], selectedFace == WEST_FACE,
			x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p2, y+p3, z+p4, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p2, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p3, y+p1, z+p4, 1.0/3, 1.0/3)
		}
	} else if orient == ORIENT_SOUTH {
		//   X
		//   XX 
		//   X  
		vertexBuffer.AddFace(NORTH_FACE, block.textures[NORTH_FACE], selectedFace == NORTH_FACE,
			x+p3, y+p1, z+p2, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p2, 1.0/3, 0.0)

		vertexBuffer.AddFace(SOUTH_FACE, block.textures[SOUTH_FACE], selectedFace == SOUTH_FACE,
			x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
			x+p4, y+p3, z+p3, 1.0/3, 0.0)

		if visible[DOWN_FACE] {
			vertexBuffer.AddFace(DOWN_FACE, block.textures[DOWN_FACE], selectedFace == DOWN_FACE,
				x+p3, y+p1, z+p3, 2.0/3, 2.0/3,
				x+p4, y+p1, z+p3, 1.0/3, 1.0/3)
		}
	}

}
