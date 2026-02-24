import os

DRIVERS = ["sqlite", "pg", "mysql", "mssql"]

TOTP_STORE_TEMPLATE = """package {driver}

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
)

type userTotpStore struct {{
	db *sql.DB
}}

func NewUserTotpStore(db *sql.DB) stores.UserTotpStore {{
	return &userTotpStore{{db: db}}
}}

func (s *userTotpStore) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserTotp, error) {{
	query := `SELECT {id_cast}, {user_id_cast}, secret, is_verified, created_at, updated_at FROM user_totps WHERE user_id = ?`
	var t models.UserTotp
	var idStr, userIDStr string
	var updatedAt {null_time_type}

	err := s.db.QueryRowContext(ctx, core.Rebind("{driver}", query), userID.String()).Scan(
		&idStr, &userIDStr, &t.Secret, &t.IsVerified, &t.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {{
		return nil, nil
	}} else if err != nil {{
		return nil, err
	}}

	uid, _ := uuid.Parse(userIDStr)
	pid, _ := uuid.Parse(idStr)
	t.ID = pid
	t.UserID = uid
	{assign_updated_at}
	
	return &t, nil
}}

func (s *userTotpStore) Upsert(ctx context.Context, totp *models.UserTotp) error {{
	{upsert_query}
}}

func (s *userTotpStore) Delete(ctx context.Context, userID uuid.UUID) error {{
	query := `DELETE FROM user_totps WHERE user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("{driver}", query), userID.String())
	return err
}}
"""

UPSERT_SQLITE = """query := `
		INSERT INTO user_totps (id, user_id, secret, is_verified, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET 
			secret=excluded.secret, 
			is_verified=excluded.is_verified, 
			updated_at=excluded.updated_at
	`
	now := time.Now()
	if totp.CreatedAt.IsZero() {
		totp.CreatedAt = now
	}
	if totp.ID == uuid.Nil {
		totp.ID = uuid.New()
	}
	totp.UpdatedAt = &now
	_, err := s.db.ExecContext(ctx, core.Rebind("{driver}", query), totp.ID.String(), totp.UserID.String(), totp.Secret, totp.IsVerified, totp.CreatedAt, totp.UpdatedAt)
	return err"""

UPSERT_PG = """query := `
		INSERT INTO user_totps (id, user_id, secret, is_verified, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (user_id) DO UPDATE SET 
			secret=EXCLUDED.secret, 
			is_verified=EXCLUDED.is_verified, 
			updated_at=EXCLUDED.updated_at
	`
	now := time.Now()
	if totp.CreatedAt.IsZero() {
		totp.CreatedAt = now
	}
	if totp.ID == uuid.Nil {
		totp.ID = uuid.New()
	}
	totp.UpdatedAt = &now
	_, err := s.db.ExecContext(ctx, core.Rebind("{driver}", query), totp.ID.String(), totp.UserID.String(), totp.Secret, totp.IsVerified, totp.CreatedAt, totp.UpdatedAt)
	return err"""

UPSERT_MYSQL = """query := `
		INSERT INTO user_totps (id, user_id, secret, is_verified, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			secret=VALUES(secret), 
			is_verified=VALUES(is_verified), 
			updated_at=VALUES(updated_at)
	`
	now := time.Now()
	if totp.CreatedAt.IsZero() {
		totp.CreatedAt = now
	}
	if totp.ID == uuid.Nil {
		totp.ID = uuid.New()
	}
	totp.UpdatedAt = &now
	_, err := s.db.ExecContext(ctx, core.Rebind("{driver}", query), totp.ID.String(), totp.UserID.String(), totp.Secret, totp.IsVerified, totp.CreatedAt, totp.UpdatedAt)
	return err"""

UPSERT_MSSQL = """// MSSQL requires MERGE statement or checking first.
	existing, err := s.GetByUserID(ctx, totp.UserID)
	if err != nil { return err }
	
	now := time.Now()
	if totp.CreatedAt.IsZero() {
		totp.CreatedAt = now
	}
	if totp.ID == uuid.Nil {
		totp.ID = uuid.New()
	}
	totp.UpdatedAt = &now
	
	if existing != nil {
		query := `UPDATE user_totps SET secret=?, is_verified=?, updated_at=? WHERE user_id=?`
		_, err = s.db.ExecContext(ctx, core.Rebind("{driver}", query), totp.Secret, totp.IsVerified, totp.UpdatedAt, totp.UserID.String())
		return err
	} else {
		query := `INSERT INTO user_totps (id, user_id, secret, is_verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = s.db.ExecContext(ctx, core.Rebind("{driver}", query), totp.ID.String(), totp.UserID.String(), totp.Secret, totp.IsVerified, totp.CreatedAt, totp.UpdatedAt)
		return err
	}"""

for drv in DRIVERS:
    if drv == "sqlite":
        upsert = UPSERT_SQLITE.replace("{driver}", drv)
        null_time_type = "sql.NullTime"
        assign_time = """if updatedAt.Valid { t.UpdatedAt = &updatedAt.Time }"""
    elif drv == "pg":
        upsert = UPSERT_PG.replace("{driver}", drv)
        null_time_type = "sql.NullTime"
        assign_time = """if updatedAt.Valid { t.UpdatedAt = &updatedAt.Time }"""
    elif drv == "mysql":
        upsert = UPSERT_MYSQL.replace("{driver}", drv)
        null_time_type = "sql.NullTime"
        assign_time = """if updatedAt.Valid { t.UpdatedAt = &updatedAt.Time }"""
    elif drv == "mssql":
        upsert = UPSERT_MSSQL.replace("{driver}", drv)
        null_time_type = "sql.NullTime"
        assign_time = """if updatedAt.Valid { t.UpdatedAt = &updatedAt.Time }"""

    content = TOTP_STORE_TEMPLATE.format(
        driver=drv,
        upsert_query=upsert,
        null_time_type=null_time_type,
        assign_updated_at=assign_time,
        id_cast="id" if drv in ["pg", "sqlite"] else "CAST(id AS CHAR(36))",
        user_id_cast="user_id" if drv in ["pg", "sqlite"] else "CAST(user_id AS CHAR(36))"
    )

    with open(f"apps/api/internal/stores/{drv}/totp_store.go", "w") as f:
        f.write(content)
