package org.xiangqian.note.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.http.MediaType;
import org.xiangqian.note.service.ImageService;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.ResourceImpl;
import org.xiangqian.note.util.Type;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

/**
 * @author xiangqian
 * @date 20:35 2024/04/23
 */
@Slf4j
public abstract class AbsService {

    /**
     * 文件目录
     */
    private final Path dir;

    protected AbsService(Environment environment) throws IOException {
        String name = null;
        if (this instanceof NoteService) {
            name = "note";
        } else if (this instanceof ImageService) {
            name = "image";
        } else {
            throw new IllegalArgumentException(this.getClass().toString());
        }

        // 文件目录
        String url = environment.getProperty("spring.datasource.url");
        File file = new File(url.trim().substring("jdbc:sqlite:".length()));
        Path dir = file.getParentFile().toPath().resolve(name);
        if (!Files.exists(dir)) {
            Files.createDirectories(dir);
        }
        this.dir = dir;
        log.debug("{}文件目录：{}", this.getClass().getSimpleName(), dir.toAbsolutePath());
    }

    protected boolean isValidId(Long id) {
        return id != null && id.longValue() > 0;
    }

    /**
     * 获取文件路径
     *
     * @param id
     * @return
     */
    public Path getPath(Long id) {
        return dir.resolve(id.toString());
    }

    protected byte[] readFile(Long id) throws IOException {
        // 文件路径
        Path path = getPath(id);

        // 判断文件是否存在
        if (!Files.exists(path)) {
            return null;
        }

        // 读取文件
        return Files.readAllBytes(path);
    }

    /**
     * 写入文件
     *
     * @param id
     * @param bytes
     * @throws IOException
     */
    protected void writeFile(Long id, byte[] bytes) throws IOException {
        // 文件路径
        Path path = getPath(id);

        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, bytes);
    }

    /**
     * 删除文件
     *
     * @param id
     * @throws IOException
     */
    protected void deleteFile(Long id) throws IOException {
        // 文件路径
        Path path = getPath(id);

        // 如果文件存在，则删除
        Files.deleteIfExists(path);
    }

    /**
     * 创建资源
     *
     * @param name
     * @param type
     * @param bytes
     * @return
     * @throws IOException
     */
    protected Resource createResource(String name, String type, byte[] bytes) throws IOException {
        if (bytes == null) {
            bytes = new byte[0];
        }

        return new ResourceImpl(name, type != null ? switch (type) {
            case Type.MD -> MediaType.TEXT_XML;
            case Type.PDF -> MediaType.APPLICATION_PDF;
            case Type.ZIP -> MediaType.APPLICATION_OCTET_STREAM;
            case Type.PNG -> MediaType.IMAGE_PNG;
            case Type.JPG -> MediaType.IMAGE_JPEG;
            case Type.GIF -> MediaType.IMAGE_GIF;
            case Type.WEBP -> MediaType.valueOf("image/webp");
            case Type.ICO -> MediaType.valueOf("image/x-icon");
            default -> null;
        } : null, bytes);
    }

}
