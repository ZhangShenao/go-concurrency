package main

import (
	"fmt"
	"sync"
)

func main() {
	counter := MutexCounter{
		Id: "1",
	}

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				counter.Incr()
			}
		}()
	}

	wg.Wait()

	fmt.Println("count=", counter.Count())
}

// MutexCounter Thread Safe Counter
type MutexCounter struct {
	Id string

	mu    sync.Mutex //embedded field
	count int
}

func (c *MutexCounter) Incr() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *MutexCounter) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
