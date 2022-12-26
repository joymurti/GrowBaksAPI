CREATE TYPE lapak_status AS ENUM ('open', 'closed');

CREATE TABLE IF NOT EXISTS lapak (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     name VARCHAR NOT NULL,
     status lapak_status DEFAULT 'closed', 
     user_id uuid NOT NULL,
     location_id uuid NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ
);