package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.xiangqian.note.entity.ImageEntity;

/**
 * 图片映射
 *
 * @author xiangqian
 * @date 23:30 2024/03/04
 */
@Mapper
public interface ImageMapper {

    /**
     * 新增图片信息
     *
     * @param imageEntity
     * @return
     */
    Boolean add(ImageEntity imageEntity);

    /**
     * 根据id更新图片信息
     *
     * @param imageEntity
     * @return
     */
    Boolean updById(ImageEntity imageEntity);

    /**
     * 根据id删除图片信息
     *
     * @param id
     * @return
     */
    Boolean delById(Long id);

    /**
     * 根据id获取图片信息
     *
     * @param id
     * @return
     */
    ImageEntity getById(Long id);

    /**
     * 获取已删除的id，以复用
     *
     * @return
     */
    Long getDeledId();

}
