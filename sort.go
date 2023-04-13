package collection

import (
	"github.com/goal-web/contracts"
	"sort"
)

func (col *Collection[T]) Len() int {
	return len(col.array)
}

func (col *Collection[T]) Swap(i, j int) {
	col.array[i], col.array[j] = col.array[j], col.array[i]
}

func (col *Collection[T]) Less(i, j int) bool {
	if col.sorter != nil {
		return col.sorter(i, j)
	}
	return i > j
}

func (col *Collection[T]) SetSorter(sorter func(i, j int) bool) contracts.Collection[T] {
	col.sorter = sorter
	return col
}

// Sort sorter 必须是接收两个参数，并且返回一个 bool 值的函数
func (col *Collection[T]) Sort(sorter func(int, int, T, T) bool) contracts.Collection[T] {
	newCollection := &Collection[T]{array: col.array}
	newCollection.SetSorter(func(i, j int) bool {
		return sorter(i, j, col.array[i], col.array[j])
	})

	sort.Stable(newCollection)

	return newCollection
}
