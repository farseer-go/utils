package operator

import (
	"github.com/farseer-go/fs/parse"
)

// GetSum 取和数
func GetSum(num int) int {
	strNum := parse.Convert(num, "")
	if len(strNum) == 0 {
		return 0
	}
	var sum int
	for _, value := range strNum {
		sum += parse.Convert(value, 0)
	}
	return sum
}

// GetTail 取数字的个位数
func GetTail(num int) int {
	//如果数值是个位数直接返回当前数值
	if num < 10 {
		return num
	}
	strNum := parse.Convert(num, "")
	return parse.Convert(strNum[(len(strNum)-1):], 0)
}

// GetHead 取数字第一位数
func GetHead(num int) int {
	strNum := parse.Convert(num, "")
	return parse.ToInt(strNum[0])
}
