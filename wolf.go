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
	wolf.stamina = 30
	wolf.energy = wolf.stamina
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
	// pos := Vectorf{self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS]}
	// gl.PushMatrix()

	// gl.Translated(self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS])
	// gl.Rotated(self.Heading(), 0.0, 1.0, 0.0)

	// WolfModel.GLDraw()

	// // // Translate to top of ground block
	// // gl.Translatef(0.0, -0.5, 0.0)
	// // pos[YAXIS] += -0.5

	// // Cuboid(pos, 1, 0.5, 1, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	// gl.PopMatrix()

	pos := Vectorf{self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS]}
	gl.PushMatrix()

	gl.Translated(self.position[XAXIS], self.position[YAXIS], self.position[ZAXIS])
	gl.Rotated(self.Heading(), 0.0, 1.0, 0.0)

	// Translate to top of ground block
	gl.Translatef(0.0, -0.5, 0.0)
	pos[YAXIS] += -0.5

	// From http://www.realcolorwheel.com/human.htm
	headHeight := float64(0.25)
	hatHeight := headHeight
	brimHeight := 0.15 * headHeight
	brimWidth := 1.5 * headHeight
	brimDepth := 1.5 * headHeight
	neckHeight := 0.25 * headHeight
	torsoWidth := 2 * headHeight
	torsoHeight := 3.25 * headHeight
	torsoDepth := 1 * headHeight
	legHeight := 8*headHeight - torsoHeight - neckHeight - headHeight
	legWidth := (torsoWidth - 0.25*headHeight) / 2
	legDepth := legWidth
	armHeight := 2.75 * headHeight
	armWidth := 0.75 * torsoDepth
	armDepth := 0.75 * torsoDepth
	// lowerArmHeight := 1.25 * headHeight
	// handHeight := 0.75 * headHeight

	var legAngle, torsoAngle, leftArmAngle, rightArmAngle, step float64

	horzSpeed := self.velocity[XAXIS]*self.velocity[XAXIS] + self.velocity[ZAXIS]*self.velocity[ZAXIS]
	legAngle = math.Sin(self.walkSequence) * (30 + 35*horzSpeed/(self.sprintSpeed*self.sprintSpeed))

	torsoAngle = -math.Abs(legAngle / 6)
	leftArmAngle = -legAngle * 1.2
	rightArmAngle = legAngle * 1.2
	step = headHeight * 0.1 * math.Pow(math.Sin(self.walkSequence), 2)

	gl.Translated(0.0, step, 0)
	pos[YAXIS] += step

	// Translate to top of leg
	// Translate to centre of front left leg
	gl.Translated(0.0, legHeight, 0)
	pos[YAXIS] += legHeight

	// Translate to centre of front left leg
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -legHeight/2, (legWidth+0.25*headHeight)/2)
	pos[YAXIS] += -legHeight / 2
	pos[ZAXIS] += (legWidth + 0.25*headHeight) / 2
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	gl.Translated(0.0, legHeight/2, -(legWidth+0.25*headHeight)/2)
	pos[YAXIS] += legHeight / 2
	pos[ZAXIS] += -(legWidth + 0.25*headHeight) / 2
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)

	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -legHeight/2, -(legWidth+0.25*headHeight)/2)
	pos[YAXIS] += -legHeight / 2
	pos[ZAXIS] += -(legWidth + 0.25*headHeight) / 2
	Cuboid(pos, legWidth, legHeight, legDepth, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	gl.Translated(0.0, legHeight/2, (legWidth+0.25*headHeight)/2)
	pos[YAXIS] += legHeight / 2
	pos[ZAXIS] += (legWidth + 0.25*headHeight) / 2
	gl.Rotated(+legAngle, 0.0, 0.0, 1.0)

	gl.Rotated(torsoAngle, 0.0, 0.0, 1.0)
	// Translate to centre of torso
	gl.Translated(0.0, torsoHeight/2, 0.0)
	pos[YAXIS] += torsoHeight / 2
	Cuboid(pos, torsoWidth, torsoHeight, torsoDepth, textures[TEXTURE_TORSO_FRONT], textures[TEXTURE_TORSO_BACK], textures[TEXTURE_TORSO_LEFT], textures[TEXTURE_TORSO_RIGHT], textures[TEXTURE_TORSO_TOP], textures[TEXTURE_TORSO_TOP], FACE_NONE)

	// Translate to shoulders
	gl.Translated(0.0, torsoHeight/2, 0.0)
	pos[YAXIS] += torsoHeight / 2

	gl.Rotated(leftArmAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -armHeight/2, torsoWidth/2+armWidth/2)
	pos[YAXIS] += -armHeight / 2
	pos[ZAXIS] += torsoWidth/2 + armWidth/2
	Cuboid(pos, armWidth, armHeight, armDepth, textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM_TOP], textures[TEXTURE_HAND], FACE_NONE)
	gl.Translated(0.0, armHeight/2, -torsoWidth/2-armWidth/2)
	pos[YAXIS] += armHeight / 2
	pos[ZAXIS] += -torsoWidth/2 + armWidth/2
	gl.Rotated(-leftArmAngle, 0.0, 0.0, 1.0)

	gl.Rotated(rightArmAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -armHeight/2, -torsoWidth/2-armWidth/2)
	pos[YAXIS] += -armHeight / 2
	pos[ZAXIS] += -torsoWidth/2 + armWidth/2
	Cuboid(pos, armWidth, armHeight, armDepth, textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM_TOP], textures[TEXTURE_HAND], FACE_NONE)
	gl.Translated(0.0, armHeight/2, torsoWidth/2+armWidth/2)
	pos[YAXIS] += armHeight / 2
	pos[ZAXIS] += torsoWidth/2 + armWidth/2
	gl.Rotated(-rightArmAngle, 0.0, 0.0, 1.0)

	// Translate to centre of head
	gl.Translated(0.0, neckHeight+headHeight/2, 0.0)
	pos[YAXIS] += neckHeight + headHeight/2

	if selectedBlockFace != nil {
		blockPos := selectedBlockFace.pos.Vectorf()
		headPos := self.position.Add(Vectorf{0, headHeight * 9, 0})

		blockDir := blockPos.Minus(headPos)
		if self.Facing(blockDir) {
			yrot := (math.Atan2(blockDir[XAXIS], blockDir[ZAXIS]) - math.Pi/2) * 180 / math.Pi
			zrot, xrot := -12.0, -12.0
			gl.Rotated(-self.Heading(), 0.0, 1.0, 0.0)
			gl.Rotated(yrot, 0.0, 1.0, 0.0)
			gl.Rotated(zrot, 0.0, 0.0, 1.0)
			gl.Rotated(xrot, 1.0, 0.0, 0.0)
		}
	}

	Cuboid(pos, headHeight, headHeight, headHeight, textures[TEXTURE_HEAD_FRONT], textures[TEXTURE_HEAD_BACK], textures[TEXTURE_HEAD_LEFT], textures[TEXTURE_HEAD_RIGHT], nil, textures[TEXTURE_HEAD_BOTTOM], FACE_NONE)

	// Translate to hat brim
	gl.Translated(0.0, headHeight/2+brimHeight/2, 0.0)
	pos[YAXIS] += headHeight/2 + brimHeight/2
	Cuboid(pos, brimWidth, brimHeight, brimDepth, textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], FACE_NONE)

	gl.Translated(0.0, brimHeight/2+hatHeight/2, 0.0)
	pos[YAXIS] += headHeight/2 + brimHeight/2
	Cuboid(pos, hatHeight, hatHeight, hatHeight, textures[TEXTURE_HAT_FRONT], textures[TEXTURE_HAT_BACK], textures[TEXTURE_HAT_LEFT], textures[TEXTURE_HAT_RIGHT], textures[TEXTURE_HAT_TOP], nil, FACE_NONE)

	gl.PopMatrix()

}

func (self *Wolf) TargetType() uint8 {
	return TARGET_WOLF
}
