package main


/*

    this is go version based on SDL version

    this version uses Go-SDL: https://github.com/banthar/Go-SDL

*/

import (
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/gl"
    "github.com/banthar/glu"
    "image"
    _ "image/png"
    "os"
    "math"
    "math/rand"  
    "flag"
    "fmt"
    "amberfell/af"
    "time"
    
)    

const piover180 = 0.0174532925
const blockSize = 1.0



var printInfo = flag.Bool("info", false, "print GL implementation information")

var T0 uint32 = 0
var Frames uint32 = 0


var view_rotx float64 = 50.0
var view_roty float64 = 50.0
var view_rotz float64 = 0.0
var gear1, gear2, gear3 uint
var angle float64 = 0.0


var (
    player *af.Player
    mapTextures [16*10]gl.Texture
    world af.World
    DebugMode bool
    screenWidth, screenHeight int
    tileWidth = 48
    screenScale int = 5 * tileWidth / 2
)


    


func main() {
    flag.Parse()
    rand.Seed(71)   
    var done bool
    var keys []uint8
    player = new(af.Player)
    player.Init(0, 10, 10, af.GroundLevel+1)
    world.Init(56,56,10)
    
    sdl.Init(sdl.INIT_VIDEO)

    var screen = sdl.SetVideoMode(800, 600, 32, sdl.OPENGL|sdl.RESIZABLE)

    if screen == nil {
        sdl.Quit()
        panic("Couldn't set GL video mode: " + sdl.GetError() + "\n")
    } 

    if gl.Init() != 0 {
        panic("gl error")   
    }

    sdl.WM_SetCaption("Amberfell", "amberfell")

    init2()
    reshape(int(screen.W), int(screen.H))

    var currentTime, accumulator int64 = 0, 0
    var t, dt int64 = 0, 1e9/40
    var drawFrame, computeFrame int64 = 0, 0
    fps := new(af.Timer)
    fps.Start()

    update := new(af.Timer)
    update.Start()

    done = false
    for !done {
        // controlForce := af.Vector{0, 0, 0}

        var vx, vz float64

        for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
            switch e.(type) {
            case *sdl.ResizeEvent:
                re := e.(*sdl.ResizeEvent)
                screen = sdl.SetVideoMode(int(re.W), int(re.H), 16,
                    sdl.OPENGL|sdl.RESIZABLE)
                if screen != nil {
                    reshape(int(screen.W), int(screen.H))
                } else {
                    panic("we couldn't set the new video mode??")
                }
                break

            case *sdl.MouseButtonEvent:
                re := e.(*sdl.MouseButtonEvent)
                if re.Button == 1 && re.State == 1 { // LEFT, DOWN
                    // println("Click:", re.X, re.Y, re.State, re.Button, re.Which)

                    // MOUSEBUTTONDOWNMASK
                    xv, yv := int(re.X), screenHeight - int(re.Y)
                    data := [4]uint8{0, 0, 0, 0}

                    draw2(true)
                    gl.ReadPixels(xv, yv, 1, 1, gl.RGBA, &data[0])
                    draw2(false)

                    fmt.Printf("pixel data: %d, %d, %d, %d\n", data[0], data[1], data[2], data[3])

                    id := uint16(data[0]) + uint16(data[1]) * 256
                    if id != 0 {
                        face := data[2]
                        dx, dy, dz := blockIdToRelativeCoordinate(id)
                        fmt.Printf("id: %d, dx: %d, dy: %d, dz: %d, face: %d\n", id, dx, dy, dz, face)
                        if ! (dx == 0 && dy == 0 && dz == 0) {
                            pos := player.IntPosition()
                            pos[af.XAXIS] += dx
                            pos[af.YAXIS] += dy
                            pos[af.ZAXIS] += dz
                            if face == 4 { // top
                                pos[af.YAXIS]++
                            } else if face == 5 { // bottom
                                pos[af.YAXIS]--
                            } else if face == 0 { // front
                                pos[af.ZAXIS]++
                            } else if face == 1 { // back
                                pos[af.ZAXIS]--
                            } else if face == 2 { // left
                                pos[af.XAXIS]++
                            } else if face == 3 { // right
                                pos[af.XAXIS]--
                            }
                            world.Set(pos[af.XAXIS], pos[af.YAXIS], pos[af.ZAXIS], 2)
                        }
                    }
                }




            case *sdl.QuitEvent:
                done = true
                break
            }
        }
        keys = sdl.GetKeyState()

        if keys[sdl.K_ESCAPE] != 0 {
            done = true
        }
        if keys[sdl.K_UP] != 0 {
            view_rotx += 5.0
            if view_rotx > 80 {
                view_rotx = 80
            }
        }
        if keys[sdl.K_DOWN] != 0 {
            view_rotx -= 5.0
            if view_rotx < 15 {
                view_rotx = 10
            }
        }
        if keys[sdl.K_LEFT] != 0 {
            view_roty += 9
            //println("view_roty:", view_roty)
        }
        if keys[sdl.K_RIGHT] != 0 {
            view_roty -= 9
        }
        if keys[sdl.K_w] != 0 {
            if !player.IsFalling() {
                vx = math.Cos(player.Heading() * math.Pi / 180)
                vz = -math.Sin(player.Heading() * math.Pi / 180)
            }

        }
        if keys[sdl.K_s] != 0 {
            if !player.IsFalling() {
                vx = -math.Cos(player.Heading() * math.Pi / 180)
                vz = math.Sin(player.Heading() * math.Pi / 180)
            }
     
        }
        if keys[sdl.K_a] != 0 {
            player.Rotate(9)

        }        
        if keys[sdl.K_SPACE] != 0 {
            if !player.IsFalling() {
                player.Accelerate(af.Vector{0, 7, 0})
            }
        } 
        if keys[sdl.K_d] != 0 {
            player.Rotate(-9)
        }        
        if keys[sdl.K_z] != 0 {
            if (sdl.GetModState() & sdl.KMOD_RSHIFT) != 0 {
                view_rotz -= 5.0
            } else {
                view_rotz += 5.0
            }
        }
        if keys[sdl.K_F3] != 0 {
            if DebugMode == true {
                DebugMode = false
            } else {
                DebugMode = true
            }
        }               

        if DebugMode {
            fmt.Printf("x:%f, z:%f\n", player.X(), player.Z())
        }

        if vx != 0 || vz != 0 {
            player.Setvx(10 * vx)
            player.Setvz(10 * vz)
        } else {
            if !player.IsFalling() {
                player.Setvx(player.Velocity()[af.XAXIS] / 2.5)
                player.Setvz(player.Velocity()[af.ZAXIS] / 2.5)
            } else {
                player.Setvx(player.Velocity()[af.XAXIS] / 1.04)
                player.Setvz(player.Velocity()[af.ZAXIS] / 1.04)

            }


        }



        newTime := time.Now().UnixNano()
        deltaTime := newTime - currentTime
        currentTime = newTime
        if deltaTime > 1e9/4 {
            deltaTime = 1e9/4
        }

        accumulator += deltaTime

        for accumulator > dt {
            accumulator -= dt
            // player.ZeroForces()
            // player.ApplyForce(controlForce)
            world.ApplyForces(player, float64(dt) / 1e9)
            player.Update(float64(dt) / 1e9)

            computeFrame++
            t += dt
        }

        //interpolate(previous, current, accumulator/dt)

        draw2(false)
        drawFrame++

        if update.GetTicks() > 1e9/2 {
            //fmt.Printf("draw fps: %f\n", float64(drawFrame) / (float64(update.GetTicks()) / float64(1e9)) )
            //fmt.Printf("compute fps: %f\n", float64(computeFrame) / (float64(update.GetTicks()) / float64(1e9)) )
            drawFrame, computeFrame = 0, 0
            update.Start()
        }

        //sdl.Delay( 1000 )
    }
    sdl.Quit()
    return

}

