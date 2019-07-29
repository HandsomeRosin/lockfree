package lockfree

import "sync/atomic"
import "unsafe"
import "fmt"

const INT64_MAX = int64(^uint64(0) >> 1)

type itemPtr = unsafe.Pointer
type ptr = unsafe.Pointer

type List struct {
	// NullVal 代表空值，在List为空时取值将会返回NullVal
	nullVal		interface{}
	head		itemPtr
	tail		itemPtr
	disabled	bool
	elNum		int64
	maxLength	int64
}

type item struct {
	next		itemPtr
	valPtr		ptr
}

// 创建一个无锁队列。defaultVal是取空链表时返回的默认值；maxSize是队列最大长度，若小于等于0，则默认其无限大。
func CreateList(defaultVal interface{}, maxSize int64) List {
	list := List {
		nullVal: 	defaultVal,
		disabled: 	false,
		elNum: 		0,
	}

	if maxSize > 0 {
		list.maxLength = maxSize
	} else {
		list.maxLength = INT64_MAX
	}

	// head和tail都初始化为同一个空节点，该节点不存实际元素。
	sentinel := itemPtr(&item{ nil, nil })
	list.head = sentinel
	list.tail = sentinel

	return list
}

// 往无锁链表中追加一个元素。第一个返回值代表元素是否插入成功；若链表状态是disabled，则第二个返回值是false。
func (this *List) PushBack(val interface{}) (bool, bool) {
	if this.disabled {
		return false, false
	}
	if this.elNum >= this.maxLength {
		return false, true
	}

	// 将新Val构造为item节点
	node := &item{
		next: nil,
		valPtr: ptr(&val),
	}
	node.valPtr = ptr(&val)

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
	atomic.AddInt64(&this.elNum, 1)

	atomic.CompareAndSwapPointer(&this.tail, itemPtr(oldp), itemPtr(node))
	return true, true
}

// 弹出一个元素。若链表状态是disabled且链表中已为空，则第二个返回值是false。
func (this *List) PopFront() (interface{}, bool) {
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
	atomic.AddInt64(&this.elNum, -1)

	return *((*interface{})(((*item)(p.next)).valPtr)), true
}

func (this *List) Disable() {
	this.disabled = true
}
func (this *List) Enable() {
	this.disabled = false
}
func (this *List) IsDisabled() bool {
	return this.disabled == true
}

// 返回当前链表的长度
func (this *List) Size() int64 {
	return this.elNum
}

// 返回当前链表是否为空
func (this *List) IsEmpty() bool {
	return (*item)(this.head).next == nil
}

// 设置链表最大长度
func (this *List) SetMaxSize(maxSize int64) {
	this.maxLength = maxSize
}

func (this *List) Print() {
	fmt.Print("size: ", this.elNum, ", data: ")
	p := (*item)(this.head)
	for p.next != nil {
		fmt.Print(*((*interface{})((*item)(p.next).valPtr)), ",")
		p = (*item)(p.next)
	}
	fmt.Println()
}