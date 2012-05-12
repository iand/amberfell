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
	"math"
)

type Inventory struct {
	visible        bool
	inventorySlots [60]uint16
	componentSlots [6]uint16
	productSlots   [18]*Recipe
	inventoryRects [60]Rect
	componentRects [6]Rect
	productRects   [18]Rect
}

func (self *Inventory) Draw(t int64) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	gl.Color4ub(0, 0, 0, 240)

	gl.Begin(gl.QUADS)
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.bplane))
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.bplane))
	gl.Vertex2f(float32(viewport.rplane), float32(viewport.tplane))
	gl.Vertex2f(float32(viewport.lplane), float32(viewport.tplane))
	gl.End()

	picker.DrawItemHighlight(t, 3)
	picker.DrawItemHighlight(t, 4)
	picker.DrawItemHighlight(t, 5)
	picker.DrawItemHighlight(t, 6)
	picker.DrawItemHighlight(t, 7)
	picker.DrawPlayerItems(t)

	const blocksize = float64(0.3)
	const COLSIZE = 12

	diam := blocksize * 2.4

	offset := diam + float64(4)*PIXEL_SCALE

	for i := 0; i < len(self.inventoryRects); i++ {
		x := float64(viewport.lplane) + float64(10)*PIXEL_SCALE + float64(i/COLSIZE)*offset
		y := float64(viewport.tplane) - float64(10)*PIXEL_SCALE - float64(i%COLSIZE)*offset
		self.inventoryRects[i] = Rect{x, y - diam, diam, diam}
		self.DrawItemSlot(t, self.inventoryRects[i])
	}

	slot := 0
	for i := 1; i < len(ThePlayer.inventory); i++ {
		if ThePlayer.inventory[i] > 0 {
			self.inventorySlots[slot] = uint16(i)
			self.DrawItem(t, ThePlayer.inventory[i], uint16(i), self.inventoryRects[slot])
			slot++
		}
	}

	for i := 0; i < len(self.componentSlots); i++ {
		x := float64(viewport.lplane) + offset*float64(2+len(self.inventoryRects)/COLSIZE) + float64(i)*offset
		y := float64(viewport.tplane) - (float64(10) * PIXEL_SCALE)
		self.componentRects[i] = Rect{x, y - diam, diam, diam}

		self.DrawItemSlot(t, self.componentRects[i])
	}

	for i := 0; i < len(self.componentSlots); i++ {
		if self.componentSlots[i] != 0 {
			self.DrawItem(t, ThePlayer.inventory[self.componentSlots[i]], self.componentSlots[i], self.componentRects[i])
		}
	}

	for i := 0; i < len(self.productSlots); i++ {
		x := float64(viewport.lplane) + offset*float64(2+len(self.inventoryRects)/COLSIZE) + offset*float64(i%len(self.componentRects))
		y := float64(viewport.tplane) - (float64(10) * PIXEL_SCALE) - offset*float64(2+float64(i/len(self.componentRects)))
		self.productRects[i] = Rect{x, y - diam, diam, diam}

		self.DrawItemSlot(t, self.productRects[i])
	}

	for i := 0; i < len(self.productSlots); i++ {
		if self.productSlots[i] != nil {
			self.DrawItem(t, self.productSlots[i].product.quantity, self.productSlots[i].product.item, self.productRects[i])
		}
	}

	var mousex, mousey int
	mousestate := sdl.GetMouseState(&mousex, &mousey)
	self.HandleMouse(mousex, mousey, mousestate)

	gl.PopMatrix()
}

func (self *Inventory) DrawItemSlot(t int64, r Rect) {
	gl.PushMatrix()
	gl.LoadIdentity()

	const blocksize = float32(0.3)

	gl.Color4ub(16, 16, 16, 255)
	gl.Begin(gl.QUADS)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.End()

	gl.Color4ub(6, 6, 6, 255)
	gl.Begin(gl.LINE)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.End()

	gl.Color4ub(64, 64, 64, 255)
	gl.Begin(gl.LINE)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.End()

	gl.PopMatrix()
}

func (self *Inventory) DrawItem(t int64, quantity uint16, blockid uint16, r Rect) {
	gl.PushMatrix()
	gl.LoadIdentity()

	const blocksize = float32(0.3)

	i := 1
	gl.Translated(r.x+r.sizex/2, r.y+r.sizey/2+4*PIXEL_SCALE, 0)

	gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 1.0, 0.0, 0.0)
	gl.Rotatef(360*float32(math.Cos(float64(t)/1e10+float64(i))), 0.0, 1.0, 0.0)
	gl.Rotatef(360*float32(math.Sin(float64(t)/1e10+float64(i))), 0.0, 0.0, 1.0)
	gl.Scalef(blocksize, blocksize, blocksize)
	gVertexBuffer.Reset()
	TerrainCube(gVertexBuffer, 0, 0, 0, [18]uint16{}, blockid, FACE_NONE)
	gVertexBuffer.RenderDirect(false)

	gl.LoadIdentity()
	gl.Translated(r.x+5*PIXEL_SCALE, r.y+2*PIXEL_SCALE, 0)
	inventoryItemFont.Print(fmt.Sprintf("%d", quantity))

	gl.PopMatrix()

}

func (self *Inventory) HasRecipeComponents(recipe *Recipe) bool {
	componentCount := 0

	for k := 0; k < len(self.componentSlots); k++ {
		if self.componentSlots[k] != 0 {
			componentCount++
		}
	}

	if componentCount != len(recipe.components) {
		return false
	}

	for j := 0; j < len(recipe.components); j++ {
		gotComponent := false
		for k := 0; k < len(self.componentSlots); k++ {
			if self.componentSlots[k] == recipe.components[j].item &&
				ThePlayer.inventory[recipe.components[j].item] >= recipe.components[j].quantity {
				gotComponent = true
				break
			}
		}

		if !gotComponent {
			return false
		}

	}

	return true
}

