package org.xiangqian.note.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.security.core.session.SessionInformation;
import org.springframework.security.core.session.SessionRegistry;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
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

    @Autowired
    @Qualifier("userEntityThreadLocal")
    private ThreadLocal<UserEntity> userEntityThreadLocal;

    @Override
    public UserDetails loadUserByUsername(String userName) throws UsernameNotFoundException {
        UserEntity entity = mapper.selectOne(new LambdaQueryWrapper<UserEntity>().last("LIMIT 1"));
        if (entity == null) {
            throw new UsernameNotFoundException(userName);
        }

        // 将用户信息设置到线程本地
        userEntityThreadLocal.set(entity);

        // 未被锁定
        if (entity.isNonLocked()) {
            // 未被限时锁定
            if (entity.isNonLimitedTimeLocked()) {
                // 用户连续错误登陆超过3次，则归零
                if (entity.getDeny() >= 3) {
                    UserEntity updEntity = new UserEntity();
                    updEntity.setId(entity.getId());
                    updEntity.setDeny(0);
                    mapper.updateById(updEntity);
                    entity.setDeny(0);
                }
            }
            // 限时锁定
            else {
                entity.setLocked(1);
            }
        }

        return entity;
    }

    @Override
    public Boolean resetPasswd(UserEntity vo) {
        UserEntity entity = SecurityUtil.getUser();

        Long id = entity.getId();
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
        if (mapper.updateById(updEntity) > 0) {
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
