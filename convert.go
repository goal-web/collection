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

func (col *Collection) ToIntArray() (results []int) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToInt(data, 0))
	}
	return
}

func (col *Collection) ToInt64Array() (results []int64) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToInt64(data, 0))
	}
	return
}
func (col *Collection) ToInterfaceArray() []interface{} {
	return col.array
}

func (col *Collection) ToFloatArray() (results []float32) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToFloat(data, 0))
	}
	return
}

func (col *Collection) ToFloat64Array() (results []float64) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToFloat64(data, 0))
	}
	return
}

func (col *Collection) ToBoolArray() (results []bool) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToBool(data, false))
	}
	return
}

func (col *Collection) ToStringArray() (results []string) {
	for _, data := range col.array {
		results = append(results, utils.ConvertToString(data, ""))
	}
	return
}

func (col *Collection) ToFields() contracts.Fields {
	fields := contracts.Fields{}
	for index, data := range col.mapData {
		fields[strconv.Itoa(index)] = data
	}
	return fields
}

func (col *Collection) ToArrayFields() []contracts.Fields {
	return col.mapData
}

func (col *Collection) ToJson() string {
	results := make([]string, 0)
	col.Map(func(data interface{}) {
		if jsonify, isJson := data.(contracts.Json); isJson {
			results = append(results, jsonify.ToJson())
			return
		}
		jsonStr, err := json.Marshal(data)
		if err != nil {
			logs.WithError(err).WithFields(col.ToFields()).Fatal("json err")
		}
		results = append(results, string(jsonStr))
	})

	return fmt.Sprintf("[%s]", strings.Join(results, ","))
}

func (col *Collection) String() string {
	return col.ToJson()
}
