-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`                   INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`                 VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname`             VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`               VARCHAR(128) NOT NULL,             -- 密码
    `rem`                  VARCHAR(256) DEFAULT '',           -- 备注
    `try`                  TINYINT      DEFAULT 0,            -- 尝试输入密码次数，超过3次账号将会被锁定
    `last_sign_in_ip`      VARCHAR(64)  DEFAULT '',           -- 上一次登录IP
    `last_sign_in_time`    INTEGER      DEFAULT 0,            -- 上一次登录时间（时间戳，单位s）
    `current_sign_in_ip`   VARCHAR(64)  DEFAULT '',           -- 当前次登录IP
    `current_sign_in_time` INTEGER      DEFAULT 0,            -- 当前登录时间（时间戳，单位s）
    `del`                  TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time`             INTEGER      DEFAULT 0,            -- 创建时间（时间戳，单位s）
    `upd_time`             INTEGER      DEFAULT 0             -- 修改时间（时间戳，单位s），也可作锁定时间
);


-- -------------------------
-- Table structure for image
-- -------------------------
DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` -- 图片信息表
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`         VARCHAR(64) DEFAULT '',            -- 名称
    `type`         VARCHAR(8)  DEFAULT '',            -- 类型
    `size`         INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `history`      TEXT        DEFAULT '',            -- 历史
    `history_size` INTEGER     DEFAULT 0,             -- 历史大小，单位：byte
    `del`          TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除，2-永久删除
    `add_time`     INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time`     INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);


-- -------------------------
-- Table structure for audio
-- -------------------------
DROP TABLE IF EXISTS `audio`;
CREATE TABLE `audio` -- 音频信息表
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`         VARCHAR(64) DEFAULT '',            -- 名称
    `type`         VARCHAR(8)  DEFAULT '',            -- 类型
    `size`         INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `history`      TEXT        DEFAULT '',            -- 历史
    `history_size` INTEGER     DEFAULT 0,             -- 历史大小，单位：byte
    `del`          TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除，2-永久删除
    `add_time`     INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time`     INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);


-- -------------------------
-- Table structure for video
-- -------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video` -- 视频信息表
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`         VARCHAR(64) DEFAULT '',            -- 名称
    `type`         VARCHAR(8)  DEFAULT '',            -- 类型
    `size`         INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `history`      TEXT        DEFAULT '',            -- 历史
    `history_size` INTEGER     DEFAULT 0,             -- 历史大小，单位：byte
    `del`          TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除，2-永久删除
    `add_time`     INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time`     INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);


-- ------------------------
-- Table structure for note
-- ------------------------
DROP TABLE IF EXISTS `note`;
CREATE TABLE `note` -- 笔记信息表
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `pid`          INTEGER     DEFAULT 0,             -- 父id
    `name`         VARCHAR(64) DEFAULT '',            -- 名称
    `type`         VARCHAR(8)  DEFAULT '',            -- 类型
    `size`         INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `history`      TEXT        DEFAULT '',            -- 历史
    `history_size` INTEGER     DEFAULT 0,             -- 历史大小，单位：byte
    `del`          TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除，2-永久删除
    `add_time`     INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time`     INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);
/*
存储folder、md、doc、pdf、html、zip等
*/

