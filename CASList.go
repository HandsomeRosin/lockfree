package lockfree

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type itemPtr = unsafe.Pointer

type casList struct {
	// NullVal 代表空值，在List为空时取值将会返回NullVal
	nullVal  interface{}
	head     itemPtr
	tail     itemPtr
	disabled bool
	// elNum    int64
}

type item struct {
	next   itemPtr
	valPtr ptr
}

// CreateCASList 创建一个无锁队列。defaultVal是取空链表时返回的默认值；
func CreateCASList(defaultVal interface{}) List {
	list := &casList{
		nullVal:  defaultVal,
		disabled: false,
		// elNum:    0,
	}

	// head和tail都初始化为同一个空节点，该节点不存实际元素。
	sentinel := itemPtr(&item{nil, nil})
	list.head = sentinel
	list.tail = sentinel

	return list
}

// PushBack 往无锁链表中追加一个元素。第一个返回值代表元素是否插入成功；若链表状态是disabled，则第二个返回值是false。
func (this *casList) PushBack(val interface{}) bool {
	if this.disabled {
		return false
	}

	// 将新Val构造为item节点
	node := &item{
		next:   nil,
		valPtr: ptr(&val),
	}

	p := (*item)(this.tail)
	oldp := p

	for p.next != nil {
		p = (*item)(p.next)
	}
	for !atomic.CompareAndSwapPointer(&p.next, nil, itemPtr(node)) {
		for p.next != nil {
			p = (*item)(p.next)
		}
	}
	// atomic.AddInt64(&this.elNum, 1)

	atomic.CompareAndSwapPointer(&this.tail, itemPtr(oldp), itemPtr(node))
	return true
}

// 弹出一个元素。若链表状态是disabled且链表中已为空，则第二个返回值是false。
func (this *casList) PopFront() (interface{}, bool) {
	p := (*item)(this.head)
	if p.next == nil {
		if this.disabled {
			// 如果已不可用，且链表中已无数据
			return this.nullVal, false
		}
		return this.nullVal, true // 如果链表为空，则返回NullVal
	}
	for !atomic.CompareAndSwapPointer(&this.head, itemPtr(p), itemPtr(p.next)) {
		p := (*item)(this.head)
		if p.next == nil {
			if this.disabled {
				// 如果已不可用，且链表中已无数据
				return this.nullVal, false
			}
			return this.nullVal, true // 如果链表为空，则返回NullVal
		}
	}
	// atomic.AddInt64(&this.elNum, -1)

	return *((*interface{})(((*item)(p.next)).valPtr)), true
}

func (this *casList) Disable() {
	this.disabled = true
}
func (this *casList) Enable() {
	this.disabled = false
}
func (this *casList) IsDisabled() bool {
	return this.disabled == true
}

// 返回当前链表是否为空
func (this *casList) IsEmpty() bool {
	return (*item)(this.head).next == nil
}

func (this *casList) Print() {
	p := (*item)(this.head)
	for p.next != nil {
		fmt.Print(*((*interface{})((*item)(p.next).valPtr)), ",")
		p = (*item)(p.next)
	}
	fmt.Println()
}
