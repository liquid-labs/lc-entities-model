CREATE TABLE locations (
  id BIGSERIAL,
  address1 VARCHAR(255) NOT NULL,
  address2 VARCHAR(255),
  city VARCHAR(255) NOT NULL,
  state VARCHAR(2) NOT NULL,
  zip VARCHAR(12) NOT NULL,
  lat DECIMAL(9,7) NOT NULL,
  lng DECIMAL(10,7) NOT NULL,
  CONSTRAINT locations_key PRIMARY KEY ( id )
);

CREATE TABLE entity_addresses (
  entity_id BIGINT NOT NULL,
  location_id BIGINT NOT NULL,
  -- MySQL
  -- idx INT(2) NOT NULL,
  -- Postgres
  idx INT NOT NULL CHECK(idx >= 0 AND idx < 100),
  label VARCHAR(64),
  CONSTRAINT entity_locations_key PRIMARY KEY (entity_id, location_id, idx),
  CONSTRAINT entity_locations_refs_entities FOREIGN KEY ( entity_id ) REFERENCES entities ( id ),
  CONSTRAINT entity_locations_refs_locations FOREIGN KEY ( location_id ) REFERENCES locations ( id )
);
