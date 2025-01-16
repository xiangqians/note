package org.xiangqian.note.service;

import org.springframework.core.io.Resource;
import org.xiangqian.note.configuration.mybatis.LazyList;
import org.xiangqian.note.entity.NoteEntity;

import java.io.IOException;
import java.nio.file.Path;
import java.util.List;

/**
 * 笔记服务
 *
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
public interface NoteService {

    /**
     * 创建目录
     *
     * @param entity
     * @return
     */
    Boolean createFolder(NoteEntity entity) throws IOException;

    /**
     * 创建 MD 文件
     *
     * @param entity
     * @return
     */
    Boolean createMd(NoteEntity entity) throws IOException;

    /**
     * 上传PDF/ZIP文件
     *
     * @param entity
     * @return
     */
    Boolean upload(NoteEntity entity) throws IOException;

    /**
     * 根据主键删除笔记信息
     *
     * @param id
     * @return
     */
    Boolean deleteById(Long id);

    /**
     * 重命名
     *
     * @param entity
     * @return
     */
    Boolean rename(NoteEntity entity);

    /**
     * 粘贴
     *
     * @param entity
     * @return
     */
    Boolean paste(NoteEntity entity);

    /**
     * 根据主键更新MD文件内容
     *
     * @param id
     * @param content
     * @return
     * @throws IOException
     */
    Boolean updateContentById(Long id, String content) throws IOException;

    /**
     * 根据主键更新目录下文件的排序规则
     *
     * @param id
     * @param content
     * @return
     * @throws IOException
     */
    Boolean updateSortById(Long id, String content) throws IOException;

    /**
     * 根据主键获取MD文件内容
     *
     * @param id
     * @return
     */
    String getContentById(Long id) throws IOException;

    /**
     * 根据主键获取笔记信息
     *
     * @param id
     * @return
     */
    NoteEntity getById(Long id);

    /**
     * 根据id获取父节点列表
     *
     * @param id
     * @return
     */
    List<NoteEntity> getParentListById(Long id);

    /**
     * 获取子节点列表
     *
     * @param entity
     * @return
     */
    LazyList<NoteEntity> getChildList(NoteEntity entity) throws IOException;

    /**
     * 获取 ZIP 文件中的内容列表
     *
     * @param id
     * @return
     */
    Path getPath(Long id);

    /**
     * 根据主键获取文件资源（除了目录外）
     *
     * @param id
     * @param name
     * @return
     * @throws IOException
     */
    Resource getResourceById(Long id, String name) throws IOException;

}
