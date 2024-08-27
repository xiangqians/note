package org.xiangqian.note.model;

import lombok.Data;

import java.util.Collections;
import java.util.List;

/**
 * 分页信息
 *
 * @author xiangqian
 * @date 00:53 2022/06/12
 */
@Data
public class Page<T> {

    /**
     * 当前页
     */
    private Integer current;

    /**
     * 页数量
     */
    private Integer size;

    /**
     * 总数
     */
    private Integer total;

    /**
     * 数据
     */
    private List<T> data;

    public Page() {
        this(1, 10);
    }

    public Page(Integer current, Integer size) {
        this.current = current;
        this.size = size;
        this.total = 0;
        this.data = Collections.emptyList();
    }

}
