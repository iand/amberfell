/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package af

import (
    "os"
    "image"
    _ "image/png"
    "github.com/banthar/gl"
)

var (
    MapTextures [16*10]gl.Texture
    TerrainCubes map[[6]byte]uint
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
            // gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
            // gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
            gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, pixels, pixels, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
            MapTextures[textureIndex].Unbind(gl.TEXTURE_2D)

        }
    }
}


func LoadTerrainCubes() {
    TerrainCubes = make(map[[6]byte]uint)
    var texture, faces byte
    for texture = 1; texture < 3; texture++ {
        for faces = 1; faces < 64; faces++ {
            listid := gl.GenLists(1)
            if listid == 0 { panic("GenLists return 0") }
            var ftexture, btexture, ltexture, rtexture, utexture, dtexture byte
            if faces & 32 == 32 { rtexture = texture }
            if faces & 16 == 16 { ltexture = texture }
            if faces &  8 ==  8 { btexture = texture }
            if faces &  4 ==  4 { ftexture = texture }
            if faces &  2 ==  2 { utexture = texture }
            if faces &  1 ==  1 { dtexture = texture }
    
            gl.NewList(listid, gl.COMPILE);
            Cuboid(1, 1, 1, ftexture, btexture, ltexture, rtexture, utexture, dtexture, 0, false)
            gl.EndList()
            CheckGLError()

            index := [6]byte{ftexture, btexture, ltexture, rtexture, utexture, dtexture}
            TerrainCubes[index] = listid
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

    if selectMode {
        Cuboid(1, 1, 1, ftexture, btexture, ltexture, rtexture, utexture, dtexture, id, selectMode)
    } else  {
        gl.CallList(TerrainCubes[[6]byte{ftexture, btexture, ltexture, rtexture, utexture, dtexture}])
    }

}




func Cuboid( bw float64, bh float64, bd float64, ftexture byte, btexture byte, ltexture byte, rtexture byte, utexture byte, dtexture byte, id uint16, selectMode bool) {
    w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2
    //gl.Materialfv(gl.FRONT, gl.EMISSION, []float32{0, 0, 0, 1});
    // Front face
    if ftexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), FRONT_FACE, 0)
        }
        MapTextures[ftexture].Bind(gl.TEXTURE_2D)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

        gl.EnableClientState( gl.VERTEX_ARRAY ); // Enable Vertex Arrays
        gl.EnableClientState( gl.TEXTURE_COORD_ARRAY ); // Enable Texture Coord Arrays
        gl.EnableClientState( gl.NORMAL_ARRAY ); // Enable Texture Coord Arrays
        gl.VertexPointer(3, 0, []float32{d, -h, -w, d,  h, -w, d,  h,  w,  d, -h,  w})
        gl.TexCoordPointer(2, 0, []float32{1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0})
        gl.NormalPointer(0, []float32{1.0, 0.0, 0.0})
        gl.DrawArrays( gl.QUADS, 0, 4 );

        // gl.Begin(gl.QUADS)
        // gl.Normal3f( 1.0, 0.0, 0.0)
        // gl.TexCoord2f(1.0, 0.0)
        // gl.Vertex3f( d, -h, -w)  // Bottom Right Of The Texture and Quad
        // gl.TexCoord2f(1.0, 1.0)
        // gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad
        // gl.TexCoord2f(0.0, 1.0)
        // gl.Vertex3f( d,  h,  w)  // Top Left Of The Texture and Quad
        // gl.TexCoord2f(0.0, 0.0)
        // gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad
        // gl.End()
        //MapTextures[ftexture].Unbind(gl.TEXTURE_2D)
        CheckGLError()
    }
    // Back Face
    if btexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), BACK_FACE, 0)
        }
        MapTextures[btexture].Bind(gl.TEXTURE_2D)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
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
        //MapTextures[btexture].Unbind(gl.TEXTURE_2D)

        CheckGLError()
    }


    // Left Face
    if ltexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), LEFT_FACE, 0)
        }
        MapTextures[ltexture].Bind(gl.TEXTURE_2D)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
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
        //MapTextures[ltexture].Unbind(gl.TEXTURE_2D)
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
        //MapTextures[rtexture].Unbind(gl.TEXTURE_2D)

        CheckGLError()
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
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
        
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 0.5)
        gl.Vertex3f(-d,  h,  0)  // Bottom Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f( 0,  h,  0)  // Bottom Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 1.0)
        gl.Vertex3f( 0,  h, -w)  // Top Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 1.0)
        gl.Vertex3f(0,  h, -w)  // Top Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f(0,  h,  0)  // Bottom Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(1.0, 0.5)
        gl.Vertex3f( d,  h,  0)  // Bottom Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 0.5)
        gl.Vertex3f(0,  h, 0)  // Top Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 0.0)
        gl.Vertex3f(0,  h,  w)  // Bottom Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d,  h,  w)  // Bottom Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(1.0, 0.5)
        gl.Vertex3f( d,  h, 0)  // Top Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 0.5)
        gl.Vertex3f(-d,  h, 0)  // Top Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.5, 0.0)
        gl.Vertex3f( 0,  h,  w)  // Bottom Right Of The Texture and Quad

        gl.Normal3f( 0.0, 1.0, 0.0)
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
        //MapTextures[utexture].Unbind(gl.TEXTURE_2D)
        CheckGLError()
    }

    // Bottom Face
    if dtexture != 0 {
        if selectMode {
            gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), BOTTOM_FACE, 0)
        }
        MapTextures[dtexture].Bind(gl.TEXTURE_2D)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
        gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
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
        //MapTextures[btexture].Unbind(gl.TEXTURE_2D)
        CheckGLError()
    }
    
     

}

