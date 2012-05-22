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

type Action uint8
type BlockId uint16
type ItemId BlockId

type Player struct {
	MobData
	currentAction     Action
	currentItem       ItemId
	equippedItems     [7]ItemId
	inventory         [MAX_ITEMS]uint16
	distanceTravelled float64
	distanceFromStart uint16
	interactingBlock  *InteractingBlockFace
}

type BlockBreakRecord struct {
	pos   Vectori
	count int
}

func (self *Player) Init(heading float64, x uint16, z uint16) {
	self.heading = heading
	self.position[XAXIS] = float64(x)
	self.position[YAXIS] = float64(TheWorld.FindSurface(x, z))
	self.position[ZAXIS] = float64(z)
	self.mass = 5
	self.stamina = 30
	self.energy = self.stamina

	self.walkingSpeed = 18
	self.sprintSpeed = 22
	self.currentAction = ACTION_HAND
	self.currentItem = ITEM_NONE

	self.equippedItems[0] = ITEM_NONE
	self.equippedItems[1] = ITEM_NONE
	self.equippedItems[2] = ITEM_NONE
	self.equippedItems[3] = ITEM_NONE
	self.equippedItems[4] = ITEM_NONE
	self.equippedItems[5] = ITEM_NONE
	self.equippedItems[6] = ITEM_NONE

}

func (self *Player) W() float64 { return 2*0.25 + 2*0.75*0.25 } // torso width + 2 arm widths
func (self *Player) H() float64 { return 8*0.25 + 0.25 }        // body+hat 
func (self *Player) D() float64 { return 0.25 }                 // torso depth

func (self *Player) Act(dt float64) {
	// noop
}

