package bplus

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTree_Insert(t *testing.T) {
	data := []struct {
		int
		string
	}{
		{60, "doing"},
		{10, "hello"},
		{20, "there."},
		{67, "fine"},
		{0, "oh!"},
		{65, "on"},
		{66, "this"},
		{40, "you"},
		{70, "day?"},
		{30, "how"},
		{35, "are"},
	}

	tr := NewTree[int, string](4)
	for _, p := range data {
		tr.Insert(p.int, p.string)
	}
	t.Log(tr.debugString())

	iterator := tr.Iterator()
	sb := &strings.Builder{}
	var keys []int
	for {
		k, v, ok := iterator.Next()
		keys = append(keys, k)
		sb.WriteString(v + " ")
		if !ok {
			break
		}
	}
	t.Log(sb.String())

	assert.IsIncreasing(t, keys)
	assert.Equal(t, len(data), len(keys))
}

func Test_find(t *testing.T) {
	type args struct {
		ordered []int
		target  int
	}
	type want struct {
		idx   int
		found bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			args: args{[]int{}, 100},
			want: want{0, false},
		},
		{
			args: args{[]int{0, 1}, 0},
			want: want{0, true},
		},
		{
			args: args{[]int{0, 1}, 1},
			want: want{1, true},
		},
		{
			args: args{[]int{0, 1}, 100},
			want: want{2, false},
		},
		{
			args: args{[]int{0, 1, 10}, 10},
			want: want{2, true},
		},
		{
			args: args{[]int{0, 1, 10}, 5},
			want: want{2, false},
		},
		{
			args: args{[]int{0, 1, 4, 7, 10, 13, 15}, 13},
			want: want{5, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIdx, gotFound := find(tt.args.ordered, tt.args.target); gotIdx != tt.want.idx || gotFound != tt.want.found {
				t.Errorf("find() = {%v %v}, want %v", gotIdx, gotFound, tt.want)
			}
		})
	}
}
