CREATE TABLE IF NOT EXISTS company
(
    id           SERIAL PRIMARY KEY,
    company_type TEXT NOT NULL,
    name         TEXT                    NOT NULL,
    phone        TEXT,
    email        TEXT,
    bank_detail  TEXT                    NOT NULL,
    created_at   TIMESTAMP DEFAULT now() NOT NULL,
    updated_at   TIMESTAMP DEFAULT now() NOT NULL
);
