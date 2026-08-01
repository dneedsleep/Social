CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    first_name Varchar(255) NOT NULL,
    last_name Varchar(255) NOT NULL,
    email citext
);