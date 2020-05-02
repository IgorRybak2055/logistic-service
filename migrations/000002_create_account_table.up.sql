CREATE TABLE IF NOT EXISTS account
(
    id         SERIAL PRIMARY KEY,
    name       TEXT,
    email      TEXT      NOT NULL UNIQUE,
    password   TEXT      NOT NULL,
    phone      TEXT,
    company_id integer,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,

  FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE RESTRICT ON DELETE CASCADE
);
