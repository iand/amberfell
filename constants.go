/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

const (
	// GroundLevel = 807

	CHUNK_WIDTH  = 8
	CHUNK_HEIGHT = 128

	MAP_DIAM       = 64000
	PLAYER_START_X = 31445 //32011 // 31767 //MAP_DIAM / 2
	PLAYER_START_Z = 32137 //31058 // 32009 //MAP_DIAM / 2
	NOISE_SCALE    = 16

	TREE_PRECIPITATION_MIN = 0.2 //0.44
	TREE_DENSITY_PCT       = 5

	VERTEX_BUFFER_CAPACITY = 20000

	MAX_ITEMS_IN_INVENTORY = 999

	XAXIS = 0
	YAXIS = 1
	ZAXIS = 2

	TILES_HORZ = 16
	TILES_VERT = 8

	TILE_WIDTH   = 48                 // Height and width of a block texture in pixels
	SCREEN_SCALE = 1.0 * TILE_WIDTH   // Width of one world coordinate unit in pixels
	PIXEL_SCALE  = 1.0 / SCREEN_SCALE // Width of one pixel in world coordinate units

	KEY_DEBOUNCE_DELAY = 3e8 // nanoseconds

	CAMPFIRE_INTENSITY = 6
	MAX_LIGHT_LEVEL    = 12

	FACE_NONE  = 6 // 
	EAST_FACE  = 0 // +ve x
	WEST_FACE  = 1 // -ve x
	NORTH_FACE = 2 // -ve z
	SOUTH_FACE = 3 // +ve z
	UP_FACE    = 4 // +ve y
	DOWN_FACE  = 5 // -ve y

	ORIENT_EAST  = 0
	ORIENT_NORTH = 1
	ORIENT_SOUTH = 2
	ORIENT_WEST  = 3

	MAX_ITEMS      = 4096
	BLOCK_AIR      = 0
	BLOCK_STONE    = 1
	BLOCK_DIRT     = 2
	BLOCK_TRUNK    = 3
	BLOCK_LEAVES   = 4
	BLOCK_LOG_WALL = 5
	BLOCK_LOG_SLAB = 6

	BLOCK_AMBERFELL    = 7
	BLOCK_COAL         = 8
	BLOCK_COPPER       = 9
	BLOCK_IRON         = 10
	BLOCK_CARVED_STONE = 11
	BLOCK_CAMPFIRE     = 12

	ITEM_NONE     = MAX_ITEMS - 1
	ITEM_FIREWOOD = 512

	ACTION_HAND   = 0
	ACTION_BREAK  = 1
	ACTION_WEAPON = 2
	ACTION_ITEM0  = 3
	ACTION_ITEM1  = 4
	ACTION_ITEM2  = 5
	ACTION_ITEM3  = 6
	ACTION_ITEM4  = 7

	// Terrain block textures
	TEXTURE_NONE         = 0
	TEXTURE_CARVED_STONE = 16
	TEXTURE_STONE        = 17
	TEXTURE_STONE_TOP    = 1
	TEXTURE_DIRT         = 18
	TEXTURE_DIRT_TOP     = 2
	TEXTURE_TRUNK        = 19
	TEXTURE_TRUNK_TOP    = 3
	TEXTURE_LEAVES       = 20
	TEXTURE_LEAVES_TOP   = 20
	TEXTURE_LOG_WALL     = 21
	TEXTURE_LOG_WALL_TOP = 5

	TEXTURE_COAL          = 7
	TEXTURE_COPPER        = 8
	TEXTURE_IRON          = 9
	TEXTURE_AMBERFELL_TOP = 6
	TEXTURE_AMBERFELL     = 22

	// Player textures

	TEXTURE_HAT_FRONT = 4096
	TEXTURE_HAT_LEFT  = 4097
	TEXTURE_HAT_BACK  = 4098
	TEXTURE_HAT_RIGHT = 4099
	TEXTURE_HAT_TOP   = 4100

	TEXTURE_HEAD_FRONT  = 4101
	TEXTURE_HEAD_LEFT   = 4102
	TEXTURE_HEAD_BACK   = 4103
	TEXTURE_HEAD_RIGHT  = 4104
	TEXTURE_HEAD_BOTTOM = 4105

	TEXTURE_TORSO_FRONT = 4106
	TEXTURE_TORSO_LEFT  = 4107
	TEXTURE_TORSO_BACK  = 4108
	TEXTURE_TORSO_RIGHT = 4109
	TEXTURE_TORSO_TOP   = 4110

	TEXTURE_LEG      = 4111
	TEXTURE_LEG_SIDE = 4112
	TEXTURE_ARM      = 4113
	TEXTURE_ARM_TOP  = 4114
	TEXTURE_HAND     = 4116
	TEXTURE_BRIM     = 4117

	// Mob textures

	// HUD textures
	TEXTURE_PICKER = 5000

	// Strength of materials
	STRENGTH_STONE       = 20
	STRENGTH_DIRT        = 3
	STRENGTH_WOOD        = 5
	STRENGTH_LEAVES      = 1
	STRENGTH_IRON        = 50
	STRENGTH_UNBREAKABLE = 255
)

var (
	NORMALS = [6]([3]float32){[3]float32{1.0, 0.0, 0.0},
		[3]float32{-1.0, 0.0, 0.0},
		[3]float32{0.0, 0.0, -1.0},
		[3]float32{0.0, 0.0, 1.0},
		[3]float32{0.0, 1.0, 0.0},
		[3]float32{0.0, -1.0, 0.0},
	}

	LIGHT_LEVELS = [13]([4]float32){[4]float32{0, 0, 0, 1},
		[4]float32{0.04, 0.04, 0.04, 1.0},
		[4]float32{0.12, 0.12, 0.12, 1.0},
		[4]float32{0.20, 0.20, 0.20, 1.0},
		[4]float32{0.28, 0.28, 0.28, 1.0},
		[4]float32{0.36, 0.36, 0.36, 1.0},
		[4]float32{0.44, 0.44, 0.44, 1.0},
		[4]float32{0.52, 0.52, 0.52, 1.0},
		[4]float32{0.60, 0.60, 0.60, 1.0},
		[4]float32{0.68, 0.68, 0.68, 1.0},
		[4]float32{0.76, 0.76, 0.76, 1.0},
		[4]float32{0.84, 0.84, 0.84, 1.0},
		[4]float32{0.92, 0.92, 0.92, 1.0},
	}
	COLOUR_WHITE = [4]float32{1.0, 1.0, 1.0, 1.0}
	COLOUR_HIGH  = [4]float32{96.0 / 255, 208.0 / 255, 96.0 / 255, 1.0}
)
