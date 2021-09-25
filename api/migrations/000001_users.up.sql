CREATE TABLE IF NOT EXISTS users (
    id            text        NOT NULL PRIMARY KEY,
    username      text        NOT NULL UNIQUE,
    password_hash text,
    created_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id         text        NOT NULL PRIMARY KEY,
    user_id    text        NOT NULL,
    expiry     timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_session_id_user_id ON sessions(id, user_id);
