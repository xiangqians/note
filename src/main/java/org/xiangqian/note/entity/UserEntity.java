package org.xiangqian.note.entity;

import lombok.Data;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.xiangqian.note.util.DateUtil;

import java.time.Duration;
import java.time.LocalDateTime;
import java.util.Collection;
import java.util.Collections;
import java.util.Objects;

/**
 * 用户信息
 *
 * @author xiangqian
 * @date 21:33 2024/02/29
 */
@Data
public class UserEntity implements UserDetails {

    private static final long serialVersionUID = 1L;

    /**
     * 用户名
     */
    private String name;

    /**
     * 密码
     */
    private String passwd;

    /**
     * 新密码
     */
    private String newPasswd;

    /**
     * 连续错误登陆次数，超过3次则锁定用户
     */
    private Integer deny;

    /**
     * 上一次登录主机
     */
    private String lastLoginHost;

    /**
     * 上一次登录时间戳（单位s）
     */
    private Long lastLoginTime;

    /**
     * 当前登录主机
     */
    private String currentLoginHost;

    /**
     * 当前登录时间戳（单位s）
     */
    private Long currentLoginTime;

    /**
     * 创建时间戳（单位s）
     */
    private Long addTime;

    /**
     * 修改时间戳（单位s）
     */
    private Long updTime;

    /**
     * 获取已被锁定时间（单位s）
     *
     * @return
     */
    public long getLockedTime() {
        return Duration.ofHours(24).toSeconds() - (DateUtil.toSecond(LocalDateTime.now()) - updTime);
    }

    /**
     * 获取用户名
     *
     * @return
     */
    @Override
    public String getUsername() {
        return name;
    }

    /**
     * 获取用户密码
     *
     * @return
     */
    @Override
    public String getPassword() {
        return passwd;
    }

    /**
     * 用户账号是否未被锁定
     *
     * @return
     */
    @Override
    public boolean isAccountNonLocked() {
        // 连续输错密码小于3次
        return deny < 3
                // 锁定24小时
                || Duration.ofSeconds(DateUtil.toSecond(LocalDateTime.now()) - updTime).toHours() >= 24;
    }

    /**
     * 用户账号是否可用
     *
     * @return
     */
    @Override
    public boolean isEnabled() {
        return true;
    }

    /**
     * 用户账号是否未过期
     *
     * @return
     */
    @Override
    public boolean isAccountNonExpired() {
        return true;
    }

    /**
     * 用户凭证（密码）是否未过期
     *
     * @return
     */
    @Override
    public boolean isCredentialsNonExpired() {
        return true;
    }

    /**
     * 用户拥有的权限
     * {@link org.springframework.security.access.expression.SecurityExpressionRoot}
     * {@link org.springframework.security.access.expression.SecurityExpressionRoot#defaultRolePrefix}
     *
     * @return
     */
    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return Collections.emptySet();
    }

    @Override
    public boolean equals(Object object) {
        if (this == object) {
            return true;
        }

        if (object == null || getClass() != object.getClass()) {
            return false;
        }

        return Objects.equals(name, ((UserEntity) object).name);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name);
    }

}
