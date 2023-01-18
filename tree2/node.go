package main

import "fmt"

type treeNode struct {
	value       int
	left, right *treeNode
}

// 因为go里没有构造函数一说
// 但我们可以用工厂函数控制构造
func createNode(value int) *treeNode {
	// 这里很明显是个局部变量，但go里可以返回出去，c++肯定不行
	return &treeNode{value: value}
	// 肯定不用考虑分配在堆上还是栈上
}

// 第一个括号里的叫 接收者
// 这个print是给这个node接收的
func (node treeNode) print() {
	fmt.Print(node.value, " ")
}

// node treeNode
// ⭐️如果这里传值不穿指针，那么我们所盖的值永远改不到节点树上
func (node *treeNode) setValue(value int) {
	node.value = value
}

func (node *treeNode) traverse() {
	if node == nil {
		return
	}
	node.left.traverse()
	node.print()
	node.right.traverse()
}

func main() {
	var root treeNode

	root = treeNode{value: 3}
	root.left = &treeNode{}
	// 不写name，就是按顺序填
	root.right = &treeNode{5, nil, nil}
	root.right.left = new(treeNode)
	root.left.right = createNode(2) // 工厂函数

	//nodes := []treeNode{
	//	{value: 3},
	//	{},
	//	{6, nil, &root},
	//}

	root.right.left.setValue(4)
	//root.left.print()
	//root.print()

	root.traverse()
}
