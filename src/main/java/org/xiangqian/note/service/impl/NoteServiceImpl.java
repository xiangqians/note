package org.xiangqian.note.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.collections4.MapUtils;
import org.apache.commons.io.IOUtils;
import org.apache.commons.io.file.PathUtils;
import org.apache.commons.lang3.BooleanUtils;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.util.PropertyPlaceholderHelper;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.ModelAndView;
import org.xiangqian.note.controller.AbsController;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.mapper.NoteMapper;
import org.xiangqian.note.service.NoteService;
import org.xiangqian.note.util.AsposeUtil;
import org.xiangqian.note.util.DateUtil;
import org.xiangqian.note.util.Md5Util;
import org.xiangqian.note.util.Type;

import java.io.IOException;
import java.io.OutputStream;
import java.net.URLEncoder;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;
import java.nio.file.attribute.FileTime;
import java.time.LocalDateTime;
import java.util.*;
import java.util.stream.Collectors;
import java.util.zip.ZipEntry;
import java.util.zip.ZipException;
import java.util.zip.ZipFile;

/**
 * @author xiangqian
 * @date 21:35 2024/03/04
 */
@Slf4j
@Service
public class NoteServiceImpl extends AbsService implements NoteService {

    @Autowired
    private NoteMapper mapper;

    @Override
    public Boolean updContentById(NoteEntity vo) throws IOException {
        Long id = vo.getId();
        NoteEntity entity = getById(id);
        if (entity == null) {
            return false;
        }

        Path path = getPath(id.toString());
        String content = vo.getContent();
        if (content == null) {
            content = "";
        }
        Files.write(path, content.getBytes(UTF_8), StandardOpenOption.TRUNCATE_EXISTING);

        long size = Files.size(path);

        NoteEntity updEntity = new NoteEntity();
        updEntity.setId(id);
        updEntity.setSize(size);
        updEntity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.updateById(updEntity) > 0;
    }

    @Override
    public org.xiangqian.note.util.List<NoteEntity> list(NoteEntity vo, Integer offset) {
        if (offset == null) {
            offset = 1;
        }
        Long pid = vo.getPid();
        org.xiangqian.note.util.List list = new org.xiangqian.note.util.List(offset);
        if (offset.intValue() <= 0 || pid == null || pid.longValue() < 0) {
            return list;
        }

        vo.setName(StringUtils.trimToNull(vo.getName()));
        vo.setType(StringUtils.trimToNull(vo.getType()));
        list = mapper.list(vo, list);
        if (BooleanUtils.isTrue(vo.getContain())) {
            if (CollectionUtils.isNotEmpty(list)) {
                Set<Long> pids = ((List<NoteEntity>) list).stream().map(NoteEntity::getPs)
                        .filter(CollectionUtils::isNotEmpty)
                        .flatMap(Collection::stream)
                        .map(NoteEntity::getId)
                        .collect(Collectors.toSet());
                if (CollectionUtils.isNotEmpty(pids)) {
                    // <pid, name>
                    Map<Long, String> pidMap = mapper.selectList(new LambdaQueryWrapper<NoteEntity>()
                                    .select(NoteEntity::getId, NoteEntity::getName)
                                    .in(NoteEntity::getId, pids))
                            .stream().collect(Collectors.toMap(NoteEntity::getId, NoteEntity::getName));
                    for (NoteEntity entity : (List<NoteEntity>) list) {
                        List<NoteEntity> ps = entity.getPs();
                        if (CollectionUtils.isNotEmpty(ps)) {
                            for (NoteEntity p : ps) {
                                p.setName(pidMap.get(p.getId()));
                            }
                        }
                    }
                }
            }
        }
        return list;
    }

    @Override
    public Object getViewById(ModelAndView modelAndView, Long id, List<String> names) throws Exception {
        NoteEntity entity = getById(id);
        if (entity == null) {
            return AbsController.errorView(modelAndView);
        }

        String type = entity.getType();
        if (Type.isText(type)) {
            return getTextViewById(modelAndView, entity);
        }

        if (Type.isDocument(type)) {
            return getDocumentViewById(modelAndView, entity);
        }

        if (Type.ZIP.equals(type)) {
            return getZipViewById(modelAndView, entity, names);
        }

        return AbsController.errorView(modelAndView);
    }

