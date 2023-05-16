// time
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

	// time
	t := ParseUnix(unix)
	now := Now()

	// format
	format := func(i18nKey string, value int64) string {
		return fmt.Sprintf(i18n.MustGetMessage(i18nKey), value)
	}

	// Duration
	duration := now.Sub(t)

	// hour
	hour := int64(duration.Hours())
	if hour >= 1 {
		if hour >= 24 {
			return FormatTime(t)
		}
		return format("i18n.nHoursAgo", hour)
	}

	// minute
	minute := int64(duration.Minutes())
	if minute >= 1 {
		return format("i18n.nMinutesAgo", minute)
	}

	// second
	second := int64(duration.Seconds())
	return format("i18n.nSecondsAgo", second)
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
	return Now().Unix()
}

// Now 当前时间
func Now() time.Time {
	return time.Now()
}
