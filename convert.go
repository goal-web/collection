package collection

import "github.com/goal-web/supports/utils"

func (this *Collection) ToIntArray() (results []int) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToInt(data, 0))
	}
	return
}

func (this *Collection) ToInt64Array() (results []int64) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToInt64(data, 0))
	}
	return
}
func (this *Collection) ToInterfaceArray() []interface{} {
	return this.array
}

func (this *Collection) ToFloatArray() (results []float32) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToFloat(data, 0))
	}
	return
}

func (this *Collection) ToFloat64Array() (results []float64) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToFloat64(data, 0))
	}
	return
}

func (this *Collection) ToBoolArray() (results []bool) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToBool(data, false))
	}
	return
}

func (this *Collection) ToStringArray() (results []string) {
	for _, data := range this.array {
		results = append(results, utils.ConvertToString(data, ""))
	}
	return
}
