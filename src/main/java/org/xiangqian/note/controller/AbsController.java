package org.xiangqian.note.controller;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpSession;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.util.DateUtil;

import java.time.LocalDateTime;

/**
 * @author xiangqian
 * @date 17:27 2024/03/02
 */
public abstract class AbsController {

    // 【request域】
    public static final String SERVLET_PATH = "servletPath";
    public static final String TIMESTAMP = "timestamp";

    // 【session域】是否已登陆
    public static final String IS_LOGGEDIN = "isLoggedin";

    // 【session域】
    public static final String USER = "user";

    // 【request域】
    public static final String VO = "vo";
    public static final String ERROR = "error";

    // 在每个请求之前设置ModelAndView值
    @ModelAttribute
    public void modelAttribute(ModelAndView modelAndView, HttpServletRequest request, HttpSession session) {
        modelAndView.addObject(SERVLET_PATH, request.getServletPath());
        modelAndView.addObject(TIMESTAMP, DateUtil.toSecond(LocalDateTime.now()));

        Object vo = session.getAttribute(VO);
        if (vo != null) {
            session.removeAttribute(VO);
            modelAndView.addObject(VO, vo);
        }

        Object error = session.getAttribute(ERROR);
        if (error != null) {
            session.removeAttribute(ERROR);
            modelAndView.addObject(ERROR, error);
        }
    }

    public static boolean getLoggedinAttribute(HttpSession session) {
        Object isLoggedin = session.getAttribute(IS_LOGGEDIN);
        if (isLoggedin instanceof Boolean) {
            return (boolean) isLoggedin;
        }
        return false;
    }

    public static void setLoggedinAttribute(HttpSession session, boolean loggedin) {
        session.setAttribute(IS_LOGGEDIN, loggedin);
    }

    public static UserEntity getUserAttribute(HttpSession session) {
        return (UserEntity) session.getAttribute(USER);
    }

    public static void setUserAttribute(HttpSession session, Object user) {
        session.setAttribute(USER, user);
    }

    protected void setVoAttribute(ModelAndView modelAndView, Object value) {
        modelAndView.addObject(VO, value);
    }

    public static void setVoAttribute(HttpSession session, Object value) {
        session.setAttribute(VO, value);
    }

    protected void setErrorAttribute(ModelAndView modelAndView, Object value) {
        modelAndView.addObject(ERROR, value);
    }

    public static void setErrorAttribute(HttpSession session, Object value) {
        session.setAttribute(ERROR, value);
    }

    protected RedirectView redirectView(String url, Object vo, Object error) {
        HttpSession session = getSession();
        setVoAttribute(session, vo);
        setErrorAttribute(session, error);
        url += (url.contains("?") ? "&" : "?") + "t=" + DateUtil.toSecond(LocalDateTime.now());
        return new RedirectView(url);
    }

    public static HttpSession getSession() {
        return getRequest().getSession();
    }

    public static HttpServletRequest getRequest() {
        return ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes()).getRequest();
    }

}
