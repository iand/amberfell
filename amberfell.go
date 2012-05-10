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
	flag_cpuprofile = flag.Bool("cpuprofile", false, "write cpu profile to file")
	flag_memprofile = flag.Bool("memprofile", false, "write memory profile to file")

	viewRadius int16 = 30
	TheWorld   *World
	ThePlayer  *Player
	viewport   Viewport

	timeOfDay float32 = 9

	worldSeed = int64(16)
	treeLine  = uint16(math.Trunc(5.0 * float64(CHUNK_HEIGHT/6.0)))
	WolfModel *mm3dmodel.Model

	consoleFont       *Font
	inventoryItemFont *Font
	pauseFont         *Font
	textures          map[uint16]*gl.Texture = make(map[uint16]*gl.Texture)
	terrainTexture    *gl.Texture
	itemsTexture      *gl.Texture
	gVertexBuffer     *VertexBuffer
	terrainBuffer     *VertexBuffer

	items map[uint16]Item

	// World elements
	lightSources []*LightSource = make([]*LightSource, 20)
	campfires    []*CampFire    = make([]*CampFire, 20)

	// HUD elements
	blockscale float32 = 0.4 // The scale at which to render blocks in the HUD
	picker     *Picker
	console    Console
	inventory  Inventory
	pause      Pause
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
	gameLoop()

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

	picker = NewPicker()

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

	gl.Enable(gl.COLOR_MATERIAL)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.ClearDepth(1.0)       // Depth Buffer Setup
	gl.Enable(gl.DEPTH_TEST) // Enables Depth Testing
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)

	gl.Enable(gl.TEXTURE_2D)
	//	LoadMapTextures()
	LoadPlayerTextures()
	InitItems()

	pauseFont = NewFont("res/Jura-DemiBold.ttf", 48, color.RGBA{255, 255, 255, 0})
	consoleFont = NewFont("res/Jura-DemiBold.ttf", 16, color.RGBA{255, 255, 255, 0})
	inventoryItemFont = NewFont("res/Jura-DemiBold.ttf", 14, color.RGBA{240, 240, 240, 0})

	textures[TEXTURE_PICKER] = loadTexture("res/dial.png")
	terrainTexture = loadTexture("tiles.png")
	itemsTexture = loadTexture("res/items.png")

	gVertexBuffer = NewVertexBuffer(10000, terrainTexture)

	//WolfModel = LoadModel("res/wolf.mm3d")

	TheWorld = new(World)
	TheWorld.Init()

	ThePlayer = new(Player)
	ThePlayer.Init(0, PLAYER_START_X, PLAYER_START_Z)

	viewport.Reshape(int(screen.W), int(screen.H))

}

func quit() {

	ttf.Quit()
	sdl.Quit()
	println("Thanks for playing.")
}

func gameLoop() {
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

	update2000ms := new(Timer)
	update2000ms.interval = 2000 * 1e6
	update2000ms.Start()

	var interactingBlock *InteractingBlockFace

	done := false
	for !done {

		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				done = true
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
				if inventory.visible {
					inventory.HandleMouseButton(re)
				} else {
					picker.HandleMouseButton(re)

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
				}

			case *sdl.KeyboardEvent:
				re := e.(*sdl.KeyboardEvent)

				if re.Keysym.Sym == sdl.K_ESCAPE {
					if re.Type == sdl.KEYDOWN {
						inventory.visible = false
						if pause.visible {
							pause.visible = false
							update2000ms.Unpause()
							update500ms.Unpause()
							update150ms.Unpause()
						} else {
							pause.visible = true
							update2000ms.Pause()
							update500ms.Pause()
							update150ms.Pause()
						}
					}
				}

				if re.Keysym.Sym == sdl.K_F3 {
					if re.Type == sdl.KEYDOWN {
						if console.visible == true {
							console.visible = false
						} else {
							console.visible = true
						}
					}
				}

				if !pause.visible {
					if re.Keysym.Sym == sdl.K_i {
						if re.Type == sdl.KEYDOWN {
							if inventory.visible == true {
								inventory.visible = false
							} else {
								inventory.visible = true
							}
						}
					}

					if inventory.visible {
						inventory.HandleKeyboard(re)
					}
				}
			}
		}

		keys := sdl.GetKeyState()

		if console.visible {
			console.HandleKeys(keys)
		}

		if pause.visible {
			pause.HandleKeys(keys)
		} else if inventory.visible {
			inventory.HandleKeys(keys)
		} else {
			viewport.HandleKeys(keys)
			ThePlayer.HandleKeys(keys)
			picker.HandleKeys(keys)
		}

		if update150ms.PassedInterval() {

			if !inventory.visible {
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
			}

			update150ms.Start()
		}

		if update500ms.PassedInterval() {
			console.fps = float64(drawFrame) / (float64(update500ms.GetTicks()) / float64(1e9))
			console.Update()

			UpdateTimeOfDay()
			UpdateCampfires()
			PreloadChunks(200)

			drawFrame, computeFrame = 0, 0
			update500ms.Start()

		}

		if update2000ms.PassedInterval() {
			CullChunks()
			UpdatePlayerStats()
			update2000ms.Start()
		}

		if !pause.visible {
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
		}
		//interpolate(previous, current, accumulator/dt)

		Draw(currentTime - startTime)
		drawFrame++
	}

}

