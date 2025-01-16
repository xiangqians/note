package org.xiangqian.note.entity;

import lombok.Data;

/**
 * 图片信息
 *
 * @author xiangqian
 * @date 23:29 2024/03/04
 */
@Data
public class ImageEntity {

    /**
     * id
     */
    private Long id;

    /**
     * 名称
     */
    private String name;

    /**
     * 类型，png、jpg、gif、webp、ico
     */
    private String type;

    /**
     * 大小，单位byte
     */
    private Long size;

    /**
     * 是否已删除，0-否，1-是
     */
    private Integer delete;

    /**
     * 创建时间戳，单位s
     */
    private Long createTime;

    /**
     * 修改时间戳，单位s
     */
    private Long updateTime;

}
