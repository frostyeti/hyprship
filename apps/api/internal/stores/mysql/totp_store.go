package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type userTotpStore struct {
	db *sql.DB
}

func NewUserTotpStore(db *sql.DB) *userTotpStore {
	return &userTotpStore{db: db}
}

func (s *userTotpStore) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserTotp, error) {
	query := `SELECT CAST(id AS CHAR(36)), CAST(user_id AS CHAR(36)), secret, is_verified, created_at, updated_at FROM user_totps WHERE user_id = ?`
	var t models.UserTotp
	var idStr, userIDStr string
	var updatedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), userID.String()).Scan(
		&idStr, &userIDStr, &t.Secret, &t.IsVerified, &t.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	uid, _ := uuid.Parse(userIDStr)
	pid, _ := uuid.Parse(idStr)
	t.ID = pid
	t.UserID = uid
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}

	return &t, nil
}

func (s *userTotpStore) Upsert(ctx context.Context, totp *models.UserTotp) error {
	query := `
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
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), totp.ID.String(), totp.UserID.String(), totp.Secret, totp.IsVerified, totp.CreatedAt, totp.UpdatedAt)
	return err
}

func (s *userTotpStore) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_totps WHERE user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), userID.String())
	return err
}
