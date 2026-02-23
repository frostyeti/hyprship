CREATE TABLE users (
    id UUID PRIMARY KEY,
    primary_email TEXT,
    primary_email_upcase TEXT,
    primary_phone TEXT,
    name TEXT,
    name_upcase TEXT,
    image_uri TEXT,
    is_banned BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE TABLE users_emails (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    email TEXT,
    email_upcase TEXT,
    email_upcase_digest TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    is_verified BOOLEAN NOT NULL,
    is_primary BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_emails_user_id ON users_emails (user_id);

CREATE TABLE users_phones (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    phone TEXT,
    phone_digest TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    is_verified BOOLEAN NOT NULL,
    is_primary BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_phones_user_id ON users_phones (user_id);

CREATE TABLE user_password_auth (
    user_id UUID,
    password_digest TEXT NOT NULL,
    password_expires_at TIMESTAMP,
    last_attempted_at TIMESTAMP,
    attempt_count INTEGER NOT NULL,
    otp_digest TEXT,
    otp_expires_at TIMESTAMP,
    otp_link_token TEXT,
    is_locked BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_password_auth_user_id ON user_password_auth (user_id);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    token TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_sessions_user_id ON user_sessions (user_id);

CREATE TABLE user_totp (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    secret_encrypted TEXT NOT NULL,
    recovery_codes_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_totp_user_id ON user_totp (user_id);

CREATE TABLE user_login_providers (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    provider_id TEXT NOT NULL,
    account_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT,
    access_token_expires_at TIMESTAMP,
    refresh_token_expires_at TIMESTAMP,
    scopes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_login_providers_user_id ON user_login_providers (user_id);

CREATE TABLE user_passkey (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    credential_id BYTEA NOT NULL,
    data JSONB NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_passkey_user_id ON user_passkey (user_id);

CREATE TABLE user_claims (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_claims_user_id ON user_claims (user_id);

CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE users_roles (
    user_id UUID,
    role_id UUID,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE role_claims (
    id SERIAL PRIMARY KEY,
    role_id UUID NOT NULL,
    type TEXT NOT NULL,
    claim_value TEXT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX ix_role_claims_role_id ON role_claims (role_id);

CREATE TABLE user_api_keys (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    name_upcase TEXT NOT NULL,
    key_hint TEXT NOT NULL,
    key_digest TEXT NOT NULL,
    expires_at TIMESTAMP,
    is_locked BOOLEAN NOT NULL,
    is_revoked BOOLEAN NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_api_keys_user_id ON user_api_keys (user_id);

CREATE TABLE user_api_keys_roles (
    user_api_key_id INTEGER,
    role_id UUID,
    PRIMARY KEY (user_api_key_id, role_id),
    FOREIGN KEY (user_api_key_id) REFERENCES user_api_keys(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_users_primary_email_upcase ON users (primary_email_upcase);

CREATE UNIQUE INDEX ix_users_name_upcase ON users (name_upcase);

CREATE UNIQUE INDEX ix_roles_name_upcase ON roles (name_upcase);

CREATE UNIQUE INDEX ix_user_api_keys_user_id_name_upcase ON user_api_keys (user_id, name_upcase);

CREATE UNIQUE INDEX uq_users_emails_primary ON users_emails (user_id) WHERE is_primary = true;

CREATE UNIQUE INDEX uq_users_phones_primary ON users_phones (user_id) WHERE is_primary = true;
