DROP FUNCTION IF EXISTS test_catalystdb_updated_trigger();
CREATE FUNCTION test_catalystdb_updated_trigger()
  RETURNS BOOLEAN AS $$
    DECLARE
      lastUpdate TIMESTAMPTZ;
    BEGIN
      INSERT INTO catalystdb (schema_hash, version_spec) VALUES ('abc', 'def');
      SELECT catalystdb.last_update INTO STRICT lastUpdate FROM catalystdb WHERE catalystdb.schema_hash = 'abc';
      -- check that we have a timestamp in place an it's very recent
      IF lastUpdate IS NULL AND EXTRACT(EPOCH FROM (NOW() - lastUpdate)) <= 2 THEN
        RAISE EXCEPTION 'last_update not set when inserting into catalystdb'
          USING ERRCODE = 'check_violation';
      END IF;
      RETURN TRUE;
    END $$ LANGUAGE 'plpgsql';

SELECT test_catalystdb_updated_trigger() AS passed;
