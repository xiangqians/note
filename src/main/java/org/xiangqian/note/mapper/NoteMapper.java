package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.xiangqian.note.configuration.mybatis.LazyList;
import org.xiangqian.note.entity.NoteEntity;

import java.util.List;

/**
 * @author xiangqian
 * @date 21:39 2024/03/03
 */
@Mapper
public interface NoteMapper {

    /**
     * 新增笔记信息
     *
     * @param entity
     * @return
     */
    Boolean create(NoteEntity entity);

    /**
     * 根据主键删除笔记信息
     *
     * @param id
     * @return
     */
    Boolean deleteById(Long id);

    /**
     * 根据主键修改笔记信息
     *
     * @param entity
     * @return
     */
    Boolean updateById(NoteEntity entity);

    /**
     * 根据主键获取笔记信息
     *
     * @param id
     * @return
     */
    NoteEntity getById(Long id);

    /**
     * 获取目录笔记列表
     * @return
     */
    List<NoteEntity> getFolderList();

    /**
     * 根据父节点主键和主键获取笔记信息
     *
     * @param pid
     * @param id
     * @return
     */
    NoteEntity getByPidAndId(@Param("pid") Long pid, @Param("id") Long id);

    /**
     * 获取已删除的笔记主键
     *
     * @return
     */
    Long getDeletedId();

    /**
     * 根据主键获取目录/文件大小
     *
     * @param id
     * @return
     */
    Long getSizeById(Long id);

    /**
     * 根据主键获取父节点列表信息（包括了当前节点信息）
     *
     * @param id
     * @return
     */
    List<NoteEntity> getParentListById(Long id);

    /**
     * 根据父节点主键获取子节点主键列表
     *
     * @param pid
     * @return
     */
    List<Long> getChildIdListByPid(Long pid);

    /**
     * 获取子节点列表
     *
     * @param lazyList
     * @param entity
     * @return
     */
    LazyList<NoteEntity> getChildList(LazyList lazyList, NoteEntity entity);

    /**
     * 根据父节点主键获取子节点数量
     *
     * @param pid
     * @return
     */
    Integer countChildrenByPid(Long pid);

    List<NoteEntity> getList();

}
