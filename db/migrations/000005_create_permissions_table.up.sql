CREATE TABLE IF NOT EXISTS `permissions` (
  `agent_id` CHAR(36) NOT NULL COMMENT "エージェントID",
  `policy_id` CHAR(36) NOT NULL COMMENT "ポリシーID",
  `effect` ENUM ("ALLOW", "DENY") NOT NULL COMMENT "効果",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  `deleted_at` DATETIME (6) COMMENT "削除日時",
  PRIMARY KEY (`agent_id`, `policy_id`),
  CONSTRAINT fk_permissions_agent_id FOREIGN KEY (`agent_id`) REFERENCES `agents` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT fk_permissions_policy_id FOREIGN KEY (`policy_id`) REFERENCES `policies` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
