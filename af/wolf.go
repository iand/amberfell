package af
import (
    "math"
    "math/rand"
    "github.com/banthar/gl"
)

type Wolf struct {
    heading float64
    position Vector
    velocity Vector
    falling bool
}

func (w *Wolf) Init(heading float64, x float32, z float32, y float32) {
    w.heading = heading
    w.position[XAXIS] = float64(x)
    w.position[YAXIS] = float64(y)
    w.position[ZAXIS] = float64(z)
}

func (w *Wolf) W() float64 { return 0.4 }
func (w *Wolf) H() float64 { return 0.8 }
func (w *Wolf) D() float64 { return 1.4 }

func (w *Wolf) Heading() float64 { return w.heading }
func (w *Wolf) X() float32 { return float32(w.position[XAXIS]) }
func (w *Wolf) Y() float32 { return float32(w.position[YAXIS]) }
func (w *Wolf) Z() float32 { return float32(w.position[ZAXIS]) }
func (w *Wolf) Velocity() Vector { return w.velocity }
func (w *Wolf) Position() Vector { return w.position }
func (w *Wolf) IntPosition() IntVector { 
    return IntVector{ int16(math.Floor(w.position[XAXIS] + 0.5)),
                      int16(math.Floor(w.position[YAXIS] + 0.5)),
                      int16(math.Floor(w.position[ZAXIS] + 0.5))}
}

func (w *Wolf) FrontBlock() IntVector { 
    ip := w.IntPosition()
    if w.heading > 337.5 || w.heading <= 22.5 {
        ip[XAXIS]++
    } else if w.heading > 22.5 && w.heading <= 67.5 {
        ip[XAXIS]++
        ip[ZAXIS]--
    } else if w.heading > 67.5 && w.heading <= 112.5 {
        ip[ZAXIS]--
    } else if w.heading > 112.5 && w.heading <= 157.5 {
        ip[XAXIS]--
        ip[ZAXIS]--
    } else if w.heading > 157.5 && w.heading <= 202.5 {
        ip[XAXIS]--
    } else if w.heading > 202.5 && w.heading <= 247.5 {
        ip[XAXIS]--
        ip[ZAXIS]++
    } else if w.heading > 247.5 && w.heading <= 292.5 {
        ip[ZAXIS]++
    } else if w.heading > 292.5 && w.heading <= 337.5 {
        ip[XAXIS]++
        ip[ZAXIS]++
    }

    return ip
}


func (w *Wolf) SetFalling(b bool) { w.falling = b }

func (w *Wolf) Rotate(angle float64) {
    w.heading += angle
    if w.heading < 0 {
        w.heading += 360
    }
    if w.heading > 360 {
        w.heading -= 360
    }
}

func (w *Wolf) Update(dt float64) {
    w.position[XAXIS] += w.velocity[XAXIS] * dt
    w.position[YAXIS] += w.velocity[YAXIS] * dt
    w.position[ZAXIS] += w.velocity[ZAXIS] * dt
    // fmt.Printf("position: %s\n", w.position)
}

func (w *Wolf) Accelerate(v Vector) {
    w.velocity[XAXIS] += v[XAXIS]
    w.velocity[YAXIS] += v[YAXIS]
    w.velocity[ZAXIS] += v[ZAXIS]
}

func (w *Wolf) IsFalling() bool {
    return w.falling
}

func (w *Wolf) Snapx(x float64, vx float64) {
    w.position[XAXIS] = x
    w.velocity[XAXIS] = vx
}

func (w *Wolf) Snapz(z float64, vz float64) {
    w.position[ZAXIS] = z
    w.velocity[ZAXIS] = vz
}

func (w *Wolf) Snapy(y float64, vy float64) {
    w.position[YAXIS] = y
    w.velocity[YAXIS] = vy
}

func (w *Wolf) Setvx(vx float64) {
    w.velocity[XAXIS] = vx
}

func (w *Wolf) Setvz(vz float64) {
    w.velocity[ZAXIS] = vz
}

func (w *Wolf) Setvy(vy float64) {
    w.velocity[YAXIS] = vy
}

func (w *Wolf) Forward(v float64) {
    w.velocity[XAXIS] = math.Cos(w.Heading() * math.Pi / 180)
    w.velocity[ZAXIS] = -math.Sin(w.Heading() * math.Pi / 180)
}


func (w *Wolf) Act(dt float64) {
    w.Rotate( rand.Float64() * 9 - 4.5 )
    w.Forward( rand.Float64() * 4 - 1)
}



func (wolf *Wolf) Draw(pos Vector, selectMode bool) {
    gl.PushMatrix()
    gl.Translatef(float32(wolf.X() - float32(pos[XAXIS])),float32(wolf.Y() - float32(pos[YAXIS])),float32(wolf.Z() - float32(pos[ZAXIS])))
    gl.Rotated(wolf.Heading(), 0.0, 1.0, 0.0)
    

    h := float32(wolf.H()) / 2
    w := float32(wolf.W()) / 2
    d := float32(wolf.D()) / 2

    gl.Translatef(0.0, h / 2 ,0.0)

    //gl.Translatef(0.0,-h,0.0)
    MapTextures[33].Bind(gl.TEXTURE_2D)
    //topTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        //gl.Color3f(0.3,0.3,0.6)
        // Front face
        gl.Normal3f( 1.0, 0.0, 0.0)
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h,  w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad

    gl.End()

    MapTextures[32].Bind(gl.TEXTURE_2D)

    // dirtTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        // Left Face
        gl.Normal3f( 0.0, 0.0, -1.0)
        gl.TexCoord2f(1.0, 0.0)        
        gl.Vertex3f(-d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Left Of The Texture and Quad


        // Right Face
        //gl.Color3f(0.5,0.5,1.0)              // Set The Color To Blue One Time Only
        gl.Normal3f( 0.0, 0.0, 1.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( -d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(  d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(  d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( -d,  h,  w)  // Top Left Of The Texture and Quad


        // Back Face
        gl.Normal3f( -1.0, 0.0, 0.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d, -h, -w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad

     //gl.Color3f(0.3,1.0,0.3)
        // Top Face
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d,  h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

        // Bottom Face
        gl.Normal3f( 0.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d, -h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d, -h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad


    gl.End();   
   gl.PopMatrix()

}
