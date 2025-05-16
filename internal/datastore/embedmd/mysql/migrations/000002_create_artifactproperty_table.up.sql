CREATE TABLE IF NOT EXISTS `ArtifactProperty` (
  `artifact_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_custom_property` tinyint(1) NOT NULL,
  `int_value` int DEFAULT NULL,
  `double_value` double DEFAULT NULL,
  `string_value` mediumtext,
  `byte_value` mediumblob,
  `proto_value` mediumblob,
  `bool_value` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`artifact_id`,`name`,`is_custom_property`),
  KEY `idx_artifact_property_int` (`name`,`is_custom_property`,`int_value`),
  KEY `idx_artifact_property_double` (`name`,`is_custom_property`,`double_value`),
  KEY `idx_artifact_property_string` (`name`,`is_custom_property`,`string_value`(255))
);
