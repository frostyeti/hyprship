CREATE TABLE user_password_history (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    password_digest NVARCHAR(1024) NOT NULL,
    created_at DATETIME2 NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX ix_user_password_history_user_id ON user_password_history (user_id);
