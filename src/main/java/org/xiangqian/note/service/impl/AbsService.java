package org.xiangqian.note.service.impl;

import org.springframework.beans.factory.annotation.Value;
import org.xiangqian.note.entity.IavEntity;
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

    private String dataDir;

    @Value("${spring.datasource.url}")
    public void setUrl(String url) {
        File file = new File(url.trim().substring("jdbc:sqlite:".length()));
        file = file.getParentFile();
        dataDir = file.getPath();
    }

    protected Path getTmpPath(Long id) throws IOException {
        return getTmpPath(id, false);
    }

    protected Path getTmpPath(Long id, boolean makeParentDirIfNotExist) throws IOException {
        return getPath(getParentDirName() + "tmp", id, makeParentDirIfNotExist);
    }

    protected Path getPath(Long id) throws IOException {
        return getPath(id, false);
    }

    protected Path getPath(Long id, boolean makeParentDirIfNotExist) throws IOException {
        return getPath(getParentDirName(), id, makeParentDirIfNotExist);
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
     * @param id                      {@link IavEntity#getId()}
     * @param parentDirName           父目录名称
     * @param makeParentDirIfNotExist 如果父目录不存在，则创建
     * @return
     */
    private Path getPath(String parentDirName, Long id, boolean makeParentDirIfNotExist) throws IOException {
        if (makeParentDirIfNotExist) {
            Path path = Paths.get(dataDir, parentDirName);
            if (!Files.exists(path)) {
                Files.createDirectories(path);
            }
        }
        return Paths.get(dataDir, parentDirName, id.toString());
    }

}
