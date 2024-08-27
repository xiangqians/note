package org.xiangqian.note.config.mybatis;

import org.apache.commons.collections4.CollectionUtils;
import org.apache.ibatis.builder.StaticSqlSource;
import org.apache.ibatis.cache.CacheKey;
import org.apache.ibatis.executor.Executor;
import org.apache.ibatis.mapping.BoundSql;
import org.apache.ibatis.mapping.MappedStatement;
import org.apache.ibatis.mapping.ParameterMapping;
import org.apache.ibatis.mapping.ResultMap;
import org.apache.ibatis.plugin.Interceptor;
import org.apache.ibatis.plugin.Intercepts;
import org.apache.ibatis.plugin.Invocation;
import org.apache.ibatis.plugin.Signature;
import org.apache.ibatis.session.Configuration;
import org.apache.ibatis.session.ResultHandler;
import org.apache.ibatis.session.RowBounds;
import org.xiangqian.note.model.LazyList;
import org.xiangqian.note.model.Page;

import java.lang.reflect.Method;
import java.sql.SQLException;
import java.util.*;

/**
 * 分页、延迟加载列表拦截器
 *
 * @author xiangqian
 * @date 21:48 2024/07/15
 */
@Intercepts({ // 配置多个拦截点（每个拦截点指定了拦截的目标方法和参数类型）
        // 拦截 Executor 类的 query 方法，参数包括 MappedStatement.class、Object.class、RowBounds.class 和 ResultHandler.class
        @Signature(type = Executor.class, method = "query", args = {MappedStatement.class, Object.class, RowBounds.class, ResultHandler.class}),
        // 拦截 Executor 类的 query 方法，参数为 MappedStatement.class、Object.class、RowBounds.class、ResultHandler.class、CacheKey.class 和 BoundSql.class
        @Signature(type = Executor.class, method = "query", args = {MappedStatement.class, Object.class, RowBounds.class, ResultHandler.class, CacheKey.class, BoundSql.class}),
})
public class MybatisInterceptor implements Interceptor {

    /**
     * 拦截目标对象的方法执行
     *
     * @param invocation 包含了被拦截的目标对象、方法名和方法参数等信息
     * @return
     * @throws Throwable
     */
    @Override
    public Object intercept(Invocation invocation) throws Throwable {
        Object target = invocation.getTarget();
        if (target instanceof Executor) {
            Method method = invocation.getMethod();
            Object[] args = invocation.getArgs();
            if ("query".equals(method.getName()) && args != null) {
                if (args.length == 4) {
                    return query(invocation, (Executor) target, (MappedStatement) args[0], args[1], (RowBounds) args[2], (ResultHandler) args[3], null, null);
                } else if (args.length == 6) {
                    return query(invocation, (Executor) target, (MappedStatement) args[0], args[1], (RowBounds) args[2], (ResultHandler) args[3], (CacheKey) args[4], (BoundSql) args[5]);
                }
            }
        }

        // 继续执行下一个拦截器
        return invocation.proceed();
    }

    /**
     * 生成目标对象的代理对象
     * <p>
     * 在 MyBatis 中，拦截器可以拦截的对象类型通常是 Executor、StatementHandler、ParameterHandler 和 ResultSetHandler。
     * 通过调用 Plugin.wrap(target, this) 来创建一个代理对象，将该拦截器应用到目标对象上。
     *
     * @param target
     * @return
     */
    @Override
    public Object plugin(Object target) {
        return Interceptor.super.plugin(target);
    }

    /**
     * 设置拦截器的属性
     * <p>
     * MyBatis 配置文件中可以通过 <plugin> 标签为拦截器配置属性，这些属性会在初始化拦截器时传递到这个方法中。
     * 你可以在这个方法中读取和设置拦截器的配置属性，以便拦截器能够根据配置表现不同的行为。
     *
     * @param properties
     */
    @Override
    public void setProperties(Properties properties) {
    }

