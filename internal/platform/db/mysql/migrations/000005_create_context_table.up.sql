CREATE TABLE IF NOT EXISTS `Context` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `external_id` varchar(255) DEFAULT NULL,
  `create_time_since_epoch` bigint NOT NULL DEFAULT '0',
  `last_update_time_since_epoch` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `type_id` (`type_id`,`name`),
  UNIQUE KEY `external_id` (`external_id`),
  KEY `idx_context_create_time_since_epoch` (`create_time_since_epoch`),
  KEY `idx_context_last_update_time_since_epoch` (`last_update_time_since_epoch`),
  KEY `idx_context_external_id` (`external_id`)
);
