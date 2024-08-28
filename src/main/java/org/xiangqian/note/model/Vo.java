package org.xiangqian.note.model;

import java.util.HashMap;

/**
 * @author xiangqian
 * @date 19:37 2024/08/27
 */
public class Vo extends HashMap<String, Object> {

    public Vo(Object error) {
        super(8, 1f);
        put("error", error);
    }

    public Vo() {
        this(null);
    }

    public Vo add(String name, Object value) {
        put(name, value);
        return this;
    }

}
