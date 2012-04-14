package main

import (
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/gl"
    "github.com/banthar/glu"
    "math"
    "math/rand"  
    "flag"
    "fmt"
    "github.com/iand/amberfell/af"
    "time"
    
)    

const piover180 = 0.0174532925



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
    wolf *af.Wolf
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

    wolf = new(af.Wolf)
    wolf.Init(0, 14, 14, af.GroundLevel+1)


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

                    Draw(true)
                    gl.ReadPixels(xv, yv, 1, 1, gl.RGBA, &data[0])
                    Draw(false)

                    fmt.Printf("pixel data: %d, %d, %d, %d\n", data[0], data[1], data[2], data[3])

                    id := uint16(data[0]) + uint16(data[1]) * 256
                    if id != 0 {
                        face := data[2]
                        dx, dy, dz := af.BlockIdToRelativeCoordinate(id)
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
            world.Simulate(float64(dt) / 1e9)

            computeFrame++
            t += dt
        }

        //interpolate(previous, current, accumulator/dt)


        Draw(false)
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
    af.LoadMapTextures()


}











func Draw(selectMode bool) {
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

    gl.Rotated(view_rotx, 1.0, 0.0, 0.0)
    gl.Rotated(-player.Heading() + view_roty, 0.0, 1.0, 0.0)
    gl.Rotated(0, 0.0, 0.0, 1.0)


    pos := player.Position()
    player.Draw(pos, selectMode)
    world.Draw(pos, selectMode)

    if !selectMode {
        sdl.GL_SwapBuffers()
        gl.Finish()
    }    
}