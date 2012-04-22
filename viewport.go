/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"math"
)

type Viewport struct {
	rotx              float64
	roty              float64
	rotz              float64
	scale             float64
	screenWidth       int
	screenHeight      int
	mousex            int
	mousey            int
	selectionDirty    bool
	selectedBlockFace *BlockFace
}

func (self *Viewport) ScreenToView(xs uint16, ys uint16) (xv float64, yv float64) {
	// xs = 0 => -float64(screenWidth) / screenScale
	// xs = screenWidth => float64(screenWidth) / screenScale

	viewWidth := 2 * float64(self.screenWidth) / float64(screenScale)
	xv = (-viewWidth/2 + viewWidth*float64(xs)/float64(self.screenWidth))

	viewHeight := 2 * float64(self.screenHeight) / float64(screenScale)
	yv = (-viewHeight/2 + viewHeight*float64(ys)/float64(self.screenHeight))

	return
}

/* new window size or exposure */
func (self *Viewport) Reshape(width int, height int) {
	self.selectionDirty = false
	self.screenWidth = width
	self.screenHeight = height

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	xmin, ymin := self.ScreenToView(0, 0)
	xmax, ymax := self.ScreenToView(uint16(width), uint16(height))

	gl.Ortho(float64(xmin), float64(xmax), float64(ymin), float64(ymax), -20, 20)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	// glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
}

func (self *Viewport) Rotx(angle float64) {
	self.selectionDirty = false
	self.rotx += angle
	if self.rotx > 75 {
		self.rotx = 75
	} else if self.rotx < 15 {
		self.rotx = 15
	}
	// self.Recalc()
}
func (self *Viewport) Roty(angle float64) {
	self.selectionDirty = false
	self.roty += angle
	// self.Recalc()
}

func (self *Viewport) Rotz(angle float64) {
	self.selectionDirty = false
	self.rotz += angle
	// self.Recalc()
}

func (self *Viewport) Zoomstd() {
	self.selectionDirty = false
	self.scale = 0.75
	// self.Recalc()
}

func (self *Viewport) Zoomin() {
	self.selectionDirty = false
	self.scale += 0.2
	if self.scale > 3 {
		self.scale = 3
	}
	// self.Recalc()
}

func (self *Viewport) Zoomout() {
	self.selectionDirty = false
	self.scale -= 0.2
	if self.scale < 0.2 {
		self.scale = 0.2
	}
	// self.Recalc()
}

func ModelMatrix() *Matrix4 {
	return NewIdentity().Rotatex(viewport.rotx).Rotatey(viewport.roty-ThePlayer.Heading()).Rotatez(viewport.rotz).Translation(-ThePlayer.position[XAXIS], -ThePlayer.position[YAXIS], -ThePlayer.position[ZAXIS])
}

func (self *Viewport) SelectedBlockFace() *BlockFace {
	var pm32 []float32 = make([]float32, 16)
	var newmousex, newmousey int
	_ = sdl.GetMouseState(&newmousex, &newmousey)

	// if self.selectionDirty || newmousex != self.mousex || newmousey != self.mousey {

	self.selectedBlockFace = nil
	self.mousex = newmousex
	self.mousey = newmousey

	gl.GetFloatv(gl.PROJECTION_MATRIX, pm32)
	var projectionMatrix64 *Matrix4 = NewMatrix(float64(pm32[0]), float64(pm32[1]), float64(pm32[2]), float64(pm32[3]), float64(pm32[4]), float64(pm32[5]), float64(pm32[6]), float64(pm32[7]), float64(pm32[8]), float64(pm32[9]), float64(pm32[10]), float64(pm32[11]), float64(pm32[12]), float64(pm32[13]), float64(pm32[14]), float64(pm32[15]))

	inverseMatrix, _ := projectionMatrix64.Multiply(ModelMatrix()).Inverse()

	x := (float64(self.mousex) - float64(self.screenWidth)/2) / (float64(self.screenWidth) / 2)
	z := (float64(self.screenHeight)/2 - float64(self.mousey)) / (float64(self.screenHeight) / 2)

	origin := inverseMatrix.Transform(&Vectorf{x, z, -1}, 1)
	norm := inverseMatrix.Transform(&Vectorf{0, 0, 1}, 0).Normalize()

	if origin != nil && norm != nil {
		pos := IntPosition(ThePlayer.position)
		ray := Ray{origin, norm}

		// See http://www.dyn-lab.com/articles/pick-selection.html
		var box *Box = nil
		distance := float64(1e9)
		face := uint8(0)
		for dy := int16(5); dy > -6; dy-- {
			for dz := int16(-5); dz < 6; dz++ {
				for dx := int16(-5); dx < 6; dx++ {
					trialDistance := math.Sqrt(math.Pow(float64(pos[XAXIS]+dx)-origin[0], 2) + math.Pow(float64(pos[YAXIS]+dy)-origin[1], 2) + math.Pow(float64(pos[ZAXIS]+dz)-origin[2], 2))
					if trialDistance < distance {
						if TheWorld.At(pos[XAXIS]+dx, pos[YAXIS]+dy, pos[ZAXIS]+dz) != BLOCK_AIR {
							trialBox := &Box{
								&Vectorf{float64(pos[XAXIS]+dx) - 0.5, float64(pos[YAXIS]+dy) - 0.5, float64(pos[ZAXIS]+dz) - 0.5},
								&Vectorf{float64(pos[XAXIS]+dx) + 0.5, float64(pos[YAXIS]+dy) + 0.5, float64(pos[ZAXIS]+dz) + 0.5}}

							hit, trialFace := ray.HitsBox(trialBox)
							if hit /*&& TheWorld.AirNeighbour(pos[XAXIS]+dx, pos[YAXIS]+dy, pos[ZAXIS]+dz, face)*/ {
								distance = trialDistance
								box = trialBox
								face = trialFace
							}

						}
					}
				}
			}
		}

		if box != nil {
			self.selectedBlockFace = &BlockFace{pos: Vectori{int16(box.min[XAXIS] + 0.5), int16(box.min[YAXIS] + 0.5), int16(box.min[ZAXIS] + 0.5)}, face: face}
			self.selectionDirty = false
		}

	}
	// }

	return self.selectedBlockFace

}
