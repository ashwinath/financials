CREATE TABLE IF NOT EXISTS users (
    id            text        NOT NULL PRIMARY KEY,
    password_hash text,
    created_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS session (
    session_id text NOT NULL PRIMARY KEY,
    user_id    text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_session_session_id_user_id ON session(session_id, user_id);
