// Package binsearch implements a binary search tree using generic types. Included as a practice/warmup exercise without
// practical purpose.
package binsearch

import (
	"errors"
	"golang.org/x/exp/constraints"
)

type Tree[T constraints.Ordered] struct {
	left  *Tree[T]
	right *Tree[T]
	value T
}

func (b *Tree[T]) Insert(v T) (*Tree[T], error) {
	var err error
	if b == nil {
		return &Tree[T]{value: v}, nil
	} else if v < b.value {
		b.left, err = b.left.Insert(v)
	} else if v > b.value {
		b.right, err = b.right.Insert(v)
	} else { // v == b.value
		err = errors.New("value already exists in tree")
	}
	return b, err
}

func (b *Tree[T]) Find(value T) bool {
	switch {
	case b == nil:
		return false
	case value < b.value:
		return b.left.Find(value)
	case value == b.value:
		return true
	case value > b.value:
		return b.right.Find(value)
	default:
		return false
	}
}

func (b *Tree[T]) PreOrder() []T {
	if b == nil {
		return nil
	}
	return append(append([]T{b.value}, b.left.PreOrder()...), b.right.PreOrder()...)
}

func (b *Tree[T]) InOrder() []T {
	if b == nil {
		return nil
	}
	return append(append(b.left.InOrder(), b.value), b.right.InOrder()...)
}
