-- Postgres
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE entities (
  id                UUID NOT NULL DEFAULT uuid_generate_v4(),
  resource_name     VARCHAR(128) NOT NULL,
  name              VARCHAR(128),
  description       TEXT,
  owner_id          UUID NOT NULL, -- should only be used for Persons (?)
  publicly_readable BOOLEAN NOT NULL,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_updated      TIMESTAMPTZ NOT NULL,
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

CREATE OR REPLACE FUNCTION trigger_entities_self_own()
  RETURNS TRIGGER AS $$
  BEGIN
    IF NEW.owner_id IS NULL THEN
      NEW.owner_id := NEW.id;
    END IF;
    RETURN new;
  END;
  $$ LANGUAGE 'plpgsql';

CREATE TRIGGER entities_self_own
  BEFORE INSERT ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_entities_self_own();
