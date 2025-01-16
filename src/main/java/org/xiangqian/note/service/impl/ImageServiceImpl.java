package org.xiangqian.note.service.impl;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.ImageEntity;
import org.xiangqian.note.mapper.ImageMapper;
import org.xiangqian.note.service.ImageService;
import org.xiangqian.note.util.TimeUtil;
import org.xiangqian.note.util.Type;

import java.io.IOException;

/**
 * @author xiangqian
 * @date 20:33 2024/04/23
 */
@Service
public class ImageServiceImpl extends AbsService implements ImageService {

    @Autowired
    private ImageMapper mapper;

    public ImageServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public ImageEntity upload(MultipartFile file) throws IOException {
        // 判断上传文件是否有效
        Assert.isTrue(file != null && !file.isEmpty(), "无效的文件");

        // 上传文件名
        String originalFilename = file.getOriginalFilename();
        String name = StringUtils.trimToEmpty(originalFilename);

        // 文件后缀名
        String suffix = null;
        int index = name.lastIndexOf(".");
        if (index > 0 && (index + 1) < name.length()) {
            suffix = name.substring(index + 1).toLowerCase().trim();
            name = name.substring(0, index);
        }

        Assert.isTrue(StringUtils.isNotEmpty(name), "文件名称不能为空");
        Assert.isTrue(StringUtils.isNotEmpty(suffix) && StringUtils.equalsAny(suffix, Type.PNG, Type.JPG, Type.GIF, Type.WEBP, Type.ICO), "上传的文件类型不受支持，请使用 PNG、JPG、GIF、WEBP 或 ICO 格式的图片");

        ImageEntity entity = new ImageEntity();
        entity.setName(name);
        entity.setType(suffix);
        entity.setSize(file.getSize());

        if (create(entity)) {
            // 写入文件
            writeFile(entity.getId(), file.getBytes());
            return entity;
        }

        return null;
    }

    @Override
    public Resource getResourceById(Long id) throws IOException {
        if (!isValidId(id)) {
            return null;
        }

        ImageEntity entity = mapper.getById(id);
        if (entity == null) {
            return null;
        }

        return createResource(String.format("%s.%s", entity.getName(), entity.getType()), entity.getType(), readFile(id));
    }

    private Boolean create(ImageEntity entity) {
        entity.setId(null);
        entity.setDelete(0);
        entity.setCreateTime(TimeUtil.now());
        entity.setUpdateTime(0L);

        // 获取已删除的图片主键，以复用
        Long deletedId = mapper.getDeletedId();
        if (deletedId != null) {
            entity.setId(deletedId);
            return mapper.updateById(entity);
        }

        return mapper.create(entity);
    }

}
