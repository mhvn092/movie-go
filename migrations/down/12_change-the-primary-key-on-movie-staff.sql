ALTER TABLE movie.movie_staff DROP CONSTRAINT pk_movie_staff_composite;

ALTER TABLE movie.movie_staff ADD CONSTRAINT pk_movie_staff PRIMARY KEY (movie_id, staff_id);
