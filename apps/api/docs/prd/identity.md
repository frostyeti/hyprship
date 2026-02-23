# Identity

## Business Case

Enable identity authentication and authorization for hyprship. Default to JWT
authentication and leave room for integration with other auth providers
at later date.

The app will need auth to keep out unauthorized users as it will contain
pontentially sensitive data.

## Overview

Use the information in this document to create an epic for identity with various
proposed issues under the apps/api/doc/issues/drafts folder.

Use Skills as needed.

Identity should use jwt for authentication and for authorization where the claims
have the scopes/permissions in the jwt. We will use short-lived Access Tokens (JWT) for stateless, performant API requests, and long-lived Refresh Tokens stored in the `user_sessions` table to track sessions and enable immediate revocation. Make security recommendations based upon common identity vulnerabilities related to jwt and session auth.  

Users may have claims for the case of groups or other claims and roles will have claims
that primarily hold permissions. Users that are in given roles will inherit the roles permission
claims.  

Users may create api keys that may or may not expire that also have roles and through those
roles an api key will get permission claims assigned to the key.

Users may have multiple api keys.

`*_upcase` columns will be used for case insentive search when search service is not available.

dependinging on configuration users may need to verify their e-mail and their phone number.

Enable MFA and TOTP and passkeys. Recommend modules based upon popularity and security.

By default use smtp for e-mail and support twilio if enabled.  

The application supports multiple database providers.

propose API route design, indexes, constraints, additional fields, etc as needed. include
concise rational for proposal.

### Propose Configuration

Add a config section (viper) for identity:

- `identity.password.min_length` (default 12)
- `identity.password.max_length` (default 128)
- `identity.password.require_upper` (default true)
- `identity.password.require_lower` (default true)
- `identity.password.require_number` (default true)
- `identity.password.require_symbol` (default true)
- `identity.password.history_count` (default 5)
- `identity.password.expire_days` (default 0 = off)
- `identity.phone.required` (default false)
- `identity.verification.email_required` (default true)
- `identity.verification.phone_required` (default false)
- `identity.lockout.max_attempts` (default 10)
- `identity.lockout.duration_minutes` (default 15)
- `identity.sessions.max_concurrent` (default 10)
- `identity.sessions.ttl_minutes` (default 1440)
- `identity.sessions.ip_required` (default false; if false, warn on missing IP)
- `identity.reset.token_ttl_minutes` (default 30)

## Proposed Permission Claims

The `apps/api/claims/permissions` package defines standard permission constants:

