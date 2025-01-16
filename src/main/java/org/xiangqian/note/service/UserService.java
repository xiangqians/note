package org.xiangqian.note.service;

import org.xiangqian.note.entity.UserEntity;

/**
 * @author xiangqian
 * @date 21:06 2024/02/29
 */
public interface UserService {

    /**
     * 重置密码
     *
     * @param entity
     * @return
     */
    Boolean resetPassword(UserEntity entity);

}
