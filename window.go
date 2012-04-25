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
	// "runtime"
)

func Draw(t int64) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	gl.Color4ub(192, 192, 192, 255)
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
	gl.Materialfv(gl.FRONT, gl.AMBIENT, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.DIFFUSE, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.SPECULAR, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.SHININESS, []float32{0.0, 0.0, 0.0, 1})
	var daylightIntensity float32 = 0.45
	var nighttimeIntensity float32 = 0.01
	if timeOfDay < 5 || timeOfDay > 21 {
		gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{0.2, 0.2, 0.2, 1})
		daylightIntensity = 0.01
	} else if timeOfDay < 6 {
		daylightIntensity = nighttimeIntensity + daylightIntensity*(timeOfDay-5)
	} else if timeOfDay > 20 {
		daylightIntensity = nighttimeIntensity + daylightIntensity*(21-timeOfDay)
	}

	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{daylightIntensity / 2.5, daylightIntensity / 2.5, daylightIntensity / 2.5, 1})

	gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{30 * float32(math.Sin(ThePlayer.Heading()*math.Pi/180)), 60, 30 * float32(math.Cos(ThePlayer.Heading()*math.Pi/180)), 0})
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})
	// gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
	// gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity * 2, daylightIntensity * 2, daylightIntensity * 2, 1})
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})

	// Torch
	ambient := float32(0.6)
	specular := float32(0.6)
	diffuse := float32(1)

	gl.Lightfv(gl.LIGHT1, gl.POSITION, []float32{float32(ThePlayer.position[XAXIS]), float32(ThePlayer.position[YAXIS] + 1), float32(ThePlayer.position[ZAXIS]), 1})
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

	// Draw HUD
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	blockscale := float32(0.4)
	gl.Disable(gl.DEPTH_TEST)
	radius := float32(90) * PIXEL_SCALE
	centrex := float32(viewport.rplane) - radius + blockscale*0.5 //+ 20 * PIXEL_SCALE
	centrey := float32(viewport.bplane) + radius - blockscale*0.5 //- 20 * PIXEL_SCALE
	//textures[TEXTURE_PICKER].Bind(gl.TEXTURE_2D)

	// gl.Begin(gl.QUADS)
	// gl.Normal3f(0.0, 0.0, 1.0)
	// gl.TexCoord2f(0.0, 0.0)
	// gl.Vertex3f(centrex - radius / 2, centrey + radius / 2, 10)
	// gl.TexCoord2f(0.0, 1.0)
	// gl.Vertex3f(centrex - radius / 2, centrey - radius / 2         , 10)
	// gl.TexCoord2f(1.0, 1.0)
	// gl.Vertex3f(centrex + radius / 2, centrey - radius / 2, 10)
	// gl.TexCoord2f(1.0, 0.0)
	// gl.Vertex3f(centrex + radius / 2, centrey + radius / 2, 10)
	// textures[TEXTURE_PICKER].Unbind(gl.TEXTURE_2D)

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 128)
	gl.Vertex2f(centrex, centrey)
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(centrex-float32(math.Sin(angle))*radius, centrey+float32(math.Cos(angle))*radius)
	}
	gl.End()

	selectionRadius := blockscale * 1.2
	actionItemRadius := radius - blockscale*1.5

	actionItemAngle := -(float64(ThePlayer.currentAction) - 1.5) * math.Pi / 4
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 228)
	gl.Vertex2f(centrex-actionItemRadius*float32(math.Sin(actionItemAngle)), centrey+actionItemRadius*float32(math.Cos(actionItemAngle)))
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(centrex-actionItemRadius*float32(math.Sin(actionItemAngle))-float32(math.Sin(angle))*selectionRadius, centrey+actionItemRadius*float32(math.Cos(actionItemAngle))+float32(math.Cos(angle))*selectionRadius)
	}
	gl.End()

	for i := 0; i < 5; i++ {

		item := ThePlayer.equippedItems[i]
		if item != ITEM_NONE {
			angle := -(float64(i) + 1.5) * math.Pi / 4
			gl.LoadIdentity()
			gl.Translatef(centrex-actionItemRadius*float32(math.Sin(angle)), centrey+actionItemRadius*float32(math.Cos(angle)), 12)
			gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 1.0, 0.0, 0.0)
			gl.Rotatef(360*float32(math.Cos(float64(t)/1e10+float64(i))), 0.0, 1.0, 0.0)
			gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 0.0, 0.0, 1.0)
			gl.Scalef(blockscale, blockscale, blockscale)
			TerrainCube(true, true, true, true, true, true, byte(item), FACE_NONE)
		}
	}

	// Draw debug console
	if DebugMode {
		h := float32(consoleFont.height) * PIXEL_SCALE
		margin := float32(3.0) * PIXEL_SCALE
		consoleHeight := 3 * h

		gl.MatrixMode(gl.MODELVIEW)

		gl.LoadIdentity()
		gl.Color4ub(0, 0, 0, 208)

		gl.Begin(gl.QUADS)
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

		gl.LoadIdentity()
		gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-3*h, 0)

		numgc := uint32(0)
		avggc := float64(0)
		var last3 [3]float64
		if metrics.mem.NumGC > 3 {
			numgc = metrics.mem.NumGC
			avggc = float64(metrics.mem.PauseTotalNs) / float64(metrics.mem.NumGC) / 1e6
			index := int(numgc) - 1
			if index > 255 {
				index = 255
			}

			last3[0] = float64(metrics.mem.PauseNs[index]) / 1e6
			last3[1] = float64(metrics.mem.PauseNs[index-1]) / 1e6
			last3[2] = float64(metrics.mem.PauseNs[index-2]) / 1e6
		}

		consoleFont.Print(fmt.Sprintf("Mem: %.1f/%.1f   GC: %.1fms [%d: %.1f, %.1f, %.1f]", float64(metrics.mem.Alloc)/(1024*1024), float64(metrics.mem.Sys)/(1024*1024), avggc, numgc, last3[0], last3[1], last3[2]))

	}

	gl.Finish()
	gl.Flush()
	sdl.GL_SwapBuffers()
	// runtime.GC()
}
