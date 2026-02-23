CREATE TABLE users (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    primary_email NVARCHAR(256),
    primary_email_upcase NVARCHAR(256),
    primary_phone NVARCHAR(32),
    name NVARCHAR(256),
    name_upcase NVARCHAR(256),
    image_uri NVARCHAR(1024),
    is_banned BIT NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2
);

CREATE TABLE users_emails (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    email NVARCHAR(256),
    email_upcase NVARCHAR(256),
    email_upcase_digest NVARCHAR(512) NOT NULL,
    is_active BIT NOT NULL,
    is_verified BIT NOT NULL,
    is_primary BIT NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_emails_user_id ON users_emails (user_id);

CREATE TABLE users_phones (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    phone NVARCHAR(32),
    phone_digest NVARCHAR(32) NOT NULL,
    is_active BIT NOT NULL,
    is_verified BIT NOT NULL,
    is_primary BIT NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_users_phones_user_id ON users_phones (user_id);

CREATE TABLE user_password_auth (
    user_id UNIQUEIDENTIFIER,
    password_digest NVARCHAR(1024) NOT NULL,
    password_expires_at DATETIME2,
    last_attempted_at DATETIME2,
    attempt_count INT NOT NULL,
    otp_digest NVARCHAR(1024),
    otp_expires_at DATETIME2,
    otp_link_token NVARCHAR(512),
    is_locked BIT NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_password_auth_user_id ON user_password_auth (user_id);

CREATE TABLE user_sessions (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    expires_at DATETIME2 NOT NULL,
    token NVARCHAR(MAX) NOT NULL,
    ip_address NVARCHAR(40),
    user_agent NVARCHAR(256),
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_sessions_user_id ON user_sessions (user_id);

CREATE TABLE user_totp (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    secret_encrypted NVARCHAR(512) NOT NULL,
    recovery_codes_hash NVARCHAR(MAX) NOT NULL,
    is_active BIT NOT NULL,
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_totp_user_id ON user_totp (user_id);

CREATE TABLE user_login_providers (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    provider_id NVARCHAR(64) NOT NULL,
    account_id NVARCHAR(128) NOT NULL,
    access_token NVARCHAR(MAX) NOT NULL,
    refresh_token NVARCHAR(MAX),
    id_token NVARCHAR(MAX),
    access_token_expires_at DATETIME2,
    refresh_token_expires_at DATETIME2,
    scopes NVARCHAR(MAX),
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_login_providers_user_id ON user_login_providers (user_id);

CREATE TABLE user_passkey (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL,
    credential_id VARBINARY(1024) NOT NULL,
    data NVARCHAR(MAX) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_passkey_user_id ON user_passkey (user_id);

CREATE TABLE user_claims (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id UNIQUEIDENTIFIER NOT NULL,
    type NVARCHAR(256) NOT NULL,
    value NVARCHAR(MAX) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_claims_user_id ON user_claims (user_id);

CREATE TABLE roles (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    name NVARCHAR(64) NOT NULL,
    name_upcase NVARCHAR(64) NOT NULL,
    description NVARCHAR(512) NOT NULL
);

CREATE TABLE users_roles (
    user_id UNIQUEIDENTIFIER,
    role_id UNIQUEIDENTIFIER,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE role_claims (
    id INT PRIMARY KEY IDENTITY(1,1),
    role_id UNIQUEIDENTIFIER NOT NULL,
    type NVARCHAR(MAX) NOT NULL,
    claim_value NVARCHAR(MAX) NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX ix_role_claims_role_id ON role_claims (role_id);

CREATE TABLE user_api_keys (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(64) NOT NULL,
    name_upcase NVARCHAR(64) NOT NULL,
    key_hint NVARCHAR(16) NOT NULL,
    key_digest NVARCHAR(1024) NOT NULL,
    expires_at DATETIME2,
    is_locked BIT NOT NULL,
    is_revoked BIT NOT NULL,
    comment NVARCHAR(512),
    created_at DATETIME2 NOT NULL,
    updated_at DATETIME2,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ix_user_api_keys_user_id ON user_api_keys (user_id);

CREATE TABLE user_api_keys_roles (
    user_api_key_id INT,
    role_id UNIQUEIDENTIFIER,
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
