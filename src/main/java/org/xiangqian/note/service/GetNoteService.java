package org.xiangqian.note.service;

import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.entity.NoteEntity;

import java.util.List;

/**
 * @author xiangqian
 * @date 16:27 2024/05/03
 */
public interface GetNoteService {

    ModelAndView getView(ModelAndView modelAndView, NoteEntity entity, List<String> names) throws Exception;

    ResponseEntity<Resource> getStream(NoteEntity entity, List<String> names) throws Exception;

    boolean isSupported(String type);

}