func screenToView(xs uint16, ys uint16) (xv float64, yv float64) {
    // xs = 0 => -float64(screenWidth) / screenScale
    // xs = screenWidth => float64(screenWidth) / screenScale

    viewWidth := 2 * float64(screenWidth) / float64(screenScale)
    xv = (-viewWidth / 2 + viewWidth * float64(xs) / float64(screenWidth))

    viewHeight := 2 * float64(screenHeight) / float64(screenScale)
    yv = (-viewHeight / 2 + viewHeight * float64(ys) / float64(screenHeight))

    return
}

/* new window size or exposure */
func reshape(width int, height int) {
    screenWidth = width
    screenHeight = height

    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    
    xmin, ymin := screenToView(0, 0)
    xmax, ymax := screenToView(uint16(width), uint16(height))
    
    gl.Ortho(float64(xmin), float64(xmax), float64(ymin), float64(ymax), -100, 100)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
}


func init2() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
//    gl.ShadeModel(gl.SMOOTH)    
    gl.Enable(gl.LIGHTING)
    gl.Enable(gl.LIGHT0)
    gl.Lightfv(0, gl.AMBIENT, []float32{0.5,0.5,0.5,1})
    gl.Lightfv(0, gl.DIFFUSE, []float32{1,1,1,1})
    gl.Lightfv(0, gl.SPECULAR, []float32{1,1,1,0.5})
    gl.Lightfv(0, gl.POSITION, []float32{-5.0, 5.0, 10.0, 0})
    gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.EMISSION )
    gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE )
    gl.Enable ( gl.COLOR_MATERIAL )



    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-12.0, 12.0, -12.0, 12.0, -10, 10.0)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0)


    gl.ClearDepth(1.0)                         // Depth Buffer Setup
    gl.Enable(gl.DEPTH_TEST)                        // Enables Depth Testing
    gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

    gl.Enable(gl.TEXTURE_2D)
    loadMapTextures()


}


