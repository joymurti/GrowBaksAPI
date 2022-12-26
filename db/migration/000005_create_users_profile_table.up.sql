CREATE TYPE gender_type AS ENUM ('laki-laki', 'perempuan', 'not-set');

CREATE TABLE IF NOT EXISTS users_profile (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     telepon VARCHAR NOT NULL,
     gender gender_type DEFAULT 'not-set',
     tanggal_lahir DATE, 
     location_id uuid NOT NULL,
     user_id uuid NOT NULL,
     created_at TIMESTAMPTZ DEFAULT now(),
     updated_at TIMESTAMPTZ
);