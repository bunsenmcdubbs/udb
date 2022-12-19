package bplus

import (
	"golang.org/x/exp/constraints"
	"strings"
)

type Tree[K constraints.Ordered, V any] struct {
	b    int // branching factor
	root Node[K, V]
}

func NewTree[K constraints.Ordered, V any](branchingFactor int) *Tree[K, V] {
	return &Tree[K, V]{
		b: branchingFactor,
		root: &LeafNode[K, V]{
			keys:   make([]K, 0, branchingFactor),
			values: make([]V, 0, branchingFactor),
		},
	}
}

func (t *Tree[K, V]) Insert(key K, value V) {
	_, newNode := t.root.Insert(t.b, key, value)
	if newNode != nil {
		t.root = &InternalNode[K, V]{[]K{newNode.key}, []Node[K, V]{t.root, newNode.node}}
	}
}

func (t *Tree[K, V]) Iterator() *Iterator[K, V] {
	curr := t.root
	for {
		switch node := curr.(type) {
		case *LeafNode[K, V]:
			return &Iterator[K, V]{node: node}
		case *InternalNode[K, V]:
			curr = node.values[0]
		}
	}
}

func (t *Tree[K, V]) debugString() string {
	sb := &strings.Builder{}
	t.root.debugString(sb)
	return sb.String()
}

type nodeSplit[K constraints.Ordered, V any] struct {
	key  K
	node Node[K, V]
}

type Node[K constraints.Ordered, V any] interface {
	Insert(cap int, key K, value V) (new bool, split *nodeSplit[K, V])
	Value(K) (V, error)
	children() []Node[K, V]
	debugString(*strings.Builder)
}

// find returns the index in the ordered array where the target element should be inserted
// to maintain order and whether the element is already present in the array or not.
func find[T constraints.Ordered](ordered []T, target T) (int, bool) {
	left, right := 0, len(ordered)-1
	if len(ordered) == 0 || target < ordered[left] {
		return 0, false
	}
	if target > ordered[right] {
		return len(ordered), false
	}
	for {
		mid := (left + right) / 2
		switch {
		case target < ordered[mid]:
			right = mid - 1
		case target == ordered[mid]:
			return mid, true
		case target > ordered[mid]:
			left = mid + 1
		}
		if left > right {
			return left, false
		}
	}
}

func insert[T any](arr []T, idx int, elem T) []T {
	arr = append(arr, elem)
	copy(arr[idx+1:], arr[idx:])
	arr[idx] = elem
	return arr
}
