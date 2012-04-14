package af

import (
    "math/rand"
    "github.com/banthar/gl"
    // "fmt"   

)


type World struct {
    W           int16
    D           int16
    H           int16
    GroundLevel int16
    Blocks      []byte
    mobs        []Mob
}

type Side struct {
    x, x1, x2, z, z1, z2, dir, y float64
}

type Rect struct {
    x, z, sizex, sizez float64
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
                world.Set(iw, ih, id, 2) // dirt
            }
            for ih = GroundLevel + 1; ih < h; ih++ {
                world.Set(iw, ih, id, 0) // air
            }
        }
    }

    numFeatures := rand.Intn(int(world.W))
    for i := 0; i < numFeatures; i++ {
        iw, id = world.RandomSquare()

        world.Set(iw, GroundLevel, id, 1) // stone
        world.Grow(iw, GroundLevel, id, 40, 40, 40, 40, 10, 30, 1)
    }
    iw, id = world.RandomSquare()

    world.Set(iw, GroundLevel, id, 0) // air
    world.Grow(iw, GroundLevel, id, 30, 30, 30, 30, 0, 30, 0)


    wolf := new(Wolf)
    wolf.Init(120, 14, 14, GroundLevel+1)
    world.mobs = append(world.mobs, wolf)


}

func (world *World) At(x int16, y int16, z int16) byte {
    x = x % world.W
    if x < 0 { x += world.W }
    z = z % world.D
    if z < 0 { z += world.D }
    if y < 0 || y > world.H-1 {
        return 0
    }
    return world.Blocks[world.W*world.D*y+world.D*x+z]
}

func (world *World) Set(x int16, y int16, z int16, b byte) {
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

func (world *World) Grow(x int16, y int16, z int16, n int, s int, w int, e int, u int, d int, texture byte) {
    if x < world.W-1 && world.At(x+1, y-1, z) != 0 && rand.Intn(100) < e {
        world.Set(x+1, y, z, texture)
        world.Grow(x+1, y, z, n, s, 0, e, u, d, texture)
    }
    if x > 0 && world.At(x-1, y-1, z) != 0 && rand.Intn(100) < w {
        world.Set(x-1, y, z, texture)
        world.Grow(x-1, y, z, n, s, w, 0, u, d, texture)
    }
    if y < world.D-1 && world.At(x, y-1, z+1) != 0 && rand.Intn(100) < s {
        world.Set(x, y, z+1, texture)
        world.Grow(x, y, z+1, n, 0, w, e, u, d, texture)
    }
    if y > 0 && world.At(x, y-1, z-1) != 0 && rand.Intn(100) < n {
        world.Set(x, y, z-1, texture)
        world.Grow(x, y, z-1, 0, s, w, e, u, d, texture)
    }
    if y < world.H-1 && rand.Intn(100) < u {
        world.Set(x, y+1, z, texture)
        world.Grow(x, y+1, z, n, s, w, e, u, 0, texture)
    }
    if y > 0 && rand.Intn(100) < d {
        world.Set(x, y-1, z, texture)
        world.Grow(x, y-1, z, n, s, w, e, 0, d, texture)
    }
}

func (world *World) AirNeighbours(x int16, z int16, y int16) (n, s, w, e, u, d bool) {
    if /* x > 0 && */world.At(x-1, y, z) == 0 {
        e = true
    }
    if /*x < world.W-1 &&*/ world.At(x+1, y, z) == 0 {
        w = true
    }
    if /*z > 0 &&*/ world.At(x, y, z-1) == 0 {
        s = true
    }
    if /*z < world.D-1 &&*/ world.At(x, y, z+1) == 0 {
        n = true
    }
    if /*y < world.H-1 &&*/ world.At(x, y+1, z) == 0 {
        u = true
    }
    return
}

// lineRectCollide( line, rect )
//
// Checks if an axis-aligned line and a bounding box overlap.
// line = { z, x1, x2 } or line = { x, z1, z2 }
// rect = { x, z, size }

func lineRectCollide( line Side, rect Rect ) (ret bool) {
    if line.z != 0  {
        ret = rect.z > line.z - rect.sizez/2 && rect.z < line.z + rect.sizez/2 && rect.x > line.x1 - rect.sizex/2 && rect.x < line.x2 + rect.sizex/2
    } else {
        ret = rect.x > line.x - rect.sizex/2 && rect.x < line.x + rect.sizex/2 && rect.z > line.z1 - rect.sizez/2 && rect.z < line.z2 + rect.sizez/2
    }
    return 
}

// rectRectCollide( r1, r2 )
//
// Checks if two rectangles (x1, y1, x2, y2) overlap.

func rectRectCollide( r1 Side, r2 Side) bool {
    if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 {
        return true
    }
    if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z1 >= r1.z1 && r2.z1 <= r1.z2 { 
        return true
    }
    if r2.x2 >= r1.x1 && r2.x2 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
        return true
    }
    if r2.x1 >= r1.x1 && r2.x1 <= r1.x2 && r2.z2 >= r1.z1 && r2.z2 <= r1.z2 {
        return true
    }
    return false
}



