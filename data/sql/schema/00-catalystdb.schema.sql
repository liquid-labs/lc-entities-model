CREATE TABLE catalystdb (
  last_update TIMESTAMPTZ,
  schema_hash CHAR(64),
  version_spec TEXT,
  CONSTRAINT catalystdb_key PRIMARY KEY ( last_update,schema_hash )
);

-- Postgres
CREATE OR REPLACE FUNCTION trigger_catalyst_db_updated()
  RETURNS TRIGGER AS '
BEGIN
  NEW.last_update := NOW();
  RETURN NEW;
END' LANGUAGE 'plpgsql';

DROP EVENT TRIGGER IF EXISTS catalystdb_updated;

CREATE TRIGGER catalystdb_updated
  BEFORE INSERT ON catalystdb
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_catalyst_db_updated();
