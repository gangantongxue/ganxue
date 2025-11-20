-- 创建数据库
CREATE DATABASE IF NOT EXISTS ganxue_server;

-- 使用数据库
USE ganxue_server;

-- 创建用户表（示例，根据你的实际需求调整）
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_name` longtext DEFAULT NULL,
  `password` longtext DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `streak_days` bigint DEFAULT NULL,
  `total_days` bigint DEFAULT NULL,
  `go_last_chapter` longtext DEFAULT NULL,
  `c_last_chapter` longtext DEFAULT NULL,
  `cpp_last_chapter` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 授权root用户远程访问（Docker环境需要）
GRANT ALL PRIVILEGES ON ganxue_server.* TO 'root'@'%';
FLUSH PRIVILEGES;