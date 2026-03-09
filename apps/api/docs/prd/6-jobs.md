# jobs

A job can be a simple cron job that runs a script, hits an endpoint, or runs a workflow.  Jobs can be scheduled, triggered by an event such as a webhook,
or run immediately.  Jobs can also be associated with devices for running scripts on those devices.  Jobs can also be associated with projects for organizational purposes.  Jobs can also have a time range for when they are active.  For example, a job may only be active for a week or a month or it may be active indefinitely.  Jobs can also have a type field to denote what type of job it is, such as endpoint, script, or workflow.  This allows us to have different fields for different types of jobs and also allows us to easily query for jobs of a certain type.

We'll create a PRD/spec for how the workflow engine and workflow definitions should look and function later, but for now we'll just have a simple content field for the workflow definition and an engine_type field to denote what type of workflow engine it is for.  By default, we'll create our own workflow engine and use yaml for the workflow definitions, but we also want to allow for other workflow engines and definitions as well.  This will allow us to be flexible and integrate with other tools and platforms that may have their own workflow engines and definitions.

## jobs - table

| Field Name      | Data Type | Attributes |
| --------------- | --------- | ---------- |
| id              | uuid      | pk         |
| name            | text(128) |            |
| name_upcase     | text(128) | uniq       |
| desc            | text(512) | nil        |
| type            | text(32)  |            |
| cron            | text(64)  | nil        |
| schedule        | text(64)  | nil        |
| ends_at         | datetime  | nil        |
| starts_at       | datetime  | nil        |
| is_active       | bit       |            |
| project_id      | uuid      | fk, nil    |
| job_endpoint_id | uuid      | fk, nil    |
| workflow_id     | uuid      | fk, nil    |
| script_id       | uuid      | fk, nil    |

type can be endpoint, script, or workflow.  if type is endpoint, job_endpoint_id must have a value. if type is workflow, workflow_id must have a value. if type is script, script_id must have a value.  cron is the cron expression for when the job should run. schedule is a human readable version of the cron expression for display purposes.

Starts at and ends at can be used to specify a time range for when the job should be active.  if starts_at is nil, the job is active immediately. if ends_at is nil, the job is active indefinitely.

## job_endpoints - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| id         | uuid      | pk         |
| url        | text(255) |            |
| method     | text(8)   |            |
| headers    | json      | nil        |
| body       | text      | nil        |
| timeout    | u16       | nil        |

## job_scripts - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| job_id     | uuid      | pk, fk     |
| script_id  | uuid      | pk, fk     |

## job_devices - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| job_id     | uuid      | pk, fk     |
| device_id  | uuid      | pk, fk     |

only needed if the job is meant to run a script on a device.  this allows us to associate a job with multiple devices and scripts and also allows us to track which devices and scripts are associated with which jobs for auditing and management purposes.

## scripts - table

| Field Name   | Data Type | Attributes |
| ------------ | --------- | ---------- |
| id           | uuid      | pk         |
| name         | text(128) |            |
| name_upcase  | text(128) | uniq       |
| desc         | text(512) | nil        |
| content      | text      | nil        |
| content_type | text(64)  | nil        |
| is_template  | bit       |            |
| exec         | text(64)  | nil        |
| is_active    | bit       |            |
| created_at   | datetime  |            |
| metadata     | json      | nil        |

## workflows - table

| Field Name  | Data Type | Attributes |
| ----------- | --------- | ---------- |
| id          | uuid      | pk         |
| name        | text(128) |            |
| name_upcase | text(128) | uniq       |
| desc        | text(512) | nil        |
| content     | text      | nil        |
| engine_type | text(32)  |            |
| is_active   | bit       |            |
| created_at  | datetime  |            |

by default this will be yaml and we'll provide our own workflow engine, but we also allow other engine types and talk to that engine.

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
