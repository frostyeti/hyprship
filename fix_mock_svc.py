import re

with open("apps/api/internal/routes/v1/auth_test.go", "r") as f:
    content = f.read()

impl = """func (m *MockIdentityService) SetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return m.SetPasswordFn(ctx, userID, newPassword)
}

func (m *MockIdentityService) CreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error) {
	return nil, nil
}"""

content = content.replace("""func (m *MockIdentityService) SetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	return m.SetPasswordFn(ctx, userID, newPassword)
}""", impl)

with open("apps/api/internal/routes/v1/auth_test.go", "w") as f:
    f.write(content)
