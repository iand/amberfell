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
	shape       byte
	transparent bool
	collectable bool
	placeable   bool
	autojump    bool
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

func NewBlockType(id uint16, name string, texture1 uint16, texture2 uint16, hitsNeeded byte, shape byte, transparent bool, collectable bool, placeable bool, autojump bool, drops *ItemQuantity) Item {
	return Item{id: id, name: name,
		texture1:    texture1,
		texture2:    texture2,
		hitsNeeded:  hitsNeeded,
		shape:       shape,
		transparent: transparent,
		collectable: collectable,
		placeable:   placeable,
		autojump:    autojump,
		drops:       drops,
	}
}

func NewItemType(id uint16, name string, texture1 uint16) Item {
	return Item{id: id, name: name,
		texture1:    texture1,
		texture2:    TEXTURE_NONE,
		collectable: true,
	}
}

func InitItems() {
	items = make(map[uint16]Item)
	items[BLOCK_AIR] = NewBlockType(BLOCK_AIR, "Air", TEXTURE_NONE, TEXTURE_NONE, STRENGTH_STONE, SHAPE_CUBE, true, false, false, false, nil)
	items[BLOCK_STONE] = NewBlockType(BLOCK_STONE, "Stone", TEXTURE_STONE, TEXTURE_STONE, STRENGTH_STONE, SHAPE_CUBE, false, true, true, true, &ItemQuantity{ITEM_RUBBLE, 6})
	items[BLOCK_DIRT] = NewBlockType(BLOCK_DIRT, "Dirt", TEXTURE_DIRT, TEXTURE_DIRT_TOP, STRENGTH_DIRT, SHAPE_CUBE, false, true, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_BURNT_GRASS] = NewBlockType(BLOCK_BURNT_GRASS, "Dirt", TEXTURE_DIRT, TEXTURE_BURNT_GRASS, STRENGTH_DIRT, SHAPE_CUBE, false, true, true, true, &ItemQuantity{BLOCK_DIRT, 1})
	items[BLOCK_TRUNK] = NewBlockType(BLOCK_TRUNK, "Tree trunk", TEXTURE_TRUNK, TEXTURE_TRUNK, STRENGTH_WOOD, SHAPE_CUBE, false, true, true, true, &ItemQuantity{BLOCK_TRUNK, 1})
	items[BLOCK_LEAVES] = NewBlockType(BLOCK_LEAVES, "Leaves", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, SHAPE_CUBE, false, false, false, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_BUSH] = NewBlockType(BLOCK_BUSH, "Bush", TEXTURE_LEAVES, TEXTURE_LEAVES, STRENGTH_LEAVES, SHAPE_CUBE, false, false, false, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_WALL] = NewBlockType(BLOCK_LOG_WALL, "Log wall", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, SHAPE_WALL, true, true, true, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_LOG_SLAB] = NewBlockType(BLOCK_LOG_SLAB, "Log slab", TEXTURE_LOG_WALL, TEXTURE_LOG_WALL_TOP, STRENGTH_WOOD, SHAPE_SLAB, true, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 1})

	items[BLOCK_STONEBRICK_WALL] = NewBlockType(BLOCK_STONEBRICK_WALL, "Stone brick wall", TEXTURE_STONE_BRICK, TEXTURE_STONE_BRICK, STRENGTH_STONE, SHAPE_WALL, true, true, true, false, &ItemQuantity{ITEM_RUBBLE, 2})
	items[BLOCK_STONEBRICK_SLAB] = NewBlockType(BLOCK_STONEBRICK_SLAB, "Stone brick slab", TEXTURE_STONE_BRICK, TEXTURE_STONE_BRICK, STRENGTH_STONE, SHAPE_SLAB, true, true, true, true, &ItemQuantity{ITEM_RUBBLE, 2})

	items[BLOCK_PLANK_WALL] = NewBlockType(BLOCK_PLANK_WALL, "Wooden wall", TEXTURE_PLANK_WALL, TEXTURE_PLANK_WALL, STRENGTH_WOOD, SHAPE_WALL, true, true, true, false, &ItemQuantity{ITEM_PLANK, 1})
	items[BLOCK_PLANK_SLAB] = NewBlockType(BLOCK_PLANK_SLAB, "Wooden slab", TEXTURE_PLANK_WALL, TEXTURE_PLANK_WALL, STRENGTH_WOOD, SHAPE_SLAB, true, true, true, true, &ItemQuantity{ITEM_PLANK, 1})

	items[BLOCK_COAL_SEAM] = NewBlockType(BLOCK_COAL_SEAM, "Coal Seam", TEXTURE_COAL, TEXTURE_COAL, STRENGTH_STONE, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_COAL, 2})
	items[BLOCK_IRON_SEAM] = NewBlockType(BLOCK_IRON_SEAM, "Haematite Seam", TEXTURE_IRON, TEXTURE_IRON, STRENGTH_STONE, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_IRON_ORE, 2})
	items[BLOCK_COPPER_SEAM] = NewBlockType(BLOCK_COPPER_SEAM, "Malachite Seam", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_COPPER_ORE, 2})
	items[BLOCK_MAGNETITE_SEAM] = NewBlockType(BLOCK_MAGNETITE_SEAM, "Magnetite Seam", TEXTURE_COPPER, TEXTURE_COPPER, STRENGTH_STONE, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_MAGNETITE_ORE, 2})
	items[BLOCK_ZINC_SEAM] = NewBlockType(BLOCK_ZINC_SEAM, "Zinc Seam", TEXTURE_ZINC, TEXTURE_ZINC, STRENGTH_STONE, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_ZINC_ORE, 2})
	items[BLOCK_QUARTZ_SEAM] = NewBlockType(BLOCK_QUARTZ_SEAM, "Quartz Seam", TEXTURE_QUARTZ_SEAM, TEXTURE_QUARTZ_SEAM, STRENGTH_WOOD, SHAPE_CUBE, false, false, false, true, &ItemQuantity{ITEM_ZINC_ORE, 2})
	items[BLOCK_AMBERFELL_SOURCE] = NewBlockType(BLOCK_AMBERFELL_SOURCE, "Amberfell Source", TEXTURE_AMBERFELL_SOURCE, TEXTURE_AMBERFELL_SOURCE_TOP, STRENGTH_UNBREAKABLE, SHAPE_CUBE, false, true, true, true, nil)
	items[BLOCK_CARVED_STONE] = NewBlockType(BLOCK_CARVED_STONE, "Carved stone", TEXTURE_CARVED_STONE, TEXTURE_STONE, STRENGTH_STONE, SHAPE_CUBE, false, true, true, true, &ItemQuantity{BLOCK_STONE, 1})
	items[BLOCK_CAMPFIRE] = NewBlockType(BLOCK_CAMPFIRE, "Campfire", TEXTURE_LOG_WALL, TEXTURE_FIRE, STRENGTH_LEAVES, SHAPE_PILE, true, true, true, false, &ItemQuantity{ITEM_FIREWOOD, 1})
	items[BLOCK_BEESNEST] = NewBlockType(BLOCK_BEESNEST, "Bees Nest", TEXTURE_BEESNEST, TEXTURE_BEESNEST_TOP, STRENGTH_LEAVES, SHAPE_CUBE, true, true, true, false, &ItemQuantity{ITEM_BEESWAX, 2})

	items[BLOCK_AMBERFELL_PUMP] = NewBlockType(BLOCK_AMBERFELL_PUMP, "Amberfell Pump", TEXTURE_COPPER_MACH_SIDE, TEXTURE_COPPER_MACH_TOP, STRENGTH_IRON, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{BLOCK_AMBERFELL_PUMP, 1})
	items[BLOCK_STEAM_GENERATOR] = NewBlockType(BLOCK_STEAM_GENERATOR, "Steam Generator", TEXTURE_IRON_MACH_SIDE, TEXTURE_IRON_MACH_TOP, STRENGTH_IRON, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{BLOCK_STEAM_GENERATOR, 1})
	items[BLOCK_CARPENTERS_BENCH] = NewBlockType(BLOCK_CARPENTERS_BENCH, "Carpenter's Bench", TEXTURE_PLANK_WALL, TEXTURE_CARPENTERS_BENCH_TOP, STRENGTH_WOOD, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 5})
	items[BLOCK_AMBERFELL_CONDENSER] = NewBlockType(BLOCK_AMBERFELL_CONDENSER, "Amberfell Condenser", TEXTURE_PLANK_WALL, TEXTURE_CARPENTERS_BENCH_TOP, STRENGTH_WOOD, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{ITEM_FIREWOOD, 5})
	items[BLOCK_FURNACE] = NewBlockType(BLOCK_FURNACE, "Furnace", TEXTURE_STONE_BRICK, TEXTURE_FURNACE_TOP, STRENGTH_STONE, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{ITEM_RUBBLE, 4})
	items[BLOCK_FORGE] = NewBlockType(BLOCK_FORGE, "Forge", TEXTURE_STONE_BRICK, TEXTURE_FORGE_TOP, STRENGTH_STONE, SHAPE_ORIENTED_CUBE, false, true, true, true, &ItemQuantity{ITEM_RUBBLE, 4})

	items[ITEM_FIREWOOD] = NewItemType(ITEM_FIREWOOD, "Firewood", TEXTURE_ITEM_FIREWOOD)
	items[ITEM_RUBBLE] = NewItemType(ITEM_RUBBLE, "Rubble", TEXTURE_ITEM_RUBBLE)
	items[ITEM_STONE_BRICK] = NewItemType(ITEM_STONE_BRICK, "Stone Brick", TEXTURE_ITEM_STONE_BRICK)
	items[ITEM_PLANK] = NewItemType(ITEM_PLANK, "Plank", TEXTURE_ITEM_PLANK)
	items[ITEM_COAL] = NewItemType(ITEM_COAL, "Coal", TEXTURE_ITEM_COAL)
	items[ITEM_IRON_ORE] = NewItemType(ITEM_IRON_ORE, "Haematite", TEXTURE_ITEM_IRON_ORE)
	items[ITEM_MAGNETITE_ORE] = NewItemType(ITEM_MAGNETITE_ORE, "Magnetite", TEXTURE_ITEM_IRON_ORE)
	items[ITEM_COPPER_ORE] = NewItemType(ITEM_COPPER_ORE, "Malachite", TEXTURE_ITEM_COPPER_ORE)
	items[ITEM_ZINC_ORE] = NewItemType(ITEM_ZINC_ORE, "Sphalerite", TEXTURE_ITEM_ZINC_ORE)
	items[ITEM_AMBERFELL] = NewItemType(ITEM_AMBERFELL, "Amberfell", TEXTURE_ITEM_AMBERFELL)
	items[ITEM_AMBERFELL_CRYSTAL] = NewItemType(ITEM_AMBERFELL_CRYSTAL, "Amberfell Crystal", TEXTURE_ITEM_AMBERFELL_CRYSTAL)

	items[ITEM_IRON_INGOT] = NewItemType(ITEM_IRON_INGOT, "Iron Ingot", TEXTURE_ITEM_IRON_INGOT)
	items[ITEM_COPPER_INGOT] = NewItemType(ITEM_COPPER_INGOT, "Copper Ingot", TEXTURE_ITEM_COPPER_INGOT)
	items[ITEM_LODESTONE] = NewItemType(ITEM_LODESTONE, "Lodestone", TEXTURE_ITEM_LODESTONE)

	items[ITEM_BRASS_INGOT] = NewItemType(ITEM_BRASS_INGOT, "Brass Ingot", TEXTURE_ITEM_BRASS_INGOT)
	items[ITEM_COPPER_PLATE] = NewItemType(ITEM_COPPER_PLATE, "Copper Plate", TEXTURE_ITEM_COPPER_PLATE)
	items[ITEM_BRASS_PLATE] = NewItemType(ITEM_BRASS_PLATE, "Brass Plate", TEXTURE_ITEM_BRASS_PLATE)
	items[ITEM_IRON_PLATE] = NewItemType(ITEM_IRON_PLATE, "Iron Plate", TEXTURE_ITEM_IRON_PLATE)

	items[ITEM_LEATHER] = NewItemType(ITEM_LEATHER, "Leather", TEXTURE_ITEM_LEATHER)
	items[ITEM_BEESWAX] = NewItemType(ITEM_BEESWAX, "Beeswax", TEXTURE_ITEM_BEESWAX)

	items[ITEM_QUARTZ] = NewItemType(ITEM_QUARTZ, "Quartz", TEXTURE_ITEM_QUARTZ)
	items[ITEM_GLASS] = NewItemType(ITEM_GLASS, "Glass", TEXTURE_ITEM_GLASS)

}
