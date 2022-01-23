package collection

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/utils"
	"github.com/shopspring/decimal"
)

// Sum struct 或者 map 情况下需要传 key
func (this *Collection) Sum(key ...string) (sum decimal.Decimal) {
	sum = decimal.NewFromInt(0)
	if len(key) == 0 {
		for _, f := range this.ToFloat64Array() {
			sum.Add(decimal.NewFromFloat(f))
		}
	} else {
		this.Map(func(fields contracts.Fields) {
			sum = sum.Add(decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])))
		})
	}
	return
}

// Avg struct 或者 map 情况下需要传 key
func (this *Collection) Avg(key ...string) (sum decimal.Decimal) {
	return this.Sum(key...).Div(decimal.NewFromInt32(int32(this.Count())))
}

// Max struct 或者 map 情况下需要传 key
func (this *Collection) Max(key ...string) (max decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range this.ToFloat64Array() {
			if max.IsZero() {
				max = decimal.NewFromFloat(f)
			} else if float := decimal.NewFromFloat(f); max.LessThan(float) {
				max = float
			}
		}
	} else {
		this.Map(func(fields contracts.Fields) {
			if max.IsZero() {
				max = decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0]))
			} else if float := decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])); max.LessThan(float) {
				max = float
			}
		})
	}
	return
}

// Min struct 或者 map 情况下需要传 key
func (this *Collection) Min(key ...string) (min decimal.Decimal) {
	if len(key) == 0 {
		for _, f := range this.ToFloat64Array() {
			if min.IsZero() {
				min = decimal.NewFromFloat(f)
			} else if float := decimal.NewFromFloat(f); float.LessThan(min) {
				min = float
			}
		}
	} else {
		this.Map(func(fields contracts.Fields) {
			if min.IsZero() {
				min = decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0]))
			} else if float := decimal.NewFromFloat(utils.GetFloat64Field(fields, key[0])); float.LessThan(min) {
				min = float
			}
		})
	}
	return
}

func (this *Collection) Count() int {
	return len(this.array)
}
