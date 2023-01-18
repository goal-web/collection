package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"github.com/shopspring/decimal"
)

// SafeSum struct 或者 map 情况下需要传 key
func (col *Collection) SafeSum(key ...string) (sum decimal.Decimal) {
	sum = decimal.NewFromInt(0)
	if len(key) == 0 {
		for _, f := range col.ToFloat64Array() {
			sum.Add(decimal.NewFromFloat(f))
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			sum = sum.Add(decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])))
		})
	}
	return
}

// SafeAvg struct 或者 map 情况下需要传 key
func (col *Collection) SafeAvg(key ...string) (sum decimal.Decimal) {
	return col.SafeSum(key...).Div(decimal.NewFromInt32(int32(col.Count())))
}

// SafeMax struct 或者 map 情况下需要传 key
func (col *Collection) SafeMax(key ...string) (max decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range col.ToFloat64Array() {
			if max.IsZero() {
				max = decimal.NewFromFloat(f)
			} else if float := decimal.NewFromFloat(f); max.LessThan(float) {
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
func (col *Collection) SafeMin(key ...string) (min decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range col.ToFloat64Array() {
			if min.IsZero() {
				min = decimal.NewFromFloat(f)
			} else if float := decimal.NewFromFloat(f); float.LessThan(min) {
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

func (col *Collection) Count() int {
	return len(col.array)
}

func (col *Collection) Sum(key ...string) (sum float64) {
	if len(key) == 0 {
		for _, f := range col.ToFloat64Array() {
			sum += f
		}
	} else {
		col.Map(func(fields contracts.Fields) {
			sum += utils.GetFloat64Field(fields, key[0])
		})
	}
	return
}

func (col *Collection) Max(key ...string) (max float64) {
	if len(key) == 0 {
		for i, f := range col.ToFloat64Array() {
			if i == 0 {
				max = f
			} else if f > max {
				max = f
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

func (col *Collection) Min(key ...string) (min float64) {
	if len(key) == 0 {
		for i, f := range col.ToFloat64Array() {
			if i == 0 {
				min = f
			} else if f < min {
				min = f
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

func (col *Collection) Avg(key ...string) float64 {
	return col.Sum(key...) / float64(col.Count())
}
