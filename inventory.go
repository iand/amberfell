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
	inventorySlots [72]ItemQuantity
	componentSlots [6]ItemQuantity
	productSlots   [18]*Recipe
	inventoryRects [72]Rect
	componentRects [6]Rect
	productRects   [18]Rect

	selectedItem *SelectedItem

	currentContainer ContainerObject
	containerRects   []Rect
	currentCrafting  CraftingObject
	drawBuffer       *VertexBuffer
}

type ItemSlot struct {
	area  uint8
	index int
}

type SelectedItem struct {
	ItemQuantity
	ItemSlot
}

const (
	AREA_INVENTORY          = 1
	AREA_HANDHELD_COMPONENT = 2
	AREA_HANDHELD_PRODUCT   = 3
	AREA_CONTAINER          = 4
)

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.drawBuffer = NewVertexBuffer(1000, terrainTexture)
	return inv
}

func (self *Inventory) Draw(t int64) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.LIGHT0)
	gl.Disable(gl.LIGHT1)

	gl.Color4ub(0, 0, 0, 208)

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
	picker.DrawPlayerItems(t, false)

	const blocksize = float64(0.3)
	const COLSIZE = 12

	slotsize := blocksize * 2.4
	slotstep := slotsize + float64(4)*PIXEL_SCALE
	slotsInRow := len(self.componentRects)

	xtools := float64(viewport.lplane) + 10.0*PIXEL_SCALE
	ytools := float64(viewport.tplane) - 10.0*PIXEL_SCALE - slotstep

	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Translated(xtools, ytools+4*PIXEL_SCALE, 0)
	inventoryItemFont.Print("Inventory")
	gl.PopMatrix()
	for i := range self.inventoryRects {
		x := xtools + float64(i%slotsInRow)*slotstep
		y := ytools - float64(i/slotsInRow)*slotstep
		self.inventoryRects[i] = Rect{x, y - slotsize, slotsize, slotsize}
		self.DrawItemSlot(t, self.inventoryRects[i])
	}

	for i := range self.inventorySlots {
		if self.inventorySlots[i].item != 0 && self.inventorySlots[i].quantity > 0 {
			self.DrawItemInSlot(t, self.inventorySlots[i].quantity, self.inventorySlots[i].item, self.inventoryRects[i])
		}
	}

	xtools += slotstep * (1.0 + float64(slotsInRow))

	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Translated(xtools, ytools+4*PIXEL_SCALE, 0)

	if self.currentCrafting == nil {
		inventoryItemFont.Print("Handcrafting")
	} else {
		inventoryItemFont.Print(self.currentCrafting.Label())
	}
	gl.PopMatrix()

	for i := range self.componentSlots {
		x := xtools + float64(i)*slotstep
		y := ytools
		self.componentRects[i] = Rect{x, y - slotsize, slotsize, slotsize}

		self.DrawItemSlot(t, self.componentRects[i])
	}

	for i, cs := range self.componentSlots {
		if cs.item != 0 {
			self.DrawItemInSlot(t, cs.quantity, cs.item, self.componentRects[i])
		}
	}

	ytools -= slotstep * 2
	for i := range self.productSlots {
		x := xtools + slotstep*float64(i%slotsInRow)
		y := ytools - slotstep*float64(i/slotsInRow)
		self.productRects[i] = Rect{x, y - slotsize, slotsize, slotsize}

		self.DrawItemSlot(t, self.productRects[i])
	}

	for i, ps := range self.productSlots {
		if ps != nil {
			self.DrawItemInSlot(t, ps.product.quantity, ps.product.item, self.productRects[i])
		}
	}

	ytools -= slotstep * float64(1+len(self.productRects)/slotsInRow)
	if self.currentContainer != nil {

		gl.PushMatrix()
		gl.LoadIdentity()
		gl.Translated(xtools, ytools+4*PIXEL_SCALE, 0)
		inventoryItemFont.Print(self.currentContainer.Label())
		gl.PopMatrix()

		for i := range self.containerRects {
			x := xtools + slotstep*float64(i%slotsInRow)
			y := ytools - slotstep*float64(i/slotsInRow)
			self.containerRects[i] = Rect{x, y - slotsize, slotsize, slotsize}

			self.DrawItemSlot(t, self.containerRects[i])
		}

		for i := 0; i < self.currentContainer.Slots(); i++ {
			if self.currentContainer.Item(i).item != 0 {
				self.DrawItemInSlot(t, self.currentContainer.Item(i).quantity, self.currentContainer.Item(i).item, self.containerRects[i])
			}
		}

	}

	var mousex, mousey int
	mousestate := sdl.GetMouseState(&mousex, &mousey)

	if self.selectedItem != nil {
		x, y := viewport.ScreenCoordsToWorld2D(uint16(mousex), uint16(mousey))
		self.DrawItem(t, self.selectedItem.quantity, self.selectedItem.item, x, y)
	}

	self.HandleMouse(mousex, mousey, mousestate)

	gl.PopMatrix()
}

