CREATE DATABASE IF NOT EXISTS `moxuevideo_chat`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE moxuevideo_chat;

SET NAMES utf8mb4;
SET collation_connection = 'utf8mb4_0900_ai_ci';

CREATE TABLE IF NOT EXISTS `chat` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `from_user_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者用户ID',
  `to_user_id` BIGINT UNSIGNED NOT NULL COMMENT '接收者用户ID',
  `msg_type` ENUM('text','picture','audio') NOT NULL DEFAULT 'text' COMMENT '消息类型(text/picture/audio)',
  `content` TEXT NOT NULL COMMENT '消息内容',
  `is_read` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读(0未读/1已读)',
  `read_time` VARCHAR(20) NULL COMMENT '阅读时间',
  `create_time` VARCHAR(20) NOT NULL COMMENT '发送时间',
  `unique` VARCHAR(64) NOT NULL COMMENT '去重唯一标识',
  `is_del` TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除(软删 0正常/1删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chat_unique` (`unique`),
  KEY `idx_chat_to_isread_ct` (`to_user_id`, `is_read`, `create_time`),
  KEY `idx_chat_from_to_ct` (`from_user_id`, `to_user_id`, `create_time`),
  KEY `idx_chat_to_from_ct` (`to_user_id`, `from_user_id`, `create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天消息表';

GRANT ALL PRIVILEGES ON `moxuevideo_chat`.* TO 'moxue'@'%';
FLUSH PRIVILEGES;