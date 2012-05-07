/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	// "fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	// "math"
)

type Pause struct {
	visible bool
}

func (self *Pause) Draw(t int64) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	gl.Color4ub(0, 0, 0, 240)

	gl.Begin(gl.QUADS)
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane))
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane))
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.tplane))
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.tplane))
	gl.End()

	str := "paused"
	h, w := pauseFont.Measure(str)

	// x := (viewport.rplane - viewport.lplane - w) / 2 
	// y := (viewport.tplane - viewport.bplane - h) / 2 
	gl.Translated(-w/2, -h/2, 0)
	pauseFont.Print(str)

	gl.PopMatrix()

}

func (self *Pause) HandleMouseButton(re *sdl.MouseButtonEvent) {
	if re.Button == 1 && re.State == 1 { // LEFT, DOWN
		// x, y := viewport.ScreenCoordsToWorld2D(re.X, re.Y)
	}
}

func (self *Pause) HandleMouse(mousex int, mousey int, mousestate uint8) {

	// x, y := viewport.ScreenCoordsToWorld2D(uint16(mousex), uint16(mousey))

}

func (self *Pause) HandleKeyboard(re *sdl.KeyboardEvent) {

}

func (self *Pause) HandleKeys(keys []uint8) {

}
