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

type Player struct {
	MobData
	Bounce float64
	// position      Vectorf
	// velocity      Vectorf
	currentAction uint8
	currentItem   uint16
	equippedItems [7]uint16
	inventory     [255]uint16
}

type BlockBreakRecord struct {
	pos   Vectori
	count int
}

func (self *Player) Init(heading float64, x int16, z int16, y int16) {
	self.heading = heading
	self.position[XAXIS] = float64(x)
	self.position[YAXIS] = float64(y)
	self.position[ZAXIS] = float64(z)
	self.walkingSpeed = 20
	self.currentAction = ACTION_HAND
	self.currentItem = ITEM_NONE

	self.equippedItems[0] = BLOCK_DIRT
	self.equippedItems[1] = BLOCK_STONE
	self.equippedItems[2] = ITEM_NONE
	self.equippedItems[3] = ITEM_NONE
	self.equippedItems[4] = ITEM_NONE
	self.equippedItems[5] = ITEM_NONE
	self.equippedItems[6] = ITEM_NONE

}

func (self *Player) W() float64 { return 0.8 }
func (self *Player) H() float64 { return 2.0 }
func (self *Player) D() float64 { return 0.6 }

func (self *Player) Act(dt float64) {
	// noop
}

func (player *Player) Draw(center Vectorf, selectedBlockFace *BlockFace) {

	gl.PushMatrix()

	gl.Translatef(float32(player.position[XAXIS]), float32(player.position[YAXIS]), float32(player.position[ZAXIS]))
	//stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
	gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)

	// Translate to top of ground block
	gl.Translatef(0.0, -0.5, 0.0)


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
	legWidth := (torsoWidth-0.25*headHeight) / 2
	legDepth := legWidth
	armHeight := 2.75 * headHeight
	armWidth := 0.75 * torsoDepth
	armDepth := 0.75 * torsoDepth
	// lowerArmHeight := 1.25 * headHeight
	// handHeight := 0.75 * headHeight

	var legAngle, torsoAngle, leftArmAngle, rightArmAngle, step float64

	if player.velocity[YAXIS] != 0 {
		legAngle = 20
		leftArmAngle = 120
		rightArmAngle = 120
	} else {
		legAngle = 40 * (math.Abs(player.velocity[XAXIS])+math.Abs(player.velocity[ZAXIS])) / player.walkingSpeed * math.Sin(player.walkSequence)
		torsoAngle = -math.Abs(legAngle / 6)
		leftArmAngle = -legAngle * 1.2
		rightArmAngle = legAngle * 1.2
		step = headHeight * 0.1 * math.Pow(math.Sin(player.walkSequence), 2)
	}


	gl.Translated(0.0, step, 0)

	// Translate to top of leg
	gl.Translated(0.0, legHeight, 0)


	// Translate to centre of leg
	gl.Rotated(legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -legHeight / 2, (legWidth + 0.25*headHeight) / 2)
	Cuboid(legWidth, legHeight, legDepth, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	gl.Translated(0.0, legHeight / 2, -(legWidth + 0.25*headHeight) / 2)
	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)

	gl.Rotated(-legAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -legHeight / 2, -(legWidth + 0.25*headHeight) / 2)
	Cuboid(legWidth, legHeight, legDepth, textures[TEXTURE_LEG], textures[TEXTURE_LEG], textures[TEXTURE_LEG_SIDE], textures[TEXTURE_LEG_SIDE], textures[32], textures[32], FACE_NONE)
	gl.Translated(0.0, legHeight / 2, (legWidth + 0.25*headHeight) / 2)
	gl.Rotated(+legAngle, 0.0, 0.0, 1.0)


	gl.Rotated(torsoAngle, 0.0, 0.0, 1.0)
	// Translate to centre of torso
	gl.Translated(0.0, torsoHeight / 2, 0.0)
	Cuboid(torsoWidth, torsoHeight, torsoDepth, textures[TEXTURE_TORSO_FRONT], textures[TEXTURE_TORSO_BACK], textures[TEXTURE_TORSO_LEFT], textures[TEXTURE_TORSO_RIGHT], textures[TEXTURE_TORSO_TOP], textures[TEXTURE_TORSO_TOP], FACE_NONE)

	// Translate to shoulders
	gl.Translated(0.0, torsoHeight / 2, 0.0)

	gl.Rotated(leftArmAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -armHeight/2, torsoWidth/2 + armWidth/2)
	Cuboid(armWidth, armHeight, armDepth, textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM_TOP], textures[TEXTURE_HAND], FACE_NONE)
	gl.Translated(0.0, armHeight/2, -torsoWidth/2 - armWidth/2)
	gl.Rotated(-leftArmAngle, 0.0, 0.0, 1.0)

	gl.Rotated(rightArmAngle, 0.0, 0.0, 1.0)
	gl.Translated(0.0, -armHeight/2, -torsoWidth/2 - armWidth/2)
	Cuboid(armWidth, armHeight, armDepth, textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM], textures[TEXTURE_ARM_TOP], textures[TEXTURE_HAND], FACE_NONE)
	gl.Translated(0.0, armHeight/2, torsoWidth/2 + armWidth/2)
	gl.Rotated(-rightArmAngle, 0.0, 0.0, 1.0)

	// Translate to centre of head
	gl.Translated(0.0, neckHeight + headHeight / 2, 0.0)

	if selectedBlockFace != nil {
		blockPos := selectedBlockFace.pos.Vectorf()
		headPos := player.position.Add(Vectorf{0,headHeight*9,0})
		
		blockDir := blockPos.Minus(headPos)

		yrot := (math.Atan2(blockDir[XAXIS], blockDir[ZAXIS]) - math.Pi/2) * 180 / math.Pi
		zrot, xrot := -12.0, -12.0
		gl.Rotated(-player.Heading(), 0.0, 1.0, 0.0)
		gl.Rotated(yrot, 0.0, 1.0, 0.0)
		gl.Rotated(zrot, 0.0, 0.0, 1.0)
		gl.Rotated(xrot, 1.0, 0.0, 0.0)
	}

	Cuboid(headHeight, headHeight, headHeight, textures[TEXTURE_HEAD_FRONT], textures[TEXTURE_HEAD_BACK], textures[TEXTURE_HEAD_LEFT], textures[TEXTURE_HEAD_RIGHT], nil, textures[TEXTURE_HEAD_BOTTOM], FACE_NONE)

	// Translate to hat brim
	gl.Translated(0.0, headHeight / 2 + brimHeight / 2, 0.0)
	Cuboid(brimWidth, brimHeight, brimDepth, textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], textures[TEXTURE_BRIM], FACE_NONE)

	gl.Translated(0.0, brimHeight / 2 + hatHeight / 2, 0.0)
	Cuboid(hatHeight, hatHeight, hatHeight, textures[TEXTURE_HAT_FRONT], textures[TEXTURE_HAT_BACK], textures[TEXTURE_HAT_LEFT], textures[TEXTURE_HAT_RIGHT], textures[TEXTURE_HAT_TOP], nil, FACE_NONE)



	gl.PopMatrix()
}

