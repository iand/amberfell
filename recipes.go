/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

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

	{product: ItemQuantity{ITEM_STONE_BRICK, 1},
		components: []ItemQuantity{
			{ITEM_RUBBLE, 2},
		}},

	{product: ItemQuantity{BLOCK_STONEBRICK_WALL, 1},
		components: []ItemQuantity{
			{ITEM_STONE_BRICK, 3},
		}},

	{product: ItemQuantity{BLOCK_STONEBRICK_SLAB, 1},
		components: []ItemQuantity{
			{ITEM_STONE_BRICK, 3},
		}},

	{product: ItemQuantity{ITEM_RUBBLE, 1},
		components: []ItemQuantity{
			{ITEM_STONE_BRICK, 1},
		}},
	{product: ItemQuantity{ITEM_RUBBLE, 2},
		components: []ItemQuantity{
			{BLOCK_STONEBRICK_WALL, 1},
		}},

	{product: ItemQuantity{ITEM_RUBBLE, 2},
		components: []ItemQuantity{
			{BLOCK_STONEBRICK_SLAB, 1},
		}},

	{product: ItemQuantity{BLOCK_FORGE, 1},
		components: []ItemQuantity{
			{ITEM_STONE_BRICK, 6},
		}},

	{product: ItemQuantity{BLOCK_FURNACE, 1},
		components: []ItemQuantity{
			{ITEM_STONE_BRICK, 6},
		}},

	{product: ItemQuantity{ITEM_RUBBLE, 8},
		components: []ItemQuantity{
			{BLOCK_FORGE, 1},
		}},

	{product: ItemQuantity{ITEM_RUBBLE, 8},
		components: []ItemQuantity{
			{BLOCK_FURNACE, 1},
		}},
}

var carpenterRecipes = []Recipe{
	{product: ItemQuantity{ITEM_PLANK, 4},
		components: []ItemQuantity{
			{BLOCK_TRUNK, 1},
		}},
	{product: ItemQuantity{BLOCK_PLANK_WALL, 1},
		components: []ItemQuantity{
			{ITEM_PLANK, 3},
		}},
	{product: ItemQuantity{BLOCK_PLANK_SLAB, 1},
		components: []ItemQuantity{
			{ITEM_PLANK, 3},
		}},

	{product: ItemQuantity{BLOCK_AMBERFELL_PUMP, 1},
		components: []ItemQuantity{
			{ITEM_PLANK, 2},
			{ITEM_BEESWAX, 1},
			{ITEM_LEATHER, 2},
			{ITEM_BRASS_PLATE, 4},
		}},

	{product: ItemQuantity{BLOCK_STEAM_GENERATOR, 1},
		components: []ItemQuantity{
			{ITEM_IRON_PLATE, 4},
			{ITEM_COPPER_PLATE, 2},
			{ITEM_STONE_BRICK, 8},
		}},

	{product: ItemQuantity{BLOCK_AMBERFELL_CONDENSER, 1},
		components: []ItemQuantity{
			{ITEM_PLANK, 2},
			{ITEM_GLASS, 6},
		}},
}

var forgeRecipes = []Recipe{

	{product: ItemQuantity{ITEM_BRASS_INGOT, 4},
		components: []ItemQuantity{
			{ITEM_COPPER_INGOT, 3},
			{ITEM_ZINC_INGOT, 1},
		}},

	{product: ItemQuantity{ITEM_COPPER_PLATE, 1},
		components: []ItemQuantity{
			{ITEM_COPPER_INGOT, 2},
		}},

	{product: ItemQuantity{ITEM_BRASS_PLATE, 1},
		components: []ItemQuantity{
			{ITEM_BRASS_INGOT, 2},
		}},

	{product: ItemQuantity{ITEM_IRON_PLATE, 1},
		components: []ItemQuantity{
			{ITEM_IRON_INGOT, 2},
		}},

	{product: ItemQuantity{ITEM_SCRAP_IRON, 2},
		components: []ItemQuantity{
			{ITEM_IRON_INGOT, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_IRON, 4},
		components: []ItemQuantity{
			{ITEM_IRON_PLATE, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_COPPER, 2},
		components: []ItemQuantity{
			{ITEM_COPPER_INGOT, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_COPPER, 4},
		components: []ItemQuantity{
			{ITEM_COPPER_PLATE, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_BRASS, 2},
		components: []ItemQuantity{
			{ITEM_BRASS_INGOT, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_BRASS, 4},
		components: []ItemQuantity{
			{ITEM_BRASS_PLATE, 1},
		}},

	{product: ItemQuantity{ITEM_SCRAP_ZINC, 2},
		components: []ItemQuantity{
			{ITEM_ZINC_INGOT, 1},
		}},
}

var furnaceRecipes = []Recipe{
	{product: ItemQuantity{ITEM_IRON_INGOT, 2},
		components: []ItemQuantity{
			{ITEM_IRON_ORE, 1},
		}},

	{product: ItemQuantity{ITEM_LODESTONE, 1},
		components: []ItemQuantity{
			{ITEM_IRON_ORE, 1},
		}},

	{product: ItemQuantity{ITEM_COPPER_INGOT, 2},
		components: []ItemQuantity{
			{ITEM_COPPER_ORE, 1},
		}},

	{product: ItemQuantity{ITEM_ZINC_INGOT, 2},
		components: []ItemQuantity{
			{ITEM_ZINC_ORE, 1},
		}},

	{product: ItemQuantity{ITEM_GLASS, 4},
		components: []ItemQuantity{
			{ITEM_QUARTZ, 1},
		}},
}