    /**
     * 查询
     *
     * @param invocation      下一个拦截器
     * @param executor        执行器
     * @param mappedStatement 表示映射的 SQL 语句以及相关的配置信息，包括 SQL 的 ID、参数映射、结果映射等
     * @param parameter       执行 SQL 时需要传递的参数对象，可以是一个简单类型、一个 JavaBean 或者一个 Map
     * @param rowBounds       用于控制结果集的返回行数和偏移量的对象，用于分页查询
     * @param resultHandler   结果处理器，MyBatis 执行 SQL 后将结果传递给它，用于自定义结果集的处理
     * @param cacheKey        缓存键对象，用于标识查询结果的缓存键，用于二级缓存
     * @param boundSql        包含了 SQL 语句及其参数映射的对象，用于记录 SQL 语句的完整信息，方便进行 SQL 的拦截和修改
     * @return
     * @throws Throwable
     */
    private Object query(Invocation invocation, Executor executor, MappedStatement mappedStatement, Object parameter, RowBounds rowBounds, ResultHandler resultHandler, CacheKey cacheKey, BoundSql boundSql) throws Throwable {
        // 获取参数中的 Page/LazyList 实例
        Page<?> page = null;
        LazyList<?> lazyList = null;
        if (parameter != null) {
            if (parameter instanceof LazyList) {
                page = (Page) parameter;
            } else if (parameter instanceof Page) {
                lazyList = (LazyList) parameter;
            } else if (parameter instanceof Map) {
                Map parameterMap = (Map) parameter;
                for (Object value : parameterMap.values()) {
                    if (value instanceof Page) {
                        page = (Page) value;
                        break;
                    } else if (value instanceof LazyList) {
                        lazyList = (LazyList) value;
                        break;
                    }
                }
            }
        }

        // 分页查询
        if (page != null) {
            return page(page, executor, mappedStatement, parameter, rowBounds, resultHandler, cacheKey, boundSql);
        }
        // 延迟加载列表查询
        else if (lazyList != null) {
            return lazyList(lazyList, executor, mappedStatement, parameter, rowBounds, resultHandler, cacheKey, boundSql);
        }

        // 继续执行下一个拦截器
        return invocation.proceed();
    }

    private Object page(Page<?> page, Executor executor, MappedStatement mappedStatement, Object parameter, RowBounds rowBounds, ResultHandler resultHandler, CacheKey cacheKey, BoundSql boundSql) throws SQLException {
        // 获取绑定 SQL
        if (boundSql == null) {
            boundSql = mappedStatement.getBoundSql(parameter);
        }

        // 查询总数
        MappedStatement countMappedStatement = getCountMappedStatement(mappedStatement, boundSql);
        List list = executor.query(countMappedStatement, parameter, rowBounds, resultHandler);
        int total = (Integer) list.get(0);
        page.setTotal(total);

        // 当前页
        int current = page.getCurrent();
        // 页数量
        int size = page.getSize();
        // 最大页数
        int number = total / size;
        if (total % size > 0) {
            number += 1;
        }

        // 查询分页数据
        if (current <= number) {
            int offset = (current - 1) * size;
            int rows = size;
            MappedStatement listMappedStatement = getListMappedStatement(mappedStatement, parameter, boundSql, offset, rows);
            list = executor.query(listMappedStatement, parameter, rowBounds, resultHandler);
            page.setData(list);
        }

        return List.of(page);
    }

    private Object lazyList(LazyList<?> lazyList, Executor executor, MappedStatement mappedStatement, Object parameter, RowBounds rowBounds, ResultHandler resultHandler, CacheKey cacheKey, BoundSql boundSql) throws SQLException {
        // 获取绑定 SQL
        if (boundSql == null) {
            boundSql = mappedStatement.getBoundSql(parameter);
        }

        // 当前页
        int current = lazyList.getCurrent();
        // 页数量
        int size = lazyList.getSize();

        // 查询列表数据
        int offset = (current - 1) * size;
        int rows = size + 1;
        MappedStatement listMappedStatement = getListMappedStatement(mappedStatement, parameter, boundSql, offset, rows);
        List list = executor.query(listMappedStatement, parameter, rowBounds, resultHandler);

        boolean next = false;
        if (CollectionUtils.size(list) > size) {
            next = true;
            list = list.subList(0, size);
        }
        lazyList.setNext(next);
        lazyList.setData(list);

        return List.of(lazyList);
    }

