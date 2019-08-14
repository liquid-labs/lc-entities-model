-- Postgres
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE entities (
  id                UUID NOT NULL DEFAULT uuid_generate_v4(),
  resource_name     VARCHAR(128),
  name              VARCHAR(128),
  description       TEXT,
  owner_id          UUID, -- should only be null for Persons
  publicly_readable BOOLEAN,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  last_updated      TIMESTAMPTZ,
  deleted_at        TIMESTAMPTZ,
  CONSTRAINT entities_key PRIMARY KEY ( id ),
  CONSTRAINT entities_owner_refs_entities FOREIGN KEY (owner_id) REFERENCES entities ( id )
);

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
