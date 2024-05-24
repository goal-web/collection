package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func (col *Collection[T]) Pluck(key string) map[string]T {
	fields := map[string]T{}

	for index, data := range col.ToArrayFields() {
		var name, ok = data[key].(string)
		if _, exists := fields[name]; ok && !exists {
			fields[name] = col.array[index]
		}
	}

	return fields
}

func (col *Collection[T]) GroupBy(key string) map[string][]T {
	list := map[string][]T{}

	for index, data := range col.ToArrayFields() {
		var value, _ = data[key].(string)
		list[value] = append(list[value], col.array[index])
	}

	return list
}

func (col *Collection[T]) Only(keys ...string) contracts.Collection[T] {
	rawResults := make([]T, 0)

	for index, data := range col.ToArrayFields() {
		fields := contracts.Fields{}
		for key, value := range data {
			if utils.IsIn(key, keys) {
				fields[key] = value
			}
		}
		rawResults = append(rawResults, col.array[index])
	}

	return New(rawResults)
}

func (col *Collection[T]) First() (T, bool) {
	var result T
	if col.Count() == 0 {
		return result, false
	}
	return col.array[0], true
}

func (col *Collection[T]) Last() (T, bool) {
	var result T
	if col.Count() == 0 {
		return result, false
	}
	return col.array[len(col.array)-1], true
}

func (col *Collection[T]) Prepend(items ...T) contracts.Collection[T] {
	return New(append(items, col.array...))
}

func (col *Collection[T]) Push(items ...T) contracts.Collection[T] {
	return New(append(col.array, items...))
}

func (col *Collection[T]) Pull(defaultValue ...T) (T, bool) {
	result, exists := col.Last()
	if exists {
		col.array = col.array[:col.Count()-1]
		return result, true
	} else if len(defaultValue) > 0 {
		return defaultValue[0], true
	}

	return result, false
}

func (col *Collection[T]) Shift(defaultValue ...T) (T, bool) {
	result, exists := col.First()
	if exists {
		col.array = col.array[1:]
		return result, true
	} else if len(defaultValue) > 0 {
		return defaultValue[0], true
	}

	return result, false
}

func (col *Collection[T]) Offset(index int, item T) contracts.Collection[T] {
	if col.Count() > index {
		col.array[index] = item
		return col
	}
	return col.Push(item)
}

func (col *Collection[T]) Put(index int, item T) contracts.Collection[T] {
	if col.Count() > index {
		return New(col.array).Offset(index, item)
	}
	return col.Push(item)
}

func (col *Collection[T]) Merge(collections ...contracts.Collection[T]) contracts.Collection[T] {
	newCollection := New(col.array)

	for _, collection := range collections {
		newCollection = newCollection.Push(collection.ToArray()...)
	}

	return newCollection
}

func (col *Collection[T]) Reverse() contracts.Collection[T] {
	newCollection := &Collection[T]{array: col.array}
	for from, to := 0, len(newCollection.array)-1; from < to; from, to = from+1, to-1 {
		newCollection.array[from], newCollection.array[to] = newCollection.array[to], newCollection.array[from]
	}
	return newCollection
}

func (col *Collection[T]) Chunk(size int, handler func(collection contracts.Collection[T], page int) error) (err error) {
	total := col.Count()
	page := 1
	for err == nil && (page-1)*size <= total {
		offset := (page - 1) * size
		endIndex := size + offset
		if endIndex > total {
			endIndex = total
		}
		newCollection := &Collection[T]{array: col.array[offset:endIndex]}

		err = handler(newCollection, page)
		page++
	}

	return
}

func (col *Collection[T]) Random(size ...uint) contracts.Collection[T] {
	num := 1
	if len(size) > 0 {
		num = int(size[0])
	}
	newCollection := &Collection[T]{}
	if col.Count() >= num {
		for _, index := range utils.RandIntArray(0, col.Count()-1, num) {
			newCollection.array = append(newCollection.array, col.array[index])
		}
	}
	return newCollection
}
