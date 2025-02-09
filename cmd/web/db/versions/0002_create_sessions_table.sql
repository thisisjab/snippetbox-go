CREATE TABLE IF NOT EXISTS sessions (
                          token TEXT PRIMARY KEY,
                          data BLOB NOT NULL,
                          expiry REAL NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry);