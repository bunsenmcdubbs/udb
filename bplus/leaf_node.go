package bplus

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

type LeafNode[K constraints.Ordered, V any] struct {
	keys   []K
	values []V
	prev   *LeafNode[K, V]
	next   *LeafNode[K, V]
}

func (n *LeafNode[K, V]) Insert(cap int, newKey K, newVal V) (new bool, split *nodeSplit[K, V]) {
	idx, exists := find(n.keys, newKey)
	if exists {
		n.values[idx] = newVal
		return false, nil
	}

	if len(n.keys) >= cap-1 { // need to split the node!
		midIdx := (len(n.keys) + 1 + 1) / 2 // add one for new element and add one to "round up"
		splitIdx := midIdx
		if idx < splitIdx {
			splitIdx -= 1
		}
		newNode := &LeafNode[K, V]{
			keys:   make([]K, len(n.keys)-splitIdx, cap),
			values: make([]V, len(n.keys)-splitIdx, cap),
			prev:   n,
			next:   n.next,
		}
		copy(newNode.keys, n.keys[splitIdx:])
		copy(newNode.values, n.values[splitIdx:])
		if n.next != nil {
			n.next.prev = newNode
		}
		n.keys = n.keys[:splitIdx]
		n.values = n.values[:splitIdx]
		n.next = newNode

		if idx < midIdx {
			n.insert(idx, newKey, newVal)
		} else {
			newNode.insert(idx-splitIdx, newKey, newVal)
		}
		return true, &nodeSplit[K, V]{newNode.keys[0], newNode}
	}

	n.insert(idx, newKey, newVal)
	return true, nil
}

func (n *LeafNode[K, V]) insert(idx int, newKey K, newVal V) {
	n.keys = insert(n.keys, idx, newKey)
	n.values = insert(n.values, idx, newVal)
}

func (n *LeafNode[K, V]) Value(key K) (V, error) {
	idx, found := find(n.keys, key)
	if !found {
		var v V
		return v, errors.New("not found")
	}
	return n.values[idx], nil
}

func (n *LeafNode[K, V]) children() []Node[K, V] {
	return nil
}

func (n *LeafNode[K, V]) debugString(sb *strings.Builder) {
	//sb.WriteString("{")
	for i := range n.keys {
		sb.WriteString(fmt.Sprintf("[%v:%v]", n.keys[i], n.values[i]))
	}
	//sb.WriteString("}")
}
