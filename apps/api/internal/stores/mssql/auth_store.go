package mssql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserPasswordAuthStore struct {
	db *sql.DB
}

func NewUserPasswordAuthStore(db *sql.DB) *UserPasswordAuthStore {
	return &UserPasswordAuthStore{db: db}
}

func (s *UserPasswordAuthStore) Get(ctx context.Context, userID uuid.UUID) (*models.UserPasswordAuth, error) {
	var a models.UserPasswordAuth
	var passwordExpiresAt, lastAttemptedAt, otpExpiresAt, updatedAt sql.NullTime
	var userIDStr string
	query := `SELECT CAST(user_id AS CHAR(36)), password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at FROM user_password_auth WHERE user_id = ?`
	
	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), userID.String()).Scan(
		&userIDStr, &a.PasswordDigest, &passwordExpiresAt, &lastAttemptedAt, &a.AttemptCount,
		&a.OtpDigest, &otpExpiresAt, &a.OtpLinkToken, &a.IsLocked, &a.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}

	a.UserID, _ = uuid.Parse(userIDStr)
	if passwordExpiresAt.Valid {
		a.PasswordExpiresAt = &passwordExpiresAt.Time
	}
	if lastAttemptedAt.Valid {
		a.LastAttemptedAt = &lastAttemptedAt.Time
	}
	if otpExpiresAt.Valid {
		a.OtpExpiresAt = &otpExpiresAt.Time
	}
	if updatedAt.Valid {
		a.UpdatedAt = &updatedAt.Time
	}
	return &a, nil
}

func (s *UserPasswordAuthStore) Create(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `INSERT INTO user_password_auth (user_id, password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), a.UserID.String(), a.PasswordDigest, a.PasswordExpiresAt, a.LastAttemptedAt, a.AttemptCount, a.OtpDigest, a.OtpExpiresAt, a.OtpLinkToken, a.IsLocked, a.CreatedAt, a.UpdatedAt)
	return err
}

func (s *UserPasswordAuthStore) Update(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `UPDATE user_password_auth SET password_digest = ?, password_expires_at = ?, last_attempted_at = ?, attempt_count = ?, otp_digest = ?, otp_expires_at = ?, otp_link_token = ?, is_locked = ?, updated_at = ? WHERE user_id = ?`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), a.PasswordDigest, a.PasswordExpiresAt, a.LastAttemptedAt, a.AttemptCount, a.OtpDigest, a.OtpExpiresAt, a.OtpLinkToken, a.IsLocked, now, a.UserID.String())
	return err
}

func (s *UserPasswordAuthStore) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_password_auth WHERE user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), userID.String())
	return err
}
