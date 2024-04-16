package org.xiangqian.note;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.EnableAspectJAutoProxy;

/**
 * @author xiangqian
 * @date 21:19 2024/02/29
 */
@SpringBootApplication
@EnableAspectJAutoProxy
public class NoteApplication {

    public static void main(String[] args) {
        SpringApplication.run(NoteApplication.class, args);
    }

}