func draw2(selectMode bool) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
    if selectMode {
        gl.Disable(gl.TEXTURE_2D)
        gl.Disable(gl.FOG)
        gl.Disable(gl.LIGHTING)
        gl.Disable(gl.LIGHT0)
        gl.Disable ( gl.COLOR_MATERIAL )
    } else {
        gl.Color4ub(255, 255, 255, 255)
        gl.Enable(gl.TEXTURE_2D)
        //gl.Enable(gl.FOG)
        gl.Enable(gl.LIGHTING)
        gl.Enable(gl.LIGHT0)
        gl.Enable ( gl.COLOR_MATERIAL )
    }



    gl.LoadIdentity()
    gl.PushMatrix()

    gl.Rotated(view_rotx, 1.0, 0.0, 0.0)
    gl.Rotated(-player.Heading() + view_roty, 0.0, 1.0, 0.0)
    gl.Rotated(0, 0.0, 0.0, 1.0)



    gl.PushMatrix()
    //stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
    gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)
    drawPlayer(selectMode)
    gl.PopMatrix()

    gl.Translatef(-player.X() * blockSize, -player.Y() * blockSize, -player.Z() * blockSize)


    ip := player.IntPosition()
    var x, y, z int16
    for x =0; x < world.W; x++ {
        //gl.Translatef(3.0,0.0,-30.0)

        for z=0; z < world.D; z++ {
            for y=0; y < world.H; y++ {
                var terrain byte = world.At(x, y, z)
                if terrain != 0 {
                    var n, s, w, e, u, d bool = world.AirNeighbours(x, z, y)
                    var id uint16 = 0

                    dx := x - ip[af.XAXIS]
                    dy := y - ip[af.YAXIS]
                    dz := z - ip[af.ZAXIS]
                    if dx >= -2 && dx <= 2 && dy >= -2 && dy <= 2 && dz >= -2 && dz <= 2 {
                        id = relativeCoordinateToBlockId(dx, dy, dz)
                    }
                    gl.PushMatrix()
                    gl.Translatef(float32(x) * blockSize,float32(y) * blockSize,float32(z) * blockSize)
                    //print ("i:", i, "j:", j, "b:", world.At(i, j, groundLevel))
                    cube(n, s, w, e, u, d, terrain, id, selectMode)
                    gl.PopMatrix()
                }
            }
        }
    }

    gl.PopMatrix()
    if !selectMode {
        sdl.GL_SwapBuffers()
        gl.Finish()
    }


}


