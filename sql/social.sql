CREATE DATABASE social;
USE social;

CREATE TABLE `friends` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
   `user_index` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
   `friend_index` BIGINT(20) unsigned NOT NULL COMMENT '好友主键',
   `comment` varchar(20) DEFAULT NULL COMMENT '好友备注',
   `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `idx_user_friend` (`user_index`, `friend_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表(冗余存储)';

CREATE TABLE `friend_requests` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
   `req_id` bigint(20) unsigned NOT NULL COMMENT '好友请求ID',
   `user_index` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
   `friend_index` BIGINT(20) unsigned NOT NULL COMMENT '好友主键',
   `req_msg` varchar(255) DEFAULT NULL COMMENT '申请信息',
   `req_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
   `handle_result` tinyint(1) DEFAULT 0 COMMENT '处理结果:0-待处理，1-同意，2-拒绝',
   `handle_msg` varchar(255) DEFAULT NULL COMMENT '处理信息',
   `handle_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '处理时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uidx_req_id` (`req_id`),
   KEY `idx_user_result` (`user_index`, `handle_result`),
   KEY `idx_friend_result` (`friend_index`, `handle_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';

CREATE TABLE `groups` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
   `group_id` BIGINT(20) unsigned NOT NULL COMMENT '群聊Id',
   `group_name` varchar(100) NOT NULL DEFAULT '' COMMENT '群聊名称',
   `user_index` BIGINT(20) unsigned NOT NULL COMMENT '群主主键',
   `group_status` tinyint(1) DEFAULT 0 COMMENT '入群状态:0-正常，1-禁言',
   `join_status` tinyint(1) DEFAULT 0 COMMENT '入群方式:0-开放，1-申请',
   `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uidx_group_id` (`group_id`),
   KEY `idx_user_index` (`user_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊表';

CREATE  TABLE `group_members`(
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
   `group_index` BIGINT(20) unsigned NOT NULL COMMENT '群聊主键',
   `user_index` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
   `user_type` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '成员类型:0-普通成员,1-管理员,2-群主',
   `user_status` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '状态：0-正常，1-禁言',
   `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
   CONSTRAINT `chk_user_type` CHECK (`user_type` IN (0, 1, 2)),
   CONSTRAINT `chk_user_status` CHECK (`user_status` IN (0, 1)),
   PRIMARY KEY (`id`),
   KEY `idx_group_index` (`group_index`),
   KEY `idx_user_index` (`user_index`),
   UNIQUE KEY `idx_group_user` (`group_index`, `user_index`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊成员表';

CREATE TABLE `group_requests` (
   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
   `user_index` BIGINT(20) unsigned NOT NULL COMMENT '用户主键',
   `group_index` BIGINT(20) unsigned NOT NULL COMMENT '群聊主键',
   `req_msg` varchar(255) DEFAULT NULL COMMENT '申请信息',
   `req_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
   `handler_index` BIGINT(20) unsigned NOT NULL COMMENT '处理人主键',
   `handle_result` tinyint(1) DEFAULT 0 COMMENT '处理结果:0-待处理，1-同意，2-拒绝',
   `handle_msg` varchar(255) DEFAULT NULL COMMENT '处理信息',
   `handle_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '处理时间',
   PRIMARY KEY (`id`),
   KEY `idx_user_group` (`user_index`, `group_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊申请表';