func (self *Inventory) DrawItemSlot(t int64, r Rect) {
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.Color4ub(16, 16, 16, 255)
	gl.Begin(gl.QUADS)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.End()

	gl.Color4ub(6, 6, 6, 255)
	gl.Begin(gl.LINES)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x, r.y)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.End()

	gl.Color4ub(64, 64, 64, 255)
	gl.Begin(gl.LINES)
	gl.Vertex2d(r.x+r.sizex, r.y)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.Vertex2d(r.x, r.y+r.sizey)
	gl.Vertex2d(r.x+r.sizex, r.y+r.sizey)
	gl.End()

	gl.PopMatrix()
}

func (self *Inventory) DrawItemInSlot(t int64, quantity uint16, itemid uint16, r Rect) {
	self.DrawItem(t, quantity, itemid, r.x+r.sizex/2, r.y+r.sizey/2+4*PIXEL_SCALE)

}

func (self *Inventory) DrawItem(t int64, quantity uint16, itemid uint16, x float64, y float64) {
	gl.PushMatrix()
	gl.LoadIdentity()

	// angle := -90.0
	const blocksize = float32(0.3)

	gl.Translated(x, y, 0)

	//gl.Rotated(angle, 1.0, 0.0, 0.0)
	gl.Rotated(90, 1.0, 0.0, 0.0)
	gl.Rotated(30*math.Sin(float64(t)/1e9+float64(itemid)/2), 0.0, 1.0, 0.0)
	gl.Scalef(blocksize, blocksize, blocksize)
	self.drawBuffer.Reset()

	if itemid < 256 {
		TerrainCube(self.drawBuffer, Vectori{}, [18]uint16{BLOCK_DIRT, BLOCK_DIRT, BLOCK_DIRT, BLOCK_DIRT, BLOCK_AIR, BLOCK_DIRT, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR, BLOCK_AIR}, NewBlockDefault(byte(itemid)), FACE_NONE)
	} else {
		RenderItemFlat(self.drawBuffer, Vectori{}, itemid)
	}
	self.drawBuffer.RenderDirect(false)

	gl.LoadIdentity()
	gl.Translated(x+2*PIXEL_SCALE-48*PIXEL_SCALE*float64(blocksize), y-7*PIXEL_SCALE-48*PIXEL_SCALE*float64(blocksize), 0)
	inventoryItemFont.Print(fmt.Sprintf("%d", quantity))

	gl.PopMatrix()

}

