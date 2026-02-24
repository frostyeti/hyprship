# Epic: Groups Management

## Description
Manage groups of users for notifications, role associations, and permissions. Groups can be linked to projects and resources, enabling granular access control.

## Issues / Tasks
1. **[Group CRUD API]** Create standard endpoints for Groups (`GET`, `POST`, `PUT`, `DELETE`).
2. **[Group Users API]** Implement endpoints for adding/removing users from a group (`groups_users` table).
3. **[Group Admins API]** Implement endpoints to manage group administrators (`groups_admins` table).
4. **[Group Roles API]** Implement endpoints for mapping roles to a group (`groups_roles` table) for claims integration.
5. **[Claims Resolution]** Update authentication/claims middleware to inject group and role arrays into the active user context.

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
