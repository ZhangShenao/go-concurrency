package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//Cond条件变量

func main() {
	//基于Locker创建一个条件变量
	c := sync.NewCond(&sync.Mutex{})

	//指定条件表达式
	ready := 0

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))

			//加锁修改条件
			c.L.Lock()
			ready++
			c.L.Unlock()

			//广播唤醒所有等待的Goroutine
			c.Broadcast()

			fmt.Printf("运动员-%d 准备就绪~\n", i)
		}(i)
	}

	c.L.Lock()
	//等待条件成立
	for ready != 10 { //wait唤醒后需要再次检查条件
		c.Wait()
		fmt.Println("裁判员被唤醒")
	}
	c.L.Unlock()

	fmt.Println("所有运动员准备就绪,比赛开始!!")

}
