CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    primary_email VARCHAR(256),
    primary_email_upcase VARCHAR(256),
    primary_phone VARCHAR(32),
    name VARCHAR(256),
    name_upcase VARCHAR(256),
    image_uri VARCHAR(1024),
    is_banned TINYINT(1) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
);

CREATE TABLE users_emails (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    email VARCHAR(256),
    email_upcase VARCHAR(256),
    email_upcase_digest VARCHAR(512) NOT NULL,
    is_active TINYINT(1) NOT NULL,
    is_verified TINYINT(1) NOT NULL,
    is_primary TINYINT(1) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_emails_user_id ON users_emails (user_id);

CREATE TABLE users_phones (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    phone VARCHAR(32),
    phone_digest VARCHAR(32) NOT NULL,
    is_active TINYINT(1) NOT NULL,
    is_verified TINYINT(1) NOT NULL,
    is_primary TINYINT(1) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_phones_user_id ON users_phones (user_id);

CREATE TABLE user_password_auth (
    user_id CHAR(36),
    password_digest VARCHAR(1024) NOT NULL,
    password_expires_at DATETIME,
    last_attempted_at DATETIME,
    attempt_count INT NOT NULL,
    otp_digest VARCHAR(1024),
    otp_expires_at DATETIME,
    otp_link_token VARCHAR(512),
    is_locked TINYINT(1) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_password_auth_user_id ON user_password_auth (user_id);

CREATE TABLE user_sessions (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    expires_at DATETIME NOT NULL,
    token TEXT NOT NULL,
    ip_address VARCHAR(40),
    user_agent VARCHAR(256),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_sessions_user_id ON user_sessions (user_id);

CREATE TABLE user_totp (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    secret_encrypted VARCHAR(512) NOT NULL,
    recovery_codes_hash TEXT NOT NULL,
    is_active TINYINT(1) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_totp_user_id ON user_totp (user_id);

CREATE TABLE user_login_providers (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    provider_id VARCHAR(64) NOT NULL,
    account_id VARCHAR(128) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    id_token TEXT,
    access_token_expires_at DATETIME,
    refresh_token_expires_at DATETIME,
    scopes TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_login_providers_user_id ON user_login_providers (user_id);

CREATE TABLE user_passkey (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    credential_id VARBINARY(1024) NOT NULL,
    data JSON NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_passkey_user_id ON user_passkey (user_id);

CREATE TABLE user_claims (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id CHAR(36) NOT NULL,
    type VARCHAR(256) NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_claims_user_id ON user_claims (user_id);

CREATE TABLE roles (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    name_upcase VARCHAR(64) NOT NULL,
    description VARCHAR(512) NOT NULL
);

CREATE TABLE users_roles (
    user_id CHAR(36),
    role_id CHAR(36),
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE role_claims (
    id INT PRIMARY KEY AUTO_INCREMENT,
    role_id CHAR(36) NOT NULL,
    type TEXT NOT NULL,
    claim_value TEXT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX ix_role_claims_role_id ON role_claims (role_id);

CREATE TABLE user_api_keys (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id CHAR(36) NOT NULL,
    name VARCHAR(64) NOT NULL,
    name_upcase VARCHAR(64) NOT NULL,
    key_hint VARCHAR(16) NOT NULL,
    key_digest VARCHAR(1024) NOT NULL,
    expires_at DATETIME,
    is_locked TINYINT(1) NOT NULL,
    is_revoked TINYINT(1) NOT NULL,
    comment VARCHAR(512),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_api_keys_user_id ON user_api_keys (user_id);

CREATE TABLE user_api_keys_roles (
    user_api_key_id INT,
    role_id CHAR(36),
    PRIMARY KEY (user_api_key_id, role_id),
    FOREIGN KEY (user_api_key_id) REFERENCES user_api_keys(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ix_users_primary_email_upcase ON users (primary_email_upcase);

CREATE UNIQUE INDEX ix_users_name_upcase ON users (name_upcase);

CREATE UNIQUE INDEX ix_roles_name_upcase ON roles (name_upcase);

CREATE UNIQUE INDEX ix_user_api_keys_user_id_name_upcase ON user_api_keys (user_id, name_upcase);
