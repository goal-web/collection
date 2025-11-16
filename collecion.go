package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"reflect"
)

type Collection[T any] struct {
	array    []T
	mapArray []contracts.Fields
	sorter   func(i, j int) bool
}

func New[T any](data []T) contracts.Collection[T] {
	return &Collection[T]{array: data}
}

func (col *Collection[T]) Index(index int) *T {
	if col.Count() > index {
		return &col.array[index]
	}
	return nil
}

func (col *Collection[T]) IsEmpty() bool {
	return len(col.array) == 0
}

func (col *Collection[T]) Map(handler any) contracts.Collection[T] {
	switch handle := handler.(type) {
	case func(fields contracts.Fields):
		for _, fields := range col.ToArrayFields() {
			handle(fields)
		}
	case func(fields contracts.Fields, index int): // 聚合函数
		for index, fields := range col.ToArrayFields() {
			handle(fields, index)
		}
	case func(fields contracts.Fields) bool: // 用于 where
		results := make([]T, 0)
		for i, data := range col.ToArrayFields() {
			if handle(data) {
				results = append(results, col.array[i])
			}
		}
		return New(results)
	case func(T):
		for _, data := range col.array {
			handle(data)
		}
	case func(i int, fields T):
		for i, data := range col.array {
			handle(i, data)
		}
	case func(i int, fields *T):
		for i, data := range col.array {
			handle(i, &data)
		}
	case func(T) bool:
		results := make([]T, 0)
		for _, data := range col.array {
			if handle(data) {
				results = append(results, data)
			}
		}
		return New(results)
	case func(int, T) bool:
		results := make([]T, 0)
		for i, data := range col.array {
			if handle(i, data) {
				results = append(results, data)
			}
		}
		return New(results)
	case func(T) T:
		for i, data := range col.array {
			col.array[i] = handle(data)
		}
		col.mapArray = nil
	case func(any):
		for _, data := range col.array {
			handle(data)
		}
	case func(any) any:
		for index, data := range col.array {
			result, ok := handle(data).(T)
			if ok {
				col.array[index] = result
				col.mapArray = nil
			}
		}
	default:
		handlerType := reflect.TypeOf(handler)
		handlerValue := reflect.ValueOf(handler)
		argsGetter := func(index int, arg any) []reflect.Value {
			return []reflect.Value{}
		}
		numOut := handlerType.NumOut()

		switch handlerType.NumIn() {
		case 1:
			argsGetter = func(_ int, arg any) []reflect.Value {
				return []reflect.Value{col.argumentConvertor(handlerType.In(0), arg)}
			}
		case 2:
			argsGetter = func(index int, arg any) []reflect.Value {
				return []reflect.Value{
					col.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
				}
			}
		case 3:
			array := reflect.ValueOf(col.array)
			argsGetter = func(index int, arg any) []reflect.Value {
				return []reflect.Value{
					col.argumentConvertor(handlerType.In(0), arg),
					reflect.ValueOf(index),
					array,
				}
			}
		}

		if numOut > 0 {
			newCollection := &Collection[T]{}
			for index, data := range col.array {
				result := handlerValue.Call(argsGetter(index, data))[0].Interface()
				newCollection.array = append(newCollection.array, result.(T))
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

func (col *Collection[T]) argumentConvertor(argType reflect.Type, arg any) reflect.Value {
	switch argType.Kind() {
	case reflect.String:
		return reflect.ValueOf(utils.ToString(arg, ""))
	case reflect.Int:
		return reflect.ValueOf(utils.ToInt(arg, 0))
	case reflect.Int64:
		return reflect.ValueOf(utils.ToInt64(arg, 0))
	case reflect.Float64:
		return reflect.ValueOf(utils.ToFloat64(arg, 0))
	case reflect.Float32:
		return reflect.ValueOf(utils.ToFloat(arg, 0))
	case reflect.Bool:
		return reflect.ValueOf(utils.ToBool(arg, false))
	}
	if reflect.TypeOf(arg).ConvertibleTo(argType) {
		return reflect.ValueOf(arg).Convert(argType)
	}
	return reflect.ValueOf(arg)

}