func (self *Inventory) UpdateProducts() {
	for i := 0; i < len(self.productSlots); i++ {
		self.productSlots[i] = nil
	}

	productIndex := 0
	for i := 0; i < len(handmadeRecipes); i++ {
		recipe := &handmadeRecipes[i]
		if self.HasRecipeComponents(recipe) {
			self.productSlots[productIndex] = recipe
			productIndex++
		}
	}

}

func (self *Inventory) HandleMouseButton(re *sdl.MouseButtonEvent) {
	if re.Button == 1 && re.State == 1 { // LEFT, DOWN
		x, y := viewport.ScreenCoordsToWorld2D(re.X, re.Y)

		for i := 0; i < len(self.inventoryRects); i++ {
			if self.inventoryRects[i].Contains(x, y) {
				if self.inventorySlots[i] != 0 {
					keys := sdl.GetKeyState()
					if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
						ThePlayer.EquipItem(self.inventorySlots[i])

					} else {
						// Add to component slots
						for j := 0; j < len(self.componentSlots); j++ {
							if self.componentSlots[j] == self.inventorySlots[i] {
								return // Already in component slot
							}
						}
						for j := 0; j < len(self.componentSlots); j++ {
							if self.componentSlots[j] == 0 {
								self.componentSlots[j] = self.inventorySlots[i]
								self.UpdateProducts()
								return
							}
						}
					}

				}
				return
			}
		}
		for i := 0; i < len(self.componentRects); i++ {
			if self.componentRects[i].Contains(x, y) {
				if self.componentSlots[i] != 0 {
					self.componentSlots[i] = 0
					self.UpdateProducts()
				}
				return
			}
		}
		for i := 0; i < len(self.productRects); i++ {
			if self.productRects[i].Contains(x, y) {
				if self.productSlots[i] != nil {
					recipe := self.productSlots[i]
					if self.HasRecipeComponents(recipe) {
						for j := 0; j < len(recipe.components); j++ {
							ThePlayer.inventory[recipe.components[j].item] -= recipe.components[j].quantity
						}
						ThePlayer.inventory[recipe.product.item] += recipe.product.quantity
						self.UpdateProducts()

						for j := 0; j < len(self.componentSlots); j++ {
							if ThePlayer.inventory[self.componentSlots[j]] <= 0 {
								self.componentSlots[j] = 0
							}
						}

					}
				}
				return
			}
		}

		hit, pos := picker.HitTest(x, y)
		if hit && pos > 2 {
			keys := sdl.GetKeyState()
			if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
				// Remove from picker
				ThePlayer.equippedItems[uint16(pos)-3] = ITEM_NONE
			}
		}

	}

	// type MouseButtonEvent struct {
	// 	Type   uint8
	// 	Which  uint8
	// 	Button uint8
	// 	State  uint8
	// 	X      uint16
	// 	Y      uint16
	// }

}

func (self *Inventory) HandleMouse(mousex int, mousey int, mousestate uint8) {

	x, y := viewport.ScreenCoordsToWorld2D(uint16(mousex), uint16(mousey))
	itemid := uint16(0)

	for i := 0; i < len(self.inventoryRects); i++ {
		if self.inventoryRects[i].Contains(x, y) {
			itemid = self.inventorySlots[i]
			break
		}
	}
	if itemid == 0 {
		for i := 0; i < len(self.componentRects); i++ {
			if self.componentRects[i].Contains(x, y) {
				itemid = self.componentSlots[i]
				break
			}
		}
	}

	if itemid == 0 {
		for i := 0; i < len(self.productRects); i++ {
			if self.productRects[i].Contains(x, y) {
				if self.productSlots[i] != nil {
					itemid = self.productSlots[i].product.item
					break
				}
			}
		}
	}

	if itemid == 0 {
		hit, pos := picker.HitTest(x, y)
		if hit {
			if pos > 2 {
				itemid = uint16(pos) - 3
			}

		}

		for i := 0; i < 5; i++ {
			angle := -(float64(i) + 1.5) * math.Pi / 4
			ix := float64(picker.x) - float64(picker.actionItemRadius)*math.Sin(angle)
			iy := float64(picker.y) + float64(picker.actionItemRadius)*math.Cos(angle)
			if x > ix-float64(picker.selectionRadius) && x < ix+float64(picker.selectionRadius) &&
				y > iy-float64(picker.selectionRadius) && y < iy+float64(picker.selectionRadius) {

				itemid = ThePlayer.equippedItems[i]
				break

			}

		}

	}

	if itemid != 0 {
		self.ShowTooltip(x, y, items[itemid].name)
	}

}

func (self *Inventory) HandleKeyboard(re *sdl.KeyboardEvent) {

}

func (self *Inventory) HandleKeys(keys []uint8) {

}

func (self *Inventory) ShowTooltip(x, y float64, str string) {
	h, w := inventoryItemFont.Measure(str)

	pad := 4 * PIXEL_SCALE
	gl.PushMatrix()

	gl.LoadIdentity()
	gl.Color4ub(0, 0, 0, 255)
	gl.Begin(gl.QUADS)
	gl.Vertex2d(x, y)
	gl.Vertex2d(x+w+pad, y)
	gl.Vertex2d(x+w+pad, y+h)
	gl.Vertex2d(x, y+h)
	gl.End()

	gl.Translated(x+pad/2, y+pad/2, 0)
	inventoryItemFont.Print(str)

	gl.PopMatrix()
}
