/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type Furnace struct {
	pos           Vectori
	coal          uint16
	ore           *ItemQuantity
	metal         *ItemQuantity
	productUnit   *ItemQuantity
	timeRemaining float64
}

func NewFurnace(pos Vectori) *Furnace {
	f := Furnace{pos: pos}
	f.ore = &ItemQuantity{}
	f.metal = &ItemQuantity{}
	return &f
}

func (self *Furnace) Update(dt float64) (completed bool) {
	if self.timeRemaining > 0 {
		self.timeRemaining -= 1 * dt
		if self.timeRemaining <= 0 && self.productUnit != nil {
			if self.metal.quantity == 0 {
				self.metal.item = self.productUnit.item
			}
			self.metal.quantity += self.productUnit.quantity
		}
	}

	if self.timeRemaining <= 0 {
		self.timeRemaining = 0
		if self.ore.quantity > 0 && self.coal > 0 {
			// Could produce some more metal

			oreMatchesProduct := false
			if self.metal.quantity > 0 {
				// Does ore match the metal currently in the furnace?
				for _, recipe := range furnaceRecipes {
					if recipe.product.item == self.metal.item {
						for _, component := range recipe.components {
							if self.ore.item == component.item {
								oreMatchesProduct = true
							}
						}
					}
				}
			} else {
				oreMatchesProduct = true
			}

			if oreMatchesProduct {
				// get started
				self.ore.quantity -= 1
				self.coal -= 1
				self.timeRemaining = 15
			}
		}
	}
	return false
}

func (self *Furnace) Label() string {
	if self.timeRemaining > 0 {
		return "Furnace (active)"
	}

	return "Furnace (inactive)"
}

func (self *Furnace) Slots() int {
	return 3
}

func (self *Furnace) Item(slot int) ItemQuantity {
	if slot == 0 && self.coal > 0 {
		return ItemQuantity{ITEM_COAL, self.coal}
	} else if slot == 1 && self.ore.quantity > 0 {
		return *self.ore
	} else if slot == 2 && self.metal.quantity > 0 {
		return *self.metal
	}

	return ItemQuantity{}
}

func (self *Furnace) Take(slot int, quantity uint16) {
	switch slot {
	case 0:
		self.coal -= quantity
		if self.coal < 0 {
			self.coal = 0
		}
	case 1:
		self.ore.quantity -= quantity
		if self.ore.quantity <= 0 {
			self.ore.quantity = 0
			self.productUnit = nil

		}
	case 2:
		self.metal.quantity -= quantity
		if self.metal.quantity < 0 {
			self.metal.quantity = 0
		}
	}
}

func (self *Furnace) Place(slot int, iq *ItemQuantity) {
	if (slot == 0 && self.coal != 0 && iq.item != ITEM_COAL) ||
		(slot == 1 && self.ore.quantity != 0 && iq.item != self.ore.item) ||
		(slot == 2 && self.metal.quantity != 0 && iq.item != self.metal.item) {
		return
	}

	// TODO bounds check
	if iq.item == ITEM_COAL {
		self.coal += iq.quantity
	} else {
		for _, recipe := range furnaceRecipes {
			if iq.item == recipe.product.item {
				self.metal.quantity += iq.quantity
				self.metal.item = iq.item
				self.productUnit = &ItemQuantity{item: recipe.product.item, quantity: recipe.product.quantity}
				return
			}

			for _, component := range recipe.components {
				if iq.item == component.item {
					self.ore.quantity += iq.quantity
					self.ore.item = iq.item
					self.productUnit = &ItemQuantity{item: recipe.product.item, quantity: recipe.product.quantity}
					return
				}
			}

		}
	}
}

func (self *Furnace) CanTake() bool {
	return true
}

func (self *Furnace) CanPlace(itemid ItemId) bool {

	if itemid == ITEM_COAL {
		return true
	} else {
		for _, recipe := range furnaceRecipes {
			if itemid == recipe.product.item {
				return true
			}

			for _, component := range recipe.components {
				if itemid == component.item {
					return true
				}
			}

		}
	}

	return false

}
