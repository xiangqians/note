-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`        INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`      VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname`  VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`    VARCHAR(128) NOT NULL,             -- 密码
    `rem`       VARCHAR(256) DEFAULT '',           -- 备注
    `try_count` TINYINT      DEFAULT 0,            -- 尝试输入次数
    `lock`      TINYINT      DEFAULT 0,            -- 锁定标识，0-正常，1-锁定
    `lock_time` INTEGER      DEFAULT 0,            -- 锁定时间（时间戳，s）
    `del`       TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time`  INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time`  INTEGER      DEFAULT 0             -- 修改时间（时间戳，s）
);