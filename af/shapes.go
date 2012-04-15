package af

import (
    "os"
    "image"
    _ "image/png"
    "github.com/banthar/gl"
)

var (
    MapTextures [16*10]gl.Texture
)

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

    for i:=0; i < 10; i++ {
        for j:=0; j < 16; j++ {
            rgba := image.NewRGBA(image.Rect(0, 0, pixels, pixels))
            for x := 0; x < pixels; x++ { 
                for y := 0; y < pixels; y++ { 
                    rgba.Set(x, y, img.At(pixels * j + x, pixels * i + y)) 
                } 
            }

            textureIndex := i*16 + j
            MapTextures[textureIndex] = gl.GenTexture()
            MapTextures[textureIndex].Bind(gl.TEXTURE_2D)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
            gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, pixels, pixels, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
            MapTextures[textureIndex].Unbind(gl.TEXTURE_2D)

        }
    }
}


func TerrainCube( n bool, s bool, w bool, e bool, u bool, d bool, texture byte, id uint16, selectMode bool) {
    var ftexture, btexture, ltexture, rtexture, utexture, dtexture byte = 0,0,0,0,0,0

    if n { rtexture = texture }
    if s { ltexture = texture }
    if e { btexture = texture }
    if w { ftexture = texture }
    if u { utexture = texture }
    if d { dtexture = texture }
    Cuboid(1, 1, 1, ftexture, btexture, ltexture, rtexture, utexture, dtexture, id, selectMode)

}




func Cuboid( bw float64, bh float64, bd float64, ftexture byte, btexture byte, ltexture byte, rtexture byte, utexture byte, dtexture byte, id uint16, selectMode bool) {
    w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2
    gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{0, 0, 0, 1});
    // Front face
    if ftexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), FRONT_FACE, 0)
        }
        MapTextures[ftexture].Bind(gl.TEXTURE_2D)
        gl.Begin(gl.QUADS)
        gl.Normal3f( 1.0, 0.0, 0.0)
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h,  w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.End()
    }

    // Back Face
    if btexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), BACK_FACE, 0)
        }
        MapTextures[btexture].Bind(gl.TEXTURE_2D)
        gl.Begin(gl.QUADS)
        gl.Normal3f( -1.0, 0.0, 0.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d, -h, -w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
        gl.End()
    }


    // Left Face
    if ltexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), LEFT_FACE, 0)
        }
        MapTextures[ltexture].Bind(gl.TEXTURE_2D)
        gl.Begin(gl.QUADS)
        gl.Normal3f( 0.0, 0.0, -1.0)
        gl.TexCoord2f(1.0, 0.0)        
        gl.Vertex3f(-d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Left Of The Texture and Quad
        gl.End()
    }

    // Right Face
    if rtexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), RIGHT_FACE, 0)
        } 
        MapTextures[rtexture].Bind(gl.TEXTURE_2D)
        gl.Begin(gl.QUADS)
        gl.Normal3f( 0.0, 0.0, 1.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( -d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(  d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(  d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( -d,  h,  w)  // Top Left Of The Texture and Quad
        gl.End()
    }


    // Top Face
    if utexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), TOP_FACE, 0)
        }
        MapTextures[utexture].Bind(gl.TEXTURE_2D)

        gl.Begin(gl.QUADS)
        // gl.Begin(gl.TRIANGLES)
        if utexture == BLOCK_STONE {
            //gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{1, 0.9, 0.9, 1});
        }

        gl.Normal3f( 0.0, 1.0, 0.0)

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
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.5)
        gl.Vertex3f(-d,  h,  0)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(0.5, 1.0)
        gl.Vertex3f( 0,  h, -w)  // Top Right Of The Texture and Quad

        gl.TexCoord2f(0.5, 1.0)
        gl.Vertex3f(0,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f(0,  h,  0)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.5)
        gl.Vertex3f( d,  h,  0)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f(0,  h, 0)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.5, 0.0)
        gl.Vertex3f(0,  h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d,  h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.5)
        gl.Vertex3f( d,  h, 0)  // Top Right Of The Texture and Quad

        gl.TexCoord2f(0.0, 0.5)
        gl.Vertex3f(-d,  h, 0)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(0.5, 0.0)
        gl.Vertex3f( 0,  h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f( 0,  h, 0)  // Top Right Of The Texture and Quad



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
    }

    // Bottom Face
    if dtexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), BOTTOM_FACE, 0)
        }
        MapTextures[dtexture].Bind(gl.TEXTURE_2D)
        gl.Begin(gl.QUADS)
        gl.Normal3f( 0.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d, -h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d, -h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.End()
    }
    
     

}



