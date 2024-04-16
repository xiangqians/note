package org.xiangqian.note.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.BooleanUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.mapper.NoteMapper;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.List;
import org.xiangqian.note.util.Type;

import java.time.LocalDateTime;
import java.util.*;
import java.util.stream.Collectors;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
@Service
public class NoteServiceImpl implements NoteService {

    @Autowired
    private NoteMapper mapper;

    @Override
    public Boolean addMdFile(NoteEntity vo) {
        vo.setType(Type.MD);
        return add(vo);
    }

    @Override
    public Boolean addFolder(NoteEntity vo) {
        vo.setType(Type.FOLDER);
        return add(vo);
    }

    @Override
    public List<NoteEntity> list(NoteEntity vo, List list) {
        Long pid = vo.getPid();
        if (list.getOffset().intValue() <= 0 || pid == null || pid.longValue() < 0) {
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
                        .filter(Objects::nonNull)
                        .map(NoteEntity::getId)
                        .collect(Collectors.toSet());
                if (CollectionUtils.isNotEmpty(pids)) {
                    // <pid, name>
                    Map<Long, String> pidMap = mapper.selectList(new LambdaQueryWrapper<NoteEntity>().select(NoteEntity::getId, NoteEntity::getName).in(NoteEntity::getId, pids))
                            .stream().collect(Collectors.toMap(NoteEntity::getId, NoteEntity::getName));
                    for (NoteEntity entity : data) {
                        java.util.List<NoteEntity> ps = entity.getPs();
                        if (CollectionUtils.isNotEmpty(ps)) {
                            for (NoteEntity p : ps) {
                                if (p != null) {
                                    Optional.ofNullable(pidMap.get(p.getId())).ifPresent(p::setName);
                                }
                            }
                        }
                    }
                }
            }
        }
        return result;
    }

    private Boolean add(NoteEntity vo) {
        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        Long pid = vo.getPid();
        Assert.notNull(pid, "pid不能为空");
        if (pid.longValue() != 0) {
            NoteEntity pEntity = mapper.selectById(pid);
            Assert.isTrue(pEntity != null && Type.FOLDER.equals(pEntity.getType()), "pid不存在");
        }

        NoteEntity addEntity = new NoteEntity();
        addEntity.setPid(vo.getPid());
        addEntity.setName(name);
        addEntity.setType(vo.getType());
        addEntity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.insert(addEntity) > 0;
    }

}
