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

    @GetMapping("/user/resetPasswd")
    public ModelAndView resetPasswd(ModelAndView modelAndView) {
        modelAndView.setViewName("user/resetPasswd");
        return modelAndView;
    }

    @PutMapping("/user/resetPasswd")
    public RedirectView resetPasswd(UserEntity vo) {
        try {
            service.resetPasswd(vo);
            return redirectView("/", null, null);
        } catch (Exception e) {
            log.error("", e);
            return redirectView("/user/resetPasswd", null, e.getMessage());
        }
    }

}
