-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`                 INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`               VARCHAR(64)  DEFAULT '',           -- 用户名
    `nickname`           VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`             VARCHAR(128) DEFAULT '',           -- 密码
    `locked`             TINYINT      DEFAULT 0,            -- 是否已锁定，0-否，1-是
    `deny`               TINYINT      DEFAULT 0,            -- 用户连续错误登陆次数，超过3次则锁定用户
    `last_login_ip`      VARCHAR(64)  DEFAULT '',           -- 上一次登录ip
    `last_login_time`    INTEGER      DEFAULT 0,            -- 上一次登录时间（时间戳，单位s）
    `current_login_ip`   VARCHAR(64)  DEFAULT '',           -- 当前登录ip
    `current_login_time` INTEGER      DEFAULT 0,            -- 当前登录时间（时间戳，单位s）
    `add_time`           INTEGER      DEFAULT 0,            -- 创建时间（时间戳，单位s）
    `upd_time`           INTEGER      DEFAULT 0             -- 修改时间（时间戳，单位s）
);

INSERT INTO `user`(`name`, `nickname`, `passwd`, `add_time`)
VALUES ('admin', '管理员', '$2a$10$ZsS2bA7B7AQtIBBpW7xz3OIw3FWU0CnXX7HZMi6vBNt9ZNcA2RNGG', STRFTIME('%s', 'now'));


-- -------------------
-- Table model for iav
-- -------------------
DROP TABLE IF EXISTS `iav`;
CREATE TABLE `iav` -- image、audio、video信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `name`     VARCHAR(64) DEFAULT '',            -- 名称
    `type`     VARCHAR(32) DEFAULT '',            -- 类型
    `size`     INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `del`      TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time` INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time` INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);


-- --------------------
-- Table model for note
-- --------------------
DROP TABLE IF EXISTS `note`;
CREATE TABLE `note` -- 笔记信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- id
    `pid`      INTEGER     DEFAULT 0,             -- 父id
    `name`     VARCHAR(64) DEFAULT '',            -- 名称
    `type`     VARCHAR(8)  DEFAULT '',            -- 类型（folder、md、doc、pdf、html、zip等）
    `size`     INTEGER     DEFAULT 0,             -- 大小，单位：byte
    `del`      TINYINT     DEFAULT 0,             -- 删除标识，0-正常，1-删除
    `add_time` INTEGER     DEFAULT 0,             -- 创建时间（时间戳，单位s）
    `upd_time` INTEGER     DEFAULT 0              -- 修改时间（时间戳，单位s）
);
