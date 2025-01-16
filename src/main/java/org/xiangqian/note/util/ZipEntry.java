package org.xiangqian.note.util;

import lombok.Data;

import java.util.List;

/**
 * @author xiangqian
 * @date 10:56 2024/11/16
 */
@Data
public class ZipEntry {

    private String name;
    private String type; // folder、text、pdf、image
    private String path;
    private Long size;
    private Long lastModifiedTime;

    private String content;
    private List<ZipEntry> parentList;
    private List<ZipEntry> childList;

}
