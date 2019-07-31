package lockfree

import "unsafe"

const INT64_MAX = int64(^uint64(0) >> 1)

type ptr = unsafe.Pointer

type List interface {
	PushBack(val interface{}) bool
	PopFront() (interface{}, bool)
	Disable()
	Enable()
	IsEmpty() bool
}
