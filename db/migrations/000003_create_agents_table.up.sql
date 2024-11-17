CREATE TABLE IF NOT EXISTS `agents` (
  `id` CHAR(36) NOT NULL COMMENT "ID",
  `user_id` CHAR(36) NOT NULL COMMENT "ユーザーID",
  `name` VARCHAR(255) NOT NULL COMMENT "エージェント名",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  `deleted_at` DATETIME (6) COMMENT "削除日時",
  PRIMARY KEY (`id`),
  UNIQUE uq_agents_user_id_and_name (`user_id`, `name`),
  CONSTRAINT fk_agents_user_id FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);