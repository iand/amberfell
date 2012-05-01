/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

const (
	GroundLevel = 4

	CHUNK_WIDTH  = 32
	CHUNK_HEIGHT = 24

	XAXIS = 0
	YAXIS = 1
	ZAXIS = 2

	TILES_HORZ = 16
	TILES_VERT = 8

	TILE_WIDTH   = 48                 // Height and width of a block texture in pixels
	SCREEN_SCALE = 1.0 * TILE_WIDTH   // Width of one world coordinate unit in pixels
	PIXEL_SCALE  = 1.0 / SCREEN_SCALE // Width of one pixel in world coordinate units

	KEY_DEBOUNCE_DELAY = 3e8 // nanoseconds

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

	BLOCK_AIR      = 0
	BLOCK_STONE    = 1
	BLOCK_DIRT     = 2
	BLOCK_TRUNK    = 3
	BLOCK_LEAVES   = 4
	BLOCK_LOG_WALL = 5
	BLOCK_LOG_SLAB = 6

	ACTION_HAND   = 0
	ACTION_BREAK  = 1
	ACTION_WEAPON = 2
	ACTION_ITEM0  = 3
	ACTION_ITEM1  = 4
	ACTION_ITEM2  = 5
	ACTION_ITEM3  = 6
	ACTION_ITEM4  = 7

	ITEM_NONE = 4096

	// Terrain block textures
	TEXTURE_NONE         = 0
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
	STRENGTH_STONE  = 20
	STRENGTH_DIRT   = 3
	STRENGTH_WOOD   = 5
	STRENGTH_LEAVES = 1
	STRENGTH_IRON   = 50
)

var (
	NORMAL_EAST  = [3]float32{1.0, 0.0, 0.0}
	NORMAL_WEST  = [3]float32{-1.0, 0.0, 0.0}
	NORMAL_NORTH = [3]float32{0.0, 0.0, -1.0}
	NORMAL_SOUTH = [3]float32{0.0, 0.0, 1.0}
	NORMAL_UP    = [3]float32{0.0, 1.0, 0.0}
	NORMAL_DOWN  = [3]float32{0.0, -1.0, 0.0}

	COLOUR_WHITE = [4]float32{1.0, 1.0, 1.0, 1.0}
	COLOUR_HIGH  = [4]float32{96.0 / 255, 208.0 / 255, 96.0 / 255, 1.0}
)
