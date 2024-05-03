package org.xiangqian.note.service.impl;

import org.apache.commons.lang3.StringUtils;
import org.springframework.core.env.Environment;
import org.springframework.core.io.Resource;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.util.AsposeUtil;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

/**
 * @author xiangqian
 * @date 16:37 2024/05/03
 */
@Service
public class GetDocumentNoteServiceImpl extends AbsGetNoteService {

    protected GetDocumentNoteServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Override
    public ModelAndView getView(ModelAndView modelAndView, NoteEntity entity, List<String> names) throws Exception {
        String type = entity.getType();
        if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX, Type.PDF)) {
            AbsController.setVoAttribute(modelAndView, entity);
            modelAndView.setViewName("note/pdf/view");
            return modelAndView;
        }

        if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
            return getXlsDocumentView(modelAndView, entity);
        }

        return AbsController.errorView(modelAndView);
    }

    @Override
    public ResponseEntity<Resource> getStream(NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path path = getDataPath(id.toString());
        if (!Files.exists(path)) {
            return notFound();
        }

        String type = entity.getType();

        // doc、docx转为pdf
        if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX)) {
            Path tmpPath = getTmpPath(id.toString());
            if (!Files.exists(tmpPath)) {
                AsposeUtil.convertDocToPdf(path.toFile(), tmpPath.toFile());
            }
            return ok(tmpPath, MediaType.APPLICATION_PDF);
        }

        // xls、xlsx转为html
        if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
            return getXlsDocumentStream(getTmpPath(id.toString()), names);
        }

        if (Type.PDF.equals(type)) {
            return ok(path, MediaType.APPLICATION_PDF);
        }

        return notFound();
    }

    @Override
    public boolean isSupported(String type) {
        return Type.isDocument(type);
    }

    private ModelAndView getXlsDocumentView(ModelAndView modelAndView, NoteEntity entity) throws Exception {
        Long id = entity.getId();
        Path htmlPath = convertXlsToHtml(String.format("/note/%s", id), getDataPath(id.toString()), getTmpPath(id.toString()));
        if (htmlPath != null) {
            String content = Files.readString(htmlPath, UTF_8);
            entity.setContent(content);
        }
        AbsController.setVoAttribute(modelAndView, entity);
        modelAndView.setViewName("note/html/view");
        return modelAndView;
    }

    private ResponseEntity<Resource> getXlsDocumentStream(Path path, List<String> names) throws Exception {
        List<Path> paths = getTmpPaths(path, names);
        path = paths.get(paths.size() - 1);
        if (path == null || !Files.exists(path) || Files.isDirectory(path)) {
            return notFound();
        }

        String type = Type.pathOf(path);
        if (Type.HTML.equals(type)) {
            return ok(path, MediaType.TEXT_HTML);
        }

        if (Type.CSS.equals(type)) {
            return ok(path, TEXT_CSS);
        }

        return notFound();
    }

}
