CREATE TABLE doctors (
    id UUID PRIMARY KEY,
    fullname TEXT NOT NULL,
    speciality TEXT NOT NULL,
    address TEXT NOT NULL,
    licensenumber TEXT NOT NULL,
    department TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL  DEFAULT now()
);

CREATE UNIQUE INDEX uniq_doctor_ln on doctors(licensenumber);