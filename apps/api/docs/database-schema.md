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
