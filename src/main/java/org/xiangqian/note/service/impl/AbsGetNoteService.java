//package org.xiangqian.note.service.impl;
//
//import org.apache.commons.collections4.CollectionUtils;
//import org.apache.commons.io.file.PathUtils;
//import org.springframework.core.env.Environment;
//import org.springframework.util.PropertyPlaceholderHelper;
//import org.xiangqian.note.service.GetNoteService;
//import org.xiangqian.note.util.AsposeUtil;
//import org.xiangqian.note.model.Type;
//
//import java.io.IOException;
//import java.nio.file.Files;
//import java.nio.file.Path;
//import java.util.ArrayList;
//import java.util.Iterator;
//import java.util.List;
//import java.util.UUID;
//import java.util.stream.Collectors;
//
///**
// * @author xiangqian
// * @date 16:47 2024/05/03
// */
//public abstract class AbsGetNoteService extends AbsService implements GetNoteService {
//
//    protected AbsGetNoteService(Environment environment) throws IOException {
//        super(environment);
//    }
//
//    protected List<Path> getTmpPaths(Path tmpPath, List<String> names) throws IOException {
//        if (!Files.exists(tmpPath)) {
//            return null;
//        }
//
//        if (CollectionUtils.isEmpty(names)) {
//            return List.of(tmpPath);
//        }
//
//        List<Path> tmpPaths = new ArrayList<>(names.size() + 1);
//        tmpPaths.add(tmpPath);
//        for (String name : names) {
//            if (tmpPath == null || !Files.isDirectory(tmpPath)) {
//                tmpPaths.add(null);
//                break;
//            }
//
//            // 获取目录下的子文件夹
//            Iterator<Path> iterator = Files.list(tmpPath).iterator();
//            while (iterator.hasNext()) {
//                tmpPath = iterator.next();
//                if (tmpPath.getFileName().toString().equals(name)) {
//                    break;
//                }
//                tmpPath = null;
//            }
//            tmpPaths.add(tmpPath);
//        }
//        return tmpPaths;
//    }
//
//    /**
//     * xls转html
//     *
//     * @param baseUrl  基础url
//     * @param xlsPath  xls文件
//     * @param htmlPath html目录
//     * @return 返回index.html
//     */
//    protected Path convertXlsToHtml(String baseUrl, Path xlsPath, Path htmlPath) throws Exception {
//        if (!Files.exists(xlsPath)) {
//            return null;
//        }
//
//        if (!Files.exists(htmlPath)) {
//            try {
//                Files.createDirectories(htmlPath);
//
//                Path indexPath = htmlPath.resolve("index");
//                AsposeUtil.convertXlsToHtml(xlsPath.toFile(), indexPath.toFile(), htmlPath.toFile());
//
//                String filesName = "_files_files";
//                Path filesPath = htmlPath.resolve(filesName);
//                if (Files.exists(filesPath)) {
//                    PropertyPlaceholderHelper propertyPlaceholderHelper = new PropertyPlaceholderHelper(String.format("\"%s", filesPath.toAbsolutePath()), "\"");
//                    PropertyPlaceholderHelper.PlaceholderResolver placeholderResolver = placeholderName -> {
//                        return String.format("\"%s/%s%s/stream\"", baseUrl, filesName, placeholderName);
//                    };
//
//                    String content = Files.readString(indexPath, UTF_8);
//                    content = propertyPlaceholderHelper.replacePlaceholders(content, placeholderResolver);
//                    // 将内容写入文件（覆盖），如果文件不存在则创建
//                    Files.write(indexPath, content.getBytes(UTF_8));
//
//                    // 获取files目录下所有html文件
//                    List<Path> htmlPaths = Files.list(filesPath)
//                            .filter(path1 -> Type.HTML.equals(Type.pathOf(path1)))
//                            .collect(Collectors.toList());
//                    if (CollectionUtils.isNotEmpty(htmlPaths)) {
//                        List<String> placeholderPrefixes = List.of("href", "src");
//                        for (String placeholderPrefix : placeholderPrefixes) {
//                            String key = UUID.randomUUID().toString().replace("-", ".");
//                            propertyPlaceholderHelper = new PropertyPlaceholderHelper(placeholderPrefix + "=\"", "\"");
//                            placeholderResolver = placeholderName -> {
//                                return String.format("%s=\"%s/%s/%s/stream\"", key, baseUrl, filesName, placeholderName);
//                            };
//
//                            for (Path htmlPath1 : htmlPaths) {
//                                content = Files.readString(htmlPath1, UTF_8);
//                                content = propertyPlaceholderHelper.replacePlaceholders(content, placeholderResolver);
//                                content = content.replace(key, placeholderPrefix);
//                                // 将内容写入文件（覆盖），如果文件不存在则创建
//                                Files.write(htmlPath1, content.getBytes(UTF_8));
//                            }
//                        }
//                    }
//                }
//            } catch (Exception e) {
//                PathUtils.deleteFile(htmlPath);
//                throw e;
//            }
//        }
//        return htmlPath.resolve("index");
//    }
//
//}
