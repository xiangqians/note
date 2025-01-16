package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.service.UserService;
import org.xiangqian.note.util.Model;

/**
 * @author xiangqian
 * @date 17:09 2024/03/02
 */
@Slf4j
@Controller
public class UserController extends AbsController {

    @Autowired
    private UserService service;

    @GetMapping("/login")
    public ModelAndView login(ModelAndView modelAndView) {
        modelAndView.setViewName("user/login");
        return modelAndView;
    }

    @GetMapping("/resetPassword")
    public ModelAndView resetPassword(ModelAndView modelAndView) {
        modelAndView.setViewName("user/reset-password");
        return modelAndView;
    }

    @PutMapping("/resetPassword")
    public RedirectView resetPassword(UserEntity entity) {
        try {
            if (service.resetPassword(entity)) {
                return redirectView("/resetPassword", null);
            }
            return redirectView("/resetPassword", Model.of(MESSAGE, "重置密码失败"));
        } catch (Exception e) {
            log.error("", e);
            return redirectView("/resetPassword", Model.of(MESSAGE, e.getMessage()));
        }
    }

}
