package bplus

import "golang.org/x/exp/constraints"

type Iterator[K constraints.Ordered, V any] struct {
	idx  int
	node *LeafNode[K, V]
}

func (i *Iterator[K, V]) Next() (K, V, bool) {
	if len(i.node.keys) == 0 {
		var k K
		var v V
		return k, v, true
	}

	k := i.node.keys[i.idx]
	v := i.node.values[i.idx]

	if i.idx < len(i.node.keys)-1 {
		i.idx += 1
	} else {
		i.node = i.node.next
		i.idx = 0
	}

	return k, v, i.node != nil
}

func (i *Iterator[K, V]) Prev() (K, V, bool) {
	if i.idx == -1 {
		i.idx = len(i.node.keys) - 1
	}
	k := i.node.keys[i.idx]
	v := i.node.values[i.idx]

	if i.idx > 0 {
		i.idx -= 1
	} else {
		i.node = i.node.prev
		i.idx = -1
	}

	return k, v, i.node != nil
}
