package operator

import (
	"github.com/farseer-go/fs/parse"
)

// GetSum 取和数
func GetSum(num int) int {
	sum := 0
	// 迭代每一位数字
	for num != 0 {
		digit := num % 10 // 每次取个位数
		sum += digit      // 将个位数加到总和上
		num /= 10         // 去掉个位数
	}
	return sum
}

// GetTail 取数字的个位数
func GetTail(num int) int {
	// gpt给的
	return num % 10
}

// GetHead 取数字第一位数
func GetHead(num int) int {
	n:=string(parse.ToString(num)[0])
	return parse.ToInt(n)
}
