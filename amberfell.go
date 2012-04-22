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
	"github.com/kierdavis/go/amberfell/mm3dmodel"
	"image/color"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	flag_cpuprofile       = flag.Bool("cpuprofile", false, "write cpu profile to file")
	flag_memprofile       = flag.Bool("memprofile", false, "write memory profile to file")
	DebugMode       bool  = false
	ViewRadius      int16 = 30
	TheWorld        *World
	ThePlayer       *Player
	viewport        Viewport

	timeOfDay float32 = 8

	WolfModel *mm3dmodel.Model

	consoleFont *Font
	metrics     Metrics
)

type Metrics struct {
	fps float64
	mem runtime.MemStats
}

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
	// gl.Enable ( gl.COLOR_MATERIAL )

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
	//LoadTerrainCubes()
	InitTerrainBlocks()

	consoleFont = NewFont("res/Jura-DemiBold.ttf", 16, color.RGBA{255, 255, 255, 0})
	// consoleFont = NewFont("res/FreeMono.ttf", 16, color.RGBA{255, 255, 255, 0})

	viewport.Reshape(int(screen.W), int(screen.H))

}

func quit() {

	ttf.Quit()
	sdl.Quit()
	println("Thanks for playing.")
}
