ALTER TABLE `user_tokens`
DROP FOREIGN KEY fk_user_tokens_user_id;

ALTER TABLE `user_tokens`
DROP INDEX uq_user_tokens_token;

DROP TABLE IF EXISTS `user_tokens`;
