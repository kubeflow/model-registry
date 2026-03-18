CREATE TABLE IF NOT EXISTS `Association` (
  `id` int NOT NULL AUTO_INCREMENT,
  `context_id` int NOT NULL,
  `execution_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `context_id` (`context_id`,`execution_id`)
);
