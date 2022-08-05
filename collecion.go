package collection

import (
	"encoding/json"
	"github.com/goal-web/contracts"
	"reflect"
)

func New[T any](rawData []T) contracts.Collection[T] {
	return Make(rawData)
}

func Make[T any](rawData []T) *Collection[T] {
	return &Collection[T]{rawData: rawData}
}

type Collection[T any] struct {
	rawData []T
	values  map[int]reflect.Value
	sorter  func(i, j int) bool
}

func (col *Collection[T]) IsEmpty() bool {
	return col.Len() == 0
}
func (col *Collection[T]) Nil() T {
	var result T
	return result
}
func (col *Collection[T]) Clone() contracts.Collection[T] {
	newCollection := New([]T{})

	for _, item := range col.RawData() {
		newCollection = newCollection.Push(item)
	}

	return newCollection
}

func (col *Collection[T]) RawData() []T {
	return col.rawData
}

func (col *Collection[T]) Map(f func(item T, index int) T) contracts.Collection[T] {
	for i, data := range col.rawData {
		col.rawData[i] = f(data, i)
	}
	return col
}

func (col *Collection[T]) Each(f func(item T, index int)) contracts.Collection[T] {
	for i, data := range col.rawData {
		f(data, i)
	}
	return col
}

func (col *Collection[T]) ToAnyArray() []any {
	array := make([]any, col.Len())
	for i, data := range col.rawData {
		array[i] = data
	}
	return array
}

func (col *Collection[T]) ToJson() any {
	array := make([]any, col.Len())
	for i, data := range col.rawData {
		array[i] = data
	}

	return array
}

func (col *Collection[T]) ToJsonString() string {
	bytes, _ := json.Marshal(col.ToJson())
	return string(bytes)
}

func (col *Collection[T]) Len() int {
	return len(col.rawData)
}

func (col *Collection[T]) Less(i, j int) bool {
	if col.sorter != nil {
		return col.sorter(i, j)
	}
	return i > j
}

func (col *Collection[T]) Swap(i, j int) {
	col.rawData[i], col.rawData[j] = col.rawData[j], col.rawData[i]
}
