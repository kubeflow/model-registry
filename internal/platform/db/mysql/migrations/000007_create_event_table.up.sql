CREATE TABLE IF NOT EXISTS `Event` (
  `id` int NOT NULL AUTO_INCREMENT,
  `artifact_id` int NOT NULL,
  `execution_id` int NOT NULL,
  `type` int NOT NULL,
  `milliseconds_since_epoch` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UniqueEvent` (`artifact_id`,`execution_id`,`type`),
  KEY `idx_event_execution_id` (`execution_id`)
);
