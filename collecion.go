package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
	"reflect"
	"strconv"
)

type Collection struct {
	array   contracts.Fields
	mapData map[string]contracts.Fields
	sorter  func(i, j int) bool
}

func New(data interface{}) (*Collection, error) {
	dataValue := reflect.ValueOf(data)

	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	switch dataValue.Kind() {
	case reflect.Array, reflect.Slice:
		collect := &Collection{mapData: nil, array: make(contracts.Fields)}
		isMapData := true

		utils.EachSlice(dataValue, func(i int, value reflect.Value) {
			switch value.Kind() {
			case reflect.Map, reflect.Struct:
				collect.array[strconv.Itoa(i)] = value.Interface()
			default:
				isMapData = false
				collect.array[strconv.Itoa(i)] = value.Interface()
			}
		})

		if isMapData {
			collect.mapData = make(map[string]contracts.Fields, 0)
			for index, item := range collect.array {
				fields, _ := utils.ConvertToFields(item)
				collect.mapData[index] = fields
			}
		}

		return collect, nil
	case reflect.Map:

	}

	return nil, Exception{
		Exception: exceptions.New("不支持的类型 "+utils.GetTypeKey(reflect.TypeOf(data)), contracts.Fields{
			"data": data,
		}),
	}
}

func MustNew(data interface{}) *Collection {
	c, err := New(data)
	exceptions.Throw(err)
	return c
}

// Filter 过滤不需要的数据 filter 返回 true 时保留
func (this *Collection) Filter(filter interface{}) *Collection {
	results := make(contracts.Fields)
	newFields := make(map[string]contracts.Fields)
	for index, data := range this.Map(filter).ToFields() {
		if utils.ConvertToBool(data, false) {
			if fields := this.mapData[index]; fields != nil {
				newFields[index] = this.mapData[index]
			}
			results[index] = this.array[index]
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

// Skip 过滤不需要的数据 filter 返回 true 时过滤
func (this *Collection) Skip(filter interface{}) *Collection {
	results := make(contracts.Fields)
	newFields := make(map[string]contracts.Fields)
	for index, data := range this.Map(filter).ToFields() {
		if !utils.ConvertToBool(data, false) {
			if fields := this.mapData[index]; fields != nil {
				newFields[index] = this.mapData[index]
			}
			results[index] = this.array[index]
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

// Where 根据条件过滤数据，支持 =,>,>=,<,<=,in,not in 等条件判断
func (this *Collection) Where(field string, args ...interface{}) *Collection {
	results := make(contracts.Fields)
	var (
		arg      interface{}
		operator = "="
	)
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		operator = args[0].(string)
		arg = args[1]
	}
	newFields := make(map[string]contracts.Fields)
	for index, data := range this.Map(func(fields contracts.Fields) bool {
		return utils.Compare(fields[field], operator, arg)
	}).ToFields() {
		if utils.ConvertToBool(data, false) {
			if fields := this.mapData[index]; fields != nil {
				newFields[index] = this.mapData[index]
			}
			results[index] = this.array[index]
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

func (this *Collection) WhereLt(field string, arg interface{}) *Collection {
	return this.Where(field, "lt", arg)
}
func (this *Collection) WhereLte(field string, arg interface{}) *Collection {
	return this.Where(field, "lte", arg)
}
func (this *Collection) WhereGt(field string, arg interface{}) *Collection {
	return this.Where(field, "gt", arg)
}
func (this *Collection) WhereGte(field string, arg interface{}) *Collection {
	return this.Where(field, "gte", arg)
}
func (this *Collection) WhereIn(field string, arg interface{}) *Collection {
	return this.Where(field, "in", arg)
}
func (this *Collection) WhereNotIn(field string, arg interface{}) *Collection {
	return this.Where(field, "not in", arg)
}

func (this *Collection) Index(index interface{}) interface{} {
	key := utils.ConvertToString(index, "")
	return this.array[key]
}

func (this *Collection) Map(handler interface{}) *Collection {

	switch handle := handler.(type) {
	case func(fields contracts.Fields):
		for _, fields := range this.mapData {
			handle(fields)
		}
	case func(fields contracts.Fields) contracts.Fields:
		for index, fields := range this.mapData {
			this.mapData[index] = handle(fields)
		}
	case func(fields contracts.Fields) bool: // 用于 where
		fields := make(contracts.Fields)
		for key, data := range this.mapData {
			fields[key] = handle(data)
		}
		return &Collection{mapData: nil, array: fields}
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
		argsGetter := func(index string, arg interface{}) []reflect.Value {
			return []reflect.Value{}
		}
		numOut := handlerType.NumOut()

		switch handlerType.NumIn() {
		case 1:
			argsGetter = func(index string, arg interface{}) []reflect.Value {
				return []reflect.Value{this.argumentConvertor(handlerType.In(0), arg)}
			}
		case 2:
			argsGetter = func(index string, arg interface{}) []reflect.Value {
				return []reflect.Value{
					this.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
				}
			}
		case 3:
			array := reflect.ValueOf(this.array)
			argsGetter = func(index string, arg interface{}) []reflect.Value {
				return []reflect.Value{
					this.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
					array,
				}
			}
		}

		if numOut > 0 {
			mapLen := len(this.mapData)
			newCollection := &Collection{mapData: make(map[string]contracts.Fields), array: make(contracts.Fields)}
			for index, data := range this.array {
				result := handlerValue.Call(argsGetter(index, data))[0].Interface()
				if mapLen > 0 {
					newCollection.mapData[index], _ = utils.ConvertToFields(result)
				}
				newCollection.array[index] = result
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
