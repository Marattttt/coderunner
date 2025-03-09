CREATE TABLE IF NOT EXISTS google_oauth (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry  TIMESTAMP NOT NULL
);
