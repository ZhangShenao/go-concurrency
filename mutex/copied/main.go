package main

import (
	"fmt"
	"sync"
)

func main() {
	copyMutex()
}

//Mutex常见错误场景:copy使用过的Mutex
//Go Runtime有死锁检测机制,可以检测出因为复制Mutex而导致死锁的情况
//也可以使用vet工具,在编译期进行校验
func copyMutex() {
	c := counter{
		value: 0,
	}

	c.Lock()
	defer c.Unlock()
	c.value++
	foo(c) //fatal error: all goroutines are asleep - deadlock!

}

func foo(c counter) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("copy mutex to foo")
}

type counter struct {
	sync.Mutex
	value int
}
