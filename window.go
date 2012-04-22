/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"math"
)

func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	gl.Color4ub(255, 255, 255, 255)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.DEPTH_TEST)
	//gl.Enable(gl.FOG)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.COLOR_MATERIAL)

	if timeOfDay < 5.3 || timeOfDay > 20.7 {
		gl.Enable(gl.LIGHT1)
	} else {
		gl.Disable(gl.LIGHT1)
	}

	// CheckGLError()
	gl.LoadIdentity()
	// gl.Rotated(0, 0.0, 0.0, 1.0)
	// gl.Rotated(viewport.rotx, 1.0, 0.0, 0.0)
	// gl.Rotated(viewport.roty, 0.0, 1.0, 0.0)
	// gl.Translatef(float32(viewport.transx), float32(viewport.transy), float32(viewport.transz))

	center := ThePlayer.Position()

	// matrix := *viewport.matrix.Float32()
	matrix := ModelMatrix().Float32()
	gl.MultMatrixf(&matrix[0])
	//gl.Translatef(-float32(center[XAXIS]), -float32(center[YAXIS]), -float32(center[ZAXIS]))

	// Sun
	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{0.4, 0.4, 0.4, 1})
	var daylightIntensity float32 = 0.4
	if timeOfDay < 5 || timeOfDay > 21 {
		daylightIntensity = 0.00
	} else if timeOfDay < 6 {
		daylightIntensity = 0.4 * (timeOfDay - 5)
	} else if timeOfDay > 20 {
		daylightIntensity = 0.4 * (21 - timeOfDay)
	}

	gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{30 * float32(math.Sin(ThePlayer.Heading()*math.Pi/180)), 60, 30 * float32(math.Cos(ThePlayer.Heading()*math.Pi/180)), 0})
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})
	// gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
	// gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})

	// Torch
	ambient := float32(0.6)
	specular := float32(0.6)
	diffuse := float32(1)

	gl.Lightfv(gl.LIGHT1, gl.POSITION, []float32{0, 1, 0, 1})
	gl.Lightfv(gl.LIGHT1, gl.AMBIENT, []float32{ambient, ambient, ambient, 1})
	gl.Lightfv(gl.LIGHT1, gl.SPECULAR, []float32{specular, specular, specular, 1})
	gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, []float32{diffuse, diffuse, diffuse, 1})
	gl.Lightf(gl.LIGHT1, gl.CONSTANT_ATTENUATION, 1.5)
	gl.Lightf(gl.LIGHT1, gl.LINEAR_ATTENUATION, 0.5)
	gl.Lightf(gl.LIGHT1, gl.QUADRATIC_ATTENUATION, 0.01)
	gl.Lightf(gl.LIGHT1, gl.SPOT_CUTOFF, 35)
	gl.Lightf(gl.LIGHT1, gl.SPOT_EXPONENT, 2.0)
	gl.Lightfv(gl.LIGHT1, gl.SPOT_DIRECTION, []float32{float32(math.Cos(ThePlayer.Heading() * math.Pi / 180)), float32(-0.7), -float32(math.Sin(ThePlayer.Heading() * math.Pi / 180))})

	// CheckGLError()
	gl.RenderMode(gl.RENDER)
	ThePlayer.Draw(center)
	// CheckGLError()

	// CheckGLError()

	selectedBlockFace := viewport.SelectedBlockFace()
	TheWorld.Draw(center, selectedBlockFace)

	if DebugMode {
		h := float32(consoleFont.height) * PIXEL_SCALE
		margin := float32(3.0) * PIXEL_SCALE
		consoleHeight := 3 * h

		gl.MatrixMode(gl.PROJECTION)
		//		gl.LoadIdentity ()
		//gl.Ortho (0, float64(viewport.screenWidth), float64(viewport.screenHeight), 0, 0, 1)
		gl.Disable(gl.DEPTH_TEST)
		gl.MatrixMode(gl.MODELVIEW)

		gl.LoadIdentity()
		gl.Color4ub(0, 0, 0, 208)

		gl.Begin(gl.QUADS)
		// gl.Vertex2f(float32(viewport.lplane)+0.5, float32(viewport.bplane)+0.5) // Bottom Left Of The Texture and Quad
		// gl.Vertex2f(float32(viewport.rplane)-0.5, float32(viewport.bplane)+0.5) // Bottom Right Of The Texture and Quad
		// gl.Vertex2f(float32(viewport.rplane)-0.5, float32(viewport.tplane)-0.5) // Top Right Of The Texture and Quad
		// gl.Vertex2f(float32(viewport.lplane)+0.5, float32(viewport.tplane)-0.5) // Top Left Of The Texture and Quad
		gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane)+consoleHeight+margin*2) // Bottom Left Of The Texture and Quad
		gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane)+consoleHeight+margin*2) // Bottom Right Of The Texture and Quad
		gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane))                        // Top Right Of The Texture and Quad
		gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane))                        // Top Left Of The Texture and Quad
		gl.End()

		gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-h, 0)
		consoleFont.Print(fmt.Sprintf("FPS: %5.2f", metrics.fps))
		gl.LoadIdentity()
		gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-2*h, 0)
		consoleFont.Print(fmt.Sprintf("X: %5.2f Y: %4.2f Z: %5.2f H: %5.2f (%s)", ThePlayer.position[XAXIS], ThePlayer.position[YAXIS], ThePlayer.position[ZAXIS], ThePlayer.heading, HeadingToCompass(ThePlayer.heading)))

	}

	gl.Finish()
	gl.Flush()
	sdl.GL_SwapBuffers()

}
