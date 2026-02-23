CREATE TABLE user_password_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    password_digest VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX ix_user_password_history_user_id ON user_password_history (user_id);
