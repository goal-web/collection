package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/utils"
	"reflect"
)

type Collection[T any] struct {
	rawData []T
	values  map[int]reflect.Value
	sorter  func(i, j int) bool
}

func New[T](data any) (contracts.Collection[T], error) {
	dataValue := reflect.ValueOf(data)

	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	switch dataValue.Kind() {
	case reflect.Array, reflect.Slice:
		collect := &Collection[T]{rawData: nil}
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

func MustNew[T](items []T) contracts.Collection[T] {
	c, _ := New[T](items)
	return c
}

func Array[T any](data []T) contracts.Collection[T] {
	collection := &Collection[T]{rawData: data}

	return collection
}

func (this *Collection[T]) Index(index int) T {
	if this.Count() > index {
		return this.rawData[index]
	}
	return nil
}

func (this *Collection[T]) IsEmpty() bool {
	return this.Count() == 0
}

func (this *Collection[T]) Map(fn func(T, int) T) contracts.Collection[T] {
	var results []T
	for index, item := range this.rawData {
		results = append(results, fn(item, index))
	}

	return Array(results)
}

func (this *Collection[T]) argumentConvertor(argType reflect.Type, arg any) reflect.Value {
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
