package color

import "io"

type Color struct {
	R, G, B, A uint8
}

func FromRGBA(r, g, b, a uint8) Color {
	return Color{r, g, b, a}
}

func FromARGB(a, r, g, b uint8) Color {
	return Color{r, g, b, a}
}

func (c Color) ToRGBA() [4]uint8 {
	return [4]uint8{c.R, c.G, c.B, c.A}
}

func (c Color) ToARGB() [4]uint8 {
	return [4]uint8{c.A, c.R, c.G, c.B}
}

func White() Color {
	return FromRGBA(255, 255, 255, 255)
}

func Black() Color {
	return FromRGBA(0, 0, 0, 255)
}

func ReadARGB(reader io.ByteReader) (Color, error) {
	a, _ := reader.ReadByte()
	r, _ := reader.ReadByte()
	g, _ := reader.ReadByte()
	// error saved for last byte
	b, e := reader.ReadByte()
	if e != nil {
		return White(), e
	}
	return FromARGB(a, r, g, b), nil
}
