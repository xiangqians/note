package org.xiangqian.note.configuration.mybatis;

import org.apache.ibatis.plugin.Interceptor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.annotation.EnableTransactionManagement;

/**
 * @author xiangqian
 * @date 19:27 2023/09/07
 */
@EnableTransactionManagement
@Configuration(proxyBeanMethods = false)
public class MyBatisConfiguration {

    /**
     * 注册延迟加载列表拦截器
     *
     * @return
     */
    @Bean
    public Interceptor mybatisInterceptor() {
        return new MybatisInterceptor();
    }

}
