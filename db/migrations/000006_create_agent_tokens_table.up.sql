CREATE TABLE IF NOT EXISTS `agent_tokens` (
  `agent_id` CHAR(36) NOT NULL COMMENT "エージェントID",
  `token` CHAR(32) NOT NULL COMMENT "トークン",
  PRIMARY KEY (`agent_id`),
  UNIQUE uq_agent_tokens_token (`token`),
  CONSTRAINT fk_agent_tokens_agent_id FOREIGN KEY (`agent_id`) REFERENCES `agents` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
