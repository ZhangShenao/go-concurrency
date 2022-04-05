package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//模拟死锁:派出所等待物业先开证明,同时物业等待派出所先开证明
	//fatal error: all goroutines are asleep - deadlock!

	var policeCert sync.Mutex   //派出所证明
	var propertyCert sync.Mutex //物业证明

	var wg sync.WaitGroup
	wg.Add(2)

	//派出所Goroutine
	go func() {
		defer wg.Done() //派出所处理完成

		//派出所检查材料
		policeCert.Lock()
		defer policeCert.Unlock()
		fmt.Println("派出所检查材料")
		time.Sleep(time.Second * 2)

		//检查物业证明
		propertyCert.Lock()
		defer propertyCert.Unlock()
	}()

	//物业Goroutine
	go func() {
		defer wg.Done() //物业处理完成

		//物业检查材料
		propertyCert.Lock()
		defer propertyCert.Unlock()
		fmt.Println("物业检查材料")
		time.Sleep(time.Second * 2)

		//检查派出所证明
		policeCert.Lock()
		defer policeCert.Unlock()
	}()

	wg.Wait()
	fmt.Println("处理完成")

}
