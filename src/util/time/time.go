// @author xiangqian
// @date 12:33 2023/02/04
package time

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"time"
)

// HumanizUnix 人性化日期时间戳（s）
func HumanizUnix(unix int64) string {
	if unix <= 0 {
		return "-"
	}

	// 解析时间戳
	t := ParseUnix(unix)

	// 当前时间
	nowTime := NowTime()

	// 当前时间 - 解析后的时间戳
	duration := nowTime.Sub(t)

	// 小时
	hour := int64(duration.Hours())
	if hour >= 24 {
		return FormatTime(t)
	}
	if hour >= 1 {
		return fmt.Sprintf(i18n.MustGetMessage("i18n.nHoursAgo"), hour)
	}

	// 分钟
	minute := int64(duration.Minutes())
	if minute >= 1 {
		return fmt.Sprintf(i18n.MustGetMessage("i18n.nMinutesAgo"), minute)
	}

	// 秒
	second := int64(duration.Seconds())
	return fmt.Sprintf(i18n.MustGetMessage("i18n.nSecondsAgo"), second)
}

// ParseUnix 解析日期时间戳（s）
func ParseUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006/01/02 15:04:05")
}

// NowUnix 当前日期时间戳（s）
func NowUnix() int64 {
	return NowTime().Unix()
}

// NowTime 当前时间
func NowTime() time.Time {
	return time.Now()
}
