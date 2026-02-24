with open("apps/api/internal/svc/identity/identity_service.go", "r") as f:
    content = f.read()

import re

# Add CreateSession to the interface
content = re.sub(
    r'(LoginWithPassword\(ctx context\.Context, email string, password string, ipAddress \*string, userAgent \*string\) \(\*models\.UserSession, error\))',
    r'\1\n\tCreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error)',
    content
)

# Add implementation
impl = """func (s *identityService) CreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error) {
	// Session management
	if s.config.Identity.Sessions.MaxConcurrent > 0 {
		sessions, err := s.sessionStore.ListByUser(ctx, userID)
		if err == nil && len(sessions) >= s.config.Identity.Sessions.MaxConcurrent {
			return nil, ErrMaxConcurrentSessions
			// Alternatively, could delete oldest session
		}
	}

	// Create Session
	session := &models.UserSession{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     uuid.New().String(), // Placeholder for token generation, possibly JWT
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(time.Duration(s.config.Identity.Sessions.TTLMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.sessionStore.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}"""

content = content + "\n\n" + impl

with open("apps/api/internal/svc/identity/identity_service.go", "w") as f:
    f.write(content)
