package org.xiangqian.note.service;

import org.xiangqian.note.entity.UserEntity;

/**
 * @author xiangqian
 * @date 21:06 2024/02/29
 */
public interface UserService {

    Boolean resetPasswd(UserEntity userEntity);

}
