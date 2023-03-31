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

	// hour
	hour := int64(now.Sub(t).Hours())
	if hour >= 1 {
		if hour >= 24 {
			return FormatUnix(unix)
		}
		return format("i18n.xHourAgo", hour)
	}

	// minute
	minute := int64(now.Sub(t).Minutes())
	if minute >= 1 {
		return format("i18n.xMinuteAgo", minute)
	}

	// second
	second := int64(now.Sub(t).Seconds())
	return format("i18n.xSecondAgo", second)
}

// FormatUnix 格式化日期时间戳（s）
func FormatUnix(unix int64) string {
	if unix <= 0 {
		return "-"
	}

	return FormatTime(ParseUnix(unix))
}

// ParseUnix 解析日期时间戳（s）
func ParseUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// NowUnix 当前日期时间戳（s）
func NowUnix() int64 {
	return Now().Unix()
}

// FormatTime 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006/01/02 15:04:05")
}

// Now 当前时间
func Now() time.Time {
	return time.Now()
}
