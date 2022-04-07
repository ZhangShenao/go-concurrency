package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := NewRWCounter()

	//multi readers
	for i := 0; i < 10; i++ {
		go func() {
			for {
				v := c.Get()
				fmt.Println("value: ", v)
				time.Sleep(time.Millisecond * 10)
			}
		}()
	}

	//one writer
	for {
		c.Incr()
		time.Sleep(time.Millisecond * 100)
	}
}

// RWCounter 读写锁分离的计数器
type RWCounter struct {
	sync.RWMutex //embedded field
	v            int
}

func NewRWCounter() *RWCounter {
	return &RWCounter{
		v: 0,
	}
}

func (c *RWCounter) Incr() {
	//write lock
	c.Lock()
	defer c.Unlock()
	c.v++
}

func (c *RWCounter) Get() int {
	//read lock
	c.RLock()
	defer c.RUnlock()
	return c.v
}
