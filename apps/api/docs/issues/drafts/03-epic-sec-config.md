# Epic: Security, Environments, and Configuration

## Description
Manage configuration files, environment variables, secrets, certificates, and SSH keys. Items can be global or scoped to a specific project. Data is encrypted at rest and supports advanced retrieval methods.

## Issues / Tasks
1. **[Environment API]** Basic CRUD for environments.
2. **[Config Files API]** Manage configuration files, templates, and history. 
3. **[Environment Variables API]** Manage environment variables and history.
4. **[Secrets API]** Manage secrets, secret names, and history. Value MUST be encrypted.
5. **[Certificates API]** Manage certificate names, history, private keys (encrypted at rest), and CSRs.
6. **[SSH Keys API]** Manage SSH keys (user and project-level), public/private keys (encrypted at rest).
7. **[Crypto Engine]** Implement AES-GCM encryption/decryption drivers within the API layer for sensitive values.

## Table Design Recommendations
* **Primary Key Consistency:** Consider standardizing all primary keys to `uuid` to prevent sequence guessing and facilitate easier cross-table joins/references. `environment`, `config_files`, `env_variables`, `secret_names`, `certificate_names`, and `ssh_key_names` are currently integers.
* **Sensitive Names:** Columns that store encrypted text (e.g. `secrets.value`, `device_ssh_auth.password`, `certificates.private_key`) should probably be appended with `_encrypted` to follow the `_digest` hashing convention mentioned in the project instructions.
* **Auditing/Tracking:** Add `created_by`, `updated_by` to history tables so that changes to configuration/secrets can be tracked securely to the exact user/process that triggered it.
* **Soft Deletes:** Ensure soft-deletion of `config_files`, `env_variables`, `secrets`, etc.

## Proposed API Routes
* `GET /api/v1/environments`
* `GET /api/v1/projects/{id}/config-files`
* `POST /api/v1/projects/{id}/config-files`
* `GET /api/v1/config-files/{id}`
* `POST /api/v1/config-files/import`
* `GET /api/v1/config-files/export`
* `GET /api/v1/projects/{id}/secrets`
* `POST /api/v1/secrets/generate`
* `POST /api/v1/certificates/generate`
* `POST /api/v1/ssh-keys/generate`

## Functionality Beyond CRUD
* **Bulk Imports/Exports:** Accept JSON, CSV, Excel, YAML via a bulk import/export process (especially utilizing ZIP archives when dealing with multiple configuration files).
* **On-the-fly Decryption Engine:** Provide headers (`Accept-Encoding: base64/plaintext`) or query parameters allowing clients to dictate whether the payload should be decrypted into memory or returned as base64-encoded bytes.
* **Template Rendering:** Create a render endpoint (`POST /api/v1/config-files/{id}/render`) that accepts context variables and outputs the rendered text via a Go text/template engine.
* **Automatic Expiration Rotation:** Send notifications or events to message queues when certificates, secrets, or SSH keys are nearing their `expires_at` date.
