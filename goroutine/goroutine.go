package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		func(i int) {
			for {
				fmt.Printf("Hello goroutine %d\n", i)
			}
		}(i)
	}
}
