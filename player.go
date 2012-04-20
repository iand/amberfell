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
	Bounce        float64
	heading       float64
	position      Vectorf
	velocity      Vectorf
	falling       bool
	currentAction uint8
	currentItem   uint16
	walkingSpeed  float64
	equippedItems [7]uint16
	inventory     [255]uint16
}

func (p *Player) Init(heading float64, x int16, z int16, y int16) {
	p.heading = heading
	p.position[XAXIS] = float64(x)
	p.position[YAXIS] = float64(y)
	p.position[ZAXIS] = float64(z)
	p.walkingSpeed = 20
	p.currentAction = ACTION_HAND
	p.currentItem = ITEM_NONE

	p.equippedItems[0] = BLOCK_DIRT
	p.equippedItems[1] = BLOCK_STONE
	p.equippedItems[2] = ITEM_NONE
	p.equippedItems[3] = ITEM_NONE
	p.equippedItems[4] = ITEM_NONE
	p.equippedItems[5] = ITEM_NONE
	p.equippedItems[6] = ITEM_NONE

}

func (p *Player) W() float64 { return 0.8 }
func (p *Player) H() float64 { return 2.0 }
func (p *Player) D() float64 { return 0.6 }

func (p *Player) Heading() float64  { return p.heading }
func (p *Player) X() float32        { return float32(p.position[XAXIS]) }
func (p *Player) Y() float32        { return float32(p.position[YAXIS]) }
func (p *Player) Z() float32        { return float32(p.position[ZAXIS]) }
func (p *Player) Velocity() Vectorf { return p.velocity }
func (p *Player) Position() Vectorf { return p.position }

func (p *Player) FrontBlock() Vectori {
	ip := IntPosition(p.position)
	if p.heading > 337.5 || p.heading <= 22.5 {
		ip[XAXIS]++
	} else if p.heading > 22.5 && p.heading <= 67.5 {
		ip[XAXIS]++
		ip[ZAXIS]--
	} else if p.heading > 67.5 && p.heading <= 112.5 {
		ip[ZAXIS]--
	} else if p.heading > 112.5 && p.heading <= 157.5 {
		ip[XAXIS]--
		ip[ZAXIS]--
	} else if p.heading > 157.5 && p.heading <= 202.5 {
		ip[XAXIS]--
	} else if p.heading > 202.5 && p.heading <= 247.5 {
		ip[XAXIS]--
		ip[ZAXIS]++
	} else if p.heading > 247.5 && p.heading <= 292.5 {
		ip[ZAXIS]++
	} else if p.heading > 292.5 && p.heading <= 337.5 {
		ip[XAXIS]++
		ip[ZAXIS]++
	}

	return ip
}

func (p *Player) SetFalling(b bool) { p.falling = b }

func (p *Player) Rotate(angle float64) {
	p.heading += angle
	if p.heading < 0 {
		p.heading += 360
	}
	if p.heading > 360 {
		p.heading -= 360
	}
}

func (p *Player) Update(dt float64) {
	p.position[XAXIS] += p.velocity[XAXIS] * dt
	p.position[YAXIS] += p.velocity[YAXIS] * dt
	p.position[ZAXIS] += p.velocity[ZAXIS] * dt

	viewport.Transx(-p.velocity[XAXIS] * dt)
	viewport.Transy(-p.velocity[YAXIS] * dt)
	viewport.Transz(-p.velocity[ZAXIS] * dt)

	// fmt.Printf("position: %s\n", p.position)
	// if math.Abs(p.velocity[XAXIS]) < 0.1 { p.velocity[XAXIS] = 0 }
	// if math.Abs(p.velocity[YAXIS]) < 0.1 { p.velocity[YAXIS] = 0 }
	// if math.Abs(p.velocity[ZAXIS]) < 0.1 { p.velocity[ZAXIS] = 0 }

	//if p.velocity[YAXIS] == 0 { p.falling = false }
}

func (p *Player) Accelerate(v Vectorf) {
	p.velocity[XAXIS] += v[XAXIS]
	p.velocity[YAXIS] += v[YAXIS]
	p.velocity[ZAXIS] += v[ZAXIS]
}

