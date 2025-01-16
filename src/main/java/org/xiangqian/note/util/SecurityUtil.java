package org.xiangqian.note.util;

import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.xiangqian.note.entity.UserEntity;

/**
 * 安全工具
 *
 * @author xiangqian
 * @date 15:26 2024/03/02
 */
public class SecurityUtil {

    /**
     * 获取当前已通过身份验证用户
     *
     * @return
     */
    public static UserEntity getAuthenticatedUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication == null) {
            return null;
        }

        // 匿名用户（表示用户未登录）
        if (authentication instanceof AnonymousAuthenticationToken) {
            return null;
        }

        // 已通过身份验证用户（已登录）
        if (authentication.isAuthenticated()) {
            return (UserEntity) authentication.getPrincipal();
        }

        return null;
    }

}
