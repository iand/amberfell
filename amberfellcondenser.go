/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type AmberfellCondenser struct {
	pos           Vectori
	firewood      uint16
	amberfell     uint16
	crystals      uint16
	timeRemaining float64
}

func NewAmberfellCondenser(pos Vectori) *AmberfellCondenser {
	ac := AmberfellCondenser{pos: pos}
	return &ac
}

func (self *AmberfellCondenser) Update(dt float64) (completed bool) {
	if self.timeRemaining > 0 {
		self.timeRemaining -= 1 * dt
		if self.timeRemaining <= 0 {
			self.crystals++
		}
	}

	if self.timeRemaining <= 0 {
		if self.amberfell > 20 && self.firewood > 5 {
			self.amberfell -= 20
			self.firewood -= 5
			self.timeRemaining = 120
		} else {
			self.timeRemaining = 0
		}
	}
	return false
}

func (self *AmberfellCondenser) Label() string {
	if self.timeRemaining > 0 {
		return "Amberfell Condenser (active)"
	}

	return "Amberfell Condenser (inactive)"
}

func (self *AmberfellCondenser) Slots() int {
	return 3
}

func (self *AmberfellCondenser) Item(slot int) ItemQuantity {
	if slot == 0 && self.firewood > 0 {
		return ItemQuantity{ITEM_FIREWOOD, self.firewood}
	} else if slot == 1 && self.amberfell > 0 {
		return ItemQuantity{ITEM_AMBERFELL, self.amberfell}
	} else if slot == 2 && self.crystals > 0 {
		return ItemQuantity{ITEM_AMBERFELL_CRYSTAL, self.crystals}
	}

	return ItemQuantity{}
}

func (self *AmberfellCondenser) Take(slot int, quantity uint16) {
	switch slot {
	case 0:
		self.firewood -= quantity
		if self.firewood < 0 {
			self.firewood = 0
		}
	case 1:
		self.amberfell -= quantity
		if self.amberfell < 0 {
			self.amberfell = 0
		}
	case 2:
		self.crystals -= quantity
		if self.crystals < 0 {
			self.crystals = 0
		}
	}
}

func (self *AmberfellCondenser) Place(slot int, iq *ItemQuantity) {
	if (slot == 0 && self.firewood != 0 && iq.item != ITEM_FIREWOOD) ||
		(slot == 1 && self.amberfell != 0 && iq.item != ITEM_AMBERFELL) ||
		(slot == 2 && self.crystals != 0 && iq.item != ITEM_AMBERFELL_CRYSTAL) {
		return
	}

	switch iq.item {
	case ITEM_FIREWOOD:
		self.firewood += iq.quantity
	case ITEM_AMBERFELL:
		self.amberfell += iq.quantity
	case ITEM_AMBERFELL_CRYSTAL:
		self.crystals += iq.quantity
	}
}

func (self *AmberfellCondenser) CanTake() bool {
	return true
}

func (self *AmberfellCondenser) CanPlace(itemid uint16) bool {
	if itemid == ITEM_FIREWOOD || itemid == ITEM_AMBERFELL || itemid == ITEM_AMBERFELL_CRYSTAL {
		return true
	}
	return false
}
