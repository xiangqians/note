package org.xiangqian.note.service.impl;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.IavEntity;
import org.xiangqian.note.mapper.IavMapper;
import org.xiangqian.note.service.IavService;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.LocalDateTime;

/**
 * @author xiangqian
 * @date 20:33 2024/04/23
 */
@Service
public class IavServiceImpl extends AbsService implements IavService {

    @Autowired
    private IavMapper mapper;

    protected IavServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Override
    public ResponseEntity<Resource> getStreamById(Long id) throws IOException {
        IavEntity entity = getById(id);
        if (entity == null) {
            return notFound();
        }

        // 文件路径
        Path path = getPath(id.toString());
        // 判断文件是否存在
        if (!Files.exists(path)) {
            return notFound();
        }

        String type = entity.getType();
        MediaType contentType = switch (type) {
            case Type.PNG -> MediaType.IMAGE_PNG;
            case Type.JPG -> MediaType.IMAGE_JPEG;
            case Type.GIF -> MediaType.IMAGE_GIF;
            case Type.WEBP -> IMAGE_WEBP;
            case Type.ICO -> IMAGE_X_ICON;
            default -> null;
        };

        if (contentType == null) {
            return notFound();
        }

        return ok(path, contentType);
    }

    @Override
    public IavEntity getById(Long id) {
        if (id == null || id.longValue() <= 0) {
            return null;
        }
        return mapper.selectById(id);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public IavEntity upload(MultipartFile file) throws IOException {
        // 判断上传文件是否有效
        Assert.isTrue(file != null && !file.isEmpty(), "无效的上传文件");

        // 判断文件名称是否有效
        String name = StringUtils.trim(file.getOriginalFilename());
        Assert.isTrue(StringUtils.isNotEmpty(name), "上传文件名称不能为空");

        // 文件类型
        String type = null;
        // 文件后缀名
        String suffix = null;
        int index = name.lastIndexOf(".");
        if (index >= 0 && (index + 1) < name.length()) {
            suffix = StringUtils.trim(name.substring(index + 1).toLowerCase());
            type = Type.suffixOf(suffix);
            name = name.substring(0, index);
        }
        Assert.isTrue(Type.isImage(type), String.format("不支持上传 %s 文件类型，请选择 png、jpg、gif、webp、ico 文件类型上传", file.getOriginalFilename()));

        // 获取上传文件字节数组
        byte[] bytes = file.getBytes();

        IavEntity entity = new IavEntity();
        entity.setName(name);
        entity.setType(type);
        entity.setSize(bytes.length + 0L);
        entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));

        // 获取已删除的id，以复用
        Long deledId = mapper.getDeledId();
        if (deledId != null) {
            entity.setId(deledId);
            entity.setDel(0);
            entity.setUpdTime(0L);
            mapper.updById(entity);
        } else {
            mapper.insert(entity);
        }

        Path path = getPath(entity.getId().toString());

        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, bytes);

        return entity;
    }

}
