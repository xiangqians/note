package org.xiangqian.note.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Mapper;
import org.xiangqian.note.entity.UserEntity;

/**
 * @author xiangqian
 * @date 21:41 2024/02/29
 */
@Mapper
public interface UserMapper extends BaseMapper<UserEntity> {
}