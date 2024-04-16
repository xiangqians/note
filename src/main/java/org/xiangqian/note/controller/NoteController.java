package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.List;
import org.xiangqian.note.util.Type;

import java.util.Map;

/**
 * @author xiangqian
 * @date 21:54 2024/02/29
 */
@Slf4j
@Controller
@RequestMapping("/note")
public class NoteController extends AbsController {

    @Autowired
    private NoteService service;

    @PostMapping("/{pid}/addMdFile")
    public RedirectView addMdFile(@PathVariable("pid") Long pid, NoteEntity vo) {
        Object error = null;
        try {
            vo.setPid(pid);
            service.addMdFile(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(pid, null, error);
    }

    @PostMapping("/{pid}/addFolder")
    public RedirectView addFolder(@PathVariable("pid") Long pid, NoteEntity vo) {
        Object error = null;
        try {
            vo.setPid(pid);
            service.addFolder(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(pid, null, error);
    }

    @GetMapping("/{pid}/list")
    public ModelAndView list(ModelAndView modelAndView, @PathVariable("pid") Long pid, NoteEntity vo, List list) {
        try {
            NoteEntity p = service.getById(pid);
            if (p == null || !Type.FOLDER.equals(p.getType())) {
                return errorView(modelAndView);
            }

            vo.setPid(pid);
            list = service.list(vo, list);
            setVoAttribute(modelAndView, Map.of("parameter", vo,
                    "types", Type.getSet(),
                    "p", p,
                    "offset", list.getOffset(),
                    "data", list.getData(),
                    "offsets", list.getOffsets()));
        } catch (Exception e) {
            log.error("", e);
            setErrorAttribute(modelAndView, e.getMessage());
        }
        modelAndView.setViewName("note/list");
        return modelAndView;
    }

    private RedirectView redirectListView(Long pid, Object vo, Object error) {
        return redirectView(String.format("/note/%s/list", pid), vo, error);
    }

}
