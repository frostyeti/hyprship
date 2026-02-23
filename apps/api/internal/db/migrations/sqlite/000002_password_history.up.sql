CREATE TABLE user_password_history (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    password_digest TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX ix_user_password_history_user_id ON user_password_history (user_id);