func drawPlayer(selectMode bool) {

    var w,h,d float32 = blockSize, blockSize, blockSize

    // topTexture.Bind(gl.TEXTURE_2D)
    // gl.Begin(gl.QUADS)                  // Start Drawing Quads

    //     // Front Face
    //     //gl.Color3f(0.5,0.5,1.0)              // Set The Color To Blue One Time Only
    //     gl.Normal3f( 0.0, 0.0, 1.0)
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f( -w, -h,  d)  // Bottom Left Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 0.0)
    //     gl.Vertex3f(  w, -h,  d)  // Bottom Right Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f(  w,  h,  d)  // Top Right Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f( -w,  h,  d)  // Top Left Of The Texture and Quad

    // gl.End();   
    
    // dirtTexture.Bind(gl.TEXTURE_2D)
    // gl.Begin(gl.QUADS)                  // Start Drawing Quads
    //        // Back Face
    //     gl.Normal3f( 0.0, 0.0, -1.0)
    //     gl.TexCoord2f(1.0, 0.0)        
    //     gl.Vertex3f(-w, -h, -d)  // Bottom Right Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f(-w,  h, -d)  // Top Right Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f( w,  h, -d)  // Top Left Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f( w, -h, -d)  // Bottom Left Of The Texture and Quad

    //     //gl.Color3f(0.3,0.3,0.6)
    //     // Right face
    //     gl.Normal3f( 1.0, 0.0, 0.0)
    //     gl.TexCoord2f(1.0, 0.0)
    //     gl.Vertex3f( w, -h, -d)  // Bottom Right Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f( w,  h, -d)  // Top Right Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f( w,  h,  d)  // Top Left Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f( w, -h,  d)  // Bottom Left Of The Texture and Quad

    //     // Left Face
    //     gl.Normal3f( -1.0, 0.0, 0.0)
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f(-w, -h, -d)  // Bottom Left Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 0.0)
    //     gl.Vertex3f(-w, -h,  d)  // Bottom Right Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f(-w,  h,  d)  // Top Right Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f(-w,  h, -d)  // Top Left Of The Texture and Quad

    //  //gl.Color3f(0.3,1.0,0.3)
    //     // Top Face
    //     gl.Normal3f( 0.0, 1.0, 0.0)
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f(-w,  h, -d)  // Top Left Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f(-w,  h,  d)  // Bottom Left Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 0.0)
    //     gl.Vertex3f( w,  h,  d)  // Bottom Right Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f( w,  h, -d)  // Top Right Of The Texture and Quad

    //     // Bottom Face
    //     gl.Normal3f( 0.0, -1.0, 0.0)
    //     gl.TexCoord2f(1.0, 1.0)
    //     gl.Vertex3f(-w, -h, -d)  // Top Right Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 1.0)
    //     gl.Vertex3f( w, -h, -d)  // Top Left Of The Texture and Quad
    //     gl.TexCoord2f(0.0, 0.0)
    //     gl.Vertex3f( w, -h,  d)  // Bottom Left Of The Texture and Quad
    //     gl.TexCoord2f(1.0, 0.0)
    //     gl.Vertex3f(-w, -h,  d)  // Bottom Right Of The Texture and Quad

    // gl.End();   


    h = float32(player.H()) * blockSize / 2
    gl.Translatef(0.0, h - blockSize / 2 ,0.0)
    w = float32(player.W()) * blockSize / 2
    d = float32(player.D()) * blockSize / 2
    //gl.Translatef(0.0,-h,0.0)
    mapTextures[33].Bind(gl.TEXTURE_2D)
    //topTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        //gl.Color3f(0.3,0.3,0.6)
        // Front face
        gl.Normal3f( 1.0, 0.0, 0.0)
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( w, -h, -d)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( w,  h, -d)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( w,  h,  d)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( w, -h,  d)  // Bottom Left Of The Texture and Quad

    gl.End()

    mapTextures[32].Bind(gl.TEXTURE_2D)

    // dirtTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        // Left Face
        gl.Normal3f( 0.0, 0.0, -1.0)
        gl.TexCoord2f(1.0, 0.0)        
        gl.Vertex3f(-w, -h, -d)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-w,  h, -d)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( w,  h, -d)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( w, -h, -d)  // Bottom Left Of The Texture and Quad


        // Right Face
        //gl.Color3f(0.5,0.5,1.0)              // Set The Color To Blue One Time Only
        gl.Normal3f( 0.0, 0.0, 1.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( -w, -h,  d)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(  w, -h,  d)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(  w,  h,  d)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( -w,  h,  d)  // Top Left Of The Texture and Quad


        // Back Face
        gl.Normal3f( -1.0, 0.0, 0.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-w, -h, -d)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-w, -h,  d)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-w,  h,  d)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-w,  h, -d)  // Top Left Of The Texture and Quad

     //gl.Color3f(0.3,1.0,0.3)
        // Top Face
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-w,  h, -d)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-w,  h,  d)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( w,  h,  d)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( w,  h, -d)  // Top Right Of The Texture and Quad

        // Bottom Face
        gl.Normal3f( 0.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-w, -h, -d)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( w, -h, -d)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( w, -h,  d)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-w, -h,  d)  // Bottom Right Of The Texture and Quad

    gl.End();   
}