func (self *Player) HandleKeys(keys []uint8) {
	if keys[sdl.K_1] != 0 {
		self.currentAction = ACTION_HAND
		self.currentItem = ITEM_NONE
	}
	if keys[sdl.K_2] != 0 {
		self.currentAction = ACTION_BREAK
		self.currentItem = ITEM_NONE
	}
	if keys[sdl.K_3] != 0 {
		self.currentAction = ACTION_WEAPON
		self.currentItem = ITEM_NONE
	}
	if keys[sdl.K_4] != 0 {
		self.currentAction = ACTION_ITEM0
		self.currentItem = self.equippedItems[0]
	}
	if keys[sdl.K_5] != 0 {
		self.currentAction = ACTION_ITEM1
		self.currentItem = self.equippedItems[1]
	}
	if keys[sdl.K_6] != 0 {
		self.currentAction = ACTION_ITEM2
		self.currentItem = self.equippedItems[2]
	}
	if keys[sdl.K_7] != 0 {
		self.currentAction = ACTION_ITEM3
		self.currentItem = self.equippedItems[3]
	}
	if keys[sdl.K_8] != 0 {
		self.currentAction = ACTION_ITEM4
		self.currentItem = self.equippedItems[4]
	}

	if keys[sdl.K_w] != 0 {
		if !self.IsFalling() {
			self.velocity[XAXIS] = math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed
			self.velocity[ZAXIS] = -math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed
		} else {
			self.velocity[XAXIS] = math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed / 3
			self.velocity[ZAXIS] = -math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed / 3
		}

	}
	if keys[sdl.K_s] != 0 {
		if !self.IsFalling() {
			self.velocity[XAXIS] = -math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed / 2
			self.velocity[ZAXIS] = math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed / 2
		} else {
			self.velocity[XAXIS] = -math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed / 6
			self.velocity[ZAXIS] = math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed / 6
		}

	}
	if keys[sdl.K_a] != 0 {
		self.Rotate(22.5 / 2)
		// viewport.Roty(-22.5 / 2)
	}

	if keys[sdl.K_d] != 0 {
		self.Rotate(-22.5 / 2)
		// viewport.Roty(22.5 / 2)
	}

	if keys[sdl.K_SPACE] != 0 {
		if !self.IsFalling() {
			self.velocity[YAXIS] = 4
		}
	}

}

