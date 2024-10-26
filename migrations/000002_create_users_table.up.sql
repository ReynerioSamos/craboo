-- Filename: migrations/000002_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    email text NOT NULL,
    fullname text NOT NULL,
    version integer NOT NULL DEFAULT 1
);