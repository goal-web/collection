package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

// Filter 过滤不需要的数据 filter 返回 true 时保留
func (col *Collection[T]) Filter(filter func(int, T) bool) contracts.Collection[T] {
	results := make([]T, 0)
	for i, item := range col.array {
		if filter(i, item) {
			results = append(results, item)
		}
	}
	return New(results)
}

// Skip 过滤不需要的数据 skipper 返回 true 时过滤
func (col *Collection[T]) Skip(skipper func(int, T) bool) contracts.Collection[T] {
	results := make([]T, 0)
	for i, item := range col.array {
		if !skipper(i, item) {
			results = append(results, item)
		}
	}
	return New(results)
}

// Foreach 遍历
func (col *Collection[T]) Foreach(handle func(int, T)) contracts.Collection[T] {
	for i, item := range col.array {
		handle(i, item)
	}
	return col
}

// Each 遍历并且返回新的集合
func (col *Collection[T]) Each(handle func(int, T) T) contracts.Collection[T] {
	results := make([]T, col.Len())
	for i, item := range col.array {
		results[i] = handle(i, item)
	}
	return New(results)
}

// Where 根据条件过滤数据，支持 =,>,>=,<,<=,in,not in 等条件判断
func (col *Collection[T]) Where(field string, args ...any) contracts.Collection[T] {
	var arg any
	var operator = "="
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		operator = args[0].(string)
		arg = args[1]
	}
	return col.Map(func(fields contracts.Fields) bool {
		return utils.Compare(fields[field], operator, arg)
	})
}

func (col *Collection[T]) WhereLt(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "lt", arg)
}
func (col *Collection[T]) WhereLte(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "lte", arg)
}
func (col *Collection[T]) WhereGt(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "gt", arg)
}
func (col *Collection[T]) WhereGte(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "gte", arg)
}
func (col *Collection[T]) WhereIn(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "in", arg)
}
func (col *Collection[T]) WhereNotIn(field string, arg any) contracts.Collection[T] {
	return col.Where(field, "not in", arg)
}
