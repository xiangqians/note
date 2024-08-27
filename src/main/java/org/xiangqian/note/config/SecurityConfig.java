package org.xiangqian.note.config;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.*;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.boot.web.servlet.ServletListenerRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.LockedException;
import org.springframework.security.config.annotation.method.configuration.EnableMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.annotation.web.configurers.HeadersConfigurer;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.session.SessionInformation;
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
import org.springframework.web.filter.HiddenHttpMethodFilter;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.mapper.UserMapper;
import org.xiangqian.note.util.DateUtil;

import java.io.IOException;
import java.time.LocalDateTime;
import java.util.List;

/**
 * @author xiangqian
 * @date 21:02 2024/02/29
 */
@EnableWebSecurity
@Configuration(proxyBeanMethods = false)
@EnableMethodSecurity(prePostEnabled = true, // 开启预处理验证
        jsr250Enabled = true, // 启用JSR250注解支持
        securedEnabled = true) // 启用 {@link org.springframework.security.access.annotation.Secured} 注解支持
public class SecurityConfig implements WebMvcConfigurer, UserDetailsService {

    @Autowired
    private UserMapper userMapper;

    private ThreadLocal<UserEntity> threadLocal = new ThreadLocal<>();

    // 处理静态资源
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        registry.addResourceHandler("/static/**").addResourceLocations("classpath:/static/");
    }

    @Override
    public UserDetails loadUserByUsername(String userName) throws UsernameNotFoundException {
        UserEntity entity = null;
        userName = StringUtils.trim(userName);
        if (StringUtils.isNotEmpty(userName)) {
            entity = userMapper.getByName(userName);
        }
        if (entity == null) {
            throw new UsernameNotFoundException(userName);
        }

        // 将用户信息设置到线程本地
        threadLocal.set(entity);

        // 未被锁定
        if (entity.isNonLocked()) {
            // 未被限时锁定
            if (entity.isNonLimitedTimeLocked()) {
                // 用户连续错误登陆超过3次，则归零
                if (entity.getDeny() >= 3) {
                    UserEntity updEntity = new UserEntity();
                    updEntity.setId(entity.getId());
                    updEntity.setDeny(0);
                    userMapper.updById(updEntity);
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

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http,
                                                   AuthenticationSuccessHandler authenticationSuccessHandler,
                                                   AuthenticationFailureHandler authenticationFailureHandler,
                                                   SessionRegistry sessionRegistry) throws Exception {
        http
                // http授权请求配置
                .authorizeHttpRequests((authorize) -> authorize
                        // 放行静态资源
                        .requestMatchers("/static/**").permitAll()
                        // 放行/login
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
                        // 当达到最大会话数时阻止新的登录，确保用户无法在多个地方同时登录
                        //.maxSessionsPreventsLogin(true)
                        // 当达到最大会话数时剔除旧登录
                        .sessionRegistry(sessionRegistry)
                        // 会话过期后跳转的URL
                        .expiredUrl("/login"))
                .headers(headers -> headers.frameOptions(HeadersConfigurer.FrameOptionsConfig::disable))
                .csrf(AbstractHttpConfigurer::disable);
        return http.build();
    }

    @Bean
    public AuthenticationSuccessHandler authenticationSuccessHandler() {
        return new SavedRequestAwareAuthenticationSuccessHandler() {
            @Autowired
            private SessionRegistry sessionRegistry;

            @Autowired
            private UserMapper mapper;

            @Override
            public void onAuthenticationSuccess(HttpServletRequest request, HttpServletResponse response, Authentication authentication) throws ServletException, IOException {
                threadLocal.remove();

                // 获取已授权用户信息
                UserEntity entity = (UserEntity) authentication.getPrincipal();

                // 获取会话
                HttpSession session = request.getSession(true);

                // 限制同一个用户只允许一端登录
                List<SessionInformation> sessions = sessionRegistry.getAllSessions(entity,
                        // 是否包括过期的会话
                        false);
                if (CollectionUtils.isNotEmpty(sessions)) {
                    String sessionId = session.getId();
                    for (SessionInformation session0 : sessions) {
                        if (!sessionId.equals(session0.getSessionId())) {
                            session0.expireNow();
                        }
                    }
                }

                // 更新用户信息
                UserEntity updEntity = new UserEntity();
                updEntity.setId(entity.getId());
                updEntity.setDeny(0);
                updEntity.setLastLoginIp(entity.getCurrentLoginIp());
                updEntity.setLastLoginTime(entity.getCurrentLoginTime());
                updEntity.setCurrentLoginIp(request.getRemoteHost());
                updEntity.setCurrentLoginTime(DateUtil.toSecond(LocalDateTime.now()));
                updEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
                mapper.updById(updEntity);

                // 更新会话中的用户信息
                entity.setDeny(updEntity.getDeny());
                entity.setLastLoginIp(updEntity.getLastLoginIp());
                entity.setLastLoginTime(updEntity.getLastLoginTime());
                entity.setCurrentLoginIp(updEntity.getCurrentLoginIp());
                entity.setCurrentLoginTime(updEntity.getCurrentLoginTime());
                entity.setUpdTime(updEntity.getUpdTime());

                // 添加用户信息到会话
                AbsController.setLoggedinAttribute(session, true);
                AbsController.setUserAttribute(session, entity);

                super.onAuthenticationSuccess(request, response, authentication);
            }
        };
    }

    @Bean
    public AuthenticationFailureHandler authenticationFailureHandler() {
        return new AuthenticationFailureHandler() {
            @Autowired
            private UserMapper mapper;

            @Override
            public void onAuthenticationFailure(HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) throws IOException, ServletException {
                UserEntity entity = threadLocal.get();
                threadLocal.remove();

                // 获取会话
                HttpSession session = request.getSession(true);

                // 将错误信息添加到会话
                String error = null;
                if (exception instanceof BadCredentialsException) {
//                    error = "用户名或密码不正确";
                    error = "密码不正确";
                    if (entity != null) {
                        UserEntity updEntity = new UserEntity();
                        updEntity.setId(entity.getId());
                        updEntity.setDeny(entity.getDeny() + 1);
                        updEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
                        mapper.updById(updEntity);
                        if (updEntity.getDeny() == 2) {
                            error = "已连续两次输错密码，如连续输错三次，用户将被锁定";
                        }
                    }
                } else if (exception instanceof LockedException) {
                    if (entity.isNonLimitedTimeLocked()) {
                        error = "用户已被锁定";
                    } else {
                        error = String.format("用户已被锁定（%s后解锁）", DateUtil.humanDurationSecond(entity.getLimitedTimeLockedTime()));
                    }
                } else {
                    error = exception.getMessage();
                }
                AbsController.setErrorAttribute(session, error);

                // 将用户信息vo存储到session
                entity = new UserEntity();
                entity.setName(request.getParameter("username"));
                AbsController.setVoAttribute(session, entity);

                // 重定向到登录页
                response.sendRedirect("/login");
            }
        };
    }

    @Bean
    public SessionRegistry sessionRegistry() {
        return new SessionRegistryImpl();
    }

    @Bean
    public ServletListenerRegistrationBean<HttpSessionListener> httpSessionListener() {
        ServletListenerRegistrationBean<HttpSessionListener> servletListenerRegistrationBean = new ServletListenerRegistrationBean<>();
        servletListenerRegistrationBean.setListener(new HttpSessionListener() {
            // 设置会话默认属性
            @Override
            public void sessionCreated(HttpSessionEvent event) {
                HttpSession session = event.getSession();
                AbsController.setLoggedinAttribute(session, false);
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
                if ("/login".equals(request.getServletPath()) && AbsController.getLoggedinAttribute(request.getSession(true))) {
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
