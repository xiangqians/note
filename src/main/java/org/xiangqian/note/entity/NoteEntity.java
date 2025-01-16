package org.xiangqian.note.entity;

import lombok.Data;
import lombok.NoArgsConstructor;
import org.apache.commons.lang3.StringUtils;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.util.Type;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * 笔记信息
 *
 * @author xiangqian
 * @date 21:35 2024/03/03
 */
@Data
@NoArgsConstructor
public class NoteEntity {

    /**
     * 主键
     */
    private Long id;

    /**
     * 主键集合
     */
    private Set<Long> ids;

    /**
     * 父主键
     */
    private Long pid;

    /**
     * 父节点集合
     */
    private List<NoteEntity> parentList;

    /**
     * 名称
     */
    private String name;

    /**
     * 类型，folder、md、pdf、zip
     */
    private String type;

    /**
     * 大小，单位byte
     */
    private Long size;

    /**
     * 内容
     */
    private String content;

    /**
     * 是否已删除，0-否，1-是
     */
    private Integer delete;

    /**
     * 创建时间戳，单位s
     */
    private Long createTime;

    /**
     * 修改时间戳，单位s
     */
    private Long updateTime;

    /**
     * 当前页
     */
    private Integer current;

    /**
     * 是否包括所有子节点
     */
    private Boolean include;

    /**
     * 上传文件
     */
    private MultipartFile file;

    private List<String> sort;

    private List<NoteEntity> children;
    private Long newId;
    private Long newPid;

    /**
     * ^(\d+)\[(\d+)\].*
     * (笔记主键)[笔记名称长度]笔记名称
     */
    private static final Pattern PATTERN = Pattern.compile("^(\\d+)\\[(\\d+)\\].*");

    public void setParentListStr(String parentListStr) {
        if (StringUtils.isEmpty(parentListStr)) {
            return;
        }

        this.parentList = new ArrayList<>();
        Matcher matcher = PATTERN.matcher(parentListStr);
        while (matcher.find()) {
            NoteEntity parentEntity = new NoteEntity();

            // 获取父笔记主键
            long id = Long.parseLong(matcher.group(1));
            parentEntity.setId(id);

            // 获取父笔记名称长度
            int length = Integer.parseInt(matcher.group(2));
            // 获取父笔记名称
            int start = matcher.end(2) + 1;
            int end = start + length;
            String name = parentListStr.substring(start, end);
            parentEntity.setName(name);

            parentEntity.setType(Type.FOLDER);

            this.parentList.add(parentEntity);

            // 上一级父笔记信息
            parentListStr = parentListStr.substring(end);
            matcher = PATTERN.matcher(parentListStr);
        }

        // 反转
        Collections.reverse(this.parentList);
    }

}
