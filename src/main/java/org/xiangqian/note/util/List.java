package org.xiangqian.note.util;

import lombok.AccessLevel;
import lombok.Data;
import lombok.Setter;

import java.util.AbstractCollection;
import java.util.Collection;
import java.util.Collections;
import java.util.Iterator;

/**
 * @author xiangqian
 * @date 22:33 2024/03/04
 */
@Data
public class List<T> extends AbstractCollection<T> {

    // 索引值，从1开始
    @Setter(AccessLevel.NONE)
    private Integer offset;

    // 行数
    @Setter(AccessLevel.NONE)
    private Integer rows;

    // 数据集
    private Collection<T> data;

    // 索引集
    private java.util.List<Integer> offsets;

    public List(Integer offset) {
        this.offset = offset;
        this.rows = 50;
        this.data = Collections.emptyList();
        this.offsets = Collections.emptyList();
    }

    public List(Collection<T> data) {
        this.data = data;
    }

    @Override
    public Iterator<T> iterator() {
        return data.iterator();
    }

    @Override
    public int size() {
        return data.size();
    }

}
