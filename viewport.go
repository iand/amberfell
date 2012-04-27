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
	lplane, rplane    float64
	bplane, tplane    float64
	near, far         float64
}

/* new window size or exposure */
func (self *Viewport) Reshape(width int, height int) {
	self.selectionDirty = false
	self.screenWidth = width
	self.screenHeight = height

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	viewWidth := float64(self.screenWidth) / float64(SCREEN_SCALE)
	viewHeight := float64(self.screenHeight) / float64(SCREEN_SCALE)

	self.lplane = -viewWidth / 2
	self.rplane = viewWidth / 2
	self.bplane = -viewHeight / 2
	self.tplane = viewHeight / 2

	gl.Ortho(self.lplane, self.rplane, self.bplane, self.tplane, -20, 20)
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

	if origin != nil {
		pos := IntPosition(ThePlayer.position)
		ray := Ray{origin, &norm}
		reach := int16(4)

		// See http://www.dyn-lab.com/articles/pick-selection.html
		var box *Box = nil
		distance := float64(1e9)
		face := uint8(0)
		for dy := reach; dy > -(reach+1); dy-- {
			for dz := -reach; dz < reach+1; dz++ {
				for dx := -reach; dx < reach+1; dx++ {
					if dy*dy + dz*dz+dx*dx <= reach*reach {
						blockDirection := Vectorf{float64(dx), float64(dy), float64(dz)}


						if ThePlayer.Facing(blockDirection) && blockDirection.Magnitude() <= float64(reach) {


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
