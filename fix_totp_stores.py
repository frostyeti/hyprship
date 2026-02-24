import os
import re

DRIVERS = ["sqlite", "pg", "mysql", "mssql"]

for drv in DRIVERS:
    path = f"apps/api/internal/stores/{drv}/totp_store.go"
    with open(path, "r") as f:
        content = f.read()
    
    content = content.replace('\t"github.com/frostyeti/hyprship/apps/api/internal/stores"\n', "")
    content = content.replace("func NewUserTotpStore(db *sql.DB) stores.UserTotpStore {", "func NewUserTotpStore(db *sql.DB) *userTotpStore {")
    
    with open(path, "w") as f:
        f.write(content)

