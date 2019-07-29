package lockfree

import (
	"fmt"
	"sync"
	"testing"
)

var defaultVal int
var queue List

func init() {
	defaultVal = -10086
	queue = CreateList(defaultVal, 20)
}

func TestBasic(t *testing.T) {
	for i := 1; i <= 50; i++ {
		suc, _ := queue.PushBack(i)
		if (!suc) {
			fmt.Println("push fail: ", i)
		} else {
			queue.Print()
		}
	}
	queue.Disable()

	for {
		val, enable := queue.PopFront()
		if !enable {
			break
		}
		if val == 0 {
			fmt.Println("pop fail")
		} else {
			fmt.Print(val, " : ")
			queue.Print()
		}
	}
}

func TestList(t *testing.T) {
	wgr := sync.WaitGroup{}
	wgw := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wgr.Add(1)
		go reader(i*1000000, &wgr)
	}
	for i := 0; i < 4; i++ {
		wgw.Add(1)
		go writter(&wgw)
	}
	wgr.Wait()
	queue.Disable()
	wgw.Wait()
	fmt.Println("total pop num: ", len(m))
	fmt.Println("END-------------------------------------------")
}

func reader(startNum int, wg *sync.WaitGroup) {
	for i := 0; i < 50; i++ {
		suc, enable := queue.PushBack(startNum + i)
		for  !suc {
			if !enable {
				goto END
			}
			suc, enable = queue.PushBack(startNum + i)
		}
		// fmt.Println("push: ", startNum + i)
	}
END:
	wg.Done()
}

var m map[int]string = make(map[int]string)
var mutex sync.Mutex = sync.Mutex{}

func writter(wg *sync.WaitGroup) {
	for {
		r, enable := queue.PopFront()
		if enable == false {
			break
		}
		if r == defaultVal {
			continue
		}
		// fmt.Println("pop: ", r)
		mutex.Lock()
		m[r.(int)] = "" // 为了核对pop出来的数据总数是否与push进去的一样，为了防止竞争错误导致的重复，这里用map来防重
		mutex.Unlock()
	}
	wg.Done()
}
