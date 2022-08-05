package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"strconv"
)

func (col *Collection[T]) ToIntArray() (results []int) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToInt(data, 0))
	}
	return
}

func (col *Collection[T]) ToInt64Array() (results []int64) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToInt64(data, 0))
	}
	return
}

func (col *Collection[T]) ToFloatArray() (results []float32) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToFloat(data, 0))
	}
	return
}

func (col *Collection[T]) ToFloat64Array() (results []float64) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToFloat64(data, 0))
	}
	return
}

func (col *Collection[T]) ToBoolArray() (results []bool) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToBool(data, false))
	}
	return
}

func (col *Collection[T]) ToStringArray() (results []string) {
	for _, data := range col.rawData {
		results = append(results, utils.ConvertToString(data, ""))
	}
	return
}

func (col *Collection[T]) ToFields() contracts.Fields {
	fields := contracts.Fields{}
	for index, data := range col.rawData {
		fields[strconv.Itoa(index)] = data
	}
	return fields
}

func (col *Collection[T]) GetFields() contracts.Fields {
	return col.ToFields()
}

func (col *Collection[T]) ToArrayFields() []contracts.Fields {
	panic("")
}

func (col *Collection[T]) String() string {
	return col.ToJsonString()
}
