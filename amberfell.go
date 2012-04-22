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
	"github.com/banthar/gl"
	"github.com/kierdavis/go/amberfell/mm3dmodel"
	"math/rand"
	"os"
	"runtime/pprof"
)

var (
	flag_profile *bool = flag.Bool("profile", false, "Output profiling information to amberfell.prof")
	DebugMode    bool  = false
	ViewRadius   int16 = 30
	TheWorld     *World
	ThePlayer    *Player
	viewport     Viewport

	tileWidth       = 48
	screenScale int = int(5 * float64(tileWidth) / 2)

	timeOfDay float32 = 8

	WolfModel *mm3dmodel.Model
)

func main() {
	flag.Parse()

	if *flag_profile {
		pfile, err := os.Create("amberfell.prof")

		if err != nil {
			panic(fmt.Sprintf("Could not create amberfell.prof:", err))
		}

		pprof.StartCPUProfile(pfile)
	}

	rand.Seed(71)

	defer quit()

	initGame()
	GameLoop()

	return

}

func initGame() {
	WolfModel = LoadModel("res/wolf.mm3d")

	TheWorld = new(World)
	TheWorld.Init()

	ThePlayer = new(Player)
	ThePlayer.Init(0, 10, 10, TheWorld.FindSurface(10, 10))

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

func quit() {
	sdl.Quit()
	println("Thanks for playing.")
}
