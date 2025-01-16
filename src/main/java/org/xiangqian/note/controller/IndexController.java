package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.view.RedirectView;

/**
 * @author xiangqian
 * @date 18:30 2024/03/03
 */
@Slf4j
@Controller
@RequestMapping("/")
public class IndexController extends AbsController {

    @RequestMapping
    public RedirectView index() {
        return redirectView("/note/0/list", null);
    }

}
