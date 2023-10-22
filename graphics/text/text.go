package text

import (
	"fmt"

	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/util/buffers"
)

type TextRenderer struct {
	// Operation used to draw text, it is not recommended to have other attrib channels than those for the textrenderer
	Operation render.Operation
	// Attribute channel for the position, should be 2 floats
	PositionChannel render.Channel
	GlyphBuffer     *GlyphBuffer
	// Glyph to draw if given glyph is not prsent
	DefaultGlyph rune
	LineSpacing  float32
	SpaceSpacing float32

	dataBuffer   render.DataBuffer
	positionData map[rune][]uint64
	indexMap     map[rune]uint32
}

func (t *TextRenderer) Init(renderer render.Renderer) {
	t.dataBuffer = renderer.MakeDataBuffer(false)
	t.dataBuffer.SetLayout(render.Input(render.InputFloat, 2))
	t.positionData = make(map[rune][]uint64)
}

func (t *TextRenderer) Clear() {
	for k := range t.positionData {
		delete(t.positionData, k)
	}
}

func (t *TextRenderer) AddText(x, y float32, text string) {
	var xOffset, yOffset float32
	for _, glyph := range text {
		if glyph == '\n' {
			yOffset += t.LineSpacing
			xOffset = 0
			continue
		}
		if glyph == ' ' {
			xOffset += t.SpaceSpacing
			continue
		}
		if glyph < 32 {
			continue
		}
		// fix glyph if not present
		if _, ok := t.GlyphBuffer.Glyphs[glyph]; !ok {
			glyph = t.DefaultGlyph
		}
		// add position
		buffer, ok := t.positionData[glyph]
		if !ok {
			buffer = make([]uint64, 0)
		}
		buffers.AddTo(&buffer, [2]float32{x + xOffset, yOffset})
		t.positionData[glyph] = buffer
		xOffset += t.GlyphBuffer.Glyphs[glyph].Advance
	}
}

func (t *TextRenderer) DrawTo(target render.RenderTarget) {
	for glyph, bufferedGlyph := range t.GlyphBuffer.Glyphs {
		_, ok := t.indexMap[glyph]
		if !ok {
			continue
		}
		//t.dataBuffer.SetData64(t.positionData[glyph])
		// buffer index
		t.Operation.SetInstanceAttribute(t.PositionChannel, t.dataBuffer, t.indexMap[glyph], 0)
		// sprite
		t.Operation.SetSprite(t.GlyphBuffer.SpriteBuffer, bufferedGlyph.Index)
		// get amount by length of positions
		t.Operation.SetAmount(uint32(len(t.positionData[glyph])))
		t.Operation.DrawTo(target)
	}
}

func (t *TextRenderer) UpdateText() {
	// calculate length
	amt := 0
	for _, data := range t.positionData {
		amt += len(data)
	}
	// merge data
	data := make([]uint64, amt)
	indexMap := make(map[rune]uint32)
	var i uint32
	for glyph, positions := range t.positionData {
		indexMap[glyph] = i
		for _, position := range positions {
			data[i] = position
			i++
		}
	}
	fmt.Println(i, amt)
	t.dataBuffer.SetData64(data)
	t.indexMap = indexMap
}
