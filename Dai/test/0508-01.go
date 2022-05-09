package main

import (
	"fmt"
	"sync"
)

func loop() {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", i)
	}
}

// main本身也是一个go routine
// func main() {
// 	// go loop() // 启动一个goroutine
// 	// loop()
// 	// loop()
// 	// make用于初始化一个slice map chan（channel，用于并发编程中多个goroutine的通信）
// 	messages := make(chan int)
// 	go func() {
// 		messages <- 1
// 	}()
// 	// expression in go must be function call
// 	// 上一行，所以必须加括号
// 	msg := <-messages
// 	fmt.Println(msg)

// }

// 关于waitgroup：https://zhuanlan.zhihu.com/p/344973865
func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			println("hello")
		}()
	}

	wg.Wait()
}
