import re

with open("apps/api/internal/models/identity.go", "r") as f:
    content = f.read()

content = re.sub(
    r'(IsLocked\s+bool\s+`json:"isLocked"`\n)',
    r'\1\tTotpSecret        *string    `json:"-"`\n\tTotpEnabled       bool       `json:"totpEnabled"`\n',
    content
)

with open("apps/api/internal/models/identity.go", "w") as f:
    f.write(content)
