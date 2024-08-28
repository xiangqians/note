package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.BooleanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.entity.UserEntity;
import org.xiangqian.note.model.Vo;
import org.xiangqian.note.service.UserService;

/**
 * @author xiangqian
 * @date 17:09 2024/03/02
 */
@Slf4j
@Controller
public class UserController extends AbsController {

    @Autowired
    private UserService userService;

    @GetMapping("/user/login")
    public ModelAndView login(ModelAndView modelAndView) {
//        add(modelAndView, "info", "信息！信息！信息！信息！")
//                .add("warn", "警告！警告！警告！")
//                .add("error", "错误！错误！错误！错误！");
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
            userService.resetPasswd(vo);
            return redirectView("/", null);
        } catch (Exception e) {
            log.error("", e);
            return redirectView("/user/resetPasswd", Vo.error(e.getMessage()));
        }
    }

}
