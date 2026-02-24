# PRD: Security, Environments, and Configuration

if the project_id is nil, then the config_file, env_variable, or secret is global and will be
assigned to all projects. projects may use the values, but not edit the global values.

Part of the api should include the ability to import and export config files, env variables, secrets, certificates, and ssh keys in bulk via csv or json. this will allow for easier management of these resources and also allow for easier onboarding of new projects and environments.

A user should be able to generate secrets, certificates, and ssh keys from the api as well. this will allow for easier management of these resources and also allow for easier onboarding of new projects and environments.

Long term, the tables for secrets, certificates, and ssh keys may be ignored when the user opts to use their own providers such as vault, aws secrets manager, acme certs, etc. iWe shouldn't add those depenencies now or directly to the api, but instead we'll need to think about how to enable modules/plugins for those whether as a cli or as a different
service that can integrate with the api.  for now, we'll just store the secrets, certificates, and ssh keys in the database and encrypt them at rest.

The values in this file must be encrypted at rest and decrypted when accessed via the api. There should be options in the api to send the data as base64 encoded, plain text, or encrypted, should the client want to handle the encryption and decryption on their end or decode to bytes in memory rather having it plain text in memory.  

## envirronment - table

| Field Name  | Data Type | Attributes |
| ----------- | --------- | ---------- |
| id          | i16       | pk auto++  |
| name        | text(128) |            |
| name_upcase | text(128) | uniq       |
| abbr        | text(8)   | nil        |
| desc        | text(512) | nil        |

## config_files - table

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | i32       | pk         |
| name           | text(256) |            |
| name_upcase    | text(256) |            |
| content        | bin       |            |
| encoding       | text(24)  |            |
| content_type   | text(128) | nil        |
| is_template    | bit       |            |
| is_active      | bit       |            |
| project_id     | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| created_at     | datetime  |            |
| updated_at     | datetime  | nil        |

uniq on project_id and name. is_template means that it can be used by gotmpl or another templating engine to generate a config file for a project. content is the raw content of the file. encoding is the encoding of the content such as utf-8, base64, etc. content_type is the mime type of the content such as text/plain, application/json, etc. if content_type is nil, it will be determined by the file extension in the name field.

## config_file_history - table

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | uuid      | pk         |
| config_file_id | i32       | fk         |
| name           | text(256) |            |
| content        | bin       |            |
| encoding       | text(24)  |            |
| content_type   | text(128) | nil        |
| project_id     | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| is_template    | bit       |            |
| created_at     | datetime  |            |

## env_variables

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | i32       | pk         |
| name           | text(255) |            |
| value          | text      |            |
| is_active      | bit       |            |
| project_id     | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| expires_at     | datetime  | nil        |
| created_at     | datetime  |            |
| updated_at     | datetime  | nil        |

## env_variable_history

| Field Name      | Data Type | Attributes |
| --------------- | --------- | ---------- |
| id              | uuid      | pk         |
| env_variable_id | i32       | fk         |
| name            | text(255) |            |
| value           | text      |            |
| project_id      | uuid      | fk, nil    |
| environment_id  | i16       | fk, nil    |
| expires_at      | datetime  | nil        |
| created_at      | datetime  |            |

## secret_names

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | i32       | pk         |
| name           | text(255) |            |
| name_upcase    | text(255) |            |
| project_id     | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| created_at     | datetime  |            |

## secrets

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | uuid      | pk         |
| secret_name_id | i32       | fk         |
| value          | text      |            |
| is_active      | bit       |            |
| is_primary     | bit       |            |
| project_id     | uuid      | fk, nil    |
| expires_at     | datetime  | nil        |
| created_at     | datetime  |            |
| updated_at     | datetime  | nil        |

history is builtin to secrets and you can reference an older uuid
if needed. value is encrypted.

## certificate_names

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | i32       | pk         |
| name           | text(255) |            |
| name_upcase    | text(255) |            |
| project_id     | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| created_at     | datetime  |            |

represents a name for a certificate.

## certificates

| Field Name          | Data Type | Attributes |
| ------------------- | --------- | ---------- |
| id                  | uuid      | pk         |
| certificate_name_id | i32       | fk         |
| subject             | text(512) | nil        |
| issuer              | text(512) | nil        |
| public_key          | text      | nil        |
| private_key         | text      | nil        |
| csr                 | text      | nil        |
| password            | text      | nil        |
| is_active           | bit       |            |
| is_primary          | bit       |            |
| expires_at          | datetime  | nil        |
| not_before_at       | datetime  | nil        |
| thumprint           | text(128) | nil        |
| created_at          | datetime  |            |
| updated_at          | datetime  | nil        |

## ssh_key_names

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | i32       | pk         |
| name           | text(255) |            |
| name_upcase    | text(255) |            |
| project_id     | uuid      | fk, nil    |
| user_id        | uuid      | fk, nil    |
| environment_id | i16       | fk, nil    |
| created_at     | datetime  |            |

A ssh key name represents a name for a ssh key. it may be associated with a project, user, or be global. if associated with a user, it is only for that user. if associated with a project, it is for all users of that project. if not associated with either, it is global and can be used by any project
or resource that requires it.

## ssh_keys

| Field Name      | Data Type | Attributes |
| --------------- | --------- | ---------- |
| id              | uuid      | pk         |
| ssh_key_name_id | i32       | fk         |
| public_key      | text      | nil        |
| private_key     | text      | nil        |
| password        | text      | nil        |
| is_active       | bit       |            |
| is_primary      | bit       |            |
| expires_at      | datetime  | nil        |
| created_at      | datetime  |            |
| updated_at      | datetime  | nil        |

password and private key are encrypted if they have values. is_primary is used to notate the primary certificate for a certificate name. this allows us to have multiple certificates for a certificate name and track history.
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
