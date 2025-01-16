package org.xiangqian.note.service;

import org.springframework.core.io.Resource;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.ImageEntity;

import java.io.IOException;

/**
 * @author xiangqian
 * @date 20:28 2024/04/23
 */
public interface ImageService {

    /**
     * 上传图片
     *
     * @param file
     * @return
     * @throws IOException
     */
    ImageEntity upload(MultipartFile file) throws IOException;

    /**
     * 根据主键获取图片文件资源
     *
     * @param id
     * @return
     * @throws IOException
     */
    Resource getResourceById(Long id) throws IOException;

}
