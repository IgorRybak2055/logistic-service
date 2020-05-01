CREATE TABLE IF NOT EXISTS company_type
(
    id   SERIAL PRIMARY KEY,
    type TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS company
(
    id           SERIAL PRIMARY KEY,
    company_type INTEGER,
    name         TEXT                    NOT NULL,
    phone        TEXT,
    email        TEXT,
    bank_detail  TEXT                    NOT NULL,
    created_at   TIMESTAMP DEFAULT now() NOT NULL,
    updated_at   TIMESTAMP DEFAULT now() NOT NULL,

    FOREIGN KEY (company_type) REFERENCES company_type (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);
