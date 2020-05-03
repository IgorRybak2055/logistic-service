CREATE TABLE truck
(
    id SERIAL PRIMARY KEY,
    company_id integer,
    name TEXT NOT NULL,
    trailer_length float4,
    trailer_width float4,
    trailer_height float4,
    carrying float4 NOT NULL,
    year integer NOT NULL,
    current_location TEXT NOT NULL,

    FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);