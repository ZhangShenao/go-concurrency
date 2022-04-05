package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {

}

// ReentrantMutex 可重入的Mutex
type ReentrantMutex struct {
	sync.Mutex
	owner          int64 //记录当前Mutex的持有者的GoroutineId
	reentrantTimes int32 //记录当前Mutex的重入次数
}

func (m *ReentrantMutex) Lock() {
	//获取当前调用者的GoroutineId
	goroutineId := m.goroutineId()

	if atomic.LoadInt64(&m.owner) == m.owner { //如果当前当前GoroutineId和owner一致,说明是重入操作,直接将重入次数+1即可
		m.reentrantTimes++
		return
	}

	//执行Lock操作
	m.Mutex.Lock()

	//记录获取Mutex的Goroutine信息
	m.owner = goroutineId
	atomic.StoreInt64(&m.owner, goroutineId) //cas
	atomic.StoreInt32(&m.reentrantTimes, 1)
}

func (m *ReentrantMutex) Unlock() {
	//获取当前调用者的GoroutineId
	goroutineId := m.goroutineId()

	//如果当前当前GoroutineId和owner一致,说明是异常操作,直接panic
	if atomic.LoadInt64(&m.owner) != goroutineId {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, goroutineId))
	}

	//重入次数-1
	m.reentrantTimes--
	if m.reentrantTimes != 0 { //如果重入次数大于0,说明Goroutine仍然持有Mutex,直接返回即可
		return
	}

	//重入次数=0,需要释放Mutex
	atomic.StoreInt64(&m.owner, -1) //清空持有者信息
	m.Mutex.Unlock()
}

//获取当前GoroutineId
func (m *ReentrantMutex) goroutineId() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
