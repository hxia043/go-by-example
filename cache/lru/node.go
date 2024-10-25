package lru

type Node struct {
	key, value int
	Pre, Next  *Node
}
