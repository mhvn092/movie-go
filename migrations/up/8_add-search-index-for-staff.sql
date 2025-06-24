ALTER TABLE staff.staff ADD COLUMN search_vector tsvector;

-- delimiter //

CREATE OR REPLACE FUNCTION staff.staff_search_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := setweight(to_tsvector('simple', COALESCE(NEW.first_name, '') || ' ' || COALESCE(NEW.last_name, '')), 'A') ||
                      setweight(to_tsvector('simple', COALESCE(NEW.bio, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql//

-- delimiter ;

CREATE INDEX IF NOT EXISTS staff_search_idx ON staff.staff USING GIN(search_vector);

