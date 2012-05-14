/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type Item struct {
	id   uint16
	name string
	//textures    [6]uint16
	texture1    uint16
	texture2    uint16
	hitsNeeded  byte
	transparent bool
	collectable bool
	placeable   bool
	drops       *ItemQuantity
}

type LightSource interface {
	Intensity() uint16
}
type TimedObject interface {
	Update(dt float64) bool
}

type ContainerObject interface {
	Label() string
	Item(slot uint16) *ItemQuantity
	// Add(item ItemQuantity)
	// Remove(item ItemQuantity)
	Slots() uint16
}

type GeneratorObject interface {
	Active() bool
}

type ItemQuantity struct {
	item     uint16
	quantity uint16
}

type Recipe struct {
	product    ItemQuantity
	components []ItemQuantity
}

func InitItems() {
	items = make(map[uint16]Item)
	items[BLOCK_AIR] = Item{BLOCK_AIR, "Air", TEXTURE_NONE, TEXTURE_NONE, STRENGTH_STONE, true, false, false, nil}
	items[BLOCK_STONE] = Item{BLOCK_STONE, "Stone", TEXTURE_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_STONE, 1}}
	items[BLOCK_DIRT] = Item{BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_DIRT_TOP, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1}}
	items[BLOCK_BURNT_GRASS] = Item{BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_BURNT_GRASS, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1}}
	items[BLOCK_TRUNK] = Item{BLOCK_TRUNK, "Tree trunk", TEXTURE_TRUNK, TEXTURE_TRUNK, STRENGTH_WOOD, false, true, true, &ItemQuantity{BLOCK_TRUNK, 1}}
	items[BLOCK_LEAVES] = Item{BLOCK_TRUNK, "Leaves", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, true, false, &ItemQuantity{ITEM_FIREWOOD, 1}}
	items[BLOCK_BUSH] = Item{BLOCK_BUSH, "Bush", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, true, false, &ItemQuantity{ITEM_FIREWOOD, 1}}
	items[BLOCK_LOG_WALL] = Item{BLOCK_LOG_WALL, "Log wall", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1}}
	items[BLOCK_LOG_SLAB] = Item{BLOCK_LOG_SLAB, "Log slab", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1}}

	items[BLOCK_COAL] = Item{BLOCK_COAL, "Coal", TEXTURE_COAL, TEXTURE_COAL, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_COAL, 1}}
	items[BLOCK_IRON] = Item{BLOCK_IRON, "Iron Ore", TEXTURE_IRON, TEXTURE_IRON, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_IRON, 1}}
	items[BLOCK_COPPER] = Item{BLOCK_COPPER, "Copper Ore", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_COPPER, 1}}
	items[BLOCK_AMBERFELL] = Item{BLOCK_AMBERFELL, "Amberfell", TEXTURE_AMBERFELL, TEXTURE_AMBERFELL_TOP, STRENGTH_UNBREAKABLE, false, true, true, nil}
	items[BLOCK_CARVED_STONE] = Item{BLOCK_CARVED_STONE, "Carved stone", TEXTURE_CARVED_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_STONE, 1}}
	items[BLOCK_CAMPFIRE] = Item{BLOCK_CAMPFIRE, "Campfire", TEXTURE_LOG_WALL, TEXTURE_FIRE, STRENGTH_LEAVES, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1}}

	items[ITEM_FIREWOOD] = Item{ITEM_FIREWOOD, "Firewood", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, STRENGTH_WOOD, true, true, false, nil}

	items[BLOCK_AMBERFELL_PUMP] = Item{BLOCK_AMBERFELL_PUMP, "Amberfell Pump", TEXTURE_COPPER_MACH_SIDE, TEXTURE_COPPER_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_AMBERFELL_PUMP, 1}}
	items[BLOCK_STEAM_GENERATOR] = Item{BLOCK_STEAM_GENERATOR, "Steam Generator", TEXTURE_IRON_MACH_SIDE, TEXTURE_IRON_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_STEAM_GENERATOR, 1}}

}

var handmadeRecipes = []Recipe{
	{product: ItemQuantity{BLOCK_LOG_WALL, 1},
		components: []ItemQuantity{
			{BLOCK_TRUNK, 1},
		}},

	{product: ItemQuantity{BLOCK_LOG_SLAB, 1},
		components: []ItemQuantity{
			{BLOCK_TRUNK, 1},
		}},

	{product: ItemQuantity{ITEM_FIREWOOD, 2},
		components: []ItemQuantity{
			{BLOCK_TRUNK, 1},
		}},

	{product: ItemQuantity{ITEM_FIREWOOD, 2},
		components: []ItemQuantity{
			{BLOCK_LOG_WALL, 1},
		}},

	{product: ItemQuantity{ITEM_FIREWOOD, 2},
		components: []ItemQuantity{
			{BLOCK_LOG_SLAB, 1},
		}},

	{product: ItemQuantity{BLOCK_CAMPFIRE, 1},
		components: []ItemQuantity{
			{ITEM_FIREWOOD, 3},
		}},
}

type CampFire struct {
	pos  Vectori
	life float64
}

func NewCampFire(pos Vectori) *CampFire {
	return &CampFire{pos: pos, life: CAMPFIRE_DURATION}
}

func (self *CampFire) Intensity() uint16 {
	return CAMPFIRE_INTENSITY
}

func (self *CampFire) Update(dt float64) (completed bool) {
	completed = false
	self.life -= 0.02 * dt
	if self.life <= 0 {
		TheWorld.Setv(self.pos, BLOCK_AIR)
		delete(TheWorld.lightSources, self.pos)
		TheWorld.InvalidateRadius(self.pos[XAXIS], self.pos[ZAXIS], CAMPFIRE_INTENSITY)
		completed = true
	}

	return completed
}

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

func (self *AmberfellPump) Slots() uint16 {
	return 1
}

func (self *AmberfellPump) Item(slot uint16) *ItemQuantity {
	if slot == 0 && self.unitsHeld > 1 {
		return &ItemQuantity{BLOCK_AMBERFELL, uint16(self.unitsHeld)}
	}

	return nil
}

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

func (self *SteamGenerator) Slots() uint16 {
	return 1
}

func (self *SteamGenerator) Item(slot uint16) *ItemQuantity {
	if slot == 0 && self.fuel > 1 {
		return &ItemQuantity{BLOCK_COAL, uint16(self.fuel)}
	}

	return nil
}

func (self *SteamGenerator) Active() bool {
	return self.life > 0
}
