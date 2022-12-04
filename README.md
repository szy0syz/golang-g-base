# golang-g-base

## 面向接口

> 从 duck typing 的概念开始学起，还将探讨其他语言中对 duck typing 的支持，由此引出接口的概念。深入理解 Go 语言接口的内部实现以及使用接口实现组合的模式。

## 函数式编程

> 函数式编程是“高级”概念，但对于 Go 语言却非常基本。函数式编程的概念并且比较其他语言函数式编程的实现方法，重点理解闭包。

## 错误处理和资源管理

> Go 语言独特的 defer/panic/recover，以及错误机制，在社区有着广泛的争论。

## 测试与性能调优

> Go 语言的测试不同于其他如 junit，Go 语言采用“表格驱动测试”的理念。
>
> Slogan: Debugging Sucks! Testing Rocks!

### 传统测试

```java
@TGest public void testAdd() {
  assertEquals(3, add(1, 2));
  assertEquals(2, add(0, 2));
  assertEquals(0, add(0, 0));
  assertEquals(0, add(-1, 1));
  assertEquals(Integer.MIN_VALUE, add(1, Integer,MAX_VALUE));
}
```

- 测试数据和测试逻辑混在一起
- 出错信息不明确

### 表格驱动测试

```go
tests := []struct {
  a, b, c int32
} {
 {1,2,3},
 {0,2,2},
 {0,0,0},
 {-1,1,0},
 {math.MaxInt32,1,math.MinInt32},
}

for _, test := range tests {
 if actual := add(test.a, test.b); actual != test.c {}
}
```

- 分离的测试数据和测试逻辑
- 明确的出错信息
- 可以部分失败
- go语言的语法使得我们更容易实践表格驱动测试

## 8.错误处理和资源管理

- go是用defer调用来管理资源的
- 确保调用在函数结束时发生
  - 就是说defer后面的代码一定会再函数结束前才会发生
- 何时使用defer调用
  - open/close
  - lock/unlock
  - PrintHeader/PrintFooter

```go
func writeFile(filename string) {
 file, err := os.Create(filename)
 if err != nil {
  panic(err)
 }
 defer file.Close()

 writer := bufio.NewWriter(file)
 defer writer.Flush()
 // 注意：defer是先进后出的栈

 for i := 0; i < 20; i++ {
  fmt.Fprintln(writer, i)
 }
}
```

> 这小段代码把`defer`的使用场景和go的设计精妙都体现出来了

### 错误处理

> 知道到底是什么样的错误，进异步处理它。

```go
file, err := os.Open(filename)
 if err != nil {
  if pathError, ok := err.(*os.PathError); ok {
   fmt.Println("@Error", pathError.Err)
  } else {
   fmt.Println("@Unknown error", err)
  }
  return
 }
```

### list-web-server

```go
func main() {
 http.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path[len("/list/"):] // /list/readme.txt

  file, err := os.Open(path)
  if err != nil {
   http.Error(w, err.Error(), http.StatusInternalServerError)
   return
  }
  defer file.Close()

  all, err := ioutil.ReadAll(file)
  if err != nil {
   panic(err)
  }

  w.Write(all)
 })

 err := http.ListenAndServe(":8888", nil)
 if err != nil {
  panic(err)
 }
}
```

- 就目前而言Error处理是业务逻辑，应该和运行流程无关，提出来

### 统一错误处理

```go
type appHandler func(w http.ResponseWriter, r *http.Request) error

func errWrapper(handler appHandler) func(w http.ResponseWriter, r *http.Request) {
 return func(w http.ResponseWriter, r *http.Request) {
  err := handler(w, r)
  if err != nil {
   log.Println("Error handling request:", err.Error())
   code := http.StatusOK
   switch {
   case os.IsNotExist(err): // 文件不存在的错误
    code = http.StatusNotFound
   case os.IsPermission(err): // 没权限读
    code = http.StatusForbidden
   default:
    code = http.StatusInternalServerError
   }
   http.Error(w, http.StatusText(code), code)
  }
 }
}
```

- 函数式编程真不错

### panic

