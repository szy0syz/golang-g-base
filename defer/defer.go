package main

import (
	"bufio"
	"fmt"
	"os"
)

func writeFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 20; i++ {
		fmt.Fprintln(writer, i)
	}
}

func writeFile1(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok {
			fmt.Println("@Error", pathError.Err)
		} else {
			fmt.Println("@Unknown error", err)
		}
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 20; i++ {
		fmt.Fprintln(writer, i)
	}
}

func main() {
	// writeFile("1.txt")
	writeFile1("1.txt")
}
