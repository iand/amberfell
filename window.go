/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	// "github.com/banthar/glu"
	"math"
	// "fmt"
)




func InitGraphics() {

	viewport.Zoomstd()
	viewport.Rotx(25)
	viewport.Roty(70)
	// viewport.Transx(-float64(ThePlayer.X()))
	// viewport.Transy(-float64(ThePlayer.Y()))
	// viewport.Transz(-float64(ThePlayer.Z()))

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
	// glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0)

	gl.ClearDepth(1.0)       // Depth Buffer Setup
	gl.Enable(gl.DEPTH_TEST) // Enables Depth Testing
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)

	gl.Enable(gl.TEXTURE_2D)
	LoadMapTextures()
	//LoadTerrainCubes()
	InitTerrainBlocks()

	viewport.Reshape(int(screen.W), int(screen.H))
}



func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	gl.Color4ub(255, 255, 255, 255)
	gl.Enable(gl.TEXTURE_2D)
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
	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{0.1, 0.1, 0.1, 1})
	var daylightIntensity float32 = 0.4
	if timeOfDay < 5 || timeOfDay > 21 {
		daylightIntensity = 0.00
	} else if timeOfDay < 6 {
		daylightIntensity = 0.4 * (timeOfDay - 5)
	} else if timeOfDay > 20 {
		daylightIntensity = 0.4 * (21 - timeOfDay)
	}

	gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{0.5, 1, 1, 0})
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

	TheWorld.Draw(center)
	// CheckGLError()

	// // var mousex, mousey int
	// // mouseState := sdl.GetMouseState(&mousex, &mousey)
	// gl.PushMatrix()
	// gl.Translatef(float32(center[XAXIS]),float32(center[YAXIS])-1,float32(center[ZAXIS]))
	// //print ("i:", i, "j:", j, "b:", World.At(i, j, groundLevel))
	// HighlightCuboidFace(1, 1, 1, TOP_FACE)
	// gl.PopMatrix()

	// if ShowOverlay {
	// 	gl.PushMatrix()
	// 	gl.LoadIdentity()
	// 	gl.Color4f(0, 0, 0, 0.25)
	// 	gl.Begin(gl.QUADS)
	// 	gl.Vertex2f(0, 0)
	// 	gl.Vertex2f(float32(screenWidth), 0)
	// 	gl.Vertex2f(float32(screenWidth), float32(screenHeight))
	// 	gl.Vertex2f(0, float32(screenHeight))
	// 	gl.End()
	// 	gl.PopMatrix()
	// }

	// var pm32 []float32 = make([]float32, 16)
	// gl.GetFloatv(gl.PROJECTION_MATRIX, pm32)
	// var projectionMatrix64 *Matrix4 = NewMatrix(float64(pm32[0]),float64(pm32[1]),float64(pm32[2]),float64(pm32[3]),float64(pm32[4]),float64(pm32[5]),float64(pm32[6]),float64(pm32[7]),float64(pm32[8]),float64(pm32[9]),float64(pm32[10]),float64(pm32[11]),float64(pm32[12]),float64(pm32[13]),float64(pm32[14]),float64(pm32[15]))

	// inverseMatrix, _ := projectionMatrix64.Multiply(ModelMatrix()).Inverse()

	// x := (float64(mousex)-float64(screenWidth)/2) / ( float64(screenWidth)/2 )
	// z := (float64(screenHeight)/2 - float64(mousey)) / ( float64(screenHeight)/2 )

	// origin := inverseMatrix.Transform(&Vectorf{x, z, 1}, 1)
	// norm := inverseMatrix.Transform(&Vectorf{0, 0, -1}, 0).Normalize()

	// fmt.Printf("Ray origin: %f, %f, %f\n", origin[0], origin[1], origin[2])
	// fmt.Printf("Ray norm: %f, %f, %f\n", norm[0], norm[1], norm[2])
	var pm32 []float32 = make([]float32, 16)
	var mousex, mousey int
	_ = sdl.GetMouseState(&mousex, &mousey)

	gl.GetFloatv(gl.PROJECTION_MATRIX, pm32)
	var projectionMatrix64 *Matrix4 = NewMatrix(float64(pm32[0]), float64(pm32[1]), float64(pm32[2]), float64(pm32[3]), float64(pm32[4]), float64(pm32[5]), float64(pm32[6]), float64(pm32[7]), float64(pm32[8]), float64(pm32[9]), float64(pm32[10]), float64(pm32[11]), float64(pm32[12]), float64(pm32[13]), float64(pm32[14]), float64(pm32[15]))

	inverseMatrix, _ := projectionMatrix64.Multiply(ModelMatrix()).Inverse()

	x := (float64(mousex) - float64(viewport.screenWidth)/2) / (float64(viewport.screenWidth) / 2)
	z := (float64(viewport.screenHeight)/2 - float64(mousey)) / (float64(viewport.screenHeight) / 2)

	origin := inverseMatrix.Transform(&Vectorf{x, z, -1}, 1)
	norm := inverseMatrix.Transform(&Vectorf{0, 0, 1}, 0).Normalize()

	if origin != nil && norm != nil {
		pos := IntPosition(ThePlayer.position)
		ray := Ray{origin, norm}
		// for dy := int16(5); dy > -6; dy-- {
		// 	for dz := int16(-5); dz < 6; dz++ {
		// 		for dx := int16(-5); dx < 6; dx++ {
		// 			if TheWorld.At(pos[XAXIS]+dx, pos[YAXIS]+dy, pos[ZAXIS]+dz) != BLOCK_AIR {
		// 				box := Box{
		// 						&Vectorf{float64(pos[XAXIS]+dx)-0.5, float64(pos[YAXIS]+dy)-0.5,float64(pos[ZAXIS]+dz)-0.5}, 
		// 						&Vectorf{float64(pos[XAXIS]+dx)+0.5, float64(pos[YAXIS]+dy)+0.5,float64(pos[ZAXIS]+dz)+0.5} }
		// 				vcenter := Vectorf{float64(pos[XAXIS]), float64(pos[YAXIS]),float64(pos[ZAXIS])}
		// 		        item := &BoxDistance{
		// 		            box:    box,
		// 		            distance: math.Sqrt(math.Pow(vcenter[0]-origin[0], 2) + math.Pow(vcenter[1]-origin[1], 2) + math.Pow(vcenter[2]-origin[2], 2)),
		// 		        }												
		// 		        heap.Push(&testBlocks, item)
		// 		    }
		// 		}
		// 	}
		// }
		// // See http://www.dyn-lab.com/articles/pick-selection.html

		var box *Box = nil
		distance := float64(1e9)
		face := int(0)
		for dy := int16(5); dy > -6; dy-- {
			for dz := int16(-5); dz < 6; dz++ {
				for dx := int16(-5); dx < 6; dx++ {
					trialDistance := math.Sqrt(math.Pow(float64(pos[XAXIS]+dx)-origin[0], 2) + math.Pow(float64(pos[YAXIS]+dy)-origin[1], 2) + math.Pow(float64(pos[ZAXIS]+dz)-origin[2], 2))
					if trialDistance < distance {
						if TheWorld.At(pos[XAXIS]+dx, pos[YAXIS]+dy, pos[ZAXIS]+dz) != BLOCK_AIR {
							trialBox := &Box{
								&Vectorf{float64(pos[XAXIS]+dx) - 0.5, float64(pos[YAXIS]+dy) - 0.5, float64(pos[ZAXIS]+dz) - 0.5},
								&Vectorf{float64(pos[XAXIS]+dx) + 0.5, float64(pos[YAXIS]+dy) + 0.5, float64(pos[ZAXIS]+dz) + 0.5}}

							hit, trialFace := ray.HitsBox(trialBox)
							if hit /*&& TheWorld.AirNeighbour(pos[XAXIS]+dx, pos[YAXIS]+dy, pos[ZAXIS]+dz, face)*/ {
								distance = trialDistance
								box = trialBox
								face = trialFace
							}

						}
					}
				}
			}
		}

		if box != nil {
			gl.PushMatrix()
			gl.LineWidth(4)
			gl.Color4ub(255, 0, 0, 32)
			gl.Begin(gl.QUADS)
			if face == UP_FACE {
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.max[YAXIS]+0.04), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.max[YAXIS]+0.04), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.max[YAXIS]+0.04), float32(box.max[ZAXIS]))
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.max[YAXIS]+0.04), float32(box.max[ZAXIS]))
			} else if face == EAST_FACE {
				gl.Vertex3f(float32(box.max[XAXIS]+0.04), float32(box.min[YAXIS]), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]+0.04), float32(box.max[YAXIS]), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]+0.04), float32(box.max[YAXIS]), float32(box.max[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]+0.04), float32(box.min[YAXIS]), float32(box.max[ZAXIS]))
			} else if face == WEST_FACE {
				gl.Vertex3f(float32(box.min[XAXIS]-0.04), float32(box.min[YAXIS]), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.min[XAXIS]-0.04), float32(box.min[YAXIS]), float32(box.max[ZAXIS]))
				gl.Vertex3f(float32(box.min[XAXIS]-0.04), float32(box.max[YAXIS]), float32(box.max[ZAXIS]))
				gl.Vertex3f(float32(box.min[XAXIS]-0.04), float32(box.max[YAXIS]), float32(box.min[ZAXIS]))
			} else if face == NORTH_FACE {
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.min[YAXIS]), float32(box.min[ZAXIS]-0.04))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.min[YAXIS]), float32(box.min[ZAXIS]-0.04))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.max[YAXIS]), float32(box.min[ZAXIS]-0.04))
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.max[YAXIS]), float32(box.min[ZAXIS]-0.04))
			} else if face == SOUTH_FACE {
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.min[YAXIS]), float32(box.max[ZAXIS]+0.04))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.min[YAXIS]), float32(box.max[ZAXIS]+0.04))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.max[YAXIS]), float32(box.max[ZAXIS]+0.04))
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.max[YAXIS]), float32(box.max[ZAXIS]+0.04))
			} else if face == DOWN_FACE {
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.min[YAXIS]-0.04), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.min[YAXIS]-0.04), float32(box.min[ZAXIS]))
				gl.Vertex3f(float32(box.max[XAXIS]), float32(box.min[YAXIS]-0.04), float32(box.max[ZAXIS]))
				gl.Vertex3f(float32(box.min[XAXIS]), float32(box.min[YAXIS]-0.04), float32(box.max[ZAXIS]))
			}
			gl.End()
			gl.PopMatrix()
		}

	}

	// gl.FeedbackBuffer(4096, gl.GL_3D_COLOR_TEXTURE, &feedbackBuffer.buffer[0])
	// gl.RenderMode(gl.FEEDBACK)
	// //ThePlayer.Draw(center, true)
	// TheWorld.Draw(center, true)
	// feedbackBuffer.size = gl.RenderMode(gl.RENDER)

	gl.Finish()
	gl.Flush()
	sdl.GL_SwapBuffers()
}
