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
)

var flag_profile *bool = flag.Bool("profile", false, "Output profiling information to amberfell.prof")

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

	defer QuitGame()
	defer QuitGraphics()

	InitGame()
	InitGraphics()

	GameLoop()

	return

}
