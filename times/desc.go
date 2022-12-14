package times

import (
	"github.com/farseer-go/fs/parse"
	"time"
)

// GetDesc 返回时间中文的描述
//
//	return "1 小时 32 分"
func GetDesc(ts time.Duration) string {
	days, hours, minutes, seconds := GetTime(ts)

	if days >= 1 {
		return parse.Convert(days, "0") + " 天 " + parse.Convert(hours, "0") + " 小时"
	}

	if hours >= 1 {
		return parse.Convert(hours, "0") + " 小时 " + parse.Convert(minutes, "0") + " 分"
	}

	if minutes >= 1 {
		return parse.Convert(minutes, "0") + " 分 " + parse.Convert(seconds, "0") + " 秒"
	}

	return parse.Convert(seconds, "0") + " 秒"
}

// GetSubDesc 两个时间相减，返回时间中文的描述
//
//	return "1 小时 32 分"
func GetSubDesc(ts1 time.Time, ts2 time.Time) string {
	return GetDesc(ts1.Sub(ts2))
}
