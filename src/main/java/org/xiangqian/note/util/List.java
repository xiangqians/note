package org.xiangqian.note.util;

import lombok.AccessLevel;
import lombok.Data;
import lombok.Setter;

import java.util.Collections;

/**
 * @author xiangqian
 * @date 22:33 2024/03/04
 */
@Data
public class List<T> {

    // 索引值，从1开始
    private Integer offset;

    // 行数
    @Setter(AccessLevel.NONE)
    private Integer rows;

    // 数据列表
    private java.util.List<T> data;

    // 索引集
    private java.util.List<Integer> offsets;

    public List() {
        this(1, 10);
    }

    public List(Integer offset, Integer rows) {
        this.offset = offset;
        this.rows = rows;
        this.offsets = Collections.emptyList();
    }

}
