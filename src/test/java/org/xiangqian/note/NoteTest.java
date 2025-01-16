package org.xiangqian.note;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.collections4.CollectionUtils;
import org.apache.commons.io.FileUtils;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.xiangqian.note.entity.NoteEntity;
import org.xiangqian.note.mapper.NoteMapper;
import org.xiangqian.note.util.Type;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.function.BiConsumer;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * @author xiangqian
 * @date 11:35 2024/11/24
 */
@SpringBootTest
public class NoteTest {

    private static final ObjectMapper objectMapper = new ObjectMapper();

    private static final String json1 = "";
    private static final String json2 = "";

    @Value("${spring.datasource.url}")
    private String url;

    @Autowired
    private NoteMapper noteMapper;

    // 步骤 1
    @Test
    public void testGet() throws Exception {
        System.out.println(url);
        get();
        System.out.println("/");
    }

    // 步骤 2
    @Test
    public void testAdd() throws Exception {
        System.out.println(url);
        add();
        System.out.println("/");
    }

    // 步骤 3
    public static void main(String[] args) throws Exception {
//        script();
//        test1();
//        test2();
        test3();
    }

    private static void test3() throws Exception {
        List<String[]> list = Files.lines(Path.of("target/mv1.sh"))
                .filter(line -> line.startsWith("mv"))
                .map(line -> {
                    String[] array = line.split(" ");
                    int length = array.length;
                    return new String[]{array[length - 2].replace("\"", ""), array[length - 1].replace("\"", "")};
                })
                .collect(Collectors.toList());
        for (String[] array : list) {
            System.out.println(Arrays.toString(array));
            // 指定源文件路径和目标文件路径
            Path sourcePath = Paths.get("F:\\note\\my\\data\\note", array[0]);
            if (Files.exists(sourcePath)) {
                Path targetPath = Paths.get("F:\\note\\my\\data\\note", array[1]);
                Files.move(sourcePath, targetPath);
            }
        }
    }

    private static void test2() throws Exception {
        List<NoteEntity> list = objectMapper.readValue(json1, new TypeReference<List<NoteEntity>>() {
        });
        for (NoteEntity entity : list) {
            System.out.println(entity);
        }
    }

    private static void test1() throws Exception {
//        String str = Files.lines(Path.of("target/mv1.sh"))
//                .filter(line->line.startsWith("mv"))
//                .map(line -> line.split(" ")[2])
//                .collect(Collectors.joining(" "));
//        System.out.println(str);

        String str = Files.lines(Path.of("target/mv1.sh"))
                .filter(line -> line.startsWith("mv"))
                .map(line -> line.split(" ")[2].replace("\"", ""))
                .collect(Collectors.joining(", "));
        System.out.println(str);

//        String str = Files.lines(Path.of("target/mv1.sh"))
//                .map(line -> {
//                    if (line.startsWith("mv")) {
//                        String[] array = line.split(" ");
//                        return array[0] + " " + array[2] + " " + array[1];
//                    }
//                    return line;
//                })
//                .collect(Collectors.joining("\n"));
//        System.out.println(str);

    }

    private static void script() throws Exception {
        List<NoteEntity> list = objectMapper.readValue(json2, new TypeReference<List<NoteEntity>>() {
        });

        StringBuilder mv1ScriptBuilder = new StringBuilder();
        mv1ScriptBuilder.append("#!/bin/bash");

        StringBuilder mv2ScriptBuilder = new StringBuilder();
        mv2ScriptBuilder.append("#!/bin/bash");
        mv2ScriptBuilder.append("\n").append("mv -n ");

        for (NoteEntity entity : list) {
            script(mv1ScriptBuilder, mv2ScriptBuilder, entity);
        }

        FileUtils.write(new File("target/mv1.sh"), mv1ScriptBuilder.toString(), StandardCharsets.UTF_8);

        mv2ScriptBuilder.append(" /data/data/com.termux/files/home/note/data/");
        FileUtils.write(new File("target/mv2.sh"), mv2ScriptBuilder.toString(), StandardCharsets.UTF_8);
    }

