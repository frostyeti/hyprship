# Epic: Devices & Remote Access

## Description
Manage device inventories, attributes, and metadata. Enable secure remote SSH/RDP connections to servers from the website directly using stored credentials. Use specific characteristics (OS, location, etc.) for targeted deployments.

## Issues / Tasks
1. **[Devices API]** Basic CRUD for devices including metadata.
2. **[Device Tags API]** Endpoints to add, remove, and update tags on devices.
3. **[Device SSH Auth API]** Manage SSH authentication details (`device_ssh_auth` table).
4. **[Device RDP Auth API]** Manage RDP authentication details (`device_rdp_auth` table).
5. **[Device SSH Proxy Service]** Implement an SSH streaming bridge using WebSockets or SSE for a terminal UI.
6. **[Device RDP Proxy Service]** Implement an RDP streaming bridge (or proxy wrapper) for a web-based RDP client.

## Table Design Recommendations
* **Auditing:** Include standard auditing fields (`created_at`, `updated_at`, `deleted_at`).
* **Metadata Type:** `devices.metadata` should utilize a `jsonb` or optimized document column type in the underlying database for fast filtering and searching.
* **Encrypted Passwords:** Ensure `device_ssh_auth.password` and `device_rdp_auth.password` are AES-GCM encrypted at rest. If internal verification is needed, use PBKDF2 as suggested (`password_digest`).
* **Last Connected:** Add `last_seen_at` or `last_connected_at` column to monitor active devices and ensure inventory freshness.

## Proposed API Routes
* `GET /api/v1/devices` (Supports dynamic filtering on metadata, tags, and OS properties)
* `GET /api/v1/devices/{id}`
* `POST /api/v1/devices`
* `PUT /api/v1/devices/{id}`
* `DELETE /api/v1/devices/{id}`
* `GET /api/v1/devices/{id}/auth/ssh`
* `PUT /api/v1/devices/{id}/auth/ssh`
* `GET /api/v1/devices/{id}/auth/rdp`
* `PUT /api/v1/devices/{id}/auth/rdp`
* `POST /api/v1/devices/{id}/tags`
* `DELETE /api/v1/devices/{id}/tags/{tag_id}`

## Functionality Beyond CRUD
* **Streaming Remote Access (Bastion/Proxy):** The API will act as a Bastion server. Users connect to the API via WebSockets (`/ws/v1/devices/{id}/ssh`), and the API proxies the TCP connection (SSH or RDP) directly to the target device using the stored credentials.
* **Golang Open Source Recommendations (Bastion/RDP):**
  * **Teleport (`gravitational/teleport`)**: The industry standard for identity-aware, Go-based access planes. It natively handles SSH, RDP, Kubernetes, and Databases with built-in RBAC and session recording. Recommending this if full enterprise features are desired.
  * **`golang.org/x/crypto/ssh`**: For a lightweight custom SSH proxy/bridge, use this package directly within the API.
  * **`citilinkru/go-rdp`**: A pure Go RDP library if the API needs to inspect or handle the RDP protocol, though Apache Guacamole is frequently the go-to protocol bridge for rendering RDP in an HTML5 canvas.
* **Device Health Check:** A scheduled task or endpoint (`/api/v1/devices/{id}/ping`) that reaches out to the device's SSH/RDP port to verify connectivity and update `is_active`.
