package org.xiangqian.note.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.NoArgsConstructor;
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
@TableName("user")
@NoArgsConstructor
public class UserEntity implements UserDetails {

    private static final long serialVersionUID = 1L;

    // 用户id
    @TableId(type = IdType.AUTO)
    private Long id;

    // 用户名
    @TableField("`name`")
    private String name;

    // 昵称
    private String nickname;

    // 密码
    private String passwd;

    // 原密码
    @TableField(exist = false)
    private String originalPasswd;

    // 新密码
    @TableField(exist = false)
    private String newPasswd;

    // 再次输入新密码
    @TableField(exist = false)
    private String reNewPasswd;

    // 是否已锁定，0-否，1-是
    private Integer locked;

    // 用户连续错误登陆次数，超过3次则锁定用户
    private Integer deny;

    // 上一次登录ip
    private String lastLoginIp;

    // 上一次登录时间（时间戳，单位s）
    private Long lastLoginTime;

    // 当前次登录ip
    private String currentLoginIp;

    // 当前登录时间（时间戳，单位s）
    private Long currentLoginTime;

    // 创建时间（时间戳，单位s）
    private Long addTime;

    // 修改时间（时间戳，单位s）
    private Long updTime;

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

    /**
     * 是否未被锁定
     *
     * @return
     */
    public boolean isNonLocked() {
        return locked == 0;
    }

    /**
     * 是否未被限时锁定
     *
     * @return
     */
    public boolean isNonLimitedTimeLocked() {
        // 连续输错密码小于3次
        return deny < 3
                // 锁定24小时
                || Duration.ofSeconds(DateUtil.toSecond(LocalDateTime.now()) - updTime).toHours() >= 24;
    }

    /**
     * 获取已被限时锁定时间（单位s）
     *
     * @return
     */
    public long getLimitedTimeLockedTime() {
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
        return isNonLocked();
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

}
