package udb

import (
	"testing"
)

func TestTable_Example(t *testing.T) {
	table := NewTable(map[ColName]DataType{
		"int": TypeInt,
		"str": TypeString,
	})

	t.Log(table.Get("str", "!!"))

	_ = table.Insert(map[ColName]any{
		"int": 3124,
		"str": "hello",
	})
	_ = table.Insert(map[ColName]any{
		"int": 9000,
		"str": "world",
	})
	_ = table.Insert(map[ColName]any{
		"int": 3,
		"str": "!!",
	})

	//_ = table.AddIndex("str")

	_ = table.Insert(map[ColName]any{
		"int": 313138,
		"str": "abc",
	})

	t.Log(table.Get("str", "!!"))
}
