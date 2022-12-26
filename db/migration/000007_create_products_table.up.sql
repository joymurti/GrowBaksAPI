CREATE TYPE kategori AS ENUM ('makanan', 'minuman');

CREATE TABLE IF NOT EXISTS products (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     name VARCHAR NOT NULL,
     stok INTEGER NOT NULL DEFAULT 0,
     product_kategori kategori NOT NULL, 
     product_img TEXT,
     lapak_id uuid NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ
);