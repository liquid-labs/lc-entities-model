CREATE OR REPLACE FUNCTION is_numeric (val varchar(255))
RETURNS BOOLEAN AS $$
SELECT val ~ '^-?[0-9]+$'
$$ LANGUAGE SQL;
