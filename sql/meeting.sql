CREATE DATABASE meeting;
USE meeting;

CREATE  TABLE `meeting_info`(
    `id` BIGINT(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `meeting_id` BIGINT(12) unsigned NOT NULL DEFAULT 0 COMMENT '会议号',
    `meeting_name` varchar(100) NOT NULL DEFAULT '' COMMENT '会议名称',
    `user_id` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
    `status` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-空闲，1-进行中',
    `join_type` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '加入方式：0-公开，1-私密',
    `meeting_password` varchar(5) NOT NULL DEFAULT '' COMMENT '会议密码',
    `start_time` timestamp NULL DEFAULT NULL COMMENT '会议开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '会议结束时间',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    CONSTRAINT `chk_meeting_status` CHECK (`status` IN (0, 1)),
    CONSTRAINT `chk_join_type` CHECK (`join_type` IN (0, 1)),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_meeting_id` (`meeting_id`),
    UNIQUE KEY `uniq_user_id` (`user_id`),
    KEY `ix_start_time` (`start_time`),
    KEY `ix_end_time` (`end_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会议信息表';


CREATE  TABLE `meeting_member`(
    `id` BIGINT(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `meeting_id` BIGINT(20) unsigned NOT NULL COMMENT '会议主键',
    `user_id` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
    `user_type` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '成员类型:0-普通成员,1-主持人',
    `user_status` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-正常，1-拉黑',
    `last_join_time` timestamp NULL DEFAULT NULL COMMENT '上次加入时间',
    CONSTRAINT `chk_user_type` CHECK (`user_type` IN (0, 1)),
    CONSTRAINT `chk_user_status` CHECK (`user_status` IN (0, 1)),
    PRIMARY KEY (`id`),
    KEY `iq_meeting_id` (`meeting_id`),
    KEY `iq_user_id` (`user_id`),
    UNIQUE KEY `idx_meeting_user` (`meeting_id`, `user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会议成员表';
