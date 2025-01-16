package org.xiangqian.note.configuration.mybatis;

import lombok.Data;

import java.util.AbstractList;
import java.util.Collections;
import java.util.List;

/**
 * 延迟加载列表
 *
 * @author xiangqian
 * @date 22:33 2024/03/04
 */
@Data
public class LazyList<T> extends AbstractList<T> {

    /**
     * 当前页
     */
    private Integer current;

    /**
     * 页数量
     */
    private Integer size;

    /**
     * 页数集
     */
    private List<Integer> pages;

    /**
     * 是否有更多页数
     */
    private Boolean more;

    /**
     * 数据
     */
    private List<T> data;

    public LazyList() {
        this(1);
    }

    public LazyList(Integer current) {
        this.current = current;
        this.size = 50;
        this.pages = Collections.emptyList();
        this.more = false;
        this.data = Collections.emptyList();
    }

    @Override
    public int size() {
        return data.size();
    }

    @Override
    public T get(int index) {
        return data.get(index);
    }

}
