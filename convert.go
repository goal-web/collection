package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"strconv"
)

func (this *Collection) ToIntArray() (results []int) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToInt(this.array[strconv.Itoa(i)], 0))
	}
	return
}

func (this *Collection) ToInt64Array() (results []int64) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToInt64(this.array[strconv.Itoa(i)], 0))
	}
	return
}
func (this *Collection) ToInterfaceArray() []interface{} {
	results := make([]interface{}, 0)
	for _, data := range this.array {
		results = append(results, data)
	}
	return results
}
func (this *Collection) ToFields() contracts.Fields {
	return this.array
}

func (this *Collection) ToFloatArray() (results []float32) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToFloat(this.array[strconv.Itoa(i)], 0))
	}
	return
}

func (this *Collection) ToFloat64Array() (results []float64) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToFloat64(this.array[strconv.Itoa(i)], 0))
	}
	return
}

func (this *Collection) ToBoolArray() (results []bool) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToBool(this.array[strconv.Itoa(i)], false))
	}
	return
}

func (this *Collection) ToStringArray() (results []string) {
	for i := 0; i < this.Len(); i++ {
		results = append(results, utils.ConvertToString(this.array[strconv.Itoa(i)], ""))
	}
	return
}