    private MappedStatement getCountMappedStatement(MappedStatement mappedStatement, BoundSql boundSql) {
        String countMappedStatementId = mappedStatement.getId() + "_count";

        String sql = boundSql.getSql();
        String countSql = "SELECT COUNT(*) FROM (" + sql + ") total";

        ResultMap resultMap = null;
        String resultMapId = "java.lang.Integer";
        Configuration configuration = mappedStatement.getConfiguration();
        if (configuration.hasResultMap(resultMapId)) {
            resultMap = configuration.getResultMap(resultMapId);
        } else {
            resultMap = new ResultMap.Builder(configuration, resultMapId, Integer.class, Collections.emptyList()).build();
            configuration.addResultMap(resultMap);
        }

        return createMappedStatement(mappedStatement, countMappedStatementId, countSql, boundSql.getParameterMappings(), resultMap);
    }

    /**
     * @param mappedStatement
     * @param boundSql
     * @param offset          从结果集的第几行开始返回数据（从0开始）
     * @param rows            返回的最大行数
     * @return
     */
    private MappedStatement getListMappedStatement(MappedStatement mappedStatement, Object parameter, BoundSql boundSql, int offset, int rows) {
        String listMappedStatementId = mappedStatement.getId() + "_list";

        // 根据数据库类型生成对应的分页SQL，例如 MySQL 使用 LIMIT，Oracle 使用 ROWNUM 等
        String sql = boundSql.getSql();
        String listSql = sql + " LIMIT ?, ?";

        List<ParameterMapping> newParameterMappings = null;
        List<ParameterMapping> parameterMappings = boundSql.getParameterMappings();
        if (CollectionUtils.isNotEmpty(parameterMappings)) {
            newParameterMappings = new ArrayList<>(parameterMappings.size() + 2);
            newParameterMappings.addAll(parameterMappings);
        } else {
            newParameterMappings = new ArrayList<>(2);
        }
        Configuration configuration = mappedStatement.getConfiguration();
        newParameterMappings.add(new ParameterMapping.Builder(configuration, "offset", Object.class).build());
        newParameterMappings.add(new ParameterMapping.Builder(configuration, "rows", Object.class).build());

        Map parameterMap = (Map) parameter;
        parameterMap.put("offset", offset);
        parameterMap.put("rows", rows);

        return createMappedStatement(mappedStatement, listMappedStatementId, listSql, newParameterMappings, null);
    }

    private MappedStatement createMappedStatement(MappedStatement mappedStatement, String id, String sql, List<ParameterMapping> parameterMappings, ResultMap resultMap) {
        Configuration configuration = mappedStatement.getConfiguration();
        return new MappedStatement.Builder(configuration, id, new StaticSqlSource(configuration, sql, parameterMappings), mappedStatement.getSqlCommandType())
                .resource(mappedStatement.getResource())
                .fetchSize(mappedStatement.getFetchSize())
                .statementType(mappedStatement.getStatementType())
                .keyGenerator(mappedStatement.getKeyGenerator())
//                .keyProperty(mappedStatement.getKeyProperties())
                .timeout(mappedStatement.getTimeout())
                .parameterMap(mappedStatement.getParameterMap())
                .resultMaps(resultMap != null ? List.of(resultMap) : mappedStatement.getResultMaps())
                .resultSetType(mappedStatement.getResultSetType())
                .cache(mappedStatement.getCache())
                .flushCacheRequired(mappedStatement.isFlushCacheRequired())
                .useCache(mappedStatement.isUseCache())
                .build();
    }

}
