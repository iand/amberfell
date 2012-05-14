/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"math"
)

type Picker struct {
	x, y, radius, actionItemRadius, selectionRadius float32
}

func NewPicker() *Picker {
	var p Picker
	p.radius = float32(90) * PIXEL_SCALE
	p.actionItemRadius = p.radius - blockscale*1.5
	p.selectionRadius = blockscale * 1.2
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
	self.DrawPlayerItems(t, true)

	gl.PopMatrix()

}

func (self *Picker) DrawItemHighlight(t int64, position uint8) {
	gl.PushMatrix()
	gl.LoadIdentity()

	actionItemAngle := -(float64(position) - 1.5) * math.Pi / 4
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4ub(0, 0, 0, 228)
	gl.Vertex2f(self.x-self.actionItemRadius*float32(math.Sin(actionItemAngle)), self.y+self.actionItemRadius*float32(math.Cos(actionItemAngle)))
	for angle := float64(0); angle <= 2*math.Pi; angle += math.Pi / 2 / 10 {
		gl.Vertex2f(self.x-self.actionItemRadius*float32(math.Sin(actionItemAngle))-float32(math.Sin(angle))*self.selectionRadius, self.y+self.actionItemRadius*float32(math.Cos(actionItemAngle))+float32(math.Cos(angle))*self.selectionRadius)
	}
	gl.End()
	gl.PopMatrix()
}

func (self *Picker) DrawPlayerItems(t int64, drawQuantities bool) {

	gl.PushMatrix()
	gl.LoadIdentity()

	for i := 0; i < 5; i++ {

		itemid := ThePlayer.equippedItems[i]
		if itemid != ITEM_NONE {
			angle := -(float64(i) + 1.5) * math.Pi / 4
			gl.LoadIdentity()
			x := self.x - self.actionItemRadius*float32(math.Sin(angle))
			y := self.y + self.actionItemRadius*float32(math.Cos(angle))
			gl.Translatef(x, y, 0)

			gl.Rotated(90, 1.0, 0.0, 0.0)
			gl.Rotated(30*math.Sin(float64(t)/1e9+float64(itemid)/2), 0.0, 1.0, 0.0)

			gl.Scalef(blockscale, blockscale, blockscale)
			gGuiBuffer.Reset()
			if itemid < 256 {
				TerrainCube(gGuiBuffer, Vectori{}, [18]uint16{BLOCK_DIRT, BLOCK_DIRT, BLOCK_DIRT, BLOCK_DIRT, BLOCK_AIR, BLOCK_DIRT, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR}, itemid, FACE_NONE)
			} else {
				RenderItemFlat(gGuiBuffer, Vectori{}, itemid)
			}
			gGuiBuffer.RenderDirect(false)

			if drawQuantities {
				gl.LoadIdentity()
				gl.Translatef(x-17*PIXEL_SCALE, y-19*PIXEL_SCALE, 20)
				consoleFont.Print(fmt.Sprintf("%d", ThePlayer.inventory[itemid]))
			}

		}
	}
	gl.PopMatrix()

}

// x and y are in screen2d coords
func (self *Picker) HitTest(x, y float64) (bool, int) {

	for i := 0; i < 8; i++ {
		angle := -(float64(i) - 1.5) * math.Pi / 4
		ix := float64(picker.x) - float64(picker.actionItemRadius)*math.Sin(angle)
		iy := float64(picker.y) + float64(picker.actionItemRadius)*math.Cos(angle)
		if x > ix-float64(picker.selectionRadius) && x < ix+float64(picker.selectionRadius) &&
			y > iy-float64(picker.selectionRadius) && y < iy+float64(picker.selectionRadius) {

			return true, i
		}

	}

	return false, 0
}

func (self *Picker) HandleMouseButton(re *sdl.MouseButtonEvent) {
	if re.Button == 1 && re.State == 1 { // LEFT, DOWN
		x, y := viewport.ScreenCoordsToWorld2D(re.X, re.Y)
		hit, pos := self.HitTest(x, y)
		if hit {
			ThePlayer.SelectAction(pos)
		}
	}
}

func (self *Picker) HandleKeys(keys []uint8) {
	if keys[sdl.K_1] != 0 {
		ThePlayer.SelectAction(0)
	}
	if keys[sdl.K_2] != 0 {
		ThePlayer.SelectAction(1)
	}
	if keys[sdl.K_3] != 0 {
		ThePlayer.SelectAction(2)
	}
	if keys[sdl.K_4] != 0 {
		ThePlayer.SelectAction(3)
	}
	if keys[sdl.K_5] != 0 {
		ThePlayer.SelectAction(4)
	}
	if keys[sdl.K_6] != 0 {
		ThePlayer.SelectAction(5)
	}
	if keys[sdl.K_7] != 0 {
		ThePlayer.SelectAction(6)
	}
	if keys[sdl.K_8] != 0 {
		ThePlayer.SelectAction(7)
	}
}
