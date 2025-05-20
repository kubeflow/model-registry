CREATE TABLE IF NOT EXISTS `EventPath` (
  `event_id` int NOT NULL,
  `is_index_step` tinyint(1) NOT NULL,
  `step_index` int DEFAULT NULL,
  `step_key` text,
  KEY `idx_eventpath_event_id` (`event_id`)
);
