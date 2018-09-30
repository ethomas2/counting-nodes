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
	var stack []node
	for {
		node := pull(nodechan, &stack)
		children := getChildren(node)
		visit(node, children)
		for _, child := range children {
			push(nodechan, &stack, child)
		}
	}
}

func test() {
	// NUMGOROUTINES := 4
	// c := make(chan node, NUMGOROUTINES+1)
	// visit := func(n node) {
	// 	fmt.Println(n)
	// 	time.Sleep(time.Second)
	// }
	// getChildren := func(n node) []node {
	// 	return []node{2 * n, 2*n + 1}
	// }
	// for i := 0; i < NUMGOROUTINES; i++ {
	// 	go traverseTree(c, visit, getChildren)
	// }
	// c <- 1
	// time.Sleep(5 * time.Second)
}
