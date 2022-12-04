package main

import (
	"fmt"
)

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

func main() {
	tryRecover()
}
