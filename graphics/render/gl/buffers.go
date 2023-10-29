package gl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/render/gl/util"
	"github.com/eliiasg/deltawing/graphics/vecsprite"
	"github.com/eliiasg/glow/enum"
)

/*
	SpriteBuffer
*/

func (r *Renderer) MakeSpriteBufferBuilder() render.SpriteBufferBuilder {
	return &spriteBufferBuilder{r.cxt, make([]*vecsprite.VecSprite, 0)}
}

// not making builders public, since they dont directly interact with gl
type spriteBufferBuilder struct {
	cxt     Context
	sprites []*vecsprite.VecSprite
}

type SpriteBuffer struct {
	render.SpriteBufferIdentifier
	cxt Context
	// Buffer of verts
	Verts any
	// Buffer of inds
	Inds any
	// Start indices for every sprite, first is 0, last is len(Inds)
	IdxPositions []uint32
}

func GLSpriteBuffer(s render.SpriteBuffer) (*SpriteBuffer, bool) {
	res, ok := s.(*SpriteBuffer)
	return res, ok
}

func (s *SpriteBuffer) Free() {
	s.cxt.DeleteBuffer(s.Verts)
	s.cxt.DeleteBuffer(s.Inds)
}

func (s *spriteBufferBuilder) AddSprite(sprite *vecsprite.VecSprite) uint32 {
	s.sprites = append(s.sprites, sprite)
	// - 1 because len will be 1 after first sprite is added
	return uint32(len(s.sprites) - 1)
}

func (s *spriteBufferBuilder) Finish() render.SpriteBuffer {
	// make sure to not modify some VAO
	s.cxt.BindVertexArray(nil)
	sb := new(SpriteBuffer)
	sb.cxt = s.cxt
	var verts, inds []uint32
	// turn sprites into arrays for verts and inds
	verts, inds, sb.IdxPositions = util.CompileVecSpriteBuffer(s.sprites)

	// vertex buffer
	sb.Verts = s.cxt.CreateBuffer()
	s.cxt.BindBuffer(enum.ARRAY_BUFFER, sb.Verts)
	// only * 4 because type is slice of uint32
	s.cxt.BufferData(enum.ARRAY_BUFFER, verts, enum.STATIC_DRAW)

	// index buffer
	sb.Inds = s.cxt.CreateBuffer()
	s.cxt.BindBuffer(enum.ELEMENT_ARRAY_BUFFER, sb.Inds)
	s.cxt.BufferData(enum.ELEMENT_ARRAY_BUFFER, inds, enum.STATIC_DRAW)
	return sb
}

/*
	DataBuffer
*/

type dataBuffer struct {
	cxt Context
	// buffer for data
	Buffer any
	// true if usgae is STATIC_DRAW
	Static bool
	// layout is part of VAO in opengl, so the glCalls should happen in the operation
	// Buffer layout of interleaved
	Layout []render.InputType
	// Size in bytes of a section of the buffer
	LayoutSize uint16
}

func (r *Renderer) MakeDataBuffer(static bool) render.DataBuffer {
	return &dataBuffer{
		cxt:    r.cxt,
		Buffer: r.cxt.CreateBuffer(),
		Static: static,
	}
}

func (d *dataBuffer) usage() uint32 {
	if d.Static {
		return enum.STATIC_DRAW
	}
	return enum.DYNAMIC_DRAW
}

func (d *dataBuffer) bind() {
	d.cxt.BindBuffer(enum.ARRAY_BUFFER, d.Buffer)
}

func (d *dataBuffer) Free() {
	d.cxt.DeleteBuffer(d.Buffer)
}

// if only generic methods were a thing...
// code repitition will work for now
func (d *dataBuffer) SetData8(data []uint8) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.usage())
	}
}

func (d *dataBuffer) SetData16(data []uint16) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.usage())
	}
}

func (d *dataBuffer) SetData32(data []uint32) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.usage())
	}
}

func (d *dataBuffer) SetData64(data []uint64) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.usage())
	}
}

func (d *dataBuffer) SetLayout(layout ...render.InputType) {
	d.Layout = layout
	d.LayoutSize = 0
	for _, elem := range d.Layout {
		d.LayoutSize += uint16(render.SizeOf(elem))
	}
}
