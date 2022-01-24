package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
	"reflect"
	"strconv"
)

type array []interface{}

type Collection struct {
	array
	mapData []contracts.Fields
	sorter  func(i, j int) bool
}

func (this *Collection) Pluck(key string) contracts.Fields {
	fields := contracts.Fields{}

	for _, data := range this.mapData {
		if _, exists := fields[key]; !exists {
			fields[key] = data
		}
	}

	return fields
}

func (this *Collection) Only(keys ...string) contracts.Collection {
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
		rawResults = append(rawResults, this.array[index])
	}

	return &Collection{mapData: arrayFields, array: rawResults}
}

func (this *Collection) First(keys ...string) interface{} {
	if this.Count() == 0 {
		return nil
	}
	if len(keys) == 0 {
		return this.array[0]
	}
	return this.mapData[0][keys[0]]
}

func (this *Collection) Last(keys ...string) interface{} {
	if this.Count() == 0 {
		return nil
	}
	if len(keys) == 0 {
		return this.array[len(this.array)-1]
	}
	return this.mapData[len(this.array)-1][keys[0]]
}

func (this *Collection) Prepend(items ...interface{}) contracts.Collection {
	newCollection := &Collection{}
	newCollection.array = append(items, this.array...)
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

func (this *Collection) Push(items ...interface{}) contracts.Collection {
	newCollection := &Collection{}
	newCollection.array = append(this.array, items...)
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

func (this *Collection) Pull(defaultValue ...interface{}) interface{} {
	if result := this.Last(); result != nil {
		this.array = this.array[:this.Count()-1]
		if len(this.mapData) > 0 {
			this.mapData = this.mapData[:this.Count()-1]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (this *Collection) Shift(defaultValue ...interface{}) interface{} {
	if result := this.First(); result != nil {
		this.array = this.array[1:]
		if len(this.mapData) > 0 {
			this.mapData = this.mapData[1:]
		}
		return result
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (this *Collection) Offset(index int, item interface{}) contracts.Collection {
	if this.Count() > index {
		this.array[index] = item
		if len(this.mapData) > 0 {
			fields, _ := utils.ConvertToFields(item)
			this.mapData[index] = fields
		}
		return this
	}
	return this.Push(item)
}

func (this *Collection) Put(index int, item interface{}) contracts.Collection {
	if this.Count() > index {
		return (&Collection{array: append(this.array), mapData: append(this.mapData)}).Offset(index, item)
	}
	return this.Push(item)
}

func (this *Collection) Merge(collections ...contracts.Collection) contracts.Collection {
	newCollection := &Collection{array: append(this.array), mapData: append(this.mapData)}

	for _, collection := range collections {
		newCollection.mapData = append(newCollection.mapData, collection.ToArrayFields()...)
		newCollection.array = append(newCollection.array, collection.ToInterfaceArray()...)
	}

	return newCollection
}

func (this *Collection) Reverse() contracts.Collection {
	newCollection := &Collection{array: append(this.array), mapData: append(this.mapData)}
	for from, to := 0, len(newCollection.array)-1; from < to; from, to = from+1, to-1 {
		newCollection.array[from], newCollection.array[to] = newCollection.array[to], newCollection.array[from]
		if len(this.mapData) > 0 {
			newCollection.mapData[from], newCollection.mapData[to] = newCollection.mapData[to], newCollection.mapData[from]
		}
	}
	return newCollection
}

func (this *Collection) Chunk(size int, handler func(collection contracts.Collection, page int) error) (err error) {
	total := this.Count()
	page := 1
	for err == nil && (page-1)*size <= total {
		newCollection := &Collection{
			array: this.array[(page-1)*size : size],
		}

		if len(this.mapData) > 0 {
			newCollection.mapData = this.mapData[(page-1)*size : size]
		}

		err = handler(newCollection, page)
		page++
	}

	return
}

func (this *Collection) Random(size ...uint) contracts.Collection {
	num := 1
	if len(size) > 0 {
		num = int(size[0])
	}
	newCollection := &Collection{}
	if this.Count() >= num {
		for _, index := range utils.RandIntArray(0, this.Count()-1, num) {
			newCollection.array = append(newCollection.array, this.array[index])
			if len(this.mapData) > 0 {
				newCollection.mapData = append(newCollection.mapData, this.mapData[index])
			}
		}
	}
	return newCollection
}

func (this *Collection) ToFields() contracts.Fields {
	fields := contracts.Fields{}
	for index, data := range this.mapData {
		fields[strconv.Itoa(index)] = data
	}
	return fields
}

func (this *Collection) ToArrayFields() []contracts.Fields {
	return this.mapData
}

func New(data interface{}) (*Collection, error) {
	dataValue := reflect.ValueOf(data)

	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	switch dataValue.Kind() {
	case reflect.Array, reflect.Slice:
		collect := &Collection{mapData: nil, array: make([]interface{}, 0)}
		isMapData := true

		utils.EachSlice(dataValue, func(_ int, value reflect.Value) {
			switch value.Kind() {
			case reflect.Map, reflect.Struct:
				collect.array = append(collect.array, value.Interface())
			default:
				isMapData = false
				collect.array = append(collect.array, value.Interface())
			}
		})

		if isMapData {
			collect.mapData = make([]contracts.Fields, 0)
			for _, item := range collect.array {
				fields, _ := utils.ConvertToFields(item)
				collect.mapData = append(collect.mapData, fields)
			}
		}

		return collect, nil
	}

	return nil, Exception{
		Exception: exceptions.New("不支持的类型 "+utils.GetTypeKey(reflect.TypeOf(data)), contracts.Fields{
			"data": data,
		}),
	}
}

func MustNew(data interface{}) contracts.Collection {
	c, err := New(data)
	exceptions.Throw(err)
	return c
}

func (this *Collection) Index(index int) interface{} {
	if this.Count() > index {
		return this.array[index]
	}
	return nil
}

func (this *Collection) IsEmpty() bool {
	return len(this.array) == 0
}
func (this *Collection) Each(handler interface{}) contracts.Collection {
	return this.Map(handler)
}

func (this *Collection) Map(handler interface{}) contracts.Collection {

	switch handle := handler.(type) {
	case func(fields contracts.Fields):
		for _, fields := range this.mapData {
			handle(fields)
		}
	case func(fields contracts.Fields, index int): // 聚合函数
		for index, fields := range this.mapData {
			handle(fields, index)
		}
	case func(fields contracts.Fields) contracts.Fields:
		for index, fields := range this.mapData {
			this.mapData[index] = handle(fields)
		}
	case func(fields contracts.Fields) bool: // 用于 where
		results := make([]interface{}, 0)
		for _, fields := range this.mapData {
			results = append(results, handle(fields))
		}
		return &Collection{mapData: nil, array: results}
	case func(fields interface{}):
		for _, data := range this.array {
			handle(data)
		}
	case func(fields interface{}) interface{}:
		for index, data := range this.array {
			this.array[index] = handle(data)
		}
	default:
		handlerType := reflect.TypeOf(handler)
		handlerValue := reflect.ValueOf(handler)
		argsGetter := func(index int, arg interface{}) []reflect.Value {
			return []reflect.Value{}
		}
		numOut := handlerType.NumOut()

		switch handlerType.NumIn() {
		case 1:
			argsGetter = func(_ int, arg interface{}) []reflect.Value {
				return []reflect.Value{this.argumentConvertor(handlerType.In(0), arg)}
			}
		case 2:
			argsGetter = func(index int, arg interface{}) []reflect.Value {
				return []reflect.Value{
					this.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
				}
			}
		case 3:
			array := reflect.ValueOf(this.array)
			argsGetter = func(index int, arg interface{}) []reflect.Value {
				return []reflect.Value{
					this.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
					array,
				}
			}
		}

		if numOut > 0 {
			mapLen := len(this.mapData)
			newCollection := &Collection{mapData: make([]contracts.Fields, mapLen), array: make([]interface{}, 0)}
			for index, data := range this.array {
				result := handlerValue.Call(argsGetter(index, data))[0].Interface()
				if mapLen > 0 {
					newCollection.mapData[index], _ = utils.ConvertToFields(result)
				}
				newCollection.array = append(newCollection.array, result)
			}
			return newCollection
		} else {
			for index, data := range this.array {
				handlerValue.Call(argsGetter(index, data))
			}
		}

	}

	return this
}

func (this *Collection) argumentConvertor(argType reflect.Type, arg interface{}) reflect.Value {
	switch argType.Kind() {
	case reflect.String:
		return reflect.ValueOf(utils.ConvertToString(arg, ""))
	case reflect.Int:
		return reflect.ValueOf(utils.ConvertToInt(arg, 0))
	case reflect.Int64:
		return reflect.ValueOf(utils.ConvertToInt64(arg, 0))
	case reflect.Float64:
		return reflect.ValueOf(utils.ConvertToFloat64(arg, 0))
	case reflect.Float32:
		return reflect.ValueOf(utils.ConvertToFloat(arg, 0))
	case reflect.Bool:
		return reflect.ValueOf(utils.ConvertToBool(arg, false))
	}
	if reflect.TypeOf(arg).ConvertibleTo(argType) {
		return reflect.ValueOf(arg).Convert(argType)
	}
	return reflect.ValueOf(arg)

}
