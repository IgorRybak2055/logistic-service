CREATE TABLE truck
(
    id SERIAL PRIMARY KEY,
    company_id integer,
    name TEXT NOT NULL,
    trailer_type trailer_type NOT NULL,
    trailer_parameters TEXT NOT NULL,
    carrying float4 NOT NULL,
    year integer NOT NULL,
    current_location TEXT NOT NULL,

    FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);