# API Routes

## Auth (`/api/v1/auth`)

- `POST /login` - Standard login with email and password.
- `POST /logout` - Logout a user by invalidating the active session.
- `POST /forgot-password` - Request a password reset link.
- `POST /reset-password` - Reset a password using the secure OTP token.

## Users (`/api/v1/users`)

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
