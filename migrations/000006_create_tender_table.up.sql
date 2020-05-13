CREATE TABLE tender
(
    id           SERIAL PRIMARY KEY,
    company_id   integer,
    start        timestamp default now(),
    close        timestamp,
    delivery_id  integer,
    price        float4,
    participants json,

    FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (delivery_id) REFERENCES delivery (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);