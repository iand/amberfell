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
	walkingSpeed  float64
	equippedItems [7]uint16
	inventory     [255]uint16
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

func (player *Player) Draw(center Vectorf) {

	gl.PushMatrix()

	gl.Translatef(float32(player.position[XAXIS]), float32(player.position[YAXIS]), float32(player.position[ZAXIS]))
	//stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
	gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)

	gl.Translatef(0.0, float32(player.H()/2)-0.5, 0.0)
	Cuboid(player.W(), player.H(), player.D(), &MapTextures[33], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32])

	gl.PopMatrix()
}

func (self *Player) HandleKeys(keys []uint8) {
	if keys[sdl.K_1] != 0 {
		self.currentAction = ACTION_HAND
	}
	if keys[sdl.K_2] != 0 {
		self.currentAction = ACTION_BREAK
	}
	if keys[sdl.K_3] != 0 {
		self.currentAction = ACTION_WEAPON
	}
	if keys[sdl.K_4] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[0]
	}
	if keys[sdl.K_5] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[1]
	}
	if keys[sdl.K_6] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[2]
	}
	if keys[sdl.K_7] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[3]
	}
	if keys[sdl.K_8] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[4]
	}
	if keys[sdl.K_9] != 0 {
		self.currentAction = ACTION_ITEM
		self.currentItem = self.equippedItems[5]
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
	if self.currentAction == ACTION_BREAK || (self.currentAction == ACTION_ITEM && self.currentItem != ITEM_NONE) {
		return true
	}
	return false
}

func (self *Player) Interact(pos Vectori, face uint8) {
	if !self.CanInteract() {
		return
	}

	println("Interacting at ", pos.String())
	switch self.currentAction {
	case ACTION_BREAK:
		block := TheWorld.Atv(pos)
		if block != BLOCK_AIR {
			println("Breaking at ", pos.String())
			TheWorld.Setv(pos, BLOCK_AIR)
			self.inventory[block]++
		}
	case ACTION_ITEM:
		if face == UP_FACE { // top
			pos[YAXIS]++
		} else if face == DOWN_FACE { // bottom
			pos[YAXIS]--
		} else if face == SOUTH_FACE { // front
			pos[ZAXIS]++
		} else if face == NORTH_FACE { // back
			pos[ZAXIS]--
		} else if face == EAST_FACE { // left
			pos[XAXIS]++
		} else if face == WEST_FACE { // right
			pos[XAXIS]--
		}
		if TheWorld.Atv(pos) == BLOCK_AIR {
			println("Setting at ", pos.String())
			TheWorld.Setv(pos, byte(self.currentItem))
		}
	}

}
