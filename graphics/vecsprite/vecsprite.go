package vecsprite

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/eliiasg/deltawing/graphics/color"
)

type VecSprite struct {
	// per vertex
	Vertices [][2]float32
	Colors   []color.Color
	Layers   []uint8
	// indices
	Indices []uint32
}

// Descrition of tris format can be found at the bottom of this readme https://github.com/EliiasG/MonoGameDrawingApp#readme
func FromBytes(reader io.ByteReader) (*VecSprite, error) {
	verts, colors, layers, e := readVerts(reader)
	if e != nil {
		return nil, e
	}
	inds := readInds(reader)
	return &VecSprite{verts, colors, layers, inds}, nil
}

func readInds(reader io.ByteReader) []uint32 {
	inds := make([]uint32, 0)
	for true {
		idx, e := readUint(reader)
		if e != nil {
			return inds
		}
		inds = append(inds, idx)
	}
	//should never happen
	return nil

}

func readVerts(reader io.ByteReader) ([][2]float32, []color.Color, []uint8, error) {
	verts := make([][2]float32, 0)
	colors := make([]color.Color, 0)
	layers := make([]uint8, 0)
	colorChanges := true

	// color is not updated every vertex, a color is only given when it changes
	var curColor color.Color
	var layer uint8 = 0
	for true {
		// ignoring error to catch at end, even if it ends at start it would be fine to continue
		if colorChanges {
			curColor, _ = color.ReadARGB(reader)
		}

		x, _ := readFloat(reader)
		y, _ := readFloat(reader)
		// byte specifying what happens next
		change, e := reader.ReadByte()
		// 0 means no color change, 1 means color change, and 2 means no more vertices, so any value over 2 is invalid
		if e != nil || change > 2 {
			return nil, nil, nil, errors.New("Invalid bytes for vector sprite")
		}

		verts = append(verts, [2]float32{x, y})
		colors = append(colors, curColor)
		layers = append(layers, layer)
		colorChanges = change == 1
		if change == 1 {
			layer++
		}

		// change is 2 when verices section is done
		if change == 2 {
			return verts, colors, layers, nil
		}
	}

	//should never happen
	return nil, nil, nil, errors.New("oops")
}

func readFloat(reader io.ByteReader) (float32, error) {
	bytes, e := readUint(reader)
	return math.Float32frombits(bytes), e
}

func readUint(reader io.ByteReader) (uint32, error) {
	var e error
	bytes := [4]byte{}
	// only catch last error, if there is any error the last call sould also give an error.
	for i := 0; i < 4; i++ {
		bytes[i], e = reader.ReadByte()
	}
	bits := binary.LittleEndian.Uint32(bytes[:])
	return bits, e
}
