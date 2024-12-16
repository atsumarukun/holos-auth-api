ALTER TABLE `permissions`
DROP FOREIGN KEY fk_permissions_agent_id;

ALTER TABLE `permissions`
DROP FOREIGN KEY fk_permissions_policy_id;

DROP TABLE IF EXISTS `permissions`;
