package org.xiangqian.note.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.List;

/**
 * @author xiangqian
 * @date 21:39 2024/03/03
 */
@Mapper
public interface NoteMapper extends BaseMapper<NoteEntity> {

    /**
     * 根据id获取目录文件大小
     *
     * @param id
     * @return
     */
    Long getSizeById(@Param("id") Long id);

    NoteEntity getById(@Param("id") Long id);

    List<NoteEntity> list(@Param("entity") NoteEntity entity, @Param("list") List list);

}
