package bplus

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInternalNode_Insert(t *testing.T) {
	type args struct {
		branchingFactor int
		key             int
	}
	type wantNode struct {
		keys   []int
		values []Node[int, string]
	}
	tests := []struct {
		name          string
		node          *InternalNode[int, string]
		args          args
		wantSplit     bool
		wantNode      wantNode
		wantSplitKey  int
		wantSplitNode wantNode
	}{
		{
			name: "no splits",
			node: &InternalNode[int, string]{
				keys:   []int{10, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{}},
			},
			args: args{
				branchingFactor: 4,
				key:             15,
			},
			wantNode: wantNode{
				keys:   []int{10, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{}},
			},
		},
		{
			name: "insert new child",
			node: &InternalNode[int, string]{
				keys:   []int{10, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{willSplit: true, splitKey: 15}, fakeNode{}},
			},
			args: args{
				branchingFactor: 4,
				key:             15,
			},
			wantNode: wantNode{
				keys:   []int{10, 15, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{willSplit: true, splitKey: 15}, fakeNode{newlyCreated: true}, fakeNode{}},
			},
		},
		{
			name: "node split left",
			node: &InternalNode[int, string]{
				keys:   []int{10, 20, 30},
				values: []Node[int, string]{fakeNode{}, fakeNode{willSplit: true, splitKey: 15}, fakeNode{}, fakeNode{}},
			},
			args: args{
				branchingFactor: 4,
				key:             15,
			},
			wantNode: wantNode{
				keys:   []int{10, 15},
				values: []Node[int, string]{fakeNode{}, fakeNode{willSplit: true, splitKey: 15}, fakeNode{newlyCreated: true}},
			},
			wantSplit:    true,
			wantSplitKey: 20,
			wantSplitNode: wantNode{
				keys:   []int{30},
				values: []Node[int, string]{fakeNode{}, fakeNode{}},
			},
		},
		{
			name: "node split right",
			node: &InternalNode[int, string]{
				keys:   []int{10, 20, 30, 40},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{}, fakeNode{willSplit: true, splitKey: 37}, fakeNode{}},
			},
			args: args{
				branchingFactor: 5,
				key:             35,
			},
			wantNode: wantNode{
				keys:   []int{10, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{}},
			},
			wantSplit:    true,
			wantSplitKey: 30,
			wantSplitNode: wantNode{
				keys:   []int{37, 40},
				values: []Node[int, string]{fakeNode{willSplit: true, splitKey: 37}, fakeNode{newlyCreated: true}, fakeNode{}},
			},
		},
		{
			name: "node split center",
			node: &InternalNode[int, string]{
				keys:   []int{10, 20, 30, 40},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{willSplit: true, splitKey: 23}, fakeNode{}, fakeNode{}},
			},
			args: args{
				branchingFactor: 5,
				key:             29,
			},
			wantNode: wantNode{
				keys:   []int{10, 20},
				values: []Node[int, string]{fakeNode{}, fakeNode{}, fakeNode{willSplit: true, splitKey: 23}},
			},
			wantSplit:    true,
			wantSplitKey: 23,
			wantSplitNode: wantNode{
				keys:   []int{30, 40},
				values: []Node[int, string]{fakeNode{newlyCreated: true}, fakeNode{}, fakeNode{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, split := tt.node.Insert(tt.args.branchingFactor, tt.args.key, "blah")
			assert.Equal(t, tt.wantSplit, split != nil)
			assert.Equal(t, tt.wantNode.keys, tt.node.keys)
			assert.Equal(t, tt.wantNode.values, tt.node.values)
			if tt.wantSplit && split != nil {
				assert.Equal(t, tt.wantSplitKey, split.key)
				splitInternalNode, ok := split.node.(*InternalNode[int, string])
				if !ok {
					t.Fatal("internal split must result in an internal node")
				}
				assert.Equal(t, tt.wantSplitNode.keys, splitInternalNode.keys)
				assert.Equal(t, tt.wantSplitNode.values, splitInternalNode.values)
			}
		})
	}
}

type fakeNode struct {
	willSplit    bool
	splitKey     int
	newlyCreated bool
}

func (f fakeNode) Insert(int, int, string) (bool, *nodeSplit[int, string]) {
	if f.willSplit {
		return true, &nodeSplit[int, string]{key: f.splitKey, node: fakeNode{newlyCreated: true}}
	}
	return true, nil
}

func (f fakeNode) Value(_ int) (string, error) {
	panic("not implemented")
}

func (f fakeNode) IteratorAt(int) *Iterator[int, string] {
	panic("not implemented")
}

func (f fakeNode) children() []Node[int, string] {
	panic("not implemented")
}

func (f fakeNode) debugString(_ *strings.Builder) {
	panic("not implemented")
}
