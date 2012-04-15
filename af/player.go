package af

import (    
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/gl"
    "math"
   // "fmt"
)



type Player struct {
    Bounce float64
    heading float64
    position Vector
    velocity Vector
    falling bool
    currentTool uint16
    walkingSpeed float64
}

func (p *Player) Init(heading float64, x float32, z float32, y float32) {
    p.heading = heading
    p.position[XAXIS] = float64(x)
    p.position[YAXIS] = float64(y)
    p.position[ZAXIS] = float64(z)
    p.walkingSpeed = 12
}

func (p *Player) W() float64 { return 0.8 }
func (p *Player) H() float64 { return 2.0 }
func (p *Player) D() float64 { return 0.6 }

func (p *Player) Heading() float64 { return p.heading }
func (p *Player) X() float32 { return float32(p.position[XAXIS]) }
func (p *Player) Y() float32 { return float32(p.position[YAXIS]) }
func (p *Player) Z() float32 { return float32(p.position[ZAXIS]) }
func (p *Player) Velocity() Vector { return p.velocity }
func (p *Player) Position() Vector { return p.position }

func (p *Player) FrontBlock() IntVector { 
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
    // fmt.Printf("position: %s\n", p.position)
    // if math.Abs(p.velocity[XAXIS]) < 0.1 { p.velocity[XAXIS] = 0 }
    // if math.Abs(p.velocity[YAXIS]) < 0.1 { p.velocity[YAXIS] = 0 }
    // if math.Abs(p.velocity[ZAXIS]) < 0.1 { p.velocity[ZAXIS] = 0 }

    //if p.velocity[YAXIS] == 0 { p.falling = false }
}

func (p *Player) Accelerate(v Vector) {
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


func (player *Player) Draw(center Vector, selectMode bool) {

    gl.PushMatrix()

    //stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
    gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)

    gl.Translatef(0.0, float32(player.H() / 2) - 0.5 ,0.0)
    Cuboid(player.W(), player.H(), player.D(), 33, 32, 32, 32, 32, 32, 0, selectMode)

    gl.PopMatrix()
}

func (p *Player) HandleKeys(keys []uint8) {
    if keys[sdl.K_1] != 0 {
        p.currentTool = TOOL_HAND
    }
    if keys[sdl.K_2] != 0 {
        p.currentTool = TOOL_DIG
    }
    if keys[sdl.K_3] != 0 {
        p.currentTool = BLOCK_DIRT
    }


    if keys[sdl.K_w] != 0 {
        if !p.IsFalling() {
            p.velocity[XAXIS] = math.Cos(p.Heading() * math.Pi / 180) * p.walkingSpeed
            p.velocity[ZAXIS] = -math.Sin(p.Heading() * math.Pi / 180) * p.walkingSpeed
        } else {
            p.velocity[XAXIS] = math.Cos(p.Heading() * math.Pi / 180) * p.walkingSpeed / 3
            p.velocity[ZAXIS] = -math.Sin(p.Heading() * math.Pi / 180) * p.walkingSpeed / 3
        }

    }
    if keys[sdl.K_s] != 0 {
        if !p.IsFalling() {
            p.velocity[XAXIS] = -math.Cos(p.Heading() * math.Pi / 180) * p.walkingSpeed
            p.velocity[ZAXIS] = math.Sin(p.Heading() * math.Pi / 180) * p.walkingSpeed
        } else {
            p.velocity[XAXIS] = -math.Cos(p.Heading() * math.Pi / 180) * p.walkingSpeed / 3
            p.velocity[ZAXIS] = math.Sin(p.Heading() * math.Pi / 180) * p.walkingSpeed / 3
        }
 
    }
    if keys[sdl.K_a] != 0 {
        p.Rotate(9)

    }    

    if keys[sdl.K_d] != 0 {
        p.Rotate(-9)
    }        


    if keys[sdl.K_SPACE] != 0 {
        if !p.IsFalling() {
            println("jump")
            p.velocity[YAXIS] += 3
        }
    } 

}

    