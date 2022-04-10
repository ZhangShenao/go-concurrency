package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//使用sync.WaitGroup实现任务编排

func main() {
	//声明WaitGroup变量,初始值为0
	var wg sync.WaitGroup

	//设置WaitGroup计数值
	wg.Add(10)
	c := NewCounter()

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done() //任务执行完成
			c.Incr()
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		}()
	}

	//检查点,调用Wait方法阻塞等待所有Goroutine执行完成
	wg.Wait()

	fmt.Println("counter value: ", c.Get())
}

type Counter struct {
	sync.Mutex
	v int
}

func NewCounter() *Counter {
	return &Counter{
		v: 0,
	}
}

func (c *Counter) Incr() {
	c.Lock()
	defer c.Unlock()
	c.v++
}

func (c *Counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.v
}
