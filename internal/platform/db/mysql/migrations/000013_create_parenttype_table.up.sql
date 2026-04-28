CREATE TABLE IF NOT EXISTS `ParentType` (
  `type_id` int NOT NULL,
  `parent_type_id` int NOT NULL,
  PRIMARY KEY (`type_id`,`parent_type_id`)
);
