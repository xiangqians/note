package org.xiangqian.note.model;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 登录记录
 *
 * @author xiangqian
 * @date 13:53 2024/08/29
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class LoginRecord {
    /**
     * 主机
     */
    private String host;

    /**
     * 时间戳（单位s）
     */
    private Long time;
}
