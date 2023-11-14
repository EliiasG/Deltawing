package buffers

import "unsafe"

func AddTo[Buffer any, Elem any](buffer *[]Buffer, elem Elem) {
	// amazing
	*buffer = append(*buffer, *(*Buffer)(unsafe.Pointer(&elem)))
}

func SetAt[Buffer any, Elem any](buffer *[]Buffer, idx int, elem Elem) {
	(*buffer)[idx] = *(*Buffer)(unsafe.Pointer(&elem))
}
