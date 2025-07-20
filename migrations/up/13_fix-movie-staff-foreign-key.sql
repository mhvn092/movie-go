ALTER TABLE "movie"."movie_staff" DROP CONSTRAINT "fk_movie_staff_and_staff_type";

ALTER TABLE "movie"."movie_staff" ADD CONSTRAINT "fk_movie_staff_and_staff_type" FOREIGN KEY ("staff_type_id") REFERENCES "staff"."staff_type" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
