# Epic: Jobs & Workflows

## Description
Manage and execute jobs ranging from cron tasks to full-fledged workflows and scripts. Jobs can be targeted globally, to projects, or specific devices.

## Issues / Tasks
1. **[Jobs API]** Basic CRUD for jobs (`GET`, `POST`, `PUT`, `DELETE`).
2. **[Job Endpoints API]** Manage endpoint jobs (`job_endpoints` table).
3. **[Job Scripts API]** Manage script jobs (`job_scripts` table) and their scripts (`scripts` table).
4. **[Job Devices API]** Associate jobs with devices (`job_devices` table) to track execution targets.
5. **[Workflows API]** Manage workflows (`workflows` table) and engine types (default YAML).
6. **[Job Executions/Runs Tracker]** Add a run tracker table and API to view logs, status, and duration of past executions.

## Table Design Recommendations
* **Auditing:** Include standard auditing fields (`created_at`, `updated_at`, `deleted_at`).
* **Run History Table:** Add a `job_runs` table containing `id`, `job_id`, `status` (`pending`, `running`, `success`, `failed`), `started_at`, `finished_at`, `output_log` (or pointer to log storage), and `error_message`.
* **UUIDs:** `script_id` and `workflow_id` should correctly map to `uuid` and be indexed for fast relational queries.
* **Content:** Consider `workflows.content` as `jsonb` or `json` if structured internally, or `text` if fully parsing external YAML definitions.

## Proposed API Routes
* `GET /api/v1/jobs`
* `GET /api/v1/jobs/{id}`
* `POST /api/v1/jobs`
* `PUT /api/v1/jobs/{id}`
* `DELETE /api/v1/jobs/{id}`
* `POST /api/v1/jobs/{id}/run`
* `GET /api/v1/jobs/{id}/runs`
* `GET /api/v1/jobs/{id}/runs/{run_id}`
* `POST /api/v1/scripts`
* `GET /api/v1/scripts/{id}`
* `POST /api/v1/workflows`
* `GET /api/v1/workflows/{id}`

## Functionality Beyond CRUD
* **Ad-hoc Execution:** Allow triggering any active job out of schedule via `/api/v1/jobs/{id}/run`.
* **Real-time Log Streaming:** Output stdout/stderr of executing scripts or workflow steps back to the UI via Server-Sent Events (SSE) or WebSockets (`/ws/v1/jobs/{id}/runs/{run_id}/logs`).
* **Webhook Triggers:** Generate unique URLs for triggering jobs remotely from an external system.
* **Golang Open Source Recommendations (Workflows):** Look into using **Temporal** (`temporalio/temporal`) or **Argo** (`argoproj/argo-workflows`) if the built-in YAML execution engine needs to grow into reliable, stateful workflow orchestration.
