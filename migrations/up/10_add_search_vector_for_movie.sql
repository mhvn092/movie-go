ALTER TABLE movie.movie
ADD COLUMN search_vector tsvector;

-- delimiter //
CREATE OR REPLACE FUNCTION movie.movie_search_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := setweight(to_tsvector('simple', COALESCE(NEW.title, '')), 'A') ||
                      setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql//
-- delimiter ;

CREATE INDEX movie_search_idx ON movie.movie USING GIN(search_vector);

