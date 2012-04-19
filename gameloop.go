/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
	"time"
)

func InitGame() {
	TheWorld = new(World)
	TheWorld.Init()


    ThePlayer = new(Player)
    ThePlayer.Init(0, 10, 10, FindSurface(10,10))




}

func QuitGame() {
	println("Thanks for playing.")
}

func GameLoop() {
	var currentTime, accumulator int64 = 0, 0
	var t, dt int64 = 0, 1e9 / 40
	var drawFrame, computeFrame int64 = 0, 0
	fps := new(Timer)
	fps.Start()

	update := new(Timer)
	update.Start()

	done := false
	for !done {

		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.ResizeEvent:
				re := e.(*sdl.ResizeEvent)
				screen := sdl.SetVideoMode(int(re.W), int(re.H), 16,
					sdl.OPENGL|sdl.RESIZABLE)
				if screen != nil {
					Reshape(int(screen.W), int(screen.H))
				} else {
					panic("we couldn't set the new video mode??")
				}
				break

            case *sdl.MouseButtonEvent:
				re := e.(*sdl.MouseButtonEvent)
				if re.Button == 1 && re.State == 1 { // LEFT, DOWN

                    println("Click:", re.X, re.Y, re.State, re.Button, re.Which)
                    println("feedbackBuffer.size:", feedbackBuffer.size)
                    feedbackBuffer.Dump()




					if ThePlayer.CanInteract() {

						// println("Click:", re.X, re.Y, re.State, re.Button, re.Which)

						// MOUSEBUTTONDOWNMASK
						xv, yv := int(re.X), screenHeight-int(re.Y)
						data := [4]uint8{0, 0, 0, 0}

						Draw(true)
						gl.ReadPixels(xv, yv, 1, 1, gl.RGBA, &data[0])
						Draw(false)

						fmt.Printf("pixel data: %d, %d, %d, %d\n", data[0], data[1], data[2], data[3])

						id := uint16(data[0]) + uint16(data[1])*256
						if id != 0 {
							face := data[2]
							dx, dy, dz := BlockIdToRelativeCoordinate(id)
							fmt.Printf("id: %d, dx: %d, dy: %d, dz: %d, face: %d\n", id, dx, dy, dz, face)
							if !(dx == 0 && dy == 0 && dz == 0) {
								pos := IntPosition(ThePlayer.Position())
								pos.Adjust(dx, dy, dz)
								ThePlayer.Interact(pos, face)

							}
						}
					}
				}
            case *sdl.MouseMotionEvent:
                //re := e.(*sdl.MouseMotionEvent)
                if ThePlayer.CanInteract() {

                    // println("Move:", re.X, re.Y, re.Xrel, re.Yrel)

                    // // MOUSEBUTTONDOWNMASK
                    // xv, yv := int(re.X), screenHeight-int(re.Y)
                    // data := [4]uint8{0, 0, 0, 0}

                    // Draw(true)
                    // gl.ReadPixels(xv, yv, 1, 1, gl.RGBA, &data[0])
                    // Draw(false)

                    // fmt.Printf("pixel data: %d, %d, %d, %d\n", data[0], data[1], data[2], data[3])
                }

			case *sdl.QuitEvent:
				done = true
				break
			}
		}
		keys := sdl.GetKeyState()

		if keys[sdl.K_ESCAPE] != 0 {
			ShowOverlay = !ShowOverlay

			// Overlay

		}
		if keys[sdl.K_UP] != 0 {

			if view_rotx < 75 {
				view_rotx += 5.0
			}
		}
		if keys[sdl.K_DOWN] != 0 {
			if view_rotx > 15.0 {
				view_rotx -= 5.0
			}
		}
		if keys[sdl.K_LEFT] != 0 {
			view_roty += 9
			//println("view_roty:", view_roty)
		}
		if keys[sdl.K_RIGHT] != 0 {
			view_roty -= 9
		}

		ThePlayer.HandleKeys(keys)

		if keys[sdl.K_F3] != 0 {
			if DebugMode == true {
				DebugMode = false
			} else {
				DebugMode = true
			}
		}

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

		if DebugMode {
			fmt.Printf("x:%f, z:%f\n", ThePlayer.X(), ThePlayer.Z())
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


        Draw(false)
		drawFrame++

		if update.GetTicks() > 1e9/2 {
			fmt.Printf("draw fps: %f\n", float64(drawFrame)/(float64(update.GetTicks())/float64(1e9)))
			fmt.Printf("compute fps: %f\n", float64(computeFrame)/(float64(update.GetTicks())/float64(1e9)))

			timeOfDay += 0.02
			if timeOfDay > 24 {
				timeOfDay -= 24
			}

			drawFrame, computeFrame = 0, 0
			update.Start()
		}

	}
}
