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
    private PasswordEncoder passwordEncoder;

    @Autowired
    private SessionRegistry sessionRegistry;

    @Autowired
    private UserMapper userMapper;

    @Override
    public Boolean resetPasswd(UserEntity userEntity) {
        // 校验密码
        String passwd = userEntity.getPasswd();
        Assert.isTrue(StringUtils.isNotEmpty(passwd), "密码不能为空");
        Assert.isTrue(passwd.length() <= 128, "密码长度不能大于128个字符");

        // 校验新密码
        String newPasswd = userEntity.getNewPasswd();
        Assert.isTrue(StringUtils.isNotEmpty(newPasswd), "新密码不能为空");
        Assert.isTrue(newPasswd.length() <= 128, "新密码长度不能大于128个字符");

        // 获取当前已通过身份验证用户
        UserEntity securityUserEntity = SecurityUtil.getUser();

        // 校验密码是否正确
        Assert.isTrue(passwordEncoder.matches(passwd, securityUserEntity.getPassword()), "密码不正确");

        // 用户名
        String name = securityUserEntity.getName();

        // 更新密码
        UserEntity updUserEntity = new UserEntity();
        updUserEntity.setName(name);
        updUserEntity.setPasswd(passwordEncoder.encode(newPasswd));
        if (userMapper.updByName(updUserEntity)) {
            expireNow(name);
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
        UserEntity userEntity = new UserEntity();
        userEntity.setName(name);

        List<SessionInformation> sessions = sessionRegistry.getAllSessions(userEntity,
                // 是否包括过期的会话
                false);
        if (CollectionUtils.isNotEmpty(sessions)) {
            sessions.forEach(SessionInformation::expireNow);
        }
    }

    /**
     * 每次启动成功重置连续错误登陆次数
     *
     * @param args
     * @throws Exception
     */
    @Override
    public void run(ApplicationArguments args) throws Exception {
        UserEntity updUserEntity = new UserEntity();
        updUserEntity.setName("admin");
        updUserEntity.setDeny(0);
        userMapper.updByName(updUserEntity);
    }

}
