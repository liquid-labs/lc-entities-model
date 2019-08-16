-- Postgres
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE entities (
  id                UUID NOT NULL DEFAULT uuid_generate_v4(),
  resource_name     VARCHAR(128) NOT NULL,
  name              VARCHAR(128),
  description       TEXT,
  owner_id          UUID NOT NULL, -- should only be used for Persons (?)
  publicly_readable BOOLEAN NOT NULL DEFAULT false,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_updated      TIMESTAMPTZ NOT NULL,
  deleted_at        TIMESTAMPTZ,
  CONSTRAINT entities_key PRIMARY KEY ( id ),
  CONSTRAINT entities_owner_refs_entities FOREIGN KEY (owner_id) REFERENCES entities ( id )
);

CREATE INDEX entities_deleted_index ON entities ( id, deleted_at );

CREATE OR REPLACE FUNCTION trigger_protect_entity_key()
  RETURNS TRIGGER AS $$
    BEGIN
      IF NEW.id != OLD.id THEN
        RAISE EXCEPTION 'Entity key cannot be changed.';
      ELSE
        RETURN NEW;
      END IF;
    END;
  $$
  LANGUAGE plpgsql VOLATILE COST 100;

CREATE TRIGGER protect_entity_key
  BEFORE UPDATE OF id ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_protect_entity_key();

CREATE OR REPLACE FUNCTION trigger_protect_entity_created_at()
  RETURNS TRIGGER AS $$
    BEGIN
      IF NEW.created_at != OLD.created_at THEN
        RAISE EXCEPTION 'Entity created timestamp cannot be changed.';        
      ELSE
        RETURN NEW;
      END IF;
    END;
  $$
  LANGUAGE plpgsql VOLATILE COST 100;

CREATE TRIGGER protect_entity_created_at
  BEFORE UPDATE OF created_at ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_protect_entity_created_at();

CREATE OR REPLACE FUNCTION trigger_entities_last_updated()
  RETURNS TRIGGER AS $$
    BEGIN
      NEW.last_updated := NOW();
      RETURN NEW;
    END;
  $$ LANGUAGE 'plpgsql';

CREATE TRIGGER entities_last_updated
  -- BEFORE INSERT OR UPDATE ON entities
  BEFORE UPDATE ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_entities_last_updated();

CREATE OR REPLACE FUNCTION trigger_entities_defaults()
  RETURNS TRIGGER AS $$
  BEGIN
    NEW.created_at := NOW();
    NEW.last_updated := NOW();
    IF NEW.id IS NULL THEN NEW.id := uuid_generate_v4(); END IF;
    IF NEW.owner_id IS NULL THEN NEW.owner_id := NEW.id; END IF;
    RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';

CREATE TRIGGER entities_defaults
  BEFORE INSERT ON entities
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_entities_defaults();
