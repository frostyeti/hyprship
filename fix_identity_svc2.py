with open("apps/api/internal/svc/identity/identity_service.go", "r") as f:
    content = f.read()

import re

# Revert the bad change
content = re.sub(
    r'func \(s \*identityService\) LoginWithPassword\(ctx context\.Context, email string, password string, ipAddress \*string, userAgent \*string\) \(\*models\.UserSession, error\)\n\tCreateSession\(ctx context\.Context, userID uuid\.UUID, ipAddress \*string, userAgent \*string\) \(\*models\.UserSession, error\) \{',
    r'func (s *identityService) LoginWithPassword(ctx context.Context, email string, password string, ipAddress *string, userAgent *string) (*models.UserSession, error) {',
    content
)

# Insert CreateSession in the interface definition correctly
content = re.sub(
    r'(LoginWithPassword\(ctx context\.Context, email string, password string, ipAddress \*string, userAgent \*string\) \(\*models\.UserSession, error\))',
    r'\1\n\tCreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error)',
    content,
    count=1 # Only replace the interface one
)

with open("apps/api/internal/svc/identity/identity_service.go", "w") as f:
    f.write(content)
