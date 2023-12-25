package str

import "time"

// ToDateTime 将时间转换为yyyy-MM-dd HH:mm:ss
func ToDateTime(t time.Time) string {
	return t.Format(time.DateTime)
}