func (world *World) ApplyForces(mob Mob, dt float64) {
    // mobBounds := mob.DesiredBoundingBox(dt)
    ip := mob.IntPosition()

    mobx := ip[XAXIS]
    moby := ip[YAXIS]
    mobz := ip[ZAXIS]
 

    // Gravity
    if mob.IsFalling() {
        mob.Accelerate( Vector{0, -0.5, 0} )
    }

    // var dx, dz, dy int16
    var x, y, z int16


    playerRect := Rect{ x: float64(mob.X()) + mob.Velocity()[XAXIS] * dt, z: float64(mob.Z()) + mob.Velocity()[ZAXIS] * dt, sizex: mob.W(), sizez: mob.D() };

    collisionCandidates := make([]Side, 0)

    // Collect XZ collision candidates
    for x = mobx - 1; x <= mobx + 1; x++ {
        for y = moby; y <= moby + 1; y++ {
            for z = mobz - 1; z <= mobz + 1; z++ {
                if world.At(x, y, z) != 0 {
                    if world.At( x - 1, y, z ) == 0  {
                        collisionCandidates = append(collisionCandidates, Side{ x: float64(x) - 0.5, dir: -1, z1: float64(z) - 0.5, z2: float64(z) + 0.5 } )
                    }
                    if world.At( x + 1, y, z ) == 0  {
                        collisionCandidates = append(collisionCandidates, Side{ x: float64(x) + 0.5, dir: 1, z1: float64(z) - 0.5, z2: float64(z) + 0.5 } )
                    }
                    if world.At( x, y, z -1 ) == 0  {
                        // fmt.Printf("float64(z) - 0.5: %f\n", float64(z) - 0.5)
                        collisionCandidates = append(collisionCandidates, Side{ z: float64(z) - 0.5, dir: -1, x1: float64(x) - 0.5, x2: float64(x) + 0.5 } )
                    }
                    if world.At( x, y, z + 1 ) == 0  {
                        collisionCandidates = append(collisionCandidates, Side{ z: float64(z) + 0.5, dir: 1, x1: float64(x) - 0.5, x2: float64(x) + 0.5 } )
                    }                  
                }
            }
        }
    }


    // Solve XZ collisions
    for _, side := range collisionCandidates {
        if lineRectCollide(side, playerRect) {
        // fmt.Printf("side.x: %f\n", side.x)
            if side.x != 0 && mob.Velocity()[XAXIS] * side.dir < 0 {
                // fmt.Printf("Snapping x\n")
                mob.Snapx(side.x + (side.dir * playerRect.sizex/2), 0)
            }
            if side.z != 0 && mob.Velocity()[ZAXIS] * side.dir < 0  {
                // fmt.Printf("Snapping z\n")
                mob.Snapz(side.z + (side.dir * playerRect.sizez/2) , 0)
            }
        }
     }



    playerFace := Side{ x1: float64(mob.X()) + mob.Velocity()[XAXIS] * dt - 0.5, z1: float64(mob.Z()) + mob.Velocity()[ZAXIS] * dt - 0.5, x2: float64(mob.X()) + mob.Velocity()[XAXIS] * dt + 0.5, z2: float64(mob.Z()) + mob.Velocity()[ZAXIS] * dt + 0.5 }

    // fmt.Printf("playerFace x1:%f, x2:%f, z1:%f, z2:%f\n", playerFace.x1, playerFace.x2, playerFace.z1, playerFace.z2)

    newBZLower := moby - 1// int16(math.Floor( float64(mob.Y()) + mob.Velocity()[YAXIS] * dt ))
    // newBZUpper := int16(math.Floor( float64(mob.Y()) + 1.7 + mob.Velocity()[YAXIS] * 1.1  * dt ))

    //fmt.Printf("newBZLower: %d, newBZUpper: %d, mob.Y(): %f\n", newBZLower, newBZUpper,  mob.Y())
    collisionCandidates = make([]Side, 0)

    for x = mobx - 1; x <= mobx + 1; x++ {
        for z = mobz - 1; z <= mobz + 1; z++ {

            if world.At( x, newBZLower, z ) != 0 {
                collisionCandidates =  append(collisionCandidates, Side{ y: float64(newBZLower) + 0.5, dir: 1, x1: float64(x) - 0.5, z1: float64(z) - 0.5, x2: float64(x) + 0.5, z2: float64(z) + 0.5 } )
            }
            // if world.At( x, newBZUpper, z ) != 0 {
            //     collisionCandidates =  append(collisionCandidates, Side{ y: float64(newBZUpper),     dir: -1, x1: float64(x) - 0.5, z1: float64(z) - 0.5, x2: float64(x) + 0.5, z2: float64(z) + 0.5 } );
            // }
        }
    }

    // Solve Y collisions
    mob.SetFalling(true)
    for _, face := range collisionCandidates {

        // fmt.Printf("face x1:%f, x2:%f, z1:%f, z2:%f\n", face.x1, face.x2, face.z1, face.z2)
        if rectRectCollide( face, playerFace ) && mob.Velocity()[YAXIS] * face.dir < 0  {
            if mob.Velocity()[YAXIS] < 0 {
                mob.SetFalling(false)
                mob.Snapy(face.y + 0.5, 0)
            } else {
                mob.Snapy(face.y + 0.5, 0)
            }
            break
        }
    }

}





