# Epic: Projects Management

## Description
Manage projects that act as organizational boundaries. Projects contain resources like secrets, configurations, environments, deployments, devices, and jobs. Groups are assigned to projects with a permission bitmask.

## Issues / Tasks
1. **[Project CRUD API]** Create standard endpoints for Projects (`GET`, `POST`, `PUT`, `DELETE`).
2. **[Project Slug API]** Enable retrieval of projects by `slug` for URLs.
3. **[Project Groups API]** Add endpoints to manage group associations (`projects_groups` table) and update their `perms` bitmask.
4. **[Permission Resolution Logic]** Implement a utility/middleware that calculates the highest bitwise permission level (`NONE`, `READ`, `WRITE`, `EXEC`, `AUDIT`, `ALL`) when a user is in multiple groups mapped to a project.

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
