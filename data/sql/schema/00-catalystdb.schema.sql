CREATE TABLE catalystdb (
  `update` INT(11),
  `schema_hash` CHAR(64),
  `version_spec` TEXT,
  CONSTRAINT `catalystdb_key` PRIMARY KEY ( `update`,`schema_hash` )
);

DELIMITER //
CREATE TRIGGER `catalystdb_updateed`
  BEFORE INSERT ON catalystdb FOR EACH ROW
    SET new.update=UNIX_TIMESTAMP();//
DELIMITER ;
