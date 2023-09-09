package condition

// IsTrue 根据条件isTrue，确认返回returnTrue 或者 returnFalse结果
func IsTrue[TReturn any](isTrue bool, returnTrue TReturn, returnFalse TReturn) TReturn {
	if isTrue {
		return returnTrue
	}
	return returnFalse
}
