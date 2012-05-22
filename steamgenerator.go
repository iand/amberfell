/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type SteamGenerator struct {
	pos  Vectori
	fuel uint16
	life float64
}

func NewSteamGenerator(pos Vectori) *SteamGenerator {
	gen := SteamGenerator{pos: pos}
	return &gen
}

func (self *SteamGenerator) Update(dt float64) (completed bool) {
	self.life -= 0.02 * dt
	if self.life <= 0 {
		if self.fuel > 0 {
			self.fuel--
			self.life = 1
		} else {
			self.life = 0
		}
	}

	return false
}

func (self *SteamGenerator) Label() string {
	if self.life > 0 {
		return "Steam Generator (active)"
	}

	return "Steam Generator (inactive)"
}

func (self *SteamGenerator) Slots() int {
	return 1
}

func (self *SteamGenerator) Item(slot int) ItemQuantity {
	if slot == 0 && self.fuel >= 1 {
		return ItemQuantity{ITEM_COAL, uint16(self.fuel)}
	}

	return ItemQuantity{}
}

func (self *SteamGenerator) Take(slot int, quantity uint16) {
	if slot == 0 {
		self.fuel -= quantity
	}

	if self.fuel < 0 {
		self.fuel = 0
	}
}

func (self *SteamGenerator) Place(slot int, iq *ItemQuantity) {
	if slot == 0 && iq.item == ITEM_COAL {
		self.fuel += iq.quantity
	}
}

func (self *SteamGenerator) Active() bool {
	return self.life > 0
}

func (self *SteamGenerator) CanTake() bool {
	return true
}

func (self *SteamGenerator) CanPlace(itemid ItemId) bool {
	if itemid == ITEM_COAL {
		return true
	}
	return false
}