| Constant           | Value          | Description                        |
|--------------------|----------------|------------------------------------|
| `ADMIN`            | `admin`        | Full administrator access          |
| `READ`             | `read`         | Read-only access except identity   |
| `WRITE`            | `write`        | Write access except identity       |
| `EXEC`             | `exec`         | Execute/access access              |
| `USERS_READ`       | `users:read`   | read user info                     |
| `USERS_WRITE       | `users:write`  | write user info                    |
| `ROLES_READ`       | `roles:read`   | read roles                         |
| `ROLES_WRITE`      | `roles:write`  |  write/delete roles                |

these will be updated as needed.

## Proposed Database Schema

### users - Table

This table should only hold what will be queried often for the user.  

| Field                     | Type       | Attributes |
| ------------------------- | ---------- | ---------- |
| id                        | uuid       | pk         |
| primary_email             | text(256)  | —          |
| primary_email_upcase      | text(256)  | uniq, ix   |
| primary_phone             | text(32)   | —          |
| name                      | text(256)  | —          |
| name_upcase               | text(256)  | uniq, ix   |
| image_uri                 | text(1024) | nil        |
| is_banned                 | boolean    | —          |
| created_at                | datetime   | —          |
| updated_at                | datetime   | nil        |

### users_emails - table

stores all historical e-mails for the user as a hash. User may remove
e-mail from the email and email_upcase fields, but the hash must
always remain. only one can be primary per user.

| Field Name         | Data Type | Attributes |
| ------------------ | --------- | ---------- |
| user_id            | uuid      | fk         |
| email              | text(256) | nil        |
| email_upcase       | text(256) | nil        |
| email_upcase_digest| text(512) | —          |
| is_active          | boolean   | —          |
| is_verified        | boolean   | —          |
| is_primary         | boolean   | —          |
| created_at         | datetime  | —          |
| updated_at         | datetime  | nil        |

### users_phones - Table

stores all historical phone for the user as a hash. User may
remove phone from phone and phone_upcase fields, but the hash
must always remain. only one can be primary per user.  

| Field Name   | Data Type | Attributes |
| ------------ | --------- | ---------- |
| user_id      | uuid      | fk         |
| phone        | text(32)  | nil        |
| phone_digest | text(32)  | —          |
| is_active    | boolean   | —          |
| is_verified  | boolean   | —          |
| is_primary   | boolean   | —          |
| created_at   | datetime  | —          |
| updated_at   | datetime  | nil        |

### user_password_auth - table

otp is one time password for reseting passwords or sending an invite to new users.  

| Field Name              | Data Type  | Attributes |
| ----------------------- | ---------- | ---------- |
| user_id                 | uuid       | fk         |
| password_digest         | text(1024) | —          |
| password_expires_at     | datetime   | nil        |
| last_attempted_at       | datetime   | —          |
| attempt_count           | int        | —          |
| otp_digest              | text(1024) | nil        |
| otp_expires_at          | datetime   | nil        |
| otp_link_token          | text(512)  | nil        |
| is_locked               | boolean    | —          |
| created_at              | datetime   | —          |
| updated_at              | datetime   | nil        |

### user_sesions - table

| Field Name  | Data Type | Attributes |
| ----------- | --------- | ---------- |
| id          | uuid      | pk auto++  |
| user_id     | uuid      | —          |
| expires_at  | datetime  | —          |
| token       | text      | —          |
| ip_address  | text(40)  | nil        |
| user_agent  | text(256) | nil        |
| created_at  | datetime  | —          |
| updated_at  | datetime  | nil        |

### user_totp - table

Stores TOTP configurations and recovery codes for Multi-Factor Authentication.

| Field Name           | Data Type | Attributes |
| -------------------- | --------- | ---------- |
| id                   | uuid      | pk         |
| user_id              | uuid      | fk         |
| secret_encrypted     | text(512) | —          |
| recovery_codes_hash  | text      | —          |
| is_active            | boolean   | —          |
| created_at           | datetime  | —          |
| updated_at           | datetime  | nil        |

### user_login_providers - table

| Field Name               | Data Type   | Attributes |
| ------------------------ | ----------- | ---------- |
| id                       | uuid        | pk auto++  |
| user_id                  | uuid        | fk         |
| provider_id              | text(64)    | —          |
| account_id               | text(128)   | —          |
| access_token             | text        | —          |
| refresh_token            | text        | —          |
| id_token                 | text        | —          |
| access_token_expires_at  | datetime    | nil        |
| refresh_token_expires_at | datetime    | nil        |
| scopes                   | text        | —          |
| created_at               | datetime    | —          |
| updated_at               | datetime    | nil        |

### user_passkey - tabe

| Field Name    | Data Type    | Attributes |
| ------------- | ------------ | ---------- |
| id            | uuid         | pk         |
| user_id       | uuid         | fk         |
| credential_id | binary(1024) | —          |
| data          | json         | —          |

### user_claims - table

| Field Name  | Data Type | Attributes       |
| ----------- | --------- | ---------------- |
| id          | i32       | pk, auto++       |
| user_id     | uuid      | fk to users.id   |
| type        | text(256) | —                |
| value       | text      | —                |

### roles - table

Defines roles that can be assigned to users.

| Field Name  | Data Type | Attributes       |
| ----------- | --------- | ---------------- |
| id          | uuid      | pk auto++        |
| name        | text(64)  | —                |
| name_upcase | text(64)  | uniq, ix         |
| desc        | text(512) | —                |

### users_roles - table

Many-to-many relationship between users and roles.

| Field Name | Data Type | Attributes       |
| ---------- | --------- | ---------------- |
| user_id    | uuid      | pk, fk           |
| role_id    | uuid      | pk, fk           |

**Primary Key:** Composite (user_id, role_id)

### role_claims - table

Claims (permissions) assigned to a role.

| Field Name  | Data Type | Attributes         |
| ----------- | --------- | ------------------ |
| id          | i32       | pk, auto ++        |
| role_id     | uuid      | fk                 |
| type        | text      | —                  |
| claim_value | text      | —                  |

### user_api_keys - table

| Field Name               | Data Type   | Attributes |
| ------------------------ | ----------- | ---------- |
| id                       | int         | pk auto++  |
| user_id                  | uuid        | fk         |
| name                     | text(64)    |            |
| name_upcase              | text(64)    |            |
| key_hint                 | text(16)    | —          |
| key_digest               | text(1024)  | —          |
| expires_at               | datetime    | nil        |

*Note on API Keys:* API Keys should be generated with a standard format (e.g., `hypr_...`). The `key_hint` field stores a short, non-sensitive prefix/suffix (e.g., `...aB3`) to display in the UI, while the actual key is hashed via PBKDF2 into `key_digest`.
| is_locked                | bit         |            |
| is_revoked               | bit         |            |
| comment                  | text(512)   | nil        |

ix for user_id, name_upcase

### user_api_keys_roles - table

Many-to-many relationship between user_api_keys and roles.

| Field Name       | Data Type | Attributes       |
| ---------------- | --------- | ---------------- |
| user_api_key_id  | int       | pk, fk           |
| role_id          | uuid      | pk, fk           |

## Routes

There should be 3 areas

- auth
- users
- roles

Auth should contain things like like login, logout, forgot-password, forgot-email, etc

Users should be more of a collection resource but also have a unique section for me .e.g
`GET api/v1/users/me`  any user that has a valid login should be able to use endpoints
that end with `/me`.

For endpoints that do no use /me: users requires `users:read` or `users:write`.  Roles requires `roles:read` or `roles:write`
permissions.

propose api routes.

## Migrations

Add golang-migrator to the project and create sub directories for mssql, mysql, pg, and sqlite.
The app should default to sqlite.  Propose and enable configuration settings for mssql, mysql, pg
and sqlite for the application.

## Databases

install modules that do not require CGO for database drivers.

Create implementations for each db provider for each store.

## Repositories

Call them stores instead of repositories.  Propose store apis for UserStore and RoleStore
interfaces.

## Resusable logic

If possible create reusable logic for:

- for imports/exports for json, csv, excel, and yaml.
- for pagination
- for converting errors and validation errors to response errors.
- for sorting
- for filtering
- for expanding

## Make recommendations

- make recommendations for indexes, contraints, additional columns.
- **Primary Constraints**: To enforce that only one email/phone can be primary per user, use a partial unique index on `(user_id) WHERE is_primary = true` at the database level for Postgres, MSSQL, and SQLite. For MySQL, enforce this logic at the application level to avoid triggers.