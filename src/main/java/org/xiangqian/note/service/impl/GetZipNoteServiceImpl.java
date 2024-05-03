package org.xiangqian.note.service.impl;

import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.collections4.MapUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.io.file.PathUtils;
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
import org.xiangqian.note.util.Md5Util;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.attribute.FileTime;
import java.util.*;
import java.util.stream.Collectors;
import java.util.zip.ZipEntry;
import java.util.zip.ZipException;
import java.util.zip.ZipFile;

/**
 * @author xiangqian
 * @date 16:39 2024/05/03
 */
@Service
public class GetZipNoteServiceImpl extends AbsGetNoteService {

    protected GetZipNoteServiceImpl(Environment environment) throws IOException {
        super(environment);
    }

    @Override
    public ModelAndView getView(ModelAndView modelAndView, NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path dataPath = getDataPath(id.toString());
        if (!Files.exists(dataPath)) {
            return view(modelAndView, entity, null);
        }

        Path tmpPath = getTmpPath(id.toString());
        unzipIfNotDecompressed(dataPath, tmpPath);

        List<Path> tmpPaths = getTmpPaths(tmpPath, names);
        tmpPath = tmpPaths.get(tmpPaths.size() - 1);
        if (tmpPath == null) {
            return view(modelAndView, entity, null);
        }

        List<NoteEntity> ps = null;
        int size = tmpPaths.size();
        if (size > 1) {
            ps = new ArrayList<>(size);
            int index = 1;
            while (index < size) {
                Path path1 = tmpPaths.get(index++);
                if (path1 == null) {
                    break;
                }
                ps.add(new NoteEntity(path1));
            }
        }
        entity.setPs(ps);

        if (Files.isDirectory(tmpPath)) {
            List<NoteEntity> childList = new ArrayList<>(16);

            // 获取目录下的子文件夹
            Iterator<Path> iterator = Files.list(tmpPath).iterator();
            while (iterator.hasNext()) {
                tmpPath = iterator.next();
                childList.add(new NoteEntity(tmpPath));
            }

            Collections.sort(childList);

            entity.setChildList(new org.xiangqian.note.util.List<>(childList));

            return view(modelAndView, entity, null);
        }

        entity = new NoteEntity(tmpPath);
        entity.setId(id);
        entity.setPs(ps);
        String type = entity.getType();
        if (Type.isText(type)) {
            // 读取文件内容
            String content = Files.readString(tmpPath, UTF_8);
            entity.setContent(content);
            return view(modelAndView, entity, "text");
        }

        if (Type.isDocument(type)) {
            if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX, Type.PDF)) {
                return view(modelAndView, entity, "pdf");
            }

