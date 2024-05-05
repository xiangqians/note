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
import java.util.*;
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

        Path path = getPath(id.toString());
        String content = vo.getContent();
        if (content == null) {
            content = "";
        }

        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, content.getBytes(UTF_8));

        long size = Files.size(path);

        entity = new NoteEntity();
        entity.setId(id);
        entity.setSize(size);
        entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        int affectedRows = mapper.updateById(entity);
        if (affectedRows > 0) {
            trigger();
            return true;
        }
        return false;
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
            return notFound();
        }

        NoteEntity entity = getById(id, false);
        if (entity == null || Type.FOLDER.equals(entity.getType())) {
            return notFound();
        }

        Path path = getPath(id.toString());
        if (!Files.exists(path)) {
            return notFound();
        }

        // 读取文件
        byte[] data = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(data);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.add(HttpHeaders.CONTENT_DISPOSITION, String.format("attachment; filename=%s", URLEncoder.encode(String.format("%s.%s", entity.getName(), entity.getType()), UTF_8)));
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

        // 注：
        // 逻辑删除，并不删除物理文件，在id没有被复用前，都可恢复

        int affectedRows = mapper.deleteById(id);
        if (affectedRows > 0) {
            trigger();
            return true;
        }

        return false;
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

        int affectedRows = mapper.updateById(entity);
        if (affectedRows > 0) {
            trigger();
            return true;
        }

        return false;
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean reUpload(NoteEntity vo) throws IOException {
        Long id = vo.getId();
        NoteEntity entity = getById(id, false);
        Assert.isTrue(entity != null && !StringUtils.equalsAny(entity.getType(), Type.FOLDER, Type.MD), "id不能是文件夹类型或者Markdown文件类型");
        return uploadOrReUpload(vo);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean upload(NoteEntity vo) throws IOException {
        vo.setId(null);
        return uploadOrReUpload(vo);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean addMd(NoteEntity vo) throws IOException {
        vo.setType(Type.MD);
        return add(vo);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean addFolder(NoteEntity vo) throws IOException {
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

        // 判断上传文件是否有效
        Assert.isTrue(file != null && !file.isEmpty(), "无效的上传文件");

        // 文件名称
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

        int affectedRows = 0;

        NoteEntity entity = new NoteEntity();
        entity.setPid(pid);
        entity.setName(name);
        entity.setType(type);
        entity.setSize(bytes.length + 0L);

        Long id = vo.getId();
        if (id != null) {
            entity.setId(id);
            entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
            affectedRows = mapper.updateById(entity);
        } else {
            entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));

            // 获取已删除的id，以复用
            Long deledId = mapper.getDeledId();
            if (deledId != null) {
                entity.setId(deledId);
                entity.setDel(0);
                entity.setUpdTime(0L);
                affectedRows = mapper.updById(entity);
            } else {
                affectedRows = mapper.insert(entity);
            }
        }
        id = entity.getId();

        Path path = getPath(id.toString());
        // 将内容写入文件（覆盖），如果文件不存在则创建
        Files.write(path, bytes);

        trigger();

        if (affectedRows > 0) {
            delTmpPath(id);
            return true;
        }

        return false;
    }

    private Boolean add(NoteEntity vo) throws IOException {
        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        Long pid = vo.getPid();
        verifyPid(pid);

        NoteEntity entity = new NoteEntity();
        entity.setPid(pid);
        entity.setName(name);
        entity.setType(vo.getType());
        entity.setSize(0L);
        entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));

        int affectedRows = 0;

        // 获取已删除的id，以复用
        Long deledId = mapper.getDeledId();
        if (deledId != null) {
            entity.setId(deledId);
            entity.setDel(0);
            entity.setUpdTime(0L);
            affectedRows = mapper.updById(entity);
        } else {
            affectedRows = mapper.insert(entity);
        }

        if (affectedRows > 0) {
            delPath(entity.getId());
            delTmpPath(entity.getId());
            return true;
        }

        return false;
    }

    private void delTmpPath(Long id) throws IOException {
        String name = id.toString();
        Path tmpPath = getTmpPath(name);
        PathUtils.delete(tmpPath);

        Path pTmpPath = tmpPath.getParent();
        if (Files.exists(pTmpPath)) {
            // 获取目录下的子文件夹
            Iterator<Path> iterator = Files.list(pTmpPath).iterator();
            while (iterator.hasNext()) {
                tmpPath = iterator.next();
                if (tmpPath.getFileName().toString().startsWith(name + "_")) {
                    PathUtils.delete(tmpPath);
                }
            }
        }
    }

    private void delPath(Long id) throws IOException {
        String name = id.toString();
        Path path = getPath(name);
        PathUtils.delete(path);
    }

    /**
     * 校验pid
     *
     * @param pid
     */
    private NoteEntity verifyPid(Long pid) {
        NoteEntity entity = getById(pid, false);
        Assert.isTrue(entity != null && Type.FOLDER.equals(entity.getType()), "pid不存在");
        return entity;
    }

    /**
     * 校验id
     *
     * @param id
     */
    private NoteEntity verifyId(Long id) {
        NoteEntity entity = getById(id, false);
        Assert.isTrue(entity != null, "id不存在");
        return entity;
    }

    /**
     * 触发【计算目录大小任务】
     */
    private void trigger() {
        synchronized (lock) {
            lock.notify();
        }
    }

}
