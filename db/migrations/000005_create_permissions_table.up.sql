CREATE TABLE IF NOT EXISTS `permissions` (
  `agent_id` CHAR(36) NOT NULL COMMENT "エージェントID",
  `policy_id` CHAR(36) NOT NULL COMMENT "ポリシーID",
  PRIMARY KEY (`agent_id`, `policy_id`),
  CONSTRAINT fk_permissions_agent_id FOREIGN KEY (`agent_id`) REFERENCES `agents` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT fk_permissions_policy_id FOREIGN KEY (`policy_id`) REFERENCES `policies` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
