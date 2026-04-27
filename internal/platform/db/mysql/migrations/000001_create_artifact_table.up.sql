CREATE TABLE IF NOT EXISTS `Artifact` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type_id` int NOT NULL,
  `uri` text,
  `state` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `external_id` varchar(255) DEFAULT NULL,
  `create_time_since_epoch` bigint NOT NULL DEFAULT '0',
  `last_update_time_since_epoch` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `external_id` (`external_id`),
  UNIQUE KEY `UniqueArtifactTypeName` (`type_id`,`name`),
  KEY `idx_artifact_uri` (`uri`(255)),
  KEY `idx_artifact_create_time_since_epoch` (`create_time_since_epoch`),
  KEY `idx_artifact_last_update_time_since_epoch` (`last_update_time_since_epoch`),
  KEY `idx_artifact_external_id` (`external_id`)
);
