package org.xiangqian.note.service.impl;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.IavEntity;
import org.xiangqian.note.mapper.IavMapper;
import org.xiangqian.note.service.IavService;
import org.xiangqian.note.util.DateUtil;

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

    @Override
    public ResponseEntity<Resource> getStreamById(Long id) throws IOException {
        if (id == null) {
            return ResponseEntity.notFound().build();
        }

        // 查库
        IavEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        if (entity == null) {
            return ResponseEntity.notFound().build();
        }

        // 文件路径
        Path path = getPath(id);

        // 判断文件是否存在
        if (!Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] bytes = Files.readAllBytes(path);

        // 将文件数据转换为字节数组资源
        ByteArrayResource resource = new ByteArrayResource(bytes);

        // 构建响应头
        HttpHeaders headers = new HttpHeaders();
        headers.setContentLength(resource.contentLength());
        headers.set(HttpHeaders.CONTENT_TYPE, entity.getType());

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public IavEntity upload(MultipartFile file) throws IOException {
        // 判断上传文件是否有效
        if (file == null || file.isEmpty()) {
            return null;
        }

        // 获取上传文件字节数组
        byte[] bytes = file.getBytes();

        IavEntity entity = new IavEntity();
        entity.setName(StringUtils.trim(file.getOriginalFilename()));
        entity.setType(file.getContentType());
        entity.setSize(bytes.length + 0L);
        entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
        mapper.insert(entity);

        Path path = getPath(entity.getId(), true);
        Files.write(path, bytes);

        return entity;
    }

}
