package lockfree

import "sync/atomic"
import "unsafe"

type itemPtr = unsafe.Pointer

type List struct {
	// NullVal 代表空值，在List为空时取值将会返回NullVal
	NuLLVal		interface{}
	head		itemPtr
	tail		itemPtr
	closed		bool
}

type item struct {
	next		itemPtr
	prev		itemPtr
	valPtr		interface{}
}

// null 代表itemPtr的默认状态，即空指针
var null itemPtr

func (this List) PushFront(val interface{}) bool {
	if this.closed {
		return false
	}
	// 将新Val构造为item节点
	node := new(item)
	node.valPtr = &val

	// 如果末尾指针为空，将其设置为该新节点。此处可能存在竞争。
	atomic.CompareAndSwapPointer(&this.tail, null, itemPtr(node))

	tmp := this.head
	node.next = tmp
	// recursively retry to set head of new item
	for !atomic.CompareAndSwapPointer(&this.head, tmp, itemPtr(node)) {
		tmp = this.head
		node.next = tmp
	}
	return true
}