CREATE TYPE person.roles
AS ENUM('admin', 'normal');

CREATE TABLE person.users (
                              id SERIAL PRIMARY KEY,
                              first_name VARCHAR(255) NOT NULL,
                              last_name VARCHAR(255) NOT NULL,
                              email VARCHAR(255),
                              password VARCHAR(255),
                              role person.roles DEFAULT 'normal',
                              phone_number VARCHAR(20),
                              created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
                              updated_at TIMESTAMP
);