- 停止当前函数执行
- 一直向上返回，执行每一次的defer
- 如果没遇到recover，程序就退出

### recover

- 仅在defer调用中使用
- 获取panic的值
- 如果无法处理，可重新panic

```go
func tryRecover() {
 defer func() {
  r := recover()
  if err, ok := r.(error); ok {
   fmt.Println("Error occurred:", err)
  } else {
   panic(fmt.Sprintf("I don't know what to do: %v", r))
  }
 }()
 // panic(errors.New("it is a big error"))
 // zero := 0
 // a := 1 / zero
 // println(a)
 panic("123")
}
```

### error VS panic

- 意料之中的: 使用error —— 如文件打不开
  - 你能想到的错误都是error，不能是panic，用panic太崩溃
- 意料之外的: 使用panic —— 如数组越界

### 综合实例

- 目前web.go代码中虽然有http的自带recover，但我们想加一些自定义Error，然后自己保护一下
- 那我们就在errWrapper里先做自己的保护
- Go语言牛的地方在于同一个接口，在两个地方使用，可以不引用，结构一样即可

## Goroutine

> 并发编程：Goroutine，协程的概念，以及背后的 Go 语言调度器。

### Coroutine(协程)

- 轻量级“协程”
  - 随手开个 1000 个协程没问题，但线程不行 🚫
- **[非抢占式]**多任务处理，由协程主动交出控制器
  - 什么交非抢占式呢？就是说由协程自己主动交出控制权
  - 但是线程的话都会被操作系统给切换，所以说线程是抢占式的任务处理
  - 线程有可能做到一半或某个语句执行到一半都会被操作系统在其中间掐掉，然后转到其他进程线程里做其他事，最后操作系统又会折回头来执行
  - 从这里就可以看出，**线程/抢占式**最大的问题就是执行到一半，操作系统需要把上下文都存下来，占用资源
  - 但非抢占式的只需要头尾切换的几个点就可以了，对资源的消耗就小很多
- 编译器/解释器/虚拟机层面(级别)的多任务，并非运行在操作系统层面
  - 这里 go 语言会有调度器去专门调度协程
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

- 上面这段代码，我们做这么一个需求：让 10 个不同的人帮我们无限循环打印 i
- 但很可惜，代码运行后无限循环执行打印 `Hello goroutine 0.. ..Hello goroutine 0 .. ..`
- 就是说变量 `i` 永远都是 0
- 对了，另外上面还有个闭包和向闭包传递指定引用的问题，不摆了，和 js 里一样么

> 其实上面这段代码就是让一个线程去事情，但内部的 for 让这个线程卡死，只能在那死循环地执行。

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
- 直接就退出的原因是我们是并发执行了，就是说 for 循环里的匿名函数和 main 函数都是同级并发执行的
- 就是说在协程里我们想打印输出，但 main 函数已经 for 完了所有的 i，然后就自己退出了！
- 在 Go 语言里，如果它的 main 函数退出，则所有的 goroutine 都会被干掉。
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
- 可以看到结果并不是先打印 10 个 0，再打印 10 个 1，而是岔开的
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

- 我们开 1000 个协程来打印，main 延迟 1 分钟再退出
- 接下来就会看到控制台疯狂输出，CPU 基本能占满
- 现在我们就稍微能体验到 Go 语言并发的`无敌之处`
- 普通编程语言开 10 个线程，干活不算难，但开 1000 个线程干活就不容易
- 但我们在 Go 语言了却很简单就能做到
- 最关键的是他非常`human`

#### 证明 goroutine 是非抢占式的

> IO 操作会触发协程的切换

```go
func main() {
 var a [10]int
 for i := 0; i < 10; i++ {
  go func(i int) {
   for {
    a[i]++
   }
  }(i)
 }
 time.Sleep(time.Millisecond)
 fmt.Println(a)
}
```

- 证明思路：

```go
// 刚开始调度器调度协程时，差异次数会很大，但把时间拉长，会相对均匀一些。
time.Sleep(time.Millisecond)
// [516 221 155 271 159 209 243 179 282 218]
time.Sleep(time.Millisecond * 100)
// [75239 76895 85530 77206 76782 80617 82376 82048 77666 83836]
time.Sleep(time.Millisecond * 1000)
// [822655 831825 832946 849206 811773 817287 846738 807789 848445 821164]
```