            if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
                Path newTmpPath = getTmpPath(id + "_" + Md5Util.encryptHex(tmpPaths.stream().map(Path::getFileName).map(Path::toString).collect(Collectors.joining("/"))));
                Path htmlPath = convertXlsToHtml(String.format("/note/%s/%s", id,
                                // subList[fromIndex, toIndex)
                                tmpPaths.subList(1, tmpPaths.size()).stream().map(Path::getFileName).map(Path::toString).collect(Collectors.joining("/"))),
                        tmpPath, newTmpPath);
                if (htmlPath != null) {
                    // 读取文件内容
                    String content = Files.readString(htmlPath, UTF_8);
                    entity.setContent(content);
                }
                return view(modelAndView, entity, "html");
            }

            return view(modelAndView, entity, "unsupported");
        }

        if (Type.isImage(type)) {
            return view(modelAndView, entity, "image");
        }

        return view(modelAndView, entity, "unsupported");
    }

    @Override
    public ResponseEntity<Resource> getStream(NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path dataPath = getDataPath(id.toString());
        if (!Files.exists(dataPath)) {
            return notFound();
        }

        List<Path> tmpPaths = null;
        Map<String, Integer> tryMap = new LinkedHashMap<>(CollectionUtils.size(names) + 1, 1f);
        tryMap.put(id.toString(), 0);
        if (CollectionUtils.isNotEmpty(names)) {
            int index = names.size();
            while (index > 0) {
                tmpPaths = getTmpPaths(getTmpPath(id.toString()),
                        // subList[fromIndex, toIndex)
                        names.subList(0, index));
                if (tmpPaths.get(tmpPaths.size() - 1) != null) {
                    tryMap.put(id + "_" + Md5Util.encryptHex(tmpPaths.stream().map(Path::getFileName).map(Path::toString).collect(Collectors.joining("/"))), index);
                }
                index--;
            }
        }

        for (Map.Entry<String, Integer> entry : tryMap.entrySet()) {
            tmpPaths = getTmpPaths(getTmpPath(entry.getKey()), names.subList(entry.getValue(), names.size()));
            if (tmpPaths != null && tmpPaths.get(tmpPaths.size() - 1) != null) {
                break;
            }
        }

        if (tmpPaths == null) {
            return notFound();
        }


        Path tmpPath = tmpPaths.get(tmpPaths.size() - 1);
        if (tmpPath == null || Files.isDirectory(tmpPath)) {
            return notFound();
        }

        entity = new NoteEntity(tmpPath);
        String type = entity.getType();
        if (type == null) {
            return notFound();
        }

        MediaType contentType = null;
        switch (type) {
            case Type.PNG -> contentType = MediaType.IMAGE_PNG;
            case Type.JPG -> contentType = MediaType.IMAGE_JPEG;
            case Type.GIF -> contentType = MediaType.IMAGE_GIF;
            case Type.WEBP -> contentType = IMAGE_WEBP;
            case Type.ICO -> contentType = IMAGE_X_ICON;
            case Type.DOC, Type.DOCX -> {
                Path newTmpPath = getTmpPath(id + "_" + Md5Util.encryptHex(tmpPaths.stream().map(Path::getFileName).map(Path::toString).collect(Collectors.joining("/"))));
                if (!Files.exists(newTmpPath)) {
                    AsposeUtil.convertDocToPdf(tmpPath.toFile(), newTmpPath.toFile());
                }
                tmpPath = newTmpPath;
                contentType = MediaType.APPLICATION_PDF;
            }
            case Type.HTML -> contentType = MediaType.TEXT_HTML;
            case Type.CSS -> contentType = TEXT_CSS;
            case Type.PDF -> contentType = MediaType.APPLICATION_PDF;
        }

        if (contentType == null) {
            return notFound();
        }

        return ok(tmpPath, contentType);
    }

    @Override
    public boolean isSupported(String type) {
        return Type.ZIP.equals(type);
    }

    private ModelAndView view(ModelAndView modelAndView, NoteEntity entity, String name) {
        AbsController.setVoAttribute(modelAndView, entity);

        String view = null;
        if (name != null) {
            view = String.format("note/zip/%s/view", name);
        } else {
            view = "note/zip/view";
        }
        modelAndView.setViewName(view);

        return modelAndView;
    }

    /**
     * 如果未解压，则解压到临时目录
     *
     * @param zipPath zip文件
     * @param dirPath 解压到指定目录
     * @return
     * @throws IOException
     */
    private void unzipIfNotDecompressed(Path zipPath, Path dirPath) throws IOException {
        ZipFile zip = null;
        Files.createDirectories(dirPath);
        try {
            try {
                zip = new ZipFile(zipPath.toFile(), UTF_8);
            } catch (ZipException e) {
                String message = e.getMessage();
                // java.util.zip.ZipException: invalid CEN header (bad entry name)
                if (StringUtils.containsIgnoreCase(message, "invalid")
                        && StringUtils.containsIgnoreCase(message, "CEN")
                        && StringUtils.containsIgnoreCase(message, "header")) {
                    zip = new ZipFile(zipPath.toFile(), GBK);
                } else {
                    throw e;
                }
            }

            Enumeration<? extends ZipEntry> entries = zip.entries();
            // <Path, LastModified>
            Map<Path, FileTime> pathLastModifiedMap = new HashMap<>(16, 1f);
            while (entries.hasMoreElements()) {
                ZipEntry entry = entries.nextElement();
                Path entryPath = Paths.get(dirPath.toString(), entry.getName());
                if (entry.isDirectory()) {
                    Files.createDirectories(entryPath);
                    pathLastModifiedMap.put(entryPath, entry.getLastModifiedTime());
                } else {
                    OutputStream outputStream = null;
                    try {
                        outputStream = Files.newOutputStream(entryPath);
                        zip.getInputStream(entry).transferTo(outputStream);
                        Files.setLastModifiedTime(entryPath, entry.getLastModifiedTime());
                    } finally {
                        IOUtils.closeQuietly(outputStream);
                    }
                }
            }

            if (MapUtils.isNotEmpty(pathLastModifiedMap)) {
                for (Map.Entry<Path, FileTime> entry : pathLastModifiedMap.entrySet()) {
                    Files.setLastModifiedTime(entry.getKey(), entry.getValue());
                }
            }
        } catch (Exception e) {
            IOUtils.closeQuietly(zip);
            zip = null;

            // 在删除目录前释放文件句柄（释放 FileInputStream、FileOutputStream、RandomAccessFile 等流），才会立即删除文件，否则会延迟删除文件
            PathUtils.deleteDirectory(dirPath);

            throw e;
        } finally {
            IOUtils.closeQuietly(zip);
        }
    }

}
