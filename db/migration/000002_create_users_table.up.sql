CREATE TABLE IF NOT EXISTS users (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     full_name VARCHAR NOT NULL,
     email VARCHAR UNIQUE NOT NULL,
     password VARCHAR NOT NULL,
     role_id uuid NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ,
     last_login_at TIMESTAMPTZ
);