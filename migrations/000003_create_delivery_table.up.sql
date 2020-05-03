create type trailer_type as enum
    (
        'тент',
        'коробка',
        'рефрижератор',
        'контейнеровоз',
        'цистерна',
        'зерновоз',
        'автовоз',
        'скотовоз',
        'площадка'
        );

CREATE TABLE IF NOT EXISTS delivery
(
    id              SERIAL PRIMARY KEY,
    shipment_date   TIMESTAMP               NOT NULL,
    shipment_place  TEXT                    NOT NULL,
    unloading_place TEXT                    NOT NULL,
    cargo           TEXT                    NOT NULL,
    weight_cargo    float4                  NOT NULL,
    volume_cargo    float4                  NOT NULL,
    trailer_type    trailer_type,
    price           float4                  NOT NULL,
    created_at      TIMESTAMP DEFAULT now() NOT NULL,
    updated_at      TIMESTAMP DEFAULT now() NOT NULL
);
