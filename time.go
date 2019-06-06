package peterstd

import (
	"sync"
	"time"
)

// TimeZone
var (
	tzShanghai     *time.Location
	tzShanghaiInit sync.Once
)

func TZShanghai() *time.Location {
	tzShanghaiInit.Do(func() {
		tzShanghai, _ = time.LoadLocation("Asia/Shanghai")
	})
	return tzShanghai
}

const (
	SecondsPerHour   = 3600
	SecondsPerMinute = 60
)

func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}
