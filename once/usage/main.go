package main

import (
	"fmt"
	"sync"
)

//Once:可以用于执行且只需要执行一次的动作,常被用于单例资源的初始化

func main() {
	counter := 0
	worker := func() {
		counter++
		fmt.Println("worker exec times: ", counter)
	}

	//仅有第一次操作会被执行
	var o sync.Once
	o.Do(worker)

	o.Do(worker) //not execute
}
