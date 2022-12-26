CREATE TABLE IF NOT EXISTS lokasi (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     provinsi VARCHAR NOT NULL,
     kota VARCHAR NOT NULL,
     daerah VARCHAR NOT NULL
);