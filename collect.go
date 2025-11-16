package collection

import "github.com/goal-web/contracts"

type collect[T, P any] struct {
	contracts.Collection[T]
}

func NewCollect[T, P any](collection contracts.Collection[T]) contracts.Collect[T, P] {
	return &collect[T, P]{collection}
}

func Collect[T, P any](list []T) contracts.Collect[T, P] {
	return &collect[T, P]{Collection: New(list)}
}

func (col collect[T, P]) ToP(mapper func(T) P) []P {
	var list = make([]P, col.Len())

	col.Foreach(func(i int, t T) {
		list = append(list, mapper(t))
	})

	return list
}
