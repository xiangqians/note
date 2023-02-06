-- ------------------------
-- Table structure for file
-- ------------------------
DROP TABLE IF EXISTS `file`;
CREATE TABLE `file` -- 文件信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`     VARCHAR(128) DEFAULT '',           -- 文件名称
    `type`     VARCHAR(64)  DEFAULT '',           -- 文件类型
    `size`     INTEGER      DEFAULT 0,            -- 文件大小，单位：byte
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `del`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- -----------------------
-- Table structure for dir
-- -----------------------
DROP TABLE IF EXISTS `dir`;
CREATE TABLE `dir` -- 目录信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `pid`      INTEGER      DEFAULT 0,            -- 父目录id
    `name`     VARCHAR(128) DEFAULT '',           -- 目录名称
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `del`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);


-- ----------------------------
-- Table structure for dir_file
-- ----------------------------
DROP TABLE IF EXISTS `dir_file`;
CREATE TABLE `dir_file` -- 目录文件信息表
(
    `dir_id`   INTEGER NOT NULL,  -- 目录id
    `file_id`  INTEGER NOT NULL,  -- 文件id
    `add_time` INTEGER DEFAULT 0, -- 创建时间（时间戳，s）
    PRIMARY KEY (`dir_id`, `file_id`)
);