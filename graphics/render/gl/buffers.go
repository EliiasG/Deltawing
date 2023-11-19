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
	// The usage passed to BufferData
	Usage uint32
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

func (s *spriteBufferBuilder) MakeBuffer(static bool) render.SpriteBuffer {
	// make sure to not modify some VAO
	s.cxt.BindVertexArray(nil)
	sb := new(SpriteBuffer)
	sb.cxt = s.cxt
	if static {
		sb.Usage = enum.STATIC_DRAW
	} else {
		sb.Usage = enum.DYNAMIC_DRAW
	}
	// data
	sb.Verts = s.cxt.CreateBuffer()
	sb.Inds = s.cxt.CreateBuffer()
	s.Reallocate(sb)
	return sb
}

func (s *spriteBufferBuilder) Reallocate(buffer render.SpriteBuffer) {
	sb, _ := GLSpriteBuffer(buffer)
	var verts, inds []uint32
	// turn sprites into arrays for verts and inds
	verts, inds, sb.IdxPositions = util.CompileVecSpriteBuffer(s.sprites)
	// vertex buffer
	s.cxt.BindBuffer(enum.ARRAY_BUFFER, sb.Verts)
	// only * 4 because type is slice of uint32
	if len(verts) > 0 {
		s.cxt.BufferData(enum.ARRAY_BUFFER, verts, sb.Usage)
	}
	// index buffer
	s.cxt.BindBuffer(enum.ELEMENT_ARRAY_BUFFER, sb.Inds)
	if len(inds) > 0 {
		s.cxt.BufferData(enum.ELEMENT_ARRAY_BUFFER, inds, sb.Usage)
	}
}

func (s *spriteBufferBuilder) Clear() {
	s.sprites = make([]*vecsprite.VecSprite, 0)
}

/*
	DataBuffer
*/

type DataBuffer struct {
	cxt Context
	// buffer for data
	Buffer any
	// true if usgae is STATIC_DRAW
	Usage uint32
	// layout is part of VAO in opengl, so the glCalls should happen in the operation
	// Buffer layout of interleaved
	Layout []render.InputType
	// Size in bytes of a section of the buffer
	LayoutSize uint16
}

func (r *Renderer) MakeDataBuffer(static bool) render.DataBuffer {
	var usage uint32
	if static {
		usage = enum.STATIC_DRAW
	} else {
		usage = enum.DYNAMIC_DRAW
	}
	return &DataBuffer{
		cxt:    r.cxt,
		Buffer: r.cxt.CreateBuffer(),
		Usage:  usage,
	}
}

func (d *DataBuffer) bind() {
	d.cxt.BindBuffer(enum.ARRAY_BUFFER, d.Buffer)
}

func (d *DataBuffer) Free() {
	d.cxt.DeleteBuffer(d.Buffer)
}

// if only generic methods were a thing...
// code repitition will work for now
func (d *DataBuffer) SetData8(data []uint8) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.Usage)
	}
}

func (d *DataBuffer) SetData16(data []uint16) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.Usage)
	}
}

func (d *DataBuffer) SetData32(data []uint32) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.Usage)
	}
}

func (d *DataBuffer) SetData64(data []uint64) {
	d.bind()
	if len(data) > 0 {
		d.cxt.BufferData(enum.ARRAY_BUFFER, data, d.Usage)
	}
}

func (d *DataBuffer) SetLayout(layout ...render.InputType) {
	d.Layout = layout
	d.LayoutSize = 0
	for _, elem := range d.Layout {
		d.LayoutSize += uint16(render.SizeOf(elem))
	}
}
