package af

import (
    "math"
    "math/rand"
    "fmt"

)


type World struct {
    W           int16
    D           int16
    H           int16
    GroundLevel int16
    Blocks      []byte
}

func (world *World) Init(w int16, d int16, h int16) {
    world.W = w
    world.D = d
    world.H = h
    world.Blocks = make([]byte, w*d*h)
    var iw, id, ih int16
    for iw = 0; iw < w; iw++ {
        for id = 0; id < d; id++ {
            for ih = 0; ih <= GroundLevel; ih++ {
                world.Set(iw, id, ih, 2) // dirt
            }
            for ih = GroundLevel + 1; ih < h; ih++ {
                world.Set(iw, id, ih, 0) // air
            }
        }
    }

    numFeatures := rand.Intn(int(world.W))
    for i := 0; i < numFeatures; i++ {
        iw, id = world.RandomSquare()

        world.Set(iw, id, GroundLevel, 1) // stone
        world.Grow(iw, id, GroundLevel, 40, 40, 40, 40, 10, 30, 1)
    }
    iw, id = world.RandomSquare()

    world.Set(iw, id, GroundLevel, 0) // air
    world.Grow(iw, id, GroundLevel, 30, 30, 30, 30, 0, 30, 0)

}

func (world *World) At(x int16, y int16, z int16) byte {
    if x < 0 || x > world.W-1 || z < 0 || z > world.D-1 || y < 0 || y > world.H-1 {
        return 0
    }
    return world.Blocks[world.W*world.D*y+world.D*x+z]
}

func (world *World) Set(x int16, z int16, y int16, b byte) {
    world.Blocks[world.W*world.D*y+world.D*x+z] = b
}

func (world *World) RandomSquare() (x int16, z int16) {
    x = int16(rand.Intn(int(world.W)))
    z = int16(rand.Intn(int(world.D)))
    return
}

// north/south = -/+ z
// east/west = +/- x
// up/down = +/- y

func (world *World) Grow(x int16, z int16, y int16, n int, s int, w int, e int, u int, d int, texture byte) {
    if x < world.W-1 && world.At(x+1, y-1, z) != 0 && rand.Intn(100) < e {
        world.Set(x+1, z, y, texture)
        world.Grow(x+1, z, y, n, s, 0, e, u, d, texture)
    }
    if x > 0 && world.At(x-1, y-1, z) != 0 && rand.Intn(100) < w {
        world.Set(x-1, z, y, texture)
        world.Grow(x-1, z, y, n, s, w, 0, u, d, texture)
    }
    if y < world.D-1 && world.At(x, y-1, z+1) != 0 && rand.Intn(100) < s {
        world.Set(x, z+1, y, texture)
        world.Grow(x, z+1, y, n, 0, w, e, u, d, texture)
    }
    if y > 0 && world.At(x, y-1, z-1) != 0 && rand.Intn(100) < n {
        world.Set(x, z-1, y, texture)
        world.Grow(x, z-1, y, 0, s, w, e, u, d, texture)
    }
    if y < world.H-1 && rand.Intn(100) < u {
        world.Set(x, z, y+1, texture)
        world.Grow(x, z, y+1, n, s, w, e, u, 0, texture)
    }
    if y > 0 && rand.Intn(100) < d {
        world.Set(x, z, y-1, texture)
        world.Grow(x, z, y-1, n, s, w, e, 0, d, texture)
    }
}

func (world *World) AirNeighbours(x int16, z int16, y int16) (n, s, w, e, u, d bool) {
    if x > 0 && world.At(x-1, y, z) == 0 {
        e = true
    }
    if x < world.W-1 && world.At(x+1, y, z) == 0 {
        w = true
    }
    if z > 0 && world.At(x, y, z-1) == 0 {
        s = true
    }
    if z < world.D-1 && world.At(x, y, z+1) == 0 {
        n = true
    }
    if y < world.H-1 && world.At(x, y+1, z) == 0 {
        u = true
    }
    return
}

