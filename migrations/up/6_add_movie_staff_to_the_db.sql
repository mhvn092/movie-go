CREATE TABLE movie.movie_staff (
                             movie_id integer NOT NULL,
                             staff_id integer NOT NULL,
                             staff_type_id integer NOT NULL,
                             CONSTRAINT pk_movie_staff PRIMARY KEY ("movie_id", "staff_id"),
                             constraint fk_movie_staff_and_staff FOREIGN KEY ("staff_id") REFERENCES "staff"."staff" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
                             constraint fk_movie_staff_and_movie FOREIGN KEY ("movie_id") REFERENCES "movie"."movie" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
                             constraint fk_movie_staff_and_staff_type FOREIGN KEY ("staff_id") REFERENCES "staff"."staff_type" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);