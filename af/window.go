/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package af
import (
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/gl"
    "github.com/banthar/glu"
    "math"
)

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
func Reshape(width int, height int) {
    screenWidth = width
    screenHeight = height

    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    
    xmin, ymin := screenToView(0, 0)
    xmax, ymax := screenToView(uint16(width), uint16(height))
    
    gl.Ortho(float64(xmin), float64(xmax), float64(ymin), float64(ymax), -40, 40)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
}




func InitGraphics() {

    sdl.Init(sdl.INIT_VIDEO)


    screen := sdl.SetVideoMode(800, 600, 32, sdl.OPENGL|sdl.RESIZABLE)

    if screen == nil {
        sdl.Quit()
        panic("Couldn't set GL video mode: " + sdl.GetError() + "\n")
    } 

    if gl.Init() != 0 {
        panic("gl error")   
    }

    sdl.WM_SetCaption("Amberfell", "amberfell")


    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    // gl.ShadeModel(gl.FLAT)    
    gl.ShadeModel(gl.SMOOTH)    
    gl.Enable(gl.LIGHTING)
    gl.Enable(gl.LIGHT0)
    gl.Enable(gl.LIGHT1)


    // gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.EMISSION )
    // gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE )
    // gl.Enable ( gl.COLOR_MATERIAL )



    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-12.0, 12.0, -12.0, 12.0, -10, 10.0)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0)


    gl.ClearDepth(1.0)                         // Depth Buffer Setup
    gl.Enable(gl.DEPTH_TEST)                        // Enables Depth Testing
    gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)

    gl.Enable(gl.TEXTURE_2D)
    LoadMapTextures()
    LoadTerrainCubes()

    Reshape(int(screen.W), int(screen.H))
}

func QuitGraphics() {
    sdl.Quit()
}


func Draw(selectMode bool) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)


    if selectMode {
        gl.Disable(gl.TEXTURE_2D)
        gl.Disable(gl.FOG)
        gl.Disable(gl.LIGHTING)
        gl.Disable(gl.LIGHT0)
        gl.Disable(gl.LIGHT1)
        gl.Disable ( gl.COLOR_MATERIAL )
    } else {
        gl.Color4ub(255, 255, 255, 255)
        gl.Enable(gl.TEXTURE_2D)
        //gl.Enable(gl.FOG)
        gl.Enable(gl.LIGHTING)
        gl.Enable(gl.LIGHT0)
        gl.Enable ( gl.COLOR_MATERIAL )

        if timeOfDay < 5.3 || timeOfDay > 20.7 {
            gl.Enable(gl.LIGHT1)
        } else {
            gl.Disable(gl.LIGHT1)
        }


    }


    CheckGLError()
    gl.LoadIdentity()
    gl.Rotated(view_rotx, 1.0, 0.0, 0.0)
    gl.Rotated(-ThePlayer.Heading() + view_roty, 0.0, 1.0, 0.0)
    gl.Rotated(0, 0.0, 0.0, 1.0)


    center := ThePlayer.Position()




    // Sun
    var daylightIntensity float32 = 0.4
    if timeOfDay < 5 || timeOfDay > 21 {
        daylightIntensity = 0.00
        gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{0.1, 0.1, 0.1,1})
    } else if timeOfDay < 6 {
        daylightIntensity = 0.4 * (timeOfDay - 5)
    } else if timeOfDay > 20 {
        daylightIntensity = 0.4 * (21 - timeOfDay)
    }



    gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{0.5, 1, 1, 0})
    gl.Lightfv(gl.LIGHT0, gl.AMBIENT, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    // gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    // gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})

    // Torch
    ambient := float32(0.6)
    specular := float32(0.6)
    diffuse := float32(1)


    gl.Lightfv(gl.LIGHT1, gl.POSITION, []float32{0,1,0, 1})
    gl.Lightfv(gl.LIGHT1, gl.AMBIENT, []float32{ambient, ambient, ambient,1})
    gl.Lightfv(gl.LIGHT1, gl.SPECULAR, []float32{specular,specular,specular,1})
    gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, []float32{diffuse,diffuse,diffuse,1})
    gl.Lightf(gl.LIGHT1, gl.CONSTANT_ATTENUATION, 1.5)
    gl.Lightf(gl.LIGHT1, gl.LINEAR_ATTENUATION, 0.5)
    gl.Lightf(gl.LIGHT1, gl.QUADRATIC_ATTENUATION, 0.01)
    gl.Lightf(gl.LIGHT1, gl.SPOT_CUTOFF, 35)
    gl.Lightf(gl.LIGHT1, gl.SPOT_EXPONENT, 2.0)
    gl.Lightfv(gl.LIGHT1, gl.SPOT_DIRECTION, []float32{float32(math.Cos(ThePlayer.Heading() * math.Pi/180)),float32(-0.7), -float32(math.Sin(ThePlayer.Heading() * math.Pi/180))})

    CheckGLError()

    ThePlayer.Draw(center, selectMode)
    CheckGLError()

    TheWorld.Draw(center, selectMode)
    CheckGLError()

    // // var mousex, mousey int
    // // mouseState := sdl.GetMouseState(&mousex, &mousey)
    // gl.PushMatrix()
    // gl.Translatef(float32(center[XAXIS]),float32(center[YAXIS])-1,float32(center[ZAXIS]))
    // //print ("i:", i, "j:", j, "b:", World.At(i, j, groundLevel))
    // HighlightCuboidFace(1, 1, 1, TOP_FACE)
    // gl.PopMatrix()

    //gl.Translatef( 0.0, -20.0, -5.0 )

    if ShowOverlay {
        gl.PushMatrix();
        gl.LoadIdentity();
        gl.Color4f(0, 0, 0, 0.25);
        gl.Begin(gl.QUADS);
        gl.Vertex2f(0, 0);
        gl.Vertex2f(float32(screenWidth), 0);
        gl.Vertex2f(float32(screenWidth), float32(screenHeight));
        gl.Vertex2f(0, float32(screenHeight));
        gl.End();
        gl.PopMatrix();
    }

    if !selectMode {
        // gl.Finish()
        // gl.Flush()
        sdl.GL_SwapBuffers()
    }    
}