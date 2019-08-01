package lockfree

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var defaultVal int
var queue List

func init() {
	defaultVal = -10086
	queue = CreateCASList(defaultVal)
}

var dataNum = 2500 // 测试数据量
var asyncNum = 4   // 测试并发协程数

func TestBasicCas(t *testing.T) {
	t1 := time.Now()
	for i := 1; i <= dataNum; i++ {
		suc := queue.PushBack(i)
		if !suc {
			// fmt.Println("push fail: ", i)
		} else {
			// queue.Print()
		}
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
			// fmt.Print(val, " : ")
			// queue.Print()
		}
	}
	fmt.Println("用时：", time.Since(t1))
}
func TestCas(t *testing.T) {
	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}
	t1 := time.Now()
	for i := 0; i < asyncNum; i++ {
		wgr.Add(1)
		go reader(i*1000000, &wgr)
	}
	for i := 0; i < asyncNum; i++ {
		wgw.Add(1)
		go writter(&wgw)
	}
	wgr.Wait()
	queue.Disable()
	wgw.Wait()
	fmt.Println("用时：", time.Since(t1))
	fmt.Println("END-------------------------------------------")
}

func TestChannel(t *testing.T) {
	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}
	t1 := time.Now()
	for i := 0; i < 8; i++ {
		wgr.Add(1)
		go chReader(i*1000000, &wgr)
	}
	for i := 0; i < 8; i++ {
		wgw.Add(1)
		go chWriter(&wgw)
	}
	wgr.Wait()
	close(ch)
	wgw.Wait()
	fmt.Println("用时：", time.Since(t1))
	fmt.Println("END-------------------------------------------")
}

var buffer int = 0
var ch chan int = make(chan int, buffer)

// var ch chan int = make(chan int, buffer)

func reader(startNum int, wg *sync.WaitGroup) {
	for i := 0; i < dataNum; i++ {
		queue.PushBack(startNum + i)
		// fmt.Println("push: ", queue.Size())
	}
	wg.Done()
}

var m map[int]string = make(map[int]string)
var mutex sync.Mutex = sync.Mutex{}

func writter(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		r, enable := queue.PopFront()
		if enable == false {
			break
		}
		if r == defaultVal {
			continue
		}

		// mutex.Lock()
		// m[r.(int)] = "" // 为了核对pop出来的数据总数是否与push进去的一样，为了防止竞争错误导致的重复，这里用map来防重
		// mutex.Unlock()
	}
}

func chReader(startNum int, wg *sync.WaitGroup) {
	for i := 0; i < dataNum; i++ {
		ch <- i + startNum
	}
	wg.Done()
}

func chWriter(wg *sync.WaitGroup) {
	for {
		_, ok := <-ch
		if !ok {
			break
		}
		// fmt.Println("pop: ", r)
	}
	wg.Done()
}
