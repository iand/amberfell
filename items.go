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
	Item(slot int) ItemQuantity
	CanTake() bool
	CanPlace(itemid uint16) bool
	Take(slot int, quantity uint16)
	Place(slot int, iq *ItemQuantity)
	Slots() int
}

type CraftingObject interface {
	Label() string
	Recipes() []Recipe
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

func NewItem(id uint16, name string, texture1 uint16, texture2 uint16, hitsNeeded byte, transparent bool, collectable bool, placeable bool, drops *ItemQuantity) Item {
	return Item{id: id, name: name,
		texture1:    texture1,
		texture2:    texture2,
		hitsNeeded:  hitsNeeded,
		transparent: transparent,
		collectable: collectable,
		placeable:   placeable,
		drops:       drops,
	}
}

func InitItems() {
	items = make(map[uint16]Item)
	items[BLOCK_AIR] = NewItem(BLOCK_AIR, "Air", TEXTURE_NONE, TEXTURE_NONE, STRENGTH_STONE, true, false, false, nil)
	items[BLOCK_STONE] = NewItem(BLOCK_STONE, "Stone", TEXTURE_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{ITEM_RUBBLE, 6})
	items[BLOCK_DIRT] = NewItem(BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_DIRT_TOP, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_BURNT_GRASS] = NewItem(BLOCK_BURNT_GRASS, "Dirt", TEXTURE_DIRT, TEXTURE_BURNT_GRASS, STRENGTH_DIRT, false, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_TRUNK] = NewItem(BLOCK_TRUNK, "Tree trunk", TEXTURE_TRUNK, TEXTURE_TRUNK, STRENGTH_WOOD, false, true, true, &ItemQuantity{BLOCK_TRUNK, 1})
	items[BLOCK_LEAVES] = NewItem(BLOCK_LEAVES, "Leaves", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, false, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_BUSH] = NewItem(BLOCK_BUSH, "Bush", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, false, false, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_WALL] = NewItem(BLOCK_LOG_WALL, "Log wall", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_SLAB] = NewItem(BLOCK_LOG_SLAB, "Log slab", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})

	items[BLOCK_STONEBRICK_WALL] = NewItem(BLOCK_STONEBRICK_WALL, "Stone brick wall", TEXTURE_STONE_BRICK, TEXTURE_STONE_BRICK, STRENGTH_STONE, true, true, true, &ItemQuantity{ITEM_RUBBLE, 2})
	items[BLOCK_STONEBRICK_SLAB] = NewItem(BLOCK_STONEBRICK_SLAB, "Stone brick slab", TEXTURE_STONE_BRICK, TEXTURE_STONE_BRICK, STRENGTH_STONE, true, true, true, &ItemQuantity{ITEM_RUBBLE, 2})

	items[BLOCK_PLANK_WALL] = NewItem(BLOCK_PLANK_WALL, "Wooden wall", TEXTURE_PLANK_WALL, TEXTURE_PLANK_WALL, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_PLANK, 1})
	items[BLOCK_PLANK_SLAB] = NewItem(BLOCK_PLANK_SLAB, "Wooden slab", TEXTURE_PLANK_WALL, TEXTURE_PLANK_WALL, STRENGTH_WOOD, true, true, true, &ItemQuantity{ITEM_PLANK, 1})

	items[BLOCK_COAL_SEAM] = NewItem(BLOCK_COAL_SEAM, "Coal Seam", TEXTURE_COAL, TEXTURE_COAL, STRENGTH_STONE, false, false, false, &ItemQuantity{ITEM_COAL, 2})
	items[BLOCK_IRON_SEAM] = NewItem(BLOCK_IRON_SEAM, "Haematite Seam", TEXTURE_IRON, TEXTURE_IRON, STRENGTH_STONE, false, false, false, &ItemQuantity{ITEM_IRON_ORE, 2})
	items[BLOCK_COPPER_SEAM] = NewItem(BLOCK_COPPER_SEAM, "Malachite Seam", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, false, false, false, &ItemQuantity{ITEM_COPPER_ORE, 2})
	items[BLOCK_MAGNETITE_SEAM] = NewItem(BLOCK_MAGNETITE_SEAM, "Magnetite Seam", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, false, false, false, &ItemQuantity{ITEM_MAGNETITE_ORE, 2})
	items[BLOCK_ZINC_SEAM] = NewItem(BLOCK_ZINC_SEAM, "Zinc Seam", TEXTURE_ZINC, TEXTURE_ZINC, STRENGTH_STONE, false, false, false, &ItemQuantity{ITEM_ZINC_ORE, 2})
	items[BLOCK_QUARTZ_SEAM] = NewItem(BLOCK_QUARTZ_SEAM, "Quartz Seam", TEXTURE_QUARTZ_SEAM, TEXTURE_QUARTZ_SEAM, STRENGTH_WOOD, false, false, false, &ItemQuantity{ITEM_ZINC_ORE, 2})
	items[BLOCK_AMBERFELL_SOURCE] = NewItem(BLOCK_AMBERFELL_SOURCE, "Amberfell Source", TEXTURE_AMBERFELL_SOURCE, TEXTURE_AMBERFELL_SOURCE_TOP, STRENGTH_UNBREAKABLE, false, true, true, nil)
	items[BLOCK_CARVED_STONE] = NewItem(BLOCK_CARVED_STONE, "Carved stone", TEXTURE_CARVED_STONE, TEXTURE_STONE, STRENGTH_STONE, false, true, true, &ItemQuantity{BLOCK_STONE, 1})
	items[BLOCK_CAMPFIRE] = NewItem(BLOCK_CAMPFIRE, "Campfire", TEXTURE_LOG_WALL, TEXTURE_FIRE, STRENGTH_LEAVES, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_BEESNEST] = NewItem(BLOCK_BEESNEST, "Bees Nest", TEXTURE_BEESNEST, TEXTURE_BEESNEST_TOP, STRENGTH_LEAVES, true, true, true, &ItemQuantity{ITEM_BEESWAX, 2})

	items[BLOCK_AMBERFELL_PUMP] = NewItem(BLOCK_AMBERFELL_PUMP, "Amberfell Pump", TEXTURE_COPPER_MACH_SIDE, TEXTURE_COPPER_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_AMBERFELL_PUMP, 1})
	items[BLOCK_STEAM_GENERATOR] = NewItem(BLOCK_STEAM_GENERATOR, "Steam Generator", TEXTURE_IRON_MACH_SIDE, TEXTURE_IRON_MACH_TOP, STRENGTH_IRON, false, true, true, &ItemQuantity{BLOCK_STEAM_GENERATOR, 1})
	items[BLOCK_CARPENTERS_BENCH] = NewItem(BLOCK_CARPENTERS_BENCH, "Carpenter's Bench", TEXTURE_PLANK_WALL, TEXTURE_CARPENTERS_BENCH_TOP, STRENGTH_WOOD, false, true, true, &ItemQuantity{ITEM_FIREWOOD, 5})
	items[BLOCK_AMBERFELL_CONDENSER] = NewItem(BLOCK_AMBERFELL_CONDENSER, "Amberfell Condenser", TEXTURE_PLANK_WALL, TEXTURE_CARPENTERS_BENCH_TOP, STRENGTH_WOOD, false, true, true, &ItemQuantity{ITEM_FIREWOOD, 5})
	items[BLOCK_FURNACE] = NewItem(BLOCK_FURNACE, "Furnace", TEXTURE_STONE_BRICK, TEXTURE_FURNACE_TOP, STRENGTH_STONE, false, true, true, &ItemQuantity{ITEM_RUBBLE, 4})
	items[BLOCK_FORGE] = NewItem(BLOCK_FORGE, "Forge", TEXTURE_STONE_BRICK, TEXTURE_FORGE_TOP, STRENGTH_STONE, false, true, true, &ItemQuantity{ITEM_RUBBLE, 4})

	items[ITEM_FIREWOOD] = NewItem(ITEM_FIREWOOD, "Firewood", TEXTURE_ITEM_FIREWOOD, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_RUBBLE] = NewItem(ITEM_RUBBLE, "Rubble", TEXTURE_ITEM_RUBBLE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_STONE_BRICK] = NewItem(ITEM_STONE_BRICK, "Stone Brick", TEXTURE_ITEM_STONE_BRICK, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_PLANK] = NewItem(ITEM_PLANK, "Plank", TEXTURE_ITEM_PLANK, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_COAL] = NewItem(ITEM_COAL, "Coal", TEXTURE_ITEM_COAL, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_IRON_ORE] = NewItem(ITEM_IRON_ORE, "Haematite", TEXTURE_ITEM_IRON_ORE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_MAGNETITE_ORE] = NewItem(ITEM_MAGNETITE_ORE, "Magnetite", TEXTURE_ITEM_IRON_ORE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_COPPER_ORE] = NewItem(ITEM_COPPER_ORE, "Malachite", TEXTURE_ITEM_COPPER_ORE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_ZINC_ORE] = NewItem(ITEM_ZINC_ORE, "Sphalerite", TEXTURE_ITEM_ZINC_ORE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_AMBERFELL] = NewItem(ITEM_AMBERFELL, "Amberfell", TEXTURE_ITEM_AMBERFELL, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_AMBERFELL_CRYSTAL] = NewItem(ITEM_AMBERFELL_CRYSTAL, "Amberfell Crystal", TEXTURE_ITEM_AMBERFELL_CRYSTAL, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)

	items[ITEM_IRON_INGOT] = NewItem(ITEM_IRON_INGOT, "Iron Ingot", TEXTURE_ITEM_IRON_INGOT, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_COPPER_INGOT] = NewItem(ITEM_COPPER_INGOT, "Copper Ingot", TEXTURE_ITEM_COPPER_INGOT, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_LODESTONE] = NewItem(ITEM_LODESTONE, "Lodestone", TEXTURE_ITEM_LODESTONE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)

	items[ITEM_BRASS_INGOT] = NewItem(ITEM_BRASS_INGOT, "Brass Ingot", TEXTURE_ITEM_BRASS_INGOT, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_COPPER_PLATE] = NewItem(ITEM_COPPER_PLATE, "Copper Plate", TEXTURE_ITEM_COPPER_PLATE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_BRASS_PLATE] = NewItem(ITEM_BRASS_PLATE, "Brass Plate", TEXTURE_ITEM_BRASS_PLATE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_IRON_PLATE] = NewItem(ITEM_IRON_PLATE, "Iron Plate", TEXTURE_ITEM_IRON_PLATE, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)

	items[ITEM_LEATHER] = NewItem(ITEM_LEATHER, "Leather", TEXTURE_ITEM_LEATHER, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_BEESWAX] = NewItem(ITEM_BEESWAX, "Beeswax", TEXTURE_ITEM_BEESWAX, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)

	items[ITEM_QUARTZ] = NewItem(ITEM_QUARTZ, "Quartz", TEXTURE_ITEM_QUARTZ, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)
	items[ITEM_GLASS] = NewItem(ITEM_GLASS, "Glass", TEXTURE_ITEM_GLASS, TEXTURE_NONE, STRENGTH_WOOD, false, true, false, nil)

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