### goroutine的定义

- 任何函数只需要加上go就能送给调度器运行
- 不需要再定义时区分是否是异步函数
  - 像某些语言Python必须是异步函数才能是协程
  - 但Go语言不需要
- 调度器在核实的点进行切换
  - 调度器会自己判断
  - 对了，协程了有IO操作才能交出控制器
  - 就是说在固定的点由调度器自己切换
  - 这个点和node.js有点不一样，但原理一样
  - 还有就是传统的协程是必须显示的写出切换点，而goroutine不需要
- 使用 -race 来检测数据访问冲突

### goroutine可能得切换点

- I/O, select
- channel
- 等待锁
- 函数调用(偶尔)
- runtime.GoSched()

## Channel

<img width="507" alt="image" src="https://user-images.githubusercontent.com/10555820/205502561-afd9a355-c824-47c0-9f52-05649dc59f94.png">

> Channel 是 Goroutine 之间通信的桥梁，它和函数一样是一等公民。

```go
func chanDemo() {
 c := make(chan int)
 c <- 1 // 我们发一个数据，没人收，就会死锁，deadlock
 c <- 2

 n := <-c
 fmt.Println(n)
}
// 🚫 fatal error: all goroutines are asleep - deadlock!
// chan 必须和 chan 相连
```

```go
func chanDemo() {
 c := make(chan int)
 go func() {
  for {
   n := <-c
   fmt.Println(n)
  }
 }()
 c <- 1
 c <- 2
}
// ok 可以收到了
// 1
// 2
```

```go
func worker(id int, c chan int) {
 for {
  n := <-c
  fmt.Printf("Worker %d received: %d\n", id, n)
 }
}

func chanDemo() {
 c := make(chan int)
 go worker(0, c)
 c <- 1
 c <- 2
}


// 进一步优化，把函数提出去
// 协程不一定是否异步函数，这一点简直绝绝子
```

```go
func chanDemo() {
 var channels [10]chan int

 // 这里只是创建chan
 for i := 0; i < 10; i++ {
  channels[i] = make(chan int)
  go worker(0, channels[i])
 }

 for i := 0; i < 10; i++ {
  channels[i] <- 'a' + i
 }

 for i := 0; i < 10; i++ {
  channels[i] <- 'A' + i
 }

 time.Sleep(time.Millisecond)
}
```

```bash
$ go run .
Worker 0 received: i
Worker 0 received: f
Worker 0 received: g
Worker 0 received: a
Worker 0 received: A
Worker 0 received: b
Worker 0 received: c
Worker 0 received: h
Worker 0 received: j
Worker 0 received: C
Worker 0 received: d
Worker 0 received: D
Worker 0 received: B
```

> 其实worker可以自己返回一个chan

```go
func createWorker(id int) chan<- int {
 c := make(chan int)
 go func() {
  for {
   fmt.Printf("Worker %d received: %c\n", id, <-c)
  }
 }()
 return c
}
```

- 限定我们的chan只能发数据，worker里就只能接数据，非常安全
- 如果我们在外边想收数据，那就编译错误了
- 我们可以限制chan的数据流向，到底是只能收，还是只能发，还是双向的

> 我们应该告诉外面调用createWorker的人怎么用我们的chan

### bufferedChannel

```go
func bufferedChannel() {
 c := make(chan int, 3)

 c <- 1
 c <- 2
 c <- 3
}

// 如果这样，缓冲区能接住，就不会deadlock
```

- Channel是一等公民
- Buffered Channel
- Channel Close
- Channel为什么要被设计成这样的？
  - 理论基础： Communication Sequential Process (CSP) 👍
  - Go语言的并发就是基础CSP模型来实现的

> ⭐️⭐️⭐️ channel实际应用中的注意点：(和以前的编程语言有非常大的区别！)

- Don't communicate by sharing memory; share memory by communicating.
- 不要通过共享内存来通信；通过通信来共享内存。

### 使用 Channel 进行树的遍历
