package lockfree

import (
	"sync"
	"sync/atomic"
)

// 自旋锁
type spinMutex struct {
	mutex int32
}

const locked = 1
const unlocked = 0

func (spin *spinMutex) lock() {
	// for !atomic.CompareAndSwapInt32(&spin.mutex, unlocked, locked) {
	// }
BEGINING:
	for spin.mutex != unlocked {
	}
	if !atomic.CompareAndSwapInt32(&spin.mutex, unlocked, locked) {
		goto BEGINING
	}
}

func (spin *spinMutex) unlock() {
	atomic.SwapInt32(&spin.mutex, unlocked)
}

type spinNcMutex struct {
	mutex int32
}

func (spin *spinNcMutex) lock() {
BEGINING:
	for spin.mutex != unlocked {
	}
	if !atomic.CompareAndSwapInt32(&spin.mutex, unlocked, locked) {
		goto BEGINING
	}
}

// 互斥锁
type mmutex struct {
	_mutex sync.Mutex
}

func (mu *mmutex) lock() {
	mu._mutex.Lock()
}

func (mu *mmutex) unlock() {
	mu._mutex.Unlock()
}

// 无效锁，仅作测试用
type emptyMutex struct{}

func (mu *emptyMutex) lock() {}

func (mu *emptyMutex) unlock() {}
