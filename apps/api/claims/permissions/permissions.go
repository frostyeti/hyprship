package permissions

// Standard permissions
const (
	Admin = "ADMIN"
	Read  = "READ"
	Write = "WRITE"
	Exec  = "EXEC"
)

// Identity specific permissions
const (
	UsersRead  = "USERS_READ"
	UsersWrite = "USERS_WRITE"
	RolesRead  = "ROLES_READ"
	RolesWrite = "ROLES_WRITE"
)
