/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

var (
	DebugMode  bool  = false
	ViewRadius int16 = 30
	TheWorld   *World
	ThePlayer  *Player
)

const piover180 = 0.0174532925

var viewport Viewport



var (
	screenWidth, screenHeight int
	tileWidth                     = 48
	screenScale               int = int(5 * float64(tileWidth) / 2)
	ShowOverlay               bool

	timeOfDay float32 = 8

	lightpos Vectorf

)
