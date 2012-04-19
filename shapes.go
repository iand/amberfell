/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/gl"
	"image"
	_ "image/png"
	"os"
)

var (
	MapTextures   [16 * 10]gl.Texture
	TerrainCubes  map[uint16]uint
	TerrainBlocks map[byte]TerrainBlock
)

type TerrainBlock struct {
	id       byte
	name     string
	utexture *gl.Texture
	dtexture *gl.Texture
	ntexture *gl.Texture
	stexture *gl.Texture
	etexture *gl.Texture
	wtexture *gl.Texture
}

type FeedbackBuffer struct {
	buffer [4096]float32
	size   int
}

func (self *FeedbackBuffer) Dump() {

	for i := 0; i < self.size; i++ {
		if self.buffer[i] == gl.PASS_THROUGH_TOKEN {
			fmt.Printf("GL_PASS_THROUGH_TOKEN\n")
			i++
			fmt.Printf("  %4.2f\n", self.buffer[i])
		} else if self.buffer[i] == gl.POINT_TOKEN {
			fmt.Printf("GL_POINT_TOKEN\n")
			i++
			i += 11
		} else if self.buffer[i] == gl.LINE_TOKEN {
			fmt.Printf("GL_LINE_TOKEN\n")
			i++
			i += 11
		} else if self.buffer[i] == gl.LINE_RESET_TOKEN {
			fmt.Printf("GL_LINE_RESET_TOKEN\n")
			i++
			i += 11
		} else if self.buffer[i] == gl.POLYGON_TOKEN {
			fmt.Printf("GL_POLYGON_TOKEN\n")
			i++
			fmt.Printf("  x: %4.2f y: %4.2f z: %4.2f\n", self.buffer[i+1], self.buffer[i+2], self.buffer[i+3])
			i += 11

		} else {
			fmt.Printf("  %4.2f\n", self.buffer[i])
		}
	}
}

func LoadMapTextures() {
	const pixels = 48

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
			rgba := image.NewRGBA(image.Rect(0, 0, pixels, pixels))
			for x := 0; x < pixels; x++ {
				for y := 0; y < pixels; y++ {
					rgba.Set(x, y, img.At(pixels*j+x, pixels*i+y))
				}
			}

			textureIndex := i*16 + j
			MapTextures[textureIndex] = gl.GenTexture()
			MapTextures[textureIndex].Bind(gl.TEXTURE_2D)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
			// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, pixels, pixels, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
			MapTextures[textureIndex].Unbind(gl.TEXTURE_2D)

		}
	}
}

func InitTerrainBlocks() {
	TerrainBlocks = make(map[byte]TerrainBlock)
	TerrainBlocks[BLOCK_STONE] = TerrainBlock{BLOCK_STONE, "Stone", &MapTextures[1], &MapTextures[17], &MapTextures[17], &MapTextures[17], &MapTextures[17], &MapTextures[17]}
	TerrainBlocks[BLOCK_DIRT] = TerrainBlock{BLOCK_DIRT, "Dirt", &MapTextures[2], &MapTextures[18], &MapTextures[18], &MapTextures[18], &MapTextures[18], &MapTextures[18]}
}

func LoadTerrainCubes() {
	TerrainCubes = make(map[uint16]uint)
	var faces byte
	for blockid, block := range TerrainBlocks {
		for faces = 1; faces < 64; faces++ {
			listid := gl.GenLists(1)
			if listid == 0 {
				panic("GenLists return 0")
			}
			var etexture, wtexture, ntexture, stexture, utexture, dtexture *gl.Texture
			if faces&32 == 32 {
				etexture = block.etexture
			}
			if faces&16 == 16 {
				wtexture = block.wtexture
			}
			if faces&8 == 8 {
				ntexture = block.ntexture
			}
			if faces&4 == 4 {
				stexture = block.stexture
			}
			if faces&2 == 2 {
				utexture = block.utexture
			}
			if faces&1 == 1 {
				dtexture = block.dtexture
			}

			gl.NewList(listid, gl.COMPILE)
			Cuboid(1, 1, 1, etexture, wtexture, ntexture, stexture, utexture, dtexture, 0, false)
			gl.EndList()
			// CheckGLError()

			var index uint16 = uint16(blockid)<<8 + uint16(faces)
			TerrainCubes[index] = listid
		}

	}
}

