package org.xiangqian.note.controller;

import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.Response;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

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

    @ResponseBody
    @PutMapping("/{id}/content")
    public Response<Boolean> updContentById(@PathVariable(name = "id") Long id, NoteEntity vo) throws Exception {
        try {
            vo.setId(id);
            return Response.ok(service.updContentById(vo));
        } catch (Exception e) {
            log.error("", e);
            return Response.error(e.getMessage());
        }
    }

    @GetMapping("/{id}/list")
    public ModelAndView list(ModelAndView modelAndView, @PathVariable("id") Long id, NoteEntity vo, @RequestParam(value = "offset", required = false) Integer offset) {
        try {
            NoteEntity entity = service.getById(id, true);
            if (entity == null || !Type.FOLDER.equals(entity.getType())) {
                return errorView(modelAndView);
            }

            vo.setPid(id);
            entity.setChildList(service.list(vo, offset));
            vo.setId(id);
            vo.setPid(null);
            setVoAttribute(modelAndView, Map.of("parameter", vo,
                    "entity", entity,
                    "types", java.util.List.of(Type.FOLDER, Type.MD, Type.DOC, Type.DOCX, Type.PDF, Type.HTML, Type.ZIP)));
        } catch (Exception e) {
            log.error("", e);
            setErrorAttribute(modelAndView, e.getMessage());
        }
        modelAndView.setViewName("note/list");
        return modelAndView;
    }

    @GetMapping("/{id}/**/view")
    public ModelAndView getViewById(HttpServletRequest request, ModelAndView modelAndView, @PathVariable(name = "id") Long id) throws Exception {
        List<String> names = getNames(request, String.format("/note/%s", id), "/view");
        return service.getViewById(modelAndView, id, names);
    }

    @GetMapping(value = "/{id}/**/stream")
    public ResponseEntity<Resource> getStreamById(HttpServletRequest request, @PathVariable("id") Long id) throws Exception {
        List<String> names = getNames(request, String.format("/note/%s", id), "/stream");
        return service.getStreamById(id, names);
    }

    @GetMapping("/{id}/download")
    public ResponseEntity<Resource> download(@PathVariable("id") Long id) throws IOException {
        return service.download(id);
    }

    @DeleteMapping("/del")
    public RedirectView delById(NoteEntity vo) {
        Object error = null;
        try {
            service.delById(vo.getId());
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PutMapping("/rename")
    public RedirectView rename(NoteEntity vo) {
        Object error = null;
        try {
            service.rename(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PutMapping("/paste")
    public RedirectView paste(NoteEntity vo) {
        Object error = null;
        try {
            service.paste(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PutMapping("/reUpload")
    public RedirectView reUpload(NoteEntity vo) {
        Object error = null;
        try {
            service.reUpload(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PostMapping("/upload")
    public RedirectView upload(NoteEntity vo) {
        Object error = null;
        try {
            service.upload(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PostMapping("/addMd")
    public RedirectView addMd(NoteEntity vo) {
        Object error = null;
        try {
            service.addMd(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @PostMapping("/addFolder")
    public RedirectView addFolder(NoteEntity vo) {
        Object error = null;
        try {
            service.addFolder(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    private RedirectView redirectListView(Long id, Object vo, Object error) {
        return redirectView(String.format("/note/%s/list", id), vo, error);
    }

    private List<String> getNames(HttpServletRequest request, String prefix, String suffix) {
        String path = request.getServletPath();
        path = path.substring(prefix.length(), path.length() - suffix.length());
        if (StringUtils.isEmpty(path)) {
            return null;
        }
        return Arrays.stream(path.split("/")).filter(StringUtils::isNotEmpty).collect(Collectors.toList());
    }

}
