package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
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

    @PutMapping("/{id}/rename")
    public RedirectView rename(@PathVariable("id") Long id, NoteEntity vo) {
        Object error = null;
        try {
            vo.setId(id);
            service.rename(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }

    @GetMapping("/{id}")
    public RedirectView get(@PathVariable("pid") Long pid, NoteEntity vo) {
        Object error = null;
        try {
            vo.setPid(pid);
            service.rename(vo);
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(pid, null, error);
    }

    @PostMapping("/uploadFile")
    public RedirectView uploadFile(@RequestParam("file") MultipartFile file, NoteEntity vo) {

        // private static String UPLOAD_FOLDER = "C:/uploads/"; // 上传文件存储的目录
        //
        //    @PostMapping("/upload")
        //    public String uploadFile() {
        //        if (file.isEmpty()) {
        //            return "请选择文件上传";
        //        }
        //
        //        try {
        //            // 获取文件的字节数组并保存到指定路径
        //            byte[] bytes = file.getBytes();
        //            Path path = Paths.get(UPLOAD_FOLDER + file.getOriginalFilename());
        //            Files.write(path, bytes);
        //            return "文件上传成功";
        //        } catch (IOException e) {
        //            e.printStackTrace();
        //            return "文件上传失败";
        //        }
        //    }
        //}

        Object error = null;
        try {
        } catch (Exception e) {
            log.error("", e);
            error = e.getMessage();
        }
        return redirectListView(vo.getPid(), null, error);
    }


    @PostMapping("/addMdFile")
    public RedirectView addMdFile(NoteEntity vo) {
        Object error = null;
        try {
            service.addMdFile(vo);
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

    @GetMapping("/{id}/list")
    public ModelAndView list(ModelAndView modelAndView, @PathVariable("id") Long id, NoteEntity vo, List list) {
        try {
            NoteEntity entity = service.getById(id);
            if (entity == null || !Type.FOLDER.equals(entity.getType())) {
                return errorView(modelAndView);
            }

            vo.setPid(id);
            list = service.list(vo, list);
            vo.setId(id);
            vo.setPid(null);
            setVoAttribute(modelAndView, Map.of("parameter", vo,
                    "entity", entity,
                    "data", list.getData(),
                    "types", Type.getSet(),
                    "offset", list.getOffset(),
                    "offsets", list.getOffsets()));
        } catch (Exception e) {
            log.error("", e);
            setErrorAttribute(modelAndView, e.getMessage());
        }
        modelAndView.setViewName("note/list");
        return modelAndView;
    }

    private RedirectView redirectListView(Long id, Object vo, Object error) {
        return redirectView(String.format("/note/%s/list", id), vo, error);
    }

}
