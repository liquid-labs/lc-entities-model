CREATE OR REPLACE FUNCTION numeric_only (val VARCHAR(255))
  RETURNS VARCHAR(255) AS $$
    DECLARE
      idx INT := 1;
      len INT := 1;
      res VARCHAR(255) := '';
      c VARCHAR(1) := '';
    BEGIN
     IF val IS NULL THEN RETURN NULL; END IF;
     IF LENGTH(val) = 0 THEN RETURN ''; END IF;

     len := LENGTH(val);
     WHILE idx <= len LOOP
       c := SUBSTRING(val FROM idx FOR 1);
       IF IS_NUMERIC(c) = 1 THEN
         res := CONCAT(res, c);
       END IF;
       idx := idx + 1;
     END LOOP;
   RETURN res;
 END $$ LANGUAGE 'plpgsql';