    @Override
    public ResponseEntity<Resource> getStreamById(Long id, List<String> names) throws Exception {
        NoteEntity entity = getById(id);
        if (entity == null) {
            return notFound();
        }

        String type = entity.getType();
        if (Type.isDocument(type)) {
            return getDocumentStreamById(entity, names);
        }

        if (Type.ZIP.equals(type)) {
            return getZipStreamById(entity, names);
        }

        return notFound();
    }

    @Override
    public ResponseEntity<Resource> download(Long id) throws IOException {
        if (id == null) {
            return ResponseEntity.notFound().build();
        }

        // 查库
        NoteEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        if (entity == null || Type.FOLDER.equals(entity.getType())) {
            return ResponseEntity.notFound().build();
        }

        Path path = getPath(id.toString());
        if (!Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] data = Files.readAllBytes(path);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(data);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.add(HttpHeaders.CONTENT_DISPOSITION, String.format("attachment; filename=%s", URLEncoder.encode(entity.getName() + "." + entity.getType(), UTF_8)));
        headers.add(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_OCTET_STREAM_VALUE);
        headers.add(HttpHeaders.CONTENT_LENGTH, String.valueOf(resource.contentLength()));

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    @Override
    public Boolean delById(Long id) {
        NoteEntity entity = verifyId(id);
        if (Type.FOLDER.equals(entity)) {
            NoteEntity child = mapper.selectOne(new LambdaQueryWrapper<NoteEntity>()
                    .select(NoteEntity::getId)
                    .eq(NoteEntity::getPid, id)
                    .last("LIMIT 1"));
            Assert.isNull(child, "无法删除非空文件夹");
        }
        return mapper.deleteById(id) > 0;
    }

    @Override
    public Boolean rename(NoteEntity vo) {
        Long id = vo.getId();
        verifyId(id);

        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        NoteEntity entity = new NoteEntity();
        entity.setId(id);
        entity.setName(name);
        entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.updateById(entity) > 0;
    }

    @Override
    public Boolean paste(NoteEntity vo) {
        Long id = vo.getId();
        verifyId(id);

        Long pid = vo.getPid();
        verifyPid(pid);

        NoteEntity entity = new NoteEntity();
        entity.setId(id);
        entity.setPid(pid);
        entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.updateById(entity) > 0;
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean reUpload(NoteEntity vo) throws IOException {
        NoteEntity entity = verifyId(vo.getId());
        Assert.isTrue(!Type.FOLDER.equals(entity.getType()) && !Type.MD.equals(entity.getType()), "id不能是文件夹类型或者Markdown文件类型");
        return uploadOrReUpload(vo);
    }

    @Transactional(rollbackFor = Exception.class)
    @Override
    public Boolean upload(NoteEntity vo) throws IOException {
        vo.setId(null);
        return uploadOrReUpload(vo);
    }

    @Override
    public Boolean addMd(NoteEntity vo) {
        vo.setType(Type.MD);
        return add(vo);
    }

    @Override
    public Boolean addFolder(NoteEntity vo) {
        vo.setType(Type.FOLDER);
        return add(vo);
    }

    private ResponseEntity<Resource> getDocumentStreamById(NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path path = getPath(id.toString());
        if (!Files.exists(path)) {
            return notFound();
        }

        String type = entity.getType();
        if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX)) {
            Path tmpPath = getTmpPath(id.toString(), true);
            if (!Files.exists(tmpPath)) {
                AsposeUtil.convertDocToPdf(path.toFile(), tmpPath.toFile());
            }
            return ok(tmpPath, MediaType.APPLICATION_PDF);
        }

        if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
            return getXlsDocumentStream(getTmpPath(id.toString()), names);
        }

        if (Type.PDF.equals(type)) {
            return ok(path, MediaType.APPLICATION_PDF);
        }

