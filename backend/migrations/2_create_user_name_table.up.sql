CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS "username"
(
    id          uuid DEFAULT public.uuid_generate_v4() UNIQUE NOT NULL PRIMARY KEY,
    user_id     uuid NOT NULL,
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES public.user(id)

);