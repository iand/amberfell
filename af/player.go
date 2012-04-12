package af

import (    
    "math"
   // "fmt"
)

const (
    XAXIS = 0
    YAXIS = 1
    ZAXIS = 2
)

type Player struct {
    Bounce float64
    heading float64
    position Vector
    velocity Vector
    falling bool
}

func (p *Player) Init(heading float64, x float32, z float32, y float32) {
    p.heading = heading
    p.position[XAXIS] = float64(x)
    p.position[YAXIS] = float64(y)
    p.position[ZAXIS] = float64(z)
}

func (p *Player) W() float64 { return 0.6 }
func (p *Player) H() float64 { return 1.9 }
func (p *Player) D() float64 { return 0.8 }

func (p *Player) Heading() float64 { return p.heading }
func (p *Player) X() float32 { return float32(p.position[XAXIS]) }
func (p *Player) Y() float32 { return float32(p.position[YAXIS]) }
func (p *Player) Z() float32 { return float32(p.position[ZAXIS]) }
func (p *Player) Velocity() Vector { return p.velocity }
func (p *Player) Position() Vector { return p.position }
func (p *Player) IntPosition() IntVector { 
    return IntVector{ int16(math.Floor(p.position[XAXIS] + 0.5)),
                      int16(math.Floor(p.position[YAXIS] + 0.5)),
                      int16(math.Floor(p.position[ZAXIS] + 0.5))}
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




// func (p *Player) BoundingBox() Bound {
//     var b Bound


//     // extentx := math.Abs(p.W() * math.Cos(p.Heading() * math.Pi / 180) / 2) + math.Abs(p.D() * math.Sin(p.Heading() * math.Pi / 180) / 2)
//     // extentz := math.Abs(p.W() * math.Sin(p.Heading() * math.Pi / 180) / 2) + math.Abs(p.D() * math.Cos(p.Heading() * math.Pi / 180) / 2)

//     //b.extent = Vector{extentx, p.H() / 2, extentz}
//     b.extent = Vector{p.W(), p.H() / 2, p.D()}
//     b.position = Vector{p.position[XAXIS], p.position[YAXIS], p.position[ZAXIS]}
//     normalx := Vector{math.Cos(p.heading * math.Pi / 180), 0, -math.Sin(p.heading * math.Pi / 180)}
//     normaly := Vector{0,1,0}
//     normalz := Vector{math.Sin(p.heading * math.Pi / 180), 0, -math.Cos(p.heading * math.Pi / 180)}
//     b.orthonormal = [3]Vector{normalx, normaly, normalz}
//     return b
// }

// func (p *Player) DesiredBoundingBox(dt float64) Bound {
//     var b Bound
//     b.extent = Vector{p.W() / 2, p.H() / 2, p.D() / 2}
//     b.position = Vector{p.position[XAXIS]+p.velocity[XAXIS]*dt, p.position[YAXIS]+p.velocity[YAXIS]*dt, p.position[ZAXIS]+p.velocity[ZAXIS]*dt}
//     normalx := Vector{math.Sin(p.heading * math.Pi / 180), 0, math.Cos(p.heading * math.Pi / 180)}
//     normaly := Vector{0,1,0}
//     normalz := Vector{math.Cos(p.heading * math.Pi / 180), 0, math.Sin(p.heading * math.Pi / 180)}
//     b.orthonormal = [3]Vector{normalx, normaly, normalz}
//     return b
// }


// func (p *Player) Reach() Vector {
//     ip := p.IntPosition()
//     if p.heading >= 360 - 22.5 && p.heading < 22.5 {
//         var x, y, z int16
//         for x = -2; x <= 2; x++ {
//             for z = -2; z <= 2; z++ {
//                 for y = -2; z <= 2; z++ {

//                 }
//             }
//         }


//         IntPosition{ip[XAXIS] + 1, ip[YAXIS]    , ip[ZAXIS]}
//         IntPosition{ip[XAXIS] + 2, ip[YAXIS]    , ip[ZAXIS]}

//         IntPosition{ip[XAXIS]    , ip[YAXIS]    , ip[ZAXIS] + 1}
//         IntPosition{ip[XAXIS]    , ip[YAXIS]    , ip[ZAXIS] + 2}
//         IntPosition{ip[XAXIS] + 1, ip[YAXIS]    , ip[ZAXIS] + 1}

//         IntPosition{ip[XAXIS]    , ip[YAXIS]    , ip[ZAXIS] - 1}
//         IntPosition{ip[XAXIS]    , ip[YAXIS]    , ip[ZAXIS] - 2}
//         IntPosition{ip[XAXIS] + 1, ip[YAXIS]    , ip[ZAXIS] - 1}

//         IntPosition{ip[XAXIS]    , ip[YAXIS] + 1, ip[ZAXIS] - 1}
//         IntPosition{ip[XAXIS]    , ip[YAXIS] + 2, ip[ZAXIS] - 1}
//         IntPosition{ip[XAXIS]    , ip[YAXIS] - 1, ip[ZAXIS] - 1}
//         IntPosition{ip[XAXIS]    , ip[YAXIS] - 2, ip[ZAXIS] - 1}



//         IntPosition{ip[XAXIS]    , ip[YAXIS] + 1, ip[ZAXIS]}
//         IntPosition{ip[XAXIS]    , ip[YAXIS] - 1, ip[ZAXIS]}
//         IntPosition{ip[XAXIS] + 1, ip[YAXIS] + 1, ip[ZAXIS]}
//         IntPosition{ip[XAXIS] + 1, ip[YAXIS] - 1, ip[ZAXIS]}
//     }
// }


