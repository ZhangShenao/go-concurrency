package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	//invoke function
	c := getFuncErrorAsync(func() error {
		defer wg.Done()
		fmt.Println("exec function")
		time.Sleep(2 * time.Second)
		return errors.New("timeout")
	})

	wg.Wait()

	//receive error from channel
	fmt.Println("error: ", <-c)

}

//异步获取函数的执行错误
func getFuncErrorAsync(f func() error) <-chan error {
	//create channel
	c := make(chan error)

	//invoke function in goroutine
	go func() {
		err := f()
		c <- err //send error to channel
	}()

	return c
}
