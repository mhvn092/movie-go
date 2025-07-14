DROP TRIGGER IF EXISTS movie_search_update ON movie.movie;

CREATE TRIGGER movie_search_update
BEFORE INSERT OR UPDATE OF title, description
ON movie.movie
FOR EACH ROW
EXECUTE FUNCTION movie.movie_search_trigger();

UPDATE movie.movie
SET search_vector = setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
                    setweight(to_tsvector('simple', COALESCE(description, '')), 'B')
WHERE search_vector IS NULL;
