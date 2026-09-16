CREATE TABLE IF NOT EXISTS public.supplier (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    company_name VARCHAR(150),
    name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(30),
    address VARCHAR(255),
    is_active SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT supplier_email_unique UNIQUE (email)
);
