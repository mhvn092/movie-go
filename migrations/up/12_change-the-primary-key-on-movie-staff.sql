ALTER TABLE movie.movie_staff DROP CONSTRAINT pk_movie_staff;

ALTER TABLE movie.movie_staff ADD CONSTRAINT pk_movie_staff_composite PRIMARY KEY (movie_id, staff_id, staff_type_id);
