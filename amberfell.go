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
	"runtime"
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

	timeOfDay     float32 = 9
	sunlightLevel int     = 7

	worldSeed = int64(20000)
	treeLine  = uint16(math.Trunc(5.0 * float64(CHUNK_HEIGHT/6.0)))
	WolfModel *mm3dmodel.Model

	consoleFont       *Font
	inventoryItemFont *Font
	pauseFont         *Font
	textures          map[uint16]*gl.Texture = make(map[uint16]*gl.Texture)
	terrainTexture    *gl.Texture
	itemsTexture      *gl.Texture
	gVertexBuffer     *VertexBuffer
	gGuiBuffer        *VertexBuffer
	terrainBuffer     *VertexBuffer

	items map[ItemId]Item

	// HUD elements
	blockscale float32 = 0.4 // The scale at which to render blocks in the HUD
	picker     *Picker
	console    Console
	inventory  *Inventory
	pause      Pause
)

func main() {
	flag.Parse()

	println("Setting GOMAXPROCS to ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *flag_cpuprofile {
		pfile, err := os.Create("amberfell.prof")

		if err != nil {
			panic(fmt.Sprintf("Could not create amberfell.prof: %s", err))
		}

		pprof.StartCPUProfile(pfile)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(71)

	defer quit()

	initGame()
	loop()
	//	gameLoop()

	if *flag_memprofile {
		pfile, err := os.Create("amberfell.prof")

		if err != nil {
			panic(fmt.Sprintf("Could not create amberfell.mprof:%s", err))
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

	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	screen := sdl.SetVideoMode(1024, 600, 32, sdl.OPENGL|sdl.RESIZABLE)

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
	LoadWolfTextures()
	InitItems()

	pauseFont = NewFont("res/Jura-DemiBold.ttf", 48, color.RGBA{255, 255, 255, 0})
	consoleFont = NewFont("res/Jura-DemiBold.ttf", 16, color.RGBA{255, 255, 255, 0})
	inventoryItemFont = NewFont("res/Jura-DemiBold.ttf", 14, color.RGBA{240, 240, 240, 0})

	textures[TEXTURE_PICKER] = loadTexture("res/dial.png")
	terrainTexture = loadTexture("res/tiles.png")
	itemsTexture = loadTexture("res/items.png")

	gVertexBuffer = NewVertexBuffer(10000, terrainTexture)
	gGuiBuffer = NewVertexBuffer(1000, terrainTexture)

	WolfModel = LoadModel("res/wolf.mm3d")

	TheWorld = NewWorld()

	ThePlayer = new(Player)
	ThePlayer.Init(0, PLAYER_START_X, PLAYER_START_Z)

	inventory = NewInventory()

	viewport.Reshape(int(screen.W), int(screen.H))
	PreloadChunks(400)

}

func quit() {

	ttf.Quit()
	sdl.Quit()
	println("Thanks for playing.")
}

func loop() {

	ticksim := time.Tick(50 * time.Millisecond)
	tickEnvironment := time.Tick(250 * time.Millisecond)
	tickUI := time.Tick(500 * time.Millisecond)
	tickHousekeeping := time.Tick(2 * time.Second)

	for {
		if HandleUserInput() {
			return
		}

		Draw(time.Now().UnixNano())

		if !pause.visible {
			select {

			case <-ticksim:
				TheWorld.Simulate()

			case <-tickEnvironment:
				MouseRepeat()
				PreloadChunks(100)

			case <-tickUI:
				console.Update()
				UpdateTimeOfDay(false)
				PreloadChunks(220)

			case <-tickHousekeeping:
				CullChunks()
				UpdatePlayerStats()

			default:
				break
			}
		}
	}
}

func HandleUserInput() bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch e.(type) {
		case *sdl.QuitEvent:
			return true
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
			HandleMouseButton(re)

		case *sdl.KeyboardEvent:
			re := e.(*sdl.KeyboardEvent)

			if re.Keysym.Sym == sdl.K_ESCAPE && re.Type == sdl.KEYDOWN {
				if inventory.visible {
					inventory.Hide()
				} else {
					if pause.visible {
						pause.visible = false
					} else {
						pause.visible = true
					}
				}
			}

			if re.Keysym.Sym == sdl.K_F3 && re.Type == sdl.KEYDOWN {
				console.visible = !console.visible
			}

			if !pause.visible {
				if re.Keysym.Sym == sdl.K_i && re.Type == sdl.KEYDOWN {
					if inventory.visible {
						inventory.Hide()
					} else {
						inventory.Show(nil, nil)
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

	return false
}

func HandleMouseButton(re *sdl.MouseButtonEvent) {
	if inventory.visible {
		inventory.HandleMouseButton(re)
	} else {
		picker.HandleMouseButton(re)
		ThePlayer.HandleMouseButton(re)
	}
}

func MouseRepeat() {
	var x, y int
	mouseState := sdl.GetMouseState(&x, &y)

	if mouseState&MOUSE_BTN_LEFT == MOUSE_BTN_LEFT {
		HandleMouseButton(&sdl.MouseButtonEvent{Type: sdl.MOUSEBUTTONDOWN, Button: sdl.BUTTON_LEFT, State: 1, X: uint16(x), Y: uint16(y)})
	}
	if mouseState&MOUSE_BTN_RIGHT == MOUSE_BTN_RIGHT {
		HandleMouseButton(&sdl.MouseButtonEvent{Type: sdl.MOUSEBUTTONDOWN, Button: sdl.BUTTON_RIGHT, State: 1, X: uint16(x), Y: uint16(y)})

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
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.LoadIdentity()

	center := ThePlayer.Position()

	matrix := ModelMatrix().Float32()
	gl.MultMatrixf(&matrix[0])

	// Sun
	gl.Materialfv(gl.FRONT, gl.AMBIENT, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.DIFFUSE, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.SPECULAR, []float32{0.1, 0.1, 0.1, 1})
	gl.Materialfv(gl.FRONT, gl.SHININESS, []float32{0.0, 0.0, 0.0, 1})

	daylightIntensity := float32(SUNLIGHT_LEVELS[sunlightLevel])

	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{daylightIntensity / 2.5, daylightIntensity / 2.5, daylightIntensity / 2.5, 1})

	gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{30 * float32(math.Sin(ThePlayer.Heading()*math.Pi/180)), 60, 30 * float32(math.Cos(ThePlayer.Heading()*math.Pi/180)), 0})
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity * 2, daylightIntensity * 2, daylightIntensity * 2, 1})
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity, 1})

	gl.RenderMode(gl.RENDER)

	var selectedBlockFace *BlockFace
	if !pause.visible && !inventory.visible {
		selectedBlockFace = viewport.SelectedBlockFace()
	}
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
	console.framesDrawn++
}

func UpdateTimeOfDay(reverse bool) {
	if reverse {
		timeOfDay -= 0.02
		if timeOfDay < 0 {
			timeOfDay += 24
		}
	} else {
		timeOfDay += 0.02
		if timeOfDay > 24 {
			timeOfDay -= 24
		}
	}

	if timeOfDay >= 21.5 || timeOfDay < 4.5 {
		sunlightLevel = 0
	} else if timeOfDay >= 4.5 && timeOfDay < 6 {
		sunlightLevel = 1 + int(6*(timeOfDay-4.5)/1.5)
	} else if timeOfDay >= 6 && timeOfDay < 20 {
		sunlightLevel = 7
	} else {
		sunlightLevel = 7 - int(6*(timeOfDay-20)/1.5)
	}

}

func PreloadChunks(maxtime int64) {
	startTicks := time.Now().UnixNano()
	center := ThePlayer.Position()
	px, pz := chunkCoordsFromWorld(uint16(center[XAXIS]), uint16(center[ZAXIS]))

	r := 1
	rmax := int(viewRadius/CHUNK_WIDTH) + 2

	x := -r
	z := -r

	var adjacents [4]*Chunk
	for time.Now().UnixNano()-startTicks < maxtime*1e6 && r < rmax {

		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)+x), uint16(int(pz)+z)-1)]; ok {
			adjacents[NORTH_FACE] = ac
		} else {
			adjacents[NORTH_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)+x), uint16(int(pz)+z)+1)]; ok {
			adjacents[SOUTH_FACE] = ac
		} else {
			adjacents[SOUTH_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)+x)+1, uint16(int(pz)+z))]; ok {
			adjacents[EAST_FACE] = ac
		} else {
			adjacents[EAST_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)+x)-1, uint16(int(pz)+z))]; ok {
			adjacents[WEST_FACE] = ac
		} else {
			adjacents[WEST_FACE] = nil
		}

		go TheWorld.GetChunk(uint16(int(px)+x), uint16(int(pz)+z)).Render(adjacents, nil, nil)

		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)-x), uint16(int(pz)-z)-1)]; ok {
			adjacents[NORTH_FACE] = ac
		} else {
			adjacents[NORTH_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)-x), uint16(int(pz)-z)+1)]; ok {
			adjacents[SOUTH_FACE] = ac
		} else {
			adjacents[SOUTH_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)-x)+1, uint16(int(pz)-z))]; ok {
			adjacents[EAST_FACE] = ac
		} else {
			adjacents[EAST_FACE] = nil
		}
		if ac, ok := TheWorld.chunks[chunkIndex(uint16(int(px)-x)-1, uint16(int(pz)-z))]; ok {
			adjacents[WEST_FACE] = ac
		} else {
			adjacents[WEST_FACE] = nil
		}

		go TheWorld.GetChunk(uint16(int(px)-x), uint16(int(pz)-z)).Render(adjacents, nil, nil)

		if z == r {
			if x == r {
				r++
				x = -r
				z = -r
			} else {
				x++
			}
		} else {
			z++
		}

	}
}
func CullChunks() {
	center := ThePlayer.Position()
	pxmin, pzmin := chunkCoordsFromWorld(uint16(center[XAXIS]-float64(viewRadius)), uint16(center[ZAXIS]-float64(viewRadius)))
	pxmax, pzmax := chunkCoordsFromWorld(uint16(center[XAXIS]+float64(viewRadius)), uint16(center[ZAXIS]+float64(viewRadius)))

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
