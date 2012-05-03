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
}

func (self *Picker) Draw(t int64) {

	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()

	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	blockscale := float32(0.4)
	radius := float32(90) * PIXEL_SCALE
	centrex := float32(viewport.rplane) - radius + blockscale*0.5 //+ 20 * PIXEL_SCALE
	centrey := float32(viewport.bplane) + radius - blockscale*0.5 //- 20 * PIXEL_SCALE

	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 128)
	gl.Vertex2f(centrex, centrey)
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(centrex-float32(math.Sin(angle))*radius, centrey+float32(math.Cos(angle))*radius)
	}
	gl.End()

	selectionRadius := blockscale * 1.2
	actionItemRadius := radius - blockscale*1.5

	actionItemAngle := -(float64(ThePlayer.currentAction) - 1.5) * math.Pi / 4
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 228)
	gl.Vertex2f(centrex-actionItemRadius*float32(math.Sin(actionItemAngle)), centrey+actionItemRadius*float32(math.Cos(actionItemAngle)))
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(centrex-actionItemRadius*float32(math.Sin(actionItemAngle))-float32(math.Sin(angle))*selectionRadius, centrey+actionItemRadius*float32(math.Cos(actionItemAngle))+float32(math.Cos(angle))*selectionRadius)
	}
	gl.End()

	for i := 0; i < 5; i++ {

		item := ThePlayer.equippedItems[i]
		if item != ITEM_NONE {
			angle := -(float64(i) + 1.5) * math.Pi / 4
			gl.LoadIdentity()
			x := centrex - actionItemRadius*float32(math.Sin(angle))
			y := centrey + actionItemRadius*float32(math.Cos(angle))
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
