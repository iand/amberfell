package af

type Mob interface {
    Heading() float64
    W() float64
    H() float64
    D() float64
    X() float32
    Y() float32
    Z() float32
    Velocity() Vector
    Forces() Vector
    Mass() float64
    ApplyForce(f Vector)
    Reaction(f Vector)
    Rotate(angle float64)
    Speed() float64
    BoundingBox() Bound
    DesiredBoundingBox(dt float64) Bound
}
