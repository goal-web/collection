package collection

import (
	"github.com/goal-web/contracts"
	"sort"
)

func (col *Collection[T]) SetSorter(sorter func(i, j int) bool) contracts.Collection[T] {
	col.sorter = sorter
	return col
}

// Sort sorter 必须是接收两个参数，并且返回一个 bool 值的函数
func (col *Collection[T]) Sort(sorter func(previous T, next T) bool) contracts.Collection[T] {
	newCollection := col.SetSorter(func(i, j int) bool {
		return sorter(col.rawData[i], col.rawData[j])
	})

	sort.Stable(newCollection)

	return newCollection
}
