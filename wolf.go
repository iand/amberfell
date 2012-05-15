/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
)

type Wolf struct {
	MobData
}

func NewWolf(heading float64, x, y, z uint16) *Wolf {
	wolf := &Wolf{}
	wolf.heading = heading
	wolf.position[XAXIS] = float64(x)
	wolf.position[YAXIS] = float64(y)
	wolf.position[ZAXIS] = float64(z)
	wolf.walkingSpeed = 14
	wolf.sprintSpeed = 25
	wolf.stamina = 30
	wolf.energy = wolf.stamina
	wolf.mass = 5
	wolf.behaviours = []MobBehaviour{
		MobBehaviour{
			behaviour:   BEHAVIOUR_PURSUE,
			targetType:  TARGET_PLAYER,
			targetRange: 12,
			targetAngle: 180,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      2,
			last:        true,
		},
		MobBehaviour{
			behaviour:   BEHAVIOUR_PURSUE,
			targetType:  TARGET_PLAYER,
			targetRange: 22,
			targetAngle: 120,
			sunlight:    SUNLIGHT_LEVELS_NIGHT,
			weight:      3,
			last:        true,
		},
		MobBehaviour{
			behaviour:   BEHAVIOUR_SEPARATE,
			targetType:  TARGET_WOLF,
			targetRange: 10,
			targetAngle: 120,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      1,
			last:        false,
		},
		MobBehaviour{
			behaviour:   BEHAVIOUR_GATHER,
			targetType:  TARGET_WOLF,
			targetRange: 30,
			targetAngle: 120,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      1,
			last:        false,
		},
		MobBehaviour{
			behaviour:   BEHAVIOUR_ALIGN,
			targetType:  TARGET_WOLF,
			targetRange: 30,
			targetAngle: 120,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      1,
			last:        false,
		},
		MobBehaviour{
			behaviour:   BEHAVIOUR_WANDER,
			targetType:  TARGET_ANY,
			targetRange: 0,
			targetAngle: 0,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      1,
			last:        false,
		},
	}

	return wolf
}

func (self *Wolf) W() float64 { return 2 }
func (self *Wolf) H() float64 { return 2 }
func (self *Wolf) D() float64 { return 1 }

func (self *Wolf) Act2(dt float64) {

	// Behaviour

	// Determine intention: hunt, flee or wander

	// Look for a threat (campfires)
	// If threat nearby then evade it 

	// Look for a quarry (e.g. player, possibly other animals)
	// If quarry nearby then pursue it (arrival)

	// Otherwise wander

	// Look for other wolves in vicinity (angle of vision, distance or neighbourhood)
	// Maintain a minimum separation from others
	// Maintain cohesion with others
	// Maintain alignment with others

	// Need angle, distance and weight for each of above 3 behaviours

	// FIVE COMPONENTS
	// Behaviour              4 bits pursue, evade, separate, cohese, align 	
	// Target block type      8 bits + 
	// Range                  6 bits
	// Angle of range         7 bits
	// Look ahead prediction time  4 bits
	// Environment (day/night)
	// Terminate sequence?
	// TOTAL: 29bits

}

func (self *Wolf) Draw(center Vectorf, selectedBlockFace *BlockFace) {
	// println("drawing at ", self.position[XAXIS], self.position[ZAXIS])
	pos := Vectorf{self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS]}
	gl.PushMatrix()

	gl.Translated(self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS])
	gl.Rotated(self.Heading(), 0.0, 1.0, 0.0)

	// Translate to top of ground block
	gl.Translatef(0.0, -0.5, 0.0)
	pos[YAXIS] += -0.5

	Cuboid(pos, 1, 0.5, 1, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	gl.PopMatrix()
}

func (self *Wolf) TargetType() uint8 {
	return TARGET_WOLF
}
