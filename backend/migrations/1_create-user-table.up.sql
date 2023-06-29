CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS "user"
(
    id              uuid DEFAULT public.uuid_generate_v4() UNIQUE NOT NULL PRIMARY KEY,
    password        VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE
);
