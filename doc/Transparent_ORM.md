# Transparent ORM

We came close to a transparent orm, but each route had a problem.


### RULE based INSERT rewrite

The problem here is `NEW` cannot be used inside the `WITH` clause (pg11).

```sql
RULE based INSERT
CREATE RULE users_join_entity_ins AS ON INSERT TO users_join_entity
DO INSTEAD
WITH
  e AS (
    INSERT INTO entities
      (id, resource_name, name, description, owner_id, publicly_readable)
      VALUES (NEW.id, NEW.resource_name, NEW.name, NEW.description, NEW.owner_id, NEW.publicly_readable)
      RETURNING id
  ),
  s AS (INSERT INTO subjects (id) SELECT id FROM e)
INSERT INTO users
  (id, auth_id, legal_id, legal_id_type, active)
  SELECT id, NEW.auth_id, NEW.legal_id, NEW.legal_id_type, NEW.active FROM e
RETURNING users.*,
  (SELECT e.created_at FROM entities e WHERE users.id=e.id),
  (SELECT e.last_updated FROM entities e WHERE users.id=e.id);
```

RULE based INSERT using FUNCTION
CREATE OR REPLACE FUNCTION users_join_entity_ins_func(id UUID, resource_name VARCHAR(128), name VARCHAR(128), description TEXT, owner_id UUID, publicly_readable BOOLEAN,
  auth_id VARCHAR(128), legal_id VARCHAR(128), legal_id_type VARCHAR(64), active BOOLEAN)
RETURNS VOID AS $$
DECLARE
  uid UUID;
BEGIN
  RAISE NOTICE 'New ID: %', NEW.id;
  INSERT INTO entities
    (id, resource_name, name, description, owner_id, publicly_readable)
    VALUES (id, resource_name, name, description, owner_id, publicly_readable)
    RETURNING "id" INTO uid;
  INSERT INTO subjects (id) VALUES (uid);
  INSERT INTO users
    (id, auth_id, legal_id, legal_id_type, active)
    VALUES (uid, NEW.auth_id, NEW.legal_id, NEW.legal_id_type, NEW.active);
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';


CREATE RULE users_join_entity_ins AS ON INSERT TO users_join_entity
DO INSTEAD
  SELECT users_join_entity_ins_func(NEW.id, NEW.resource_name, NEW.name, NEW.description, NEW.owner_id, NEW.publicly_readable,
    NEW.auth_id, NEW.legal_id, CAST(NEW.legal_id_type AS VARCHAR(64)), NEW.active);




CREATE OR REPLACE FUNCTION users_join_entity_update()
RETURNS TRIGGER AS $$
  DECLARE updateU BOOLEAN := FALSE;
  DECLARE updateE BOOLEAN := FALSE;
BEGIN
  IF NEW.auth_id != OLD.auth_id
       OR NEW.legal_id != OLD.legal_id
       OR NEW.legal_id_type != OLD.legal_id_type
       OR NEW.active != OLD.active THEN
    UPDATE users
      SET auth_id=NEW.auth_id, legal_id=NEW.legal_ID, legal_id_type=NEW.legal_id_type, active=NEW.active
      WHERE id=NEW.id;
  END IF;
  IF NEW.name != OLD.name
       OR NEW.description != OLD.description
       OR NEW.owner_id != OLD.owner_id
       OR NEW.publicly_readable != OLD.publicly_readable THEN
    UPDATE entities
      SET name=NEW.NAME, description=NEW.description, owner_id=NEW.owner_id, publicly_readable=NEW.owner_id
      WHERE id=NEW.id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER users_join_entity_update_trigger
INSTEAD OF UPDATE ON users_join_entity
FOR EACH ROW
EXECUTE PROCEDURE users_join_entity_update();

CREATE OR REPLACE FUNCTION users_join_entity_insert()
RETURNS TRIGGER AS $$
DECLARE
  uid UUID;
BEGIN
  RAISE NOTICE 'New ID: %', NEW.id;
  INSERT INTO entities
    (id, resource_name, name, description, owner_id, publicly_readable)
    VALUES (NEW.id, NEW.resource_name, NEW.name, NEW.description, NEW.owner_id, NEW.publicly_readable)
    RETURNING "id" INTO uid;
  INSERT INTO subjects (id) VALUES (uid);
  INSERT INTO users
    (id, auth_id, legal_id, legal_id_type, active)
    VALUES (uid, NEW.auth_id, NEW.legal_id, NEW.legal_id_type, NEW.active);
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER users_join_entity_insert_trigger
INSTEAD OF INSERT ON users_join_entity
FOR EACH ROW
EXECUTE PROCEDURE users_join_entity_insert();
