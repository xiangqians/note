package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.BooleanUtils;
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
import org.xiangqian.note.util.*;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

    @PostMapping("/folder")
    public RedirectView addFolder(NoteEntity entity) {
        Object message = null;
        try {
            service.createFolder(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @PostMapping("/md")
    public RedirectView addMd(NoteEntity entity) {
        Object message = null;
        try {
            service.createMd(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @GetMapping("/{pid}/list")
    public ModelAndView list(ModelAndView modelAndView, @PathVariable("pid") Long pid, NoteEntity entity) {
        try {
            // 当前页
            if (entity.getCurrent() == null) {
                entity.setCurrent(1);
            } else if (entity.getCurrent().intValue() <= 0) {
                return errorView(modelAndView);
            }

            entity.setPid(pid);
            entity.setName(StringUtils.trimToNull(entity.getName()));
            entity.setContent(StringUtils.trimToNull(entity.getContent()));
            entity.setType(StringUtils.trimToNull(entity.getType()));

            NoteEntity parentEntity = service.getById(pid);
            if (parentEntity == null || !Type.FOLDER.equals(parentEntity.getType())) {
                return errorView(modelAndView);
            }

            modelAndView.addObject("entity", entity);
            modelAndView.addObject("parentEntity", parentEntity);
            modelAndView.addObject("parentList", service.getParentListById(pid));
            modelAndView.addObject("childList", service.getChildList(entity));
            modelAndView.addObject("types", List.of(Type.FOLDER, Type.MD, Type.PDF, Type.ZIP));
        } catch (Exception e) {
            log.error("", e);
            modelAndView.addObject(MESSAGE, e.getMessage());
        }
        modelAndView.setViewName("note/list");
        return modelAndView;
    }

    @GetMapping("/{id}/sort")
    public ModelAndView getSortById(ModelAndView modelAndView, @PathVariable("id") Long id) {
        try {
            NoteEntity entity = service.getById(id);
            if (entity == null || !Type.FOLDER.equals(entity.getType())) {
                return errorView(modelAndView);
            }
            modelAndView.addObject("id", id);
            modelAndView.addObject("name", entity.getName());
            modelAndView.addObject("content", service.getContentById(id));
        } catch (Exception e) {
            log.error("", e);
            modelAndView.addObject(MESSAGE, e.getMessage());
        }
        modelAndView.setViewName("note/sort");
        return modelAndView;
    }

    @ResponseBody
    @PutMapping("/{id}/sort")
    public Map<String, Object> updateSortById(@PathVariable(name = "id") Long id, @RequestBody(required = false) String content) throws Exception {
        try {
            return Model.of("result", service.updateSortById(id, content));
        } catch (Exception e) {
            log.error("", e);
            return Model.of(MESSAGE, e.getMessage());
        }
    }

    @PostMapping("/upload")
    public RedirectView upload(NoteEntity entity) {
        Object message = null;
        try {
            service.upload(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @PutMapping("/{id}/reUpload")
    public RedirectView reUpload(@PathVariable("id") Long id, NoteEntity entity) {
        Object message = null;
        try {
            entity.setId(id);
            service.upload(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @DeleteMapping("/{id}")
    public RedirectView deleteById(@PathVariable("id") Long id, NoteEntity entity) {
        Object message = null;
        try {
            service.deleteById(id);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @PutMapping("/{id}/rename")
    public RedirectView rename(@PathVariable("id") Long id, NoteEntity entity) {
        Object message = null;
        try {
            entity.setId(id);
            service.rename(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    @PutMapping("/paste")
    public RedirectView paste(NoteEntity entity) {
        Object message = null;
        try {
            service.paste(entity);
        } catch (Exception e) {
            log.error("", e);
            message = e.getMessage();
        }
        return redirectListView(entity.getPid(), message);
    }

    private final Pattern PATTERN = Pattern.compile("([^/]+/?)");

    @GetMapping("/{id}/view")
    public ModelAndView getViewById(ModelAndView modelAndView, @PathVariable("id") Long id,
                                    @RequestParam(required = false) String name,
                                    @RequestParam(required = false) Boolean special) throws IOException {
        NoteEntity entity = service.getById(id);
        if (entity == null || Type.FOLDER.equals(entity.getType())) {
            return errorView(modelAndView);
        }

        if (Type.MD.equals(entity.getType())) {
            modelAndView.addObject("id", entity.getId());
            modelAndView.addObject("name", String.format("%s.%s", entity.getName(), entity.getType()));
            modelAndView.addObject("content", service.getContentById(id));
            modelAndView.setViewName("note/md/view");
            return modelAndView;
        }

        if (Type.PDF.equals(entity.getType())) {
            modelAndView.addObject("name", String.format("%s.%s", entity.getName(), entity.getType()));
            modelAndView.addObject("url", String.format("/note/%s?t=%s", entity.getId(), TimeUtil.now()));
            modelAndView.setViewName("note/pdf/view");
            return modelAndView;
        }

        if (Type.ZIP.equals(entity.getType())) {
            if (name == null) {
                name = "";
            }

            String type = FileUtil.getType(name);
            if (Type.PDF.equals(type) && BooleanUtils.isTrue(special)) {
                String name1 = null;
                int index = name.lastIndexOf("/");
                if (index >= 0) {
                    name1 = name.substring(index + 1);
                } else {
                    name1 = name;
                }
                modelAndView.addObject("name", name1);
                modelAndView.addObject("url", String.format("/note/%s?name=%s&t=%s", id, name, TimeUtil.now()));
                modelAndView.setViewName("note/pdf/view");
                return modelAndView;
            }

            ZipEntry zipEntry = new ZipEntry();
            zipEntry.setPath(name);
            zipEntry.setType(type);

            if ("".equals(name) || Type.FOLDER.equals(type)) {
                zipEntry.setType(Type.FOLDER);
                List<ZipEntry> childList = FileUtil.getZipEntryList(service.getPath(id), name);
                zipEntry.setChildList(childList);
            } else {
                byte[] content = FileUtil.getZipEntryContent(service.getPath(id), name);
                zipEntry.setContent(new String(content, StandardCharsets.UTF_8));
            }

            if (!"".equals(name)) {
                List<ZipEntry> parentList = new ArrayList<>();
                Matcher matcher = PATTERN.matcher(name);
                while (matcher.find()) {
                    ZipEntry parentZipEntry = new ZipEntry();

                    String fullName = matcher.group(1);
                    type = FileUtil.getType(fullName);
                    String name1 = null;
                    if (fullName.endsWith("/")) {
                        name1 = fullName.substring(0, fullName.length() - 1);
                    } else {
                        name1 = fullName;
                    }

                    String path = null;
                    if (parentList.size() > 0) {
                        path = parentList.get(parentList.size() - 1).getPath() + fullName;
                    } else {
                        path = fullName;
                    }

                    parentZipEntry.setName(name1);
                    parentZipEntry.setType(type);
                    parentZipEntry.setPath(path);
                    parentList.add(parentZipEntry);
                }
                zipEntry.setParentList(parentList);

                zipEntry.setName(parentList.get(parentList.size() - 1).getName());
            }

            modelAndView.addObject("id", entity.getId());
            modelAndView.addObject("name", String.format("%s.%s", entity.getName(), entity.getType()));
            modelAndView.addObject("zipEntry", zipEntry);
            modelAndView.setViewName("note/zip/view");
            return modelAndView;
        }

        return errorView(modelAndView);
    }

    @ResponseBody
    @PutMapping("/{id}/content")
    public Map<String, Object> updateContentById(@PathVariable(name = "id") Long id, @RequestBody(required = false) String content) throws Exception {
        try {
            return Model.of("result", service.updateContentById(id, content));
        } catch (Exception e) {
            log.error("", e);
            return Model.of(MESSAGE, e.getMessage());
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<Resource> getResourceById(@PathVariable("id") Long id, @RequestParam(required = false) String name) throws IOException {
        return responseEntity("inline", service.getResourceById(id, name));
    }

    @GetMapping("/{id}/download")
    public ResponseEntity<Resource> download(@PathVariable("id") Long id, @RequestParam(required = false) String name) throws IOException {
        return responseEntity("attachment", service.getResourceById(id, name));
    }

    private RedirectView redirectListView(Long id, Object message) {
        return redirectView(String.format("/note/%s/list", id), Model.of(MESSAGE, message));
    }

}
