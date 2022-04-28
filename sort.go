package collection

import (
	"github.com/goal-web/contracts"
	"sort"
)

func (this *Collection[T]) Len() int {
	return this.Count()
}

func (this *Collection[T]) Swap(i, j int) {
	this.rawData[i], this.rawData[j] = this.rawData[j], this.rawData[i]
}

func (this *Collection[T]) Less(i, j int) bool {
	if this.sorter != nil {
		return this.sorter(i, j)
	}
	return i > j
}

func (this *Collection[T]) SetSorter(sorter func(i, j int) bool) contracts.Collection[T] {
	this.sorter = sorter
	return this
}

// Sort sorter 必须是接收两个参数，并且返回一个 bool 值的函数
func (this *Collection[T]) Sort(sorter func(previous T, next T) bool) contracts.Collection {

	newCollection := (&Collection[T]{rawData: this.rawData}).SetSorter(func(i, j int) bool {
		return sorter(this.rawData[i], this.rawData[j])
	})

	sort.Stable(newCollection)

	return newCollection
}
