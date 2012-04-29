/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	// "fmt"
	"github.com/banthar/Go-SDL/sdl"
	// "github.com/banthar/gl"
	"runtime"
	"time"
)

func GameLoop() {
	var startTime int64 = time.Now().UnixNano()
	var currentTime, accumulator int64 = 0, 0
	var t, dt int64 = 0, 1e9 / 40
	var drawFrame, computeFrame int64 = 0, 0

	update500ms := new(Timer)
	update500ms.interval = 500 * 1e6
	update500ms.Start()

	update150ms := new(Timer)
	update150ms.interval = 50 * 1e6
	update150ms.Start()

	debugModekeyLock := false

	var interactingBlock *InteractingBlockFace

	done := false
	for !done {

		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.ResizeEvent:
				re := e.(*sdl.ResizeEvent)
				screen := sdl.SetVideoMode(int(re.W), int(re.H), 16,
					sdl.OPENGL|sdl.RESIZABLE)
				if screen != nil {
					viewport.Reshape(int(screen.W), int(screen.H))
				} else {
					panic("Could not set video mode")
				}
				break

			case *sdl.MouseButtonEvent:
				re := e.(*sdl.MouseButtonEvent)
				if re.Button == 1 && re.State == 1 { // LEFT, DOWN
					if ThePlayer.CanInteract() {
						selectedBlockFace := viewport.SelectedBlockFace()
						if selectedBlockFace != nil {
							if interactingBlock == nil || interactingBlock.blockFace.pos != selectedBlockFace.pos {
								interactingBlock = new(InteractingBlockFace)
								interactingBlock.blockFace = selectedBlockFace
								interactingBlock.hitCount = 0
							}
							ThePlayer.Interact(interactingBlock)
						}
						// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
					}
				}
			case *sdl.QuitEvent:
				done = true
			case *sdl.KeyboardEvent:
				re := e.(*sdl.KeyboardEvent)
				if re.Keysym.Sym == sdl.K_F3 {
					if !debugModekeyLock && re.Type == sdl.KEYDOWN {
						debugModekeyLock = true
						if DebugMode == true {
							DebugMode = false
						} else {
							DebugMode = true
						}
					} else if re.Type == sdl.KEYUP {
						debugModekeyLock = false
					}
				}
			}
		}
		keys := sdl.GetKeyState()

		if keys[sdl.K_ESCAPE] != 0 {
			// ShowOverlay = !ShowOverlay

			// Overlay

		}
		if keys[sdl.K_UP] != 0 {
			viewport.Rotx(5)
		}
		if keys[sdl.K_DOWN] != 0 {
			viewport.Rotx(-5)
		}
		if keys[sdl.K_LEFT] != 0 {
			viewport.Roty(9)
		}
		if keys[sdl.K_RIGHT] != 0 {
			viewport.Roty(-9)
		}

		ThePlayer.HandleKeys(keys)

		if keys[sdl.K_o] != 0 {
			timeOfDay += 0.1
			if timeOfDay > 24 {
				timeOfDay -= 24
			}
		}
		if keys[sdl.K_l] != 0 {
			timeOfDay -= 0.1
			if timeOfDay < 0 {
				timeOfDay += 24
			}
		}

		if keys[sdl.K_i] != 0 {
			if ViewRadius < 90 {
				ViewRadius++
				println("ViewRadius: ", ViewRadius)
			}
		}
		if keys[sdl.K_k] != 0 {
			if ViewRadius > 10 {
				ViewRadius--
				println("ViewRadius: ", ViewRadius)
			}
		}

		if keys[sdl.K_u] != 0 {
			viewport.Zoomin()

		}
		if keys[sdl.K_j] != 0 {
			viewport.Zoomout()
		}

		if update150ms.PassedInterval() {
			debugModekeyLock = false

			// If player is breaking a block then allow them to hold mouse down to continue action
			if interactingBlock != nil && ThePlayer.currentAction == ACTION_BREAK {
				mouseState := sdl.GetMouseState(nil, nil)
				if mouseState == 1 {
					if ThePlayer.CanInteract() {
						selectedBlockFace := viewport.SelectedBlockFace()
						if selectedBlockFace != nil {
							if interactingBlock == nil || !interactingBlock.blockFace.pos.Equals(&selectedBlockFace.pos) {
								interactingBlock = new(InteractingBlockFace)
								interactingBlock.blockFace = selectedBlockFace
								interactingBlock.hitCount = 0
							}
							ThePlayer.Interact(interactingBlock)
						}
						// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
					}
				}
			}
			update150ms.Start()
		}

		if update500ms.PassedInterval() {
			metrics.fps = float64(drawFrame) / (float64(update500ms.GetTicks()) / float64(1e9))
			runtime.ReadMemStats(&metrics.mem)
			timeOfDay += 0.02
			if timeOfDay > 24 {
				timeOfDay -= 24
			}

			drawFrame, computeFrame = 0, 0
			update500ms.Start()
		}

		newTime := time.Now().UnixNano()
		deltaTime := newTime - currentTime
		currentTime = newTime
		if deltaTime > 1e9/4 {
			deltaTime = 1e9 / 4
		}

		accumulator += deltaTime

		for accumulator > dt {
			accumulator -= dt

			TheWorld.ApplyForces(ThePlayer, float64(dt)/1e9)

			ThePlayer.Update(float64(dt) / 1e9)
			TheWorld.Simulate(float64(dt) / 1e9)

			computeFrame++
			t += dt
		}

		//interpolate(previous, current, accumulator/dt)

		Draw(currentTime - startTime)
		drawFrame++
	}

}
