package org.xiangqian.note.model;

import java.util.HashMap;

/**
 * @author xiangqian
 * @date 19:37 2024/08/27
 */
public class Vo extends HashMap<String, Object> {

    private Vo() {
        super(8, 1f);
    }

    public Vo add(String name, Object value) {
        put(name, value);
        return this;
    }

    public static Vo none() {
        return new Vo();
    }

    public static Vo info(Object value) {
        return new Vo().add("info", value);
    }

    public static Vo error(Object value) {
        return new Vo().add("error", value);
    }

}
