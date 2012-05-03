/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/gl"
	"math"
)

type Picker struct {
	x, y, radius float32
}

func NewPicker() *Picker {
	var p Picker
	p.radius = float32(90) * PIXEL_SCALE
	return &p
}

func (self *Picker) Draw(t int64) {

	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()

	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 128)
	gl.Vertex2f(self.x, self.y)
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(self.x-float32(math.Sin(angle))*self.radius, self.y+float32(math.Cos(angle))*self.radius)
	}
	gl.End()

	self.DrawItemHighlight(t, ThePlayer.currentAction)
	self.DrawPlayerItems(t)

	gl.PopMatrix()

}

func (self *Picker) DrawItemHighlight(t int64, position uint8) {
	gl.PushMatrix()
	gl.LoadIdentity()

	selectionRadius := blockscale * 1.2
	actionItemRadius := self.radius - blockscale*1.5

	actionItemAngle := -(float64(position) - 1.5) * math.Pi / 4
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 228)
	gl.Vertex2f(self.x-actionItemRadius*float32(math.Sin(actionItemAngle)), self.y+actionItemRadius*float32(math.Cos(actionItemAngle)))
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(self.x-actionItemRadius*float32(math.Sin(actionItemAngle))-float32(math.Sin(angle))*selectionRadius, self.y+actionItemRadius*float32(math.Cos(actionItemAngle))+float32(math.Cos(angle))*selectionRadius)
	}
	gl.End()
	gl.PopMatrix()
}

func (self *Picker) DrawPlayerItems(t int64) {
	actionItemRadius := self.radius - blockscale*1.5

	gl.PushMatrix()
	gl.LoadIdentity()

	for i := 0; i < 5; i++ {

		item := ThePlayer.equippedItems[i]
		if item != ITEM_NONE {
			angle := -(float64(i) + 1.5) * math.Pi / 4
			gl.LoadIdentity()
			x := self.x - actionItemRadius*float32(math.Sin(angle))
			y := self.y + actionItemRadius*float32(math.Cos(angle))
			gl.Translatef(x, y, 0)

			gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 1.0, 0.0, 0.0)
			gl.Rotatef(360*float32(math.Cos(float64(t)/1e10+float64(i))), 0.0, 1.0, 0.0)
			gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 0.0, 0.0, 1.0)
			gl.Scalef(blockscale, blockscale, blockscale)
			gVertexBuffer.Reset()
			TerrainCube(gVertexBuffer, 0, 0, 0, [6]uint16{0, 0, 0, 0, 0, 0}, byte(item), FACE_NONE)
			gVertexBuffer.RenderDirect()

			gl.LoadIdentity()
			gl.Translatef(x-17*PIXEL_SCALE, y-19*PIXEL_SCALE, 20)
			consoleFont.Print(fmt.Sprintf("%d", ThePlayer.inventory[item]))

		}
	}
	gl.PopMatrix()

}
