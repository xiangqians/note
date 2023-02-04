// time
// @author xiangqian
// @date 12:33 2023/02/04
package util

import (
	"fmt"
	"github.com/gin-contrib/i18n"
	"time"
)

// UnixToTime 日期时间戳（s）转时间
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// FormatUnix 格式化日期时间戳（s）
func FormatUnix(unix int64) string {
	return UnixToTime(unix).Format("2006/01/02 15:04:05")
}

// 人性化日期时间戳（s）
func HumanizUnix(unix int64) string {
	if unix <= 0 {
		return "-"
	}

	t := UnixToTime(unix)
	now := time.Now()

	hour := int64(now.Sub(t).Hours())
	if hour >= 1 {
		if hour >= 24 {
			return FormatUnix(unix)
		}
		return fmt.Sprintf(i18n.MustGetMessage("i18n.xHourAgo"), hour)
	}

	minute := int64(now.Sub(t).Minutes())
	if minute >= 1 {
		return fmt.Sprintf(i18n.MustGetMessage("i18n.xMinuteAgo"), minute)
	}

	second := int64(now.Sub(t).Seconds())
	return fmt.Sprintf(i18n.MustGetMessage("i18n.xSecondAgo"), second)
}

// 人性化秒
func HumanizSecond(second int64) string {
	if second < 0 {
		return "-"
	}

	if second == 0 {
		return fmt.Sprintf("0%s", i18n.MustGetMessage("i18n.secondAbbrev"))
	}

	hour := second / (60 * 60)
	if hour > 24 {
		return "error"
	}

	mods := second % (60 * 60) // s, mod second
	minute := mods / 60
	second = mods % 60

	hourStr := ""
	minuteStr := ""
	secondStr := ""
	if hour != 0 {
		hourStr = fmt.Sprintf("%v%s", hour, i18n.MustGetMessage("i18n.hourAbbrev"))
	}
	if minute != 0 {
		minuteStr = fmt.Sprintf("%v%s", minute, i18n.MustGetMessage("i18n.minuteAbbrev"))
	}
	if second != 0 {
		secondStr = fmt.Sprintf("%v%s", second, i18n.MustGetMessage("i18n.secondAbbrev"))
	}

	return fmt.Sprintf("%s%s%s", hourStr, minuteStr, secondStr)
}
