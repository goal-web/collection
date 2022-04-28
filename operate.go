package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func (this *Collection[T]) Pluck(key string) contracts.Fields {
	fields := contracts.Fields{}

	for index, data := range this.rawData {
		var name, ok = data[key].(string)
		if _, exists := fields[name]; ok && !exists {
			fields[name] = this.rawData[index]
		}
	}

	return fields
}

func (this *Collection[T]) Only(keys ...string) contracts.Collection[T] {
	arrayFields := make([]contracts.Fields, 0)
	rawResults := make([]interface{}, 0)

	for index, data := range this.mapData {
		fields := contracts.Fields{}
		for key, value := range data {
			if utils.IsIn(key, keys) {
				fields[key] = value
			}
		}
		arrayFields = append(arrayFields, fields)
		rawResults = append(rawResults, this.rawData[index])
	}

	return &Collection{mapData: arrayFields, array: rawResults}
}

func (this *Collection[T]) First() T {
	if this.IsEmpty() {
		return nil
	}
	return this.rawData[0]
}

func (this *Collection[T]) Last() T {
	if this.IsEmpty() {
		return nil
	}
	return this.rawData[len(this.rawData)-1]
}

func (this *Collection[T]) Prepend(items ...interface{}) contracts.Collection[T] {
	newCollection := &Collection{}
	newCollection.array = append(items, this.rawData...)
	if len(this.mapData) > 0 {
		newMaps := make([]contracts.Fields, 0)
		for _, item := range items {
			fields, _ := utils.ConvertToFields(item)
			newMaps = append(newMaps, fields)
		}
		newCollection.mapData = append(newMaps, this.mapData...)
	}
	return newCollection
}

func (this *Collection[T]) Push(items ...interface{}) contracts.Collection[T] {
	newCollection := &Collection{}
	newCollection.array = append(this.rawData, items...)
	if len(this.mapData) > 0 {
		newMaps := make([]contracts.Fields, 0)
		for _, item := range items {
			fields, _ := utils.ConvertToFields(item)
			newMaps = append(newMaps, fields)
		}
		newCollection.mapData = append(this.mapData, newMaps...)
	}
	return newCollection
}

func (this *Collection[T]) Pull(defaultValue ...interface{}) interface{} {
	if result := this.Last(); result != nil {
		this.rawData = this.rawData[:this.Count()-1]
		if len(this.mapData) > 0 {
			this.mapData = this.mapData[:this.Count()-1]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (this *Collection[T]) Shift(defaultValue ...interface{}) interface{} {
	if result := this.First(); result != nil {
		this.rawData = this.rawData[1:]
		if len(this.mapData) > 0 {
			this.mapData = this.mapData[1:]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (this *Collection[T]) Offset(index int, item T) contracts.Collection[T] {
	if this.Count() > index {
		this.rawData[index] = item
		if len(this.mapData) > 0 {
			fields, _ := utils.ConvertToFields(item)
			this.mapData[index] = fields
		}
		return this
	}
	return this.Push(item)
}

func (this *Collection[T]) Put(index int, item interface{}) contracts.Collection[T] {
	if this.Count() > index {
		return (&Collection{array: append(this.rawData), mapData: append(this.mapData)}).Offset(index, item)
	}
	return this.Push(item)
}

func (this *Collection[T]) Merge(collections ...contracts.Collection) contracts.Collection[T] {
	newCollection := &Collection{array: append(this.rawData), mapData: append(this.mapData)}

	for _, collection := range collections {
		newCollection.mapData = append(newCollection.mapData, collection.ToArrayFields()...)
		newCollection.array = append(newCollection.array, collection.ToInterfaceArray()...)
	}

	return newCollection
}

func (this *Collection[T]) Reverse() contracts.Collection[T] {
	newCollection := &Collection{array: append(this.rawData), mapData: append(this.mapData)}
	for from, to := 0, len(newCollection.array)-1; from < to; from, to = from+1, to-1 {
		newCollection.array[from], newCollection.array[to] = newCollection.array[to], newCollection.array[from]
		if len(this.mapData) > 0 {
			newCollection.mapData[from], newCollection.mapData[to] = newCollection.mapData[to], newCollection.mapData[from]
		}
	}
	return newCollection
}

func (this *Collection[T]) Chunk(size int, handler func(collection contracts.Collection, page int) error) (err error) {
	total := this.Count()
	page := 1
	for err == nil && (page-1)*size <= total {
		offset := (page - 1) * size
		endIndex := size + offset
		if endIndex > total {
			endIndex = total
		}
		newCollection := &Collection{rawData: this.rawData[offset:endIndex]}
		if len(this.mapData) > 0 {
			newCollection.mapData = this.mapData[offset:endIndex]
		}

		err = handler(newCollection, page)
		page++
	}

	return
}

func (this *Collection[T]) Random(size ...uint) contracts.Collection[T] {
	num := 1
	if len(size) > 0 {
		num = int(size[0])
	}
	newCollection := &Collection{}
	if this.Count() >= num {
		for _, index := range utils.RandIntArray(0, this.Count()-1, num) {
			newCollection.array = append(newCollection.rawData, this.rawData[index])
			if len(this.mapData) > 0 {
				newCollection.mapData = append(newCollection.mapData, this.mapData[index])
			}
		}
	}
	return newCollection
}
