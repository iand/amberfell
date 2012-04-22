/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"github.com/banthar/Go-SDL/sdl"
)

var flag_profile *bool = flag.Bool("profile", false, "Output profiling information to amberfell.prof")

var (
	DebugMode  bool  = false
	ViewRadius int16 = 30
	TheWorld   *World
	ThePlayer  *Player
)

var viewport Viewport

var (
	tileWidth                     = 48
	screenScale               int = int(5 * float64(tileWidth) / 2)

	timeOfDay float32 = 8

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

	InitGame()
	InitGraphics()

	GameLoop()

	return

}


func quit() {
	sdl.Quit()
	println("Thanks for playing.")
}