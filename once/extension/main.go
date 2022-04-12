package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//Once的自定义扩展:
//1.捕获并返回执行错误
//2.提供方法判断是f否已经执行完成

func main() {
	var o ExtensionOnce
	counter := 0
	fmt.Println("done: ", o.Done())

	f := func() error {
		counter++
		fmt.Println("exec f, times: ", counter)
		return nil
	}

	err := o.Do(f)
	if err != nil {
		panic(err)
	}
	fmt.Println("done: ", o.Done())

	err = o.Do(f)
}

type ExtensionOnce struct {
	done       uint32 //执行标记
	sync.Mutex        //互斥锁
}

// Do 执行动作
func (o *ExtensionOnce) Do(f func() error) error {
	//原子操作,判断是否已经执行过
	//如果未执行,则调用doSlow方法
	if atomic.LoadUint32(&o.done) == 1 { //fast path
		return nil
	}
	return o.doSlow(f)
}

// Done 判断动作是否已经执行过
func (o *ExtensionOnce) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

func (o *ExtensionOnce) doSlow(f func() error) error {
	//加互斥锁
	o.Lock()
	defer o.Unlock()

	var err error
	if o.done == 0 { //双重检查
		//执行动作f,捕获执行异常
		err = f()

		//执行成功后再修改标记
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}

	return err
}
