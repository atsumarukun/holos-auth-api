ALTER TABLE `agent_tokens`
DROP FOREIGN KEY fk_agent_tokens_agent_id;

ALTER TABLE `agent_tokens`
DROP INDEX uq_agent_tokens_token;

DROP TABLE IF EXISTS `agent_tokens`;