func (self *Player) Draw(center Vectorf, selectedBlockFace *BlockFace) {

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
	legAngle = math.Sin(self.walkSequence) * (20 + 45*horzSpeed/(self.sprintSpeed*self.sprintSpeed))

	// legAngle = 55 * (math.Abs(self.velocity[XAXIS]) + math.Abs(self.velocity[ZAXIS])) / self.sprintSpeed * math.Sin(self.walkSequence)
	torsoAngle = -math.Abs(legAngle / 6)
	leftArmAngle = -legAngle * 1.2
	rightArmAngle = legAngle * 1.2
	step = headHeight * 0.1 * math.Pow(math.Sin(self.walkSequence), 2)

	gl.Translated(0.0, step, 0)
	pos[YAXIS] += step

	// Translate to top of leg
	gl.Translated(0.0, legHeight, 0)
	pos[YAXIS] += legHeight

	// Translate to centre of leg
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

func (self *Player) HandleKeys(keys []uint8) {

	if keys[sdl.K_w] != 0 {
		if self.IsFalling() {
			self.velocity[XAXIS] = math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed / 2
			self.velocity[ZAXIS] = -math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed / 2
		} else {
			speed := self.walkingSpeed
			if self.energy > 5 && (keys[sdl.K_LSHIFT] != 0 || keys[sdl.K_RSHIFT] != 0) {
				speed = self.sprintSpeed
			}
			self.velocity[XAXIS] = math.Cos(self.Heading()*math.Pi/180) * speed
			self.velocity[ZAXIS] = -math.Sin(self.Heading()*math.Pi/180) * speed
		}

	}
	if keys[sdl.K_s] != 0 {
		if self.IsFalling() {
			self.velocity[XAXIS] = -math.Cos(self.Heading()*math.Pi/180) * self.walkingSpeed / 4
			self.velocity[ZAXIS] = math.Sin(self.Heading()*math.Pi/180) * self.walkingSpeed / 4
		} else {
			speed := self.walkingSpeed
			if self.energy > 5 && (keys[sdl.K_LSHIFT] != 0 || keys[sdl.K_RSHIFT] != 0) {
				speed = self.sprintSpeed
			}
			self.velocity[XAXIS] = -math.Cos(self.Heading()*math.Pi/180) * speed / 2
			self.velocity[ZAXIS] = math.Sin(self.Heading()*math.Pi/180) * speed / 2
		}

	}

	if keys[sdl.K_q] != 0 {
		if self.IsFalling() {
			self.velocity[XAXIS] = math.Cos((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 3
			self.velocity[ZAXIS] = -math.Sin((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 3
		} else {
			self.velocity[XAXIS] = math.Cos((self.Heading()+90)*math.Pi/180) * self.walkingSpeed
			self.velocity[ZAXIS] = -math.Sin((self.Heading()+90)*math.Pi/180) * self.walkingSpeed
		}

	}
	if keys[sdl.K_e] != 0 {
		if self.IsFalling() {
			self.velocity[XAXIS] = -math.Cos((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 6
			self.velocity[ZAXIS] = math.Sin((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 6
		} else {
			self.velocity[XAXIS] = -math.Cos((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 2
			self.velocity[ZAXIS] = math.Sin((self.Heading()+90)*math.Pi/180) * self.walkingSpeed / 2
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
			self.velocity[YAXIS] = 7
		}
	}

}

func (self *Player) CanInteract() bool {
	if self.currentAction == ACTION_HAND || self.currentAction == ACTION_BREAK || self.currentItem != ITEM_NONE {
		return true
	}
	return false
}

func (self *Player) Interact(interactingBlockFace *InteractingBlockFace) {
	if !self.CanInteract() {
		return
	}

	selectedBlockFace := interactingBlockFace.blockFace
	// println("Interacting at ", selectedBlockFace.pos.String())
	switch self.currentAction {

	case ACTION_HAND:
		blockid := TheWorld.Atv(selectedBlockFace.pos)
		switch blockid {
		case BLOCK_AMBERFELL_PUMP, BLOCK_STEAM_GENERATOR, BLOCK_AMBERFELL_CONDENSER, BLOCK_FURNACE, BLOCK_BEESNEST:
			if obj, ok := TheWorld.containerObjects[selectedBlockFace.pos]; ok {
				inventory.Show(obj, nil)
			}
		case BLOCK_CARPENTERS_BENCH:
			inventory.Show(nil, NewCarpentersBench(selectedBlockFace.pos))
		case BLOCK_FORGE:
			inventory.Show(nil, NewForge(selectedBlockFace.pos))

		}

	case ACTION_BREAK:
		block := TheWorld.AtBv(selectedBlockFace.pos)
		blocktype := items[ItemId(block.id)]
		blockid := block.id
		if blockid != BLOCK_AIR {
			interactingBlockFace.hitCount++
			if blocktype.hitsNeeded != STRENGTH_UNBREAKABLE && interactingBlockFace.hitCount >= blocktype.hitsNeeded/2 {

				if interactingBlockFace.hitCount >= blocktype.hitsNeeded {
					TheWorld.Setv(selectedBlockFace.pos, BLOCK_AIR)
					TheWorld.InvalidateRadius(selectedBlockFace.pos[XAXIS], selectedBlockFace.pos[ZAXIS], 1)

					switch blockid {
					case BLOCK_CAMPFIRE:
						delete(TheWorld.campfires, selectedBlockFace.pos)
						delete(TheWorld.lightSources, selectedBlockFace.pos)
						delete(TheWorld.timedObjects, selectedBlockFace.pos)

						TheWorld.InvalidateRadius(selectedBlockFace.pos[XAXIS], selectedBlockFace.pos[ZAXIS], uint16(CAMPFIRE_INTENSITY))

					case BLOCK_AMBERFELL_PUMP:
						delete(TheWorld.timedObjects, selectedBlockFace.pos)
						delete(TheWorld.containerObjects, selectedBlockFace.pos)

					case BLOCK_STEAM_GENERATOR:
						delete(TheWorld.timedObjects, selectedBlockFace.pos)
						delete(TheWorld.containerObjects, selectedBlockFace.pos)
						delete(TheWorld.generatorObjects, selectedBlockFace.pos)

					case BLOCK_AMBERFELL_CONDENSER:
						delete(TheWorld.timedObjects, selectedBlockFace.pos)
						delete(TheWorld.containerObjects, selectedBlockFace.pos)

					case BLOCK_FURNACE:
						delete(TheWorld.timedObjects, selectedBlockFace.pos)
						delete(TheWorld.containerObjects, selectedBlockFace.pos)

					case BLOCK_BEESNEST:
						delete(TheWorld.timedObjects, selectedBlockFace.pos)
						delete(TheWorld.containerObjects, selectedBlockFace.pos)

					}

					if blocktype.drops != nil {
						droppedItem := blocktype.drops.item
						if self.inventory[droppedItem] < MAX_ITEMS_IN_INVENTORY {
							self.inventory[droppedItem]++
							if items[droppedItem].placeable {
								self.EquipItem(droppedItem)
							}
						}
					}
					interactingBlockFace.hitCount = 0
				} else if !block.Damaged() {
					block.SetDamaged(true)
					TheWorld.SetBv(selectedBlockFace.pos, block)

				}
			}

		}
	case ACTION_ITEM0, ACTION_ITEM1, ACTION_ITEM2, ACTION_ITEM3, ACTION_ITEM4:
		if self.inventory[self.currentItem] > 0 && items[self.currentItem].placeable {
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
			if TheWorld.Atv(selectedBlockFace.pos) == BLOCK_AIR && self.currentItem < 256 {
				blockid := self.currentItem

				switch blockid {
				case BLOCK_CAMPFIRE:
					// Add a light source

					campfire := NewCampFire(selectedBlockFace.pos)
					TheWorld.lightSources[selectedBlockFace.pos] = campfire
					TheWorld.timedObjects[selectedBlockFace.pos] = campfire
					TheWorld.campfires[selectedBlockFace.pos] = campfire

					TheWorld.InvalidateRadius(selectedBlockFace.pos[XAXIS], selectedBlockFace.pos[ZAXIS], uint16(CAMPFIRE_INTENSITY))

				case BLOCK_AMBERFELL_PUMP:
					sourced := false
					if selectedBlockFace.pos[YAXIS] > 0 && TheWorld.At(selectedBlockFace.pos[XAXIS], selectedBlockFace.pos[YAXIS]-1, selectedBlockFace.pos[ZAXIS]) == BLOCK_AMBERFELL_SOURCE {
						sourced = true
					}

					pump := NewAmberfellPump(selectedBlockFace.pos, sourced, false)
					TheWorld.timedObjects[selectedBlockFace.pos] = pump
					TheWorld.containerObjects[selectedBlockFace.pos] = pump

				case BLOCK_STEAM_GENERATOR:
					gen := NewSteamGenerator(selectedBlockFace.pos)
					TheWorld.timedObjects[selectedBlockFace.pos] = gen
					TheWorld.containerObjects[selectedBlockFace.pos] = gen
					TheWorld.generatorObjects[selectedBlockFace.pos] = gen

				case BLOCK_AMBERFELL_CONDENSER:
					obj := NewAmberfellCondenser(selectedBlockFace.pos)
					TheWorld.timedObjects[selectedBlockFace.pos] = obj
					TheWorld.containerObjects[selectedBlockFace.pos] = obj

				case BLOCK_FURNACE:
					obj := NewFurnace(selectedBlockFace.pos)
					TheWorld.timedObjects[selectedBlockFace.pos] = obj
					TheWorld.containerObjects[selectedBlockFace.pos] = obj
				}

				orientation := HeadingToOrientation(self.heading)
				block := NewBlock(BlockId(blockid), false, orientation)

				TheWorld.SetBv(selectedBlockFace.pos, block)
				self.inventory[self.currentItem]--

			}
		}
	}

}

func (self *Player) HandleMouseButton(re *sdl.MouseButtonEvent) {
	if re.Button == 1 && re.State == 1 { // LEFT, DOWN
		if self.CanInteract() {
			selectedBlockFace := viewport.SelectedBlockFace()
			if selectedBlockFace != nil {
				if self.interactingBlock == nil || self.interactingBlock.blockFace.pos != selectedBlockFace.pos {
					self.interactingBlock = new(InteractingBlockFace)
					self.interactingBlock.blockFace = selectedBlockFace
					self.interactingBlock.hitCount = 0
				}
				self.Interact(self.interactingBlock)
			}
			// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
		}
	}

}

func (self *Player) HandleKeyboard(re *sdl.KeyboardEvent) {

}

func (self *Player) SelectAction(action int) {
	switch action {
	case 0:
		self.currentAction = ACTION_HAND
		self.currentItem = ITEM_NONE
	case 1:
		self.currentAction = ACTION_BREAK
		self.currentItem = ITEM_NONE
	case 2:
		self.currentAction = ACTION_WEAPON
		self.currentItem = ITEM_NONE
	case 3:
		self.currentAction = ACTION_ITEM0
		self.currentItem = self.equippedItems[0]
	case 4:
		self.currentAction = ACTION_ITEM1
		self.currentItem = self.equippedItems[1]
	case 5:
		self.currentAction = ACTION_ITEM2
		self.currentItem = self.equippedItems[2]
	case 6:
		self.currentAction = ACTION_ITEM3
		self.currentItem = self.equippedItems[3]
	case 7:
		self.currentAction = ACTION_ITEM4
		self.currentItem = self.equippedItems[4]
	}
}

func (self *Player) EquipItem(itemid ItemId) {

	// Check to see if this item is already equipped
	for j := 0; j < 5; j++ {
		if self.equippedItems[j] == itemid {
			return
		}
	}

	// Place it in the first empty slot
	for j := 0; j < 5; j++ {
		if self.equippedItems[j] == ITEM_NONE {
			self.equippedItems[j] = itemid
			return
		}
	}

}

func (self *Player) Update(dt float64) (completed bool) {
	self.distanceTravelled += dt * math.Sqrt(math.Pow(self.velocity[XAXIS], 2)+math.Pow(self.velocity[ZAXIS], 2))
	self.MobData.Update(dt)
	return false
}

func (self *Player) TargetType() uint8 {
	return TARGET_PLAYER
}

func (self *Player) DrawWolf(center Vectorf, selectedBlockFace *BlockFace) {
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
