package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.xiangqian.note.entity.UserEntity;

/**
 * @author xiangqian
 * @date 21:41 2024/02/29
 */
@Mapper
public interface UserMapper {

    Boolean updByName(UserEntity userEntity);

    UserEntity getByName(String name);

}
