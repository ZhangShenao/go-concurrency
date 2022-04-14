package main

import (
	"fmt"
	"sync"
)

func main() {
	concurrentReadWriteMap()
}

//map类型的基本使用
//尽量使用内建的基本类型作为map的key
//如果想要使用struct作为key，必须保证struct是immutable的

//自定义struct作为map的key
type mapkey struct {
	key string
}

//如果想要使用struct作为key，必须保证struct是immutable的,否则可能导致查询不到元素
func useStructAsKey() {
	m := make(map[mapkey]int)
	key := mapkey{
		key: "key1",
	}
	m[key] = 1

	value, existed := m[key] //推荐使用comma-ok语法获取map中的元素
	fmt.Printf("before key modified, value=%d\texisted=%v\n", value, existed)

	//修改map的key之后,value就查不到了
	key.key = "key111"
	value, existed = m[key]
	fmt.Printf("after key modified, value=%d\texisted=%v\n", value, existed)
}

//Go语言内建的map类型不是并发安全的,对map的并发读写会导致panic
//fatal error: concurrent map read and map write
func concurrentReadWriteMap() {
	m := make(map[string]int)
	var wg sync.WaitGroup
	wg.Add(2)

	//write goroutine
	go func() {
		defer wg.Done()
		for {
			m["k1"] = 1
			//time.Sleep(time.Millisecond * 20)
		}
	}()

	//read goroutine
	go func() {
		defer wg.Done()
		for {
			_ = m["k2"]
			//time.Sleep(time.Millisecond * 20)
		}
	}()

	wg.Wait()

}
