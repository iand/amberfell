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