func cube( n bool, s bool, w bool, e bool, u bool, d bool, texture byte, id uint16, selectMode bool) {
    mapTextures[texture].Bind(gl.TEXTURE_2D)
    
    gl.Begin(gl.QUADS)                  // Start Drawing Quads

        if n {
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 0, 0)
            } 
            // Front Face
            gl.Normal3f( 0.0, 0.0, 1.0)
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( blockSize/2,  blockSize/2,  blockSize/2)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2,  blockSize/2)  // Top Left Of The Texture and Quad
        }

        if s {
            // Back Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 1, 0)
            }
            gl.Normal3f( 0.0, 0.0, -1.0)
            gl.TexCoord2f(1.0, 0.0)        
            gl.Vertex3f(-blockSize/2, -blockSize/2, -blockSize/2)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2, -blockSize/2)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( blockSize/2,  blockSize/2, -blockSize/2)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( blockSize/2, -blockSize/2, -blockSize/2)  // Bottom Left Of The Texture and Quad
        }

        //gl.Color3f(0.3,0.3,0.6)
        if w {
            // Right face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 2, 0)
            }
            gl.Normal3f( 1.0, 0.0, 0.0)
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( blockSize/2, -blockSize/2, -blockSize/2)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( blockSize/2,  blockSize/2, -blockSize/2)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( blockSize/2,  blockSize/2,  blockSize/2)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Left Of The Texture and Quad
        }

        if e {
            // Left Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 3, 0)
            }
            gl.Normal3f( -1.0, 0.0, 0.0)
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-blockSize/2, -blockSize/2, -blockSize/2)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f(-blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2,  blockSize/2)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2, -blockSize/2)  // Top Left Of The Texture and Quad
        }
    gl.End();   
    
    mapTextures[texture].Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        //gl.Color3f(0.3,1.0,0.3)
        if u {
            // Top Face
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 4, 0)
            }
            gl.Normal3f( 0.0, 1.0, 0.0)
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2, -blockSize/2)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f(-blockSize/2,  blockSize/2,  blockSize/2)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f( blockSize/2,  blockSize/2,  blockSize/2)  // Bottom Right Of The Texture and Quad
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f( blockSize/2,  blockSize/2, -blockSize/2)  // Top Right Of The Texture and Quad
           }
     
        if d {
            if selectMode {
                gl.Color4ub(uint8(id & 0x00FF), uint8(id & 0xFF00 >> 8), 5, 0)
            }
            // Bottom Face
            gl.Normal3f( 0.0, -1.0, 0.0)
            gl.TexCoord2f(1.0, 1.0)
            gl.Vertex3f(-blockSize/2, -blockSize/2, -blockSize/2)  // Top Right Of The Texture and Quad
            gl.TexCoord2f(0.0, 1.0)
            gl.Vertex3f( blockSize/2, -blockSize/2, -blockSize/2)  // Top Left Of The Texture and Quad
            gl.TexCoord2f(0.0, 0.0)
            gl.Vertex3f( blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Left Of The Texture and Quad
            gl.TexCoord2f(1.0, 0.0)
            gl.Vertex3f(-blockSize/2, -blockSize/2,  blockSize/2)  // Bottom Right Of The Texture and Quad
        }

    gl.End();   

}

func loadMapTextures() {
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
            mapTextures[textureIndex] = gl.GenTexture()
            mapTextures[textureIndex].Bind(gl.TEXTURE_2D)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
            gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
            gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, pixels, pixels, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
            mapTextures[textureIndex].Unbind(gl.TEXTURE_2D)

        }
    }
}

// relative coordinates range from -3 to +3
func relativeCoordinateToBlockId(dx int16, dy int16, dz int16) (id uint16) {
    id =  0
    id |= uint16(dx + 3)
    id |= uint16(dy + 3) << 3
    id |= uint16(dz + 3) << 6
    return 
}   

func blockIdToRelativeCoordinate(id uint16) (dx int16, dy int16, dz int16) {
    dx = int16(id & 0x0007 - 3)
    dy = int16((id & 0x0038) >> 3 - 3)
    dz = int16((id & 0x01C0) >> 6 - 3)
    return
}