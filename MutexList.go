package lockfree

type myMutex interface {
	lock()
	unlock()
}

type mutexList struct {
	nullVal  interface{}
	head     *mutexItem
	tail     *mutexItem
	disabled bool
	elNum    int64
	mutex    myMutex
}

type mutexItem struct {
	next   *mutexItem
	valPtr ptr
}

// CreateMutexList 创建一个队列。
func CreateMutexList(defaultVal interface{}, spinlock bool) List {
	list := &mutexList{
		nullVal:  defaultVal,
		disabled: false,
		elNum:    0,
	}

	if spinlock {
		list.mutex = new(spinMutex)
	} else {
		list.mutex = new(mmutex)
	}

	// head和tail都初始化为同一个空节点，该节点不存实际元素。
	sentinel := (&mutexItem{nil, nil})
	list.head = sentinel
	list.tail = sentinel

	return list
}

func (list *mutexList) PushBack(val interface{}) bool {
	if list.disabled {
		return false
	}

	node := &mutexItem{
		next:   nil,
		valPtr: ptr(&val),
	}

	list.mutex.lock()
	list.tail.next = node
	list.tail = node
	list.mutex.unlock()
	return true
}

func (list *mutexList) PopFront() (interface{}, bool) {

	list.mutex.lock()
	p := list.head
	if p.next == nil {
		list.mutex.unlock()
		if list.disabled {
			return list.nullVal, false
		}
		return list.nullVal, true
	}
	list.head = p.next
	list.mutex.unlock()

	return *((*interface{})(p.next.valPtr)), true
}
func (list *mutexList) Disable() {
	list.disabled = true
}
func (list *mutexList) Enable() {
	list.disabled = false
}
func (list *mutexList) IsEmpty() bool {
	return list.head.next == nil
}
