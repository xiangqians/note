package org.xiangqian.note.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.io.file.PathUtils;
import org.apache.commons.lang3.BooleanUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.core.env.Environment;
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
import org.xiangqian.note.service.GetNoteService;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.LocalDateTime;
import java.util.Collection;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
@Slf4j
@Service
public class NoteServiceImpl extends AbsService implements NoteService, ApplicationRunner, Runnable {

    @Autowired
    private NoteMapper mapper;

    @Autowired
    private List<GetNoteService> getNoteServices;

    private final Object lock = new Object();

    protected NoteServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Override
    public Boolean updContentById(NoteEntity vo) throws IOException {
        Long id = vo.getId();
        NoteEntity entity = getById(id, false);
        if (entity == null) {
            return false;
        }

        Path path = getDataPath(id.toString());
        String content = vo.getContent();
        if (content == null) {
            content = "";
        }

        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, content.getBytes(UTF_8));

        long size = Files.size(path);

        NoteEntity updEntity = new NoteEntity();
        updEntity.setId(id);
        updEntity.setSize(size);
        updEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        boolean result = mapper.updateById(updEntity) > 0;
        if (result) {
            synchronized (lock) {
                lock.notify();
            }
        }
        return result;
    }

    @Override
    public org.xiangqian.note.util.List<NoteEntity> list(NoteEntity vo, Integer offset) {
        if (offset == null) {
            offset = 1;
        }
        Long pid = vo.getPid();
        org.xiangqian.note.util.List list = new org.xiangqian.note.util.List(offset);
        if (offset.intValue() <= 0 || pid == null || pid.longValue() < 0) {
            return list;
        }

        vo.setName(StringUtils.trimToNull(vo.getName()));
        vo.setType(StringUtils.trimToNull(vo.getType()));
        list = mapper.list(vo, list);
        if (BooleanUtils.isTrue(vo.getContain())) {
            if (CollectionUtils.isNotEmpty(list)) {
                Set<Long> pids = ((List<NoteEntity>) list).stream().map(NoteEntity::getPs)
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
                    for (NoteEntity entity : (List<NoteEntity>) list) {
                        List<NoteEntity> ps = entity.getPs();
                        if (CollectionUtils.isNotEmpty(ps)) {
                            for (NoteEntity p : ps) {
                                p.setName(pidMap.get(p.getId()));
                            }
                        }
                    }
                }
            }
        }
        return list;
    }

    @Override
    public ModelAndView getViewById(ModelAndView modelAndView, Long id, List<String> names) throws Exception {
        NoteEntity entity = getById(id, false);
        if (entity == null) {
            return AbsController.errorView(modelAndView);
        }

        String type = entity.getType();
        for (GetNoteService getNoteService : getNoteServices) {
            if (getNoteService.isSupported(type)) {
                return getNoteService.getView(modelAndView, entity, names);
            }
        }

        return AbsController.errorView(modelAndView);
    }

    @Override
    public ResponseEntity<Resource> getStreamById(Long id, List<String> names) throws Exception {
        NoteEntity entity = getById(id, false);
        if (entity == null) {
            return notFound();
        }

        String type = entity.getType();
        for (GetNoteService getNoteService : getNoteServices) {
            if (getNoteService.isSupported(type)) {
                return getNoteService.getStream(entity, names);
            }
        }

        return notFound();
    }

    @Override
    public NoteEntity getById(Long id, boolean isGetPs) {
        if (id == null || id.longValue() < 0) {
            return null;
        }

        if (id.longValue() == 0) {
            NoteEntity entity = new NoteEntity();
            entity.setId(0L);
            entity.setType(Type.FOLDER);
            return entity;
        }

        if (!isGetPs) {
            return mapper.selectById(id);
        }

        NoteEntity entity = mapper.getById(id);
        if (entity != null) {
            List<NoteEntity> ps = entity.getPs();
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

        Path path = getDataPath(id.toString());
        if (!Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] data = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(data);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.add(HttpHeaders.CONTENT_DISPOSITION, String.format("attachment; filename=%s", URLEncoder.encode(entity.getName() + "." + entity.getType(), UTF_8)));
        headers.add(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_OCTET_STREAM_VALUE);
        headers.add(HttpHeaders.CONTENT_LENGTH, String.valueOf(resource.contentLength()));

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    @Override
    public Boolean delById(Long id) throws IOException {
        NoteEntity entity = verifyId(id);
        if (Type.FOLDER.equals(entity.getType())) {
            NoteEntity child = mapper.selectOne(new LambdaQueryWrapper<NoteEntity>()
                    .select(NoteEntity::getId)
                    .eq(NoteEntity::getPid, id)
                    .last("LIMIT 1"));
            Assert.isNull(child, "无法删除非空文件夹");
        }

        boolean result = mapper.deleteById(id) > 0;
        if (result && !Type.FOLDER.equals(entity.getType())) {
            PathUtils.deleteFile(getDataPath(id.toString()));
        }
        if (result) {
            synchronized (lock) {
                lock.notify();
            }
        }
        return result;
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
        boolean result = mapper.updateById(entity) > 0;
        if (result) {
            synchronized (lock) {
                lock.notify();
            }
        }
        return result;
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

    @Override
    public void run(ApplicationArguments args) throws Exception {
        new Thread(this).start();
    }

    @Override
    public void run() {
        while (true) {
            try {
                long startTime = System.currentTimeMillis();
                log.debug("开始执行【计算目录大小任务】...");

                List<NoteEntity> entities = mapper.selectList(new LambdaQueryWrapper<NoteEntity>()
                        .select(NoteEntity::getId)
                        .eq(NoteEntity::getType, Type.FOLDER));
                if (CollectionUtils.isEmpty(entities)) {
                    return;
                }

                for (NoteEntity entity : entities) {
                    Long id = entity.getId();
                    Long size = mapper.getSizeById(id);

                    NoteEntity updEntity = new NoteEntity();
                    updEntity.setId(id);
                    updEntity.setSize(size);
                    mapper.updateById(updEntity);
                }

                long endTime = System.currentTimeMillis();
                log.debug("执行【计算目录大小任务】完成，耗时：{}ms", endTime - startTime);

                synchronized (lock) {
                    log.debug("【计算目录大小任务】沉睡");
                    lock.wait();
                    log.debug("【计算目录大小任务】醒来");
                }
            } catch (Exception e) {
                log.error("执行【计算目录大小任务】发生异常", e);
            }
        }
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
            type = Type.suffixOf(suffix);
            name = name.substring(0, index);
        }
        Assert.isTrue(StringUtils.equalsAny(type,
                        Type.DOC, Type.DOCX,
                        Type.XLS, Type.XLSX,
                        Type.PDF,
                        Type.HTML,
                        Type.ZIP),
                String.format("不支持上传 %s 文件类型，请选择 doc、docx、xls、xlsx、pdf、html、zip 文件类型上传", file.getOriginalFilename()));

        Long pid = vo.getPid();
        verifyPid(pid);

        byte[] bytes = file.getBytes();

        NoteEntity entity = new NoteEntity();
        entity.setPid(pid);
        entity.setName(name);
        entity.setType(type);
        entity.setSize(bytes.length + 0L);
        if (vo.getId() == null) {
            entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
            mapper.insert(entity);
        } else {
            entity.setId(vo.getId());
            entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
            mapper.updateById(entity);
        }

        Path path = getDataPath(entity.getId().toString());
        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, bytes);

        synchronized (lock) {
            lock.notify();
        }
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
