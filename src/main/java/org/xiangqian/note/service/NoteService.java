package org.xiangqian.note.service;

import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.entity.NoteEntity;

import java.io.IOException;
import java.util.List;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
public interface NoteService {

    Boolean updContentById(NoteEntity vo) throws IOException;

    org.xiangqian.note.util.List<NoteEntity> list(NoteEntity vo, Integer offset);

    Object getViewById(ModelAndView modelAndView, Long id, List<String> names) throws Exception;

    ResponseEntity<Resource> getStreamById(Long id, List<String> names) throws Exception;

    ResponseEntity<Resource> download(Long id) throws IOException;

    Boolean delById(Long id);

    Boolean rename(NoteEntity vo);

    Boolean paste(NoteEntity vo);

    Boolean reUpload(NoteEntity vo) throws IOException;

    Boolean upload(NoteEntity vo) throws IOException;

    Boolean addMd(NoteEntity vo);

    Boolean addFolder(NoteEntity vo);

}
