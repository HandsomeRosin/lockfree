package lockfree

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNoMutex(t *testing.T) {
	queue = CreateMutexList(defaultVal, false)
	mq := queue.(*mutexList)
	mq.mutex = &emptyMutex{}
	t1 := time.Now()
	for i := 0; i < dataNum; i++ {
		queue.PushBack(i)
	}
	queue.Disable()
	for {
		val, enable := queue.PopFront()
		if !enable {
			break
		}
		if val == 0 {
			// fmt.Println("pop fail")
		} else {
			// fmt.Println(val)
		}
	}
	fmt.Println("用时：", time.Since(t1))
}

func TestBasicMutex(t *testing.T) {
	queue = CreateMutexList(defaultVal, false)
	t1 := time.Now()
	for i := 0; i < dataNum; i++ {
		queue.PushBack(i)
	}
	queue.Disable()
	for {
		val, enable := queue.PopFront()
		if !enable {
			break
		}
		if val == 0 {
			// fmt.Println("pop fail")
		} else {
			// fmt.Println(val)
		}
	}
	fmt.Println("用时：", time.Since(t1))
}

func TestBasicSpin(t *testing.T) {
	queue = CreateMutexList(defaultVal, true)
	t1 := time.Now()
	for i := 0; i < dataNum; i++ {
		queue.PushBack(i)
	}
	queue.Disable()
	for {
		val, enable := queue.PopFront()
		if !enable {
			break
		}
		if val == 0 {
			// fmt.Println("pop fail")
		} else {
			// fmt.Println(val)
		}
	}
	fmt.Println("用时：", time.Since(t1))
}

func TestMutex(t *testing.T) {
	queue = CreateMutexList(defaultVal, false)
	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}
	t1 := time.Now()
	for i := 0; i < asyncNum; i++ {
		wgr.Add(1)
		go readerMutex(i*1000000, &wgr)
	}
	for i := 0; i < asyncNum; i++ {
		wgw.Add(1)
		go writterMutex(&wgw)
	}
	wgr.Wait()
	queue.Disable()
	wgw.Wait()
	fmt.Println("用时：", time.Since(t1))
}

func TestSpin(t *testing.T) {
	queue = CreateMutexList(defaultVal, true)
	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}
	t1 := time.Now()
	for i := 0; i < asyncNum; i++ {
		wgr.Add(1)
		go readerMutex(i*1000000, &wgr)
	}
	for i := 0; i < asyncNum; i++ {
		wgw.Add(1)
		go writterMutex(&wgw)
	}
	wgr.Wait()
	queue.Disable()
	wgw.Wait()
	fmt.Println("用时：", time.Since(t1))
}

func readerMutex(startNum int, wg *sync.WaitGroup) {
	for i := 0; i < dataNum; i++ {
		suc := queue.PushBack(startNum + i)
		for !suc {
			suc = queue.PushBack(startNum + i)
		}
		// fmt.Println("push: ", queue.Size())
	}
	wg.Done()
}

func writterMutex(wg *sync.WaitGroup) {
	for {
		r, enable := queue.PopFront()
		if enable == false {
			break
		}
		if r == defaultVal {
			continue
		}
	}
	wg.Done()
}
