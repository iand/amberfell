/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

const (
	BEESNEST_DENSITY_PCT      = 0.1
	BEESNEST_UNITS_PER_SECOND = 0.05
	BEESNEST_CAPACITY         = 4
)

type BeesNest struct {
	pos       Vectori
	unitsHeld float64
}

func NewBeesNest(pos Vectori) *BeesNest {
	nest := BeesNest{pos: pos}
	return &nest
}

func (self *BeesNest) Update(dt float64) (completed bool) {
	self.unitsHeld += BEESNEST_UNITS_PER_SECOND * dt
	if self.unitsHeld > BEESNEST_CAPACITY {
		self.unitsHeld = BEESNEST_CAPACITY
	}

	return false
}

func (self *BeesNest) Label() string {
	return items[BLOCK_BEESNEST].name
}

func (self *BeesNest) Slots() int {
	return 1
}

func (self *BeesNest) Item(slot int) ItemQuantity {
	if slot == 0 && self.unitsHeld >= 1 {
		return ItemQuantity{ITEM_BEESWAX, uint16(self.unitsHeld)}
	}

	return ItemQuantity{}
}
func (self *BeesNest) Take(slot int, quantity uint16) {
	if slot == 0 {
		self.unitsHeld -= float64(quantity)
	}

	if self.unitsHeld < 0 {
		self.unitsHeld = 0
	}
}

func (self *BeesNest) Place(slot int, iq *ItemQuantity) {
	// NOOP
}

func (self *BeesNest) CanTake() bool {
	return true
}

func (self *BeesNest) CanPlace(itemid ItemId) bool {
	return false
}
