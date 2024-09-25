-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `name`               TEXT UNIQUE NOT NULL, -- 用户名
    `passwd`             TEXT        NOT NULL, -- 密码
    `deny`               INTEGER DEFAULT 0,    -- 连续错误登陆次数，超过3次则锁定用户
    `last_login_host`    TEXT    DEFAULT '',   -- 上一次登录主机
    `last_login_time`    INTEGER DEFAULT 0,    -- 上一次登录时间戳（单位s）
    `current_login_host` TEXT    DEFAULT '',   -- 当前登录主机
    `current_login_time` INTEGER DEFAULT 0,    -- 当前登录时间戳（单位s）
    `add_time`           INTEGER DEFAULT 0,    -- 创建时间戳（单位s）
    `upd_time`           INTEGER DEFAULT 0     -- 修改时间戳（单位s）
);

INSERT INTO `user`(`name`, `passwd`, `add_time`)
VALUES ('admin', '$2a$10$ZsS2bA7B7AQtIBBpW7xz3OIw3FWU0CnXX7HZMi6vBNt9ZNcA2RNGG', 1675393080);


-- ------------------------
-- Table structure for note
-- ------------------------
DROP TABLE IF EXISTS `note`;
CREATE TABLE `note` -- 笔记信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键
    `pid`      INTEGER DEFAULT 0,                 -- 父id
    `name`     TEXT    DEFAULT '',                -- 名称
    `type`     TEXT    DEFAULT '',                -- 类型（folder、md、doc、docx、xls、xlsx、pdf、html、zip）
    `size`     INTEGER DEFAULT 0,                 -- 大小（单位byte）
    `del`      INTEGER DEFAULT 0,                 -- 是否已删除（0-否，1-是）
    `add_time` INTEGER DEFAULT 0,                 -- 创建时间戳（单位s）
    `upd_time` INTEGER DEFAULT 0                  -- 修改时间戳（单位s）
);


-- -------------------------
-- Table structure for image
-- -------------------------
DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` -- 图片信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键
    `name`     TEXT    DEFAULT '',                -- 名称
    `type`     TEXT    DEFAULT '',                -- 类型（png、jpg、gif、webp、ico）
    `size`     INTEGER DEFAULT 0,                 -- 大小（单位byte）
    `del`      INTEGER DEFAULT 0,                 -- 是否已删除（0-否，1-是）
    `add_time` INTEGER DEFAULT 0,                 -- 创建时间戳（单位s）
    `upd_time` INTEGER DEFAULT 0                  -- 修改时间戳（单位s）
);
