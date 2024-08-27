-- ------------------------
-- Table structure for user
-- ------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` -- 用户信息表
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键
    `name`     TEXT    DEFAULT '',                -- 用户名
    `nickname` TEXT    DEFAULT '',                -- 昵称
    `passwd`   TEXT    DEFAULT '',                -- 密码
    `locked`   INTEGER DEFAULT 0,                 -- 是否已锁定，0-否，1-是
    `deny`     INTEGER DEFAULT 0,                 -- 连续错误登陆次数，超过3次则锁定用户
    `history`  TEXT    DEFAULT '',                -- 登录历史 [{"ip": "localhost", "time": 1709211720}]
    `add_time` INTEGER DEFAULT 0,                 -- 创建时间（时间戳，单位s）
    `upd_time` INTEGER DEFAULT 0                  -- 修改时间（时间戳，单位s）
);


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
    `size`     INTEGER DEFAULT 0,                 -- 大小，单位：byte
    `del`      INTEGER DEFAULT 0,                 -- 是否已删除，0-否，1-是
    `add_time` INTEGER DEFAULT 0,                 -- 创建时间（时间戳，单位s）
    `upd_time` INTEGER DEFAULT 0                  -- 修改时间（时间戳，单位s）
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
    `size`     INTEGER DEFAULT 0,                 -- 大小，单位：byte
    `del`      INTEGER DEFAULT 0,                 -- 是否已删除，0-否，1-是
    `add_time` INTEGER DEFAULT 0,                 -- 创建时间（时间戳，单位s）
    `upd_time` INTEGER DEFAULT 0                  -- 修改时间（时间戳，单位s）
);
