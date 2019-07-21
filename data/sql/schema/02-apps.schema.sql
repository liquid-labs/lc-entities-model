CREATE TABLE apps (
  id BIGSERIAL,
  name VARCHAR(128) NOT NULL,
  type VARCHAR(64) NOT NULL,
  CONSTRAINT apps_key PRIMARY KEY ( id ),
  CONSTRAINT apps_ref_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);
