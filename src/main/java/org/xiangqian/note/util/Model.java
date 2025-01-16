package org.xiangqian.note.util;

import java.util.HashMap;
import java.util.Map;

/**
 * @author xiangqian
 * @date 11:38 2024/11/10
 */
public class Model {

    public static <K, V> Map<K, V> of(K k1, V v1) {
        Map<K, V> map = new HashMap<>(1, 1f);
        map.put(k1, v1);
        return map;
    }

    public static <K, V> Map<K, V> of(K k1, V v1, K k2, V v2) {
        Map<K, V> map = new HashMap<>(2, 1f);
        map.put(k1, v1);
        map.put(k2, v2);
        return map;
    }

    public static <K, V> Map<K, V> of(K k1, V v1, K k2, V v2, K k3, V v3) {
        Map<K, V> map = new HashMap<>(3, 1f);
        map.put(k1, v1);
        map.put(k2, v2);
        map.put(k3, v3);
        return map;
    }

}
