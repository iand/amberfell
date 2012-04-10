package af

import (    
    "math"
   "fmt"
)

const (
    XAXIS = 0
    YAXIS = 1
    ZAXIS = 2
)

type State struct {
    heading float64
    position Vector
    momentum Vector
    velocity Vector
    mass float64
}

func (s *State) Recalculate() {
    s.velocity = s.momentum.Scale(1 / s.mass)
    if math.Abs(s.velocity[0]) < 0.2 {
        s.velocity[0] = 0
    }
    if math.Abs(s.velocity[1]) < 0.2 {
        s.velocity[1] = 0
    }
    if math.Abs(s.velocity[2]) < 0.2 {
        s.velocity[2] = 0
    }

    //s.velocity[ZAXIS] = 5
}

func (s *State) Update(forces Vector, dt float64) {
    fmt.Printf("forces: %s\n", forces)

    s.position[XAXIS] += s.velocity[XAXIS] * dt
    s.position[YAXIS] += s.velocity[YAXIS] * dt
    s.position[ZAXIS] += s.velocity[ZAXIS] * dt

    s.momentum[XAXIS] += forces[XAXIS] * dt
    s.momentum[YAXIS] += forces[YAXIS] * dt
    s.momentum[ZAXIS] += forces[ZAXIS] * dt

    s.Recalculate()
}

func (s *State) Rotate(angle float64) {
    s.heading += angle
    if s.heading < 0 {
        s.heading += 360
    }
    if s.heading > 360 {
        s.heading -= 360
    }
}


type Player struct {
    Bounce float64
    forces Vector

    current, previous State
}

func (p *Player) Init(heading float64, x float32, z float32, y float32) {
    p.current.heading = heading
    p.current.position[XAXIS] = float64(x)
    p.current.position[YAXIS] = float64(y)
    p.current.position[ZAXIS] = float64(z)
    p.current.mass = 200

    p.previous = p.current
}

func (p *Player) W() float64 { return 0.5 }
func (p *Player) H() float64 { return 1.0 }
func (p *Player) D() float64 { return 0.7 }

func (p *Player) Heading() float64 { return p.current.heading }
func (p *Player) X() float32 { return float32(p.current.position[XAXIS]) }
func (p *Player) Y() float32 { return float32(p.current.position[YAXIS]) }
func (p *Player) Z() float32 { return float32(p.current.position[ZAXIS]) }
func (p *Player) Velocity() Vector { return p.current.velocity }
func (p *Player) Forces() Vector { return p.forces }
func (p *Player) Mass() float64 { return p.current.mass }

func (p *Player) Speed() float64 { return p.current.velocity.Magnitude() }

func (p *Player) Rotate(angle float64) {
    p.current.Rotate(angle)
}

func (p *Player) Snap(state State) {
    p.current = state
    p.previous = state
}


func (p *Player) BoundingBox() Bound {
    var b Bound


    // extentx := math.Abs(p.W() * math.Cos(p.Heading() * math.Pi / 180) / 2) + math.Abs(p.D() * math.Sin(p.Heading() * math.Pi / 180) / 2)
    // extentz := math.Abs(p.W() * math.Sin(p.Heading() * math.Pi / 180) / 2) + math.Abs(p.D() * math.Cos(p.Heading() * math.Pi / 180) / 2)

    //b.extent = Vector{extentx, p.H() / 2, extentz}
    b.extent = Vector{p.W(), p.H() / 2, p.D()}
    b.position = Vector{p.current.position[XAXIS], p.current.position[YAXIS], p.current.position[ZAXIS]}
    normalx := Vector{math.Cos(p.current.heading * math.Pi / 180), 0, -math.Sin(p.current.heading * math.Pi / 180)}
    normaly := Vector{0,1,0}
    normalz := Vector{math.Sin(p.current.heading * math.Pi / 180), 0, -math.Cos(p.current.heading * math.Pi / 180)}
    b.orthonormal = [3]Vector{normalx, normaly, normalz}
    return b
}

func (p *Player) DesiredBoundingBox(dt float64) Bound {
    var b Bound
    b.extent = Vector{p.W() / 2, p.H() / 2, p.D() / 2}
    b.position = Vector{p.current.position[XAXIS]+p.current.velocity[XAXIS]*dt, p.current.position[YAXIS]+p.current.velocity[YAXIS]*dt, p.current.position[ZAXIS]+p.current.velocity[ZAXIS]*dt}
    normalx := Vector{math.Sin(p.current.heading * math.Pi / 180), 0, math.Cos(p.current.heading * math.Pi / 180)}
    normaly := Vector{0,1,0}
    normalz := Vector{math.Cos(p.current.heading * math.Pi / 180), 0, math.Sin(p.current.heading * math.Pi / 180)}
    b.orthonormal = [3]Vector{normalx, normaly, normalz}
    return b
}


func (p *Player) ApplyForce(f Vector) {
    //fmt.Printf("Force: %s\n", f)
    p.forces[XAXIS] += f[XAXIS]
    p.forces[YAXIS] += f[YAXIS]
    p.forces[ZAXIS] += f[ZAXIS]
}

func (p *Player) Reaction(n Vector) {
    fmt.Printf("normal: %s\n", n)

    if n[XAXIS] > 0 {
        p.forces[XAXIS] = math.Min(p.forces[XAXIS], 0)
        p.current.velocity[XAXIS] = math.Min(p.current.velocity[XAXIS], 0)
        p.current.momentum[XAXIS] = math.Min(p.current.momentum[XAXIS], 0)
    } else if n[XAXIS] < 0 {
        p.forces[XAXIS] = -p.forces[XAXIS]
        p.current.velocity[XAXIS] = math.Max(p.current.velocity[XAXIS], 0)
        p.current.momentum[XAXIS] = math.Max(p.current.momentum[XAXIS], 0)
    }
    if n[YAXIS] > 0 {
        p.forces[YAXIS] = math.Min(p.forces[YAXIS], 0)
        p.current.velocity[YAXIS] = math.Min(p.current.velocity[YAXIS], 0)
        p.current.momentum[YAXIS] = math.Min(p.current.momentum[YAXIS], 0)
    } else if n[YAXIS] < 0 {
        p.forces[YAXIS] = math.Max(p.forces[YAXIS], 0)
        p.current.velocity[YAXIS] = math.Max(p.current.velocity[YAXIS], 0)
        p.current.momentum[YAXIS] = math.Max(p.current.momentum[YAXIS], 0)
    }
    if n[ZAXIS] > 0 {
        p.forces[ZAXIS] = math.Min(p.forces[ZAXIS], 0)
        p.current.velocity[ZAXIS] = math.Min(p.current.velocity[ZAXIS], 0)
        p.current.momentum[ZAXIS] = math.Min(p.current.momentum[ZAXIS], 0)
    } else if n[ZAXIS] < 0 {
        p.forces[ZAXIS] = math.Max(p.forces[ZAXIS], 0)
        p.current.velocity[ZAXIS] = math.Max(p.current.velocity[ZAXIS], 0)
        p.current.momentum[ZAXIS] = math.Max(p.current.momentum[ZAXIS], 0)
    }
}


func (p *Player) ZeroForces() {
    p.forces[XAXIS] = 0
    p.forces[YAXIS] = 0
    p.forces[ZAXIS] = 0
    //p.velocity[XAXIS] = 0
    //p.velocity[YAXIS] = 0
    //p.velocity[ZAXIS] = 0
}


func (p *Player) Update(dt float64) {
    p.previous = p.current
    p.current.Update(p.forces, dt)
}

