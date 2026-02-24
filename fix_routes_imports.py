import re

with open("apps/api/internal/routes/v1/auth.go", "r") as f:
    content = f.read()

content = re.sub(
    r'"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"',
    r'"github.com/frostyeti/hyprship/apps/api/internal/stores"\n\t"github.com/frostyeti/hyprship/apps/api/internal/svc/identity"',
    content
)
with open("apps/api/internal/routes/v1/auth.go", "w") as f:
    f.write(content)

with open("apps/api/internal/routes/v1/routes.go", "r") as f:
    content = f.read()

content = content.replace("RegisterAuthRoutes(authGroup, identitySvc, passkeySvc)", "RegisterAuthRoutes(authGroup, identitySvc, userStore, passkeySvc)")

with open("apps/api/internal/routes/v1/routes.go", "w") as f:
    f.write(content)
