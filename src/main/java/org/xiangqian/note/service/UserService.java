package org.xiangqian.note.service;

import org.springframework.security.core.userdetails.UserDetailsService;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.entity.IavEntity;

import java.util.List;

/**
 * @author xiangqian
 * @date 21:06 2024/02/29
 */
public interface UserService extends UserDetailsService {

    Boolean resetPasswd(UserEntity vo);

}
