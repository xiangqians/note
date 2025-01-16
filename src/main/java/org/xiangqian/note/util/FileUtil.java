package org.xiangqian.note.util;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.StringUtils;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.FileTime;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Enumeration;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.zip.ZipException;
import java.util.zip.ZipFile;

/**
 * @author xiangqian
 * @date 11:49 2024/03/09
 */
@Slf4j
public class FileUtil {

    /**
     * 人性化字节
     * <p>
     * 1B(Byte) = 8b(bit)
     * 1KB = 1024B
     * 1MB = 1024KB
     * 1GB = 1024MB
     * 1TB = 1024GB
     *
     * @param b Byte
     * @return
     */
    public static String humanSize(Long b) {
        if (b == null) {
            return "";
        }

        if (b <= 0) {
            return b + "B";
        }

        StringBuilder builder = new StringBuilder();

        // GB
        long gb = b / (1024 * 1024 * 1024);
        if (gb > 0) {
            builder.append(gb).append("GB");
            b = b % (1024 * 1024 * 1024);
        }

        // MB
        long mb = b / (1024 * 1024);
        if (mb > 0) {
            builder.append(mb).append("MB");
            b = b % (1024 * 1024);
        }

        // KB
        long kb = b / 1024;
        if (kb > 0) {
            builder.append(kb).append("KB");
            b = b % 1024;
            return builder.toString();
        }

        // B
        if (b > 0) {
            builder.append(b).append("B");
        }

        return builder.toString();
    }

    private static ZipFile getZipFile(Path path) throws IOException {
        File file = path.toFile();
        ZipFile zipFile = null;
        try {
            zipFile = new ZipFile(file, StandardCharsets.UTF_8);
        } catch (ZipException e) {
            String message = e.getMessage();
            // java.util.zip.ZipException: invalid CEN header (bad entry name)
            if (StringUtils.containsIgnoreCase(message, "invalid")
                    && StringUtils.containsIgnoreCase(message, "CEN")
                    && StringUtils.containsIgnoreCase(message, "header")) {
                zipFile = new ZipFile(file, Charset.forName("GBK"));
            } else {
                throw e;
            }
        }
        return zipFile;
    }

    /**
     * 获取 ZIP 文件中的条目的列表
     *
     * @param path
     * @param name
     * @return
     */
    public static List<ZipEntry> getZipEntryList(Path path, String name) throws IOException {
        if ((!"".equals(name) && !name.endsWith("/")) || !Files.exists(path)) {
            return null;
        }

        ZipFile zipFile = null;
        try {
            zipFile = getZipFile(path);

            // ^test/([^/]+/?)
            // ^test/：匹配的字符串必须以 test/ 开头
            // ([^/]+)：匹配一个或多个非 / 的字符，确保在 test/ 后面有内容（没有 / 作为分隔符）
            // /?：匹配 零个或一个 / 字符

            Pattern pattern = Pattern.compile("^" + name + "([^/]+/?)");

            List<ZipEntry> list = new ArrayList<>();
            Enumeration<? extends java.util.zip.ZipEntry> zipEntries = zipFile.entries();
            while (zipEntries.hasMoreElements()) {
                java.util.zip.ZipEntry zipEntry = zipEntries.nextElement();
                String path1 = zipEntry.getName();
                Matcher matcher = pattern.matcher(path1);
                if (matcher.matches()) {
                    ZipEntry zipEntry1 = new ZipEntry();

                    String type = null;

                    String fullName = matcher.group(1);
                    String name1 = null;
                    if (fullName.endsWith("/")) {
                        name1 = fullName.substring(0, fullName.length() - 1);
                        type = Type.FOLDER;
                    } else {
                        name1 = fullName;
                        type = getType(name1);
                    }

                    zipEntry1.setName(name1);
                    zipEntry1.setType(type);
                    zipEntry1.setPath(path1);
                    zipEntry1.setSize(zipEntry.getSize());
                    FileTime lastModifiedTime = zipEntry.getLastModifiedTime();
                    if (lastModifiedTime != null) {
                        zipEntry1.setLastModifiedTime(lastModifiedTime.toMillis() / 1000);
                    }
                    list.add(zipEntry1);
                }
            }

            Collections.sort(list, (o1, o2) -> {
                boolean isFolder1 = Type.FOLDER.equals(o1.getType());
                boolean isFolder2 = Type.FOLDER.equals(o2.getType());

                // 目录排在前面，文件排在后面
                if (isFolder1 && !isFolder2) {
                    // o1 是目录，排前面
                    return -1;
                } else if (!isFolder1 && isFolder2) {
                    // o2 是目录，排前面
                    return 1;
                } else {
                    // 如果都是目录或者都是文件，则按照字母顺序排序
                    return o1.getName().compareTo(o2.getName());
                }
            });

            return list;
        } finally {
            IOUtils.closeQuietly(zipFile);
        }
    }

    /**
     * 获取 ZIP 文件中的条目
     *
     * @param path
     * @param name
     * @return
     * @throws IOException
     */
    public static byte[] getZipEntryContent(Path path, String name) throws IOException {
        if (name == null || name.endsWith("/") || !Files.exists(path)) {
            return null;
        }

        ZipFile zipFile = null;
        InputStream inputStream = null;
        try {
            zipFile = getZipFile(path);
            java.util.zip.ZipEntry entry = zipFile.getEntry(name);
            inputStream = zipFile.getInputStream(entry);
            return inputStream.readAllBytes();
        } finally {
            IOUtils.closeQuietly(inputStream, zipFile);
        }
    }

    public static final String TEXT = "text";
    public static final String IMAGE = "image";

    public static String getType(String name) {
        if (name == null) {
            return null;
        }

        if (name.endsWith("/")) {
            return Type.FOLDER;
        }

        name = name.toLowerCase();

        if (name.endsWith(".pdf")) {
            return Type.PDF;
        }

        if (StringUtils.endsWithAny(name, ".txt", ".xml", ".svg", ".html", ".css", ".js", ".ts",
                ".java", ".c", ".cpp", ".cs", ".py", ".php", ".rb", ".go", ".swift", ".rs",
                ".sql", ".properties", ".yaml", ".yml", ".json", ".gitignore",
                ".cmd", ".bat", ".sh", ".iml", ".ftl", ".log", ".cer",
                "/dockerfile", "/license") || StringUtils.equalsAny(name, "dockerfile", "license")) {
            return TEXT;
        }

        if (StringUtils.endsWithAny(name, ".png", ".jpg", ".gif", ".webp", ".ico")) {
            return IMAGE;
        }

        return null;
    }

}
