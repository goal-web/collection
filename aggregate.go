package collection

import (
	"github.com/shopspring/decimal"
)

// SafeSum struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeSum() (sum decimal.Decimal) {
	sum = decimal.NewFromFloat(0)
	for _, f := range col.ToFloat64Array() {
		sum.Add(decimal.NewFromFloat(f))
	}
	return
}

// SafeAvg struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeAvg() (sum decimal.Decimal) {
	return col.SafeSum().Div(decimal.NewFromInt32(int32(col.Len())))
}

// SafeMax struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeMax() (max decimal.Decimal) {
	for _, f := range col.ToFloat64Array() {
		if max.IsZero() {
			max = decimal.NewFromFloat(f)
		} else if float := decimal.NewFromFloat(f); max.LessThan(float) {
			max = float
		}
	}
	return
}

// SafeMin struct 或者 map 情况下需要传 key
func (col *Collection[T]) SafeMin() (min decimal.Decimal) {
	for _, f := range col.ToFloat64Array() {
		if min.IsZero() {
			min = decimal.NewFromFloat(f)
		} else if float := decimal.NewFromFloat(f); float.LessThan(min) {
			min = float
		}
	}
	return
}

func (col *Collection[T]) Sum() (sum float64) {
	for _, f := range col.ToFloat64Array() {
		sum += f
	}
	return
}

func (col *Collection[T]) Max() (max float64) {
	for i, f := range col.ToFloat64Array() {
		if i == 0 {
			max = f
		} else if f > max {
			max = f
		}
	}
	return
}

func (col *Collection[T]) Min() (min float64) {
	for i, f := range col.ToFloat64Array() {
		if i == 0 {
			min = f
		} else if f < min {
			min = f
		}
	}
	return
}

func (col *Collection[T]) Avg() float64 {
	return col.Sum() / float64(col.Len())
}
