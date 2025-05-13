CREATE TABLE IF NOT EXISTS `Attribution` (
  `id` int NOT NULL AUTO_INCREMENT,
  `context_id` int NOT NULL,
  `artifact_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `context_id` (`context_id`,`artifact_id`)
);
