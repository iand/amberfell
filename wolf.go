/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/gl"
	"math"
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
	wolf.walkingSpeed = 12
	wolf.sprintSpeed = 25

	wolf.fullEnergy = 30
	wolf.energy = wolf.fullEnergy
	wolf.fullHealth = 100
	wolf.health = wolf.fullHealth
	wolf.healingRate = 1
	wolf.attackStrength = 5

	wolf.mass = 5
	wolf.behaviours = []MobBehaviour{

		MobBehaviour{
			behaviour:   BEHAVIOUR_EVADE,
			targetType:  TARGET_CAMPFIRE,
			targetRange: 14,
			targetAngle: 180,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      2,
			last:        false,
		},

		MobBehaviour{
			behaviour:   BEHAVIOUR_PURSUE,
			targetType:  TARGET_PLAYER,
			targetRange: 12,
			targetAngle: 180,
			sunlight:    SUNLIGHT_LEVELS_ANY,
			weight:      2,
			last:        false,
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
			weight:      2,
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
	pos := Vectorf{self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS]}
	gl.PushMatrix()

	gl.Translated(self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS])
	gl.Rotated(self.Heading(), 0.0, 1.0, 0.0)

	// Translate to top of ground block
	gl.Translatef(0.0, -0.5, 0.0)
	pos[YAXIS] += -0.5

	headHeight := 0.25
	headWidth := headHeight
	headDepth := headHeight * 2.0
	neckHeight := 0.0
	torsoWidth := 1.5 * headHeight
	torsoHeight := 1.5 * headHeight
	torsoDepth := 5 * headHeight
	legHeight := 5*headHeight - torsoHeight - neckHeight - headHeight
	legWidth := (torsoWidth - 0.25*headHeight) / 2
	legDepth := legWidth
	// lowerArmHeight := 1.25 * headHeight
	// handHeight := 0.75 * headHeight

	var legAngle, step float64

	horzSpeed := self.velocity[XAXIS]*self.velocity[XAXIS] + self.velocity[ZAXIS]*self.velocity[ZAXIS]
	legAngle = math.Sin(self.walkSequence) * (15 + 55*horzSpeed/(self.sprintSpeed*self.sprintSpeed))
	headAngle := 30.0
	// torsoAngle = -math.Abs(legAngle / 6)
	step = headHeight * 0.3 * math.Pow(math.Sin(self.walkSequence), 2)

	gl.Translated(0.0, step, 0)
	pos[YAXIS] += step

	// Translate to top of leg
	// Translate to centre of front left leg
	gl.Translated(0.0, legHeight, 0)
	pos[YAXIS] += legHeight

	legDepthOffset := torsoDepth/2 - legWidth/2
	legHeightOffset := -legHeight / 2
	legWidthOffset := (legWidth + 0.25*headHeight) / 2

	// Translate to centre of front left leg
	gl.Translated(legDepthOffset, 0, legWidthOffset)
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0, legHeightOffset, 0)
	pos[XAXIS] += legDepthOffset
	pos[YAXIS] += legHeightOffset
	pos[ZAXIS] += legWidthOffset
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], FACE_NONE)
	pos[XAXIS] -= legDepthOffset
	pos[YAXIS] -= legHeightOffset
	pos[ZAXIS] -= legWidthOffset
	gl.Translated(0, -legHeightOffset, 0)
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(-legDepthOffset, 0, -legWidthOffset)

	legWidthOffset = -legWidthOffset
	if horzSpeed <= self.walkingSpeed*self.walkingSpeed {
		legAngle = -legAngle
	}

	gl.Translated(legDepthOffset, 0, legWidthOffset)
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0, legHeightOffset, 0)
	pos[XAXIS] += legDepthOffset
	pos[YAXIS] += legHeightOffset
	pos[ZAXIS] += legWidthOffset
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], FACE_NONE)
	pos[XAXIS] -= legDepthOffset
	pos[YAXIS] -= legHeightOffset
	pos[ZAXIS] -= legWidthOffset
	gl.Translated(0, -legHeightOffset, 0)
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(-legDepthOffset, 0, -legWidthOffset)

	legDepthOffset = -legDepthOffset
	legWidthOffset = -legWidthOffset

	if horzSpeed > self.walkingSpeed*self.walkingSpeed {
		legAngle = -legAngle
	}

	gl.Translated(legDepthOffset, 0, legWidthOffset)
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0, legHeightOffset, 0)
	pos[XAXIS] += legDepthOffset
	pos[YAXIS] += legHeightOffset
	pos[ZAXIS] += legWidthOffset
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], FACE_NONE)
	pos[XAXIS] -= legDepthOffset
	pos[YAXIS] -= legHeightOffset
	pos[ZAXIS] -= legWidthOffset
	gl.Translated(0, -legHeightOffset, 0)
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(-legDepthOffset, 0, -legWidthOffset)

	legWidthOffset = -legWidthOffset
	if horzSpeed <= self.walkingSpeed*self.walkingSpeed {
		legAngle = -legAngle
	}

	gl.Translated(legDepthOffset, 0, legWidthOffset)
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0, legHeightOffset, 0)
	pos[XAXIS] += legDepthOffset
	pos[YAXIS] += legHeightOffset
	pos[ZAXIS] += legWidthOffset
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], textures[TEXTURE_WOLF_LEG], FACE_NONE)
	pos[XAXIS] -= legDepthOffset
	pos[YAXIS] -= legHeightOffset
	pos[ZAXIS] -= legWidthOffset
	gl.Translated(0, -legHeightOffset, 0)
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(-legDepthOffset, 0, -legWidthOffset)

	//gl.Rotated(torsoAngle, 0.0, 0.0, 1.0)
	// Translate to centre of torso
	gl.Translated(0.0, torsoHeight/2, 0.0)
	pos[YAXIS] += torsoHeight / 2
	Cuboid(pos, torsoWidth, torsoHeight, torsoDepth, textures[TEXTURE_WOLF_TORSO_FRONT], textures[TEXTURE_WOLF_TORSO_BACK], textures[TEXTURE_WOLF_TORSO_SIDE], textures[TEXTURE_WOLF_TORSO_SIDE], textures[TEXTURE_WOLF_TORSO_TOP], textures[TEXTURE_WOLF_TORSO_TOP], FACE_NONE)

	// Translate to shoulders
	gl.Translated(0.0, torsoHeight/2, 0.0)
	pos[YAXIS] += torsoHeight / 2

	// Translate to centre of head
	gl.Translated(torsoDepth/2+headDepth*0.5, 0.0, 0.0)
	pos[XAXIS] += torsoDepth/2 + headDepth*0.5
	pos[YAXIS] += 0.0

	gl.Rotated(-headAngle, 0.0, 0.0, 1.0)
	Cuboid(pos, headWidth, headHeight, headDepth, textures[TEXTURE_WOLF_HEAD_FRONT], textures[TEXTURE_WOLF_HEAD_BACK], textures[TEXTURE_WOLF_HEAD_SIDE], textures[TEXTURE_WOLF_HEAD_SIDE], textures[TEXTURE_WOLF_HEAD_TOP], textures[TEXTURE_WOLF_HEAD_BOTTOM], FACE_NONE)

	gl.PopMatrix()

}

func (self *Wolf) TargetType() uint8 {
	return TARGET_WOLF
}
