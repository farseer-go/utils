package operator

import "reflect"

// IsTrue 三元操作运算
// conditional == true，return trueResult
// conditional == false，return falseResult
func IsTrue[TResult any](conditional bool, trueResult TResult, falseResult TResult) TResult {
	if conditional {
		return trueResult
	}
	return falseResult
}

// NotEmpty 三元操作运算。
// result != ""，return result
// else return emptyResult
func NotEmpty(result string, emptyResult string) string {
	if result != "" {
		return result
	}
	return emptyResult
}

// 取出第一个非nil值
func GetNotNil[T any](ags ...T) T {
	for _, ag := range ags {
		v := reflect.ValueOf(ag)
		// 检查变量是否有效且不为 nil
		if v.IsValid() && !v.IsNil() {
			return ag
		}
	}
	var zero T
	return zero
}

// 存在false值则返回true
func ExistsFalse(ags ...bool) bool {
	for _, ag := range ags {
		if ag == false {
			return true
		}
	}
	return false
}

// 存在true值则返回true
func ExistsTrue(ags ...bool) bool {
	for _, ag := range ags {
		if ag == true {
			return true
		}
	}
	return false
}
