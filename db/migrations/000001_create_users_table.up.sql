CREATE TABLE IF NOT EXISTS `users` (
  `id` BINARY(16) NOT NULL COMMENT "ID",
  `name` VARCHAR(255) NOT NULL COMMENT "ユーザー名",
  `password` VARCHAR(60) NOT NULL COMMENT "パスワード",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  `deleted_at` DATETIME (6) COMMENT "削除日時",
  PRIMARY KEY (`id`),
  UNIQUE uq_users_name (`name`)
);
