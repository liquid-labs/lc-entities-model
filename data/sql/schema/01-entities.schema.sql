-- Postgres
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- MySQL:
-- pub_id note: It seems that we set 'pub_id's to 'NOT NULL' and it's OK when
-- we run the setup scripts, but when the app tries to insert, it causes errors.
-- last confirmed 2018-07-13
CREATE TABLE entities (
  id                BIGSERIAL,
  -- MySQL
  -- pub_id            CHAR(36) NOT NULL,
  -- Postgres
  pub_id            UUID NOT NULL DEFAULT uuid_generate_v4(),
  owner_id          BIGINT, -- should only be null for Persons
  publicly_readable BOOLEAN,
  containers        BIGINT[],
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  last_updated      TIMESTAMPTZ,
  deleted_at        TIMESTAMPTZ,
  CONSTRAINT entities_key PRIMARY KEY ( id ),
  CONSTRAINT entities_pub_id_unique UNIQUE (pub_id),
  CONSTRAINT entities_owner_refs_entities FOREIGN KEY (owner_id) REFERENCES entities ( id )
);

CREATE UNIQUE INDEX entities_pub_id_index ON entities (pub_id);

-- MySQL
-- DELIMITER //
-- CREATE TRIGGER `entities_public_id`
--   BEFORE INSERT ON entities FOR EACH ROW
--     BEGIN
--       IF new.pub_id IS NULL THEN
--         SET new.pub_id=UPPER(UUID());
--       ELSIF new.pub_id NOT SIMILAR TO '^[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$' THEN
--        RAISE EXCEPTION 'Invalid UUID format.' USING ERRCODE
--      END IF;
--    END//

-- CREATE TRIGGER `entities_last_updated`
--   BEFORE INSERT OR UPDATE ON entities FOR EACH ROW
--     SET new.last_updated=UNIX_TIMESTAMP();//
-- DELIMITER ;

-- Postgres
CREATE OR REPLACE FUNCTION trigger_entities_last_updated()
  RETURNS TRIGGER AS '
BEGIN
  NEW.last_updated := NOW();
  RETURN NEW;
END' LANGUAGE 'plpgsql';

DROP EVENT TRIGGER IF EXISTS entities_last_updated;

CREATE TRIGGER entities_last_updated
  BEFORE INSERT OR UPDATE ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_entities_last_updated();

CREATE TABLE containers (
  id  BIGINT,
  CONSTRAINT containers_key PRIMARY KEY ( id ),
  CONSTRAINT containers_ref_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);
