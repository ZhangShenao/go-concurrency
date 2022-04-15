package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	cancelByGoroutine()
}

func withValue() {
	ctx := context.Background()

	//WithValue会基于parent context生成一个新的context,并保存一个key-value键值对,常常用于传递上下文
	ctx1 := context.WithValue(ctx, "k1", "v1")
	ctx2 := context.WithValue(ctx1, "k2", "v2")
	fmt.Println("k2=", ctx2.Value("k2"))
}

//Context典型用法:利用context取消一个Goroutine的执行
//Context也被称为Goroutine生命周期范围（goroutine-scoped）的 Context
func cancelByGoroutine() {
	//创建一个可取消的Context
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer func() {
			fmt.Println("goroutine exit")
		}()

		fmt.Println("goroutine running...")
		for {
			select {
			case <-ctx.Done(): //监听Done Channel
				fmt.Println("goroutine done")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(time.Second)
	cancel()
	time.Sleep(time.Second * 2)
}
