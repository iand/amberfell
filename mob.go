/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"math"
)

type MobData struct {
	heading      float64
	position     Vectorf
	velocity     Vectorf
	falling      bool
	walkingSpeed float64
	walkSequence float64
}

func (self *MobData) Heading() float64 { return self.heading }

func (self *MobData) Normal() *Vectorf {
	return &Vectorf{math.Cos(self.heading * math.Pi / 180), 0, -math.Sin(self.heading * math.Pi / 180)}
}

func (self *MobData) Velocity() Vectorf { return self.velocity }
func (self *MobData) Position() Vectorf { return self.position }

func (self *MobData) SetFalling(b bool) { self.falling = b }
func (self *MobData) IsFalling() bool   { return self.falling }
func (self *MobData) IsMoving() bool {
	return self.velocity[XAXIS] != 0 || self.velocity[YAXIS] != 0 || self.velocity[ZAXIS] != 0
}

func (self *MobData) Rotate(angle float64) {
	self.heading += angle
	if self.heading < 0 {
		self.heading += 360
	}
	if self.heading > 360 {
		self.heading -= 360
	}
}

func (self *MobData) Forward(v float64) {
	self.velocity[XAXIS] = math.Cos(self.Heading() * math.Pi / 180)
	self.velocity[ZAXIS] = -math.Sin(self.Heading() * math.Pi / 180)
}

func (self *MobData) Snapx(x float64, vx float64) {
	self.position[XAXIS] = x
	self.velocity[XAXIS] = vx
}

func (self *MobData) Snapz(z float64, vz float64) {
	self.position[ZAXIS] = z
	self.velocity[ZAXIS] = vz
}

func (self *MobData) Snapy(y float64, vy float64) {
	self.position[YAXIS] = y
	self.velocity[YAXIS] = vy
}

func (self *MobData) Setvx(vx float64) {
	self.velocity[XAXIS] = vx
}

func (self *MobData) Setvz(vz float64) {
	self.velocity[ZAXIS] = vz
}

func (self *MobData) Setvy(vy float64) {
	self.velocity[YAXIS] = vy
}

func (self *MobData) FrontBlock() Vectori {
	ip := IntPosition(self.Position())
	if self.heading > 337.5 || self.heading <= 22.5 {
		ip[XAXIS]++
	} else if self.heading > 22.5 && self.heading <= 67.5 {
		ip[XAXIS]++
		ip[ZAXIS]--
	} else if self.heading > 67.5 && self.heading <= 112.5 {
		ip[ZAXIS]--
	} else if self.heading > 112.5 && self.heading <= 157.5 {
		ip[XAXIS]--
		ip[ZAXIS]--
	} else if self.heading > 157.5 && self.heading <= 202.5 {
		ip[XAXIS]--
	} else if self.heading > 202.5 && self.heading <= 247.5 {
		ip[XAXIS]--
		ip[ZAXIS]++
	} else if self.heading > 247.5 && self.heading <= 292.5 {
		ip[ZAXIS]++
	} else if self.heading > 292.5 && self.heading <= 337.5 {
		ip[XAXIS]++
		ip[ZAXIS]++
	}

	return ip
}

// Return true if the mob is facing the point indicated by the vector
func (self *MobData) Facing(v Vectorf) bool {
	return self.Normal().Dot(v) > 0
}

func (self *MobData) DistanceTo(v Vectorf) float64 {
	return self.position.Minus(v).Magnitude()
}

func (self *MobData) Update(dt float64) (completed bool) {
	self.position[XAXIS] += self.velocity[XAXIS] * dt
	self.position[YAXIS] += self.velocity[YAXIS] * dt
	self.position[ZAXIS] += self.velocity[ZAXIS] * dt

	for i := 0; i < 3; i++ {
		if math.Abs(self.velocity[i]) < 0.02 {
			self.velocity[i] = 0
		}
	}

	if self.velocity[XAXIS] != 0 || self.velocity[ZAXIS] != 0 {
		self.walkSequence += 2 * math.Pi * dt
		if self.walkSequence > 2*math.Pi {
			self.walkSequence -= 2 * math.Pi
		}
	} else {
		self.walkSequence = 0
	}

	return false
}

type Mob interface {
	W() float64
	H() float64
	D() float64

	IsFalling() bool
	Velocity() Vectorf
	Position() Vectorf
	Snapx(x float64, vx float64)
	Snapy(y float64, vy float64)
	Snapz(z float64, vz float64)
	Setvx(vx float64)
	Setvy(vy float64)
	Setvz(vz float64)
	SetFalling(b bool)
	Rotate(angle float64)
	Act(dt float64)
	Draw(pos Vectorf, selectedBlockFace *BlockFace)
	Update(dt float64) bool
}
