package org.xiangqian.note.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.BooleanUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.configuration.mybatis.LazyList;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.mapper.NoteMapper;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.FileUtil;
import org.xiangqian.note.util.TimeUtil;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
@Slf4j
@Service
public class NoteServiceImpl extends AbsService implements NoteService, ApplicationRunner {

    @Autowired
    private NoteMapper mapper;

    public NoteServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean createFolder(NoteEntity entity) throws IOException {
        entity.setType(Type.FOLDER);
        entity.setSize(0L);
        if (create(entity)) {
            deleteFile(entity.getId());
            return true;
        }
        return false;
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean createMd(NoteEntity entity) throws IOException {
        entity.setType(Type.MD);
        entity.setSize(0L);
        if (create(entity)) {
            deleteFile(entity.getId());
            return true;
        }
        return false;
    }

    @Override
    public Boolean upload(NoteEntity entity) throws IOException {
        MultipartFile file = entity.getFile();

        // 判断上传文件是否有效
        Assert.isTrue(file != null && !file.isEmpty(), "无效的文件");

        // 判断文件名称是否有效
        String name = StringUtils.trim(file.getOriginalFilename());
        Assert.isTrue(StringUtils.isNotEmpty(name), "文件名称不能为空");

        // 文件后缀名
        String suffix = null;
        int index = name.lastIndexOf(".");
        if (index >= 0 && (index + 1) < name.length()) {
            suffix = StringUtils.trim(name.substring(index + 1).toLowerCase());
            name = name.substring(0, index);
        }
        Assert.isTrue(StringUtils.equalsAny(suffix, Type.PDF, Type.ZIP), String.format("不支持上传 “%s” 文件类型，请选择 pdf、zip 文件类型上传", file.getOriginalFilename()));

        Long id = entity.getId();
        if (id != null) {
            NoteEntity dbEntity = mapper.getById(id);
            Assert.isTrue(dbEntity != null && StringUtils.equalsAny(dbEntity.getType(), Type.PDF, Type.ZIP), "id不存在");
        }

        entity.setName(name);
        entity.setType(suffix);
        entity.setSize(file.getSize());

        Boolean result = false;
        if (id != null) {
            entity.setUpdateTime(TimeUtil.now());
            result = mapper.updateById(entity);
        } else {
            result = create(entity);
        }

        if (result) {
            // 写入文件
            writeFile(entity.getId(), file.getBytes());

            // 计算目录大小
            calculateFolderSize(entity.getId());

            return true;
        }

        return false;
    }

    @Override
    public Boolean deleteById(Long id) {
        if (!isValidId(id)) {
            return false;
        }

        NoteEntity entity = mapper.getById(id);
        if (entity == null) {
            return false;
        }

        if (Type.FOLDER.equals(entity.getType())) {
            Assert.isTrue(mapper.countChildrenByPid(id).intValue() == 0, "无法删除非空文件夹");
        }

        // 注：逻辑删除，并不删除物理文件，在id没有被复用前，都可恢复

        if (mapper.deleteById(id)) {
            // 计算目录大小
            calculateFolderSize(entity.getPid());

            return true;
        }

        return false;
    }

    @Override
    public Boolean rename(NoteEntity entity) {
        if (!isValidId(entity.getId())) {
            return false;
        }

        // 校验名称
        String name = StringUtils.trim(entity.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");
        entity.setName(name);

        NoteEntity updateEntity = new NoteEntity();
        updateEntity.setId(entity.getId());
        updateEntity.setName(name);
        updateEntity.setUpdateTime(TimeUtil.now());
        return mapper.updateById(updateEntity);
    }

    @Override
    public Boolean paste(NoteEntity entity) {
        Long id = entity.getId();
        if (!isValidId(id)) {
            return false;
        }

        NoteEntity storedEntity = mapper.getById(id);
        if (storedEntity == null) {
            return false;
        }

        Long oldPid = storedEntity.getPid();

        // 校验父节点主键
        Long pid = entity.getPid();
        validatePid(pid);

        if (id.equals(pid)) {
            return false;
        }

        entity = mapper.getByPidAndId(id, pid);
        if (entity != null) {
            if (pid.equals(entity.getPid())) {
                return false;
            }
            Assert.isTrue(!Type.FOLDER.equals(entity.getType()), "目标文件夹是源文件夹的子文件夹");
        }

        NoteEntity updateEntity = new NoteEntity();
        updateEntity.setId(id);
        updateEntity.setPid(pid);
        updateEntity.setUpdateTime(TimeUtil.now());
        if (mapper.updateById(updateEntity)) {
            // 计算目录大小
            calculateFolderSize(oldPid);
            calculateFolderSize(id);

            return true;
        }

        return false;
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean updateContentById(Long id, String content) throws IOException {
        if (!isValidId(id)) {
            return false;
        }

        NoteEntity entity = mapper.getById(id);
        if (entity == null || !Type.MD.equals(entity.getType())) {
            return false;
        }

        byte[] bytes = null;
        if (StringUtils.isNotEmpty(content)) {
            bytes = content.getBytes(StandardCharsets.UTF_8);
        } else {
            bytes = new byte[0];
        }

        NoteEntity updateEntity = new NoteEntity();
        updateEntity.setId(id);
        updateEntity.setSize(Long.valueOf(bytes.length));
        updateEntity.setUpdateTime(TimeUtil.now());
        mapper.updateById(updateEntity);

        writeFile(id, bytes);

        // 计算目录大小
        calculateFolderSize(entity.getId());

        return true;
    }

    @Override
    public Boolean updateSortById(Long id, String content) throws IOException {
        NoteEntity entity = getById(id);
        if (entity == null || !Type.FOLDER.equals(entity.getType())) {
            return false;
        }

        byte[] bytes = null;
        if (StringUtils.isNotEmpty(content)) {
            bytes = content.getBytes(StandardCharsets.UTF_8);
        } else {
            bytes = new byte[0];
        }

        writeFile(id, bytes);
        return true;
    }

    @Override
    public String getContentById(Long id) throws IOException {
        byte[] bytes = readFile(id);
        if (bytes == null) {
            return null;
        }

        return new String(bytes, StandardCharsets.UTF_8);
    }

    @Override
    public NoteEntity getById(Long id) {
        if (id == null || id.longValue() < 0) {
            return null;
        }

        if (id.longValue() == 0) {
            NoteEntity entity = new NoteEntity();
            entity.setId(0L);
            entity.setName("根目录");
            entity.setType(Type.FOLDER);
            return entity;
        }

        return mapper.getById(id);
    }

    @Override
    public List<NoteEntity> getParentListById(Long id) {
        List<NoteEntity> list = mapper.getParentListById(id);
        list.sort((o1, o2) -> {
            if (o1.getId().equals(o2.getPid())) {
                return -1;
            }

            if (o1.getPid().equals(o2.getId())) {
                return 1;
            }

            return 0;
        });
        return list;
    }

    @Override
    public LazyList<NoteEntity> getChildList(NoteEntity entity) throws IOException {
        LazyList lazyList = new LazyList(entity.getCurrent());

        String content = StringUtils.trim(entity.getContent());
        if (StringUtils.isNotEmpty(content)) {
            List<Long> childIdList = mapper.getChildIdListByPid(entity.getPid());
            if (CollectionUtils.isEmpty(childIdList)) {
                return lazyList;
            }

            String content1 = content.toLowerCase();

            Set<Long> ids = new HashSet<>();
            for (Long childId : childIdList) {
                Path path = getPath(childId);
                if (Files.exists(path)) {
                    boolean isMatch = false;
                    try {
                        isMatch = Files.lines(path, StandardCharsets.UTF_8)
                                .parallel()  // 启用并行流来处理大文件或多个行，从而加快查找过程
                                .anyMatch(line -> line.toLowerCase().contains(content1));
                    } catch (UncheckedIOException e) {
                        log.error(path.toString(), e);
                    }
                    if (isMatch) {
                        ids.add(childId);
                    }
                }
            }
            if (CollectionUtils.isEmpty(ids)) {
                return lazyList;
            }
            entity.setIds(ids);
        }

        if (!BooleanUtils.isTrue(entity.getInclude())) {
            String sort = Optional.ofNullable(readFile(entity.getPid())).map(bytes -> new String(bytes, StandardCharsets.UTF_8)).orElse(null);
            if (StringUtils.isNotEmpty(sort)) {
                // (?<!\\)，
                // 使用正则表达式以中文逗号分割，但忽略被转义的中文逗号
                entity.setSort(Arrays.stream(sort.split("(?<!\\\\)，")).toList());
            }
        }

        return mapper.getChildList(lazyList, entity);
    }

    @Override
    public Resource getResourceById(Long id, String name) throws IOException {
        if (!isValidId(id)) {
            return null;
        }

        NoteEntity entity = mapper.getById(id);
        if (entity == null || Type.FOLDER.equals(entity.getType())) {
            return null;
        }

        if (Type.ZIP.equals(entity.getType()) && StringUtils.isNotEmpty(name)) {
            String name1 = null;
            int index = name.lastIndexOf("/");
            if (index >= 0 && (index + 1) < name.length()) {
                name1 = name.substring(index + 1);
            } else {
                name1 = name;
            }
            return createResource(name1, FileUtil.getType(name), FileUtil.getZipEntryContent(getPath(id), name));
        }

        return createResource(String.format("%s.%s", entity.getName(), entity.getType()), entity.getType(), readFile(id));
    }

    private Boolean create(NoteEntity entity) {
        // 校验名称
        String name = StringUtils.trim(entity.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");
        entity.setName(name);

        // 校验父节点主键
        validatePid(entity.getPid());

        entity.setId(null);
        entity.setDelete(0);
        entity.setCreateTime(TimeUtil.now());
        entity.setUpdateTime(0L);

        // 获取已删除的id，以复用
        Long deletedId = mapper.getDeletedId();
        if (deletedId != null) {
            entity.setId(deletedId);
            return mapper.updateById(entity);
        }

        return mapper.create(entity);
    }

    private void validatePid(Long pid) {
        NoteEntity parentEntity = getById(pid);
        Assert.isTrue(parentEntity != null && Type.FOLDER.equals(parentEntity.getType()), "上级目录不存在");
    }

    /**
     * 计算目录大小
     *
     * @param id
     */
    private void calculateFolderSize(Long id) {
        List<NoteEntity> parentList = mapper.getParentListById(id);
        for (NoteEntity parent : parentList) {
            if (Type.FOLDER.equals(parent.getType())) {
                NoteEntity updateEntity = new NoteEntity();
                updateEntity.setId(parent.getId());
                updateEntity.setSize(mapper.getSizeById(parent.getId()));
                mapper.updateById(updateEntity);
            }
        }
    }

    @Override
    public void run(ApplicationArguments args) throws Exception {
        new Thread(() -> {
            long startTime = System.currentTimeMillis();
            log.debug("【计算全部目录大小】...");

            List<NoteEntity> list = mapper.getFolderList();
            if (CollectionUtils.isNotEmpty(list)) {
                for (NoteEntity entity : list) {
                    Long id = entity.getId();
                    calculateFolderSize(id);
                }
            }

            long endTime = System.currentTimeMillis();
            log.debug("【计算全部目录大小】耗时 {} ms", endTime - startTime);
        }).start();
    }

}
