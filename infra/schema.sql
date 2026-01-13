CREATE DATABASE IF NOT EXISTS `moxuevideo`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE moxuevideo;

SET NAMES utf8mb4;
SET collation_connection = 'utf8mb4_0900_ai_ci';

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` VARCHAR(64) NOT NULL COMMENT '用户名',
  `email` VARCHAR(255) NULL COMMENT '邮箱',
  `phone` VARCHAR(32) NULL COMMENT '手机号',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `avatar_url` VARCHAR(512) NULL COMMENT '头像URL',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态',
  `password_updated_at` DATETIME(3) NULL COMMENT '密码更新时间',
  `last_login_at` DATETIME(3) NULL COMMENT '最近登录时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`),
  UNIQUE KEY `uk_users_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `user_sessions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `refresh_token_hash` CHAR(64) NOT NULL COMMENT '刷新令牌哈希',
  `user_agent` VARCHAR(255) NULL COMMENT 'User-Agent',
  `ip` VARCHAR(45) NULL COMMENT 'IP地址',
  `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
  `revoked_at` DATETIME(3) NULL COMMENT '撤销时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_sessions_refresh_token_hash` (`refresh_token_hash`),
  KEY `idx_user_sessions_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_user_sessions_expires_at` (`expires_at`),
  CONSTRAINT `fk_user_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话(刷新令牌)表';

CREATE TABLE IF NOT EXISTS `password_reset_tokens` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `token_hash` CHAR(64) NOT NULL COMMENT '重置令牌哈希',
  `expires_at` DATETIME(3) NOT NULL COMMENT '过期时间',
  `used_at` DATETIME(3) NULL COMMENT '使用时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_password_reset_tokens_token_hash` (`token_hash`),
  KEY `idx_password_reset_tokens_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_password_reset_tokens_expires_at` (`expires_at`),
  CONSTRAINT `fk_password_reset_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='密码重置令牌表';

CREATE TABLE IF NOT EXISTS `follows` (
  `follower_id` BIGINT UNSIGNED NOT NULL COMMENT '关注者用户ID',
  `followee_id` BIGINT UNSIGNED NOT NULL COMMENT '被关注者用户ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`follower_id`, `followee_id`),
  KEY `idx_follows_followee_id_created_at` (`followee_id`, `created_at`),
  CONSTRAINT `fk_follows_follower` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_follows_followee` FOREIGN KEY (`followee_id`) REFERENCES `users` (`id`),
  CONSTRAINT `ck_follows_not_self` CHECK (`follower_id` <> `followee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='关注关系表';

CREATE TABLE IF NOT EXISTS `videos` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者用户ID',
  `title` VARCHAR(200) NOT NULL COMMENT '标题',
  `cover_url` VARCHAR(512) NULL COMMENT '封面URL',
  `play_url` VARCHAR(512) NOT NULL COMMENT '播放URL',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态',
  `published_at` DATETIME(3) NULL COMMENT '发布时间',
  `like_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `favorite_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏数',
  `comment_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数',
  `share_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分享数',
  `view_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '播放量',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_videos_author_id_published_at` (`author_id`, `published_at`),
  KEY `idx_videos_published_at` (`published_at`),
  CONSTRAINT `fk_videos_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频表';

CREATE TABLE IF NOT EXISTS `video_likes` (
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '点赞时间',
  PRIMARY KEY (`user_id`, `video_id`),
  KEY `idx_video_likes_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_likes_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频点赞表';

CREATE TABLE IF NOT EXISTS `video_favorites` (
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '收藏时间',
  PRIMARY KEY (`user_id`, `video_id`),
  KEY `idx_video_favorites_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_favorites_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频收藏表';

CREATE TABLE IF NOT EXISTS `video_shares` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
  `share_channel` VARCHAR(32) NULL COMMENT '分享渠道',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_video_shares_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_video_shares_video_id_created_at` (`video_id`, `created_at`),
  CONSTRAINT `fk_video_shares_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_video_shares_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频分享记录表';

CREATE TABLE IF NOT EXISTS `watch_histories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
  `watched_at` DATETIME(3) NOT NULL COMMENT '观看时间',
  `watch_seconds` INT UNSIGNED NULL COMMENT '观看秒数',
  `device` VARCHAR(64) NULL COMMENT '设备标识',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_watch_histories_user_id_watched_at` (`user_id`, `watched_at`),
  KEY `idx_watch_histories_video_id_watched_at` (`video_id`, `watched_at`),
  CONSTRAINT `fk_watch_histories_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_watch_histories_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='观看历史表';

CREATE TABLE IF NOT EXISTS `comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `parent_id` BIGINT UNSIGNED NULL COMMENT '父评论ID',
  `root_id` BIGINT UNSIGNED NULL COMMENT '根评论ID',
  `content` TEXT NOT NULL COMMENT '内容',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态',
  `like_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_comments_video_id_created_at` (`video_id`, `created_at`),
  KEY `idx_comments_user_id_created_at` (`user_id`, `created_at`),
  KEY `idx_comments_parent_id_created_at` (`parent_id`, `created_at`),
  CONSTRAINT `fk_comments_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表';

CREATE TABLE IF NOT EXISTS `comment_likes` (
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `comment_id` BIGINT UNSIGNED NOT NULL COMMENT '评论ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '点赞时间',
  PRIMARY KEY (`user_id`, `comment_id`),
  KEY `idx_comment_likes_comment_id_created_at` (`comment_id`, `created_at`),
  CONSTRAINT `fk_comment_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_comment_likes_comment` FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论点赞表';

CREATE TABLE IF NOT EXISTS `dm_threads` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_a_id` BIGINT UNSIGNED NOT NULL COMMENT '会话用户A ID',
  `user_b_id` BIGINT UNSIGNED NOT NULL COMMENT '会话用户B ID',
  `user_low_id` BIGINT UNSIGNED GENERATED ALWAYS AS (LEAST(`user_a_id`, `user_b_id`)) STORED COMMENT '会话较小用户ID(生成列)',
  `user_high_id` BIGINT UNSIGNED GENERATED ALWAYS AS (GREATEST(`user_a_id`, `user_b_id`)) STORED COMMENT '会话较大用户ID(生成列)',
  `last_message_id` BIGINT UNSIGNED NULL COMMENT '最后一条消息ID',
  `last_message_at` DATETIME(3) NULL COMMENT '最后消息时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dm_threads_pair` (`user_low_id`, `user_high_id`),
  KEY `idx_dm_threads_user_a_id` (`user_a_id`),
  KEY `idx_dm_threads_user_b_id` (`user_b_id`),
  CONSTRAINT `fk_dm_threads_user_a` FOREIGN KEY (`user_a_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_dm_threads_user_b` FOREIGN KEY (`user_b_id`) REFERENCES `users` (`id`),
  CONSTRAINT `ck_dm_threads_not_self` CHECK (`user_a_id` <> `user_b_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='私信会话表';

CREATE TABLE IF NOT EXISTS `dm_messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `thread_id` BIGINT UNSIGNED NOT NULL COMMENT '会话ID',
  `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者用户ID',
  `msg_type` TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型',
  `content` TEXT NULL COMMENT '内容',
  `payload_json` JSON NULL COMMENT '消息扩展数据(JSON)',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_dm_messages_thread_id_id` (`thread_id`, `id`),
  KEY `idx_dm_messages_sender_id_created_at` (`sender_id`, `created_at`),
  CONSTRAINT `fk_dm_messages_thread` FOREIGN KEY (`thread_id`) REFERENCES `dm_threads` (`id`),
  CONSTRAINT `fk_dm_messages_sender` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='私信消息表';

CREATE TABLE IF NOT EXISTS `user_notifications` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `type` VARCHAR(32) NOT NULL COMMENT '通知类型',
  `title` VARCHAR(200) NULL COMMENT '标题',
  `content` TEXT NULL COMMENT '内容',
  `data_json` JSON NULL COMMENT '通知数据(JSON)',
  `is_read` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读',
  `read_at` DATETIME(3) NULL COMMENT '阅读时间',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_notifications_user_id_is_read_created_at` (`user_id`, `is_read`, `created_at`),
  KEY `idx_user_notifications_user_id_created_at` (`user_id`, `created_at`),
  CONSTRAINT `fk_user_notifications_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户通知表';