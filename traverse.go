package main

type node = TNode
type nodeInfo = []node // right now it's just an array of children
type visitFn = func(node, nodeInfo)

func push(nodechan chan node, stack *[]node, n node) {
	select {
	case nodechan <- n:
	default:
		*stack = append(*stack, n)
	}
}

func pull(nodechan chan node, stack *[]node) node {
	if len(*stack) > 0 {
		n := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		return n
	}
	return <-nodechan
}

func traverse(nodechan chan node, visit visitFn, getChildren func(node) []node) {
	stack := make([]node, 0, 50)
	for {
		node := pull(nodechan, &stack)
		children := getChildren(node)
		visit(node, children)
		for _, child := range children {
			push(nodechan, &stack, child)
		}
	}
}
