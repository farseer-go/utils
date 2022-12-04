package str

// CutRight 裁剪末尾标签
// str：要裁剪的原字符串
// lastTag：裁剪的字符串
func CutRight(str string, lastTag string) string {
	if str[len(str)-len(lastTag):] == lastTag {
		return str[0 : len(str)-len(lastTag)]
	}
	return str
}
