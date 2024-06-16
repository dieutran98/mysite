-- Create extension for uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create an enum type
CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');

CREATE TABLE IF NOT EXISTS "user_account" (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "user_name" varchar(200) NOT NULL,
    "password" varchar(200) NOT NULL,
    "is_active" boolean NOT NULL,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "user_info" (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_account_id UUID NOT NULL,
    "name" varchar(100) NOT NULL,
    "phone" varchar(15),
    "address" varchar(200),
    "gender" gender_type,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    CONSTRAINT user_info_user_account_fk FOREIGN KEY (user_account_id) REFERENCES user_account(id)
);
