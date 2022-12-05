package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	// 我们先来一个chan，这个chan返回int
	out := make(chan int)
	// 上goroutine
	// 告诉编译器这个是个协程，记得给调度器
	go func() {
		// 这里可是有闭包
		i := 0
		for {
			// 延迟哈这个协程
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			// 发射！！！
			out <- i
			// 走人，下个轮回见
			i++
		}
	}()
	return out
}

func main() {
	// 生成两个g，会不停的产生信号
	var c1, c2 = generator(), generator()
	// 主程序阻塞式的接收chan消息
	// 如果 c1就打印c1收到，反之
	for {
		select {
		case n := <-c1:
			fmt.Println("Received from c1:", n)
		case n := <-c2:
			fmt.Println("Received from c2:", n)
		}
	}
}
