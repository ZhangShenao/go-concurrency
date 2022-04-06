package main

import (
	"fmt"
	"sync"
)

func main() {
	q := NewSyncQueue(10)

	for i := 0; i < 20; i++ {
		go func(x int) {
			q.Enqueue(x)
			fmt.Printf("queue size: %d\n", q.Size())
		}(i)
	}

	var wg sync.WaitGroup
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			if !q.IsEmpty() {
				v := q.Dequeue()
				fmt.Printf("dequeue: %v\n", v)
			}
		}()
	}

	wg.Wait()

	fmt.Printf("finally queue size: %d\n", q.Size())
}

// SyncQueue 线程安全的队列
type SyncQueue struct {
	sync.Mutex               //互斥锁
	data       []interface{} //基于Slice维护队列内部元素
}

// NewSyncQueue 创建线程安全的队列
func NewSyncQueue(capacity int) *SyncQueue {
	return &SyncQueue{
		data: make([]interface{}, 0, capacity),
	}
}

// Enqueue 入队
func (q *SyncQueue) Enqueue(v interface{}) {
	q.Lock()
	defer q.Unlock()
	q.data = append(q.data, v)
}

// Dequeue 出队
func (q *SyncQueue) Dequeue() interface{} {
	q.Lock()

	if len(q.data) == 0 { //queue is empty
		q.Unlock()
		return nil
	}

	defer q.Unlock()

	v := q.data[0]
	q.data = q.data[1:]
	return v
}

// IsEmpty 判断队列是否为空
func (q *SyncQueue) IsEmpty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.data) == 0
}

// Size 获取队列中的元素数量
func (q *SyncQueue) Size() int {
	q.Lock()
	defer q.Unlock()
	return len(q.data)
}
