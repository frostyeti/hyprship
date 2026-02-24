with open("apps/api/internal/svc/identity/passkey_service.go", "r") as f:
    content = f.read()

import re

# Update DeletePasskey signature
content = re.sub(
    r'(DeletePasskey\(ctx context\.Context, passkeyID uuid\.UUID\) error)',
    r'DeletePasskey(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error',
    content,
    count=1
)

new_delete = """func (s *passkeyService) DeletePasskey(ctx context.Context, userID uuid.UUID, passkeyID uuid.UUID) error {
	pk, err := s.store.Get(ctx, passkeyID)
	if err != nil {
		return err
	}
	if pk.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.store.Delete(ctx, passkeyID)
}"""

content = re.sub(
    r'func \(s \*passkeyService\) DeletePasskey\(ctx context\.Context, passkeyID uuid\.UUID\) error \{\n\treturn s\.store\.Delete\(ctx, passkeyID\)\n\}',
    new_delete,
    content
)

with open("apps/api/internal/svc/identity/passkey_service.go", "w") as f:
    f.write(content)

with open("apps/api/internal/routes/v1/users.go", "r") as f:
    content = f.read()

# Fix the handler to actually pass userID
content = re.sub(
    r'if err := h\.passkeySvc\.DeletePasskey\(c\.Request\.Context\(\), passkeyID\); err != nil \{',
    r'if err := h.passkeySvc.DeletePasskey(c.Request.Context(), userID.(uuid.UUID), passkeyID); err != nil {',
    content
)

with open("apps/api/internal/routes/v1/users.go", "w") as f:
    f.write(content)
