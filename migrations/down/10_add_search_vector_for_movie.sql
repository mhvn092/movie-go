DROP INDEX IF EXISTS movie.movie_search_idx;

ALTER TABLE movie.movie DROP COLUMN IF EXISTS search_vector;

DROP FUNCTION IF EXISTS movie.movie_search_trigger;