func Draw(t int64) {
	console.culledVertices = 0
	console.vertices = 0

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

	if pause.visible {
		pause.Draw(t)
	} else if inventory.visible {
		inventory.Draw(t)
	} else {
		picker.Draw(t)

	}

	if console.visible {
		console.Draw(t)
	}

	gl.Finish()
	gl.Flush()
	sdl.GL_SwapBuffers()
	// runtime.GC()
}

func UpdateTimeOfDay() {
	timeOfDay += 0.02
	if timeOfDay > 24 {
		timeOfDay -= 24
	}
}

func UpdateCampfires() {
	// Age any campfires
	for i := 0; i < len(campfires); i++ {
		if campfires[i] != nil {
			campfires[i].life -= 0.02
			if campfires[i].life <= 0 {
				lightSources[campfires[i].lightSourceIndex] = nil
				campfires[i] = nil
			}
		}
	}
}

func PreloadChunks(maxtime int64) {

	if !ThePlayer.IsMoving() {
		// Load some chunks around where the player is headed
		norm := ThePlayer.Normal()
		center := ThePlayer.Position()
		px, _, pz := chunkCoordsFromWorld(uint16(center[XAXIS]), uint16(center[YAXIS]), uint16(center[ZAXIS]))

		d := 0
		r := uint16(viewRadius/CHUNK_WIDTH) + 1
		rmax := r + 4
		startTicks := time.Now().UnixNano()
		for time.Now().UnixNano()-startTicks < maxtime*1e6 && r < rmax {
			switch d {
			case 0:
				TheWorld.GetChunk(uint16(float64(px)+float64(r)*norm[XAXIS]), 0, pz).PreRender(nil)
			case 1:
				TheWorld.GetChunk(px, 0, uint16(float64(pz)+float64(r)*norm[ZAXIS])).PreRender(nil)
			case 2:
				TheWorld.GetChunk(uint16(float64(px)+float64(r)*norm[XAXIS]), 0, uint16(float64(pz)+float64(r)*norm[ZAXIS])).PreRender(nil)

			}
			d++
			if d > 2 {
				d = 0
				r++
			}
		}
	}
}
func CullChunks() {
	center := ThePlayer.Position()
	pxmin, _, pzmin := chunkCoordsFromWorld(uint16(center[XAXIS]-float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]-float64(viewRadius)))
	pxmax, _, pzmax := chunkCoordsFromWorld(uint16(center[XAXIS]+float64(viewRadius)), uint16(center[YAXIS]), uint16(center[ZAXIS]+float64(viewRadius)))

	// Cull chunks more than 10 chunks away from view radius
	for chunkIndex, chunk := range TheWorld.chunks {
		if chunk.x > pxmax+6 || chunk.x < pxmin-6 || chunk.z > pzmax+6 || chunk.z < pzmin-6 {
			delete(TheWorld.chunks, chunkIndex)
		}
	}

}

func UpdatePlayerStats() {
	center := ThePlayer.Position()
	// Update player stats
	distanceFromStart := uint16(math.Sqrt(math.Pow(center[XAXIS]-float64(PLAYER_START_X), 2) + math.Pow(center[ZAXIS]-float64(PLAYER_START_Z), 2)))
	if distanceFromStart > ThePlayer.distanceFromStart {
		ThePlayer.distanceFromStart = distanceFromStart
	}

}
