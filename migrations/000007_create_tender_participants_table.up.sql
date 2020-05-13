CREATE TABLE tender_participants
(
    id           SERIAL PRIMARY KEY,
    company_id   integer,
    delivery_id  integer,
    price        float4,

    FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (delivery_id) REFERENCES delivery (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);