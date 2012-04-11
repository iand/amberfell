package af

type Mob interface {
    Heading() float64
    W() float64
    H() float64
    D() float64
    X() float32
    Y() float32
    Z() float32
    IsFalling() bool
    Velocity() Vector
    Position() Vector
    Snapx(x float64, vx float64)
    Snapy(y float64, vy float64)
    Snapz(z float64, vz float64)
    SetFalling(b bool)
    // Forces() Vector
    // Mass() float64
    Accelerate(v Vector)
    // ApplyForce(f Vector)
    // Reaction(f Vector)
    Rotate(angle float64)
    // Speed() float64
    // BoundingBox() Bound
    // DesiredBoundingBox(dt float64) Bound
}
