package org.xiangqian.note.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;
import org.apache.commons.lang3.StringUtils;

import java.util.Arrays;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

/**
 * 笔记信息
 *
 * @author xiangqian
 * @date 21:35 2024/03/03
 */
@Data
@TableName("note")
public class NoteEntity {

    // id
    @TableId(type = IdType.AUTO)
    private Long id;

    // 父id
    private Long pid;

    // 父节点集合
    @TableField(exist = false)
    private List<NoteEntity> ps;

    // 名称
    @TableField("`name`")
    private String name;

    // 类型
    private String type;

    // 大小，单位：byte
    private Long size;

    // 删除标识，0-正常，1-删除
    @TableLogic
    private Integer del;

    // 创建时间（时间戳，单位s）
    private Long addTime;

    // 修改时间（时间戳，单位s）
    private Long updTime;

    // 包括子节点
    @TableField(exist = false)
    private Boolean contain;

    public void setPids(String pids) {
        if (StringUtils.isNotEmpty(pids)) {
            this.ps = Arrays.stream(pids.split(","))
                    .map(pid -> StringUtils.isNotEmpty(pid) ? Long.parseLong(pid) : null)
                    .filter(Objects::nonNull)
                    .map(pid -> {
                        NoteEntity p = new NoteEntity();
                        p.setId(pid);
                        return p;
                    }).collect(Collectors.toList());
        }
    }


    

}
