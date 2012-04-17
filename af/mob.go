/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
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
    Setvx(vx float64) 
    Setvy(vy float64) 
    Setvz(vz float64) 
    SetFalling(b bool)
    Accelerate(v Vector)
    Rotate(angle float64)
    Act(dt float64)
    Draw(pos Vector, selectMode bool)
    Update(dt float64)
}
