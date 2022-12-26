CREATE TABLE IF NOT EXISTS blacklist (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     reason VARCHAR(50) NOT NULL,
     user_id uuid NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ
);