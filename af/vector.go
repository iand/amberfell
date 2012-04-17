/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package af
import (
    "fmt"
    "math"
)

type Vector [3]float64
func (a Vector) Minus(b Vector) Vector {
    return Vector{a[0]-b[0], a[1]-b[1], a[2]-b[2]}
}

func (a Vector) Dot(b Vector) float64 {
    return a[0] * b[0] + a[1] * b[1] + a[2] * b[2]
}

func (a Vector) Scale(f float64) Vector {
    return Vector{a[0] * f, a[1] * f, a[2] * f}
}

func (a Vector) Magnitude() float64 {
    return math.Sqrt(math.Pow(a[0], 2) + math.Pow(a[1], 2)  + math.Pow(a[2], 2) )
}


func (a Vector) String() string {
    return fmt.Sprintf("Vector(%f,%f,%f)", a[0], a[1], a[2]) 
}


type IntVector [3]int16

type Rect struct {
    x, z, sizex, sizez float64
}

func (wp IntVector) North() IntVector { return IntVector{wp[XAXIS], wp[YAXIS], wp[ZAXIS]-1} }
func (wp IntVector) South() IntVector { return IntVector{wp[XAXIS], wp[YAXIS], wp[ZAXIS]+1} }
func (wp IntVector) East()  IntVector { return IntVector{wp[XAXIS]+1, wp[YAXIS], wp[ZAXIS]} }
func (wp IntVector) West()  IntVector { return IntVector{wp[XAXIS]-1, wp[YAXIS], wp[ZAXIS]} }
func (wp IntVector) Up()  IntVector   { return IntVector{wp[XAXIS], wp[YAXIS]+1, wp[ZAXIS]} }
func (wp IntVector) Down()  IntVector { return IntVector{wp[XAXIS], wp[YAXIS]-1, wp[ZAXIS]} }

func (wp IntVector) HRect() Rect { return Rect{float64(wp[XAXIS]) - 0.5, float64(wp[ZAXIS]) - 0.5, 1, 1} }

func (r1 Rect) Intersects(r2 Rect) bool {
    if r2.x >= r1.x && r2.x <= r1.x + r1.sizex && r2.z >= r1.z && r2.z <= r1.z + r1.sizez {
        return true
    }
    if r2.x + r2.sizex >= r1.x && r2.x + r2.sizex <= r1.x + r1.sizex && r2.z >= r1.z && r2.z <= r1.z + r1.sizez { 
        return true
    }
    if r2.x + r2.sizex >= r1.x && r2.x + r2.sizex <= r1.x + r1.sizex && r2.z + r2.sizez >= r1.z && r2.z + r2.sizez <= r1.z + r1.sizez {
        return true
    }
    if r2.x >= r1.x && r2.x <= r1.x + r1.sizex && r2.z + r2.sizez >= r1.z && r2.z + r2.sizez <= r1.z + r1.sizez {
        return true
    }
    return false

}


