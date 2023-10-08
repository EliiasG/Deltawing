package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/vecsprite"
	"github.com/eliiasg/glow/v3.3-core/gl"

	g "github.com/eliiasg/deltawing/internal/rendering/gl"
)

/*
	SpriteBuffer
*/

func (r *renderer) MakeSpriteBufferBuilder() render.SpriteBufferBuilder {
	return &spriteBufferBuilder{make([]*vecsprite.VecSprite, 0)}
}

type spriteBufferBuilder struct {
	sprites []*vecsprite.VecSprite
}

type spriteBuffer struct {
	render.SpriteBufferIdentifier
	vertsID       uint32
	indsID        uint32
	vertPositions []uint32
	idxPositions  []uint32
}

func (s *spriteBuffer) Free() {
	// doing multiple because function expects pointer and this seems easier that getting a pointer for both of them
	gl.DeleteBuffers(1, &s.vertsID)
	gl.DeleteBuffers(1, &s.indsID)
}

func (s *spriteBufferBuilder) AddSprite(sprite *vecsprite.VecSprite) uint32 {
	s.sprites = append(s.sprites, sprite)
	// - 1 because len will be 1 after first sprite is added
	return uint32(len(s.sprites) - 1)
}

func (s *spriteBufferBuilder) Finish() render.SpriteBuffer {
	// make sure to not modify some VAO
	gl.BindVertexArray(0)
	sb := new(spriteBuffer)
	var verts, inds []uint32
	// turn sprites into arrays for verts and inds
	verts, inds, sb.vertPositions, sb.idxPositions = g.CompileVecSpriteBuffer(s.sprites)

	// vertex buffer
	gl.GenBuffers(1, &sb.vertsID)
	gl.BindBuffer(gl.ARRAY_BUFFER, sb.vertsID)
	// only * 4 because type is slice of uint32
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	// index buffer
	gl.GenBuffers(1, &sb.indsID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sb.indsID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(inds)*4, gl.Ptr(inds), gl.STATIC_DRAW)
	return sb
}

/*
	DataBuffer
*/

func (r *renderer) MakeDataBuffer(static bool) render.DataBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	return &dataBuffer{
		id:     id,
		static: static,
	}
}

type dataBuffer struct {
	id     uint32
	static bool
	// layout is part of VAO in opengl, so the glCalls should happen in the operation
	layout     []render.InputType
	layoutSize uint16
}

func (d *dataBuffer) usage() uint32 {
	if d.static {
		return gl.STATIC_DRAW
	}
	return gl.DYNAMIC_DRAW
}

func (d *dataBuffer) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, d.id)
}

func (d *dataBuffer) Free() {
	gl.DeleteBuffers(1, &d.id)
}

// if only generic methods were a thing...
// code repitition will work for now
func (d *dataBuffer) SetData8(data []uint8) {
	d.bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*1, gl.Ptr(data), d.usage())
}

func (d *dataBuffer) SetData16(data []uint16) {
	d.bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*2, gl.Ptr(data), d.usage())
}

func (d *dataBuffer) SetData32(data []uint32) {
	d.bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), d.usage())
}

func (d *dataBuffer) SetData64(data []uint64) {
	d.bind()
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*8, gl.Ptr(data), d.usage())
}

func (d *dataBuffer) SetLayout(layout ...render.InputType) {
	d.layout = layout
	d.layoutSize = 0
	for _, elem := range d.layout {
		d.layoutSize += uint16(render.SizeOf(elem))
	}
}
