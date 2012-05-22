/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

const (
	AMBERFELL_UNITS_PER_SECOND_UNPOWERED = 0.01
	AMBERFELL_UNITS_PER_SECOND_POWERED   = 0.2
	AMBERFELL_PUMP_CAPACITY              = 5
)

type AmberfellPump struct {
	pos            Vectori
	sourced        bool
	unitsPerSecond float64
	unitsHeld      float64
}

func NewAmberfellPump(pos Vectori, sourced bool, powered bool) *AmberfellPump {
	pump := AmberfellPump{pos: pos, sourced: sourced}
	return &pump
}

func (self *AmberfellPump) Update(dt float64) (completed bool) {
	if self.sourced {
		self.unitsPerSecond = AMBERFELL_UNITS_PER_SECOND_UNPOWERED
		for face := 0; face < 6; face++ {
			npos := Vectori{self.pos[XAXIS], self.pos[YAXIS], self.pos[ZAXIS]}
			switch face {
			case NORTH_FACE:
				npos[ZAXIS]--
			case SOUTH_FACE:
				npos[ZAXIS]++
			case EAST_FACE:
				npos[XAXIS]++
			case WEST_FACE:
				npos[XAXIS]--
			case UP_FACE:
				npos[YAXIS]++
			case DOWN_FACE:
				npos[YAXIS]--
			}

			if gen, ok := TheWorld.generatorObjects[npos]; ok && gen.Active() {
				self.unitsPerSecond = AMBERFELL_UNITS_PER_SECOND_POWERED
				break
			}
		}

		self.unitsHeld += self.unitsPerSecond * dt
		if self.unitsHeld > AMBERFELL_PUMP_CAPACITY {
			self.unitsHeld = AMBERFELL_PUMP_CAPACITY
		}
	}

	return false
}

func (self *AmberfellPump) Label() string {
	if self.sourced {
		if self.unitsPerSecond == AMBERFELL_UNITS_PER_SECOND_POWERED {
			return "Amberfell Pump (powered)"
		} else {
			return "Amberfell Pump (unpowered)"

		}
	}

	return "Amberfell Pump (inactive)"
}

func (self *AmberfellPump) Slots() int {
	return 1
}

func (self *AmberfellPump) Item(slot int) ItemQuantity {
	if slot == 0 && self.unitsHeld >= 1 {
		return ItemQuantity{ITEM_AMBERFELL, uint16(self.unitsHeld)}
	}

	return ItemQuantity{}
}
func (self *AmberfellPump) Take(slot int, quantity uint16) {
	if slot == 0 {
		self.unitsHeld -= float64(quantity)
	}

	if self.unitsHeld < 0 {
		self.unitsHeld = 0
	}
}

func (self *AmberfellPump) Place(slot int, iq *ItemQuantity) {
	// NOOP
}

func (self *AmberfellPump) CanTake() bool {
	return true
}

func (self *AmberfellPump) CanPlace(itemid ItemId) bool {
	return false
}