func (w *World) Simulate(dt float64) {
    for _, v := range w.mobs {
        v.Act(dt)
        w.ApplyForces(v, dt)
        v.Update(dt)
    }
}


func (world *World) Draw(pos Vector, selectMode bool) {
    for _, v := range world.mobs {
        v.Draw(pos, selectMode)
    }    

    gl.Translatef(-float32(pos[XAXIS]), -float32(pos[YAXIS]), -float32(pos[ZAXIS]))

    var px, py, pz = int16(pos[XAXIS]), int16(pos[YAXIS]), int16(pos[ZAXIS])
    var x, y, z int16

    for x = px - 30; x < px + 30; x++ {
        for z=pz - 30; z < pz + 30; z++ {
            for y=0; y < world.H; y++ {
                dx := x - px
                dy := y - py
                dz := z - pz

                var terrain byte = world.At(x, y, z)
                if terrain != 0 {
                    var n, s, w, e, u, d bool = world.AirNeighbours(x, z, y)
                    var id uint16 = 0

                    if dx >= -2 && dx <= 2 && dy >= -2 && dy <= 2 && dz >= -2 && dz <= 2 {
                        id = RelativeCoordinateToBlockId(dx, dy, dz)
                    }
                    gl.PushMatrix()
                    gl.Translatef(float32(x),float32(y),float32(z))
                    //print ("i:", i, "j:", j, "b:", world.At(i, j, groundLevel))
                    Cube(n, s, w, e, u, d, terrain, id, selectMode)
                    gl.PopMatrix()
                }
            }
        }
    }

}

