ALTER TABLE `agents`
DROP FOREIGN KEY fk_agents_user_id;

ALTER TABLE `agents`
DROP INDEX uq_agents_user_id_and_name;

DROP TABLE IF EXISTS `agents`;
