package org.xiangqian.note.service.impl;

import org.springframework.core.env.Environment;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.xiangqian.note.service.GetNoteService;
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
    protected final MediaType IMAGE_X_ICON = MediaType.valueOf("image/x-icon");
    protected final MediaType TEXT_CSS = MediaType.valueOf("text/css");

    protected final java.nio.charset.Charset UTF_8 = java.nio.charset.Charset.forName("UTF-8");
    protected final java.nio.charset.Charset GBK = java.nio.charset.Charset.forName("GBK");

    // 数据路径
    private final Path dataPath;

    // 临时路径
    private final Path tmpPath;

    protected AbsService(Environment environment) throws IOException {
        String name = null;
        if (this instanceof NoteService || this instanceof GetNoteService) {
            name = "note";
        } else if (this instanceof IavService) {
            name = "iav";
        } else {
            throw new IllegalArgumentException(this.getClass().toString());
        }

        // 数据路径
        String url = environment.getProperty("spring.datasource.url");
        File file = new File(url.trim().substring("jdbc:sqlite:".length()));
        Path path = file.getParentFile().toPath().resolve(name);
        createDirectoriesIfNotExist(path);
        this.dataPath = path;

        // 临时路径
        path = Paths.get("tmp", name);
        createDirectoriesIfNotExist(path);
        this.tmpPath = path;
    }

    private void createDirectoriesIfNotExist(Path path) throws IOException {
        if (!Files.exists(path)) {
            Files.createDirectories(path);
        }
    }

    protected Path getDataPath(String name) throws IOException {
        return dataPath.resolve(name);
    }

    protected Path getTmpPath(String name) throws IOException {
        return tmpPath.resolve(name);
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
