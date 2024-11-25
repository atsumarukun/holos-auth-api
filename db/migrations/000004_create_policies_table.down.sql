ALTER TABLE `policies`
DROP FOREIGN KEY fk_policies_user_id;

DROP TABLE IF EXISTS `policies`;
