package org.xiangqian.note.util;

import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;

/**
 * 时间工具
 *
 * @author xiangqian
 * @date 21:24 2022/08/01
 */
public class TimeUtil {

    private static final DateTimeFormatter FORMATTER = DateTimeFormatter.ofPattern("yyyy/MM/dd HH:mm:ss");

    public static long now() {
        return LocalDateTime.now().atZone(ZoneId.systemDefault()).toEpochSecond();
    }

    /**
     * 人性化日期时间戳（单位s）
     *
     * @param second
     * @return
     */
    public static String human(Long second) {
        if (second == null || second <= 0) {
            return "";
        }

        LocalDateTime dateTime = LocalDateTime.ofInstant(Instant.ofEpochSecond(second), ZoneId.systemDefault());
        Duration duration = Duration.between(dateTime, LocalDateTime.now());

        // 天
        if (duration.toDays() > 0) {
            return FORMATTER.format(dateTime);
        }

        // 小时
        long hour = duration.toHours();
        if (hour > 0) {
            return hour + "小时前";
        }

        // 分钟
        long minute = duration.toMinutes();
        if (minute > 0) {
            return minute + "分钟前";
        }

        // 秒
        second = duration.toSeconds();
        return second + "秒前";
    }

    /**
     * 人性化间隔时间
     *
     * @param second
     * @return
     */
    public static String humanDuration(Long second) {
        if (second == null) {
            return "";
        }

        if (second <= 0) {
            return second + "秒";
        }

        StringBuilder builder = new StringBuilder();
        Duration duration = Duration.ofSeconds(second);

        // 小时
        long hour = duration.toHours();
        if (hour > 0) {
            builder.append(hour).append("小时");
        }

        // 分钟
        int minute = duration.toMinutesPart();
        if (minute > 0) {
            builder.append(minute).append("分钟");
        }

        // 秒
        int s = duration.toSecondsPart();
        if (s > 0) {
            builder.append(s).append("秒");
        }

        return builder.toString();
    }

}
