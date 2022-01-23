package collection

import "github.com/goal-web/supports/exceptions"

type Exception struct {
	exceptions.Exception
}

// MapException 遍历参数异常
type MapException struct {
	exceptions.Exception
}

// SortException 排序参数异常
type SortException struct {
	exceptions.Exception
}
