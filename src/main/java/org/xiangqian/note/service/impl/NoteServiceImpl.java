package org.xiangqian.note.service.impl;

import com.aspose.words.Document;
import com.aspose.words.SaveFormat;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.BooleanUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.mapper.NoteMapper;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.List;
import org.xiangqian.note.util.Type;

import java.io.*;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.attribute.FileTime;
import java.time.LocalDateTime;
import java.util.*;
import java.util.stream.Collectors;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
@Slf4j
@Service
public class NoteServiceImpl extends AbsService implements NoteService {

    @Autowired
    private NoteMapper mapper;

    @Override
    public List<NoteEntity> list(NoteEntity vo, List list) {
        Integer offset = list.getOffset();
        Long pid = vo.getPid();
        if (offset == null || offset.intValue() <= 0 || pid == null || pid.longValue() < 0) {
            return list;
        }

        vo.setName(StringUtils.trimToNull(vo.getName()));
        vo.setType(StringUtils.trimToNull(vo.getType()));
        List<NoteEntity> result = mapper.list(vo, list);
        if (BooleanUtils.isTrue(vo.getContain())) {
            java.util.List<NoteEntity> data = result.getData();
            if (CollectionUtils.isNotEmpty(data)) {
                Set<Long> pids = data.stream().map(NoteEntity::getPs)
                        .filter(CollectionUtils::isNotEmpty)
                        .flatMap(Collection::stream)
                        .map(NoteEntity::getId)
                        .collect(Collectors.toSet());
                if (CollectionUtils.isNotEmpty(pids)) {
                    // <pid, name>
                    Map<Long, String> pidMap = mapper.selectList(new LambdaQueryWrapper<NoteEntity>()
                                    .select(NoteEntity::getId, NoteEntity::getName)
                                    .in(NoteEntity::getId, pids))
                            .stream().collect(Collectors.toMap(NoteEntity::getId, NoteEntity::getName));
                    for (NoteEntity entity : data) {
                        java.util.List<NoteEntity> ps = entity.getPs();
                        if (CollectionUtils.isNotEmpty(ps)) {
                            for (NoteEntity p : ps) {
                                p.setName(pidMap.get(p.getId()));
                            }
                        }
                    }
                }
            }
        }
        return result;
    }

    @Override
    public NoteEntity getById(Long id) {
        Assert.notNull(id, "id不能为空");
        Assert.isTrue(id.longValue() >= 0, "id不能小于0");
        if (id.longValue() == 0) {
            NoteEntity entity = new NoteEntity();
            entity.setId(0L);
            entity.setType(Type.FOLDER);
            return entity;
        }

        NoteEntity entity = mapper.getById(id);
        if (entity != null) {
            java.util.List<NoteEntity> ps = entity.getPs();
            if (CollectionUtils.isNotEmpty(ps)) {
                Set<Long> pids = ps.stream().map(NoteEntity::getId).collect(Collectors.toSet());
                // <pid, name>
                Map<Long, String> pidMap = mapper.selectList(new LambdaQueryWrapper<NoteEntity>()
                                .select(NoteEntity::getId, NoteEntity::getName)
                                .in(NoteEntity::getId, pids))
                        .stream().collect(Collectors.toMap(NoteEntity::getId, NoteEntity::getName));
                for (NoteEntity p : ps) {
                    p.setName(pidMap.get(p.getId()));
                }
            }
        }
        return entity;
    }

