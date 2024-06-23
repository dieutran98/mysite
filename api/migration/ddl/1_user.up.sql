-- Create an enum type
CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');

CREATE TABLE IF NOT EXISTS "user_account" (
    "id" serial PRIMARY KEY,
    "user_name" varchar(200) NOT NULL,
    "password" varchar(200) NOT NULL,
    "is_active" boolean NOT NULL,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "user_info" (
    "id" serial PRIMARY KEY,
    user_account_id integer NOT NULL,
    "name" varchar(100),
    "phone" varchar(15),
    "email" varchar(200),
    "gender" gender_type,
    "membership_id" integer,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    CONSTRAINT user_info_user_account_fk FOREIGN KEY (user_account_id) REFERENCES user_account(id)
);
