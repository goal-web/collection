package collection

import (
	"errors"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"reflect"
	"sort"
)

func (col *Collection) Len() int {
	return len(col.array)
}

func (col *Collection) Swap(i, j int) {
	col.array[i], col.array[j] = col.array[j], col.array[i]
	if len(col.mapData) > 0 {
		col.mapData[i], col.mapData[j] = col.mapData[j], col.mapData[i]
	}
}

func (col *Collection) Less(i, j int) bool {
	if col.sorter != nil {
		return col.sorter(i, j)
	}
	return i > j
}

func (col *Collection) SetSorter(sorter func(i, j int) bool) contracts.Collection {
	col.sorter = sorter
	return col
}

// Sort sorter 必须是接收两个参数，并且返回一个 bool 值的函数
func (col *Collection) Sort(sorter interface{}) contracts.Collection {
	sorterType := reflect.TypeOf(sorter)

	if sorterType.Kind() != reflect.Func || sorterType.NumIn() != 2 || sorterType.NumOut() != 1 || sorterType.Out(0).Kind() != reflect.Bool {
		exceptions.Throw(SortException{
			Err: errors.New("参数类型异常：sorter 必须是接收两个参数，并且返回一个 bool 值的函数"),
		})
	}
	sorterValue := reflect.ValueOf(sorter)

	newCollection := (&Collection{
		mapData: col.mapData,
		array:   col.array,
	}).SetSorter(func(i, j int) bool {
		return sorterValue.Call([]reflect.Value{
			col.argumentConvertor(sorterType.In(0), col.array[i]),
			col.argumentConvertor(sorterType.In(1), col.array[j]),
		})[0].Bool()
	})

	sort.Stable(newCollection)

	return newCollection
}
