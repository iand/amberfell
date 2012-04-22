/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
)

type Viewport struct {
	rotx         float64
	roty         float64
	rotz         float64
	x            float64
	y            float64
	z            float64
	scale        float64
	screenWidth  int
	screenHeight int
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
	self.rotx += angle
	if self.rotx > 75 {
		self.rotx = 75
	} else if self.rotx < 15 {
		self.rotx = 15
	}
	// self.Recalc()
}
func (self *Viewport) Roty(angle float64) {
	self.roty += angle
	// self.Recalc()
}

// func (self *Viewport) Transx(d float64) {
// 	self.transx += d
// 	// self.Recalc()
// }

func (self *Viewport) Rotz(angle float64) {
	self.rotz += angle
	// self.Recalc()
}

// func (self *Viewport) Transy(d float64) {
// 	self.transy += d
// 	self.Recalc()
// }

// func (self *Viewport) Transz(d float64) {
// 	self.transz += d
// 	self.Recalc()
// }

func (self *Viewport) Zoomstd() {
	self.scale = 0.75
	// self.Recalc()
}

func (self *Viewport) Zoomin() {
	self.scale += 0.2
	if self.scale > 3 {
		self.scale = 3
	}
	// self.Recalc()
}

func (self *Viewport) Zoomout() {
	self.scale -= 0.2
	if self.scale < 0.2 {
		self.scale = 0.2
	}
	// self.Recalc()
}

func ModelMatrix() *Matrix4 {
	return NewIdentity().Rotatex(viewport.rotx).Rotatey(viewport.roty-ThePlayer.Heading()).Rotatez(viewport.rotz).Translation(-ThePlayer.position[XAXIS], -ThePlayer.position[YAXIS], -ThePlayer.position[ZAXIS])
}
