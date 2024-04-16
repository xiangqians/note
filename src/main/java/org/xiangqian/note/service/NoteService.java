package org.xiangqian.note.service;

import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.List;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
public interface NoteService {

    Boolean addMdFile(NoteEntity vo);

    Boolean addFolder(NoteEntity vo);

    List<NoteEntity> list(NoteEntity vo, List list);

}
