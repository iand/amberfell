/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/gl"
	"github.com/kierdavis/go/amberfell/mm3dmodel"
	"math"
	"math/rand"
	"os"
)

var WolfModel *mm3dmodel.Model

func init() {
	WolfModel = LoadModel("res/wolf.mm3d")
}

type Wolf struct {
	heading  float64
	position Vectorf
	velocity Vectorf
	falling  bool
	bounce   float64
}

func (w *Wolf) Init(heading float64, x float32, z float32, y float32) {
	w.heading = heading
	w.position[XAXIS] = float64(x)
	w.position[YAXIS] = float64(y)
	w.position[ZAXIS] = float64(z)
}

func (w *Wolf) W() float64 { return 1 }
func (w *Wolf) H() float64 { return 2 }
func (w *Wolf) D() float64 { return 3 }

func (w *Wolf) Heading() float64  { return w.heading }
func (w *Wolf) X() float32        { return float32(w.position[XAXIS]) }
func (w *Wolf) Y() float32        { return float32(w.position[YAXIS]) }
func (w *Wolf) Z() float32        { return float32(w.position[ZAXIS]) }
func (w *Wolf) Velocity() Vectorf { return w.velocity }
func (w *Wolf) Position() Vectorf { return w.position }

func (w *Wolf) FrontBlock() Vectori {
	ip := IntPosition(w.Position())
	if w.heading > 337.5 || w.heading <= 22.5 {
		ip[XAXIS]++
	} else if w.heading > 22.5 && w.heading <= 67.5 {
		ip[XAXIS]++
		ip[ZAXIS]--
	} else if w.heading > 67.5 && w.heading <= 112.5 {
		ip[ZAXIS]--
	} else if w.heading > 112.5 && w.heading <= 157.5 {
		ip[XAXIS]--
		ip[ZAXIS]--
	} else if w.heading > 157.5 && w.heading <= 202.5 {
		ip[XAXIS]--
	} else if w.heading > 202.5 && w.heading <= 247.5 {
		ip[XAXIS]--
		ip[ZAXIS]++
	} else if w.heading > 247.5 && w.heading <= 292.5 {
		ip[ZAXIS]++
	} else if w.heading > 292.5 && w.heading <= 337.5 {
		ip[XAXIS]++
		ip[ZAXIS]++
	}

	return ip
}

func (w *Wolf) SetFalling(b bool) { w.falling = b }

func (w *Wolf) Rotate(angle float64) {
	w.heading += angle
	if w.heading < 0 {
		w.heading += 360
	}
	if w.heading > 360 {
		w.heading -= 360
	}
}

func (w *Wolf) Update(dt float64) {
	w.position[XAXIS] += w.velocity[XAXIS] * dt
	w.position[YAXIS] += w.velocity[YAXIS] * dt
	w.position[ZAXIS] += w.velocity[ZAXIS] * dt
	// fmt.Printf("position: %s\n", w.position)
}

func (w *Wolf) Accelerate(v Vectorf) {
	w.velocity[XAXIS] += v[XAXIS]
	w.velocity[YAXIS] += v[YAXIS]
	w.velocity[ZAXIS] += v[ZAXIS]
}

func (w *Wolf) IsFalling() bool {
	return w.falling
}

func (w *Wolf) Snapx(x float64, vx float64) {
	w.position[XAXIS] = x
	w.velocity[XAXIS] = vx
}

func (w *Wolf) Snapz(z float64, vz float64) {
	w.position[ZAXIS] = z
	w.velocity[ZAXIS] = vz
}

func (w *Wolf) Snapy(y float64, vy float64) {
	w.position[YAXIS] = y
	w.velocity[YAXIS] = vy
}

func (w *Wolf) Setvx(vx float64) {
	w.velocity[XAXIS] = vx
}

func (w *Wolf) Setvz(vz float64) {
	w.velocity[ZAXIS] = vz
}

func (w *Wolf) Setvy(vy float64) {
	w.velocity[YAXIS] = vy
}

func (w *Wolf) Forward(v float64) {
	w.velocity[XAXIS] = math.Cos(w.Heading() * math.Pi / 180)
	w.velocity[ZAXIS] = -math.Sin(w.Heading() * math.Pi / 180)
}

func (w *Wolf) Act(dt float64) {
	w.Rotate(rand.Float64()*9 - 4.5)
	w.Forward(rand.Float64()*4 - 1)
	w.bounce += 360 * dt
	w.velocity[YAXIS] = 0.8 * math.Abs(math.Sin(w.bounce*math.Pi/180))
}

func (wolf *Wolf) Draw(center Vectorf) {
	gl.PushMatrix()
	gl.Translatef(float32(wolf.X()), float32(wolf.Y()), float32(wolf.Z()))
	gl.Rotated(wolf.Heading(), 0.0, 1.0, 0.0)
	//Cuboid(wolf.W(), wolf.H(), wolf.D(), 33, 32, 32, 32, 32, 32, 0, selectMode)
	//	Cuboid(0.3, 0.5, 1.2, &MapTextures[33], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32])
	WolfModel.GLDraw()
	// gl.Translatef(0.8, 0.3, 0)
	// gl.Rotated(-10, 0.0, 0.0, 1.0)
	// Cuboid(0.3, 0.3, 0.4, 33, 32, 32, 32, 32, 32, 0, selectMode)
	gl.PopMatrix()
}