    private static void script(StringBuilder mv1ScriptBuilder, StringBuilder mv2ScriptBuilder, NoteEntity entity) {
        System.out.format("### %s\n", entity);
        if (!Type.FOLDER.equals(entity.getType())) {
            mv1ScriptBuilder.append("\n").append("mv -n")
                    .append(" \"").append(entity.getId()).append("\"")
                    .append(" \"").append(entity.getNewId()).append("\"");

            mv2ScriptBuilder.append(" \"").append(entity.getNewId()).append("\"");
        }

        List<NoteEntity> children = entity.getChildren();
        if (CollectionUtils.isNotEmpty(children)) {
            for (NoteEntity child : children) {
                script(mv1ScriptBuilder, mv2ScriptBuilder, child);
            }
        }
    }

    private void add() throws Exception {
        List<NoteEntity> list = objectMapper.readValue(json1, new TypeReference<List<NoteEntity>>() {
        });

        for (NoteEntity entity : list) {
            add(entity, 0L);
        }

        String json = objectMapper.writeValueAsString(list);
        System.out.println(json);
    }

    private int index = 0;

    private void add(NoteEntity entity, Long pid) {
        System.out.format("### [%s] %s | %s\n", ++index, pid, entity);

        NoteEntity createEntity = new NoteEntity();
        createEntity.setPid(pid);
        createEntity.setName(entity.getName());
        createEntity.setType(entity.getType());
        createEntity.setSize(entity.getSize());
        createEntity.setCreateTime(entity.getCreateTime());
        createEntity.setUpdateTime(entity.getUpdateTime());
        noteMapper.create(createEntity);

        Long id = createEntity.getId();
        entity.setNewId(id);
        entity.setNewPid(pid);

        List<NoteEntity> children = entity.getChildren();
        if (CollectionUtils.isNotEmpty(children)) {
            for (NoteEntity child : children) {
                add(child, id);
            }
        }
    }

    private void get() throws Exception {
        List<NoteEntity> list = noteMapper.getList();
        System.out.println(list);

        List<NoteEntity> result = associateNodes(list,
                NoteEntity::getId,
                NoteEntity::getPid,
                (entity, noteEntities) -> entity.setChildren(noteEntities),
                100);
        System.out.println(result);

        String json = objectMapper.writeValueAsString(result);
        System.out.println(json);
    }

    private <T> List<T> associateNodes(List<T> candidateNodes,
                                       Function<T, Object> getIdFunction,
                                       Function<T, Object> getParentIdFunction,
                                       BiConsumer<T, List<T>> setChildNodesConsumer,
                                       int depth) {
        if (CollectionUtils.isEmpty(candidateNodes)) {
            return Collections.emptyList();
        }

        List<T> candidateRootNodes = candidateNodes.stream().collect(Collectors.toList());
        for (T candidateRootNode : candidateRootNodes) {
            associateNodes(candidateRootNode, candidateNodes, getIdFunction, getParentIdFunction, setChildNodesConsumer, depth);
        }
        return candidateNodes;
    }

    /**
     * 递归关联节点关系
     *
     * @param node                  节点
     * @param candidateNodes        候选节点集
     * @param getIdFunction         获取节点唯一标识id
     * @param getParentIdFunction   获取节点父节点id
     * @param setChildNodesConsumer 设置节点的子节点集
     * @param depth                 递归深度
     * @param <T>
     */
    private <T> void associateNodes(T node,
                                    List<T> candidateNodes,
                                    Function<T, Object> getIdFunction,
                                    Function<T, Object> getParentIdFunction,
                                    BiConsumer<T, List<T>> setChildNodesConsumer,
                                    int depth) {
        if (depth <= 0 || CollectionUtils.isEmpty(candidateNodes)) {
            return;
        }

        // 获取子节点集
        List<T> childNodes = null;
        Object id = getIdFunction.apply(node);
        Iterator<T> iterator = candidateNodes.iterator();
        while (iterator.hasNext()) {
            T candidateNode = iterator.next();
            Object parentId = getParentIdFunction.apply(candidateNode);
            if (Objects.equals(id, parentId)) {
                if (childNodes == null) {
                    childNodes = new ArrayList<>(8);
                }
                childNodes.add(candidateNode);
                iterator.remove();
            }
        }
        if (childNodes == null) {
            return;
        }
        setChildNodesConsumer.accept(node, childNodes);

        // 关联子节点
        depth--;
        for (T childNode : childNodes) {
            associateNodes(childNode, candidateNodes, getIdFunction, getParentIdFunction, setChildNodesConsumer, depth);
        }
    }

}
