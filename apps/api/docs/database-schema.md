# Database Schema

## Identity

### `groups`
| Field          | Common Type | Attributes      | Notes                           |
|----------------|-------------|-----------------|---------------------------------|
| id             | uuid        | pk              | Primary key                     |
| name           | text(256)   | uniq            | Group name                      |
| name_upcase    | text(256)   | uniq, ix        | Uppercased group name           |
| description    | text        | nil             | Group description               |
| image_uri      | text        | nil             | Group avatar or icon URI        |
| email          | text(256)   | uniq, nil       | Group email address             |
| email_upcase   | text(256)   | uniq, ix, nil   | Uppercased group email address  |
| is_active      | boolean     | df(1)           | Whether the group is active     |
| created_at     | datetime    |                 | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil             | Timestamp of last update (UTC)  |

- `ix_groups_name_upcase` (`name_upcase`)
- `ix_groups_email_upcase` (`email_upcase`)

### `groups_users`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| group_id       | uuid        | pk, fk cascade(delete) | Group ID                        |
| user_id        | uuid        | pk, fk cascade(delete) | User ID                         |

- `fk` to `groups(id)`
- `fk` to `users(id)`

### `groups_admins`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| group_id       | uuid        | pk, fk cascade(delete) | Group ID                        |
| user_id        | uuid        | pk, fk cascade(delete) | Admin User ID                   |

- `fk` to `groups(id)`
- `fk` to `users(id)`

### `groups_roles`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| group_id       | uuid        | pk, fk cascade(delete) | Group ID                        |
| role_id        | uuid        | pk, fk cascade(delete) | Role ID                         |

- `fk` to `groups(id)`
- `fk` to `roles(id)`

### `projects`
| Field          | Common Type | Attributes      | Notes                           |
|----------------|-------------|-----------------|---------------------------------|
| id             | uuid        | pk              | Primary key                     |
| name           | text(256)   |                 | Project name                    |
| name_upcase    | text(256)   | ix              | Uppercased project name         |
| slug           | text(256)   | uniq, ix        | URL-friendly identifier         |
| description    | text        | nil             | Project description             |
| is_active      | boolean     | df(1)           | Whether the project is active   |
| created_at     | datetime    |                 | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil             | Timestamp of last update (UTC)  |

- `ix_projects_name_upcase` (`name_upcase`)
- `ix_projects_slug` (`slug`)

### `projects_groups`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| project_id     | uuid        | pk, fk cascade(delete) | Project ID                      |
| group_id       | uuid        | pk, fk cascade(delete) | Group ID                        |
| permissions    | bigint      | df(0)                  | Bitmask of group permissions    |

- `fk` to `projects(id)`
- `fk` to `groups(id)`

## Environments & Configuration

### `environments`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | Environment name                |
| name_upcase    | text(256)   |                        | Uppercased environment name     |
| description    | text        | nil                    | Environment description         |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_environments_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `config_file_names`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | Config file logical name        |
| name_upcase    | text(256)   |                        | Uppercased config file name     |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_config_file_names_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `config_files`
| Field               | Common Type | Attributes             | Notes                           |
|---------------------|-------------|------------------------|---------------------------------|
| id                  | uuid        | pk                     | Primary key                     |
| config_file_name_id | uuid        | fk cascade(delete)     | Logical config file ID          |
| environment_id      | uuid        | fk cascade(delete)     | Environment ID                  |
| content             | text        |                        | The raw configuration content   |
| created_at          | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at          | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_config_files_name_env` (`config_file_name_id`, `environment_id`) uniq
- `fk` to `config_file_names(id)`
- `fk` to `environments(id)`

### `env_variable_names`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | Env variable logical name       |
| name_upcase    | text(256)   |                        | Uppercased env variable name    |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_env_variable_names_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `env_variables`
| Field                | Common Type | Attributes             | Notes                           |
|----------------------|-------------|------------------------|---------------------------------|
| id                   | uuid        | pk                     | Primary key                     |
| env_variable_name_id | uuid        | fk cascade(delete)     | Logical env variable ID         |
| environment_id       | uuid        | fk cascade(delete)     | Environment ID                  |
| value                | text        |                        | Env variable string value       |
| created_at           | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at           | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_env_variables_name_env` (`env_variable_name_id`, `environment_id`) uniq
- `fk` to `env_variable_names(id)`
- `fk` to `environments(id)`

### `secret_names`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | Secret logical name             |
| name_upcase    | text(256)   |                        | Uppercased secret name          |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_secret_names_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `secrets`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| secret_name_id | uuid        | fk cascade(delete)     | Logical secret ID               |
| environment_id | uuid        | fk cascade(delete)     | Environment ID                  |
| value          | text        |                        | Encrypted secret value          |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_secrets_name_env` (`secret_name_id`, `environment_id`) uniq
- `fk` to `secret_names(id)`
- `fk` to `environments(id)`

### `certificate_names`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | Certificate logical name        |
| name_upcase    | text(256)   |                        | Uppercased certificate name     |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_certificate_names_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `certificates`
| Field               | Common Type | Attributes             | Notes                           |
|---------------------|-------------|------------------------|---------------------------------|
| id                  | uuid        | pk                     | Primary key                     |
| certificate_name_id | uuid        | fk cascade(delete)     | Logical certificate ID          |
| environment_id      | uuid        | fk cascade(delete)     | Environment ID                  |
| certificate_data    | text        |                        | Public certificate data         |
| private_key_secret  | text        | nil                    | Encrypted private key data      |
| created_at          | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at          | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_certificates_name_env` (`certificate_name_id`, `environment_id`) uniq
- `fk` to `certificate_names(id)`
- `fk` to `environments(id)`

### `ssh_key_names`
| Field          | Common Type | Attributes             | Notes                           |
|----------------|-------------|------------------------|---------------------------------|
| id             | uuid        | pk                     | Primary key                     |
| project_id     | uuid        | fk cascade(delete)     | Project ID                      |
| name           | text(256)   |                        | SSH key logical name            |
| name_upcase    | text(256)   |                        | Uppercased SSH key name         |
| created_at     | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at     | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_ssh_key_names_project_id_name_upcase` (`project_id`, `name_upcase`) uniq
- `fk` to `projects(id)`

### `ssh_keys`
| Field               | Common Type | Attributes             | Notes                           |
|---------------------|-------------|------------------------|---------------------------------|
| id                  | uuid        | pk                     | Primary key                     |
| ssh_key_name_id     | uuid        | fk cascade(delete)     | Logical SSH key ID              |
| environment_id      | uuid        | fk cascade(delete)     | Environment ID                  |
| public_key          | text        |                        | Public key data                 |
| private_key_secret  | text        |                        | Encrypted private key data      |
| created_at          | datetime    |                        | Timestamp of creation (UTC)     |
| updated_at          | datetime    | nil                    | Timestamp of last update (UTC)  |

- `ix_ssh_keys_name_env` (`ssh_key_name_id`, `environment_id`) uniq
- `fk` to `ssh_key_names(id)`
- `fk` to `environments(id)`
