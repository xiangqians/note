package org.xiangqian.note.util;

import lombok.extern.slf4j.Slf4j;

/**
 * @author xiangqian
 * @date 11:49 2024/03/09
 */
@Slf4j
public class OsUtil {

    /**
     * 人性化字节
     * <p>
     * 1B(Byte) = 8b(bit)
     * 1KB = 1024B
     * 1MB = 1024KB
     * 1GB = 1024MB
     * 1TB = 1024GB
     *
     * @param b    Byte
     * @param unit 单位：GB、MB、KB、B
     * @return
     */
    public static String humanByte(Long b, String unit) {
        if (b == null) {
            return "";
        }

        if (b <= 0) {
            return b + "B";
        }

        StringBuilder builder = new StringBuilder();

        // GB
        long gb = b / (1024 * 1024 * 1024);
        if (gb > 0) {
            builder.append(gb).append("GB");
            b = b % (1024 * 1024 * 1024);
            if ("GB".equals(unit)) {
                return builder.toString();
            }
        }

        // MB
        long mb = b / (1024 * 1024);
        if (mb > 0) {
            builder.append(mb).append("MB");
            b = b % (1024 * 1024);
            if ("MB".equals(unit)) {
                return builder.toString();
            }
        }

        // KB
        long kb = b / 1024;
        if (kb > 0) {
            builder.append(kb).append("KB");
            b = b % 1024;
            if ("KB".equals(unit)) {
                return builder.toString();
            }
        }

        // B
        if (b > 0) {
            builder.append(b).append("B");
        }

        return builder.toString();
    }

}
