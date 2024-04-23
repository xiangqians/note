package org.xiangqian.note.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.xiangqian.note.entity.IavEntity;
import org.xiangqian.note.service.IavService;

import java.io.IOException;

/**
 * @author xiangqian
 * @date 20:24 2024/04/23
 */
@Slf4j
@Controller
@RequestMapping("/iav")
public class IavController {

    @Autowired
    private IavService service;

    @GetMapping("/{id}/stream")
    public ResponseEntity<Resource> getStreamById(@PathVariable("id") Long id) throws IOException {
        return service.getStreamById(id);
    }

    @ResponseBody
    @PostMapping("/upload")
    public IavEntity upload(@RequestParam("file") MultipartFile file) throws IOException {
        return service.upload(file);
    }

}
