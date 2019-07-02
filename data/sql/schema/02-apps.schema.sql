CREATE TABLE apps (
  `id` int(10) NOT NULL auto_increment,
  `name` VARCHAR(128) NOT NULL,
  `type` VARCHAR(64) NOT NULL,
  CONSTRAINT `apps_key` PRIMARY KEY ( `id` ),
  CONSTRAINT `apps_ref_entities` FOREIGN KEY ( `id` ) REFERENCES `entities` ( `id` )
);