func (world *World) ApplyForces(mob Mob, dt float64) {
    mobBounds := mob.DesiredBoundingBox(dt)
    mobx := int16(math.Floor(mobBounds.position[0] + 0.5))
    moby := int16(math.Floor(mobBounds.position[1] + 0.5))
    mobz := int16(math.Floor(mobBounds.position[2] + 0.5))
 

    // Gravity
    mob.ApplyForce( Vector{0, -400, 0} )

    const frictionConstant = -5.5
    friction := mob.Velocity().Scale(mob.Mass() * frictionConstant)
    mob.ApplyForce( Vector{friction[0], 0, friction[2]} )

    const c = -600.0
    const k = -800.0
    const b = -500.0   
    const f = 3.0

    var dx, dz, dy int16
    for dx = -1; dx < 2; dx++ {
        for dz = -1; dz < 2; dz++ {
            for dy = -1; dy < 2; dy++ {
                //if dy != 0 && dx != 00 && dz != 0 {
                    if world.At(mobx+dx, moby+dy, mobz+dz) != 0 {
                        block := BlockBound(mobx+dx, moby+dy, mobz+dz)

                        // normal := Vector{ -float64(dx), -float64(dy), -float64(dz)  }
                        // //separation := block.Distance(mobBounds) - normal.Dot(mobBounds.extent) - normal.Dot(block.extent)

                        // sepx := math.Abs(mobBounds.position[XAXIS] - block.position[XAXIS]) - mobBounds.extent[XAXIS] - block.extent[XAXIS]
                        // sepy := math.Abs(mobBounds.position[YAXIS] - block.position[YAXIS]) - mobBounds.extent[YAXIS] - block.extent[YAXIS]
                        // sepz := math.Abs(mobBounds.position[ZAXIS] - block.position[ZAXIS]) - mobBounds.extent[ZAXIS] - block.extent[ZAXIS]

                        // if sepy < 0 {
                        //     if sepx < 0 {
                        //         mob.Reaction(Vector{float64(dx), 0, 0})
                        //     }
                        //     if sepz < 0 {
                        //         mob.Reaction(Vector{0, 0, float64(dz)})
                        //     }

                        //     fmt.Printf("normal: %s\n", normal)
                        //     fmt.Printf("sepx: %f\n", sepx)
                        //     fmt.Printf("sepy: %f\n", sepy)
                        //     fmt.Printf("sepz: %f\n", sepz)
                                
                        // }

                        if block.Overlaps(mobBounds) {
                            //fmt.Printf("block.y+0.5: %f, mob.y-0.5:%f\n", float64(moby+dy), mobBounds.position[1])

                            // penetration := block.Distance(mobBounds) // an approximation
                            normal := Vector{ float64(dx), float64(dy), float64(dz)  }
                            fmt.Printf("normal: %s\n", normal)
                            // relativeSpeed := normal.Dot(mob.Velocity())
                            // /reaction := normal.Dot(mob.Forces())

                            //reactionForce := normal.Scale(reaction * 5);
                            //fmt.Printf("reactionForce: %f\n", reactionForce)
                            //mob.ApplyForce(reactionForce)
                            mob.Reaction(normal)
                            
                            // fmt.Printf("relativeSpeed: %f\n", relativeSpeed)
                            // fmt.Printf("penetration: %f\n", penetration)

                            // if relativeSpeed > 0 {
                            //     collisionForce := normal.Scale(relativeSpeed * c);
                            //     fmt.Printf("collisionForce: %f\n", collisionForce)
                            //      mob.ApplyForce(collisionForce)
                            // }

                            // penaltyForce := normal.Scale(penetration * k);
                            // fmt.Printf("penaltyForce: %f\n", penaltyForce)
                            // mob.ApplyForce(penaltyForce)

                            // dampingForce := normal.Scale(relativeSpeed * penetration * b);
                            // fmt.Printf("dampingForce: %f\n", dampingForce)
                            // mob.ApplyForce(dampingForce)
                        }
                    }
                //}
            }
        }
    }
}




type Bound struct {
    extent Vector
    position Vector
    orthonormal [3]Vector
}

func (a Bound) Overlaps (b Bound) bool {
    return Overlaps(a.extent, a.position, a.orthonormal, b.extent, b.position, b.orthonormal)
}

func (a Bound) Distance(b Bound) float64 {
    return math.Sqrt(math.Pow(a.position[0]-b.position[0], 2) + math.Pow(a.position[1]-b.position[1], 2) + math.Pow(a.position[2]-b.position[2], 2))
}


func BlockBound(x, y, z int16) Bound {
    var b Bound
    b.extent = Vector{0.5, 0.5, 0.5}
    b.position = Vector{float64(x), float64(y), float64(z)}
    normalx := Vector{1,0,0}
    normaly := Vector{0,1,0}
    normalz := Vector{0,0,1}
    b.orthonormal = [3]Vector{normalx, normaly, normalz}
    return b
}



