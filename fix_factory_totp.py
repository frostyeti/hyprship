with open("apps/api/internal/stores/factory.go", "r") as f:
    content = f.read()

import re

# Add to StoreFactory
content = re.sub(
    r'(UserPasskeyStore\s+UserPasskeyStore)',
    r'\1\n\tUserTotpStore            UserTotpStore',
    content
)

# Add to implementations
content = re.sub(
    r'(UserPasskeyStore:\s+sqlite\.NewUserPasskeyStore\(db\),)',
    r'\1\n\t\t\tUserTotpStore:            sqlite.NewUserTotpStore(db),',
    content
)

content = re.sub(
    r'(UserPasskeyStore:\s+pg\.NewUserPasskeyStore\(db\),)',
    r'\1\n\t\t\tUserTotpStore:            pg.NewUserTotpStore(db),',
    content
)

content = re.sub(
    r'(UserPasskeyStore:\s+mysql\.NewUserPasskeyStore\(db\),)',
    r'\1\n\t\t\tUserTotpStore:            mysql.NewUserTotpStore(db),',
    content
)

content = re.sub(
    r'(UserPasskeyStore:\s+mssql\.NewUserPasskeyStore\(db\),)',
    r'\1\n\t\t\tUserTotpStore:            mssql.NewUserTotpStore(db),',
    content
)

with open("apps/api/internal/stores/factory.go", "w") as f:
    f.write(content)
