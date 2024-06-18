ALTER TABLE movie.movie_staff DROP CONSTRAINT fk_movie_staff_and_staff;
ALTER TABLE movie.movie_staff DROP CONSTRAINT fk_movie_staff_and_movie;
ALTER TABLE movie.movie_staff DROP CONSTRAINT fk_movie_staff_and_staff_type;
DROP TABLE movie.movie_staff;