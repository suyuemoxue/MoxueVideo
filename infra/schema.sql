CREATE DATABASE IF NOT EXISTS `moxuevideo`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE `moxuevideo`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL,
  `email` VARCHAR(255) NULL,
  `phone` VARCHAR(32) NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `avatar_url` VARCHAR(512) NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `password_updated_at` DATETIME(3) NULL,
  `last_login_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`),
  UNIQUE KEY `uk_users_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_sessions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `refresh_token_hash` CHAR(64) NOT NULL,
  `user_agent` VARCHAR(255) NULL,
  `ip` VARCHAR(45) NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `revoked_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_sessions_refresh_token_hash` (`refresh_token_hash`),
  KEY `idx_user_sessions_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_user_sessions_expires_at` (`expires_at`),
  CONSTRAINT `fk_user_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `password_reset_tokens` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `token_hash` CHAR(64) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `used_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_password_reset_tokens_token_hash` (`token_hash`),
  KEY `idx_password_reset_tokens_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_password_reset_tokens_expires_at` (`expires_at`),
  CONSTRAINT `fk_password_reset_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `follows` (
  `follower_id` BIGINT UNSIGNED NOT NULL,
  `followee_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`follower_id`, `followee_id`),
  KEY `idx_follows_followee_id_created_at` (`followee_id`, `created_at`),
  CONSTRAINT `fk_follows_follower` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_follows_followee` FOREIGN KEY (`followee_id`) REFERENCES `users` (`id`),
  CONSTRAINT `ck_follows_not_self` CHECK (`follower_id` <> `followee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `videos` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `author_id` BIGINT UNSIGNED NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `cover_url` VARCHAR(512) NULL,
  `play_url` VARCHAR(512) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `published_at` DATETIME(3) NULL,
  `like_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `favorite_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `comment_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `share_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `view_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_videos_author_id_published_at` (`author_id`, `published_at`),
  KEY `idx_videos_published_at` (`published_at`),
  CONSTRAINT `fk_videos_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `video_likes` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `video_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `video_id`),
  KEY `idx_video_likes_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_likes_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `video_favorites` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `video_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `video_id`),
  KEY `idx_video_favorites_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_favorites_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `video_shares` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `video_id` BIGINT UNSIGNED NOT NULL,
  `share_channel` VARCHAR(32) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_video_shares_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_video_shares_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_shares_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `watch_histories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `video_id` BIGINT UNSIGNED NOT NULL,
  `watched_at` DATETIME(3) NOT NULL,
  `watch_seconds` INT UNSIGNED NULL,
  `device` VARCHAR(64) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_watch_histories_user_id_watched_at` (`user_id`, `watched_at`),
  KEY `idx_watch_histories_video_id_watched_at` (`video_id`, `watched_at`),
  CONSTRAINT `fk_watch_histories_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_watch_histories_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `video_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `parent_id` BIGINT UNSIGNED NULL,
  `root_id` BIGINT UNSIGNED NULL,
  `content` TEXT NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `like_count` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_comments_video_id_created_at` (`video_id`, `created_at`),
  KEY `idx_comments_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_comments_parent_id_created_at` (`parent_id`, `created_at`),
  CONSTRAINT `fk_comments_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `comment_likes` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `comment_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `comment_id`),
  KEY `idx_comment_likes_comment_id_created_at` (`comment_id`, `created_at`),
  CONSTRAINT `fk_comment_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_comment_likes_comment` FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `dm_threads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_a_id` BIGINT UNSIGNED NOT NULL,
  `user_b_id` BIGINT UNSIGNED NOT NULL,
  `user_low_id` BIGINT UNSIGNED GENERATED ALWAYS AS (LEAST(`user_a_id`, `user_b_id`)) STORED,
  `user_high_id` BIGINT UNSIGNED GENERATED ALWAYS AS (GREATEST(`user_a_id`, `user_b_id`)) STORED,
  `last_message_id` BIGINT UNSIGNED NULL,
  `last_message_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_threads_pair` (`user_low_id`, `user_high_id`),
  KEY `idx_dm_threads_user_a_id` (`user_a_id`),
  KEY `idx_dm_threads_user_b_id` (`user_b_id`),
  CONSTRAINT `fk_dm_threads_user_a` FOREIGN KEY (`user_a_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_dm_threads_user_b` FOREIGN KEY (`user_b_id`) REFERENCES `users` (`id`),
  CONSTRAINT `ck_dm_threads_not_self` CHECK (`user_a_id` <> `user_b_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `dm_messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `thread_id` BIGINT UNSIGNED NOT NULL,
  `sender_id` BIGINT UNSIGNED NOT NULL,
  `msg_type` TINYINT NOT NULL DEFAULT 1,
  `content` TEXT NULL,
  `payload_json` JSON NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_dm_messages_thread_id_id` (`thread_id`, `id`),
  KEY `idx_dm_messages_sender_id_created_at` (`sender_id`, `created_at`),
  CONSTRAINT `fk_dm_messages_thread` FOREIGN KEY (`thread_id`) REFERENCES `dm_threads` (`id`),
  CONSTRAINT `fk_dm_messages_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `user_notifications` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `type` VARCHAR(32) NOT NULL,
  `title` VARCHAR(200) NULL,
  `content` TEXT NULL,
  `data_json` JSON NULL,
  `is_read` TINYINT NOT NULL DEFAULT 0,
  `read_at` DATETIME(3) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_user_notifications_user_id_is_read_created_at` (`user_id`, `is_read`, `created_at`),
  KEY `idx_user_notifications_user_id_created_at` (`user_id`, `created_at`),
  CONSTRAINT `fk_user_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

