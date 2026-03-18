CREATE TABLE IF NOT EXISTS `ParentContext` (
  `context_id` int NOT NULL,
  `parent_context_id` int NOT NULL,
  PRIMARY KEY (`context_id`,`parent_context_id`),
  KEY `idx_parentcontext_parent_context_id` (`parent_context_id`)
);
