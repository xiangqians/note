-- -------------------------
-- Table structure for image
-- -------------------------
DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` -- 图片信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`      VARCHAR(64) DEFAULT '',            -- 名称
    `type`      VARCHAR(8)  DEFAULT '',            -- 类型
    `size`      INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `hist`      TEXT        DEFAULT '',            -- history（历史记录）
    `hist_size` INTEGER     DEFAULT 0,             -- history（历史记录）大小，单位：byte
    `del`       TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER     DEFAULT 0,             -- 创建时间（时间戳，s）
    `upd_time`  INTEGER     DEFAULT 0              -- 修改时间（时间戳，s）
);


-- -------------------------
-- Table structure for audio
-- -------------------------
DROP TABLE IF EXISTS `audio`;
CREATE TABLE `audio` -- 音频信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`      VARCHAR(64) DEFAULT '',            -- 名称
    `type`      VARCHAR(8)  DEFAULT '',            -- 类型
    `size`      INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `hist`      TEXT        DEFAULT '',            -- history（历史记录）
    `hist_size` INTEGER     DEFAULT 0,             -- history（历史记录）大小，单位：byte
    `del`       TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER     DEFAULT 0,             -- 创建时间（时间戳，s）
    `upd_time`  INTEGER     DEFAULT 0              -- 修改时间（时间戳，s）
);


-- -------------------------
-- Table structure for video
-- -------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video` -- 视频信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`      VARCHAR(64) DEFAULT '',            -- 名称
    `type`      VARCHAR(8)  DEFAULT '',            -- 类型
    `size`      INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `hist`      TEXT        DEFAULT '',            -- history（历史记录）
    `hist_size` INTEGER     DEFAULT 0,             -- history（历史记录）大小，单位：byte
    `del`       TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER     DEFAULT 0,             -- 创建时间（时间戳，s）
    `upd_time`  INTEGER     DEFAULT 0              -- 修改时间（时间戳，s）
);


-- ------------------------
-- Table structure for note
-- ------------------------
DROP TABLE IF EXISTS `note`;
CREATE TABLE `note` -- 笔记信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `pid`       INTEGER     DEFAULT 0,             -- 父id
    `name`      VARCHAR(64) DEFAULT '',            -- 笔记名称
    `type`      VARCHAR(8)  DEFAULT '',            -- 笔记类型
    `size`      INTEGER     DEFAULT 0,             -- 笔记大小，单位：byte
    `hist`      TEXT        DEFAULT '',            -- history（历史记录）
    `hist_size` INTEGER     DEFAULT 0,             -- history（历史记录）笔记大小，单位：byte
    `del`       TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER     DEFAULT 0,             -- 创建时间（时间戳，s）
    `upd_time`  INTEGER     DEFAULT 0              -- 修改时间（时间戳，s）
);
-- 存储dir、md、pdf、html、zip等

