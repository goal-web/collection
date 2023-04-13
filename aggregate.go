package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"github.com/shopspring/decimal"
)

// SafeSum struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeSum(key ...string) (sum decimal.Decimal) {
	sum = decimal.NewFromInt(0)
	if len(key) == 0 {
		for _, f := range col.ToArray() {
			sum.Add(decimal.NewFromFloat(utils.ToFloat64(f, 0)))
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			sum = sum.Add(decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])))
		})
	}
	return
}

// SafeAvg struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeAvg(key ...string) (sum decimal.Decimal) {
	return col.SafeSum(key...).Div(decimal.NewFromInt32(int32(col.Count())))
}

// SafeMax struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeMax(key ...string) (max decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range col.ToArray() {
			value := utils.ToFloat64(f, 0)
			if max.IsZero() {
				max = decimal.NewFromFloat(value)
			} else if float := decimal.NewFromFloat(value); max.LessThan(float) {
				max = float
			}
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			if max.IsZero() {
				max = decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0]))
			} else if float := decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])); max.LessThan(float) {
				max = float
			}
		})
	}
	return
}

// SafeMin struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeMin(key ...string) (min decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range col.ToArray() {
			value := utils.ToFloat64(f, 0)
			if min.IsZero() {
				min = decimal.NewFromFloat(value)
			} else if float := decimal.NewFromFloat(value); float.LessThan(min) {
				min = float
			}
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			if min.IsZero() {
				min = decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0]))
			} else if float := decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])); float.LessThan(min) {
				min = float
			}
		})
	}
	return
}

func (col *Collection[T]) Count() int {
	if col != nil {
		return len(col.array)
	}
	return 0
}

func (col *Collection[T]) Sum(key ...string) (sum float64) {
	if len(key) == 0 {
		for _, f := range col.ToArray() {
			sum += utils.ToFloat64(f, 0)
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			sum += utils.GetFloat64Field(fields, key[0])
		})
	}
	return
}

func (col *Collection[T]) Max(key ...string) (max float64) {
	if len(key) == 0 {
		for i, f := range col.ToArray() {
			value := utils.ToFloat64(f, 0)
			if i == 0 {
				max = value
			} else if value > max {
				max = value
			}
		}
	} else {
		col.Map(func(fields contracts.Fields, index int) {
			if index == 0 {
				max = utils.GetFloat64Field(fields, key[0])
			} else if float := utils.GetFloat64Field(fields, key[0]); float > max {
				max = float
			}
		})
	}
	return
}

func (col *Collection[T]) Min(key ...string) (min float64) {
	if len(key) == 0 {
		for i, f := range col.ToArray() {
			value := utils.ToFloat64(f, 0)
			if i == 0 {
				min = value
			} else if value < min {
				min = value
			}
		}
	} else {
		col.Map(func(fields contracts.Fields, index int) {
			if index == 0 {
				min = utils.GetFloat64Field(fields, key[0])
			} else if float := utils.GetFloat64Field(fields, key[0]); float < min {
				min = float
			}
		})
	}
	return
}

func (col *Collection[T]) Avg(key ...string) float64 {
	return col.Sum(key...) / float64(col.Count())
}
