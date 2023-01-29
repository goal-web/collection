package collection

import (
	"github.com/goal-web/contracts"
)

type Exception struct {
	Err      error
	previous contracts.Exception
}

func (e Exception) Error() string {
	return e.Err.Error()
}

func (e Exception) GetPrevious() contracts.Exception {
	return e.previous
}

// MapException 遍历参数异常
type MapException = Exception

// SortException 排序参数异常
type SortException = Exception