func (self *Player) CanInteract() bool {
	if self.currentAction == ACTION_BREAK || self.currentItem != ITEM_NONE {
		return true
	}
	return false
}

func (self *Player) Interact(interactingBlockFace *InteractingBlockFace) {
	if !self.CanInteract() {
		return
	}

	selectedBlockFace := interactingBlockFace.blockFace
	println("Interacting at ", selectedBlockFace.pos.String())
	switch self.currentAction {
	case ACTION_BREAK:
		blockid := TheWorld.Atv(selectedBlockFace.pos)
		if blockid != BLOCK_AIR {

			println("Hitting ", selectedBlockFace.pos.String())
			interactingBlockFace.hitCount++
			if interactingBlockFace.hitCount >= TerrainBlocks[uint16(blockid)].hitsNeeded {
				println("Breaking at ", selectedBlockFace.pos.String())
				TheWorld.Setv(selectedBlockFace.pos, BLOCK_AIR)
				self.inventory[blockid]++
				interactingBlockFace.hitCount = 0
			}

		}
	case ACTION_ITEM0, ACTION_ITEM1, ACTION_ITEM2, ACTION_ITEM3, ACTION_ITEM4:
		if selectedBlockFace.face == UP_FACE { // top
			selectedBlockFace.pos[YAXIS]++
		} else if selectedBlockFace.face == DOWN_FACE { // bottom
			selectedBlockFace.pos[YAXIS]--
		} else if selectedBlockFace.face == SOUTH_FACE { // front
			selectedBlockFace.pos[ZAXIS]++
		} else if selectedBlockFace.face == NORTH_FACE { // back
			selectedBlockFace.pos[ZAXIS]--
		} else if selectedBlockFace.face == EAST_FACE { // left
			selectedBlockFace.pos[XAXIS]++
		} else if selectedBlockFace.face == WEST_FACE { // right
			selectedBlockFace.pos[XAXIS]--
		}
		if TheWorld.Atv(selectedBlockFace.pos) == BLOCK_AIR {
			println("Setting at ", selectedBlockFace.pos.String())
			TheWorld.Setv(selectedBlockFace.pos, byte(self.currentItem))
		}
	}

}
