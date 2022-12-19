package bplus

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

type InternalNode[K constraints.Ordered, V any] struct {
	keys   []K
	values []Node[K, V]
}

func (n *InternalNode[K, V]) child(key K) Node[K, V] {
	idx, exists := find(n.keys, key)
	if !exists {
		return n.values[idx]
	}
	return n.values[idx+1]
}

func (n *InternalNode[K, V]) Insert(cap int, newKey K, newVal V) (bool, *nodeSplit[K, V]) {
	newWrite, split := n.child(newKey).Insert(cap, newKey, newVal)
	if split == nil {
		return newWrite, nil
	}

	idx, _ := find(n.keys, split.key)
	n.keys = insert(n.keys, idx, split.key)
	n.values = insert(n.values, idx+1, split.node)

	if len(n.values) > cap { // split
		midIdx := (len(n.keys) + 1) / 2

		newNode := &InternalNode[K, V]{
			keys:   make([]K, len(n.keys)-midIdx, cap),
			values: make([]Node[K, V], len(n.values)-midIdx, cap),
		}
		copy(newNode.keys, n.keys[midIdx:])
		copy(newNode.values, n.values[midIdx:])

		promotedKey := n.keys[midIdx-1]
		n.keys = n.keys[:midIdx-1]
		n.values = n.values[:midIdx]
		return true, &nodeSplit[K, V]{promotedKey, newNode}
	}
	return true, nil
}

func (n *InternalNode[K, V]) Value(key K) (V, error) {
	return n.child(key).Value(key)
}

func (n *InternalNode[K, V]) children() []Node[K, V] {
	return n.values
}

func (n *InternalNode[K, V]) debugString(sb *strings.Builder) {
	sb.WriteString("{")
	for i, k := range n.keys {
		n.values[i].debugString(sb)
		sb.WriteString(fmt.Sprintf(" (%v) ", k))
	}
	n.values[len(n.values)-1].debugString(sb)
	sb.WriteString("}")
}
