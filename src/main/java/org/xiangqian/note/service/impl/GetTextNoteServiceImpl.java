package org.xiangqian.note.service.impl;

import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

/**
 * @author xiangqian
 * @date 16:31 2024/05/03
 */
@Service
public class GetTextNoteServiceImpl extends AbsGetNoteService {

    protected GetTextNoteServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Override
    public ModelAndView getView(ModelAndView modelAndView, NoteEntity entity, List<String> names) throws Exception {
        String type = entity.getType();
        String name = switch (type) {
            case Type.MD -> "md";
            case Type.HTML -> "html";
            default -> null;
        };

        if (name == null) {
            return AbsController.errorView(modelAndView);
        }

        return view(modelAndView, entity, name);
    }

    @Override
    public ResponseEntity<Resource> getStream(NoteEntity entity, List<String> names) throws Exception {
        return notFound();
    }

    @Override
    public boolean isSupported(String type) {
        return Type.isText(type);
    }

    private ModelAndView view(ModelAndView modelAndView, NoteEntity entity, String name) throws IOException {
        Long id = entity.getId();
        Path path = getPath(id.toString());
        if (Files.exists(path)) {
            String content = Files.readString(path, UTF_8);
            entity.setContent(content);
        }

        AbsController.setVoAttribute(modelAndView, entity);
        modelAndView.setViewName(String.format("note/%s/view", name));
        return modelAndView;
    }

}
