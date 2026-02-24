# devices

The use case for storing device information is not only for deployments but managing devices, primarily servers, within a given customer's
environment. We'll also enable it to start remote ssh sessions in the website directly to the server.

The user should also be able to use the API to pull ssh info or rdp info information and connect to the device via ssh or rdp.  this will allow for easier management of devices 

The various attributes in the devices table will allow us to have a better inventory of the devices and also allow for better targeting of deployments and jobs to specific devices based on their attributes.  for example, a user may want to deploy a certain app only to devices that have a certain os_platform or os_version.  this also allows for better management of devices and tracking of their attributes over time.

The metadata will be used to store key facts about the device that isn't already captured in the other fields.  this can include things like location, owner, purpose, etc.  this will allow for better organization and management of devices and also allow for better targeting of deployments and jobs to specific devices based on their metadata.

Inventory of apps and other things shouldn't be stored in the device metadata. We'll create a separate table later for tracking what apps and other things are on the device.  this will allow for better tracking and management of the apps and other things on the device and also allow for better targeting of deployments and jobs to specific devices based on the apps and other things they have.

## devices - table

| Field Name     | Data Type | Attributes |
| -------------- | --------- | ---------- |
| id             | uuid      | pk         |
| name           | text(128) |            |
| name_upcase    | text(128) | uniq       |
| desc           | text(512) | nil        |
| hostname       | text(128) | nil        |
| ipv4_address   | text(15)  | nil        |
| ipv6_address   | text(39)  | nil        |
| fqdn           | text(255) | nil        |
| host_id_type   | text(8)   |            |
| ssh_port       | u16       | nil        |
| rdp_port       | u16       | nil        |
| is_active      | bit       |            |
| mac_address    | text(17)  | nil        |
| os_platform    | text(64)  | nil        |
| os_version     | text(64)  | nil        |
| os_arch        | text(64)  | nil        |
| os_codename    | text(64)  | nil        |
| os_id          | text(64)  | nil        |
| os_variant     | text(64)  | nil        |
| os_name        | text(128) | nil        |
| os_pretty_name | text(128) |            |
| metadata       | json      | nil        |

host_id_type can be hostname, ipv4, ipv6, or fqdn.  this allows us to have multiple devices with the same hostname as long as they have different host_id_types and values.  metadata can be used to store any additional information about the device such as location, owner, purpose, etc.

## device_ssh_auth - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| device_id  | uuid      | pk, fk     |
| ssh_key_id | uuid      | pk, fk     |
| username   | text(64)  | nil        |
| password   | text(128) | nil        |

password is encrypted if it has a value. this allows us to store credentials for devices that require them. if username and password are nil, it means that the device can be accessed without credentials or with a certificate.

## device_tags - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| device_id  | uuid      | pk, fk     |
| tag_id     | uuid      | pk, fk     |
| value      | text(128) | nil        |

## device_rdp_auth - table

| Field Name | Data Type | Attributes |
| ---------- | --------- | ---------- |
| device_id  | uuid      | pk, fk     |
| username   | text(64)  | nil        |
| password   | text(128) | nil        |
| domain     | text(128) | nil        |
| rdpfile    | text      | nil        |

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
