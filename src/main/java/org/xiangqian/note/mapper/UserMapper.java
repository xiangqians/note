package org.xiangqian.note.mapper;

import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.model.LazyList;

/**
 * @author xiangqian
 * @date 21:41 2024/02/29
 */
@Mapper
public interface UserMapper {

    Boolean add(UserEntity entity);

    Boolean delById(Long id);

    Boolean updById(UserEntity entity);

    UserEntity getById(Long id);

    UserEntity getByName(String name);

}
