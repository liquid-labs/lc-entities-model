CREATE TABLE `locations` (
  `id` int(10) NOT NULL auto_increment,
  `address1` varchar(255) NOT NULL,
  `address2` varchar(255),
  `city` varchar(255) NOT NULL,
  `state` varchar(2) NOT NULL,
  `zip` varchar(12) NOT NULL,
  `lat` decimal(9,7) NOT NULL,
  `lng` decimal(10,7) NOT NULL,
  CONSTRAINT `locations_key` PRIMARY KEY ( `id` )
);

CREATE TABLE `entity_addresses` (
  `entity_id` int(10) NOT NULL,
  `location_id` int(10) NOT NULL,
  `idx` int(2) NOT NULL,
  `label` varchar(64),
  CONSTRAINT `entity_locations_key` PRIMARY KEY (`entity_id`, `location_id`, `idx`),
  CONSTRAINT `entity_locations_refs_entities` FOREIGN KEY ( `entity_id` ) REFERENCES `entities` ( `id` ),
  CONSTRAINT `entity_locations_refs_locations` FOREIGN KEY ( `location_id` ) REFERENCES `locations` ( `id` )
);
