# API Routes

## Auth (`/api/v1/auth`)

- `POST /login` - Standard login with email and password.
- `POST /logout` - Logout a user by invalidating the active session.
- `POST /forgot-password` - Request a password reset link.
- `POST /reset-password` - Reset a password using the secure OTP token.
- `POST /forgot-email` - Request a forgotten email retrieval SMS.

## Users (`/api/v1/users`)

- `GET /me` - Get current authenticated user profile.
- `PUT /me` - Update current authenticated user profile.
- `GET /` - List users.
- `GET /{id}` - Get user by ID.
- `POST /` - Create a new user.
- `PUT /{id}` - Update a user by ID.
- `DELETE /{id}` - Delete a user by ID.

## Roles (`/api/v1/roles`)

- `GET /` - List roles.
- `GET /{id}` - Get a role by ID.
- `POST /` - Create a new role.
- `PUT /{id}` - Update a role by ID.
- `DELETE /{id}` - Delete a role by ID.

## Groups (`/api/v1/groups`)

- `GET /` - List groups.
- `GET /{id}` - Get a group by ID.
- `POST /` - Create a new group.
- `PUT /{id}` - Update a group by ID.
- `DELETE /{id}` - Delete a group by ID.

## Projects (`/api/v1/projects`)

- `GET /` - List projects.
- `GET /{id}` - Get a project by ID or slug.
- `POST /` - Create a new project.
- `PUT /{id}` - Update a project by ID.
- `DELETE /{id}` - Delete a project by ID.
- `GET /{id}/groups` - List groups associated with a project.
- `POST /{id}/groups/{groupId}` - Add a group to a project with permissions.
- `PUT /{id}/groups/{groupId}` - Update a group's permissions in a project.
- `DELETE /{id}/groups/{groupId}` - Remove a group from a project.

### Project Environments (`/api/v1/projects/{id}/environments`)
- `GET /` - List environments for a project.
- `GET /{envId}` - Get an environment by ID or name.
- `POST /` - Create a new environment.
- `PUT /{envId}` - Update an environment by ID.
- `DELETE /{envId}` - Delete an environment by ID.

### Project Configuration Names (`/api/v1/projects/{id}/config-files` & `/env-variables`)
- `GET /config-files` - List configuration file names for a project.
- `POST /config-files` - Create a new configuration file name.
- `DELETE /config-files/{nameId}` - Delete a configuration file name.
- `GET /env-variables` - List environment variable names for a project.
- `POST /env-variables` - Create a new environment variable name.
- `DELETE /env-variables/{nameId}` - Delete an environment variable name.

### Project Secrets Names (`/api/v1/projects/{id}/secrets` & `/certificates` & `/ssh-keys`)
- `GET /secrets` - List secret names for a project.
- `POST /secrets` - Create a new secret name.
- `DELETE /secrets/{nameId}` - Delete a secret name.
- `GET /certificates` - List certificate names for a project.
- `POST /certificates` - Create a new certificate name.
- `DELETE /certificates/{nameId}` - Delete a certificate name.
- `GET /ssh-keys` - List SSH key names for a project.
- `POST /ssh-keys` - Create a new SSH key name.
- `DELETE /ssh-keys/{nameId}` - Delete an SSH key name.

### Environment Specific Values (`/api/v1/projects/{id}/environments/{envId}`)
- `GET /config-files/{nameId}` - Get the contents of a configuration file for this environment.
- `PUT /config-files/{nameId}` - Create or update a configuration file for this environment.
- `GET /env-variables/{nameId}` - Get the value of an environment variable.
- `PUT /env-variables/{nameId}` - Create or update an environment variable.
- `GET /secrets/{nameId}` - Get a secret value.
- `PUT /secrets/{nameId}` - Create or update a secret value.
- `GET /certificates/{nameId}` - Get a certificate.
- `PUT /certificates/{nameId}` - Create or update a certificate.
- `GET /ssh-keys/{nameId}` - Get an SSH key.
- `PUT /ssh-keys/{nameId}` - Create or update an SSH key.