func (self *Inventory) HasRecipeComponents(recipe *Recipe) bool {
	componentCount := 0

	for k := range self.componentSlots {
		if self.componentSlots[k].quantity != 0 {
			componentCount++
		}
	}

	if componentCount != len(recipe.components) {
		return false
	}

	for _, rc := range recipe.components {
		gotComponent := false
		for k := range self.componentSlots {
			if self.componentSlots[k].item == rc.item && self.componentSlots[k].quantity >= rc.quantity {
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
	for i := range self.productSlots {
		self.productSlots[i] = nil
	}

	recipes := handmadeRecipes
	if self.currentCrafting != nil {
		recipes = self.currentCrafting.Recipes()
	}

	productIndex := 0
	for i := range recipes {
		recipe := &recipes[i]
		if self.HasRecipeComponents(recipe) {
			self.productSlots[productIndex] = recipe
			productIndex++
		}
	}

}

func (self *Inventory) HandleMouseButton(re *sdl.MouseButtonEvent) {
	if re.Button == 1 && re.State == 1 { // LEFT, DOWN
		x, y := viewport.ScreenCoordsToWorld2D(re.X, re.Y)

		// Nothing selected, left click == pick one up
		// Nothing selected, shift left click == pick all up

		// Item selected, same source, same item in slot, left click == pick another up
		// Item selected, same source, same item in slot, shift left click == pick all up
		// Item selected, same source, same item in slot, ctrl left click == drop one
		// Item selected, same source, same item in slot, ctrl shift left click == drop all

		// Item selected, same source, empty slot, left click == drop one
		// Item selected, same source, empty slot, shift left click == drop all
		// Item selected, same source, empty slot, ctrl left click == N/A
		// Item selected, same source, empty slot, ctrl shift left click == N/A

		// Item selected, diff source, same item in slot, left click == drop one
		// Item selected, diff source, same item in slot, shift left click == drop all
		// Item selected, diff source, same item in slot, ctrl left click == pick one up
		// Item selected, diff source, same item in slot, ctrl shift left click == pick all up

		// Item selected, diff source, empty slot, left click == drop one
		// Item selected, diff source, empty slot, shift left click == drop all
		// Item selected, diff source, empty slot, ctrl left click == pick one up
		// Item selected, diff source, empty slot, ctrl shift left click == pick all up

		keys := sdl.GetKeyState()

		bulk := false
		invert := false
		if keys[sdl.K_LSHIFT] != 0 || keys[sdl.K_RSHIFT] != 0 {
			bulk = true
		}
		if keys[sdl.K_LCTRL] != 0 || keys[sdl.K_RCTRL] != 0 {
			invert = true
		}

		for i := range self.inventoryRects {
			if self.inventoryRects[i].Contains(x, y) {
				slot := &self.inventorySlots[i]

				if self.selectedItem == nil {
					// pick up item, if any
					if slot.item != 0 && slot.quantity > 0 {
						if bulk {
							// Act on all items in slot
							self.selectedItem = &SelectedItem{ItemQuantity{slot.item, slot.quantity}, ItemSlot{AREA_INVENTORY, i}}
							slot.quantity = 0
							slot.item = 0
						} else {
							// Act on a single item
							self.selectedItem = &SelectedItem{ItemQuantity{slot.item, 1}, ItemSlot{AREA_INVENTORY, i}}
							slot.quantity -= 1
							if slot.quantity == 0 {
								slot.item = 0
							}
						}
					}
				} else if slot.item == self.selectedItem.item {
					// Clicked on slot containing same item as current selection
					// Is this the same slot as the source of items
					if self.selectedItem.area == AREA_INVENTORY && self.selectedItem.index == i {
						// Same slot
						if invert {
							// drop back into same slot
							if bulk {
								// Act on all items in slot
								slot.quantity += self.selectedItem.quantity
								self.selectedItem.quantity = 0
							} else {
								// Act on a single item
								self.selectedItem.quantity -= 1
								slot.quantity += 1
							}

						} else {
							// pick more up
							if bulk {
								// Act on all items in slot
								self.selectedItem.quantity += slot.quantity
								slot.quantity = 0
							} else {
								// Act on a single item
								self.selectedItem.quantity += 1
								slot.quantity -= 1
							}
						}
					} else {
						// Different slot
						if invert {
							if bulk {
								self.selectedItem.quantity += slot.quantity
								slot.quantity = 0
							} else {
								self.selectedItem.quantity += 1
								slot.quantity -= 1
							}
						} else {
							if bulk {
								slot.quantity += self.selectedItem.quantity
								self.selectedItem.quantity = 0
							} else {
								self.selectedItem.quantity -= 1
								slot.quantity += 1
							}
						}
					}
				} else if slot.item == 0 {
					quantity := uint16(1)
					if bulk {
						quantity = self.selectedItem.quantity
					}

					self.selectedItem.quantity -= quantity
					slot.quantity += quantity
					slot.item = self.selectedItem.item
				}

				if self.selectedItem != nil && self.selectedItem.quantity == 0 {
					self.selectedItem = nil
				}

				if slot.quantity == 0 {
					slot.item = 0
				}

				return
			}
		}
		for i := range self.componentRects {
			if self.componentRects[i].Contains(x, y) {
				slot := &self.componentSlots[i]
				slotQuantity := slot.quantity
				if self.selectedItem == nil {
					// pick up item, if any
					if slot.item != 0 && slot.quantity > 0 {
						if bulk {
							// Act on all items in slot
							self.selectedItem = &SelectedItem{ItemQuantity{slot.item, slot.quantity}, ItemSlot{AREA_HANDHELD_COMPONENT, i}}
							slot.quantity = 0
							slot.item = 0
						} else {
							// Act on a single item
							self.selectedItem = &SelectedItem{ItemQuantity{slot.item, 1}, ItemSlot{AREA_HANDHELD_COMPONENT, i}}
							slot.quantity -= 1
							if slot.quantity == 0 {
								slot.item = 0
							}
						}
					}
				} else if slot.item == self.selectedItem.item {
					// Clicked on slot containing same item as current selection
					// Is this the same slot as the source of items
					if self.selectedItem.area == AREA_HANDHELD_COMPONENT && self.selectedItem.index == i {
						// Same slot
						if invert {
							// drop back into same slot
							if bulk {
								// Act on all items in slot
								slot.quantity += self.selectedItem.quantity
								self.selectedItem.quantity = 0
							} else {
								// Act on a single item
								self.selectedItem.quantity -= 1
								slot.quantity += 1
							}

						} else {
							// pick more up
							if bulk {
								// Act on all items in slot
								self.selectedItem.quantity += slot.quantity
								slot.quantity = 0
							} else {
								// Act on a single item
								self.selectedItem.quantity += 1
								slot.quantity -= 1
							}
						}
					} else if self.selectedItem.area != AREA_HANDHELD_COMPONENT {
						// Different slot
						if invert {
							if bulk {
								self.selectedItem.quantity += slot.quantity
								slot.quantity = 0
							} else {
								self.selectedItem.quantity += 1
								slot.quantity -= 1
							}
						} else {
							if bulk {
								slot.quantity += self.selectedItem.quantity
								self.selectedItem.quantity = 0
							} else {
								self.selectedItem.quantity -= 1
								slot.quantity += 1
							}
						}

					}
				} else if slot.item == 0 {
					// Check for duplicates
					for j := range self.componentSlots {
						if j != i && self.componentSlots[j].item == self.selectedItem.item {
							return
						}
					}

					quantity := uint16(1)
					if bulk {
						quantity = self.selectedItem.quantity
					}

					self.selectedItem.quantity -= quantity
					slot.quantity += quantity
					slot.item = self.selectedItem.item
				}

				if self.selectedItem != nil && self.selectedItem.quantity == 0 {
					self.selectedItem = nil
				}

				if slot.quantity == 0 {
					slot.item = 0
				}

				if slotQuantity != slot.quantity {
					self.UpdateProducts()
				}

				return

			}
		}
		for i := range self.productRects {
			if self.productRects[i].Contains(x, y) {
				if self.productSlots[i] != nil {
					recipe := self.productSlots[i]

					itemid := recipe.product.item
					if self.selectedItem == nil {
						if self.HasRecipeComponents(recipe) {
							for _, rc := range recipe.components {
								for k := range self.componentSlots {
									if self.componentSlots[k].item == rc.item {
										self.componentSlots[k].quantity -= rc.quantity
									}
								}
							}
							self.selectedItem = &SelectedItem{ItemQuantity{recipe.product.item, recipe.product.quantity}, ItemSlot{AREA_HANDHELD_PRODUCT, i}}
							self.UpdateProducts()

						}
					} else if itemid == self.selectedItem.item {
						if self.HasRecipeComponents(recipe) {
							for _, rc := range recipe.components {
								for k := range self.componentSlots {
									if self.componentSlots[k].item == rc.item {
										self.componentSlots[k].quantity -= rc.quantity
									}
									if self.componentSlots[k].quantity == 0 {
										self.componentSlots[k].item = 0
									}
								}
							}
							self.selectedItem.quantity += recipe.product.quantity
							if self.selectedItem.quantity > MAX_ITEMS_IN_INVENTORY {
								self.selectedItem.quantity = MAX_ITEMS_IN_INVENTORY
							}
							self.UpdateProducts()
						}
					}

					return
				}
			}
		}

		if self.currentContainer != nil {

			for i := range self.containerRects {
				if self.containerRects[i].Contains(x, y) {

					if self.selectedItem == nil {
						if self.currentContainer.CanTake() {
							// pick up item, if any
							if self.currentContainer.Item(i).item != 0 && self.currentContainer.Item(i).quantity > 0 {
								if bulk {
									// Act on all items in slot
									self.selectedItem = &SelectedItem{ItemQuantity{self.currentContainer.Item(i).item, self.currentContainer.Item(i).quantity}, ItemSlot{AREA_CONTAINER, i}}
									self.currentContainer.Take(i, self.currentContainer.Item(i).quantity)
								} else {
									// Act on a single item
									self.selectedItem = &SelectedItem{ItemQuantity{self.currentContainer.Item(i).item, 1}, ItemSlot{AREA_CONTAINER, i}}
									self.currentContainer.Take(i, 1)
								}
							}
						}
					} else if self.currentContainer.Item(i).item == self.selectedItem.item {
						// Clicked on slot containing same item as current selection
						// Is this the same slot as the source of items
						if self.selectedItem.area == AREA_CONTAINER && self.selectedItem.index == i {
							if self.currentContainer.CanTake() {
								// Same slot
								if invert {
									// drop back into same slot
									if bulk {
										// Act on all items in slot
										self.currentContainer.Place(i, &ItemQuantity{self.selectedItem.item, self.selectedItem.quantity})
										self.selectedItem.quantity = 0
									} else {
										// Act on a single item
										self.selectedItem.quantity -= 1
										self.currentContainer.Place(i, &ItemQuantity{self.selectedItem.item, 1})
									}

								} else {
									// pick more up
									if bulk {
										// Act on all items in slot
										self.selectedItem.quantity += self.currentContainer.Item(i).quantity
										self.currentContainer.Take(i, self.currentContainer.Item(i).quantity)
									} else {
										// Act on a single item
										self.selectedItem.quantity += 1
										self.currentContainer.Take(i, 1)
									}
								}
							}
						} else {
							// Different slot
							if invert {
								if self.currentContainer.CanTake() {
									if bulk {
										// Act on all items in slot
										self.selectedItem.quantity += self.currentContainer.Item(i).quantity
										self.currentContainer.Take(i, self.currentContainer.Item(i).quantity)
									} else {
										// Act on a single item
										self.selectedItem.quantity += 1
										self.currentContainer.Take(i, 1)
									}
								}
							} else {
								if self.currentContainer.CanPlace(self.selectedItem.item) {
									if bulk {
										// Act on all items in slot
										self.currentContainer.Place(i, &ItemQuantity{self.selectedItem.item, self.selectedItem.quantity})
										self.selectedItem.quantity = 0
									} else {
										// Act on a single item
										self.selectedItem.quantity -= 1
										self.currentContainer.Place(i, &ItemQuantity{self.selectedItem.item, 1})
									}
								}
							}
						}
					} else if self.currentContainer.Item(i).item == 0 {
						if self.currentContainer.CanPlace(self.selectedItem.item) {
							quantity := uint16(1)
							if bulk {
								quantity = self.selectedItem.quantity
							}

							self.selectedItem.quantity -= quantity
							self.currentContainer.Place(i, &ItemQuantity{self.selectedItem.item, quantity})
						}
					}

					if self.selectedItem != nil && self.selectedItem.quantity == 0 {
						self.selectedItem = nil
					}

					return
				}
			}
		}

		// Add/remove items from picker
		hit, pos := picker.HitTest(x, y)
		if hit && pos > 2 {
			if self.selectedItem == nil {
				// Remove from picker
				ThePlayer.equippedItems[uint16(pos)-3] = ITEM_NONE
			} else {
				// Set the picker
				ThePlayer.equippedItems[uint16(pos)-3] = self.selectedItem.item
				if self.selectedItem.area == AREA_INVENTORY {
					self.inventorySlots[self.selectedItem.index].quantity += self.selectedItem.quantity
				} else {
					println("Looking for somewhere in inventory to place some ", self.selectedItem.item)
					for i := range self.inventorySlots {
						if self.inventorySlots[i].item == self.selectedItem.item {
							println("Found some ", self.selectedItem.item, " in inventory slot ", i)
							if self.inventorySlots[i].quantity+self.selectedItem.quantity < MAX_ITEMS_IN_INVENTORY {
								self.inventorySlots[i].quantity += self.selectedItem.quantity
								break
							} else {
								self.selectedItem.quantity -= MAX_ITEMS_IN_INVENTORY - self.inventorySlots[i].quantity
								self.inventorySlots[i].quantity = MAX_ITEMS_IN_INVENTORY
							}
						} else if self.inventorySlots[i].item == 0 {
							self.inventorySlots[i].item = self.selectedItem.item
							self.inventorySlots[i].quantity = self.selectedItem.quantity
							break
						}
					}
				}
				self.selectedItem = nil

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

	for i, ir := range self.inventoryRects {
		if ir.Contains(x, y) {
			itemid = self.inventorySlots[i].item
			break
		}
	}
	if itemid == 0 {
		for i, cr := range self.componentRects {
			if cr.Contains(x, y) {
				itemid = self.componentSlots[i].item
				break
			}
		}
	}

	if itemid == 0 {
		for i, pr := range self.productRects {
			if pr.Contains(x, y) && self.productSlots[i] != nil {
				itemid = self.productSlots[i].product.item
				break
			}
		}
	}

	if itemid == 0 {

		if hit, pos := picker.HitTest(x, y); hit && pos > 2 {
			itemid = uint16(pos) - 3
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

func (self *Inventory) Hide() {
	self.currentContainer = nil
	self.currentCrafting = nil
	self.containerRects = nil

	for i := range ThePlayer.inventory {
		ThePlayer.inventory[i] = 0
	}

	for i := range self.inventorySlots {
		ThePlayer.inventory[self.inventorySlots[i].item] += self.inventorySlots[i].quantity
		self.inventorySlots[i].quantity = 0
		self.inventorySlots[i].item = 0
	}

	if self.selectedItem != nil {
		ThePlayer.inventory[self.selectedItem.item] += self.selectedItem.quantity
		self.selectedItem = nil
	}

	for i := range self.componentSlots {
		ThePlayer.inventory[self.componentSlots[i].item] += self.componentSlots[i].quantity
		self.componentSlots[i].quantity = 0
		self.componentSlots[i].item = 0
	}

	for i := range ThePlayer.inventory {
		if ThePlayer.inventory[i] > MAX_ITEMS_IN_INVENTORY {
			ThePlayer.inventory[i] = MAX_ITEMS_IN_INVENTORY
		}
	}

	self.visible = false

}

func (self *Inventory) Show(container ContainerObject, crafting CraftingObject) {
	self.currentContainer = container
	if self.currentContainer != nil {
		self.containerRects = make([]Rect, container.Slots())
	}

	self.currentCrafting = crafting

	self.selectedItem = nil

	slot := 0
	for i := range ThePlayer.inventory {
		if ThePlayer.inventory[i] > 0 {
			self.inventorySlots[slot].item = uint16(i)
			self.inventorySlots[slot].quantity = ThePlayer.inventory[i]
			slot++
		}
	}

	self.visible = true

}
