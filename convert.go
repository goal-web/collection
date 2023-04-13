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

func (col *Collection[T]) ToArray() []T {
	return col.array
}

func (col *Collection[T]) ToAnyArray() []any {
	var array = make([]any, col.Len())
	for i, data := range col.array {
		array[i] = data
	}
	return array
}

func (col *Collection[T]) ToFields() contracts.Fields {
	fields := contracts.Fields{}
	for index, data := range col.array {
		fields[strconv.Itoa(index)] = data
	}
	return fields
}

func (col *Collection[T]) ToArrayFields() []contracts.Fields {
	if col.mapArray == nil {
		col.mapArray = make([]contracts.Fields, col.Len())
		for i, data := range col.array {
			col.mapArray[i], _ = utils.ToFields(data)
		}
	}

	return col.mapArray
}

func (col *Collection[T]) ToJson() string {
	results := make([]string, 0)
	col.Map(func(data any) {
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

func (col *Collection[T]) String() string {
	return col.ToJson()
}
