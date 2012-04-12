package af
import (
    "fmt"
    "math"
)

type IntVector [3]int16
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