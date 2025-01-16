package org.xiangqian.note.configuration;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.HttpSession;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
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
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.filter.HiddenHttpMethodFilter;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.mapper.UserMapper;
import org.xiangqian.note.util.Model;
import org.xiangqian.note.util.SecurityUtil;
import org.xiangqian.note.util.TimeUtil;

import java.io.IOException;

/**
 * @author xiangqian
 * @date 21:02 2024/02/29
 */
@EnableWebSecurity
@Configuration(proxyBeanMethods = false)
public class SecurityConfiguration implements WebMvcConfigurer {

    private final String USER = "user";

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
                // http授权请求
                .authorizeHttpRequests((authorize) -> authorize
                        // 允许未经授权访问静态资源请求
                        .requestMatchers("/static/**").permitAll()
                        // 允许未经授权访问登录请求
                        .requestMatchers("/login").permitAll()
                        // 其他请求需要授权才能访问
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
            private UserMapper mapper;

            @Override
            public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
                UserEntity entity = mapper.get();

                HttpServletRequest request = ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes()).getRequest();
                request.setAttribute(USER, entity);

                return entity;
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
                request.removeAttribute(USER);

                // 获取已授权用户信息
                UserEntity entity = (UserEntity) authentication.getPrincipal();

                // 更新数据库用户信息
                UserEntity updateEntity = new UserEntity();
                updateEntity.setDeny(0);
                updateEntity.setLastLoginHost(entity.getCurrentLoginHost());
                updateEntity.setLastLoginTime(entity.getCurrentLoginTime());
                updateEntity.setCurrentLoginHost(request.getRemoteHost());
                updateEntity.setCurrentLoginTime(TimeUtil.now());
                updateEntity.setUpdateTime(TimeUtil.now());
                userMapper.update(updateEntity);

                // 更新会话中用户信息
                entity.setDeny(updateEntity.getDeny());
                entity.setLastLoginHost(updateEntity.getLastLoginHost());
                entity.setLastLoginTime(updateEntity.getLastLoginTime());
                entity.setCurrentLoginHost(updateEntity.getCurrentLoginHost());
                entity.setCurrentLoginTime(updateEntity.getCurrentLoginTime());
                entity.setUpdateTime(updateEntity.getUpdateTime());

                super.onAuthenticationSuccess(request, response, authentication);
            }
        };
    }

    @Bean
    public AuthenticationFailureHandler authenticationFailureHandler() {
        return new SimpleUrlAuthenticationFailureHandler("/login") {
            @Autowired
            private UserMapper mapper;

            @Override
            public void onAuthenticationFailure(HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) throws IOException, ServletException {
                UserEntity entity = (UserEntity) request.getAttribute(USER);
                request.removeAttribute(USER);

                String message = null;

                // 凭据错误异常
                if (exception instanceof BadCredentialsException) {
                    if (entity != null) {
                        UserEntity updateEntity = new UserEntity();

                        int deny = entity.getDeny();
                        if (entity.getLockedTime() <= 0) {
                            deny = 0;
                        }
                        deny += 1;
                        updateEntity.setDeny(deny);

                        updateEntity.setUpdateTime(TimeUtil.now());
                        mapper.update(updateEntity);
                        if (deny == 2) {
                            message = "已连续两次输错密码，如连续输错三次，系统将被锁定";
                        }
                    }
                    if (message == null) {
                        message = "密码不正确";
                    }
                }
                // 锁定异常
                else if (exception instanceof LockedException) {
                    message = String.format("系统已被锁定（%s后解锁）", TimeUtil.humanDuration(entity.getLockedTime()));
                } else {
                    message = exception.getMessage();
                }

                HttpSession session = request.getSession(true);
                session.setAttribute(AbsController.MODEL, Model.of(AbsController.MESSAGE, message));

                super.onAuthenticationFailure(request, response, exception);
            }
        };
    }

    /**
     * <form method="post" action="/your-endpoint">
     * <input type="hidden" name="_method" value="PUT">
     * <button type="submit">提交</button>
     * </form>
     * 由于 HTML 表单只支持 GET 和 POST 请求方法，因此需要进行一些额外的处理来实现 PUT 和 DELETE 请求
     * 配置过滤器 {@link HiddenHttpMethodFilter} 来解析这个名为 _method 的隐藏字段，并将请求方法修改为 PUT 或 DELETE
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
                if ("/login".equals(request.getServletPath()) && SecurityUtil.getAuthenticatedUser() != null) {
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