func terrainCubeIndex(n bool, s bool, w bool, e bool, u bool, d bool, blockid byte) uint16 {
	var index uint16 = uint16(blockid) << 8
	if e {
		index += 32
	}
	if w {
		index += 16
	}
	if n {
		index += 8
	}
	if w {
		index += 4
	}
	if u {
		index += 2
	}
	if d {
		index += 1
	}

	return index
}

func TerrainCube(n bool, s bool, w bool, e bool, u bool, d bool, blockid byte, id uint16, selectMode bool) {
	var ntexture, stexture, etexture, wtexture, utexture, dtexture *gl.Texture

	block := TerrainBlocks[blockid]

	if n {
		ntexture = block.ntexture
	}
	if s {
		stexture = block.stexture
	}
	if e {
		etexture = block.etexture
	}
	if w {
		wtexture = block.wtexture
	}
	if u {
		utexture = block.utexture
	}
	// if d { dtexture = block.dtexture }

	// if selectMode {
	Cuboid(1, 1, 1, etexture, wtexture, ntexture, stexture, utexture, dtexture, id, selectMode)
	// } else {
	// gl.CallList(TerrainCubes[terrainCubeIndex(n, s, e, w, u, d, blockid)])
	// }

}

func Cuboid(bw float64, bh float64, bd float64, etexture *gl.Texture, wtexture *gl.Texture, ntexture *gl.Texture, stexture *gl.Texture, utexture *gl.Texture, dtexture *gl.Texture, id uint16, selectMode bool) {
	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2

	// Front face
	if etexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + EAST_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), EAST_FACE, 0)
		}
		etexture.Bind(gl.TEXTURE_2D)

		// gl.EnableClientState(gl.VERTEX_ARRAY)        // Enable Vertex Arrays
		// gl.EnableClientState(gl.TEXTURE_COORD_ARRAY) // Enable Texture Coord Arrays
		// gl.EnableClientState(gl.NORMAL_ARRAY)        // Enable Texture Coord Arrays
		// gl.VertexPointer(3, 0, []float32{d, -h, -w, d, h, -w, d, h, w, d, -h, w})
		// gl.TexCoordPointer(2, 0, []float32{1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0})
		// gl.NormalPointer(0, []float32{1.0, 0.0, 0.0})
		// gl.DrawArrays(gl.QUADS, 0, 4)

		gl.Begin(gl.QUADS)
		gl.Normal3f(1.0, 0.0, 0.0)
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(d, -h, -w) // Bottom Right Of The Texture and Quad
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(d, h, -w) // Top Right Of The Texture and Quad
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(d, h, w) // Top Left Of The Texture and Quad
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(d, -h, w) // Bottom Left Of The Texture and Quad
		gl.End()
		// etexture.Unbind(gl.TEXTURE_2D)
		// CheckGLError()
	}
	// Back Face
	if wtexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + WEST_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), WEST_FACE, 0)
		}
		wtexture.Bind(gl.TEXTURE_2D)
		gl.Begin(gl.QUADS)
		gl.Normal3f(-1.0, 0.0, 0.0)
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(-d, -h, -w) // Bottom Left Of The Texture and Quad
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(-d, -h, w) // Bottom Right Of The Texture and Quad
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(-d, h, w) // Top Right Of The Texture and Quad
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(-d, h, -w) // Top Left Of The Texture and Quad
		gl.End()
		// wtexture.Unbind(gl.TEXTURE_2D)

		// CheckGLError()
	}

	// Left Face
	if ntexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + NORTH_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), NORTH_FACE, 0)
		}
		ntexture.Bind(gl.TEXTURE_2D)
		gl.Begin(gl.QUADS)
		gl.Normal3f(0.0, 0.0, -1.0)
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(-d, -h, -w) // Bottom Right Of The Texture and Quad
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(-d, h, -w) // Top Right Of The Texture and Quad
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(d, h, -w) // Top Left Of The Texture and Quad
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(d, -h, -w) // Bottom Left Of The Texture and Quad
		gl.End()
		// ntexture.Unbind(gl.TEXTURE_2D)
	}

	// Right Face
	if stexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + SOUTH_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), SOUTH_FACE, 0)
		}
		stexture.Bind(gl.TEXTURE_2D)
		gl.Begin(gl.QUADS)
		gl.Normal3f(0.0, 0.0, 1.0)
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(-d, -h, w) // Bottom Left Of The Texture and Quad
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(d, -h, w) // Bottom Right Of The Texture and Quad
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(d, h, w) // Top Right Of The Texture and Quad
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(-d, h, w) // Top Left Of The Texture and Quad
		gl.End()
		// stexture.Unbind(gl.TEXTURE_2D)

		// CheckGLError()
	}

	// Top Face
	if utexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + UP_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), UP_FACE, 0)
		}
		utexture.Bind(gl.TEXTURE_2D)

		gl.Begin(gl.QUADS)
		// gl.Begin(gl.TRIANGLES)

		// gl.TexCoord2f(0.0, 1.0)
		// gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 0.0)
		// gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 0.0)
		// gl.Vertex3f( d,  h,  w)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1.0, 1.0)
		// gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

		//  -d/-w   -d/0   -d/+w
		//
		//  0/-w     0/0     0/+w
		//
		//  +d/-w   +d/0   +d/+w

		// Texture
		// 0.0/1.0    0.0/0.5   0.0/0.0
		// 0.5/1.0    0.5/0.5   0.5/0.0
		// 1.0/1.0    1.0/0.5   1.0/0.0

		// 2x2 Subsquares
		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(-d, h, -w) // Top Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 0.5)
		gl.Vertex3f(-d, h, 0) // Bottom Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.5)
		gl.Vertex3f(0, h, 0) // Bottom Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 1.0)
		gl.Vertex3f(0, h, -w) // Top Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 1.0)
		gl.Vertex3f(0, h, -w) // Top Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.5)
		gl.Vertex3f(0, h, 0) // Bottom Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(1.0, 0.5)
		gl.Vertex3f(d, h, 0) // Bottom Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(d, h, -w) // Top Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.5)
		gl.Vertex3f(0, h, 0) // Top Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.0)
		gl.Vertex3f(0, h, w) // Bottom Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(d, h, w) // Bottom Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(1.0, 0.5)
		gl.Vertex3f(d, h, 0) // Top Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 0.5)
		gl.Vertex3f(-d, h, 0) // Top Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(-d, h, w) // Bottom Left Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.0)
		gl.Vertex3f(0, h, w) // Bottom Right Of The Texture and Quad

		gl.Normal3f(0.0, 1.0, 0.0)
		gl.TexCoord2f(0.5, 0.5)
		gl.Vertex3f(0, h, 0) // Top Right Of The Texture and Quad

		// Triangles
		// gl.TexCoord2f(0.0, 1.0)
		// gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 0.0)
		// gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(0.5, 0.5)
		// gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad

		// gl.TexCoord2f(0.0, 0.0)
		// gl.Vertex3f(-d,  h,  w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 0.0)
		// gl.Vertex3f(d,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(0.5, 0.5)
		// gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad

		// gl.TexCoord2f(1.0, 0.0)
		// gl.Vertex3f(d,  h,  w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 1.0)
		// gl.Vertex3f(d,  h,  -w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(0.5, 0.5)
		// gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad

		// gl.TexCoord2f(1.0, 1.0)
		// gl.Vertex3f(d,  h,  -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 1.0)
		// gl.Vertex3f(-d,  h,  -w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(0.5, 0.5)
		// gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad

		// DRAW 3x3 Subsquares

		//  -d:-w   -d:-w/3     -d:w/3     -d:+w
		//
		//  -d/3:-w  -d/3:-w/3  -d/3:w/3   -d/3:w
		//
		//  +d/3:-w  +d/3:-w/3  +d/3:w/3   +d/3:w
		//
		//  +d:-w  +d:-w/3  +d:w/3   +d:w

		// Texture
		// 0.0:1.0    0.0:2/3   0.0:1/3   0.0:0.0
		// 1/3:1.0    1/3:2/3   1/3:1/3   1/3:0.0
		// 2/3:1.0    2/3:2/3   2/3:1/3   2/3:0.0
		// 1.0:1.0    1.0:2/3   1.0:1/3   1.0:0.0

		// FIRST
		// gl.TexCoord2f(0.0, 1.0)
		// gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 2/3)
		// gl.Vertex3f(-d,  h,  -w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 2/3)
		// gl.Vertex3f(-d/3,  h, -w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1/3, 1.0)
		// gl.Vertex3f(-d/3,  h, -w)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(0.0, 2/3)
		// gl.Vertex3f(-d,  h, -w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 1/3)
		// gl.Vertex3f(-d,  h,  w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 1/3)
		// gl.Vertex3f(-d/3,  h, w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1/3, 2/3)
		// gl.Vertex3f(-d/3,  h, -w/3)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(0.0, 1/3)
		// gl.Vertex3f(-d,  h, w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(0.0, 0.0)
		// gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 0.0)
		// gl.Vertex3f(-d/3,  h, w)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1/3, 1/3)
		// gl.Vertex3f(-d/3,  h, w/3)  // Top Right Of The Texture and Quad

		// // SECOND
		// gl.TexCoord2f(1/3, 1.0)
		// gl.Vertex3f(-d/3,  h, -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 2/3)
		// gl.Vertex3f(-d/3,  h,  -w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 2/3)
		// gl.Vertex3f( d/3,  h, -w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(2/3, 1.0)
		// gl.Vertex3f( d/3,  h, -w)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(1/3, 2/3)
		// gl.Vertex3f(-d/3,  h, -w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 1/3)
		// gl.Vertex3f(-d/3,  h,  w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 1/3)
		// gl.Vertex3f( d/3,  h, w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(2/3, 2/3)
		// gl.Vertex3f( d/3,  h, -w/3)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(1/3, 1/3)
		// gl.Vertex3f(-d/3,  h, w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(1/3, 0.0)
		// gl.Vertex3f(-d/3,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 0.0)
		// gl.Vertex3f( d/3,  h, w)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(2/3, 1/3)
		// gl.Vertex3f( d/3,  h, w/3)  // Top Right Of The Texture and Quad

		// // THIRD
		// gl.TexCoord2f(2/3, 1/3)
		// gl.Vertex3f(d/3,  h, w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 0.0)
		// gl.Vertex3f(d/3,  h,  w)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 0.0)
		// gl.Vertex3f( d,  h, w)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1.0, 1/3)
		// gl.Vertex3f( d,  h, w/3)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(2/3, 2/3)
		// gl.Vertex3f(d/3,  h, -w/3)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 1/3)
		// gl.Vertex3f(d/3,  h,  w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 1/3)
		// gl.Vertex3f( d,  h, w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1.0, 2/3)
		// gl.Vertex3f( d,  h, -w/3)  // Top Right Of The Texture and Quad

		// gl.TexCoord2f(2/3, 1.0)
		// gl.Vertex3f(d/3,  h, -w)  // Top Left Of The Texture and Quad
		// gl.TexCoord2f(2/3, 2/3)
		// gl.Vertex3f(d/3,  h,  -w/3)  // Bottom Left Of The Texture and Quad
		// gl.TexCoord2f(1.0, 2/3)
		// gl.Vertex3f( d,  h, -w/3)  // Bottom Right Of The Texture and Quad
		// gl.TexCoord2f(1.0, 1.0)
		// gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

		gl.End()
		// utexture.Unbind(gl.TEXTURE_2D)
		// CheckGLError()
	}

	// Bottom Face
	if dtexture != nil {
		if selectMode {
			gl.PassThrough(float32(id*65336) + DOWN_FACE)
			// gl.Color4ub(uint8(id&0x00FF), uint8(id&0xFF00>>8), DOWN_FACE, 0)
		}
		dtexture.Bind(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.Begin(gl.QUADS)
		gl.Normal3f(0.0, -1.0, 0.0)
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(-d, -h, -w) // Top Right Of The Texture and Quad
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(d, -h, -w) // Top Left Of The Texture and Quad
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(d, -h, w) // Bottom Left Of The Texture and Quad
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(-d, -h, w) // Bottom Right Of The Texture and Quad
		gl.End()
		// dtexture.Unbind(gl.TEXTURE_2D)
		// CheckGLError()
	}

}

func HighlightCuboidFace(bw float64, bh float64, bd float64, face int) {
	w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2
	// might need to use glPolygonOffset
	gl.Color4ub(255, 32, 32, 0)
	if face == UP_FACE {
		gl.Enable(gl.POLYGON_OFFSET_LINE)
		gl.PolygonOffset(-3, -1.0)
		gl.LineWidth(1.6)
		gl.Begin(gl.LINE_LOOP)
		gl.Normal3f(0.0, 1.0, 0.0)
		gl.Vertex3f(-d, h+0.8, -w) // Top Left Of The Texture and Quad
		gl.Vertex3f(d, h+0.8, -w)  // Top Right Of The Texture and Quad
		gl.Vertex3f(d, h+0.8, w)   // Bottom Right Of The Texture and Quad
		gl.Vertex3f(-d, h+0.8, w)  // Bottom Left Of The Texture and Quad

		gl.End()
		gl.Disable(gl.POLYGON_OFFSET_LINE)
	}
	gl.Color4ub(255, 128, 128, 128)
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
