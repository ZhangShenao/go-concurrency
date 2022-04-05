package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	reentrant()
}

//Mutex常见错误场景:Lock/Unlock不是成对出现
func unlock() {
	var mu sync.Mutex
	fmt.Println("do some processing...")
	time.Sleep(time.Second)
	mu.Unlock() //panic: unlock of unlocked mutex
}

//Mutex常见错误场景:重入
//Mutex是不可重入锁
func reentrant() {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("do something...")
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("do something again...")
}
