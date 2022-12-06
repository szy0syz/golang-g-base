package main

import (
	"fmt"
	"os"
)

type point struct {
	i, j int
}

func (p point) add(r point) point {
	// 我们不要原地加，还是new一个出去
	// 说的事传不传指针问题
	return point{p.i + r.i, p.j + r.j}
}

// 在point上继续抽象struct 抽抽抽
func (p point) at(grid [][]int) (int, bool) {
	// 第i行是否越界
	if p.i < 0 || p.j >= len(grid) {
		return 0, false
	}

	// 第j列是否又越界了
	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}

	return grid[p.i][p.j], true
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

		for _, dir := range dirs {
			// go语言没有操作符重载或太麻烦，不如直接写个方法
			// next := cur + dir

			// 当前点加方向就是我们下一步要探索的点
			next := cur.add(dir)

			// 我们得确保 maze at next is 0 可走
			// and steps at next is 0
			// and next != start
			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}

			val, ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}

			if next == start {
				continue
			}
		}
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
