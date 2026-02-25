IF OBJECT_ID('user_api_keys_roles', 'U') IS NOT NULL DROP TABLE user_api_keys_roles;
IF OBJECT_ID('user_api_keys', 'U') IS NOT NULL DROP TABLE user_api_keys;
IF OBJECT_ID('user_passkey', 'U') IS NOT NULL DROP TABLE user_passkey;

CREATE TABLE user_api_keys (
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL FOREIGN KEY REFERENCES users(id) ON DELETE CASCADE,
    name NVARCHAR(64) NOT NULL,
    name_upcase NVARCHAR(64) NOT NULL,
    key_hint NVARCHAR(16) NOT NULL,
    key_digest NVARCHAR(1024) NOT NULL,
    expires_at DATETIME2,
    is_locked BIT NOT NULL DEFAULT 0,
    is_revoked BIT NOT NULL DEFAULT 0,
    comment NVARCHAR(512),
    created_at DATETIME2 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME2
);

CREATE TABLE user_passkeys (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    user_id UNIQUEIDENTIFIER NOT NULL FOREIGN KEY REFERENCES users(id) ON DELETE CASCADE,
    credential_id VARBINARY(1024) NOT NULL,
    data NVARCHAR(MAX) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT CURRENT_TIMESTAMP
);
