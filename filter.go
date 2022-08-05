package collection

import (
	"github.com/goal-web/contracts"
)

// Filter 过滤不需要的数据 filter 返回 true 时保留
func (col *Collection[T]) Filter(filter func(T, int) bool) contracts.Collection[T] {
	array := make([]T, 0)

	for i, data := range col.rawData {
		if filter(data, i) {
			array = append(array, data)
		}
	}

	return New(array)
}

// Skip 过滤不需要的数据 filter 返回 true 时过滤
func (col *Collection[T]) Skip(skipper func(item T, index int) bool) contracts.Collection[T] {
	array := make([]T, 0)

	for i, data := range col.rawData {
		if !skipper(data, i) {
			array = append(array, data)
		}
	}

	return New(array)
}
