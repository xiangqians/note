package org.xiangqian.note.configuration;

import com.baomidou.mybatisplus.core.toolkit.PluginUtils;
import com.baomidou.mybatisplus.extension.plugins.MybatisPlusInterceptor;
import com.baomidou.mybatisplus.extension.plugins.inner.InnerInterceptor;
import org.apache.ibatis.executor.Executor;
import org.apache.ibatis.mapping.BoundSql;
import org.apache.ibatis.mapping.MappedStatement;
import org.apache.ibatis.plugin.Interceptor;
import org.apache.ibatis.session.ResultHandler;
import org.apache.ibatis.session.RowBounds;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.annotation.EnableTransactionManagement;
import org.xiangqian.note.util.List;

import java.sql.SQLException;
import java.util.Map;

/**
 * @author xiangqian
 * @date 19:27 2023/09/07
 */
@EnableTransactionManagement
@Configuration(proxyBeanMethods = false)
public class MyBatisConfiguration {

    /**
     * {@link com.baomidou.mybatisplus.extension.plugins.MybatisPlusInterceptor} MybatisPlus插件
     * <p>
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.PaginationInnerInterceptor} 自动分页插件
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.TenantLineInnerInterceptor} 多租户插件
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.DynamicTableNameInnerInterceptor} 动态表名插件
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.OptimisticLockerInnerInterceptor} 乐观锁插件
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.IllegalSQLInnerInterceptor} 性能规范插件
     * {@link com.baomidou.mybatisplus.extension.plugins.inner.BlockAttackInnerInterceptor} 防止全表更新与删除插件
     * <p>
     * {@link com.baomidou.mybatisplus.core.MybatisConfiguration}
     * {@link com.baomidou.mybatisplus.core.MybatisMapperRegistry}
     * {@link com.baomidou.mybatisplus.core.override.MybatisMapperProxyFactory}
     * {@link com.baomidou.mybatisplus.core.override.MybatisMapperProxy}
     * {@link com.baomidou.mybatisplus.core.override.MybatisMapperMethod}
     * {@link com.baomidou.mybatisplus.core.override.MybatisMapperMethod#execute(org.apache.ibatis.session.SqlSession, java.lang.Object[])}
     *
     * @return
     */
    @Bean
    public Interceptor interceptor() {
        MybatisPlusInterceptor mybatisPlusInterceptor = new MybatisPlusInterceptor();
        mybatisPlusInterceptor.addInnerInterceptor(new InnerInterceptor() {
            /**
             * 1）该方法是在 SQL 执行前被调用的
             * 2）在该方法中，可以对查询操作进行预处理，例如修改查询条件、修改查询语句等
             * 3）该方法返回一个布尔值，用于控制是否执行原始的查询操作。如果返回 true，则会继续执行原始的查询操作；如果返回 false，则会终止后续的查询操作
             * 4）如果有多个拦截器同时设置了 willDoQuery 方法，那么它们会按照拦截器注册的顺序依次执行
             *
             * @param executor
             * @param ms
             * @param parameter
             * @param rowBounds
             * @param resultHandler
             * @param boundSql
             * @return
             * @throws SQLException
             */
            @Override
            public boolean willDoQuery(Executor executor, MappedStatement ms, Object parameter, RowBounds rowBounds, ResultHandler resultHandler, BoundSql boundSql) throws SQLException {
                return InnerInterceptor.super.willDoQuery(executor, ms, parameter, rowBounds, resultHandler, boundSql);
            }

            /**
             * 1）该方法是在 SQL 执行前被调用的，与 willDoQuery 方法类似
             * 2）在该方法中，也可以对查询操作进行预处理，例如修改查询条件、修改查询语句等
             * 3）与 willDoQuery 方法不同的是，beforeQuery 方法没有返回值
             * 4）如果有多个拦截器同时设置了 beforeQuery 方法，那么它们会按照拦截器注册的顺序依次执行
             *
             * @param executor
             * @param ms
             * @param parameter
             * @param rowBounds
             * @param resultHandler
             * @param boundSql
             * @throws SQLException
             */
            @Override
            public void beforeQuery(Executor executor, MappedStatement ms, Object parameter, RowBounds rowBounds, ResultHandler resultHandler, BoundSql boundSql) throws SQLException {
                List list = null;
                if (parameter == null) {
                    list = null;
                } else if (parameter.getClass() == List.class) {
                    list = (List) parameter;
                } else if (parameter instanceof Map) {
                    Map<?, ?> parameterMap = (Map) parameter;
                    for (Object value : parameterMap.values()) {
                        if (value instanceof List) {
                            list = (List) value;
                            break;
                        }
                    }
                }

                if (list == null) {
                    return;
                }

                int rows = list.getRows();
                int offset = (list.getOffset() - 1) * rows;
                // 预加载5页
                rows *= 5;
                PluginUtils.mpBoundSql(boundSql).sql(String.format("%s LIMIT %s,%s",
                        boundSql.getSql(),
                        // 指定从结果集的第几行开始返回数据
                        offset,
                        // 指定返回的最大行数
                        rows));
            }
        });
        return mybatisPlusInterceptor;
    }

}