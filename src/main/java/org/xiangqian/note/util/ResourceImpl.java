package org.xiangqian.note.util;

import lombok.Getter;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.http.MediaType;

/**
 * @author xiangqian
 * @date 11:38 2024/11/10
 */
public class ResourceImpl extends ByteArrayResource {

    @Getter
    private String name;

    @Getter
    private MediaType type;

    public ResourceImpl(String name, MediaType type, byte[] bytes) {
        super(bytes);
        this.name = name;
        this.type = type;
    }

}
