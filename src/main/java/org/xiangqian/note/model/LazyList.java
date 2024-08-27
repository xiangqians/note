package org.xiangqian.note.model;

import lombok.Data;

import java.util.Collections;
import java.util.List;

/**
 * 延迟加载列表信息
 *
 * @author xiangqian
 * @date 22:33 2024/03/04
 */
@Data
public class LazyList<T> {

    /**
     * 当前页
     */
    private Integer current;

    /**
     * 页数量
     */
    private Integer size;

    /**
     * 是否有下一页数据
     */
    private Boolean next;

    /**
     * 数据
     */
    private List<T> data;

    public LazyList() {
        this(1, 10);
    }

    public LazyList(Integer current, Integer size) {
        this.current = current;
        this.size = size;
        this.next = false;
        this.data = Collections.emptyList();
    }

}
