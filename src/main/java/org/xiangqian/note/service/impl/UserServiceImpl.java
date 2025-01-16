package org.xiangqian.note.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
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
public class UserServiceImpl implements UserService, ApplicationRunner {

    @Autowired
    private UserMapper mapper;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private SessionRegistry sessionRegistry;

    @Override
    public Boolean resetPassword(UserEntity entity) {
        // 校验密码
        String password = entity.getPassword();
        Assert.isTrue(StringUtils.isNotEmpty(password), "密码不能为空");
        Assert.isTrue(password.length() <= 128, "密码长度不能大于128个字符");

        // 校验新密码
        String newPassword = entity.getNewPassword();
        Assert.isTrue(StringUtils.isNotEmpty(newPassword), "新密码不能为空");
        Assert.isTrue(newPassword.length() <= 128, "新密码长度不能大于128个字符");

        // 获取当前已通过身份验证用户
        UserEntity authenticatedEntity = SecurityUtil.getAuthenticatedUser();

        // 校验密码是否正确
        Assert.isTrue(passwordEncoder.matches(password, authenticatedEntity.getPassword()), "密码不正确");

        // 更新密码
        UserEntity updateEntity = new UserEntity();
        updateEntity.setPassword(passwordEncoder.encode(newPassword));
        if (mapper.update(updateEntity)) {
            expireNow(authenticatedEntity.getName());
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
                true);
        if (CollectionUtils.isNotEmpty(sessions)) {
            sessions.forEach(SessionInformation::expireNow);
        }
    }

    /**
     * 每次启动后重置连续错误登陆次数为 0
     *
     * @param args
     * @throws Exception
     */
    @Override
    public void run(ApplicationArguments args) throws Exception {
        mapper.resetDeny();
    }

}
