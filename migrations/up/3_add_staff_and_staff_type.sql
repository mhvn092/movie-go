CREATE TABLE staff.staff_type (
                              id SERIAL PRIMARY KEY,
                              title VARCHAR(255) NOT NULL,
                              created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
                              updated_at TIMESTAMP
);


CREATE TABLE staff.staff (
                                  id SERIAL PRIMARY KEY,
                                  first_name VARCHAR(255) NOT NULL,
                                  last_name VARCHAR(255) NOT NULL,
                                  bio VARCHAR(255) NOT NULL,
                                  staff_type_id integer NOT NULL,
                                  birth_date  TIMESTAMP NOT NULL,
                                  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
                                  updated_at TIMESTAMP,
                                  constraint fk_staff_type_and_staff FOREIGN KEY ("staff_type_id") REFERENCES "staff"."staff_type" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);
