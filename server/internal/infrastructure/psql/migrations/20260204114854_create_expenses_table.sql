-- migrate:up
CREATE TABLE expenses (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS expenses;
