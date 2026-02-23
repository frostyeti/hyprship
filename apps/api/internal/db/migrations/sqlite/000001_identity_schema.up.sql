CREATE TABLE users (
    id TEXT PRIMARY KEY,
    primary_email TEXT,
    primary_email_upcase TEXT,
    primary_phone TEXT,
    name TEXT,
    name_upcase TEXT,
    image_uri TEXT,
    is_banned INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER
);

CREATE TABLE users_emails (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    email TEXT,
    email_upcase TEXT,
    email_upcase_digest TEXT NOT NULL,
    is_active INTEGER NOT NULL,
    is_verified INTEGER NOT NULL,
    is_primary INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_emails_user_id ON users_emails (user_id);

CREATE TABLE users_phones (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    phone TEXT,
    phone_digest TEXT NOT NULL,
    is_active INTEGER NOT NULL,
    is_verified INTEGER NOT NULL,
    is_primary INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_phones_user_id ON users_phones (user_id);

CREATE TABLE user_password_auth (
    user_id TEXT,
    password_digest TEXT NOT NULL,
    password_expires_at INTEGER,
    last_attempted_at INTEGER,
    attempt_count INTEGER NOT NULL,
    otp_digest TEXT,
    otp_expires_at INTEGER,
    otp_link_token TEXT,
    is_locked INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_password_auth_user_id ON user_password_auth (user_id);

CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    token TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_sessions_user_id ON user_sessions (user_id);

CREATE TABLE user_totp (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    secret_encrypted TEXT NOT NULL,
    recovery_codes_hash TEXT NOT NULL,
    is_active INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_totp_user_id ON user_totp (user_id);

CREATE TABLE user_login_providers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT,
    access_token_expires_at INTEGER,
    refresh_token_expires_at INTEGER,
    scopes TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_login_providers_user_id ON user_login_providers (user_id);

CREATE TABLE user_passkey (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    credential_id BLOB NOT NULL,
    data TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_passkey_user_id ON user_passkey (user_id);

CREATE TABLE user_claims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_claims_user_id ON user_claims (user_id);

CREATE TABLE roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE users_roles (
    user_id TEXT,
    role_id TEXT,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE role_claims (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id TEXT NOT NULL,
    type TEXT NOT NULL,
    claim_value TEXT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX ix_role_claims_role_id ON role_claims (role_id);

CREATE TABLE user_api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    key_hint TEXT NOT NULL,
    key_digest TEXT NOT NULL,
    expires_at INTEGER,
    is_locked INTEGER NOT NULL,
    is_revoked INTEGER NOT NULL,
    comment TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_api_keys_user_id ON user_api_keys (user_id);

CREATE TABLE user_api_keys_roles (
    user_api_key_id INTEGER,
    role_id TEXT,
    PRIMARY KEY (user_api_key_id, role_id),
    FOREIGN KEY (user_api_key_id) REFERENCES user_api_keys(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_users_primary_email_upcase ON users (primary_email_upcase);

CREATE UNIQUE INDEX ix_users_name_upcase ON users (name_upcase);

CREATE UNIQUE INDEX ix_roles_name_upcase ON roles (name_upcase);

CREATE UNIQUE INDEX ix_user_api_keys_user_id_name_upcase ON user_api_keys (user_id, name_upcase);

CREATE UNIQUE INDEX uq_users_emails_primary ON users_emails (user_id) WHERE is_primary = 1;

CREATE UNIQUE INDEX uq_users_phones_primary ON users_phones (user_id) WHERE is_primary = 1;
