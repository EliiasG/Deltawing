package gl

import (
	"github.com/eliiasg/deltawing/graphics/vecsprite"
	"github.com/eliiasg/deltawing/util/buffers"
)

// returns verts, inds, start positions for verts, start positions for inds
// inds are 32 bit unsigned
// verts are: 32 bit float x, 32 bit float y, 24 bit color (rgb), 8 bit layer
func CompileVecSpriteBuffer(sprites []*vecsprite.VecSprite) ([]uint32, []uint32, []uint32, []uint32) {
	// counting at start for optimization
	// counting only iterates sprites, not verts and inds
	// maybe this would be called premature optimization
	numVerts, numInds := countSizes(sprites)
	// slices are of uint32, bits of original values are then added
	verts := make([]uint32, 0, numVerts)
	inds := make([]uint32, 0, numInds)
	var vertPos, idxPos uint32
	vertPositions := make([]uint32, 0, len(sprites)+1)
	idxPositions := make([]uint32, 0, len(sprites)+1)

	for _, sprite := range sprites {
		addSprite(&verts, &inds, sprite)
		vertPositions = append(vertPositions, vertPos)
		idxPositions = append(idxPositions, idxPos)
		// incrementing after since positions start at 0
		vertPos += uint32(len(sprite.Vertices))
		idxPos += uint32(len(sprite.Indices))
	}

	// appending at end because it should be possible to get the size of the last sprite
	vertPositions = append(vertPositions, vertPos)
	idxPositions = append(idxPositions, idxPos)

	return verts, inds, vertPositions, idxPositions
}

func addSprite(verts, inds *[]uint32, sprite *vecsprite.VecSprite) {
	// add vertices
	for i, vert := range sprite.Vertices {
		// adding different types to the int array using unsafe
		buffers.AddTo(verts, vert.X)
		buffers.AddTo(verts, vert.Y)
		col := sprite.Colors[i]
		// r, g, a, layer as uint8s
		buffers.AddTo(verts, [4]uint8{col.R, col.G, col.B, sprite.Layers[i]})
	}
	// add indices
	for _, idx := range sprite.Indices {
		// maybe it should just be added directly, since the slice is same type
		buffers.AddTo(inds, idx)
	}
}

func countSizes(sprites []*vecsprite.VecSprite) (verts uint32, inds uint32) {
	for _, sprite := range sprites {
		verts += uint32(len(sprite.Vertices))
		inds += uint32(len(sprite.Indices))
	}
	return
}
