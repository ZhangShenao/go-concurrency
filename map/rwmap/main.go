package main

import "sync"

func main() {
	//使用并发安全的RWMap
	rwMap := NewRWMap()
	var wg sync.WaitGroup
	wg.Add(2)

	//write goroutine
	go func() {
		defer wg.Done()
		for {
			rwMap.Put(1, 1)
		}
	}()

	//read goroutine
	go func() {
		defer wg.Done()
		for {
			_, _ = rwMap.Get(2)
		}
	}()

	wg.Wait()
}

// RWMap 基于RWMutex实现的并发安全的map类型,可以实现读写锁的分离
type RWMap struct {
	sync.RWMutex
	m map[int]int
}

func NewRWMap() *RWMap {
	return &RWMap{
		m: make(map[int]int),
	}
}

// Get 根据Key查找元素
func (m *RWMap) Get(k int) (int, bool) {
	//加读锁
	m.RLock()
	defer m.RUnlock()

	value, existed := m.m[k]
	return value, existed
}

// Put 保存元素
func (m *RWMap) Put(k, v int) {
	//加写锁
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

// Delete 删除指定元素
func (m *RWMap) Delete(k int) {
	//加写锁
	m.Lock()
	defer m.Unlock()

	//调用内建的delete函数
	delete(m.m, k)
}

// Len 获取map的长度
func (m *RWMap) Len() int {
	//加读锁
	m.RLock()
	defer m.RUnlock()

	return len(m.m)
}

// ForEach 遍历Map
func (m *RWMap) ForEach(f func(k, v int) bool) {
	//加读锁
	m.RLock()
	//直到遍历完再释放锁
	defer m.RUnlock()

	//遍历
	for key, val := range m.m {
		if !f(key, val) {
			return
		}
	}
}
