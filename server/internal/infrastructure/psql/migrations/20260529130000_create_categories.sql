-- migrate:up
CREATE TABLE categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    icon       TEXT NOT NULL,
    color      TEXT NOT NULL,
    budget     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    keywords   TEXT[] NOT NULL DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0
);

-- Category rows are seeded as test data (see server/scripts/seed.go), not here.

-- migrate:down
DROP TABLE IF EXISTS categories;
