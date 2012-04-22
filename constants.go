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

	TILE_WIDTH   = 48                 // Height and width of a block texture in pixels
	SCREEN_SCALE = 1.0 * TILE_WIDTH   // Width of one world coordinate unit in pixels
	PIXEL_SCALE  = 1.0 / SCREEN_SCALE // Width of one pixel in world coordinate units

	KEY_DEBOUNCE_DELAY = 3e8 // nanoseconds

	FACE_NONE  = 0 // +ve x
	EAST_FACE  = 1 // +ve x
	WEST_FACE  = 2 // -ve x
	NORTH_FACE = 3 // -ve z
	SOUTH_FACE = 4 // +ve z
	UP_FACE    = 5 // +ve y
	DOWN_FACE  = 6 // -ve y

	BLOCK_AIR   = 0
	BLOCK_STONE = 1
	BLOCK_DIRT  = 2

	ACTION_HAND   = 0
	ACTION_BREAK  = 1
	ACTION_WEAPON = 2
	ACTION_ITEM   = 3

	ITEM_NONE = 4096
)
