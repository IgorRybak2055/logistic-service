CREATE TABLE IF NOT EXISTS subscribers
(
    company_id integer,
    subscriber_id integer,

    FOREIGN KEY (company_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (subscriber_id) REFERENCES company (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE UNIQUE INDEX track_topic_id_date_uindex
    ON subscribers (company_id, subscriber_id);