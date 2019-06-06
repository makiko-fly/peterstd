// 本文件定义了时区相关的常量

package datetime

import "time"

// TimeZone
var (
	tzShanghai *time.Location
)

func init() {
	tzShanghai, _ = time.LoadLocation("Asia/Shanghai")
}

func TZShanghai() *time.Location {
	return tzShanghai
}
