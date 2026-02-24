# Groups

Enable groups of users for the purpose of notifications and permissions.  A user can be in multiple groups and a group can have multiple users.  A group can also have multiple admins who can manage the group and its users.  Groups can also have roles associated with them for the purpose of claims.  Groups may also be associated with a distro list for notifications.  Groups may also be used for permissions on projects and resources.  For example, a project may have a group that has read access to it and another group that has write access to it.  This allows for more granular control over who has access to what resources and projects and other resources like servers.

## groups - table

| Field Name  | Data Type | Attributes |
| ----------- | --------- | ---------- |
| id          | uuid      | pk         |
| name        | text(128) |            |
| name_upcase | text(128) | uniq       |
| key         | text(32)  | uniq       |
| desc        | text(512) | nil        |
| email       | text(512) | nil        |

key field must be a 32 characters or less short name for the group.
must start with alpha and may container lower alpha, digit, or hypen.
must not end with hyphen.

The key field will be used in claims.

The group may be associated with a distro list hosted somewhere else.

## groups_admins - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| group_id   | uuid      | pk, fk     |
| user_id    | uuid      | pk, fk     |

## groups_users - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| group_id   | uuid      | pk, fk     |
| user_id    | uuid      | pk, fk     |

## groups_roles - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| group_id   | uuid      | pk, fk     |
| role_id    | uuid      | pk, fk     |

useful for applying roles to a group of users for claims.

## Table Design Recommendations
* **Auditing:** Add standard auditing fields to `groups` (`created_at`, `updated_at`, `deleted_at`).
* **Keys:** The `groups.key` column is marked `uniq` and limited to 32 chars. Consider adding a DB-level regular expression constraint or check to enforce the starting alpha requirement and prevent trailing hyphens.

## Proposed API Routes
* `GET /api/v1/groups` (Supports filtering, sorting, paging)
* `GET /api/v1/groups/{id}`
* `POST /api/v1/groups`
* `PUT /api/v1/groups/{id}`
* `DELETE /api/v1/groups/{id}`
* `POST /api/v1/groups/{id}/users`
* `DELETE /api/v1/groups/{id}/users/{user_id}`
* `POST /api/v1/groups/{id}/admins`
* `DELETE /api/v1/groups/{id}/admins/{user_id}`
* `POST /api/v1/groups/{id}/roles`
* `DELETE /api/v1/groups/{id}/roles/{role_id}`

## Functionality Beyond CRUD
* **Bulk User Assignment:** Accept lists of UUIDs or Emails to add/remove users from a group in batch.
* **Distro Integration Sync:** Support webhooks or periodic syncs to update external distro lists based on group membership changes.