    @Override
    public Object getViewById(ModelAndView modelAndView, Long id, java.util.List<String> names) {
        try {
            NoteEntity entity = verifyId(id);

            String type = entity.getType();
            String model = switch (type) {
                case Type.MD -> "md";
                case Type.DOC, Type.DOCX, Type.PDF -> "pdf";
                case Type.HTML -> "html";
                case Type.ZIP -> "zip";
                default -> null;
            };
            if (model == null) {
                return AbsController.errorView(modelAndView);
            }

            if (StringUtils.equalsAny(type, Type.MD, Type.HTML)) {
                String content = null;
                Path path = getPath(id);
                if (Files.exists(path)) {
                    content = Files.readString(path, StandardCharsets.UTF_8);
                }
                entity.setContent(content);

            } else if (Type.ZIP.equals(type)) {
                Path path = getPath(id);
                if (Files.exists(path)) {
                    Path tmpPath = getTmpPath(id);
                    if (!Files.exists(tmpPath)) {
                        Files.createDirectories(tmpPath);

                        ZipFile zip = new ZipFile(path.toFile());
                        Enumeration<? extends ZipEntry> entries = zip.entries();
                        Stack<Path1> stack = new Stack<>();
                        while (entries.hasMoreElements()) {
                            ZipEntry entry = entries.nextElement();
                            Path entryPath = Paths.get(tmpPath.toString(), entry.getName());
                            if (entry.isDirectory()) {
                                Files.createDirectories(entryPath);
                                stack.push(new Path1(entryPath, entry.getLastModifiedTime()));
                            } else {
                                OutputStream outputStream = Files.newOutputStream(entryPath);
                                zip.getInputStream(entry).transferTo(outputStream);
                                Files.setLastModifiedTime(entryPath, entry.getLastModifiedTime());
                            }
                        }

                        while (!stack.empty()) {
                            Path1 path1 = stack.pop();
                            Files.setLastModifiedTime(path1.getPath(), path1.getLastModified());
                        }
                    }

                    java.util.List<NoteEntity> ps = null;
                    if (CollectionUtils.isNotEmpty(names)) {
                        for (String name : names) {
                            if (tmpPath == null || !Files.isDirectory(tmpPath)) {
                                break;
                            }

                            // 获取目录下的子文件夹
                            Iterator<Path> iterator = Files.list(tmpPath).iterator();
                            while (iterator.hasNext()) {
                                tmpPath = iterator.next();
                                if (tmpPath.getFileName().toString().equals(name)) {
                                    if (ps == null) {
                                        ps = new ArrayList<>(names.size());
                                    }
                                    ps.add(new NoteEntity(tmpPath));
                                    break;
                                }
                                tmpPath = null;
                            }
                        }
                    }

                    if (CollectionUtils.isNotEmpty(ps)) {
                        entity.setPs(ps);
                    }

                    if (tmpPath != null) {
                        if (Files.isDirectory(tmpPath)) {
                            java.util.List<NoteEntity> children = new ArrayList<>(16);

                            // 获取目录下的子文件夹
                            Iterator<Path> iterator = Files.list(tmpPath).iterator();
                            while (iterator.hasNext()) {
                                tmpPath = iterator.next();
                                children.add(new NoteEntity(tmpPath));
                            }

                            Collections.sort(children);
                            entity.setChildren(children);
                        } else {
                            entity = new NoteEntity(tmpPath);

                            // 读取文件内容
                            String content = Files.readString(tmpPath, StandardCharsets.UTF_8);
                            entity.setContent(content);

                            model = "text";
                        }
                    }
                }
            }

            AbsController.setVoAttribute(modelAndView, entity);

            modelAndView.setViewName(String.format("note/%s/view", model));

            return modelAndView;
        } catch (Exception e) {
            log.error("", e);
            return AbsController.errorView(modelAndView);
        }
    }

    @Data
    @AllArgsConstructor
    static class Path1 {
        private Path path;
        private FileTime lastModified;
    }

