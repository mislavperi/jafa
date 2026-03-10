-- migrate:up
ALTER TABLE users
    ALTER COLUMN username TYPE VARCHAR(50),
    ALTER COLUMN first_name TYPE VARCHAR(100),
    ALTER COLUMN last_name TYPE VARCHAR(100),
    ALTER COLUMN email TYPE VARCHAR(255);

-- migrate:down
ALTER TABLE users
    ALTER COLUMN username TYPE TEXT,
    ALTER COLUMN first_name TYPE TEXT,
    ALTER COLUMN last_name TYPE TEXT,
    ALTER COLUMN email TYPE TEXT;
