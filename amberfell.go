/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"flag"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/Go-SDL/ttf"
	"github.com/banthar/gl"
	"github.com/kierdavis/go/mm3dmodel"
	"image/color"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var (
	flag_cpuprofile       = flag.Bool("cpuprofile", false, "write cpu profile to file")
	flag_memprofile       = flag.Bool("memprofile", false, "write memory profile to file")
	DebugMode       bool  = false
	InventoryMode   bool  = false
	ViewRadius      int16 = 30
	TheWorld        *World
	ThePlayer       *Player
	viewport        Viewport

	timeOfDay float32 = 8

	WolfModel *mm3dmodel.Model

	consoleFont    *Font
	textures       map[uint16]*gl.Texture = make(map[uint16]*gl.Texture)
	terrainTexture *gl.Texture
	gVertexBuffer  *VertexBuffer

	// HUD elements
	picker  Picker
	console Console
)

func main() {
	flag.Parse()

	if *flag_cpuprofile {
		pfile, err := os.Create("amberfell.prof")

		if err != nil {
			panic(fmt.Sprintf("Could not create amberfell.prof:", err))
		}

		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(71)

	defer quit()

	initGame()
	GameLoop()

	if *flag_memprofile {
		pfile, err := os.Create("amberfell.prof")

		if err != nil {
			panic(fmt.Sprintf("Could not create amberfell.mprof:", err))
		}

		pprof.WriteHeapProfile(pfile)
		pfile.Close()
	}

	return

}

func initGame() {

	viewport.Zoomstd()
	viewport.Rotx(25)
	viewport.Roty(70)
	// viewport.Transx(-float64(ThePlayer.X()))
	// viewport.Transy(-float64(ThePlayer.Y()))
	// viewport.Transz(-float64(ThePlayer.Z()))

	sdl.Init(sdl.INIT_VIDEO)
	if ttf.Init() != 0 {
		panic("Could not initalize fonts")
	}

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ShadeModel(gl.SMOOTH)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.LIGHT1)

	// gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.EMISSION )
	// gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE )
	gl.Enable(gl.COLOR_MATERIAL)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	// gl.Ortho(-12.0, 12.0, -12.0, 12.0, -10, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	// glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0)

	gl.ClearDepth(1.0)       // Depth Buffer Setup
	gl.Enable(gl.DEPTH_TEST) // Enables Depth Testing
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)

	gl.Enable(gl.TEXTURE_2D)
	LoadMapTextures()
	LoadPlayerTextures()
	//LoadTerrainCubes()
	InitTerrainBlocks()

	consoleFont = NewFont("res/Jura-DemiBold.ttf", 16, color.RGBA{255, 255, 255, 0})
	// consoleFont = NewFont("res/FreeMono.ttf", 16, color.RGBA{255, 255, 255, 0})

	textures[TEXTURE_PICKER] = loadTexture("res/dial.png")
	terrainTexture = loadTexture("tiles.png")

	gVertexBuffer = NewVertexBuffer(10000, terrainTexture)

	WolfModel = LoadModel("res/wolf.mm3d")

	TheWorld = new(World)
	TheWorld.Init()

	ThePlayer = new(Player)
	ThePlayer.Init(0, 32760, 32760, TheWorld.FindSurface(10, 10))

	viewport.Reshape(int(screen.W), int(screen.H))

}

func quit() {

	ttf.Quit()
	sdl.Quit()
	println("Thanks for playing.")
}

