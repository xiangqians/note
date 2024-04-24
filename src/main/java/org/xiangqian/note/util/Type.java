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
    public static final String MD = "md"; // Markdown
    public static final String DOC = "doc";
    public static final String DOCX = "docx";
    public static final String PDF = "pdf";
    public static final String HTML = "html"; // HTML
    public static final String ZIP = "zip";

    // text
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
    public static final String PROPERTIES = "properties";
    public static final String YAML = "yaml"; // YAML
    public static final String JSON = "json"; // JSON
    public static final String GITIGNORE = "gitignore"; // .gitignore
    public static final String LICENSE = "license"; // LICENSE
    public static final String CMD = "cmd"; // CMD
    public static final String SH = "sh"; // Shell
    public static final String DOCKERFILE = "dockerfile"; // Dockerfile
    public static final String IML = "iml"; // IntelliJ IDEA项目文件，是一种XML格式的文件
    public static final String FTL = "ftl"; // FreeMarker模板文件，包含HTML代码

    // image
    public static final String PNG = "png";
    public static final String JPG = "jpg";

    static {
        set(MD, "md", "markdown");
        set(DOC, "doc");
        set(DOCX, "docx");
        set(PDF, "pdf");
        set(HTML, "html");
        set(ZIP, "zip");

        set(XML, "xml");
        set(SVG, "svg");
        set(CSS, "css");
        set(JS, "js");
        set(TS, "ts");
        set(JAVA, "java");
        set(C, "c");
        set(CPP, "cpp");
        set(CS, "cs");
        set(PY, "py");
        set(PHP, "php");
        set(RB, "rb");
        set(GO, "go");
        set(SWIFT, "swift");
        set(RS, "rs");
        set(SQL, "sql");
        set(PROPERTIES, "properties");
        set(YAML, "yaml", "yml");
        set(JSON, "json");
        set(GITIGNORE, "gitignore");
        set(LICENSE, "license");
        set(CMD, "cmd");
        set(SH, "sh");
        set(DOCKERFILE, "dockerfile");
        set(IML, "iml");
        set(FTL, "ftl");

        set(PNG, "png");
        set(JPG, "jpg");
    }

    /**
     * 设置文件类型后缀名集
     *
     * @param type     文件类型
     * @param suffixes 文件后缀名集
     */
    private static void set(String type, String... suffixes) {
        for (String suffix : suffixes) {
            map.put(suffix, type);
        }
    }

    /**
     * 根据文件后缀名获取文件类型
     *
     * @param suffix 文件后缀名
     * @return
     */
    public static String suffixOf(String suffix) {
        return map.get(suffix);
    }

}
