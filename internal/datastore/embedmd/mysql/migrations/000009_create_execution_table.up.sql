CREATE TABLE IF NOT EXISTS `Execution` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type_id` int NOT NULL,
  `last_known_state` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `external_id` varchar(255) DEFAULT NULL,
  `create_time_since_epoch` bigint NOT NULL DEFAULT '0',
  `last_update_time_since_epoch` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `external_id` (`external_id`),
  UNIQUE KEY `UniqueExecutionTypeName` (`type_id`,`name`),
  KEY `idx_execution_create_time_since_epoch` (`create_time_since_epoch`),
  KEY `idx_execution_last_update_time_since_epoch` (`last_update_time_since_epoch`),
  KEY `idx_execution_external_id` (`external_id`)
);
