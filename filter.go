package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
)

// Filter 过滤不需要的数据 filter 返回 true 时保留
func (this *Collection[T]) Filter(filter func(item T, index int) bool) contracts.Collection[T] {
	results := make([]T, 0)
	for index, item := range this.rawData {
		if filter(item, index) {
			results = append(results, item)
		}
	}

	return Array(results)
}

// Skip 过滤不需要的数据 filter 返回 true 时过滤
func (this *Collection[T]) Skip(skip func(item T, index int) bool) contracts.Collection[T] {
	results := make([]T, 0)
	for index, item := range this.rawData {
		if !skip(item, index) {
			results = append(results, item)
		}
	}

	return Array(results)
}

// Skip 过滤不需要的数据 filter 返回 true 时过滤
func (this *Collection[T]) getField(index int, field string) any {
}

// Where 根据条件过滤数据，支持 =,>,>=,<,<=,in,not in 等条件判断
func (this *Collection[T]) Where(field string, args ...interface{}) contracts.Collection {
	results := make([]interface{}, 0)
	var (
		arg      interface{}
		results  []T
		operator = "="
	)
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		operator = args[0].(string)
		arg = args[1]
	}
	return this.Filter(func(item T, index int) bool {
		return utils.Compare()
	})
}

func (this *Collection[T]) WhereLt(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "lt", arg)
}
func (this *Collection[T]) WhereLte(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "lte", arg)
}
func (this *Collection[T]) WhereGt(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "gt", arg)
}
func (this *Collection[T]) WhereGte(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "gte", arg)
}
func (this *Collection[T]) WhereIn(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "in", arg)
}
func (this *Collection[T]) WhereNotIn(field string, arg interface{}) contracts.Collection {
	return this.Where(field, "not in", arg)
}
