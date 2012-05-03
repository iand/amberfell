/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/gl"
	"runtime"
)

type Console struct {
	fps       float64
	mem       runtime.MemStats
	cubecount int
	vertices  int
}

func (self *Console) Update() {
	runtime.ReadMemStats(&self.mem)
}

func (self *Console) Draw(t int64) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	h := float32(consoleFont.height) * PIXEL_SCALE
	margin := float32(3.0) * PIXEL_SCALE
	consoleHeight := 3 * h

	gl.MatrixMode(gl.MODELVIEW)

	gl.LoadIdentity()
	gl.Color4ub(0, 0, 0, 208)

	gl.Begin(gl.QUADS)
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane)+consoleHeight+margin*2) // Bottom Left Of The Texture and Quad
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane)+consoleHeight+margin*2) // Bottom Right Of The Texture and Quad
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane))                        // Top Right Of The Texture and Quad
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane))                        // Top Left Of The Texture and Quad
	gl.End()

	gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-h, 0)
	consoleFont.Print(fmt.Sprintf("FPS: %5.2f CC: %d V: %d", self.fps, self.cubecount, self.vertices))
	gl.LoadIdentity()
	gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-2*h, 0)

	consoleFont.Print(fmt.Sprintf("X: %5.2f Y: %4.2f Z: %5.2f H: %5.2f (%s)", ThePlayer.position[XAXIS], ThePlayer.position[YAXIS], ThePlayer.position[ZAXIS], ThePlayer.heading, HeadingToCompass(ThePlayer.heading)))

	gl.LoadIdentity()
	gl.Translatef(float32(viewport.lplane)+margin, float32(viewport.bplane)+consoleHeight+margin-3*h, 0)

	numgc := uint32(0)
	avggc := float64(0)
	var last3 [3]float64
	if self.mem.NumGC > 3 {
		numgc = self.mem.NumGC
		avggc = float64(self.mem.PauseTotalNs) / float64(self.mem.NumGC) / 1e6
		index := int(numgc) - 1
		if index > 255 {
			index = 255
		}

		last3[0] = float64(self.mem.PauseNs[index]) / 1e6
		last3[1] = float64(self.mem.PauseNs[index-1]) / 1e6
		last3[2] = float64(self.mem.PauseNs[index-2]) / 1e6
	}

	consoleFont.Print(fmt.Sprintf("Mem: %.1f/%.1f   GC: %.1fms [%d: %.1f, %.1f, %.1f]", float64(self.mem.Alloc)/(1024*1024), float64(self.mem.Sys)/(1024*1024), avggc, numgc, last3[0], last3[1], last3[2]))

	gl.PopMatrix()
}
