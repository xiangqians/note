package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.xiangqian.note.entity.UserEntity;

/**
 * @author xiangqian
 * @date 21:41 2024/02/29
 */
@Mapper
public interface UserMapper {

    /**
     * 更新用户信息
     *
     * @param entity
     * @return
     */
    Boolean update(UserEntity entity);

    /**
     * 重置连续错误登陆次数为 0
     *
     * @return
     */
    Boolean resetDeny();

    /**
     * 获取用户信息
     *
     * @return
     */
    UserEntity get();

}
