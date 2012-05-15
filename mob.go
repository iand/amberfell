/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"math"
	"math/rand"
)

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
	TargetType() uint8
}

type MobData struct {
	heading           float64
	position          Vectorf
	velocity          Vectorf
	falling           bool
	walkingSpeed      float64
	sprintSpeed       float64
	mass              float64
	stamina           float64
	energy            float64
	walkSequence      float64
	behaviours        []MobBehaviour
	dominantBehaviour uint8
}

type MobBehaviour struct {
	behaviour   uint8
	targetType  uint8
	targetRange uint8
	targetAngle uint8
	weight      uint8
	sunlight    uint8
	last        bool
}

const SUNLIGHT_LEVELS_LOWER_MASK = 0xF0
const SUNLIGHT_LEVELS_UPPER_MASK = 0x0F
const SUNLIGHT_LEVELS_ANY = 0x08
const SUNLIGHT_LEVELS_NIGHT = 0x05
const SUNLIGHT_LEVELS_FULL_NIGHT = 0x01
const SUNLIGHT_LEVELS_DAY = 0x58

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
	if self.velocity.Magnitude() > self.walkingSpeed {
		self.energy -= 1 * dt
	} else {
		self.energy += 1 * dt
	}

	if self.energy <= 0 {
		self.energy = 0
	} else if self.energy > self.stamina {
		self.energy = self.stamina
	}

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

