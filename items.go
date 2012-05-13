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
	// Items() []ItemQuantity
	// Add(item ItemQuantity)
	// Remove(item ItemQuantity)
	// Capacity(itemid uint16)
	Interact() // Just temporay for testing
}

type ItemQuantity struct {
	item     uint16
	quantity uint16
}

type Recipe struct {
	product    ItemQuantity
	components []ItemQuantity
}

func NewItem(id uint16, name string, texture1 uint16, texture2 uint16, hitsNeeded byte, transparent bool, collectable bool, placeable bool, drops *ItemQuantity) Item {
	return Item{id, name,
		//	[6]uint16{e, w, n, s, u, d},
		texture1,
		texture2,
		hitsNeeded,
		transparent,
		collectable,
		placeable,
		drops,
	}
}

func InitItems() {
	items = make(map[uint16]Item)
	items[BLOCK_AIR] = NewItem(BLOCK_AIR, "Air", TEXTURE_NONE, TEXTURE_NONE, STRENGTH_STONE, true, false, false, nil)
	items[BLOCK_STONE] = NewItem(BLOCK_STONE, "Stone", TEXTURE_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_STONE, 1})
	items[BLOCK_DIRT] = NewItem(BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_DIRT_TOP, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_BURNT_GRASS] = NewItem(BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_BURNT_GRASS, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_TRUNK] = NewItem(BLOCK_TRUNK, "Tree trunk", TEXTURE_TRUNK, TEXTURE_TRUNK, STRENGTH_WOOD, false, true, true, &ItemQuantity{BLOCK_TRUNK, 1})
	items[BLOCK_LEAVES] = NewItem(BLOCK_TRUNK, "Leaves", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, true, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_BUSH] = NewItem(BLOCK_BUSH, "Bush", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, true, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_WALL] = NewItem(BLOCK_LOG_WALL, "Log wall", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_SLAB] = NewItem(BLOCK_LOG_SLAB, "Log slab", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})

	items[BLOCK_COAL] = NewItem(BLOCK_COAL, "Coal", TEXTURE_COAL, TEXTURE_COAL, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_COAL, 1})
	items[BLOCK_IRON] = NewItem(BLOCK_IRON, "Iron Ore", TEXTURE_IRON, TEXTURE_IRON, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_IRON, 1})
	items[BLOCK_COPPER] = NewItem(BLOCK_COPPER, "Copper Ore", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_COPPER, 1})
	items[BLOCK_AMBERFELL] = NewItem(BLOCK_AMBERFELL, "Amberfell", TEXTURE_AMBERFELL, TEXTURE_AMBERFELL_TOP, STRENGTH_UNBREAKABLE, false, true, true, nil)
	items[BLOCK_CARVED_STONE] = NewItem(BLOCK_CARVED_STONE, "Carved stone", TEXTURE_CARVED_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_STONE, 1})
	items[BLOCK_CAMPFIRE] = NewItem(BLOCK_CAMPFIRE, "Campfire", TEXTURE_LOG_WALL, TEXTURE_FIRE, STRENGTH_LEAVES, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})

	items[ITEM_FIREWOOD] = NewItem(ITEM_FIREWOOD, "Firewood", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL, STRENGTH_WOOD, true, true, false, nil)

	items[BLOCK_AMBERFELL_PUMP] = NewItem(BLOCK_AMBERFELL_PUMP, "Amberfell Pump", TEXTURE_COPPER_MACH_SIDE, TEXTURE_COPPER_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_AMBERFELL_PUMP, 1})
	items[BLOCK_STEAM_GENERATOR] = NewItem(BLOCK_STEAM_GENERATOR, "Steam Generator", TEXTURE_IRON_MACH_SIDE, TEXTURE_IRON_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_STEAM_GENERATOR, 1})

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
	sourced        bool
	unitsPerSecond float64
	unitsHeld      float64
}

func NewAmberfellPump(sourced bool, powered bool) *AmberfellPump {
	pump := AmberfellPump{}
	pump.sourced = sourced
	if powered {
		pump.unitsPerSecond = AMBERFELL_UNITS_PER_SECOND_POWERED
	} else {
		pump.unitsPerSecond = AMBERFELL_UNITS_PER_SECOND_UNPOWERED
	}
	return &pump
}

func (self *AmberfellPump) Update(dt float64) (completed bool) {
	if self.sourced {
		self.unitsHeld += self.unitsPerSecond * dt
		if self.unitsHeld > AMBERFELL_PUMP_CAPACITY {
			self.unitsHeld = AMBERFELL_PUMP_CAPACITY
		}
	}

	return false
}

func (self *AmberfellPump) Interact() {
	println("Amberfell pump: ", self.unitsHeld)

}
