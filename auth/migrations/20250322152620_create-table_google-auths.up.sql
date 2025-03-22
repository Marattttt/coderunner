CREATE TABLE IF NOT EXISTS google_oauths (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry  TIMESTAMP NOT NULL
);
