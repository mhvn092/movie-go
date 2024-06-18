CREATE TABLE movie.genre (
                                  id SERIAL PRIMARY KEY,
                                  title VARCHAR(255) NOT NULL,
                                  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
                                  updated_at TIMESTAMP
);