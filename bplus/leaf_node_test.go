package bplus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLeafNode_Insert(t *testing.T) {
	type args struct {
		branchingFactor int
		key             int
		value           string
	}
	type wantNode struct {
		keys   []int
		values []string
	}
	tests := []struct {
		name      string
		args      args
		node      *LeafNode[int, string]
		wantNew   bool
		wantSplit bool
		wantNode1 wantNode
		wantNode2 wantNode
	}{
		{
			name: "insert without filling node",
			node: &LeafNode[int, string]{
				keys:   []int{-48, 30},
				values: []string{"abcdefg", "world"},
			},
			args: args{
				branchingFactor: 4,
				key:             4,
				value:           "hello",
			},
			wantNew: true,
			wantNode1: wantNode{
				keys:   []int{-48, 4, 30},
				values: []string{"abcdefg", "hello", "world"},
			},
		},
		{
			name: "update existing value",
			node: &LeafNode[int, string]{
				keys:   []int{-100, 30},
				values: []string{"hello", "world"},
			},
			args: args{
				branchingFactor: 4,
				key:             -100,
				value:           "hello_123",
			},
			wantNode1: wantNode{
				keys:   []int{-100, 30},
				values: []string{"hello_123", "world"},
			},
		},
		{
			name: "split node; b=5; new data on left",
			node: &LeafNode[int, string]{
				keys:   []int{10, 20, 30, 40},
				values: []string{"element 10", "element 20", "element 30", "element 40"},
			},
			args: args{
				branchingFactor: 5,
				key:             25,
				value:           "new element!!!",
			},
			wantNew:   true,
			wantSplit: true,
			wantNode1: wantNode{
				keys:   []int{10, 20, 25},
				values: []string{"element 10", "element 20", "new element!!!"},
			},
			wantNode2: wantNode{
				keys:   []int{30, 40},
				values: []string{"element 30", "element 40"},
			},
		},
		{
			name: "split node; b=5; new data on right",
			node: &LeafNode[int, string]{
				keys:   []int{10, 20, 30, 40},
				values: []string{"element 10", "element 20", "element 30", "element 40"},
			},
			args: args{
				branchingFactor: 5,
				key:             35,
				value:           "new element!!!",
			},
			wantNew:   true,
			wantSplit: true,
			wantNode1: wantNode{
				keys:   []int{10, 20, 30},
				values: []string{"element 10", "element 20", "element 30"},
			},
			wantNode2: wantNode{
				keys:   []int{35, 40},
				values: []string{"new element!!!", "element 40"},
			},
		},
		{
			name: "split node; b=6; new data on left",
			node: &LeafNode[int, string]{
				keys:   []int{10, 20, 30, 40, 50},
				values: []string{"element 10", "element 20", "element 30", "element 40", "element 50"},
			},
			args: args{
				branchingFactor: 6,
				key:             25,
				value:           "new element!!!",
			},
			wantNew:   true,
			wantSplit: true,
			wantNode1: wantNode{
				keys:   []int{10, 20, 25},
				values: []string{"element 10", "element 20", "new element!!!"},
			},
			wantNode2: wantNode{
				keys:   []int{30, 40, 50},
				values: []string{"element 30", "element 40", "element 50"},
			},
		},
		{
			name: "split node; b=6; new data on right",
			node: &LeafNode[int, string]{
				keys:   []int{10, 20, 30, 40, 50},
				values: []string{"element 10", "element 20", "element 30", "element 40", "element 50"},
			},
			args: args{
				branchingFactor: 6,
				key:             35,
				value:           "new element!!!",
			},
			wantNew:   true,
			wantSplit: true,
			wantNode1: wantNode{
				keys:   []int{10, 20, 30},
				values: []string{"element 10", "element 20", "element 30"},
			},
			wantNode2: wantNode{
				keys:   []int{35, 40, 50},
				values: []string{"new element!!!", "element 40", "element 50"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newWrite, split := tt.node.Insert(tt.args.branchingFactor, tt.args.key, tt.args.value)
			assert.Equal(t, tt.wantNew, newWrite)
			assert.Equal(t, tt.wantSplit, split != nil, "unexpected node split")
			assert.EqualValues(t, tt.wantNode1.keys, tt.node.keys)
			assert.EqualValues(t, tt.wantNode1.values, tt.node.values)
			if tt.wantSplit { // hack because generics with different nil aren't equal
				assert.Equal(t, split.node, tt.node.next, "unexpected next node")
				if split != nil {
					assert.EqualValues(t, tt.wantNode2.keys, tt.node.next.keys)
					assert.EqualValues(t, tt.wantNode2.values, tt.node.next.values)
				}
			}
		})
	}
}

func TestLeafNode_InsertUpdate(t *testing.T) {
	b := 5
	l := &LeafNode[int, string]{}
	inputs := []struct {
		k int
		v string
	}{
		{1, "hello"},
		{3, "michael"},
		{2, "world"},
		{0, "oh"},
	}
	for _, in := range inputs {
		new, split := l.Insert(b, in.k, in.v)
		assert.True(t, new, "should add new key")
		assert.Nil(t, split, "should not split node")
	}

	got, err := l.Value(1)
	assert.Equal(t, "hello", got)
	assert.NoError(t, err)

	l.Insert(b, 1, "hello 4")
	got, err = l.Value(1)
	assert.Equal(t, "hello 4", got)
	assert.NoError(t, err)

	got, err = l.Value(4)
	assert.Equal(t, "", got)
	assert.Error(t, err)
}
