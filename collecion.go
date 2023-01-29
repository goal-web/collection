package collection

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
	"reflect"
)

type array []interface{}

type Collection struct {
	array
	mapData []contracts.Fields
	sorter  func(i, j int) bool
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
		Err: errors.New("不支持的类型 " + utils.GetTypeKey(reflect.TypeOf(data))),
	}
}

func MustNew(data interface{}) contracts.Collection {
	c, err := New(data)
	exceptions.Throw(err)
	return c
}

func FromFieldsSlice(data []contracts.Fields) contracts.Collection {
	collection := &Collection{mapData: data}
	for _, fields := range data {
		collection.array = append(collection.array, fields)
	}
	return collection
}

func (col *Collection) Index(index int) interface{} {
	if col.Count() > index {
		return col.array[index]
	}
	return nil
}

func (col *Collection) IsEmpty() bool {
	return len(col.array) == 0
}
func (col *Collection) Each(handler interface{}) contracts.Collection {
	return col.Map(handler)
}

func (col *Collection) Map(handler interface{}) contracts.Collection {

	switch handle := handler.(type) {
	case func(fields contracts.Fields):
		for _, fields := range col.mapData {
			handle(fields)
		}
	case func(fields contracts.Fields, index int): // 聚合函数
		for index, fields := range col.mapData {
			handle(fields, index)
		}
	case func(fields contracts.Fields) contracts.Fields:
		for index, fields := range col.mapData {
			col.mapData[index] = handle(fields)
		}
	case func(fields contracts.Fields) bool: // 用于 where
		results := make([]interface{}, 0)
		for _, fields := range col.mapData {
			results = append(results, handle(fields))
		}
		return &Collection{mapData: nil, array: results}
	case func(fields interface{}):
		for _, data := range col.array {
			handle(data)
		}
	case func(fields interface{}) interface{}:
		for index, data := range col.array {
			col.array[index] = handle(data)
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
				return []reflect.Value{col.argumentConvertor(handlerType.In(0), arg)}
			}
		case 2:
			argsGetter = func(index int, arg interface{}) []reflect.Value {
				return []reflect.Value{
					col.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
				}
			}
		case 3:
			array := reflect.ValueOf(col.array)
			argsGetter = func(index int, arg interface{}) []reflect.Value {
				return []reflect.Value{
					col.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
					array,
				}
			}
		}

		if numOut > 0 {
			mapLen := len(col.mapData)
			newCollection := &Collection{mapData: make([]contracts.Fields, mapLen), array: make([]interface{}, 0)}
			for index, data := range col.array {
				result := handlerValue.Call(argsGetter(index, data))[0].Interface()
				if mapLen > 0 {
					newCollection.mapData[index], _ = utils.ConvertToFields(result)
				}
				newCollection.array = append(newCollection.array, result)
			}
			return newCollection
		} else {
			for index, data := range col.array {
				handlerValue.Call(argsGetter(index, data))
			}
		}

	}

	return col
}

func (col *Collection) argumentConvertor(argType reflect.Type, arg interface{}) reflect.Value {
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
