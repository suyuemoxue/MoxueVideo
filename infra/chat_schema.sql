CREATE DATABASE IF NOT EXISTS `moxuevideo_chat`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE `moxuevideo_chat`;

CREATE TABLE IF NOT EXISTS `dm_threads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_low_id` BIGINT UNSIGNED NOT NULL,
  `user_high_id` BIGINT UNSIGNED NOT NULL,
  `last_message_id` BIGINT UNSIGNED NULL,
  `last_message_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_threads_pair` (`user_low_id`, `user_high_id`),
  KEY `idx_dm_threads_last_message_at` (`last_message_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `dm_messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `thread_id` BIGINT UNSIGNED NOT NULL,
  `sender_id` BIGINT UNSIGNED NOT NULL,
  `receiver_id` BIGINT UNSIGNED NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_dm_messages_thread_id_id` (`thread_id`, `id`),
  KEY `idx_dm_messages_receiver_id_id` (`receiver_id`, `id`),
  KEY `idx_dm_messages_sender_id_created_at` (`sender_id`, `created_at`),
  CONSTRAINT `fk_dm_messages_thread` FOREIGN KEY (`thread_id`) REFERENCES `dm_threads` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `dm_message_reads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `message_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `read_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_message_reads_message_user` (`message_id`, `user_id`),
  KEY `idx_dm_message_reads_user_id_read_at` (`user_id`, `read_at`),
  CONSTRAINT `fk_dm_message_reads_message` FOREIGN KEY (`message_id`) REFERENCES `dm_messages` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