func GameLoop() {
	var startTime int64 = time.Now().UnixNano()
	var currentTime, accumulator int64 = 0, 0
	var t, dt int64 = 0, 1e9 / 40
	var drawFrame, computeFrame int64 = 0, 0

	update500ms := new(Timer)
	update500ms.interval = 500 * 1e6
	update500ms.Start()

	update150ms := new(Timer)
	update150ms.interval = 50 * 1e6
	update150ms.Start()

	debugModekeyLock := false
	inventoryModekeyLock := false

	var interactingBlock *InteractingBlockFace

	done := false
	for !done {

		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.ResizeEvent:
				re := e.(*sdl.ResizeEvent)
				screen := sdl.SetVideoMode(int(re.W), int(re.H), 16,
					sdl.OPENGL|sdl.RESIZABLE)
				if screen != nil {
					viewport.Reshape(int(screen.W), int(screen.H))
				} else {
					panic("Could not set video mode")
				}
				break

			case *sdl.MouseButtonEvent:
				re := e.(*sdl.MouseButtonEvent)
				if re.Button == 1 && re.State == 1 { // LEFT, DOWN
					if ThePlayer.CanInteract() {
						selectedBlockFace := viewport.SelectedBlockFace()
						if selectedBlockFace != nil {
							if interactingBlock == nil || interactingBlock.blockFace.pos != selectedBlockFace.pos {
								interactingBlock = new(InteractingBlockFace)
								interactingBlock.blockFace = selectedBlockFace
								interactingBlock.hitCount = 0
							}
							ThePlayer.Interact(interactingBlock)
						}
						// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
					}
				}
			case *sdl.QuitEvent:
				done = true
			case *sdl.KeyboardEvent:
				re := e.(*sdl.KeyboardEvent)
				if re.Keysym.Sym == sdl.K_F3 {
					if !debugModekeyLock && re.Type == sdl.KEYDOWN {
						debugModekeyLock = true
						if DebugMode == true {
							DebugMode = false
						} else {
							DebugMode = true
						}
					} else if re.Type == sdl.KEYUP {
						debugModekeyLock = false
					}
				}
				if re.Keysym.Sym == sdl.K_i {
					if !inventoryModekeyLock && re.Type == sdl.KEYDOWN {
						inventoryModekeyLock = true
						if InventoryMode == true {
							InventoryMode = false
						} else {
							InventoryMode = true
						}
					} else if re.Type == sdl.KEYUP {
						inventoryModekeyLock = false
					}
				}
			}
		}
		keys := sdl.GetKeyState()

		if keys[sdl.K_ESCAPE] != 0 {
			// ShowOverlay = !ShowOverlay

			// Overlay

		}
		if keys[sdl.K_UP] != 0 {
			if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
				viewport.Zoomin()
			} else {
				viewport.Rotx(5)
			}
		}
		if keys[sdl.K_DOWN] != 0 {
			if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
				viewport.Zoomout()
			} else {
				viewport.Rotx(-5)
			}
		}
		if keys[sdl.K_LEFT] != 0 {
			viewport.Roty(9)
		}
		if keys[sdl.K_RIGHT] != 0 {
			viewport.Roty(-9)
		}

		ThePlayer.HandleKeys(keys)

		if keys[sdl.K_o] != 0 {
			timeOfDay += 0.1
			if timeOfDay > 24 {
				timeOfDay -= 24
			}
		}
		if keys[sdl.K_l] != 0 {
			timeOfDay -= 0.1
			if timeOfDay < 0 {
				timeOfDay += 24
			}
		}

		if update150ms.PassedInterval() {
			debugModekeyLock = false

			// If player is breaking a block then allow them to hold mouse down to continue action
			if interactingBlock != nil && ThePlayer.currentAction == ACTION_BREAK {
				mouseState := sdl.GetMouseState(nil, nil)
				if mouseState == 1 {
					if ThePlayer.CanInteract() {
						selectedBlockFace := viewport.SelectedBlockFace()
						if selectedBlockFace != nil {
							if interactingBlock == nil || !interactingBlock.blockFace.pos.Equals(&selectedBlockFace.pos) {
								interactingBlock = new(InteractingBlockFace)
								interactingBlock.blockFace = selectedBlockFace
								interactingBlock.hitCount = 0
							}
							ThePlayer.Interact(interactingBlock)
						}
						// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
					}
				}
			}
			update150ms.Start()
		}

		if update500ms.PassedInterval() {
			console.fps = float64(drawFrame) / (float64(update500ms.GetTicks()) / float64(1e9))
			console.Update()
			timeOfDay += 0.02
			if timeOfDay > 24 {
				timeOfDay -= 24
			}

			drawFrame, computeFrame = 0, 0
			update500ms.Start()
		}

		newTime := time.Now().UnixNano()
		deltaTime := newTime - currentTime
		currentTime = newTime
		if deltaTime > 1e9/4 {
			deltaTime = 1e9 / 4
		}

		accumulator += deltaTime

		for accumulator > dt {
			accumulator -= dt

			TheWorld.ApplyForces(ThePlayer, float64(dt)/1e9)

			ThePlayer.Update(float64(dt) / 1e9)
			TheWorld.Simulate(float64(dt) / 1e9)

			computeFrame++
			t += dt
		}

		//interpolate(previous, current, accumulator/dt)

		Draw(currentTime - startTime)
		drawFrame++
	}

}

func Draw(t int64) {
	gVertexBuffer.Reset()

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

	gl.RenderMode(gl.RENDER)

	selectedBlockFace := viewport.SelectedBlockFace()
	ThePlayer.Draw(center, selectedBlockFace)
	TheWorld.Draw(center, selectedBlockFace)

	if !InventoryMode {
		picker.Draw(t)
	} else {

	}
	if DebugMode {
		console.Draw(t)
	}

	gl.Finish()
	gl.Flush()
	sdl.GL_SwapBuffers()
	// runtime.GC()
}
