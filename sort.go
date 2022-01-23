package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/exceptions"
	"reflect"
	"sort"
	"strconv"
)

func (this *Collection) Len() int {
	return len(this.array)
}

func (this *Collection) Swap(i, j int) {
	preIndex := strconv.Itoa(i)
	nextIndex := strconv.Itoa(j)
	this.array[preIndex], this.array[nextIndex] = this.array[nextIndex], this.array[preIndex]
	if len(this.mapData) > 0 {
		this.mapData[preIndex], this.mapData[nextIndex] = this.mapData[nextIndex], this.mapData[preIndex]
	}
}

func (this *Collection) Less(i, j int) bool {
	if this.sorter != nil {
		return this.sorter(i, j)
	}
	return i > j
}

func (this *Collection) SetSorter(sorter func(i, j int) bool) *Collection {
	this.sorter = sorter
	return this
}

// Sort sorter 必须是接收两个参数，并且返回一个 bool 值的函数
func (this *Collection) Sort(sorter interface{}) *Collection {
	sorterType := reflect.TypeOf(sorter)

	if sorterType.Kind() != reflect.Func || sorterType.NumIn() != 2 || sorterType.NumOut() != 1 || sorterType.Out(0).Kind() != reflect.Bool {
		exceptions.Throw(SortException{
			exceptions.New(
				"参数类型异常：sorter 必须是接收两个参数，并且返回一个 bool 值的函数",
				contracts.Fields{
					"sorter": sorter,
				})})
	}
	sorterValue := reflect.ValueOf(sorter)

	newCollection := (&Collection{mapData: this.mapData, array: this.array}).SetSorter(func(i, j int) bool {
		return sorterValue.Call([]reflect.Value{
			this.argumentConvertor(sorterType.In(0), this.array[strconv.Itoa(i)]),
			this.argumentConvertor(sorterType.In(1), this.array[strconv.Itoa(j)]),
		})[0].Bool()
	})

	sort.Sort(newCollection)

	return newCollection
}
