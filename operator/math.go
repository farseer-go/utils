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

// GetTail 取尾数
func GetTail(num int) int {
	strNum := parse.Convert(num, "")
	if len(strNum) == 0 {
		return 0
	}
	return parse.Convert(strNum[(len(strNum)-1):], 0)
}
