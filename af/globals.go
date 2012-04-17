/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package af

var (
    DebugMode bool = false
    ViewRadius int16 = 30
    TheWorld *World
    ThePlayer *Player
)


const piover180 = 0.0174532925




var T0 uint32 = 0
var Frames uint32 = 0


var view_rotx float64 = 50.0
var view_roty float64 = 50.0
var view_rotz float64 = 0.0
var gear1, gear2, gear3 uint
var angle float64 = 0.0



var (
    screenWidth, screenHeight int
    tileWidth = 48
    screenScale int = int(5 * float64(tileWidth) / 2)
    ShowOverlay bool

    timeOfDay float32 = 19

    lightpos Vector


)