CREATE DATABASE user;
USE user;

CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint(12) unsigned NOT NULL DEFAULT '' COMMENT '用户ID',
    `username` varchar(20) NOT NULL DEFAULT '' COMMENT '用户昵称',
    `password` varchar(32) NOT NULL DEFAULT '' COMMENT '用户密码，MD5加密',
    `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
    `meeting_id` varchar(12) NOT NULL DEFAULT '' COMMENT '个人会议号',
    `last_login_time` timestamp NULL DEFAULT NULL COMMENT '上次登录时间',
    `last_off_time` timestamp NULL DEFAULT NULL COMMENT '上次下线时间',
    `sex` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    `status` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '状态：1-禁用，0-正常',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_user_id` (`user_id`),
    UNIQUE KEY `uniq_username` (`username`),
    UNIQUE KEY `uniq_email` (`email`),
    UNIQUE KEY `uniq_meeting_id` (`meeting_id`),
    KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';