// See http://www.gamasutra.com/view/feature/131790/simple_intersection_tests_for_games.php?print=1
// and http://www.geometrictools.com/Documentation/DynamicCollisionDetection.pdf
// a - extents of a
// pa - position of a
// A - orthonormal basis of a
// b - extents of b
// pb - position of b
// B - orthonormal basis of b
func Overlaps(a Vector, pa Vector, A [3]Vector, b Vector, pb Vector, B [3]Vector) bool {
    //translation, in parent frame
    v := pb.Minus(pa) 

    // all calculations are now done in a's frame

    //translation, in A's frame
    T := Vector{v.Dot(A[0]), v.Dot(A[1]), v.Dot(A[2])}

    var ra, rb, t float64

    //B's basis with respect to A's local frame
    var R [3][3]float64

    //calculate rotation matrix
    for i:=0 ; i<3 ; i++ {
        for k:=0 ; k<3 ; k++ {
            R[i][k] = A[i].Dot(B[k]) 
        }
    }


    // In the following, t is the separation between the centres of each box
    // ra + rb - t gives the penetration depth

    /*ALGORITHM: Use the separating axis test for all 15 potential
    separating axes. If a separating axis could not be found, the two
    boxes overlap. */

    //A's basis vectors
    for i:=0 ; i<3 ; i++ {
        ra = a[i]
        rb = b[0]*math.Abs(R[i][0]) + b[1]*math.Abs(R[i][1]) + b[2]*math.Abs(R[i][2])
        t = math.Abs( T[i] )
        if t > ra + rb {
            return false
        }
    } 

    //B's basis vectors
    for k:=0 ; k<3 ; k++ {
        ra = a[0]*math.Abs(R[0][k]) + a[1]*math.Abs(R[1][k]) + a[2]*math.Abs(R[2][k])
        rb = b[k]

        t = math.Abs( T[0]*R[0][k] + T[1]*R[1][k] + T[2]*R[2][k] )

        if t > ra + rb {
            return false
        }
    } 

    //9 cross products

    //L = A0 x B0
    ra = a[1]*math.Abs(R[2][0]) + a[2]*math.Abs(R[1][0])
    rb = b[1]*math.Abs(R[0][2]) + b[2]*math.Abs(R[0][1])
    t = math.Abs( T[2]*R[1][0] - T[1]*R[2][0] )
    if t >= ra + rb {
        return false
    }


    //L = A0 x B1
    ra = a[1]*math.Abs(R[2][1]) + a[2]*math.Abs(R[1][1])
    rb = b[0]*math.Abs(R[0][2]) + b[2]*math.Abs(R[0][0])
    t = math.Abs( T[2]*R[1][1] - T[1]*R[2][1] )
    if t > ra + rb {
        return false
    }

    //L = A0 x B2
    ra = a[1]*math.Abs(R[2][2]) + a[2]*math.Abs(R[1][2])
    rb = b[0]*math.Abs(R[0][1]) + b[1]*math.Abs(R[0][0])
    t = math.Abs( T[2]*R[1][2] - T[1]*R[2][2] )
    if t > ra + rb {
        return false
    }

    //L = A1 x B0
    ra = a[0]*math.Abs(R[2][0]) + a[2]*math.Abs(R[0][0])
    rb = b[1]*math.Abs(R[1][2]) + b[2]*math.Abs(R[1][1])
    t = math.Abs( T[0]*R[2][0] - T[2]*R[0][0] )
    if t > ra + rb {
        return false
    }

    //L = A1 x B1
    ra = a[0]*math.Abs(R[2][1]) + a[2]*math.Abs(R[0][1])
    rb = b[0]*math.Abs(R[1][2]) + b[2]*math.Abs(R[1][0])
    t = math.Abs( T[0]*R[2][1] - T[2]*R[0][1] )
    if t > ra + rb {
        return false
    }

    //L = A1 x B2
    ra = a[0]*math.Abs(R[2][2]) + a[2]*math.Abs(R[0][2])
    rb = b[0]*math.Abs(R[1][1]) + b[1]*math.Abs(R[1][0])
    t = math.Abs( T[0]*R[2][2] - T[2]*R[0][2] )
    if t > ra + rb {
        return false
    }

    //L = A2 x B0
    ra = a[0]*math.Abs(R[1][0]) + a[1]*math.Abs(R[0][0])
    rb = b[1]*math.Abs(R[2][2]) + b[2]*math.Abs(R[2][1])
    t = math.Abs( T[1]*R[0][0] - T[0]*R[1][0] )
    if t > ra + rb {
        return false
    }

    //L = A2 x B1
    ra = a[0]*math.Abs(R[1][1]) + a[1]*math.Abs(R[0][1])
    rb = b[0]*math.Abs(R[2][2]) + b[2]*math.Abs(R[2][0])
    t = math.Abs( T[1]*R[0][1] - T[0]*R[1][1] )
    if t > ra + rb {
        return false
    }

    //L = A2 x B2
    ra = a[0]*math.Abs(R[1][2]) + a[1]*math.Abs(R[0][2])
    rb = b[0]*math.Abs(R[2][1]) + b[1]*math.Abs(R[2][0])
    t = math.Abs( T[1]*R[0][2] - T[0]*R[1][2] )
    if t > ra + rb {
        return false
    }

    /*no separating axis found,
    the two boxes overlap */

    return true
}
