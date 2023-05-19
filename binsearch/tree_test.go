package binsearch

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestBST(t *testing.T) {
	testData := []int{9351, 1121, 5072, 7593, 9110, 176, 9883, 9914, 8839, 1897}
	t.Log("Original", testData)

	var bst *Tree[int]
	var err error
	for _, d := range testData {
		bst, err = bst.Insert(d)
		assert.NoError(t, err)
	}
	t.Log("PreOrder", bst.PreOrder())

	for _, d := range testData {
		assert.True(t, bst.Find(d))
	}
	assert.False(t, bst.Find(9001))

	t.Log("InOrder", bst.InOrder())
	sort.Ints(testData)
	sorted := bst.InOrder()
	assert.Equal(t, testData, sorted)
}
