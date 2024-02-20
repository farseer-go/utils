package operator

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
