CREATE TABLE IF NOT EXISTS `user_tokens` (
  `user_id` CHAR(36) NOT NULL COMMENT "ユーザーID",
  `token` CHAR(32) NOT NULL COMMENT "トークン",
  `expires_at` DATETIME (6) NOT NULL COMMENT "有効期限",
  PRIMARY KEY (`user_id`),
  UNIQUE uq_user_tokens_token (`token`),
  CONSTRAINT fk_user_tokens_user_id FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
