CREATE DATABASE IF NOT EXISTS `moxuevideo_chat`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE moxuevideo_chat;

SET NAMES utf8mb4;
SET collation_connection = 'utf8mb4_0900_ai_ci';

CREATE TABLE IF NOT EXISTS `dm_threads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_low_id` BIGINT UNSIGNED NOT NULL COMMENT '会话较小用户ID(生成列)',
  `user_high_id` BIGINT UNSIGNED NOT NULL COMMENT '会话较大用户ID(生成列)',
  `last_message_id` BIGINT UNSIGNED NULL COMMENT '最后一条消息ID',
  `last_message_at` DATETIME(3) NULL COMMENT '最后消息时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_threads_pair` (`user_low_id`, `user_high_id`),
  KEY `idx_dm_threads_last_message_at` (`last_message_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='私信会话表';

CREATE TABLE IF NOT EXISTS `dm_messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `thread_id` BIGINT UNSIGNED NOT NULL COMMENT '会话ID',
  `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者用户ID',
  `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收者用户ID',
  `content` TEXT NOT NULL COMMENT '内容',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_dm_messages_thread_id_id` (`thread_id`, `id`),
  KEY `idx_dm_messages_receiver_id_id` (`receiver_id`, `id`),
  KEY `idx_dm_messages_sender_id_created_at` (`sender_id`, `created_at`),
  CONSTRAINT `fk_dm_messages_thread` FOREIGN KEY (`thread_id`) REFERENCES `dm_threads` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='私信消息表';

CREATE TABLE IF NOT EXISTS `dm_message_reads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `message_id` BIGINT UNSIGNED NOT NULL COMMENT '消息ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `read_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '阅读时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_message_reads_message_user` (`message_id`, `user_id`),
  KEY `idx_dm_message_reads_user_id_read_at` (`user_id`, `read_at`),
  CONSTRAINT `fk_dm_message_reads_message` FOREIGN KEY (`message_id`) REFERENCES `dm_messages` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='私信消息已读表';

GRANT ALL PRIVILEGES ON `moxuevideo_chat`.* TO 'moxue'@'%';
FLUSH PRIVILEGES;
