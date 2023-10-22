-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户id
    `name`     VARCHAR(64)  NOT NULL,             -- 用户名
    `nickname` VARCHAR(64)  DEFAULT '',           -- 昵称
    `passwd`   VARCHAR(128) NOT NULL,             -- 密码
    `rem`      VARCHAR(256) DEFAULT '',           -- 备注
    `try`      TINYINT      DEFAULT 0,            -- 尝试输入密码次数，超过3次账号将会被锁定
    `del`      TINYINT      DEFAULT 0,            -- 删除标识，0-正常，1-删除
    `add_time` INTEGER      DEFAULT 0,            -- 创建时间（时间戳，s）
    `upd_time` INTEGER      DEFAULT 0             -- 修改时间（时间戳，s），也可作锁定时间
);
