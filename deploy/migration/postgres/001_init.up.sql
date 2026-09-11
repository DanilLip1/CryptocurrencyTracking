CREATE TABLE coins_tracked (
    title TEXT primary key
);

CREATE TABLE coin_prices (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    creation_time TIMESTAMP NOT NULL
);

CREATE INDEX idx_coin_prices_title_creation_time
    ON coin_prices (title, price, creation_time DESC);