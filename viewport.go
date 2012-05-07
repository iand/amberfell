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
	// "fmt"
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
	self.bplane = -viewHeight / 4
	self.tplane = 3 * viewHeight / 4

	// println("self.lplane:", self.lplane, "self.rplane", self.rplane)
	gl.Ortho(self.lplane, self.rplane, self.bplane, self.tplane, -40, 40)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	picker.x = float32(viewport.rplane) - picker.radius + blockscale*0.5
	picker.y = float32(viewport.bplane) + picker.radius - blockscale*0.5

}

func (self *Viewport) Rotx(angle float64) {
	self.selectionDirty = false
	self.rotx += angle
	if self.rotx > 75 {
		self.rotx = 75
	} else if self.rotx < 5 {
		self.rotx = 5
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
	self.scale = 1.0
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
	return NewIdentity().Scale(viewport.scale).Rotatex(viewport.rotx).Rotatey(viewport.roty-ThePlayer.Heading()).Rotatez(viewport.rotz).Translation(-ThePlayer.position[XAXIS], -ThePlayer.position[YAXIS], -ThePlayer.position[ZAXIS])
}

func (self *Viewport) ProjectPoint(p *Vectorf) (point *Vectorf, normal *Vectorf) {
	var pm32 []float32 = make([]float32, 16)
	gl.GetFloatv(gl.PROJECTION_MATRIX, pm32)
	var projectionMatrix64 *Matrix4 = NewMatrix(float64(pm32[0]), float64(pm32[1]), float64(pm32[2]), float64(pm32[3]), float64(pm32[4]), float64(pm32[5]), float64(pm32[6]), float64(pm32[7]), float64(pm32[8]), float64(pm32[9]), float64(pm32[10]), float64(pm32[11]), float64(pm32[12]), float64(pm32[13]), float64(pm32[14]), float64(pm32[15]))
	inverseMatrix, _ := projectionMatrix64.Multiply(ModelMatrix()).Inverse()
	point = inverseMatrix.Transform(p, 1)
	normalv := inverseMatrix.Transform(&Vectorf{0, 0, 1}, 0).Normalize()
	return point, &normalv
}

func (self *Viewport) ClipPlanes() *[6][4]float32 {

	var pm32 []float64 = make([]float64, 16)
	gl.GetDoublev(gl.PROJECTION_MATRIX, pm32)
	var projectionMatrix64 *Matrix4 = NewMatrix(pm32[0], pm32[1], pm32[2], pm32[3], pm32[4], pm32[5], pm32[6], pm32[7], pm32[8], pm32[9], pm32[10], pm32[11], pm32[12], pm32[13], pm32[14], pm32[15])
	mvpmatrix := projectionMatrix64.Multiply(ModelMatrix())

	var planes32 [6][4]float32
	planes64 := [6][4]float64{
		[4]float64{mvpmatrix[3] + mvpmatrix[0], mvpmatrix[7] + mvpmatrix[4], mvpmatrix[11] + mvpmatrix[8], mvpmatrix[15] + mvpmatrix[12]},
		[4]float64{mvpmatrix[3] - mvpmatrix[0], mvpmatrix[7] - mvpmatrix[4], mvpmatrix[11] - mvpmatrix[8], mvpmatrix[15] - mvpmatrix[12]},
		[4]float64{mvpmatrix[3] + mvpmatrix[1], mvpmatrix[7] + mvpmatrix[5], mvpmatrix[11] + mvpmatrix[9], mvpmatrix[15] + mvpmatrix[13]},
		[4]float64{mvpmatrix[3] - mvpmatrix[1], mvpmatrix[7] - mvpmatrix[5], mvpmatrix[11] - mvpmatrix[9], mvpmatrix[15] - mvpmatrix[13]},
		[4]float64{mvpmatrix[3] + mvpmatrix[2], mvpmatrix[7] + mvpmatrix[6], mvpmatrix[11] + mvpmatrix[10], mvpmatrix[15] + mvpmatrix[14]},
		[4]float64{mvpmatrix[3] - mvpmatrix[2], mvpmatrix[7] - mvpmatrix[6], mvpmatrix[11] - mvpmatrix[10], mvpmatrix[15] - mvpmatrix[14]},
	}

	for p := 0; p < 6; p++ {
		length := math.Sqrt(math.Pow(planes64[p][0], 2) + math.Pow(planes64[p][1], 2) + math.Pow(planes64[p][2], 2) + math.Pow(planes64[p][3], 2))
		planes32[p] = [4]float32{float32(planes64[p][0] / length), float32(planes64[p][1] / length), float32(planes64[p][2] / length), float32(planes64[p][3] / length)}
	}

	return &planes32
}

func (self *Viewport) SelectedBlockFace() *BlockFace {
	var newmousex, newmousey int
	_ = sdl.GetMouseState(&newmousex, &newmousey)

	self.selectedBlockFace = nil
	self.mousex = newmousex
	self.mousey = newmousey

	x := (float64(self.mousex) - float64(self.screenWidth)/2) / (float64(self.screenWidth) / 2)
	z := (float64(self.screenHeight)/2 - float64(self.mousey)) / (float64(self.screenHeight) / 2)

	origin, norm := self.ProjectPoint(&Vectorf{x, z, -1})

	if origin != nil {
		pos := IntPosition(ThePlayer.position)
		ray := Ray{origin, norm}
		reach := int16(5)

		// See http://www.dyn-lab.com/articles/pick-selection.html
		var box *Box = nil
		distance := float64(1e9)
		face := uint8(0)
		for dy := reach; dy > -(reach + 1); dy-- {
			for dz := -reach; dz < reach+1; dz++ {
				for dx := -reach; dx < reach+1; dx++ {
					if dy*dy+dz*dz+dx*dx <= reach*reach {
						blockDirection := Vectorf{float64(dx), float64(dy), float64(dz)}

						if /* ThePlayer.Facing(blockDirection) && */ blockDirection.Magnitude() <= float64(reach) {

							posTest := pos.Translate(dx, dy, dz)
							trialDistance := math.Sqrt(math.Pow(float64(posTest[XAXIS])-origin[0], 2) + math.Pow(float64(posTest[YAXIS])-origin[1], 2) + math.Pow(float64(posTest[ZAXIS])-origin[2], 2))
							if trialDistance < distance {

								if TheWorld.Atv(posTest) != BLOCK_AIR {
									trialBox := &Box{
										&Vectorf{float64(posTest[XAXIS]) - 0.5, float64(posTest[YAXIS]) - 0.5, float64(posTest[ZAXIS]) - 0.5},
										&Vectorf{float64(posTest[XAXIS]) + 0.5, float64(posTest[YAXIS]) + 0.5, float64(posTest[ZAXIS]) + 0.5}}

									hit, trialFace := ray.HitsBox(trialBox)
									if hit {
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
			self.selectedBlockFace = &BlockFace{pos: Vectori{uint16(box.min[XAXIS] + 0.5), uint16(box.min[YAXIS] + 0.5), uint16(box.min[ZAXIS] + 0.5)}, face: face}
			self.selectionDirty = false
		}

	}
	// }

	return self.selectedBlockFace

}

func (self *Viewport) HandleKeys(keys []uint8) {
	if keys[sdl.K_UP] != 0 {
		if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
			self.Zoomin()
		} else {
			self.Rotx(5)
		}
	}
	if keys[sdl.K_DOWN] != 0 {
		if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
			self.Zoomout()
		} else {
			self.Rotx(-5)
		}
	}
	if keys[sdl.K_LEFT] != 0 {
		if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
			viewRadius -= 4
			if viewRadius < 8 {
				viewRadius = 8
			}
		} else {
			self.Roty(9)
		}
	}
	if keys[sdl.K_RIGHT] != 0 {
		if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
			viewRadius += 4

		} else {
			self.Roty(-9)
		}
	}
}

func (self *Viewport) ScreenCoordsToWorld2D(sx, sy uint16) (x, y float64) {
	x = (float64(self.lplane) + float64(sx)*PIXEL_SCALE)
	y = (float64(self.tplane) - float64(sy)*PIXEL_SCALE)

	return
}
