/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/Go-SDL/ttf"
	"github.com/banthar/gl"
	"image"
	"image/color"
)

type Font struct {
	height   int
	textures map[rune]gl.Texture
	widths   map[rune]int
}

func NewFont(filename string, size int, c color.Color) *Font {

	var font Font
	font.textures = make(map[rune]gl.Texture)
	font.widths = make(map[rune]int)

	extfont := ttf.OpenFont(filename, size)
	if extfont == nil {
		panic("Could not load font")
	}
	defer extfont.Close()

	font.height = extfont.Ascent()

	for _, ch := range "abcdefghijklmnopqrdstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 :;@'<>,.?/~#{}[]!£$%^&*()_-+=\"\\|" {
		minx, _, _, maxy, advance, err := extfont.GlyphMetrics(uint16(ch))
		if err != 0 {
			panic("Could not get glyph metrics")
		}
		_, _ = minx, maxy

		// Create a bitmap with width=advance, height=font.height
		surface := sdl.CreateRGBSurface(sdl.SWSURFACE|sdl.SRCALPHA, advance, font.height, 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
		//surface.FillRect(&sdl.Rect{0,0, uint16(advance), uint16(font.height)}, 0x0)
		// rect := sdl.Rect{0,0, uint16(advance), uint16(ascent)}
		// rect := sdl.Rect{int16(minx), int16(ascent)-int16(maxy), 0, 0}

		fontSurface := ttf.RenderText_Blended(extfont, string(ch), sdl.ColorFromGoColor(c))
		fontSurface.Blit(nil, surface, nil)

		rgba := image.NewRGBA(image.Rect(0, 0, advance, font.height))
		for x := 0; x < advance; x++ {
			for y := 0; y < font.height; y++ {
				rgba.Set(x, y, fontSurface.At(x, font.height-y))
			}
		}

		font.widths[ch] = advance
		font.textures[ch] = gl.GenTexture()
		font.textures[ch].Bind(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, advance, font.height, 0, gl.RGBA, gl.UNSIGNED_BYTE, &rgba.Pix[0])
		font.textures[ch].Unbind(gl.TEXTURE_2D)
	}
	return &font
}

func (self *Font) Print(str string) {
	for _, ch := range str {
		self.textures[ch].Bind(gl.TEXTURE_2D)
		h := float32(self.height) * PIXEL_SCALE
		w := float32(self.widths[ch]) * PIXEL_SCALE
		gl.Color4ub(255, 255, 255, 255)
		gl.Begin(gl.QUADS)
		gl.TexCoord2d(0, 0)
		gl.Vertex2f(0, 0) // Bottom Left Of The Texture and Quad
		gl.TexCoord2d(1, 0)
		gl.Vertex2f(w, 0) // Bottom Right Of The Texture and Quad
		gl.TexCoord2d(1, 1)
		gl.Vertex2f(w, h) // Top Right Of The Texture and Quad
		gl.TexCoord2d(0, 1)
		gl.Vertex2f(0, h) // Top Left Of The Texture and Quad
		gl.End()
		gl.Translatef(w, 0, 0)
		self.textures[ch].Unbind(gl.TEXTURE_2D)
	}
}
