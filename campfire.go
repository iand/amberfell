/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

const (
	CAMPFIRE_INTENSITY = 6
	CAMPFIRE_DURATION  = 10
)

type CampFire struct {
	pos  Vectori
	life float64
}

func NewCampFire(pos Vectori) *CampFire {
	return &CampFire{pos: pos, life: CAMPFIRE_DURATION}
}

func (self *CampFire) Intensity() uint16 {
	return CAMPFIRE_INTENSITY
}

func (self *CampFire) Update(dt float64) (completed bool) {
	completed = false
	self.life -= 0.02 * dt
	if self.life <= 0 {
		TheWorld.Setv(self.pos, BLOCK_AIR)
		delete(TheWorld.lightSources, self.pos)
		TheWorld.InvalidateRadius(self.pos[XAXIS], self.pos[ZAXIS], CAMPFIRE_INTENSITY)
		completed = true
	}

	return completed
}

func (self *CampFire) TargetType() uint8 {
	return TARGET_CAMPFIRE
}

func (self *CampFire) Position() Vectorf {
	return self.pos.Vectorf()
}

func (self *CampFire) Velocity() Vectorf {
	return Vectorf{}
}

func (self *CampFire) ApplyDamage(damage float64) {
	// NOOP
}
