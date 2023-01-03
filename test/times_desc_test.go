package test

import (
	"github.com/farseer-go/utils/times"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetSubDesc(t *testing.T) {
	ts1 := time.Date(2022, 8, 1, 20, 54, 12, 0, time.Local)
	ts2 := time.Date(2022, 8, 1, 19, 22, 12, 0, time.Local)

	desc := times.GetSubDesc(ts1, ts2)
	assert.Equal(t, "1 小时 32 分", desc)

	ts1 = time.Date(2022, 8, 2, 20, 54, 12, 0, time.Local)
	ts2 = time.Date(2022, 8, 1, 19, 22, 12, 0, time.Local)
	desc = times.GetSubDesc(ts1, ts2)
	assert.Equal(t, "1 天 1 小时", desc)

	ts1 = time.Date(2022, 8, 1, 19, 54, 19, 0, time.Local)
	ts2 = time.Date(2022, 8, 1, 19, 22, 12, 0, time.Local)
	desc = times.GetSubDesc(ts1, ts2)
	assert.Equal(t, "32 分 7 秒", desc)

	ts1 = time.Date(2022, 8, 1, 19, 22, 19, 99, time.Local)
	ts2 = time.Date(2022, 8, 1, 19, 22, 12, 0, time.Local)
	desc = times.GetSubDesc(ts1, ts2)
	assert.Equal(t, "7 秒", desc)
}
