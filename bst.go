package data

import (
	"errors"
	"golang.org/x/exp/constraints"
)

type BST[T constraints.Ordered] struct {
	left  *BST[T]
	right *BST[T]
	value T
}

func (b *BST[T]) Insert(v T) (*BST[T], error) {
	var err error
	if b == nil {
		return &BST[T]{value: v}, nil
	} else if v < b.value {
		b.left, err = b.left.Insert(v)
	} else if v > b.value {
		b.right, err = b.right.Insert(v)
	} else { // v == b.value
		err = errors.New("value already exists in tree")
	}
	return b, err
}

func (b *BST[T]) Find(value T) bool {
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

func (b *BST[T]) PreOrder() []T {
	if b == nil {
		return nil
	}
	return append(append([]T{b.value}, b.left.PreOrder()...), b.right.PreOrder()...)
}

func (b *BST[T]) InOrder() []T {
	if b == nil {
		return nil
	}
	return append(append(b.left.InOrder(), b.value), b.right.InOrder()...)
}
