CREATE DATABASE IF NOT EXISTS `moxuevideo_chat`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE moxuevideo_chat;

SET NAMES utf8mb4;
SET collation_connection = 'utf8mb4_0900_ai_ci';

CREATE TABLE IF NOT EXISTS `chat` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `from_user_id` BIGINT UNSIGNED NOT NULL,
  `to_user_id` BIGINT UNSIGNED NOT NULL,
  `msg_type` ENUM('text','picture','audio') NOT NULL DEFAULT 'text',
  `content` TEXT NOT NULL,
  `is_read` TINYINT NOT NULL DEFAULT 0,
  `uniqued` VARCHAR(64) NOT NULL,
  `create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `read_time` DATETIME(3) NULL,
  `is_del` TINYINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_chat_uniqued` (`uniqued`),
  KEY `idx_chat_to_isread_ct` (`to_user_id`, `is_read`, `create_time`),
  KEY `idx_chat_from_to_ct` (`from_user_id`, `to_user_id`, `create_time`),
  KEY `idx_chat_to_from_ct` (`to_user_id`, `from_user_id`, `create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

GRANT ALL PRIVILEGES ON `moxuevideo_chat`.* TO 'moxue'@'%';
FLUSH PRIVILEGES;