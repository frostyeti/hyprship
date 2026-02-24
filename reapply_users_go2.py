import re

with open("apps/api/internal/routes/v1/users.go", "r") as f:
    content = f.read()

# Add PasskeySvc and APIKeySvc to UserHandler struct
content = re.sub(
    r'type UserHandler struct \{\n\tsvc          identity\.IdentityService\n\tstore        stores\.UserStore\n\tsessionStore stores\.UserSessionStore\n\}',
    r'type UserHandler struct {\n\tsvc          identity.IdentityService\n\tstore        stores.UserStore\n\tsessionStore stores.UserSessionStore\n\tapiKeySvc    identity.APIKeyService\n\tpasskeySvc   identity.PasskeyService\n}',
    content
)

# Fix RegisterUserRoutes
content = re.sub(
    r'func RegisterUserRoutes\(r \*gin\.RouterGroup, svc identity\.IdentityService, store stores\.UserStore, sessionStore stores\.UserSessionStore\) \{',
    r'func RegisterUserRoutes(r *gin.RouterGroup, svc identity.IdentityService, store stores.UserStore, sessionStore stores.UserSessionStore, apiKeySvc identity.APIKeyService, passkeySvc identity.PasskeyService) {',
    content
)

content = re.sub(
    r'h := &UserHandler\{svc: svc, store: store, sessionStore: sessionStore\}',
    r'h := &UserHandler{svc: svc, store: store, sessionStore: sessionStore, apiKeySvc: apiKeySvc, passkeySvc: passkeySvc}',
    content
)

with open("apps/api/internal/routes/v1/users.go", "w") as f:
    f.write(content)
