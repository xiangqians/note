package org.xiangqian.note.util;

import java.util.HashMap;
import java.util.Map;

/**
 * 文件类型
 *
 * @author xiangqian
 * @date 19:54 2024/04/16
 */
public class Type {

    // <后缀名，文件类型>
    private static final Map<String, String> map = new HashMap<>(64, 1f);

    public static final String FOLDER = "folder";
    public static final String MD = "md";
    public static final String DOC = "doc";
    public static final String DOCX = "docx";
    public static final String PDF = "pdf";
    public static final String HTML = "html"; // HTML
    public static final String ZIP = "zip";

    public static final String XML = "xml"; // XML
    public static final String SVG = "svg"; // SVG
    public static final String CSS = "css"; // CSS
    public static final String JS = "js"; // javascript
    public static final String TS = "ts"; // typescript
    public static final String JAVA = "java"; // Java
    public static final String C = "c"; // C
    public static final String CPP = "cpp"; // C++
    public static final String CS = "cs"; // C#
    public static final String PY = "py"; // Python
    public static final String PHP = "php"; // PHP
    public static final String RB = "rb"; // Ruby
    public static final String GO = "go"; // Go
    public static final String SWIFT = "swift"; // Swift
    public static final String RS = "rs"; // Rust
    public static final String SQL = "sql"; // SQL
    public static final String YAML = "yaml"; // YAML
    public static final String JSON = "json"; // JSON

    static {
        add(MD, "md", "markdown");
        add(DOC, "doc");
        add(DOCX, "docx");
        add(PDF, "pdf");
        add(HTML, "html");
        add(ZIP, "zip");

        add(XML, "xml");
        add(SVG, "svg");
        add(CSS, "css");
        add(JS, "js");
        add(TS, "ts");
        add(JAVA, "java");
        add(C, "c");
        add(CPP, "cpp");
        add(CS, "cs");
        add(PY, "py");
        add(PHP, "php");
        add(RB, "rb");
        add(GO, "go");
        add(SWIFT, "swift");
        add(RS, "rs");
        add(SQL, "sql");
        add(YAML, "yaml", "yml");
        add(JSON, "json");
    }

    /**
     * @param type     文件类型
     * @param suffixes 文件后缀名集
     */
    private static void add(String type, String... suffixes) {
        for (String suffix : suffixes) {
            map.put(suffix, type);
        }
    }

    public static String suffixOf(String suffix) {
        return map.get(suffix);
    }

}
