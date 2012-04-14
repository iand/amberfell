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


func Cube( n bool, s bool, w bool, e bool, u bool, d bool, texture byte, id uint16, selectMode bool) {
    MapTextures[texture].Bind(gl.TEXTURE_2D)
    
    const width = 0.5


    gl.Begin(gl.QUADS)                  // Start Drawing Quads

        if n {
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 0, 0)
            } 
            // Front Face
            gl.Normal3f( 0.0, 0.0, 1.0)
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-width, -width,  width)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( width, -width,  width)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( width,  width,  width)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-width,  width,  width)  // Top Left Of The Texture and Quad
        }

        if s {
            // Back Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 1, 0)
            }
            gl.Normal3f( 0.0, 0.0, -1.0)
            gl.TexCoord2f(1.0, 0.0)        
            gl.Vertex3f(-width, -width, -width)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-width,  width, -width)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( width,  width, -width)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( width, -width, -width)  // Bottom Left Of The Texture and Quad
        }

        //gl.Color3f(0.3,0.3,0.6)
        if w {
            // Right face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 2, 0)
            }
            gl.Normal3f( 1.0, 0.0, 0.0)
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( width, -width, -width)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( width,  width, -width)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( width,  width,  width)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( width, -width,  width)  // Bottom Left Of The Texture and Quad
        }

        if e {
            // Left Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 3, 0)
            }
            gl.Normal3f( -1.0, 0.0, 0.0)
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-width, -width, -width)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f(-width, -width,  width)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-width,  width,  width)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-width,  width, -width)  // Top Left Of The Texture and Quad
        }
    gl.End();   
    
    MapTextures[texture].Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        //gl.Color3f(0.3,1.0,0.3)
        if u {
            // Top Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 4, 0)
            }
            gl.Normal3f( 0.0, 1.0, 0.0)
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-width,  width, -width)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-width,  width,  width)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( width,  width,  width)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( width,  width, -width)  // Top Right Of The Texture and Quad
           }
     
        if d {
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 5, 0)
            }
            // Bottom Face
            gl.Normal3f( 0.0, -1.0, 0.0)
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-width, -width, -width)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( width, -width, -width)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( width, -width,  width)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f(-width, -width,  width)  // Bottom Right Of The Texture and Quad
        }

    gl.End();   

}
