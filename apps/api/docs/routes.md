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
