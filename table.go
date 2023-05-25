package udb

import (
	"errors"
	"fmt"
	"github.com/bunsenmcdubbs/udb/bplus"
)

type rowID uint64
type ColName string

type Table struct {
	nextRowID  rowID
	rows       *bplus.Tree[rowID, []any]
	indexes    map[ColName]*bplus.Tree[string, rowID]
	schema     map[ColName]DataType
	fieldOrder []ColName
}

func NewTable(schema map[ColName]DataType) *Table {
	var fieldOrder []ColName
	for col := range schema {
		fieldOrder = append(fieldOrder, col)
	}
	return &Table{
		rows:       bplus.NewTree[rowID, []any](10),
		indexes:    make(map[ColName]*bplus.Tree[string, rowID]),
		schema:     schema,
		fieldOrder: fieldOrder,
	}
}

func (t *Table) AddIndex(name ColName) error {
	if _, exists := t.indexes[name]; exists {
		return nil
	}
	fieldIdx := t.fieldIdx(name)
	if fieldIdx == -1 {
		return errors.New("no such column")
	}

	idx := bplus.NewTree[string, rowID](10)
	iter := t.rows.Iterator()
	id, row, done := iter.Next()
	for ; !done; id, row, done = iter.Next() {
		idx.Insert(serialize(row[fieldIdx]), id)
	}
	t.indexes[name] = idx
	return nil
}

func (t *Table) fieldIdx(name ColName) int {
	for idx, col := range t.fieldOrder {
		if col == name {
			return idx
		}
	}
	return -1
}

func (t *Table) Insert(row map[ColName]any) error {
	if len(row) != len(t.fieldOrder) {
		return errors.New("incorrect number of fields")
	}

	tup := make([]any, len(row))
	for fieldIdx, colName := range t.fieldOrder {
		tup[fieldIdx] = row[colName]
	}

	id := t.nextRowID
	t.rows.Insert(id, tup)
	t.nextRowID += 1

	if len(t.indexes) != 0 {
		for col, idx := range t.indexes {
			idx.Insert(serialize(tup[t.fieldIdx(col)]), id)
		}
	}
	return nil
}

func (t *Table) Get(col ColName, want any) ([]any, error) {
	if _, ok := t.schema[col]; !ok {
		return nil, errors.New("unknown column")
	}

	// index lookup if available
	if idx, ok := t.indexes[col]; ok {
		id, err := idx.Value(serialize(want))
		if err != nil {
			return nil, err
		}
		return t.rows.Value(id)
	}

	// full table scan as a fallback
	iter := t.rows.Iterator()
	fieldIdx := t.fieldIdx(col)
	_, row, done := iter.Next()
	for ; !done; _, row, done = iter.Next() {
		if row[fieldIdx] == want {
			return row, nil
		}
	}
	return nil, nil
}

// TODO add vectorized versions of all operations

func serialize(a any) string {
	if ser, ok := a.(fmt.Stringer); ok {
		ser.String()
	}
	panic("unknown type for serialization")
}
