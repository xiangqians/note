package org.xiangqian.note.model;

import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author xiangqian
 * @date 15:07 2024/05/03
 */
@Data
@NoArgsConstructor
public class Response<T> {

    private String code;
    private String msg;
    private T data;

    private Response(String code, String msg, T data) {
        this.code = code;
        this.msg = msg;
        this.data = data;
    }

    public static <T> Response<T> ok(T data) {
        return new Response<>("ok", null, data);
    }

    public static <T> Response<T> error(String msg) {
        return new Response<>("error", msg, null);
    }

}
