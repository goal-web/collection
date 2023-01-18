package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

// Filter 过滤不需要的数据 filter 返回 true 时保留
func (col *Collection) Filter(filter interface{}) contracts.Collection {
	results := make([]interface{}, 0)
	newFields := make([]contracts.Fields, 0)
	for index, data := range col.Map(filter).ToInterfaceArray() {
		if utils.ConvertToBool(data, false) {
			if fields := col.mapData[index]; fields != nil {
				newFields = append(newFields, fields)
			}
			results = append(results, col.array[index])
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

// Skip 过滤不需要的数据 filter 返回 true 时过滤
func (col *Collection) Skip(filter interface{}) contracts.Collection {
	results := make([]interface{}, 0)
	newFields := make([]contracts.Fields, 0)
	for index, data := range col.Map(filter).ToInterfaceArray() {
		if !utils.ConvertToBool(data, false) {
			if fields := col.mapData[index]; fields != nil {
				newFields = append(newFields, fields)
			}
			results = append(results, col.array[index])
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

// Where 根据条件过滤数据，支持 =,>,>=,<,<=,in,not in 等条件判断
func (col *Collection) Where(field string, args ...interface{}) contracts.Collection {
	results := make([]interface{}, 0)
	var (
		arg      interface{}
		operator = "="
	)
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		operator = args[0].(string)
		arg = args[1]
	}
	newFields := make([]contracts.Fields, 0)
	for index, data := range col.Map(func(fields contracts.Fields) bool {
		return utils.Compare(fields[field], operator, arg)
	}).ToInterfaceArray() {
		if utils.ConvertToBool(data, false) {
			if fields := col.mapData[index]; fields != nil {
				newFields = append(newFields, fields)
			}
			results = append(results, col.array[index])
		}
	}
	return &Collection{
		mapData: newFields,
		array:   results,
	}
}

func (col *Collection) WhereLt(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "lt", arg)
}
func (col *Collection) WhereLte(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "lte", arg)
}
func (col *Collection) WhereGt(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "gt", arg)
}
func (col *Collection) WhereGte(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "gte", arg)
}
func (col *Collection) WhereIn(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "in", arg)
}
func (col *Collection) WhereNotIn(field string, arg interface{}) contracts.Collection {
	return col.Where(field, "not in", arg)
}