func (p *Player) IsFalling() bool {
	return p.falling
}

func (p *Player) Snapx(x float64, vx float64) {
	p.position[XAXIS] = x
	p.velocity[XAXIS] = vx	
}

func (p *Player) Snapz(z float64, vz float64) {
	p.position[ZAXIS] = z
	p.velocity[ZAXIS] = vz
}

func (p *Player) Snapy(y float64, vy float64) {
	p.position[YAXIS] = y
	p.velocity[YAXIS] = vy
}

func (p *Player) Setvx(vx float64) {
	p.velocity[XAXIS] = vx
}

func (p *Player) Setvz(vz float64) {
	p.velocity[ZAXIS] = vz
}

func (p *Player) Setvy(vy float64) {
	p.velocity[YAXIS] = vy
}

func (p *Player) Act(dt float64) {
	// noop
}

func (player *Player) Draw(center Vectorf, selectMode bool) {

	gl.PushMatrix()

	gl.Translatef(float32(player.X()), float32(player.Y()), float32(player.Z()))
	//stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
	gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)

	gl.Translatef(0.0, float32(player.H()/2)-0.5, 0.0)
	Cuboid(player.W(), player.H(), player.D(), &MapTextures[33], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32], &MapTextures[32], 0, selectMode)

	gl.PopMatrix()
}

func (p *Player) HandleKeys(keys []uint8) {
	if keys[sdl.K_1] != 0 {
		p.currentAction = ACTION_HAND
	}
	if keys[sdl.K_2] != 0 {
		p.currentAction = ACTION_BREAK
	}
	if keys[sdl.K_3] != 0 {
		p.currentAction = ACTION_WEAPON
	}
	if keys[sdl.K_4] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[0]
	}
	if keys[sdl.K_5] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[1]
	}
	if keys[sdl.K_6] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[2]
	}
	if keys[sdl.K_7] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[3]
	}
	if keys[sdl.K_8] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[4]
	}
	if keys[sdl.K_9] != 0 {
		p.currentAction = ACTION_ITEM
		p.currentItem = p.equippedItems[5]
	}

	if keys[sdl.K_w] != 0 {
		if !p.IsFalling() {
			p.velocity[XAXIS] = math.Cos(p.Heading()*math.Pi/180) * p.walkingSpeed
			p.velocity[ZAXIS] = -math.Sin(p.Heading()*math.Pi/180) * p.walkingSpeed
		} else {
			p.velocity[XAXIS] = math.Cos(p.Heading()*math.Pi/180) * p.walkingSpeed / 3
			p.velocity[ZAXIS] = -math.Sin(p.Heading()*math.Pi/180) * p.walkingSpeed / 3
		}

	}
	if keys[sdl.K_s] != 0 {
		if !p.IsFalling() {
			p.velocity[XAXIS] = -math.Cos(p.Heading()*math.Pi/180) * p.walkingSpeed / 2
			p.velocity[ZAXIS] = math.Sin(p.Heading()*math.Pi/180) * p.walkingSpeed / 2
		} else {
			p.velocity[XAXIS] = -math.Cos(p.Heading()*math.Pi/180) * p.walkingSpeed / 6
			p.velocity[ZAXIS] = math.Sin(p.Heading()*math.Pi/180) * p.walkingSpeed / 6
		}

	}
	if keys[sdl.K_a] != 0 {
		p.Rotate(22.5 / 2)
		viewport.Roty(-22.5 / 2)
	}

	if keys[sdl.K_d] != 0 {
		p.Rotate(-22.5 / 2)
		viewport.Roty(22.5 / 2)
	}

	if keys[sdl.K_SPACE] != 0 {
		if !p.IsFalling() {
			p.velocity[YAXIS] = 4
		}
	}

}

func (p *Player) CanInteract() bool {
	if p.currentAction == ACTION_BREAK || (p.currentAction == ACTION_ITEM && p.currentItem != ITEM_NONE) {
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
