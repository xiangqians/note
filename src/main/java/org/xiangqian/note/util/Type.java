package org.xiangqian.note.util;

import java.util.LinkedHashSet;
import java.util.Set;

/**
 * 文件类型
 *
 * @author xiangqian
 * @date 19:54 2024/04/16
 */
public class Type {

    public static final String FOLDER = "folder";
    public static final String MD = "md";
    public static final String DOC = "doc";
    public static final String DOCX = "docx";
    public static final String PDF = "pdf";
    public static final String ZIP = "zip";

    private static final Set<String> set;

    static {
        set = new LinkedHashSet<>(6, 1f);
        set.add(Type.FOLDER);
        set.add(Type.MD);
        set.add(Type.DOC);
        set.add(Type.DOCX);
        set.add(Type.PDF);
        set.add(Type.ZIP);
    }

    public static Set<String> getSet() {
        return set;
    }

}
