package org.xiangqian.note.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.apache.commons.lang3.StringUtils;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
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
@NoArgsConstructor
public class NoteEntity implements Comparable<NoteEntity> {

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

    // 文本内容
    @TableField(exist = false)
    private String content;

    // 包括子节点
    @TableField(exist = false)
    private Boolean contain;

    // 子节点集
    @TableField(exist = false)
    private org.xiangqian.note.util.List<NoteEntity> childList;

    // 上传文件
    @TableField(exist = false)
    private MultipartFile file;

    // 删除标识，0-正常，1-删除
    @TableLogic
    private Integer del;

    // 创建时间（时间戳，单位s）
    private Long addTime;

    // 修改时间（时间戳，单位s）
    private Long updTime;

    public NoteEntity(Path path) throws IOException {
        this.name = StringUtils.trim(path.getFileName().toString());
        this.type = Type.pathOf(path);
        this.size = Files.size(path);
        this.updTime = Files.getLastModifiedTime(path).toMillis() / 1000;
    }

    public void setPids(String pids) {
        if (StringUtils.isNotEmpty(pids)) {
            this.ps = Arrays.stream(pids.split(","))
                    .map(pid -> StringUtils.isNotEmpty(pid) ? Long.parseLong(pid) : null)
                    .filter(Objects::nonNull)
                    .map(pid -> {
                        NoteEntity p = new NoteEntity();
                        p.setId(pid);
                        p.setType(Type.FOLDER);
                        return p;
                    }).collect(Collectors.toList());
        }
    }

    @Override
    public int compareTo(NoteEntity other) {
        if (other == null) {
            return 1;
        }

        if (Type.FOLDER.equals(type)) {
            if (Type.FOLDER.equals(other.type)) {
                return name.toLowerCase().compareTo(other.name.toLowerCase());
            }
            return -1;
        }

        if (Type.FOLDER.equals(other.type)) {
            return 1;
        }

        return name.toLowerCase().compareTo(other.name.toLowerCase());
    }

}
