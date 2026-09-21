BEGIN;

DROP INDEX IF EXISTS idx_coin_prices_title_creation_time;
DROP TABLE IF EXISTS coin_prices;
DROP TABLE IF EXISTS coins_tracked;

COMMIT;