        return notFound();
    }

    private ResponseEntity<Resource> getXlsDocumentStream(Path path, List<String> names) throws Exception {
        List<Path> paths = getPathsByNames(path, names);
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

    private ResponseEntity<Resource> getZipStreamById(NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path path = getPath(id.toString());
        if (!Files.exists(path)) {
            return ResponseEntity.notFound().build();
        }

        List<Path> tmpPaths = getZipTmpPathsById(entity, names);
        Path tmpPath = tmpPaths.get(tmpPaths.size() - 1);
        if (tmpPath == null || Files.isDirectory(tmpPath)) {
            return ResponseEntity.notFound().build();
        }

        entity = new NoteEntity(tmpPath);
        String type = entity.getType();
        if (type == null) {
            return ResponseEntity.notFound().build();
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
            return ResponseEntity.notFound().build();
        }

        // 读取文件
        byte[] bytes = Files.readAllBytes(tmpPath);

        // 将文件数据转换为资源
        ByteArrayResource resource = new ByteArrayResource(bytes);

        // 响应头
        HttpHeaders headers = new HttpHeaders();
        headers.setContentLength(resource.contentLength());
        headers.setContentType(contentType);

        // 响应
        return new ResponseEntity<>(resource, headers, HttpStatus.OK);
    }

    private ModelAndView getTextViewById(ModelAndView modelAndView, NoteEntity entity) throws IOException {
        String type = entity.getType();
        String view = switch (type) {
            case Type.MD -> "md";
            case Type.HTML -> "html";
            default -> null;
        };

        if (view == null) {
            return AbsController.errorView(modelAndView);
        }

        Long id = entity.getId();
        Path path = getPath(id.toString());
        if (Files.exists(path)) {
            String content = Files.readString(path, UTF_8);
            entity.setContent(content);
        }

        AbsController.setVoAttribute(modelAndView, entity);
        modelAndView.setViewName(String.format("note/%s/view", view));
        return modelAndView;
    }

    private ModelAndView getDocumentViewById(ModelAndView modelAndView, NoteEntity entity) throws Exception {
        String type = entity.getType();
        if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX, Type.PDF)) {
            AbsController.setVoAttribute(modelAndView, entity);
            modelAndView.setViewName("note/pdf/view");
            return modelAndView;
        }

        if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
            return getXlsDocumentViewById(modelAndView, entity);
        }

        return AbsController.errorView(modelAndView);
    }

    private ModelAndView getXlsDocumentViewById(ModelAndView modelAndView, NoteEntity entity) throws Exception {
        Long id = entity.getId();
        Path htmlPath = convertXlsToHtml(id, getPath(id.toString()), getTmpPath(id.toString()));
        if (htmlPath != null) {
            String content = Files.readString(htmlPath, UTF_8);
            entity.setContent(content);
        }
        AbsController.setVoAttribute(modelAndView, entity);
        modelAndView.setViewName("note/html/view");
        return modelAndView;
    }

    /**
     * xls转html
     *
     * @param id          {@link NoteEntity#getId()}
     * @param xlsPath     xls文件
     * @param htmlDirPath html目录
     * @return 返回index.html
     */
    private Path convertXlsToHtml(Long id, Path xlsPath, Path htmlDirPath) throws Exception {
        if (!Files.exists(xlsPath)) {
            return null;
        }

        if (!Files.exists(htmlDirPath)) {
            try {
                Files.createDirectories(htmlDirPath);

                Path indexPath = htmlDirPath.resolve("index");
                AsposeUtil.convertXlsToHtml(xlsPath.toFile(), indexPath.toFile(), htmlDirPath.toFile());

                String filesName = "_files_files";
                Path filesPath = htmlDirPath.resolve(filesName);
                if (Files.exists(filesPath)) {
                    PropertyPlaceholderHelper propertyPlaceholderHelper = new PropertyPlaceholderHelper(String.format("\"%s", filesPath.toAbsolutePath()), "\"");
                    PropertyPlaceholderHelper.PlaceholderResolver placeholderResolver = placeholderName -> {
                        return String.format("\"/note/%s/%s%s/stream\"", id, filesName, placeholderName);
                    };

                    String content = Files.readString(indexPath, UTF_8);
                    content = propertyPlaceholderHelper.replacePlaceholders(content, placeholderResolver);
                    Files.write(indexPath, content.getBytes(UTF_8), StandardOpenOption.TRUNCATE_EXISTING);

                    // 获取files目录下所有html文件
                    List<Path> htmlPaths = Files.list(filesPath)
                            .filter(path1 -> Type.HTML.equals(Type.pathOf(path1)))
                            .collect(Collectors.toList());
                    if (CollectionUtils.isNotEmpty(htmlPaths)) {
                        List<String> placeholderPrefixes = List.of("href", "src");
                        for (String placeholderPrefix : placeholderPrefixes) {
                            String key = UUID.randomUUID().toString().replace("-", ".");
                            propertyPlaceholderHelper = new PropertyPlaceholderHelper(placeholderPrefix + "=\"", "\"");
                            placeholderResolver = placeholderName -> {
                                return String.format("%s=\"/note/%s/%s/%s/stream\"", key, id, filesName, placeholderName);
                            };

                            for (Path htmlPath : htmlPaths) {
                                content = Files.readString(htmlPath, UTF_8);
                                content = propertyPlaceholderHelper.replacePlaceholders(content, placeholderResolver);
                                content = content.replace(key, placeholderPrefix);
                                Files.write(htmlPath, content.getBytes(UTF_8), StandardOpenOption.TRUNCATE_EXISTING);
                            }
                        }
                    }
                }
            } catch (Exception e) {
                PathUtils.deleteFile(htmlDirPath);
                throw e;
            }
        }
        return htmlDirPath.resolve("index");
    }

    private ModelAndView getZipViewById(ModelAndView modelAndView, NoteEntity entity, List<String> names) throws Exception {
        Long id = entity.getId();
        Path path = getPath(id.toString());
        if (!Files.exists(path)) {
            return zipView(modelAndView, entity, null);
        }

        Path zipPath = getPath(id.toString());
        Path dirPath = getTmpPath(id.toString());
        unzipIfNotDecompressed(zipPath, dirPath);
        List<Path> paths = getPathsByNames(dirPath, names);
        path = paths.get(paths.size() - 1);

        int size = paths.size();
        if (size > 1) {
            List<NoteEntity> ps = new ArrayList<>(size);
            int index = 1;
            while (index < size) {
                Path path1 = paths.get(index++);
                if (path1 == null) {
                    break;
                }
                ps.add(new NoteEntity(path1));
            }
            entity.setPs(ps);
        }

        if (path == null) {
            return zipView(modelAndView, entity, null);
        }

        if (Files.isDirectory(path)) {
            List<NoteEntity> childList = new ArrayList<>(16);

            // 获取目录下的子文件夹
            Iterator<Path> iterator = Files.list(path).iterator();
            while (iterator.hasNext()) {
                path = iterator.next();
                childList.add(new NoteEntity(path));
            }

            Collections.sort(childList);

            entity.setChildList(new org.xiangqian.note.util.List<>(childList));

            return zipView(modelAndView, entity, null);
        }

        entity = new NoteEntity(path);
        entity.setId(id);
        String type = entity.getType();
        if (Type.isText(type)) {
            // 读取文件内容
            String content = Files.readString(path, UTF_8);
            entity.setContent(content);
            return zipView(modelAndView, entity, "text");
        }

        if (Type.isDocument(type)) {
            if (StringUtils.equalsAny(type, Type.DOC, Type.DOCX, Type.PDF)) {
                return zipView(modelAndView, entity, "pdf");
            }

            if (StringUtils.equalsAny(type, Type.XLS, Type.XLSX)) {
                Path newTmpPath = getTmpPath(id + "_" + Md5Util.encryptHex(paths.stream().map(Path::getFileName).map(Path::toString).collect(Collectors.joining("/"))));
                Path htmlPath = convertXlsToHtml(id, path, newTmpPath);
                if (htmlPath != null) {
                    // 读取文件内容
                    String content = Files.readString(htmlPath, UTF_8);
                    entity.setContent(content);
                }
                return zipView(modelAndView, entity, "html");
            }

            return zipView(modelAndView, entity, "unsupported");
        }

        if (Type.isImage(type)) {
            return zipView(modelAndView, entity, "image");
        }

        return zipView(modelAndView, entity, "unsupported");
    }

    private ModelAndView zipView(ModelAndView modelAndView, NoteEntity entity, String subview) {
        AbsController.setVoAttribute(modelAndView, entity);

        String view = null;
        if (subview != null) {
            view = String.format("note/zip/%s/view", subview);
        } else {
            view = "note/zip/view";
        }
        modelAndView.setViewName(view);

        return modelAndView;
    }

    private List<Path> getZipTmpPathsById(NoteEntity entity, List<String> names) throws IOException {
        Long id = entity.getId();
        Path zipPath = getPath(id.toString());
        Path dirPath = getTmpPath(id.toString());
        unzipIfNotDecompressed(zipPath, dirPath);
        return getPathsByNames(dirPath, names);
    }

    private List<Path> getPathsByNames(Path path, List<String> names) throws IOException {
        if (CollectionUtils.isEmpty(names)) {
            return List.of(path);
        }

        List<Path> paths = new ArrayList<>(names.size() + 1);
        paths.add(path);
        for (String name : names) {
            if (path == null || !Files.isDirectory(path)) {
                break;
            }

            // 获取目录下的子文件夹
            Iterator<Path> iterator = Files.list(path).iterator();
            while (iterator.hasNext()) {
                path = iterator.next();
                if (path.getFileName().toString().equals(name)) {
                    break;
                }
                path = null;
            }
            paths.add(path);
        }
        return paths;
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
        if (!Files.exists(zipPath) || Files.exists(dirPath)) {
            return;
        }

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

    private Boolean uploadOrReUpload(NoteEntity vo) throws IOException {
        // 文件
        MultipartFile file = vo.getFile();
        Assert.isTrue(!file.isEmpty(), "上传文件不能为空");

        // 文件名称
        String name = StringUtils.trim(file.getOriginalFilename());
        Assert.isTrue(StringUtils.isNotEmpty(name), "上传文件名称不能为空");
        // 文件后缀名
        String suffix = null;
        // 文件类型
        String type = null;
        int index = name.lastIndexOf(".");
        if (index >= 0 && (index + 1) < name.length()) {
            suffix = StringUtils.trim(name.substring(index + 1).toLowerCase());
            type = Type.suffixOf(suffix);
            name = name.substring(0, index);
        }
        Assert.isTrue(StringUtils.equalsAny(type,
                        Type.DOC, Type.DOCX,
                        Type.XLS, Type.XLSX,
                        Type.PDF,
                        Type.HTML,
                        Type.ZIP),
                String.format("不支持上传 %s 文件类型，请选择 doc、docx、pdf、html、zip 文件类型上传", suffix));

        Long pid = vo.getPid();
        verifyPid(pid);

        byte[] bytes = file.getBytes();

        NoteEntity entity = new NoteEntity();
        entity.setPid(pid);
        entity.setName(name);
        entity.setType(type);
        entity.setSize(bytes.length + 0L);
        if (vo.getId() == null) {
            entity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
            mapper.insert(entity);
        } else {
            entity.setId(vo.getId());
            entity.setUpdTime(DateUtil.toSecond(LocalDateTime.now()));
            mapper.updateById(entity);
        }

        Path path = getPath(entity.getId().toString(), true);
        Files.write(path, bytes);
        return true;
    }

    private Boolean add(NoteEntity vo) {
        String name = StringUtils.trim(vo.getName());
        Assert.isTrue(StringUtils.isNotEmpty(name), "名称不能为空");

        Long pid = vo.getPid();
        verifyPid(pid);

        NoteEntity addEntity = new NoteEntity();
        addEntity.setPid(pid);
        addEntity.setName(name);
        addEntity.setType(vo.getType());
        addEntity.setAddTime(DateUtil.toSecond(LocalDateTime.now()));
        return mapper.insert(addEntity) > 0;
    }

    private NoteEntity verifyId(Long id) {
        Assert.notNull(id, "id不能为空");

        NoteEntity entity = null;
        if (id.longValue() > 0) {
            entity = mapper.selectById(id);
        }
        Assert.notNull(entity, "id不存在");

        return entity;
    }

    private void verifyPid(Long pid) {
        Assert.notNull(pid, "pid不能为空");

        // 根节点
        if (pid.longValue() == 0) {
            return;
        }

        NoteEntity entity = null;
        if (pid.longValue() > 0) {
            entity = mapper.selectById(pid);
        }
        Assert.isTrue(entity != null && Type.FOLDER.equals(entity.getType()), "pid不存在");
    }

    private NoteEntity getById(Long id) {
        if (id == null || id.longValue() <= 0) {
            return null;
        }
        return mapper.getById(id);
    }

    public NoteEntity getById1(Long id) {
        Assert.notNull(id, "id不能为空");
        Assert.isTrue(id.longValue() >= 0, "id不能小于0");
        if (id.longValue() == 0) {
            NoteEntity entity = new NoteEntity();
            entity.setId(0L);
            entity.setType(Type.FOLDER);
            return entity;
        }

        NoteEntity entity = mapper.getById(id);
        if (entity != null) {
            List<NoteEntity> ps = entity.getPs();
            if (CollectionUtils.isNotEmpty(ps)) {
                Set<Long> pids = ps.stream().map(NoteEntity::getId).collect(Collectors.toSet());
                // <pid, name>
                Map<Long, String> pidMap = mapper.selectList(new LambdaQueryWrapper<NoteEntity>()
                                .select(NoteEntity::getId, NoteEntity::getName)
                                .in(NoteEntity::getId, pids))
                        .stream().collect(Collectors.toMap(NoteEntity::getId, NoteEntity::getName));
                for (NoteEntity p : ps) {
                    p.setName(pidMap.get(p.getId()));
                }
            }
        }
        return entity;
    }

}