func (self *MobData) Act(dt float64) {

	previousBehaviour := self.dominantBehaviour
	self.dominantBehaviour = BEHAVIOUR_WANDER

	const MAX_FORCE = 100 // to be calibrated
	force := Vectorf{}

	normal := self.Normal()
	for _, behaviour := range self.behaviours {

		lowerLightLevel := int(behaviour.sunlight & SUNLIGHT_LEVELS_LOWER_MASK)
		upperLightLevel := int(behaviour.sunlight & SUNLIGHT_LEVELS_UPPER_MASK)

		if sunlightLevel >= lowerLightLevel && sunlightLevel < upperLightLevel {
			weight := float64(behaviour.weight)
			if weight == 0 {
				weight = 1
			}

			switch behaviour.behaviour {
			case BEHAVIOUR_WANDER:
				offset := self.walkingSpeed / 2
				angle := rand.Float64() * 2 * math.Pi
				angleDir := Vectorf{math.Cos(angle),
					0,
					-math.Sin(angle),
				}

				turnDir := normal.Scale(offset).Add(angleDir).Normalize()
				force = force.Add(turnDir.Scale(self.walkingSpeed / 5).Scale(weight))
				if behaviour.last {
					break
				}

			case BEHAVIOUR_PURSUE:
				if previousBehaviour == BEHAVIOUR_PURSUE && self.energy > 5 || self.energy > 15 {
					offset := ThePlayer.Position().Minus(self.position)
					separation := offset.Magnitude()
					direction := offset.Normalize()
					angle := normal.AngleNormalized(direction) * 180 / math.Pi

					if (angle >= 360-float64(behaviour.targetAngle) || angle <= float64(behaviour.targetAngle)) && separation <= float64(behaviour.targetRange) {
						pos := ThePlayer.position.Add(ThePlayer.velocity.Scale(separation * 0.01))
						desiredVelocity := pos.Minus(self.position).Normalize().Scale(self.sprintSpeed)
						force = force.Add(desiredVelocity.Minus(self.velocity).Scale(weight))
						self.dominantBehaviour = BEHAVIOUR_PURSUE
						if behaviour.last {
							break
						}
					}
				}

			case BEHAVIOUR_EVADE:
				if previousBehaviour == BEHAVIOUR_EVADE && self.energy > 5 || self.energy > 15 {
					offset := ThePlayer.Position().Minus(self.position)
					separation := offset.Magnitude()
					direction := offset.Normalize()
					angle := normal.AngleNormalized(direction) * 180 / math.Pi
					if (angle >= 360-float64(behaviour.targetAngle) || angle <= float64(behaviour.targetAngle)) && separation <= float64(behaviour.targetRange) {
						pos := ThePlayer.position.Add(ThePlayer.velocity.Scale(separation * 0.01))
						desiredVelocity := pos.Minus(self.position).Normalize().Scale(-self.sprintSpeed)
						force = force.Add(desiredVelocity.Minus(self.velocity).Scale(weight))
						self.dominantBehaviour = BEHAVIOUR_EVADE
						if behaviour.last {
							break
						}
					}
				}

			case BEHAVIOUR_SEPARATE:
				for _, mob := range TheWorld.mobs {
					offset := mob.Position().Minus(self.position)
					separation := offset.Magnitude()
					direction := offset.Normalize()
					angle := normal.AngleNormalized(direction) * 180 / math.Pi
					if separation > 0 && (angle >= 360-float64(behaviour.targetAngle) || angle <= float64(behaviour.targetAngle)) && separation <= float64(behaviour.targetRange) {
						force = force.Add(direction.Scale(-self.walkingSpeed / separation).Scale(weight))
					}
					if behaviour.last {
						break
					}
				}

			case BEHAVIOUR_GATHER:
				x, z, count := 0.0, 0.0, 0.0

				for _, mob := range TheWorld.mobs {
					if mob.TargetType() == behaviour.targetType {
						pos := mob.Position()
						offset := pos.Minus(self.position)
						separation := offset.Magnitude()
						direction := offset.Normalize()
						angle := normal.AngleNormalized(direction) * 180 / math.Pi
						if (angle >= 360-float64(behaviour.targetAngle) || angle <= float64(behaviour.targetAngle)) && separation <= float64(behaviour.targetRange) {
							x += pos[XAXIS]
							z += pos[ZAXIS]
							count++
						}
					}
				}

				if count > 0 {
					x /= count
					z /= count
					direction := Vectorf{x, 0, z}.Minus(self.position).Normalize()
					force = force.Add(direction.Scale(self.walkingSpeed).Scale(weight))
					if behaviour.last {
						break
					}
				}

			case BEHAVIOUR_ALIGN:
				x, z, count := 0.0, 0.0, 0.0

				for _, mob := range TheWorld.mobs {
					if mob.TargetType() == behaviour.targetType {
						pos := mob.Position()
						offset := pos.Minus(self.position)
						separation := offset.Magnitude()
						direction := offset.Normalize()
						angle := normal.AngleNormalized(direction) * 180 / math.Pi
						if (angle >= 360-float64(behaviour.targetAngle) || angle <= float64(behaviour.targetAngle)) && separation <= float64(behaviour.targetRange) {
							x += mob.Velocity()[XAXIS]
							z += mob.Velocity()[ZAXIS]
							count++
						}
					}
				}

				if count > 0 {
					x /= count
					z /= count
					force = force.Add(Vectorf{x, 0, z}.Minus(self.velocity).Scale(weight))
					if behaviour.last {
						break
					}
				}

			}
		}

	}

	// Force into 2 dimensions
	force[YAXIS] = 0

	force_magnitude := force.Magnitude()

	if force_magnitude > 0 {

		if force_magnitude > MAX_FORCE {
			force = force.Normalize().Scale(MAX_FORCE)
		}

		accel := force.Scale(1.0 / self.mass)
		self.velocity = self.velocity.Add(accel)
		velocity_magnitude := self.velocity.Magnitude()
		if velocity_magnitude > self.sprintSpeed {
			self.velocity = self.velocity.Normalize().Scale(self.sprintSpeed * dt)
		}
	}

	normalizedVel := self.velocity.Normalize()
	if normalizedVel[XAXIS] != 0 {
		self.heading = math.Atan(-normalizedVel[ZAXIS]/normalizedVel[XAXIS]) * 180 / math.Pi
	}

}
