CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    city TEXT NOT NULL,
    created_at BIGINT,
    updated_at BIGINT,
    frequency_type INTEGER NOT NULL,
    token TEXT NOT NULL,
    status integer NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_updated_at ON subscriptions (frequency_type, updated_at);

CREATE INDEX IF NOT EXISTS idx_token ON subscriptions (token);

CREATE INDEX IF NOT EXISTS idx_email ON subscriptions (email);
