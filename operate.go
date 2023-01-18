package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

func (col *Collection) Pluck(key string) contracts.Fields {
	fields := contracts.Fields{}

	for index, data := range col.mapData {
		var name, ok = data[key].(string)
		if _, exists := fields[name]; ok && !exists {
			fields[name] = col.array[index]
		}
	}

	return fields
}

func (col *Collection) Only(keys ...string) contracts.Collection {
	arrayFields := make([]contracts.Fields, 0)
	rawResults := make([]interface{}, 0)

	for index, data := range col.mapData {
		fields := contracts.Fields{}
		for key, value := range data {
			if utils.IsIn(key, keys) {
				fields[key] = value
			}
		}
		arrayFields = append(arrayFields, fields)
		rawResults = append(rawResults, col.array[index])
	}

	return &Collection{mapData: arrayFields, array: rawResults}
}

func (col *Collection) First(keys ...string) interface{} {
	if col.Count() == 0 {
		return nil
	}
	if len(keys) == 0 {
		return col.array[0]
	}
	return col.mapData[0][keys[0]]
}

func (col *Collection) Last(keys ...string) interface{} {
	if col.Count() == 0 {
		return nil
	}
	if len(keys) == 0 {
		return col.array[len(col.array)-1]
	}
	return col.mapData[len(col.array)-1][keys[0]]
}

func (col *Collection) Prepend(items ...interface{}) contracts.Collection {
	newCollection := &Collection{}
	newCollection.array = append(items, col.array...)
	if len(col.mapData) > 0 {
		newMaps := make([]contracts.Fields, 0)
		for _, item := range items {
			fields, _ := utils.ConvertToFields(item)
			newMaps = append(newMaps, fields)
		}
		newCollection.mapData = append(newMaps, col.mapData...)
	}
	return newCollection
}

func (col *Collection) Push(items ...interface{}) contracts.Collection {
	newCollection := &Collection{}
	newCollection.array = append(col.array, items...)
	if len(col.mapData) > 0 {
		newMaps := make([]contracts.Fields, 0)
		for _, item := range items {
			fields, _ := utils.ConvertToFields(item)
			newMaps = append(newMaps, fields)
		}
		newCollection.mapData = append(col.mapData, newMaps...)
	}
	return newCollection
}

func (col *Collection) Pull(defaultValue ...interface{}) interface{} {
	if result := col.Last(); result != nil {
		col.array = col.array[:col.Count()-1]
		if len(col.mapData) > 0 {
			col.mapData = col.mapData[:col.Count()-1]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (col *Collection) Shift(defaultValue ...interface{}) interface{} {
	if result := col.First(); result != nil {
		col.array = col.array[1:]
		if len(col.mapData) > 0 {
			col.mapData = col.mapData[1:]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (col *Collection) Offset(index int, item interface{}) contracts.Collection {
	if col.Count() > index {
		col.array[index] = item
		if len(col.mapData) > 0 {
			fields, _ := utils.ConvertToFields(item)
			col.mapData[index] = fields
		}
		return col
	}
	return col.Push(item)
}

func (col *Collection) Put(index int, item interface{}) contracts.Collection {
	if col.Count() > index {
		return (&Collection{array: append(col.array), mapData: append(col.mapData)}).Offset(index, item)
	}
	return col.Push(item)
}

func (col *Collection) Merge(collections ...contracts.Collection) contracts.Collection {
	newCollection := &Collection{array: append(col.array), mapData: append(col.mapData)}

	for _, collection := range collections {
		newCollection.mapData = append(newCollection.mapData, collection.ToArrayFields()...)
		newCollection.array = append(newCollection.array, collection.ToInterfaceArray()...)
	}

	return newCollection
}

func (col *Collection) Reverse() contracts.Collection {
	newCollection := &Collection{array: append(col.array), mapData: append(col.mapData)}
	for from, to := 0, len(newCollection.array)-1; from < to; from, to = from+1, to-1 {
		newCollection.array[from], newCollection.array[to] = newCollection.array[to], newCollection.array[from]
		if len(col.mapData) > 0 {
			newCollection.mapData[from], newCollection.mapData[to] = newCollection.mapData[to], newCollection.mapData[from]
		}
	}
	return newCollection
}

func (col *Collection) Chunk(size int, handler func(collection contracts.Collection, page int) error) (err error) {
	total := col.Count()
	page := 1
	for err == nil && (page-1)*size <= total {
		offset := (page - 1) * size
		endIndex := size + offset
		if endIndex > total {
			endIndex = total
		}
		newCollection := &Collection{array: col.array[offset:endIndex]}
		if len(col.mapData) > 0 {
			newCollection.mapData = col.mapData[offset:endIndex]
		}

		err = handler(newCollection, page)
		page++
	}

	return
}

func (col *Collection) Random(size ...uint) contracts.Collection {
	num := 1
	if len(size) > 0 {
		num = int(size[0])
	}
	newCollection := &Collection{}
	if col.Count() >= num {
		for _, index := range utils.RandIntArray(0, col.Count()-1, num) {
			newCollection.array = append(newCollection.array, col.array[index])
			if len(col.mapData) > 0 {
				newCollection.mapData = append(newCollection.mapData, col.mapData[index])
			}
		}
	}
	return newCollection
}
