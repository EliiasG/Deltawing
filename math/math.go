package math

type Float interface {
	~float32 | ~float64
}

type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Signed interface {
	Float | SignedInt
}
