package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.ImageEntity;
import org.xiangqian.note.service.ImageService;
import org.xiangqian.note.util.Model;

import java.io.IOException;
import java.util.Map;

/**
 * @author xiangqian
 * @date 20:24 2024/04/23
 */
@Slf4j
@Controller
@RequestMapping("/image")
public class ImageController extends AbsController {

    @Autowired
    private ImageService service;

    @ResponseBody
    @PostMapping("/upload")
    public Map<String, Object> upload(@RequestParam("file") MultipartFile file) {
        try {
            ImageEntity entity = service.upload(file);
            return Model.of("id", entity.getId(),
                    "name", entity.getName(),
                    "type", entity.getType());
        } catch (Exception e) {
            log.error("", e);
            return Model.of(MESSAGE, e.getMessage());
        }
    }

    @GetMapping("/{id}")
    public ResponseEntity<Resource> getResourceById(@PathVariable("id") Long id) throws IOException {
        return responseEntity("inline", service.getResourceById(id));
    }

}
