//package org.xiangqian.note.entity;
//
//import com.baomidou.mybatisplus.annotation.*;
//import lombok.Data;
//
///**
// * image、audio、video信息
// *
// * @author xiangqian
// * @date 23:29 2024/03/04
// */
//@Data
//@TableName("iav")
//public class IavEntity {
//
//    // id
//    @TableId(type = IdType.AUTO)
//    private Long id;
//
//    // 名称
//    @TableField("`name`")
//    private String name;
//
//    // 类型
//    private String type;
//
//    // 大小，单位：byte
//    private Long size;
//
//    // 删除标识，0-正常，1-删除
//    @TableLogic
//    private Integer del;
//
//    // 创建时间（时间戳，单位s）
//    private Long addTime;
//
//    // 修改时间（时间戳，单位s）
//    private Long updTime;
//
//}
