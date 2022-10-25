# golang-g-base

## 面向接口

> 从 duck typing 的概念开始学起，还将探讨其他语言中对 duck typing 的支持，由此引出接口的概念。深入理解 Go 语言接口的内部实现以及使用接口实现组合的模式。

## 函数式编程

> 函数式编程是“高级”概念，但对于 Go 语言却非常基本。函数式编程的概念并且比较其他语言函数式编程的实现方法，重点理解闭包。

## 错误处理和资源管理

> Go 语言独特的 defer/panic/recover，以及错误机制，在社区有着广泛的争论。

## 测试与性能调优

> Go 语言的测试不同于其他如 junit，Go 语言采用“表格驱动测试”的理念。

## Goroutine

> 并发编程：Goroutine，协程的概念，以及背后的 Go 语言调度器。

### Coroutine(协程)

- 轻量级“协程”
- **[非抢占式]**多任务处理，由协程主动交出控制器
- 编译器/解释器/虚拟机层面的多任务
- 多个协程可能在一个活多个线程上运行

#### 代码演示

```go
func main() {
	for i := 0; i < 10; i++ {
		func(i int) {
			for {
				fmt.Printf("Hello goroutine %d\n", i)
			}
		}(i)
	}
}
```

- 上面这段代码，我们做这么一个需求：让10个不同的人帮我们无限循环打印i
- 但很可惜，代码运行后无限循环执行打印 `Hello goroutine 0.. ..Hello goroutine 0 .. ..`
- 就是说变量 `i` 永远都是0
- 对了，另外上面还有个闭包和向闭包传递指定引用的问题，不摆了，和js里一样么

> 其实上面这段代码就是让一个线程去事情，但内部的for让这个线程卡死，只能在那死循环地执行。

现在我们在匿名函数前加个 `go` 关键字，让系统给我们分配协程去做事！

```go
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				fmt.Printf("Hello goroutine %d\n", i)
			}
		}(i)
	}
}
```

- 但执行上面这段代码你会发现，控制台什么都没输出，直接就退出了！
- 直接就退出的原因是我们是并发执行了，就是说for循环里的匿名函数和main函数都是同级并发执行的
- 就是说在协程里我们想打印输出，但main函数已经for完了所有的i，然后就自己退出了！
- 在Go语言里，如果它的main函数退出，则所有的goroutine都会被干掉。
- 所以说我们肯定看不到输出

```go
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				fmt.Printf("Hello goroutine %d\n", i)
			}
		}(i)
	}
	time.Sleep(time.Millisecond)
}
```

- 先临时加个延时看下输出结果
- 可以看到结果并不是先打印10个0，再打印10个1，而是岔开的
- 意思就是不同的人在同时看活，而且是没有顺序的
- 接着继续升级

```go
func main() {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			for {
				fmt.Printf("Hello goroutine %d\n", i)
			}
		}(i)
	}
	time.Sleep(time.Minute)
}
```

- 我们开1000个协程来打印，main延迟1分钟再退出
- 接下来就会看到控制台疯狂输出，CPU基本能占满
- 现在我们就稍微能体验到Go语言并发的`无敌之处`
- 普通编程语言开10个线程，干活不算难，但开1000个线程干活就不容易
- 但我们在Go语言了却很简单就能做到
- 最关键的是他非常`human`

## Channel

> Channel 是 Goroutine 之间通信的桥梁，它和函数一样是一等公民。
