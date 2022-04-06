package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

//Mutex功能扩展
//1.实现TryLock
//2.获取等待者的数量等统计指标

func main() {
	printMutexState()
}

//TryLock示例
func tryLockDemo() {
	var mu Mutex
	go func() {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	}()

	time.Sleep(time.Second)

	//try lock
	if mu.TryLock() {
		fmt.Println("got lock")
		mu.Unlock()
		return
	}
	fmt.Println("can't get lock")
}

//打印锁的状态信息
func printMutexState() {
	var mu Mutex

	for i := 0; i < 10000; i++ {
		go func() {
			mu.Lock()
			defer mu.Unlock()
			time.Sleep(time.Second)
		}()
	}

	time.Sleep(time.Second * 10)

	// 输出锁的信息
	fmt.Printf("waiterCount: %d, isLocked: %t, woken: %t, starving: %t\n", mu.WaiterCount(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
}

// 复制Mutex定义的常量
const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

// Mutex 扩展一个Mutex结构
type Mutex struct {
	sync.Mutex
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool {
	// 如果能成功抢到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// 尝试在竞争的状态下请求锁
	newState := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, newState)
}

// WaiterCount 获取等待者的数量
func (m *Mutex) WaiterCount() int {
	//获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v>>mutexWaiterShift + (v & mutexLocked)
	return int(v)
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}
