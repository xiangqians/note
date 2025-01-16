package org.xiangqian.note.controller;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpSession;
import org.springframework.core.io.Resource;
import org.springframework.http.ContentDisposition;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.util.ResourceImpl;
import org.xiangqian.note.util.TimeUtil;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * @author xiangqian
 * @date 17:27 2024/03/02
 */
public abstract class AbsController {

    public static final String MODEL = "model";
    public static final String MESSAGE = "message";

    /**
     * 在每个请求之前设置ModelAndView值
     *
     * @param modelAndView
     * @param request
     * @param session
     */
    @ModelAttribute
    public final void modelAttribute(ModelAndView modelAndView, HttpServletRequest request, HttpSession session) {
        Map<String, Object> model = modelAndView.getModel();
        Object value = session.getAttribute(MODEL);
        if (value instanceof Map) {
            model.putAll((Map) value);
            session.removeAttribute(MODEL);
        }

        // 上下文路径
        modelAndView.addObject("contextPath", request.getContextPath());
        // servlet路径
        modelAndView.addObject("servletPath", request.getServletPath());
        // 上下文路径 + servlet路径
        modelAndView.addObject("requestURI", request.getRequestURI());
        // 时间戳
        modelAndView.addObject("timestamp", TimeUtil.now());
    }

    public static ModelAndView errorView(ModelAndView modelAndView) {
        modelAndView.setViewName("error");
        return modelAndView;
    }

    public static RedirectView redirectView(String url, Map<String, Object> model) {
        HttpServletRequest request = ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes()).getRequest();
        HttpSession session = request.getSession();
        session.setAttribute(MODEL, model);
        url += "?t=" + TimeUtil.now();
        return new RedirectView(url);
    }

    /**
     * @param type     设置Content-Disposition，指定文件以附件方式下载。
     *                 "attachment" 表示文件将作为附件下载。
     *                 还可以使用 "inline" 来指定浏览器直接展示文件（例如PDF文件直接在浏览器中显示）。
     * @param resource
     * @return
     * @throws IOException
     */
    public static ResponseEntity<Resource> responseEntity(String type, Resource resource) throws IOException {
        if (resource == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        }

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        // 设置Content-Disposition，指定文件以附件方式下载，并且给文件一个名字。
        headers.setContentDisposition(ContentDisposition.builder(type)
                .filename(URLEncoder.encode(((ResourceImpl) resource).getName(), StandardCharsets.UTF_8))
                .build());
        // 设置 Content-Type
        headers.setContentType(((ResourceImpl) resource).getType());
        // 设置 Content-Length
        headers.setContentLength(resource.contentLength());

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

}
