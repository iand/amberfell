/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
	"math/rand"
)

type Wolf struct {
	MobData
}

func (self *Wolf) Init(heading float64, x float32, z float32, y float32) {
	self.heading = heading
	self.position[XAXIS] = float64(x)
	self.position[YAXIS] = float64(y)
	self.position[ZAXIS] = float64(z)
}

func (self *Wolf) W() float64 { return 2 }
func (self *Wolf) H() float64 { return 2 }
func (self *Wolf) D() float64 { return 1 }

func (self *Wolf) Act(dt float64) {
	self.Rotate(rand.Float64()*9 - 4.5)
	self.Forward(rand.Float64()*4 - 1)
}

func (self *Wolf) Draw(center Vectorf) {
	gl.PushMatrix()
	gl.Translatef(float32(self.position[XAXIS]), float32(self.position[YAXIS]), float32(self.position[ZAXIS]))
	gl.Rotated(self.Heading(), 0.0, 1.0, 0.0)
	WolfModel.GLDraw()
	gl.PopMatrix()
}
