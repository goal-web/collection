package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func (col *Collection[T]) First() (T, bool) {
	if col.Len() == 0 {
		return col.Nil(), false
	}
	return col.rawData[0], true
}

func (col *Collection[T]) Last() (T, bool) {
	if col.Len() == 0 {
		return col.Nil(), false
	}

	return col.rawData[len(col.rawData)-1], true
}

func (col *Collection[T]) Prepend(items ...T) contracts.Collection[T] {
	return New(append(items, col.rawData...))
}

func (col *Collection[T]) Push(items ...T) contracts.Collection[T] {
	return New(append(col.rawData, items...))

}

func (col *Collection[T]) Pull(defaultValue ...T) (T, bool) {
	if col.Len() > 0 {
		result, exists := col.Last()
		if exists {
			col.rawData = col.rawData[:col.Len()-1]
		}
		return result, exists
	} else if len(defaultValue) > 0 {
		return defaultValue[0], false
	}

	return col.Nil(), false
}

func (col *Collection[T]) Shift(defaultValue ...T) (T, bool) {
	if col.Len() > 0 {
		result, exists := col.First()
		if exists {
			col.rawData = col.rawData[1:]
		}
		return result, exists
	} else if len(defaultValue) > 0 {
		return defaultValue[0], false
	}

	return col.Nil(), false
}

func (col *Collection[T]) Offset(index int, item T) contracts.Collection[T] {
	if col.Len() > index {
		col.rawData[index] = item
		return col
	}
	return col.Push(item)
}

func (col *Collection[T]) Put(index int, item T) contracts.Collection[T] {
	if col.Len() > index {
		return col.Clone().Offset(index, item)
	}
	return col.Push(item)
}

func (col *Collection[T]) Merge(collections ...contracts.Collection[T]) contracts.Collection[T] {
	for _, collection := range collections {
		col.rawData = append(col.rawData, collection.RawData()...)
	}

	return col
}

func (col *Collection[T]) Reverse() contracts.Collection[T] {
	for from, to := 0, len(col.rawData)-1; from < to; from, to = from+1, to-1 {
		col.rawData[from], col.rawData[to] = col.rawData[to], col.rawData[from]
	}
	return col
}

func (col *Collection[T]) Chunk(size int, handler func(collection contracts.Collection[T], page int) error) (err error) {
	total := col.Len()
	page := 1
	for err == nil && (page-1)*size <= total {
		offset := (page - 1) * size
		endIndex := size + offset
		if endIndex > total {
			endIndex = total
		}
		newCollection := New(col.rawData[offset:endIndex])

		err = handler(newCollection, page)
		page++
	}

	return
}

// Random 返回新的集合
func (col *Collection[T]) Random(size ...uint) contracts.Collection[T] {
	num := 1
	if len(size) > 0 {
		num = int(size[0])
	}
	var array []T
	if col.Len() >= num {
		for _, index := range utils.RandIntArray(0, col.Len()-1, num) {
			array = append(array, col.rawData[index])
		}
	}
	return New(array)
}
