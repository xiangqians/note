package org.xiangqian.note.config;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.boot.web.servlet.ServletListenerRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.LockedException;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.annotation.web.configurers.HeadersConfigurer;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.session.SessionRegistry;
import org.springframework.security.core.session.SessionRegistryImpl;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.security.web.authentication.SavedRequestAwareAuthenticationSuccessHandler;
import org.springframework.security.web.authentication.SimpleUrlAuthenticationFailureHandler;
import org.springframework.web.filter.HiddenHttpMethodFilter;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.mapper.UserMapper;
import org.xiangqian.note.model.Vo;
import org.xiangqian.note.util.DateUtil;

import java.io.IOException;
import java.time.LocalDateTime;

/**
 * @author xiangqian
 * @date 21:02 2024/02/29
 */
@EnableWebSecurity
@Configuration(proxyBeanMethods = false)
public class SecurityConfig implements WebMvcConfigurer {

    private final String IS_LOGGED_IN = "isLoggedIn";

    private final ThreadLocal<UserEntity> threadLocal = new ThreadLocal<>();

    /**
     * 处理静态资源
     *
     * @param registry
     */
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        registry.addResourceHandler("/static/**").addResourceLocations("classpath:/static/");
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http,
                                                   AuthenticationSuccessHandler authenticationSuccessHandler,
                                                   AuthenticationFailureHandler authenticationFailureHandler,
                                                   SessionRegistry sessionRegistry) throws Exception {
        http
                // http授权请求配置
                .authorizeHttpRequests((authorize) -> authorize
                        // 放行静态资源请求
                        .requestMatchers("/static/**").permitAll()
                        // 放行登录请求
                        .requestMatchers("/login").permitAll()
                        // 其他请求需要授权
                        .anyRequest().authenticated())
                // 自定义表单登录
                .formLogin(configurer -> configurer
                        // 自定义登录页面
                        .loginPage("/login")
                        // 登录成功处理器
                        .successHandler(authenticationSuccessHandler)
                        // 登录失败处理器
                        .failureHandler(authenticationFailureHandler)
                        // 放行资源
                        .permitAll())
                // 处理权限不足的请求
                .exceptionHandling(configurer -> configurer.accessDeniedPage("/error"))
                // session会话管理：限制同一个用户只允许一端登录
                .sessionManagement(configurer -> configurer
                        // 设置每个用户最大会话数为1，这意味着同一个用户只允许一端登录
                        .maximumSessions(1)
                        // 如果用户已经达到了最大会话数量限制（maximumSessions），则新登录请求将会被拒绝，用户无法登录
                        //.maxSessionsPreventsLogin(true)
                        // 如果用户已经达到了最大会话数量限制（maximumSessions），允许新登录请求，旧的会话将会被失效
                        .maxSessionsPreventsLogin(false)
                        // 当达到最大会话数时剔除旧登录
                        .sessionRegistry(sessionRegistry)
                        // 会话过期后跳转的URL
                        .expiredUrl("/login"))
                .headers(headers -> headers.frameOptions(HeadersConfigurer.FrameOptionsConfig::disable))
                .csrf(AbstractHttpConfigurer::disable);
        return http.build();
    }

    @Bean
    public UserDetailsService userDetailsService() {
        return new UserDetailsService() {
            @Autowired
            private UserMapper userMapper;

            @Override
            public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
                username = "admin";
                UserEntity userEntity = userMapper.getByName(username);
                if (userEntity == null) {
                    throw new UsernameNotFoundException(username);
                }

                // 将用户信息设置到线程本地
                threadLocal.set(userEntity);

                return userEntity;
            }
        };
    }

    @Bean
    public SessionRegistry sessionRegistry() {
        return new SessionRegistryImpl();
    }

    @Bean
    public AuthenticationSuccessHandler authenticationSuccessHandler() {
        return new SavedRequestAwareAuthenticationSuccessHandler() {
            @Autowired
            private UserMapper userMapper;

            @Override
            public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws ServletException, IOException {
                threadLocal.remove();

                // 获取已授权用户信息
                UserEntity userEntity = (UserEntity) authentication.getPrincipal();

                // 更新数据库用户信息
                UserEntity updUserEntity = new UserEntity();
                updUserEntity.setId(userEntity.getId());
                updUserEntity.setDeny(0);
                updUserEntity.setLastLoginHost(userEntity.getCurrentLoginHost());
                updUserEntity.setLastLoginTime(userEntity.getCurrentLoginTime());
                updUserEntity.setCurrentLoginHost(request.getRemoteHost());
                updUserEntity.setCurrentLoginTime(DateUtil.toSecond(LocalDateTime.now()));
                updUserEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
                userMapper.updById(updUserEntity);

                // 更新会话中用户信息
                userEntity.setDeny(updUserEntity.getDeny());
                userEntity.setDeny(0);
                userEntity.setLastLoginHost(updUserEntity.getLastLoginHost());
                userEntity.setLastLoginTime(updUserEntity.getLastLoginTime());
                userEntity.setCurrentLoginHost(updUserEntity.getCurrentLoginHost());
                userEntity.setCurrentLoginTime(updUserEntity.getCurrentLoginTime());
                userEntity.setUpdTime(updUserEntity.getUpdTime());

                // 获取会话
                HttpSession session = request.getSession(true);
                // 设置会话用户为已登录状态
                session.setAttribute(IS_LOGGED_IN, true);
                // 设置会话用户信息
                session.setAttribute("user", userEntity);

                super.onAuthenticationSuccess(request, response, authentication);
            }
        };
    }

    @Bean
    public AuthenticationFailureHandler authenticationFailureHandler() {
        return new SimpleUrlAuthenticationFailureHandler("/login") {
            @Autowired
            private UserMapper userMapper;

            @Override
            public void onAuthenticationFailure(HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) throws IOException, ServletException {
                UserEntity userEntity = threadLocal.get();
                threadLocal.remove();

                String error = null;

                // 凭据错误异常
                if (exception instanceof BadCredentialsException) {
                    if (userEntity != null) {
                        UserEntity updUserEntity = new UserEntity();
                        updUserEntity.setId(userEntity.getId());

                        int deny = userEntity.getDeny();
                        if (userEntity.getLimitedTimeLockedTime() <= 0) {
                            deny = 0;
                        }
                        deny += 1;
                        updUserEntity.setDeny(deny);

                        updUserEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
                        userMapper.updById(updUserEntity);
                        if (deny == 2) {
                            error = "已连续两次输错密码，如连续输错三次，系统将被锁定";
                        }
                    }
                    if (error == null) {
                        error = "密码不正确";
                    }
                }
                // 锁定异常
                else if (exception instanceof LockedException) {
                    if (userEntity != null && userEntity.isNonLocked() && !userEntity.isNonLimitedTimeLocked()) {
                        error = String.format("系统已被锁定（%s后解锁）", DateUtil.humanDurationSecond(userEntity.getLimitedTimeLockedTime()));
                    } else {
                        error = "系统已被锁定";
                    }
                } else {
                    error = exception.getMessage();
                }

                HttpSession session = request.getSession(true);
                session.setAttribute(AbsController.VO, Vo.error(error));

                super.onAuthenticationFailure(request, response, exception);
            }
        };
    }

    @Bean
    public ServletListenerRegistrationBean<HttpSessionListener> httpSessionListener() {
        ServletListenerRegistrationBean<HttpSessionListener> servletListenerRegistrationBean = new ServletListenerRegistrationBean<>();
        servletListenerRegistrationBean.setListener(new HttpSessionListener() {
            // 设置会话默认属性
            @Override
            public void sessionCreated(HttpSessionEvent event) {
                HttpSession session = event.getSession();
                session.setAttribute(IS_LOGGED_IN, false);
            }
        });
        return servletListenerRegistrationBean;
    }

    /**
     * <form method="post" action="/your-endpoint">
     * <input type="hidden" name="_method" value="PUT">
     * <button type="submit">提交</button>
     * </form>
     * 由于HTML表单只支持GET和POST请求方法，因此需要进行一些额外的处理来实现PUT和DELETE请求
     * 配置过滤器 {@link HiddenHttpMethodFilter} 来解析这个名为 _method 的隐藏字段，并将请求方法修改为PUT或DELETE
     *
     * @return
     */
    @Bean
    public FilterRegistrationBean<HiddenHttpMethodFilter> hiddenHttpMethodFilter() {
        FilterRegistrationBean<HiddenHttpMethodFilter> filterRegistrationBean = new FilterRegistrationBean<>(new HiddenHttpMethodFilter());
        filterRegistrationBean.addUrlPatterns("/*");
        return filterRegistrationBean;
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(new HandlerInterceptor() {
            @Override
            public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
                if ("/login".equals(request.getServletPath()) && (boolean) request.getSession(true).getAttribute(IS_LOGGED_IN)) {
                    // 已登录，重定向到首页
                    response.sendRedirect("/");
                    // 不继续执行后续的拦截器
                    return false;
                }
                // 继续执行后续的拦截器
                return true;
            }
        });
    }

}
