package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.List;

/**
 * @author xiangqian
 * @date 18:30 2024/03/03
 */
@Slf4j
@Controller
@RequestMapping("/")
public class IndexController extends AbsController {

    @Autowired
    private NoteController noteController;

    @RequestMapping
    public ModelAndView index(ModelAndView modelAndView) {
        return noteController.list(modelAndView, 0L, new NoteEntity(), new List());
    }

}