    @Override
    public ResponseEntity<Resource> getStreamById(Long id) throws Exception {
        if (id == null) {
            return ResponseEntity.notFound().build();
        }

        // 查库
        NoteEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        if (entity == null) {
            return ResponseEntity.notFound().build();
        }

        Path path = null;
        MediaType contentType = null;
        switch (entity.getType()) {
            case Type.DOC, Type.DOCX -> {
                InputStream inputStream = null;
                OutputStream outputStream = null;
                try {
                    inputStream = new FileInputStream(getPath(id).toFile());
                    Document document = new Document(inputStream);

                    // 支持DOC, DOCX, OOXML, RTF HTML, OpenDocument, PDF, EPUB, XPS, SWF 相互转换
                    path = getTmpPath(id, true);
                    outputStream = new FileOutputStream(path.toFile());
                    document.save(outputStream, SaveFormat.PDF);
                } finally {
                    IOUtils.closeQuietly(outputStream, inputStream);
                }
                contentType = MediaType.APPLICATION_PDF;
            }
            case Type.PDF -> {
                path = getPath(id);
                contentType = MediaType.APPLICATION_PDF;

            }
            case Type.HTML -> {
                path = getPath(id);
                contentType = MediaType.TEXT_HTML;
            }
        }
        if (path == null || !Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] bytes = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(bytes);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.setContentLength(resource.contentLength());
        headers.setContentType(contentType);

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<Resource> download(Long id) throws IOException {
        if (id == null) {
            return ResponseEntity.notFound().build();
        }

        // 查库
        NoteEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        if (entity == null || Type.FOLDER.equals(entity.getType())) {
            return ResponseEntity.notFound().build();
        }

        Path path = getPath(id);
        if (!Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] data = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(data);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.add(HttpHeaders.CONTENT_DISPOSITION, String.format("attachment; filename=%s", URLEncoder.encode(entity.getName() + "." + entity.getType(), StandardCharsets.UTF_8)));
        headers.add(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_OCTET_STREAM_VALUE);
        headers.add(HttpHeaders.CONTENT_LENGTH, String.valueOf(resource.contentLength()));

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    @Override
    public Boolean delById(Long id) {
        NoteEntity entity = verifyId(id);
        if (Type.FOLDER.equals(entity)) {
            NoteEntity child = mapper.selectOne(new LambdaQueryWrapper<NoteEntity>()
                    .select(NoteEntity::getId)
                    .eq(NoteEntity::getPid, id)
                    .last("LIMIT 1"));
            Assert.isNull(child, "无法删除非空文件夹");
        }
        return mapper.deleteById(id) > 0;
    }

    @Override
    public Boolean rename(NoteEntity vo) {
        Long id = vo.getId();
        verifyId(id);

        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        NoteEntity entity = new NoteEntity();
        entity.setId(id);
        entity.setName(name);
        entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.updateById(entity) > 0;
    }

    @Override
    public Boolean paste(NoteEntity vo) {
        Long id = vo.getId();
        verifyId(id);

        Long pid = vo.getPid();
        verifyPid(pid);

        NoteEntity entity = new NoteEntity();
        entity.setId(id);
        entity.setPid(pid);
        entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.updateById(entity) > 0;
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean reUpload(NoteEntity vo) throws IOException {
        NoteEntity entity = verifyId(vo.getId());
        Assert.isTrue(!Type.FOLDER.equals(entity.getType()) && !Type.MD.equals(entity.getType()), "id不能是文件夹类型或者Markdown文件类型");
        return uploadOrReUpload(vo);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean upload(NoteEntity vo) throws IOException {
        vo.setId(null);
        return uploadOrReUpload(vo);
    }

    @Override
    public Boolean addMd(NoteEntity vo) {
        vo.setType(Type.MD);
        return add(vo);
    }

    @Override
    public Boolean addFolder(NoteEntity vo) {
        vo.setType(Type.FOLDER);
        return add(vo);
    }

    private Boolean uploadOrReUpload(NoteEntity vo) throws IOException {
        // 文件
        MultipartFile file = vo.getFile();
        Assert.isTrue(!file.isEmpty(), "上传文件不能为空");

        // 文件名称
        String name = StringUtils.trim(file.getOriginalFilename());
        Assert.isTrue(StringUtils.isNotEmpty(name), "上传文件名称不能为空");
        // 文件后缀名
        String suffix = null;
        // 文件类型
        String type = null;
        int index = name.lastIndexOf(".");
        if (index >= 0 && (index + 1) < name.length()) {
            suffix = StringUtils.trim(name.substring(index + 1).toLowerCase());
            name = name.substring(0, index);
            type = switch (suffix) {
                case Type.DOC -> Type.DOC;
                case Type.DOCX -> Type.DOCX;
                case Type.PDF -> Type.PDF;
                case Type.HTML -> Type.HTML;
                case Type.ZIP -> Type.ZIP;
                default -> null;
            };
        }
        Assert.notNull(type, String.format("不支持上传 %s 文件类型，请选择 doc、docx、pdf、html、zip 文件类型上传", suffix));

        Long pid = vo.getPid();
        verifyPid(pid);

        byte[] bytes = file.getBytes();

        NoteEntity entity = new NoteEntity();
        entity.setPid(pid);
        entity.setName(name);
        entity.setType(type);
        entity.setSize(bytes.length + 0L);
        entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
        if (vo.getId() == null) {
            mapper.insert(entity);
        } else {
            entity.setId(vo.getId());
            mapper.updateById(entity);
        }

        Path path = getPath(entity.getId(), true);
        Files.write(path, bytes);
        return true;
    }

    private Boolean add(NoteEntity vo) {
        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        Long pid = vo.getPid();
        verifyPid(pid);

        NoteEntity addEntity = new NoteEntity();
        addEntity.setPid(pid);
        addEntity.setName(name);
        addEntity.setType(vo.getType());
        addEntity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.insert(addEntity) > 0;
    }

    private NoteEntity verifyId(Long id) {
        Assert.notNull(id, "id不能为空");

        NoteEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        Assert.notNull(entity, "id不存在");

        return entity;
    }

    private void verifyPid(Long pid) {
        Assert.notNull(pid, "pid不能为空");

        // 根节点
        if (pid.longValue() == 0) {
            return;
        }

        NoteEntity entity = null;
        if (pid.longValue() > 0) {
            entity = mapper.selectById(pid);
        }
        Assert.isTrue(entity != null && Type.FOLDER.equals(entity.getType()), "pid不存在");
    }

}
