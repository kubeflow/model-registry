CREATE TABLE IF NOT EXISTS `TypeProperty` (
  `type_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `data_type` int DEFAULT NULL,
  PRIMARY KEY (`type_id`,`name`)
);
