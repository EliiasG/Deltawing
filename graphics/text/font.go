package text

import (
	"errors"

	"github.com/eliiasg/deltawing/graphics/color"
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/vecsprite"
	"github.com/eliiasg/trifont"
)

// A glyph on in a GlyphBuffer
type BufferedGlyph struct {
	Index   uint32
	Advance float32
}

// Basically a font, this contains a spritebuffer with glyphs.
type GlyphBuffer struct {
	SpriteBuffer render.SpriteBuffer
	Glyphs       map[rune]BufferedGlyph
}

// Makes a sprite from a trifont.Char
func SpriteFromGlyph(char *trifont.Char) *vecsprite.VecSprite {
	sprite := new(vecsprite.VecSprite)
	// vertices
	sprite.Vertices = make([][2]float32, len(char.Vertices))
	sprite.Colors = make([]color.Color, len(char.Vertices))
	sprite.Layers = make([]uint8, len(char.Vertices))
	for i, vert := range char.Vertices {
		sprite.Vertices[i] = vert
		sprite.Colors[i] = color.Black()
		sprite.Layers[i] = 0
	}
	// indices
	sprite.Indices = make([]uint32, len(char.Indices))
	for i, index := range char.Indices {
		sprite.Indices[i] = uint32(index)
	}
	return sprite
}

// Make GlyphBuffer from font, glyphSet indicates what glyphs to inclue, if glyphSet is empty, all glyphs will be included.
// Will cause error if a glyph the font does not support is in glyphSet
func LoadFont(font *trifont.Font, glyphSet string, renderer render.Renderer) (*GlyphBuffer, error) {
	glyphs := getGlyphs(font, glyphSet)
	glyphMap := make(map[rune]BufferedGlyph)
	builder := renderer.MakeSpriteBufferBuilder()
	// glyphs
	for _, glyph := range glyphs {
		char, ok := font.Chars[glyph]
		if !ok {
			return nil, errors.New("Glyph not in font: '" + string(glyph) + "'")
		}
		idx := builder.AddSprite(SpriteFromGlyph(&char))
		glyphMap[glyph] = BufferedGlyph{idx, char.Advance}
	}
	return &GlyphBuffer{builder.Finish(), glyphMap}, nil
}

func getGlyphs(font *trifont.Font, glyphSet string) []rune {
	if len(glyphSet) == 0 {
		// returun all glyphs of font if glyphSet is empty
		keys := make([]rune, 0, len(font.Chars))
		for k := range font.Chars {
			keys = append(keys, k)
		}
		return keys
	}
	return []rune(glyphSet)
}