func HighlightCuboidFace(bw float64, bh float64, bd float64, face int) {
    w, h, d := float32(bw)/2, float32(bh)/2, float32(bd)/2
    // might need to use glPolygonOffset
    gl.Color4ub(255, 32, 32, 0)
    if face == TOP_FACE {
        gl.Enable(gl.POLYGON_OFFSET_LINE)
        gl.PolygonOffset(-3, -1.0)
        gl.LineWidth(1.6);
        gl.Begin(gl.LINE_LOOP)
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.Vertex3f(-d,  h+0.8, -w)  // Top Left Of The Texture and Quad
        gl.Vertex3f( d,  h+0.8, -w)  // Top Right Of The Texture and Quad
        gl.Vertex3f( d,  h+0.8,  w)  // Bottom Right Of The Texture and Quad
        gl.Vertex3f(-d,  h+0.8,  w)  // Bottom Left Of The Texture and Quad

        gl.End()
        gl.Disable(gl.POLYGON_OFFSET_LINE)
    }
    gl.Color4ub(255, 128, 128, 128)
}

func CheckGLError() {
    glerr := gl.GetError()
    if glerr & gl.INVALID_ENUM == gl.INVALID_ENUM { println("gl error: INVALID_ENUM") }
    if glerr & gl.INVALID_VALUE != 0 { println("gl error: INVALID_VALUE") }
    if glerr & gl.INVALID_OPERATION != 0 { println("gl error: INVALID_OPERATION") }
    if glerr & gl.STACK_OVERFLOW != 0 { println("gl error: STACK_OVERFLOW") }
    if glerr & gl.STACK_UNDERFLOW != 0 { println("gl error: STACK_UNDERFLOW") }
    if glerr & gl.OUT_OF_MEMORY != 0 { println("gl error: OUT_OF_MEMORY") }
    if glerr & gl.TABLE_TOO_LARGE  != 0 { println("gl error: TABLE_TOO_LARGE ") }

    if glerr != gl.NO_ERROR {
        panic("Got an OpenGL Error")
    }


}