func (self *CampFire) TargetType() uint8 {
	return TARGET_CAMPFIRE
}

func (self *CampFire) Position() Vectorf {
	return self.pos.Vectorf()
}

func (self *CampFire) Velocity() Vectorf {
	return Vectorf{}
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

func (self *AmberfellPump) CanPlace(itemid uint16) bool {
	return false
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

func (self *SteamGenerator) CanPlace(itemid uint16) bool {
	if itemid == ITEM_COAL {
		return true
	}
	return false
}

type CarpentersBench struct {
	pos Vectori
}

func NewCarpentersBench(pos Vectori) *CarpentersBench {
	obj := CarpentersBench{pos: pos}
	return &obj
}

func (self *CarpentersBench) Label() string {
	return "Carpenter's Bench"
}

func (self *CarpentersBench) Recipes() []Recipe {
	return carpenterRecipes
}

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

type Forge struct {
	pos Vectori
}

func NewForge(pos Vectori) *Forge {
	obj := Forge{pos: pos}
	return &obj
}

func (self *Forge) Label() string {
	return "Forge"
}

func (self *Forge) Recipes() []Recipe {
	return forgeRecipes
}

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

func (self *Furnace) CanPlace(itemid uint16) bool {

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

func (self *BeesNest) CanPlace(itemid uint16) bool {
	return false
}
