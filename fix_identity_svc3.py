with open("apps/api/internal/svc/identity/identity_service.go", "r") as f:
    lines = f.readlines()

out = []
found = False
for line in lines:
    if line.strip() == "CreateSession(ctx context.Context, userID uuid.UUID, ipAddress *string, userAgent *string) (*models.UserSession, error)":
        if found:
            continue
        found = True
    out.append(line)

with open("apps/api/internal/svc/identity/identity_service.go", "w") as f:
    f.writelines(out)
