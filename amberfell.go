/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main


import (
    "math/rand"  
    "flag"
    "github.com/iand/amberfell/af"
)    


var printInfo = flag.Bool("info", false, "print GL implementation information")

func main() {
    flag.Parse()
    rand.Seed(71)   


    defer af.QuitGame()
    defer af.QuitGraphics()

    af.InitGame()
    af.InitGraphics()
    af.GameLoop()
    
    return

}













