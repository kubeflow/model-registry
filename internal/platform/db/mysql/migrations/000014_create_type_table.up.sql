CREATE TABLE IF NOT EXISTS `Type` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `type_kind` tinyint(1) NOT NULL,
  `description` text,
  `input_type` text,
  `output_type` text,
  `external_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `external_id` (`external_id`),
  KEY `idx_type_name` (`name`),
  KEY `idx_type_external_id` (`external_id`)
);
