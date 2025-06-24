CREATE TRIGGER staff_search_update
BEFORE INSERT OR UPDATE OF first_name, last_name, bio
ON staff.staff
FOR EACH ROW
EXECUTE FUNCTION staff.staff_search_trigger();

UPDATE staff.staff
SET search_vector = setweight(to_tsvector('simple', COALESCE(first_name, '') || ' ' || COALESCE(last_name, '')), 'A') ||
                    setweight(to_tsvector('simple', COALESCE(bio, '')), 'B')
WHERE search_vector IS NULL;
