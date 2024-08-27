package org.xiangqian.note.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.session.SessionInformation;
import org.springframework.security.core.session.SessionRegistry;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.mapper.UserMapper;
import org.xiangqian.note.service.UserService;
import org.xiangqian.note.util.SecurityUtil;

import java.util.List;

/**
 * @author xiangqian
 * @date 17:06 2024/02/29
 */
@Slf4j
@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private SessionRegistry sessionRegistry;

    @Autowired
    private UserMapper mapper;

    @Override
    public Boolean resetPasswd(UserEntity vo) {
        UserEntity entity = SecurityUtil.getUser();

        Integer id = entity.getId();
        Assert.notNull(id, "用户id为空");

        String originalPasswd = StringUtils.trim(vo.getOriginalPasswd());
        Assert.isTrue(StringUtils.isNotEmpty(originalPasswd), "原密码不能为空");

        String newPasswd = StringUtils.trim(vo.getNewPasswd());
        String reNewPasswd = StringUtils.trim(vo.getReNewPasswd());
        Assert.isTrue(StringUtils.equals(newPasswd, reNewPasswd), "新密码两次输入不一致");

        Assert.isTrue(passwordEncoder.matches(originalPasswd, entity.getPassword()), "原密码不正确");

        UserEntity updEntity = new UserEntity();
        updEntity.setId(id);
        updEntity.setPasswd(passwordEncoder.encode(newPasswd));
        if (mapper.updById(updEntity)) {
            expireNow(entity.getName());
            return true;
        }
        return false;
    }

    /**
     * 使指定用户session过期
     *
     * @param name {@link UserEntity#getName()}
     */
    private void expireNow(String name) {
        UserEntity entity = new UserEntity();
        entity.setName(name);

        List<SessionInformation> sessions = sessionRegistry.getAllSessions(entity,
                // 是否包括过期的会话
                false);
        if (CollectionUtils.isNotEmpty(sessions)) {
            sessions.forEach(SessionInformation::expireNow);
        }
    }

}
