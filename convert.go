package collection

import (
	"encoding/json"
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"strconv"
	"strings"
)

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

func (this *Collection) ToJson() string {
	results := make([]string, 0)
	this.Map(func(data interface{}) {
		if jsonify, isJson := data.(contracts.Json); isJson {
			results = append(results, jsonify.ToJson())
			return
		}
		jsonStr, err := json.Marshal(data)
		if err != nil {
			logs.WithError(err).WithFields(this.ToFields()).Fatal("json err")
		}
		results = append(results, string(jsonStr))
	})

	return fmt.Sprintf("[%s]", strings.Join(results, ","))
}

func (this *Collection) String() string {
	return this.ToJson()
}
