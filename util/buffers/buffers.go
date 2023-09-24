package buffers

import "unsafe"

func AddTo[Buffer any, Elem any](buffer *[]Buffer, elem Elem) {
	// amazing
	*buffer = append(*buffer, *(*Buffer)(unsafe.Pointer(&elem)))
}
