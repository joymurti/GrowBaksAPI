CREATE TABLE IF NOT EXISTS roles (
      id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
      name VARCHAR UNIQUE NOT NULL,
      created_at TIMESTAMPTZ DEFAULT now(),
      updated_at TIMESTAMPTZ
);