-- ------------------------
-- Table structure for file
-- ------------------------
DROP TABLE IF EXISTS `file`;
CREATE TABLE `file` -- 文件信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `pid`      INTEGER      DEFAULT 0,            -- 父id
    `name`     VARCHAR(128) DEFAULT '',           -- 文件名称
    `type`     VARCHAR(64)  DEFAULT '',           -- 文件类型，d表示目录，其他则是文件
    `size`     INTEGER      DEFAULT 0,            -- 文件大小，单位：byte
    `del`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- -----------------------
-- Table structure for img
-- -----------------------
DROP TABLE IF EXISTS `img`;
CREATE TABLE `img` -- 图片信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`     VARCHAR(128) DEFAULT '',           -- 图片名称
    `type`     VARCHAR(64)  DEFAULT '',           -- 图片类型
    `size`     INTEGER      DEFAULT 0,            -- 图片大小，单位：byte
    `del`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);