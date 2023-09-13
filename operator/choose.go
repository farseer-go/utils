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
