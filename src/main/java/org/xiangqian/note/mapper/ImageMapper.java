package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.xiangqian.note.entity.ImageEntity;

/**
 * 图片映射器
 *
 * @author xiangqian
 * @date 23:30 2024/03/04
 */
@Mapper
public interface ImageMapper {

    /**
     * 新增图片信息
     *
     * @param entity
     * @return
     */
    Boolean create(ImageEntity entity);

    /**
     * 根据主键删除图片信息
     *
     * @param id
     * @return
     */
    Boolean deleteById(Long id);

    /**
     * 根据主键更新图片信息
     *
     * @param entity
     * @return
     */
    Boolean updateById(ImageEntity entity);

    /**
     * 根据主键获取图片信息
     *
     * @param id
     * @return
     */
    ImageEntity getById(Long id);

    /**
     * 获取已删除的图片主键
     *
     * @return
     */
    Long getDeletedId();

}
