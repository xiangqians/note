package org.xiangqian.note.service.impl;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.xiangqian.note.service.IavService;
import org.xiangqian.note.service.NoteService;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/**
 * @author xiangqian
 * @date 20:35 2024/04/23
 */
public abstract class AbsService {

    protected final MediaType IMAGE_WEBP = MediaType.valueOf("image/webp");
    protected final MediaType TEXT_CSS = MediaType.valueOf("text/css");
    protected final MediaType IMAGE_X_ICON = MediaType.valueOf("image/x-icon");

    protected final java.nio.charset.Charset UTF_8 = java.nio.charset.Charset.forName("UTF-8");
    protected final java.nio.charset.Charset GBK = java.nio.charset.Charset.forName("GBK");

    private String dataDir;

    @Value("${spring.datasource.url}")
    public void setUrl(String url) {
        File file = new File(url.trim().substring("jdbc:sqlite:".length()));
        file = file.getParentFile();
        dataDir = file.getPath();
    }

    protected Path getTmpPath(String name) throws IOException {
        return getTmpPath(name, false);
    }

    protected Path getTmpPath(String name, boolean makeParentDirIfNotExist) throws IOException {
        return getPath(getParentDirName() + "tmp", name, makeParentDirIfNotExist);
    }

    protected Path getPath(String name) throws IOException {
        return getPath(name, false);
    }

    protected Path getPath(String name, boolean makeParentDirIfNotExist) throws IOException {
        return getPath(getParentDirName(), name, makeParentDirIfNotExist);
    }

    private String getParentDirName() {
        if (this instanceof NoteService) {
            return "note";
        }
        if (this instanceof IavService) {
            return "iav";
        }
        throw new IllegalArgumentException(this.getClass().toString());
    }

    /**
     * 获取文件路径
     *
     * @param name                    文件/目录名称
     * @param parentDirName           父目录名称
     * @param makeParentDirIfNotExist 如果父目录不存在，则创建
     * @return
     */
    private Path getPath(String parentDirName, String name, boolean makeParentDirIfNotExist) throws IOException {
        if (makeParentDirIfNotExist) {
            Path path = Paths.get(dataDir, parentDirName);
            if (!Files.exists(path)) {
                Files.createDirectories(path);
            }
        }
        return Paths.get(dataDir, parentDirName, name);
    }

    protected ResponseEntity<Resource> notFound() {
        return ResponseEntity.notFound().build();
    }

    protected ResponseEntity<Resource> ok(Path path, MediaType contentType) throws IOException {
        // 读取文件
        byte[] bytes = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(bytes);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.setContentLength(resource.contentLength());
        headers.setContentType(contentType);

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

}
