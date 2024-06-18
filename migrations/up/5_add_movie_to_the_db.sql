CREATE TABLE movie.movie (
                             id SERIAL PRIMARY KEY,
                             title VARCHAR(255) NOT NULL,
                             description VARCHAR(255) NOT NULL,
                             production_year integer NOT NULL,
                             director_id integer NOT NULL,
                             genre_id integer NOT NULL,
                             created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
                             updated_at TIMESTAMP,
                             constraint fk_movie_and_staff FOREIGN KEY ("director_id") REFERENCES "staff"."staff" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
                             constraint fk_movie_and_genre FOREIGN KEY ("genre_id") REFERENCES "movie"."genre" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);