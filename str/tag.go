package str

import (
	"fmt"
	"strings"
)

// ParseTag 在字符串内，根据开始标记、结束标记，找到标记内的内容
func ParseTag(body string, startTag, endTag string) (string, error) {
	index := strings.Index(body, startTag)
	if index == -1 {
		return "", fmt.Errorf("ParseTag：没有找到开始标记：%s", startTag)
	}
	body = body[index+len(startTag):]

	if endTag != "" {
		index = strings.Index(body, endTag)
		if index == -1 {
			return body, fmt.Errorf("ParseTag：没有找到结束标记：%s", endTag)
		}
		body = body[:index]
	}
	return body, nil
}

// Substring 获取从开始到endTag的内容
func Substring(body string, endTag string) (string, error) {
	index := strings.Index(body, endTag)
	if index == -1 {
		return body, fmt.Errorf("Substring：没有找到标记：%s", endTag)
	}
	return body[:index], nil
}
