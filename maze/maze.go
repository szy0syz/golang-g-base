package main

import (
	"fmt"
	"os"
)

type point struct {
	i, j int
}

var dirs = [4]point{
	// 上			 右👉🏻			 下		  👈🏻左
	{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

func walk(maze [][]int, start, end point) {
	// 它不是二维数组的意思，是个切片
	steps := make([][]int, len(maze))

	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}

	// 我靠，这里竟然是:
	// 创建一个队列，”一维数组“ 初始值是 [start]
	Q := []point{start}

	// 队列不能为空才能走下去
	// 否则就是死路了
	for len(Q) > 0 {
		// 取队首
		cur := Q[0]
		// 提取队首后面的元素
		Q = Q[1:]
	}
}

func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var row, col int
	fmt.Fscanf(file, "%d %d", &row, &col)

	maze := make([][]int, row)

	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}

	return maze
}

func main() {
	maze := readMaze("maze.in")

	for _, row := range maze {
		for _, val := range row {
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}

	walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
}
