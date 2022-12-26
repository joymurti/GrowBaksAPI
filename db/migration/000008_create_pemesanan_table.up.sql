CREATE TYPE pemesanan_status AS ENUM ('reserve', 'pending', 'done', 'cancel');

CREATE TABLE IF NOT EXISTS pemesanan (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     name VARCHAR NOT NULL,
     status pemesanan_status NOT NULL,
     user_id uuid NOT NULL,
     product_id uuid NOT NULL, 
     qty INTEGER NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ
);