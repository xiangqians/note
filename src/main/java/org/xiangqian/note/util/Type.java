package org.xiangqian.note.util;

import org.apache.commons.lang3.StringUtils;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;

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
    public static final String XLS = "xls";
    public static final String XLSX = "xlsx";
    public static final String PDF = "pdf";
    public static final String HTML = "html"; // HTML
    public static final String ZIP = "zip";

    // 文档格式集
    private static final Set<String> documents = Set.of(DOC, DOCX, XLS, XLSX, PDF);

    // text
    public static final String TXT = "txt"; // TXT
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
    public static final String BAT = "bat"; // batch
    public static final String SH = "sh"; // Shell
    public static final String DOCKERFILE = "dockerfile"; // Dockerfile
    public static final String IML = "iml"; // IntelliJ IDEA项目文件，是一种XML格式的文件
    public static final String FTL = "ftl"; // FreeMarker模板文件，包含HTML代码
    public static final String LOG = "log";
    public static final String CER = "cer";

    private static final Set<String> texts = Set.of(MD, HTML,
            TXT, XML, SVG, CSS, JS, TS, JAVA, C, CPP, CS, PY, PHP, RB, GO, SWIFT, RS, SQL, PROPERTIES, YAML, JSON, GITIGNORE, LICENSE, CMD, BAT, SH, DOCKERFILE, IML, FTL, LOG, CER);

    // image
    public static final String PNG = "png";
    public static final String JPG = "jpg";
    public static final String GIF = "gif";
    public static final String WEBP = "webp";
    public static final String ICO = "ico";


    private static final Set<String> images = Set.of(PNG, JPG, GIF, WEBP, ICO);

    static {
        set(MD, "md", "markdown");
        set(DOC, "doc");
        set(DOCX, "docx");
        set(XLS, "xls");
        set(XLSX, "xlsx");
        set(PDF, "pdf");
        set(HTML, "html", "htm");
        set(ZIP, "zip");

        set(TXT, "txt");
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
        set(BAT, "bat");
        set(SH, "sh");
        set(DOCKERFILE, "dockerfile");
        set(IML, "iml");
        set(FTL, "ftl");
        set(LOG, "log");
        set(CER, "cer");

        set(PNG, "png");
        set(JPG, "jpg");
        set(GIF, "gif");
        set(WEBP, "webp");
        set(ICO, "ico");
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

    /**
     * 根据Path获取文件类型
     *
     * @param path
     * @return
     */
    public static String pathOf(Path path) {
        if (Files.isDirectory(path)) {
            return Type.FOLDER;
        }

        String name = StringUtils.trim(path.getFileName().toString());
        String suffix = null;
        int index = name.lastIndexOf(".");
        if (index != -1 && (index + 1) < name.length()) {
            suffix = StringUtils.trim(name.substring(index + 1)).toLowerCase();
        } else {
            suffix = name.trim().toLowerCase();
        }
        return Type.suffixOf(suffix);
    }

    /**
     * 是否是文本类型
     *
     * @param type
     * @return
     */
    public static boolean isText(String type) {
        if (type == null) {
            return false;
        }
        return texts.contains(type);
    }

    /**
     * 是否是文档格式
     *
     * @param type
     * @return
     */
    public static boolean isDocument(String type) {
        if (type == null) {
            return false;
        }
        return documents.contains(type);
    }

    /**
     * 是否是图片类型
     *
     * @param type
     * @return
     */
    public static boolean isImage(String type) {
        if (type == null) {
            return false;
        }
        return images.contains(type);
    }

}
