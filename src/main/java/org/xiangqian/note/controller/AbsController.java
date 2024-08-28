package org.xiangqian.note.controller;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpSession;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.model.Vo;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.ServletUtil;

import java.time.LocalDateTime;
import java.util.Map;

/**
 * @author xiangqian
 * @date 17:27 2024/03/02
 */
public abstract class AbsController {

    private static final String VO = "vo";

    // 在每个请求之前设置ModelAndView值
    @ModelAttribute
    public final void addObject(ModelAndView modelAndView, HttpServletRequest request, HttpSession session) {
        Vo vo = null;
        Object value = session.getAttribute(VO);
        if (value instanceof Vo) {
            vo = (Vo) value;
            session.removeAttribute(VO);
        }

        if (vo == null) {
            vo = Vo.none();
        }

        vo.add("contextPath", request.getContextPath())
                .add("servletPath", request.getServletPath())
                .add("requestURI", request.getRequestURI())
                .add("timestamp", DateUtil.toSecond(LocalDateTime.now()));

        modelAndView.addObject(VO, vo);
    }

    public static Vo add(ModelAndView modelAndView, String name, Object value) {
        Map<String, Object> map = modelAndView.getModel();
        Vo vo = (Vo) map.get(VO);
        vo.add(name, value);
        return vo;
    }

    public static ModelAndView errorView(ModelAndView modelAndView) {
        modelAndView.setViewName("error");
        return modelAndView;
    }

    public static RedirectView redirectView(String url, Vo vo) {
        HttpSession session = ServletUtil.getSession();
        session.setAttribute(VO, vo);
        url += "?t=" + DateUtil.toSecond(LocalDateTime.now());
        return new RedirectView(url);
    }

}
