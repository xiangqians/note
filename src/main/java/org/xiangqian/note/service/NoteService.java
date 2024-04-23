package org.xiangqian.note.service;

import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.List;

import java.io.IOException;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
public interface NoteService {

    List<NoteEntity> list(NoteEntity vo, List list);

    NoteEntity getById(Long id);

    ModelAndView getViewById(ModelAndView modelAndView, Long id) throws IOException;

    ResponseEntity<Resource> getStreamById(Long id) throws Exception;

    ResponseEntity<Resource> download(Long id) throws IOException;

    Boolean delById(Long id);

    Boolean rename(NoteEntity vo);

    Boolean paste(NoteEntity vo);

    Boolean reUpload(NoteEntity vo) throws IOException;

    Boolean upload(NoteEntity vo) throws IOException;

    Boolean addMd(NoteEntity vo);

    Boolean addFolder(NoteEntity vo);

}
