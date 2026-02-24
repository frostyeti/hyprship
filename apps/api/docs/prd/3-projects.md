# projects

Projects will have resources added to it later such as secrets,
configuration, environments, deployment, devices and jobs. 

## projects - table

| Field Name  | Data Type | Attributes |
| ----------- | --------- | ---------- |
| id          | uuid      | pk         |
| name        | text(128) |            |
| name_upcase | text(128) | uniq       |
| slug        | text(128) | uniq       |
| desc        | text(512) | nil        |

slug will be used for urls

## projects_groups - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| project_id | uuid      | pk, fk     |
| group_id   | uuid      | pk, fk     |
| perms      | u16       |            |

perms would be an enum.  if a user is in multiple groups, the
highest permission set wins out.

```golang
const (
   NONE = 0
   READ = 1 << 0,  // 1
   WRITE = 1 << 1,  // 2
   EXEC = 1 << 2,  // 4
   AUDIT = 1 << 3,  // 8
   ALL = READ | WRITE | EXEC | AUDIT
)
```
## Table Design Recommendations
* **Auditing:** Include standard auditing fields (`created_at`, `updated_at`, `deleted_at`).
* **Bitmask Constraint:** The `perms` column in `projects_groups` should enforce valid enum values at the database level if possible (e.g. `CHECK (perms IN (...))`), or rely on the application layer strictly validating the value before insert/update.

## Proposed API Routes
* `GET /api/v1/projects` (Supports filtering, sorting, paging)
* `GET /api/v1/projects/{id}`
* `GET /api/v1/projects/slug/{slug}`
* `POST /api/v1/projects`
* `PUT /api/v1/projects/{id}`
* `DELETE /api/v1/projects/{id}`
* `POST /api/v1/projects/{id}/groups` (Assign a group with a specific `perms` level)
* `PUT /api/v1/projects/{id}/groups/{group_id}` (Update `perms` for a group)
* `DELETE /api/v1/projects/{id}/groups/{group_id}`

## Functionality Beyond CRUD
* **Evaluated Permissions Endpoint:** An endpoint (e.g., `GET /api/v1/projects/{id}/my-permissions`) that quickly calculates and returns a bitmask (or list of string constants) of the requesting user's effective permissions based on their active groups.
* **Archiving/Deactivation:** A soft delete strategy that immediately disables read/write functionality on underlying resources (secrets, jobs) when a project is deleted/archived.
